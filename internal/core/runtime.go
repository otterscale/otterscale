package core

import (
	"context"
	"io"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// ---------------------------------------------------------------------------
// Interfaces
// ---------------------------------------------------------------------------

// RuntimeRepo abstracts Kubernetes runtime operations (logs, exec,
// scale, restart, port-forward). All methods accept a cluster name so
// that the underlying implementation can route requests through the
// correct tunnel.
type RuntimeRepo interface {
	// PodLogs opens a streaming reader for container log output.
	PodLogs(ctx context.Context, cluster, namespace, name string, opts PodLogOptions) (io.ReadCloser, error)
	// Exec starts an exec session and blocks until it completes.
	Exec(ctx context.Context, cluster, namespace, name string, opts *ExecOptions) error
	// UpdateScale sets the desired replica count via the /scale subresource
	// and returns the updated value.
	UpdateScale(ctx context.Context, cluster string, gvr schema.GroupVersionResource, namespace, name string, replicas int32) (int32, error)
	// Restart triggers a rolling restart by patching the pod template annotation.
	Restart(ctx context.Context, cluster string, gvr schema.GroupVersionResource, namespace, name string) error
	// PortForward opens a port-forward session and copies data
	// bidirectionally until the context is canceled or the
	// connection closes.
	PortForward(ctx context.Context, cluster, namespace, name string, opts PortForwardOptions) error
	// SubResourceAction invokes a PUT or POST action on a named
	// subresource (e.g. KubeVirt VM start/stop/restart). Returns
	// the response body, if any.
	SubResourceAction(ctx context.Context, cluster string, gvr schema.GroupVersionResource,
		namespace, name, subresource, method string, body []byte) (map[string]any, error)
	// VNC opens a VNC WebSocket session to a KubeVirt VMI and copies
	// data bidirectionally until the context is canceled or the
	// connection closes.
	VNC(ctx context.Context, cluster, namespace, name string, opts VNCOptions) error
}

// HelmRepo abstracts server-side Helm chart repository operations.
// Unlike RuntimeRepo (which routes through the tunnel), HelmRepo
// executes directly on the server using the Helm Go SDK.
type HelmRepo interface {
	// ShowChart retrieves the default values.yaml and README.md content
	// from a chart in a remote Helm repository (HTTP or OCI).
	ShowChart(ctx context.Context, repoURL, chartName, version string) (values, readme []byte, err error)
}

// ---------------------------------------------------------------------------
// Options types
// ---------------------------------------------------------------------------

// PodLogOptions mirrors the fields of corev1.PodLogOptions that are
// exposed through the RuntimeService proto.
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

// ExecOptions holds parameters for an interactive exec session.
type ExecOptions struct {
	Container string
	Command   []string
	TTY       bool
	Stdin     io.Reader
	Stdout    io.Writer
	Stderr    io.Writer
	SizeQueue TerminalSizer
}

// StartExecParams collects the parameters for starting an interactive
// exec session. This avoids a long parameter list on
// RuntimeUseCase.StartExec.
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

// PortForwardOptions holds parameters for a port-forward session.
type PortForwardOptions struct {
	Port   int32
	Stdin  io.Reader
	Stdout io.Writer
}

// VNCOptions holds parameters for a VNC session.
type VNCOptions struct {
	Stdin  io.Reader
	Stdout io.Writer
}

// ---------------------------------------------------------------------------
// Use case
// ---------------------------------------------------------------------------

// RuntimeUseCase provides application-level runtime operations with
// session management for exec and port-forward.
type RuntimeUseCase struct {
	discovery DiscoveryClient
	runtime   RuntimeRepo
	helm      HelmRepo
	sessions  *SessionStore
}

// NewRuntimeUseCase returns a RuntimeUseCase wired to the given
// discovery, runtime, and session store backends. The SessionStore is
// injected rather than created internally so that callers can supply
// alternative implementations for testing or monitoring.
func NewRuntimeUseCase(discovery DiscoveryClient, runtime RuntimeRepo, helm HelmRepo, sessions *SessionStore) *RuntimeUseCase {
	return &RuntimeUseCase{
		discovery: discovery,
		runtime:   runtime,
		helm:      helm,
		sessions:  sessions,
	}
}

// ---------------------------------------------------------------------------
// Session ownership
// ---------------------------------------------------------------------------

// sessionOwner returns the authenticated subject a new session belongs
// to. Sessions are addressed by an identifier the client is handed back
// and later replays, so they are bound to their creator at birth.
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

// ownedBy reports whether the caller in ctx is the subject that opened
// the session.
func ownedBy(ctx context.Context, owner string) bool {
	user, ok := UserInfoFromContext(ctx)
	return ok && user.Subject != "" && user.Subject == owner
}

// execSession looks up an exec session the caller is allowed to use. A
// session belonging to somebody else is reported as missing rather than
// forbidden, so a leaked identifier cannot be used to probe for other
// users' sessions.
func (uc *RuntimeUseCase) execSession(ctx context.Context, sessionID string) (*ExecSession, error) {
	sess, ok := uc.sessions.GetExec(sessionID)
	if !ok || !ownedBy(ctx, sess.Owner) {
		return nil, &ErrSessionNotFound{Resource: "exec-session", ID: sessionID}
	}
	return sess, nil
}

// portForwardSession looks up a port-forward session the caller is
// allowed to use. See execSession.
func (uc *RuntimeUseCase) portForwardSession(ctx context.Context, sessionID string) (*PortForwardSession, error) {
	sess, ok := uc.sessions.GetPortForward(sessionID)
	if !ok || !ownedBy(ctx, sess.Owner) {
		return nil, &ErrSessionNotFound{Resource: "portforward-session", ID: sessionID}
	}
	return sess, nil
}

// vncSession looks up a VNC session the caller is allowed to use.
// See execSession.
func (uc *RuntimeUseCase) vncSession(ctx context.Context, sessionID string) (*VNCSession, error) {
	sess, ok := uc.sessions.GetVNC(sessionID)
	if !ok || !ownedBy(ctx, sess.Owner) {
		return nil, &ErrSessionNotFound{Resource: "vnc-session", ID: sessionID}
	}
	return sess, nil
}

// StartPodLogs validates the request and opens a streaming log reader.
func (uc *RuntimeUseCase) StartPodLogs(ctx context.Context, cluster, namespace, name string, opts PodLogOptions) (io.ReadCloser, error) {
	if name == "" {
		return nil, &ErrInvalidInput{Field: "name", Message: "pod name is required"}
	}
	return uc.runtime.PodLogs(ctx, cluster, namespace, name, opts)
}

// StartExec creates an exec session, starts the exec in a background
// goroutine, and returns the session together with stdout and stderr
// readers that the caller can stream from.
func (uc *RuntimeUseCase) StartExec(ctx context.Context, params *StartExecParams) (session *ExecSession, stdoutReader, stderrReader io.ReadCloser, err error) {
	if params.Name == "" {
		return nil, nil, nil, &ErrInvalidInput{Field: "name", Message: "pod name is required"}
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

	// Send initial terminal size.
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

	// Register the session BEFORE launching the goroutine to avoid
	// wasting resources if the session store is full.
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

// WriteExec writes stdin data to an active exec session. The write is
// performed in a background goroutine so that the caller's context can
// cancel a blocking pipe write during graceful shutdown or if the exec
// session has already finished.
func (uc *RuntimeUseCase) WriteExec(ctx context.Context, sessionID string, data []byte) error {
	sess, err := uc.execSession(ctx, sessionID)
	if err != nil {
		return err
	}

	// Fast-path: if the session goroutine has already exited, the
	// pipe reader is closed and Write would return immediately with
	// an error. Check Done first to return a clearer error.
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

// ResizeExec sends a terminal resize event to an active exec session.
func (uc *RuntimeUseCase) ResizeExec(ctx context.Context, sessionID string, rows, cols uint16) error {
	sess, err := uc.execSession(ctx, sessionID)
	if err != nil {
		return err
	}
	sess.SizeQueue.Set(cols, rows)
	return nil
}

// CleanupExec stops an exec session and removes it from the store.
// RemoveExec is used instead of separate Get+Delete to atomically
// claim ownership, preventing a double-close race with
// ReapStaleSessions.
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

// StartPortForward creates a port-forward session, starts the
// forwarding in a background goroutine, and returns the session
// together with a reader for data coming from the pod.
func (uc *RuntimeUseCase) StartPortForward(ctx context.Context, cluster, namespace, name string, port int32) (*PortForwardSession, io.ReadCloser, error) {
	if name == "" {
		return nil, nil, &ErrInvalidInput{Field: "name", Message: "pod name is required"}
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

	// Register the session BEFORE launching the goroutine to avoid
	// wasting resources if the session store is full.
	if err := uc.sessions.PutPortForward(sess); err != nil {
		cancel()
		dataInW.Close()
		dataInR.Close()
		dataOutW.Close()
		return nil, nil, err
	}

	// Close the read end as soon as ctx is canceled. Without this the
	// adapter's "client → pod" copy loop stays blocked reading a pipe
	// that only the goroutine below would close, and that goroutine is
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

// closeStdinOnCancel closes r once ctx is done, unblocking any read in
// flight. Pipe reads are not interruptible, so a session whose client
// sends nothing would otherwise keep its adapter goroutine alive after
// cancellation — and with it the session entry, which the reaper can
// only collect once the session's Done channel fires.
func closeStdinOnCancel(ctx context.Context, r io.Closer) {
	go func() {
		<-ctx.Done()
		if err := r.Close(); err != nil {
			slog.Debug("failed to close session stdin after cancellation", "error", err)
		}
	}()
}

// WritePortForward writes data to an active port-forward session. The
// write is performed in a background goroutine so that the caller's
// context can cancel a blocking pipe write during graceful shutdown.
func (uc *RuntimeUseCase) WritePortForward(ctx context.Context, sessionID string, data []byte) error {
	sess, err := uc.portForwardSession(ctx, sessionID)
	if err != nil {
		return err
	}

	// Fast-path: if the session goroutine has already exited, the
	// pipe reader is closed and Write would return immediately with
	// an error. Check Done first to return a clearer error.
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

// CleanupPortForward stops a port-forward session and removes it from
// the store. RemovePortForward is used instead of separate Get+Delete
// to atomically claim ownership, preventing a double-close race with
// ReapStaleSessions.
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

// StartVNC creates a VNC session, starts the VNC connection in a
// background goroutine, and returns the session together with a reader
// for data coming from the VMI.
func (uc *RuntimeUseCase) StartVNC(ctx context.Context, cluster, namespace, name string) (*VNCSession, io.ReadCloser, error) {
	if name == "" {
		return nil, nil, &ErrInvalidInput{Field: "name", Message: "VMI name is required"}
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

	// See StartPortForward: the adapter's "client → VMI" copy loop must
	// not be left blocked on a pipe this goroutine can only close after
	// the loop has finished.
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

// WriteVNC writes data to an active VNC session.
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

// CleanupVNC stops a VNC session and removes it from the store.
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

// Scale validates the inputs, looks up the GVR, updates the desired
// replica count, and returns the new value.
func (uc *RuntimeUseCase) Scale(ctx context.Context, id *ResourceIdentifier, replicas int32) (int32, error) {
	if id.Name == "" {
		return 0, &ErrInvalidInput{Field: "name", Message: "resource name is required"}
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

// StartSessionReaper launches a background goroutine that
// periodically scans for stale sessions (finished but not cleaned up)
// and removes them. It blocks until ctx is canceled.
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

// Restart validates the inputs, looks up the GVR, and triggers a
// rolling restart.
func (uc *RuntimeUseCase) Restart(ctx context.Context, id *ResourceIdentifier) error {
	if id.Name == "" {
		return &ErrInvalidInput{Field: "name", Message: "resource name is required"}
	}
	gvr, err := id.lookupGVR(ctx, uc.discovery)
	if err != nil {
		return err
	}
	return uc.runtime.Restart(ctx, id.Cluster, gvr, id.Namespace, id.Name)
}

// allowedSubResourceMethods are the HTTP methods permitted for
// SubResourceAction. DELETE and PATCH have dedicated RPCs.
var allowedSubResourceMethods = map[string]bool{"PUT": true, "POST": true}

// SubResourceAction validates the inputs and forwards a PUT/POST
// subresource action to the Kubernetes API via impersonation. This
// covers use-cases such as KubeVirt VM start/stop/restart/migrate.
func (uc *RuntimeUseCase) SubResourceAction(ctx context.Context, id *ResourceIdentifier, method string, body []byte) (map[string]any, error) {
	if id.Name == "" {
		return nil, &ErrInvalidInput{Field: "name", Message: "resource name is required"}
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

// ShowChart validates the inputs and delegates to HelmRepo to
// retrieve chart values and readme from a remote repository.
func (uc *RuntimeUseCase) ShowChart(ctx context.Context, repoURL, chartName, version string) (values, readme []byte, err error) {
	if repoURL == "" {
		return nil, nil, &ErrInvalidInput{Field: "repo_url", Message: "repository URL is required"}
	}
	if chartName == "" {
		return nil, nil, &ErrInvalidInput{Field: "chart_name", Message: "chart name is required"}
	}
	return uc.helm.ShowChart(ctx, repoURL, chartName, version)
}
