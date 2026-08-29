package core

import (
	"context"
	"io"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// RuntimeRepo abstracts Kubernetes runtime operations. Every method takes a
// cluster name so the implementation can route through the right tunnel.
type RuntimeRepo interface {
	PodLogs(ctx context.Context, cluster, namespace, name string, opts PodLogOptions) (io.ReadCloser, error)
	// Exec blocks until the session completes.
	Exec(ctx context.Context, cluster, namespace, name string, opts *ExecOptions) error
	// UpdateScale writes through the /scale subresource and returns the new value.
	UpdateScale(ctx context.Context, cluster string, gvr schema.GroupVersionResource, namespace, name string, replicas int32) (int32, error)
	// Restart patches the pod template annotation to trigger a rolling restart.
	Restart(ctx context.Context, cluster string, gvr schema.GroupVersionResource, namespace, name string) error
	// PortForward copies bidirectionally until ctx is canceled or the connection closes.
	PortForward(ctx context.Context, cluster, namespace, name string, opts PortForwardOptions) error
	// SubResourceAction invokes a PUT/POST action on a named subresource
	// (e.g. KubeVirt VM start/stop/restart) and returns the response body, if any.
	SubResourceAction(ctx context.Context, cluster string, gvr schema.GroupVersionResource,
		namespace, name, subresource, method string, body []byte) (map[string]any, error)
	// VNC copies bidirectionally until ctx is canceled or the connection closes.
	VNC(ctx context.Context, cluster, namespace, name string, opts VNCOptions) error
}

// HelmRepo runs directly on the server via the Helm Go SDK, not through the tunnel.
type HelmRepo interface {
	// ShowChart returns values.yaml and README.md from a remote HTTP or OCI repository.
	ShowChart(ctx context.Context, repoURL, chartName, version string) (values, readme []byte, err error)
}

// PodLogOptions mirrors the corev1.PodLogOptions fields exposed through the proto.
type PodLogOptions struct {
	Container    string
	Follow       bool
	TailLines    *int64
	SinceSeconds *int64
	SinceTime    *time.Time
	Previous     bool
	Timestamps   bool
	LimitBytes   *int64
}

type ExecOptions struct {
	Container string
	Command   []string
	TTY       bool
	Stdin     io.Reader
	Stdout    io.Writer
	Stderr    io.Writer
	SizeQueue TerminalSizer
}

// StartExecParams keeps RuntimeUseCase.StartExec off a long parameter list.
type StartExecParams struct {
	Cluster   string
	Namespace string
	Name      string
	Container string
	Command   []string
	TTY       bool
	Rows      uint16
	Cols      uint16
}

type PortForwardOptions struct {
	Port   int32
	Stdin  io.Reader
	Stdout io.Writer
}

type VNCOptions struct {
	Stdin  io.Reader
	Stdout io.Writer
}

// RuntimeUseCase provides runtime operations with session management for
// exec, port-forward, and VNC.
type RuntimeUseCase struct {
	discovery DiscoveryClient
	runtime   RuntimeRepo
	helm      HelmRepo
	sessions  *SessionStore
}

// NewRuntimeUseCase takes the SessionStore as a dependency rather than creating
// one so callers can substitute it for testing or monitoring.
func NewRuntimeUseCase(discovery DiscoveryClient, runtime RuntimeRepo, helm HelmRepo, sessions *SessionStore) *RuntimeUseCase {
	return &RuntimeUseCase{
		discovery: discovery,
		runtime:   runtime,
		helm:      helm,
		sessions:  sessions,
	}
}

// sessionOwner returns the subject a new session belongs to. Sessions are
// addressed by an identifier the client replays later, so they are bound to
// their creator at birth.
func sessionOwner(ctx context.Context) (string, error) {
	user, ok := UserInfoFromContext(ctx)
	if !ok || user.Subject == "" {
		return "", &DomainError{
			Code:    ErrorCodeUnauthenticated,
			Message: "user info not found in context",
		}
	}
	return user.Subject, nil
}

func ownedBy(ctx context.Context, owner string) bool {
	user, ok := UserInfoFromContext(ctx)
	return ok && user.Subject != "" && user.Subject == owner
}

// execSession looks up a session the caller may use. Somebody else's session is
// reported as missing rather than forbidden, so a leaked identifier cannot be
// used to probe for other users' sessions.
func (uc *RuntimeUseCase) execSession(ctx context.Context, sessionID string) (*ExecSession, error) {
	sess, ok := uc.sessions.GetExec(sessionID)
	if !ok || !ownedBy(ctx, sess.Owner) {
		return nil, &ErrSessionNotFound{Resource: "exec-session", ID: sessionID}
	}
	return sess, nil
}

func (uc *RuntimeUseCase) portForwardSession(ctx context.Context, sessionID string) (*PortForwardSession, error) {
	sess, ok := uc.sessions.GetPortForward(sessionID)
	if !ok || !ownedBy(ctx, sess.Owner) {
		return nil, &ErrSessionNotFound{Resource: "portforward-session", ID: sessionID}
	}
	return sess, nil
}

func (uc *RuntimeUseCase) vncSession(ctx context.Context, sessionID string) (*VNCSession, error) {
	sess, ok := uc.sessions.GetVNC(sessionID)
	if !ok || !ownedBy(ctx, sess.Owner) {
		return nil, &ErrSessionNotFound{Resource: "vnc-session", ID: sessionID}
	}
	return sess, nil
}

func (uc *RuntimeUseCase) StartPodLogs(ctx context.Context, cluster, namespace, name string, opts PodLogOptions) (io.ReadCloser, error) {
	if name == "" {
		return nil, &ErrInvalidInput{Field: fieldName, Message: msgPodNameRequired}
	}
	return uc.runtime.PodLogs(ctx, cluster, namespace, name, opts)
}

// StartExec runs the exec in a background goroutine and returns the session
// plus stdout and stderr readers for the caller to stream from.
func (uc *RuntimeUseCase) StartExec(ctx context.Context, params *StartExecParams) (session *ExecSession, stdoutReader, stderrReader io.ReadCloser, err error) {
	if params.Name == "" {
		return nil, nil, nil, &ErrInvalidInput{Field: fieldName, Message: msgPodNameRequired}
	}
	if len(params.Command) == 0 {
		return nil, nil, nil, &ErrInvalidInput{Field: "command", Message: "command is required"}
	}

	owner, err := sessionOwner(ctx)
	if err != nil {
		return nil, nil, nil, err
	}

	stdinR, stdinW := io.Pipe()
	stdoutR, stdoutW := io.Pipe()
	stderrR, stderrW := io.Pipe()
	sizeQueue := NewTerminalSizeQueue()

	if params.Rows > 0 && params.Cols > 0 {
		sizeQueue.Set(params.Cols, params.Rows)
	}

	ctx, cancel := context.WithCancel(ctx)
	done := make(chan struct{})

	session = &ExecSession{
		ID:        uuid.New().String(),
		Owner:     owner,
		Stdin:     stdinW,
		SizeQueue: sizeQueue,
		Cancel:    cancel,
		Done:      done,
	}

	// Register before launching the goroutine, so a full session store costs
	// nothing.
	if err := uc.sessions.PutExec(session); err != nil {
		cancel()
		stdinW.Close()
		stdinR.Close()
		stdoutW.Close()
		stderrW.Close()
		sizeQueue.Close()
		return nil, nil, nil, err
	}

	go func() {
		defer close(done)
		defer stdinR.Close()
		defer stdoutW.Close()
		defer stderrW.Close()
		defer sizeQueue.Close()

		var stderr io.Writer
		if !params.TTY {
			stderr = stderrW
		}

		session.Err = uc.runtime.Exec(ctx, params.Cluster, params.Namespace, params.Name, &ExecOptions{
			Container: params.Container,
			Command:   params.Command,
			TTY:       params.TTY,
			Stdin:     stdinR,
			Stdout:    stdoutW,
			Stderr:    stderr,
			SizeQueue: sizeQueue,
		})
	}()

	return session, stdoutR, stderrR, nil
}

// WaitExec blocks until the session finishes and returns the error it ended
// with, or ctx.Err() if the caller gives up first.
//
// Reading sess.Err is only safe once Done is closed: the session goroutine
// assigns it as its last act, and the deferred close of Done is what publishes
// that write to observers — which is also why Done is closed rather than sent to.
func (uc *RuntimeUseCase) WaitExec(ctx context.Context, sess *ExecSession) error {
	select {
	case <-sess.Done:
		return sess.Err
	case <-ctx.Done():
		return ctx.Err()
	}
}

// WaitPortForward: see WaitExec.
func (uc *RuntimeUseCase) WaitPortForward(ctx context.Context, sess *PortForwardSession) error {
	select {
	case <-sess.Done:
		return sess.Err
	case <-ctx.Done():
		return ctx.Err()
	}
}

// WaitVNC: see WaitExec.
func (uc *RuntimeUseCase) WaitVNC(ctx context.Context, sess *VNCSession) error {
	select {
	case <-sess.Done:
		return sess.Err
	case <-ctx.Done():
		return ctx.Err()
	}
}

// WriteExec writes stdin data in a background goroutine so the caller's context
// can cancel a blocking pipe write during graceful shutdown.
func (uc *RuntimeUseCase) WriteExec(ctx context.Context, sessionID string, data []byte) error {
	sess, err := uc.execSession(ctx, sessionID)
	if err != nil {
		return err
	}

	// Once the session goroutine has exited the pipe reader is closed and Write
	// fails with an opaque error; check Done first for a clearer one.
	select {
	case <-sess.Done:
		return &ErrSessionNotFound{Resource: "exec-session", ID: sessionID}
	default:
	}

	errCh := make(chan error, 1)
	go func() {
		_, err := sess.Stdin.Write(data)
		errCh <- err
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}

func (uc *RuntimeUseCase) ResizeExec(ctx context.Context, sessionID string, rows, cols uint16) error {
	sess, err := uc.execSession(ctx, sessionID)
	if err != nil {
		return err
	}
	sess.SizeQueue.Set(cols, rows)
	return nil
}

// CleanupExec removes the session with RemoveExec rather than Get+Delete, to
// atomically claim ownership and avoid a double-close race with ReapStaleSessions.
func (uc *RuntimeUseCase) CleanupExec(ctx context.Context, sessionID string) {
	if _, err := uc.execSession(ctx, sessionID); err != nil {
		return
	}
	sess := uc.sessions.RemoveExec(sessionID)
	if sess == nil {
		return
	}
	sess.Cancel()
	sess.Stdin.Close()
}

// StartPortForward forwards in a background goroutine and returns the session
// plus a reader for data coming from the pod.
func (uc *RuntimeUseCase) StartPortForward(ctx context.Context, cluster, namespace, name string, port int32) (*PortForwardSession, io.ReadCloser, error) {
	if name == "" {
		return nil, nil, &ErrInvalidInput{Field: fieldName, Message: msgPodNameRequired}
	}
	if port <= 0 || port > 65535 {
		return nil, nil, &ErrInvalidInput{Field: "port", Message: "must be between 1 and 65535"}
	}

	owner, err := sessionOwner(ctx)
	if err != nil {
		return nil, nil, err
	}

	dataInR, dataInW := io.Pipe()
	dataOutR, dataOutW := io.Pipe()

	ctx, cancel := context.WithCancel(ctx)
	done := make(chan struct{})

	sess := &PortForwardSession{
		ID:     uuid.New().String(),
		Owner:  owner,
		Writer: dataInW,
		Cancel: cancel,
		Done:   done,
	}

	// Register before launching the goroutine, so a full session store costs
	// nothing.
	if err := uc.sessions.PutPortForward(sess); err != nil {
		cancel()
		dataInW.Close()
		dataInR.Close()
		dataOutW.Close()
		return nil, nil, err
	}

	// Without this the adapter's "client → pod" copy loop stays blocked reading
	// a pipe that only the goroutine below would close — and that goroutine is
	// itself waiting for the copy loop to finish.
	closeStdinOnCancel(ctx, dataInR)

	go func() {
		defer close(done)
		defer dataInR.Close()
		defer dataOutW.Close()
		sess.Err = uc.runtime.PortForward(ctx, cluster, namespace, name, PortForwardOptions{
			Port:   port,
			Stdin:  dataInR,
			Stdout: dataOutW,
		})
	}()

	return sess, dataOutR, nil
}

// closeStdinOnCancel closes r once ctx is done, unblocking any read in flight.
// Pipe reads are not interruptible, so a session whose client sends nothing
// would otherwise keep its adapter goroutine alive after cancellation — and
// with it the session entry, which the reaper can only collect once the
// session's Done channel fires.
func closeStdinOnCancel(ctx context.Context, r io.Closer) {
	go func() {
		<-ctx.Done()
		if err := r.Close(); err != nil {
			slog.Debug("failed to close session stdin after cancellation", "error", err)
		}
	}()
}

// WritePortForward writes in a background goroutine so the caller's context can
// cancel a blocking pipe write during graceful shutdown.
func (uc *RuntimeUseCase) WritePortForward(ctx context.Context, sessionID string, data []byte) error {
	sess, err := uc.portForwardSession(ctx, sessionID)
	if err != nil {
		return err
	}

	// See WriteExec.
	select {
	case <-sess.Done:
		return &ErrSessionNotFound{Resource: "portforward-session", ID: sessionID}
	default:
	}

	errCh := make(chan error, 1)
	go func() {
		_, err := sess.Writer.Write(data)
		errCh <- err
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}

// CleanupPortForward: see CleanupExec.
func (uc *RuntimeUseCase) CleanupPortForward(ctx context.Context, sessionID string) {
	if _, err := uc.portForwardSession(ctx, sessionID); err != nil {
		return
	}
	sess := uc.sessions.RemovePortForward(sessionID)
	if sess == nil {
		return
	}
	sess.Cancel()
	sess.Writer.Close()
}

// StartVNC connects in a background goroutine and returns the session plus a
// reader for data coming from the VMI.
func (uc *RuntimeUseCase) StartVNC(ctx context.Context, cluster, namespace, name string) (*VNCSession, io.ReadCloser, error) {
	if name == "" {
		return nil, nil, &ErrInvalidInput{Field: fieldName, Message: "VMI name is required"}
	}

	owner, err := sessionOwner(ctx)
	if err != nil {
		return nil, nil, err
	}

	dataInR, dataInW := io.Pipe()
	dataOutR, dataOutW := io.Pipe()

	ctx, cancel := context.WithCancel(ctx)
	done := make(chan struct{})

	sess := &VNCSession{
		ID:     uuid.New().String(),
		Owner:  owner,
		Writer: dataInW,
		Cancel: cancel,
		Done:   done,
	}

	if err := uc.sessions.PutVNC(sess); err != nil {
		cancel()
		dataInW.Close()
		dataInR.Close()
		dataOutR.Close()
		dataOutW.Close()
		return nil, nil, err
	}

	// See StartPortForward.
	closeStdinOnCancel(ctx, dataInR)

	go func() {
		defer close(done)
		defer dataInR.Close()
		defer dataOutW.Close()
		sess.Err = uc.runtime.VNC(ctx, cluster, namespace, name, VNCOptions{
			Stdin:  dataInR,
			Stdout: dataOutW,
		})
	}()

	return sess, dataOutR, nil
}

func (uc *RuntimeUseCase) WriteVNC(ctx context.Context, sessionID string, data []byte) error {
	sess, err := uc.vncSession(ctx, sessionID)
	if err != nil {
		return err
	}

	select {
	case <-sess.Done:
		if sess.Err != nil {
			return sess.Err
		}
		return &ErrSessionNotFound{Resource: "vnc-session", ID: sessionID}
	default:
	}

	errCh := make(chan error, 1)
	go func() {
		_, err := sess.Writer.Write(data)
		errCh <- err
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}

// CleanupVNC: see CleanupExec.
func (uc *RuntimeUseCase) CleanupVNC(ctx context.Context, sessionID string) {
	if _, err := uc.vncSession(ctx, sessionID); err != nil {
		return
	}
	sess := uc.sessions.RemoveVNC(sessionID)
	if sess == nil {
		return
	}
	sess.Cancel()
	sess.Writer.Close()
}

func (uc *RuntimeUseCase) Scale(ctx context.Context, id *ResourceIdentifier, replicas int32) (int32, error) {
	if id.Name == "" {
		return 0, &ErrInvalidInput{Field: fieldName, Message: msgResourceNameRequired}
	}
	if replicas < 0 {
		return 0, &ErrInvalidInput{Field: "replicas", Message: "must be non-negative"}
	}
	gvr, err := id.lookupGVR(ctx, uc.discovery)
	if err != nil {
		return 0, err
	}
	return uc.runtime.UpdateScale(ctx, id.Cluster, gvr, id.Namespace, id.Name, replicas)
}

// StartSessionReaper periodically removes sessions that finished but were never
// cleaned up. It blocks until ctx is canceled.
func (uc *RuntimeUseCase) StartSessionReaper(ctx context.Context, interval time.Duration) {
	log := slog.Default().With("component", "session-reaper")
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if n := uc.sessions.ReapStaleSessions(); n > 0 {
				log.Info("reaped stale sessions", "count", n)
			}
		}
	}
}

func (uc *RuntimeUseCase) Restart(ctx context.Context, id *ResourceIdentifier) error {
	if id.Name == "" {
		return &ErrInvalidInput{Field: fieldName, Message: msgResourceNameRequired}
	}
	gvr, err := id.lookupGVR(ctx, uc.discovery)
	if err != nil {
		return err
	}
	return uc.runtime.Restart(ctx, id.Cluster, gvr, id.Namespace, id.Name)
}

// DELETE and PATCH have dedicated RPCs.
var allowedSubResourceMethods = map[string]bool{"PUT": true, "POST": true}

// SubResourceAction forwards a PUT/POST subresource action to the Kubernetes API
// via impersonation, covering cases such as KubeVirt VM start/stop/restart/migrate.
func (uc *RuntimeUseCase) SubResourceAction(ctx context.Context, id *ResourceIdentifier, method string, body []byte) (map[string]any, error) {
	if id.Name == "" {
		return nil, &ErrInvalidInput{Field: fieldName, Message: msgResourceNameRequired}
	}
	if id.SubResource == "" {
		return nil, &ErrInvalidInput{Field: "subresource", Message: "subresource is required"}
	}
	if !allowedSubResourceMethods[method] {
		return nil, &ErrInvalidInput{Field: "method", Message: "must be PUT or POST"}
	}

	gvr, err := id.lookupGVR(ctx, uc.discovery)
	if err != nil {
		return nil, err
	}

	return uc.runtime.SubResourceAction(ctx, id.Cluster, gvr, id.Namespace, id.Name, id.SubResource, method, body)
}

func (uc *RuntimeUseCase) ShowChart(ctx context.Context, repoURL, chartName, version string) (values, readme []byte, err error) {
	if repoURL == "" {
		return nil, nil, &ErrInvalidInput{Field: "repo_url", Message: "repository URL is required"}
	}
	if chartName == "" {
		return nil, nil, &ErrInvalidInput{Field: "chart_name", Message: "chart name is required"}
	}
	return uc.helm.ShowChart(ctx, repoURL, chartName, version)
}
