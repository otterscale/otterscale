package core

import (
	"context"
	"errors"
	"strings"
	"testing"
)

// mockTunnelProvider implements TunnelProvider for testing.
type mockTunnelProvider struct {
	links       map[string]Link
	caCertPEM   []byte
	regEndpoint string
	regCertPEM  []byte
	regErr      error

	// registered records every RegisterLink call, so a test can assert
	// that a rejected request never reached the tunnel provider.
	registered []string
}

func (m *mockTunnelProvider) CACertPEM() []byte { return m.caCertPEM }
func (m *mockTunnelProvider) ListLinks() map[string]Link {
	if m.links == nil {
		return map[string]Link{}
	}
	return m.links
}

func (m *mockTunnelProvider) RegisterLink(_ context.Context, cluster, _, _ string, _ []byte) (endpoint string, certPEM []byte, err error) {
	m.registered = append(m.registered, cluster)
	return m.regEndpoint, m.regCertPEM, m.regErr
}

func (m *mockTunnelProvider) ResolveAddress(_ context.Context, _ string) (string, error) {
	return "", nil
}

// testEnrolmentSecret backs the tokens used by the link tests.
const testEnrolmentSecret = "test-root-secret"

func newTestLinkUseCase(t *testing.T, tp TunnelProvider) *LinkUseCase {
	t.Helper()
	return NewLinkUseCase(tp, "v1.0.0", newTestEnrolment(t, testEnrolmentSecret))
}

// validRegistration returns a request that passes every check, so a
// test can vary the one field it is about.
func validRegistration(t *testing.T, cluster string) *RegistrationRequest {
	t.Helper()
	return &RegistrationRequest{
		Cluster:        cluster,
		AgentID:        "agent-1",
		AgentVersion:   "v1",
		EnrolmentToken: newTestEnrolment(t, testEnrolmentSecret).Token(cluster),
		CSRPEM:         []byte("csr-data"),
	}
}

func TestLinkUseCase_ListLinks(t *testing.T) {
	links := map[string]Link{
		"prod": {Host: "127.0.0.1", AgentVersion: "v1"},
		"dev":  {Host: "127.0.0.2", AgentVersion: "v2"},
	}
	tp := &mockTunnelProvider{links: links}
	uc := newTestLinkUseCase(t, tp)

	got := uc.ListLinks(t.Context())
	if len(got) != 2 {
		t.Fatalf("expected 2 clusters, got %d", len(got))
	}
}

func TestLinkUseCase_RegisterCluster_Validation(t *testing.T) {
	tp := &mockTunnelProvider{regEndpoint: "127.0.0.1:8080", regCertPEM: []byte("cert")}
	uc := newTestLinkUseCase(t, tp)

	tests := []struct {
		name    string
		cluster string
		agentID string
		csr     []byte
		wantErr string
	}{
		{"empty cluster", "", "agent", []byte("csr"), "cluster"},
		{"cluster too long", strings.Repeat("a", 64), "agent", []byte("csr"), "must not exceed"},
		{"invalid cluster name", "UPPER", "agent", []byte("csr"), "must match"},
		{"empty agent_id", "valid-cluster", "", []byte("csr"), "agent_id"},
		{"empty csr", "valid-cluster", "agent", nil, "csr"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validRegistration(t, tt.cluster)
			req.AgentID = tt.agentID
			req.CSRPEM = tt.csr

			_, err := uc.RegisterCluster(t.Context(), req)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			var invalidInput *ErrInvalidInput
			if !isErrInvalidInput(err, &invalidInput) {
				t.Fatalf("expected ErrInvalidInput, got %T: %v", err, err)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("expected error containing %q, got %q", tt.wantErr, err.Error())
			}
		})
	}
}

func TestLinkUseCase_RegisterCluster_Success(t *testing.T) {
	tp := &mockTunnelProvider{
		regEndpoint: "127.0.0.1:8080",
		regCertPEM:  []byte("signed-cert"),
		caCertPEM:   []byte("ca-cert"),
	}
	uc := newTestLinkUseCase(t, tp)

	reg, err := uc.RegisterCluster(t.Context(), validRegistration(t, "my-cluster"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if reg.Endpoint != "127.0.0.1:8080" {
		t.Errorf("endpoint = %q, want %q", reg.Endpoint, "127.0.0.1:8080")
	}
	if string(reg.Certificate) != "signed-cert" {
		t.Errorf("certificate = %q, want %q", reg.Certificate, "signed-cert")
	}
	if string(reg.CACertificate) != "ca-cert" {
		t.Errorf("ca certificate = %q, want %q", reg.CACertificate, "ca-cert")
	}
	if reg.ServerVersion != "v1.0.0" {
		t.Errorf("server version = %q, want %q", reg.ServerVersion, "v1.0.0")
	}
}

// isErrInvalidInput checks if err is *ErrInvalidInput using the
// standard errors.As mechanism.
func isErrInvalidInput(err error, target **ErrInvalidInput) bool {
	return errors.As(err, target)
}

// TestLinkUseCase_RegisterCluster_RejectsBadToken is the regression
// test for unauthenticated registration. Registering replaces whatever
// holds the cluster name, so the important assertion is not just that
// the call fails, but that it fails before the tunnel provider is
// touched: a rejected request must not disturb the agent currently
// serving that cluster.
func TestLinkUseCase_RegisterCluster_RejectsBadToken(t *testing.T) {
	tests := []struct {
		name  string
		token func(t *testing.T) string
	}{
		{
			name:  "no token",
			token: func(*testing.T) string { return "" },
		},
		{
			name:  "made-up token",
			token: func(*testing.T) string { return "bm90LWEtcmVhbC10b2tlbi1idXQtdmFsaWQtYmFzZTY0LWhlcmU" },
		},
		{
			name: "another cluster's token",
			token: func(t *testing.T) string {
				t.Helper()
				return newTestEnrolment(t, testEnrolmentSecret).Token("staging")
			},
		},
		{
			name: "token from a rotated secret",
			token: func(t *testing.T) string {
				t.Helper()
				return newTestEnrolment(t, "some-other-secret").Token("my-cluster")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tp := &mockTunnelProvider{regEndpoint: "127.0.0.1:8080", regCertPEM: []byte("signed-cert")}
			uc := newTestLinkUseCase(t, tp)

			req := validRegistration(t, "my-cluster")
			req.EnrolmentToken = tt.token(t)

			_, err := uc.RegisterCluster(t.Context(), req)
			if err == nil {
				t.Fatal("expected an error, got nil")
			}
			if code, ok := DomainErrorCode(err); !ok || code != ErrorCodeUnauthenticated {
				t.Errorf("code = %v (domain=%v), want ErrorCodeUnauthenticated", code, ok)
			}
			if len(tp.registered) != 0 {
				t.Errorf("the tunnel provider was called for %v; a rejected registration must not change any state", tp.registered)
			}
		})
	}
}
