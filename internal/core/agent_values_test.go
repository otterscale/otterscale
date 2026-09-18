package core

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

// recordingRenderer captures the resolved values, so a test can assert on what
// the domain decided rather than on YAML.
type recordingRenderer struct {
	calls int
	last  *AgentValues
}

func (r *recordingRenderer) Render(values *AgentValues) (string, error) {
	r.calls++
	r.last = values
	return "cluster: " + values.Cluster + "\n", nil
}

// recordingHarbor counts the robots it is asked to provision. The secret it
// was given is asserted through the renderer, which is where it has to arrive.
type recordingHarbor struct {
	calls int
	err   error
}

func (h *recordingHarbor) EnsureRobotAccount(
	_ context.Context, cluster, secret string,
) (HarborRobotCredentials, error) {
	h.calls++
	if h.err != nil {
		return HarborRobotCredentials{}, h.err
	}
	return HarborRobotCredentials{Name: "robot$" + cluster, Secret: secret}, nil
}

const testJoinSecret = "root-secret"

// fakeTicketStore stands in for the real store, which lives in the providers
// layer that core must not import.
type fakeTicketStore struct {
	puts    int
	tickets map[string]*AgentValuesRequest
	err     error
}

func newFakeTicketStore() *fakeTicketStore {
	return &fakeTicketStore{tickets: make(map[string]*AgentValuesRequest)}
}

func (s *fakeTicketStore) Put(
	_ context.Context, id string, req *AgentValuesRequest, _ time.Time,
) error {
	if s.err != nil {
		return s.err
	}
	s.puts++
	s.tickets[id] = req
	return nil
}

func (s *fakeTicketStore) Get(_ context.Context, id string) (*AgentValuesRequest, error) {
	req, ok := s.tickets[id]
	if !ok {
		return nil, ErrAgentValuesTicketNotFound
	}
	return req, nil
}

func testAgentChartVersions() ChartVersions {
	return ChartVersions{
		AgentFlux: "0.1.2",
		Agent:     "1.2.3",
		Flux:      "2.3.4",
	}
}

func fullAgentValuesConfig() *AgentValuesConfig {
	return &AgentValuesConfig{
		ExternalURL:     "https://otterscale.example.com/api/",
		TunnelServerURL: "https://192.0.2.1:30300",
		HarborURL:       "https://harbor.example.com:8443",
		TrustedCASecret: TrustedCASecretName,
		TrustedCAKey:    DefaultTrustedCAKey,
	}
}

func validAgentValuesRequest() *AgentValuesRequest {
	return &AgentValuesRequest{
		Cluster: "prod",
		Subject: "2ae606b7-ff7c-4c47-be76-8371963ee3c5",
		ClusterInfo: &AgentClusterInfo{
			ExternalAddress: "192.0.2.1",
			NodePortRange:   "30000-32767",
			InferenceURL:    "https://inference.example.com",
		},
	}
}

type agentValuesFixture struct {
	useCase  *AgentValuesUseCase
	renderer *recordingRenderer
	harbor   *recordingHarbor
	tickets  *fakeTicketStore
}

func newAgentValuesFixture(t *testing.T, cfg *AgentValuesConfig) agentValuesFixture {
	t.Helper()

	join, err := NewJoinAuthority(testJoinSecret)
	if err != nil {
		t.Fatalf("NewJoinAuthority() error = %v", err)
	}
	renderer := &recordingRenderer{}
	harbor := &recordingHarbor{}
	tickets := newFakeTicketStore()

	return agentValuesFixture{
		useCase:  NewAgentValuesUseCase(cfg, testAgentChartVersions(), join, tickets, renderer, harbor),
		renderer: renderer,
		harbor:   harbor,
		tickets:  tickets,
	}
}

func TestAgentValuesUseCase_Issue(t *testing.T) {
	f := newAgentValuesFixture(t, fullAgentValuesConfig())

	result, err := f.useCase.Issue(t.Context(), validAgentValuesRequest())
	if err != nil {
		t.Fatalf("Issue() error = %v", err)
	}
	if result.YAML == "" {
		t.Error("YAML is empty")
	}
	if want := "https://otterscale.example.com/api/link/values/"; !strings.HasPrefix(result.URL, want) {
		t.Errorf("URL = %q, want the prefix %q", result.URL, want)
	}
	if want := testAgentChartVersions().AgentFlux; result.ChartVersion != want {
		t.Errorf("ChartVersion = %q, want %q", result.ChartVersion, want)
	}

	got := f.renderer.last
	if got.JoinToken == "" {
		t.Fatal("the join token was not resolved")
	}
	// The token has to be the one registration verifies against, which no test
	// of the rendering alone would catch.
	join, err := NewJoinAuthority(testJoinSecret)
	if err != nil {
		t.Fatalf("NewJoinAuthority() error = %v", err)
	}
	if err := join.Verify("prod", got.JoinToken); err != nil {
		t.Errorf("the embedded join token does not verify for its own cluster: %v", err)
	}
	if got.Harbor.ModulesRepoURL != "oci://harbor.example.com:8443/modules" {
		t.Errorf("modules repo = %q", got.Harbor.ModulesRepoURL)
	}
	if got.Harbor.Robot.Name != "robot$prod" {
		t.Errorf("robot name = %q, want robot$prod", got.Harbor.Robot.Name)
	}
	if f.harbor.calls != 1 {
		t.Errorf("harbor calls = %d, want 1", f.harbor.calls)
	}
}

// TestAgentValuesUseCase_IssueAndRenderFromTicketAgree is the regression test
// for the defect that shaped this design: a Harbor-generated secret meant the
// URL carried a different credential than the RPC, and fetching it revoked the
// one in use.
func TestAgentValuesUseCase_IssueAndRenderFromTicketAgree(t *testing.T) {
	f := newAgentValuesFixture(t, fullAgentValuesConfig())

	result, err := f.useCase.Issue(t.Context(), validAgentValuesRequest())
	if err != nil {
		t.Fatalf("Issue() error = %v", err)
	}
	issuedValues := f.renderer.last

	_, id, found := strings.Cut(result.URL, "/link/values/")
	if !found {
		t.Fatalf("URL %q carries no id", result.URL)
	}

	callsAfterIssue := f.harbor.calls
	fromURL, err := f.useCase.RenderFromTicket(t.Context(), id)
	if err != nil {
		t.Fatalf("RenderFromTicket() error = %v", err)
	}

	if fromURL != result.YAML {
		t.Errorf("the URL serves different bytes than the RPC returned:\n%s\nvs\n%s", fromURL, result.YAML)
	}
	if f.harbor.calls != callsAfterIssue {
		t.Errorf("harbor calls = %d, want %d: the public endpoint must not write to Harbor",
			f.harbor.calls, callsAfterIssue)
	}
	if f.renderer.last.Harbor.Robot.Secret != issuedValues.Harbor.Robot.Secret {
		t.Error("the robot secret differs between the RPC and the URL")
	}
}

func TestAgentValuesUseCase_RenderFromTicketRejectsUnknownID(t *testing.T) {
	f := newAgentValuesFixture(t, fullAgentValuesConfig())

	if _, err := f.useCase.RenderFromTicket(t.Context(), "never-issued"); err == nil {
		t.Fatal("expected an error, got nil")
	}
	if f.renderer.calls != 0 {
		t.Error("an unknown id must not reach the renderer")
	}
	if f.harbor.calls != 0 {
		t.Error("an unknown id must not reach Harbor")
	}
}

// TestAgentValuesUseCase_URLIsShortAndOpaque is why the parameters are stored
// rather than carried: the URL travels through shell history, scrollback and
// proxy logs.
func TestAgentValuesUseCase_URLIsShortAndOpaque(t *testing.T) {
	f := newAgentValuesFixture(t, fullAgentValuesConfig())
	req := validAgentValuesRequest()

	result, err := f.useCase.Issue(t.Context(), req)
	if err != nil {
		t.Fatalf("Issue() error = %v", err)
	}

	_, id, _ := strings.Cut(result.URL, "/link/values/")
	if want := 22; len(id) != want {
		t.Errorf("id length = %d, want %d", len(id), want)
	}
	for _, secret := range []string{
		req.Cluster, req.Subject,
		req.ClusterInfo.ExternalAddress, req.ClusterInfo.InferenceURL,
	} {
		if strings.Contains(id, secret) {
			t.Errorf("the id leaks %q", secret)
		}
	}

	// Two issues must not collide, or the second would overwrite the first.
	second, err := f.useCase.Issue(t.Context(), validAgentValuesRequest())
	if err != nil {
		t.Fatalf("Issue() error = %v", err)
	}
	if second.URL == result.URL {
		t.Error("two issues produced the same id")
	}
}

// No URL may exist for a render that never succeeded.
func TestAgentValuesUseCase_NoTicketForAFailedIssue(t *testing.T) {
	f := newAgentValuesFixture(t, fullAgentValuesConfig())
	f.harbor.err = errors.New("harbor: create robot: unexpected status 500")

	if _, err := f.useCase.Issue(t.Context(), validAgentValuesRequest()); err == nil {
		t.Fatal("expected an error, got nil")
	}
	if f.tickets.puts != 0 {
		t.Error("a failed issue must not store a ticket")
	}
}

// A full store has to reach the caller, not yield a URL that serves nothing.
func TestAgentValuesUseCase_PropagatesStoreFailure(t *testing.T) {
	f := newAgentValuesFixture(t, fullAgentValuesConfig())
	f.tickets.err = &DomainError{Code: ErrorCodeResourceExhausted, Message: "too many"}

	if _, err := f.useCase.Issue(t.Context(), validAgentValuesRequest()); err == nil {
		t.Fatal("expected an error, got nil")
	}
}

// Mirrors the ordering RegisterCluster is tested for: a malformed request must
// not touch a running cluster's robot account.
func TestAgentValuesUseCase_ValidatesBeforeSideEffects(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(req *AgentValuesRequest)
		field  string
	}{
		{
			name:   "cluster info omitted",
			mutate: func(req *AgentValuesRequest) { req.ClusterInfo = nil },
			field:  fieldClusterInfo,
		},
		{
			name:   "cluster name empty",
			mutate: func(req *AgentValuesRequest) { req.Cluster = "" },
			field:  fieldCluster,
		},
		{
			name:   "cluster name is not a label value",
			mutate: func(req *AgentValuesRequest) { req.Cluster = "Prod Cluster" },
			field:  fieldCluster,
		},
		{
			name:   "subject empty",
			mutate: func(req *AgentValuesRequest) { req.Subject = "" },
			field:  "subject",
		},
		{
			name:   "external address empty",
			mutate: func(req *AgentValuesRequest) { req.ClusterInfo.ExternalAddress = "" },
			field:  fieldExternalAddress,
		},
		{
			name:   "external address carries a scheme",
			mutate: func(req *AgentValuesRequest) { req.ClusterInfo.ExternalAddress = "https://192.0.2.1" },
			field:  fieldExternalAddress,
		},
		{
			name:   "external address carries a port",
			mutate: func(req *AgentValuesRequest) { req.ClusterInfo.ExternalAddress = "192.0.2.1:8443" },
			field:  fieldExternalAddress,
		},
		{
			name:   "node port range malformed",
			mutate: func(req *AgentValuesRequest) { req.ClusterInfo.NodePortRange = "30000..32767" },
			field:  fieldNodePortRange,
		},
		{
			name:   "node port range inverted",
			mutate: func(req *AgentValuesRequest) { req.ClusterInfo.NodePortRange = "32767-30000" },
			field:  fieldNodePortRange,
		},
		{
			// The shape check passes: it is all digits either side of the
			// hyphen. Only the conversion catches it.
			name: "node port range overflows an int",
			mutate: func(req *AgentValuesRequest) {
				req.ClusterInfo.NodePortRange = "99999999999999999999-99999999999999999999"
			},
			field: fieldNodePortRange,
		},
		{
			name:   "inference url is not absolute",
			mutate: func(req *AgentValuesRequest) { req.ClusterInfo.InferenceURL = "inference.example.com" },
			field:  fieldInferenceURL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := newAgentValuesFixture(t, fullAgentValuesConfig())
			req := validAgentValuesRequest()
			tt.mutate(req)

			_, err := f.useCase.Issue(t.Context(), req)
			if err == nil {
				t.Fatal("expected an error, got nil")
			}

			var invalid *ErrInvalidInput
			if !errors.As(err, &invalid) {
				t.Fatalf("error = %v, want an *ErrInvalidInput", err)
			}
			if invalid.Field != tt.field {
				t.Errorf("field = %q, want %q", invalid.Field, tt.field)
			}
			if f.harbor.calls != 0 {
				t.Error("validation must run before the Harbor call")
			}
			if f.renderer.calls != 0 {
				t.Error("validation must run before the render")
			}
		})
	}
}

// The one cluster-info field with a genuine default rather than an example.
func TestAgentValuesUseCase_DefaultsNodePortRange(t *testing.T) {
	f := newAgentValuesFixture(t, fullAgentValuesConfig())
	req := validAgentValuesRequest()
	req.ClusterInfo.NodePortRange = ""

	if _, err := f.useCase.Issue(t.Context(), req); err != nil {
		t.Fatalf("Issue() error = %v", err)
	}
	if got := f.renderer.last.ClusterInfo.NodePortRange; got != defaultNodePortRange {
		t.Errorf("node port range = %q, want %q", got, defaultNodePortRange)
	}
}

// An empty URL has to survive as empty, because the chart default it overrides
// is an example address.
func TestAgentValuesUseCase_KeepsEmptyInferenceURL(t *testing.T) {
	f := newAgentValuesFixture(t, fullAgentValuesConfig())
	req := validAgentValuesRequest()
	req.ClusterInfo.InferenceURL = ""

	if _, err := f.useCase.Issue(t.Context(), req); err != nil {
		t.Fatalf("Issue() error = %v", err)
	}
	if got := f.renderer.last.ClusterInfo.InferenceURL; got != "" {
		t.Errorf("inference url = %q, want empty", got)
	}
}

func TestAgentValuesUseCase_ClusterAdminUsers(t *testing.T) {
	tests := []struct {
		name       string
		subject    string
		extraUsers []string
		want       []string
	}{
		{
			name:    "caller only",
			subject: "caller",
			want:    []string{"caller"},
		},
		{
			name:       "caller comes first",
			subject:    "caller",
			extraUsers: []string{"alice", "bob"},
			want:       []string{"caller", "alice", "bob"},
		},
		{
			name:       "caller naming themselves is not repeated",
			subject:    "caller",
			extraUsers: []string{"alice", "caller"},
			want:       []string{"caller", "alice"},
		},
		{
			name:       "duplicates and blanks are dropped",
			subject:    "caller",
			extraUsers: []string{"alice", "", "alice"},
			want:       []string{"caller", "alice"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := newAgentValuesFixture(t, fullAgentValuesConfig())
			req := validAgentValuesRequest()
			req.Subject = tt.subject
			req.ExtraUsers = tt.extraUsers

			if _, err := f.useCase.Issue(t.Context(), req); err != nil {
				t.Fatalf("Issue() error = %v", err)
			}

			got := f.renderer.last.ClusterAdminUsers
			if len(got) != len(tt.want) {
				t.Fatalf("users = %v, want %v", got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Fatalf("users = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

// Why the Wire providers accept empty settings: an unconfigured deployment
// fails this procedure, not the server's startup.
func TestAgentValuesUseCase_ReportsMissingConfig(t *testing.T) {
	tests := []struct {
		name  string
		clear func(cfg *AgentValuesConfig)
	}{
		{name: "no external url", clear: func(cfg *AgentValuesConfig) { cfg.ExternalURL = "" }},
		{name: "no tunnel url", clear: func(cfg *AgentValuesConfig) { cfg.TunnelServerURL = "" }},
		{name: "no harbor url", clear: func(cfg *AgentValuesConfig) { cfg.HarborURL = "" }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := fullAgentValuesConfig()
			tt.clear(cfg)
			f := newAgentValuesFixture(t, cfg)

			_, err := f.useCase.Issue(t.Context(), validAgentValuesRequest())
			if err == nil {
				t.Fatal("expected an error, got nil")
			}
			if code, ok := DomainErrorCode(err); !ok || code != ErrorCodeFailedPrecondition {
				t.Errorf("code = %v (domain=%t), want %v", code, ok, ErrorCodeFailedPrecondition)
			}
			if f.renderer.calls != 0 {
				t.Error("an unconfigured deployment must not reach the renderer")
			}
			if f.harbor.calls != 0 {
				t.Error("an unconfigured deployment must not reach Harbor")
			}
		})
	}
}

func TestAgentValuesUseCase_PropagatesHarborFailure(t *testing.T) {
	f := newAgentValuesFixture(t, fullAgentValuesConfig())
	f.harbor.err = errors.New("harbor: create robot: unexpected status 500")

	_, err := f.useCase.Issue(t.Context(), validAgentValuesRequest())
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if f.renderer.calls != 0 {
		t.Error("a failed robot provision must not produce a values file")
	}
}

// The properties the design relies on: stable per input, distinct per cluster,
// and acceptable to Harbor's complexity policy.
func TestRobotSecret(t *testing.T) {
	first := newAgentValuesFixture(t, fullAgentValuesConfig()).useCase
	second := newAgentValuesFixture(t, fullAgentValuesConfig()).useCase

	prod := first.robotSecret("prod")
	if prod != second.robotSecret("prod") {
		t.Error("the same join secret and cluster must yield the same robot secret")
	}
	if prod == first.robotSecret("staging") {
		t.Error("different clusters must not share a robot secret")
	}

	if !strings.ContainsFunc(prod, func(r rune) bool { return r >= 'A' && r <= 'Z' }) {
		t.Errorf("secret %q has no upper-case character", prod)
	}
	if !strings.ContainsFunc(prod, func(r rune) bool { return r >= 'a' && r <= 'z' }) {
		t.Errorf("secret %q has no lower-case character", prod)
	}
	if !strings.ContainsFunc(prod, func(r rune) bool { return r >= '0' && r <= '9' }) {
		t.Errorf("secret %q has no digit", prod)
	}
}

// TestRobotSecret_IsStableAcrossChanges pins the bytes, not just the
// properties above. A joined cluster's tenant operator authenticates with the
// secret it was issued, so changing this derivation breaks every one of them
// with no error that points at the cause — the same hazard the join token's
// own golden test guards.
func TestRobotSecret_IsStableAcrossChanges(t *testing.T) {
	const (
		secret = "dev-only-secret"
		want   = "OtdRa4c-aAAxPuFcgLgp7pmGbmIqfRdXts0pXfz36TNxI0"
	)

	join, err := NewJoinAuthority(secret)
	if err != nil {
		t.Fatalf("NewJoinAuthority() error = %v", err)
	}
	uc := &AgentValuesUseCase{join: join}

	if got := uc.robotSecret("dev"); got != want {
		t.Errorf("robotSecret() = %q, want %q\n"+
			"the derivation changed: every joined cluster's Harbor robot secret is now wrong",
			got, want)
	}
}

// The two derivations share a root secret, so the namespacing has to keep them
// apart: neither may ever be usable as the other.
func TestRobotSecretIsNotAJoinToken(t *testing.T) {
	join, err := NewJoinAuthority(testJoinSecret)
	if err != nil {
		t.Fatalf("NewJoinAuthority() error = %v", err)
	}
	uc := &AgentValuesUseCase{join: join}

	if uc.robotSecret("prod") == join.Token("prod") {
		t.Error("the robot secret and the join token for one cluster are the same value")
	}
	if err := join.Verify("prod", uc.robotSecret("prod")); err == nil {
		t.Error("a robot secret was accepted as a join token")
	}
}

func TestRobotSecret_DiffersByJoinSecret(t *testing.T) {
	makeUseCase := func(secret string) *AgentValuesUseCase {
		join, err := NewJoinAuthority(secret)
		if err != nil {
			t.Fatalf("NewJoinAuthority() error = %v", err)
		}
		return NewAgentValuesUseCase(
			fullAgentValuesConfig(), testAgentChartVersions(), join,
			newFakeTicketStore(), &recordingRenderer{}, &recordingHarbor{},
		)
	}

	if makeUseCase("one").robotSecret("prod") == makeUseCase("two").robotSecret("prod") {
		t.Error("rotating the join secret must change the derived robot secret")
	}
}

func TestOCIRepoURL(t *testing.T) {
	tests := []struct {
		name      string
		harborURL string
		want      string
		wantErr   bool
	}{
		{
			name:      "bare address with a port",
			harborURL: "https://192.168.196.222:8443",
			want:      "oci://192.168.196.222:8443/modules",
		},
		{
			name:      "domain without a port",
			harborURL: "https://harbor.example.com",
			want:      "oci://harbor.example.com/modules",
		},
		{
			name:      "domain with a port",
			harborURL: "https://harbor.example.com:8443",
			want:      "oci://harbor.example.com:8443/modules",
		},
		{
			name:      "trailing slash",
			harborURL: "https://harbor.example.com:8443/",
			want:      "oci://harbor.example.com:8443/modules",
		},
		{
			// The common misconfiguration. Accepting it beats emitting
			// oci:///modules, which fails much later and much less clearly.
			name:      "no scheme",
			harborURL: "harbor.example.com:8443",
			want:      "oci://harbor.example.com:8443/modules",
		},
		{ //nolint:gosec // G101: a credential in the URL is the point of this case
			// The rendered values are pasted around and land in a HelmRelease,
			// so credentials must not survive into the repository URL.
			name:      "userinfo is dropped",
			harborURL: "https://robot:hunter2@harbor.example.com:8443",
			want:      "oci://harbor.example.com:8443/modules",
		},
		{
			name:      "ipv6 with a port",
			harborURL: "https://[2001:db8::1]:8443",
			want:      "oci://[2001:db8::1]:8443/modules",
		},
		{
			name:      "empty",
			harborURL: "",
			wantErr:   true,
		},
		{
			name:      "not a URL at all",
			harborURL: "http://[::1",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ociRepoURL(tt.harborURL, modulesProject)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected an error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("ociRepoURL() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("got = %q, want %q", got, tt.want)
			}
		})
	}
}
