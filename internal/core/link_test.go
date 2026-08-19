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
}

func (m *mockTunnelProvider) CACertPEM() []byte { return m.caCertPEM }
func (m *mockTunnelProvider) ListLinks() map[string]Link {
	if m.links == nil {
		return map[string]Link{}
	}
	return m.links
}

func (m *mockTunnelProvider) RegisterLink(_ context.Context, _, _, _ string, _ []byte) (endpoint string, certPEM []byte, err error) {
	return m.regEndpoint, m.regCertPEM, m.regErr
}

func (m *mockTunnelProvider) ResolveAddress(_ context.Context, _ string) (string, error) {
	return "", nil
}

// mockManifestRenderer implements ManifestRenderer for testing.
type mockManifestRenderer struct {
	result string
	err    error
}

func (m *mockManifestRenderer) RenderAgentManifest(_ *ManifestParams) (string, error) {
	return m.result, m.err
}

func testLinkConfig() AgentManifestConfig {
	return AgentManifestConfig{
		ServerURL: "https://server.example.com",
		TunnelURL: "https://tunnel.example.com:8300",
		HMACKey:   []byte("test-hmac-key-must-be-32-bytes!!"),
	}
}

func newTestLinkUseCase(t *testing.T, tp TunnelProvider, renderer ManifestRenderer) *LinkUseCase {
	t.Helper()
	uc, err := NewLinkUseCase(tp, "v1.0.0", testLinkConfig(), renderer, nil)
	if err != nil {
		t.Fatalf("NewLinkUseCase: %v", err)
	}
	return uc
}

func TestNewLinkUseCase_ValidationErrors(t *testing.T) {
	tp := &mockTunnelProvider{}
	renderer := &mockManifestRenderer{}

	tests := []struct {
		name    string
		cfg     AgentManifestConfig
		wantErr string
	}{
		{
			name:    "missing server URL",
			cfg:     AgentManifestConfig{TunnelURL: "x", HMACKey: []byte("k")},
			wantErr: "server URL is required",
		},
		{
			name:    "missing tunnel URL",
			cfg:     AgentManifestConfig{ServerURL: "x", HMACKey: []byte("k")},
			wantErr: "tunnel URL is required",
		},
		{
			name:    "missing HMAC key",
			cfg:     AgentManifestConfig{ServerURL: "x", TunnelURL: "x"},
			wantErr: "HMAC key is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewLinkUseCase(tp, "v1.0.0", tt.cfg, renderer, nil)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("expected error containing %q, got %q", tt.wantErr, err.Error())
			}
		})
	}
}

func TestLinkUseCase_ListLinks(t *testing.T) {
	links := map[string]Link{
		"prod": {Host: "127.0.0.1", AgentVersion: "v1"},
		"dev":  {Host: "127.0.0.2", AgentVersion: "v2"},
	}
	tp := &mockTunnelProvider{links: links}
	uc := newTestLinkUseCase(t, tp, &mockManifestRenderer{})

	got := uc.ListLinks(t.Context())
	if len(got) != 2 {
		t.Fatalf("expected 2 clusters, got %d", len(got))
	}
}

func TestLinkUseCase_RegisterCluster_Validation(t *testing.T) {
	tp := &mockTunnelProvider{regEndpoint: "127.0.0.1:8080", regCertPEM: []byte("cert")}
	uc := newTestLinkUseCase(t, tp, &mockManifestRenderer{})

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
			_, err := uc.RegisterCluster(t.Context(), tt.cluster, tt.agentID, "v1", tt.csr)
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
	uc := newTestLinkUseCase(t, tp, &mockManifestRenderer{})

	reg, err := uc.RegisterCluster(t.Context(), "my-cluster", "agent-1", "v1", []byte("csr-data"))
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

func TestLinkUseCase_ManifestToken_IssueAndVerify(t *testing.T) {
	tp := &mockTunnelProvider{}
	uc := newTestLinkUseCase(t, tp, &mockManifestRenderer{})

	url, err := uc.IssueManifestURL(t.Context(), "test-cluster", "user@example.com", []string{"bob@example.com"})
	if err != nil {
		t.Fatalf("IssueManifestURL: %v", err)
	}

	// Extract token from URL.
	parts := strings.SplitN(url, "/link/manifest/", 2)
	if len(parts) != 2 || parts[1] == "" {
		t.Fatalf("unexpected URL format: %q", url)
	}
	token := parts[1]

	claims, err := uc.VerifyManifestToken(t.Context(), token)
	if err != nil {
		t.Fatalf("VerifyManifestToken: %v", err)
	}
	if claims.Cluster != "test-cluster" {
		t.Errorf("cluster = %q, want %q", claims.Cluster, "test-cluster")
	}
	if claims.Sub != "user@example.com" {
		t.Errorf("userName = %q, want %q", claims.Sub, "user@example.com")
	}
	if len(claims.ExtraUsers) != 1 || claims.ExtraUsers[0] != "bob@example.com" {
		t.Errorf("extraUsers = %v, want [bob@example.com]", claims.ExtraUsers)
	}
}

func TestLinkUseCase_VerifyManifestToken_MalformedToken(t *testing.T) {
	tp := &mockTunnelProvider{}
	uc := newTestLinkUseCase(t, tp, &mockManifestRenderer{})

	tests := []struct {
		name  string
		token string
	}{
		{"no dots", "nodots"},
		{"empty", ""},
		{"bad base64 payload", "!!!.YWJj"},
		{"bad base64 signature", "YWJj.!!!"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := uc.VerifyManifestToken(t.Context(), tt.token)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}

func TestLinkUseCase_VerifyManifestToken_TamperedSignature(t *testing.T) {
	tp := &mockTunnelProvider{}
	uc := newTestLinkUseCase(t, tp, &mockManifestRenderer{})

	url, err := uc.IssueManifestURL(t.Context(), "test-cluster", "user@example.com", nil)
	if err != nil {
		t.Fatalf("IssueManifestURL: %v", err)
	}

	parts := strings.SplitN(url, "/link/manifest/", 2)
	token := parts[1]

	// Tamper with the signature.
	tokenParts := strings.SplitN(token, ".", 2)
	tampered := tokenParts[0] + ".dGFtcGVyZWQ"

	_, err = uc.VerifyManifestToken(t.Context(), tampered)
	if err == nil {
		t.Fatal("expected error for tampered token")
	}
	if !strings.Contains(err.Error(), "invalid or expired token") {
		t.Errorf("expected 'invalid or expired token' error, got: %v", err)
	}
}

func TestLinkUseCase_GenerateAgentManifest_Validation(t *testing.T) {
	tp := &mockTunnelProvider{}
	renderer := &mockManifestRenderer{result: "manifest-yaml"}
	uc := newTestLinkUseCase(t, tp, renderer)
	ctx := t.Context()

	tests := []struct {
		name     string
		cluster  string
		userName string
		wantErr  string
	}{
		{"empty cluster", "", "user", "cluster"},
		{"invalid cluster", "INVALID!", "user", "must match"},
		{"empty user", "valid", "", "user_name"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := uc.GenerateAgentManifest(ctx, tt.cluster, tt.userName, nil)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}

func TestLinkUseCase_GenerateAgentManifest_Success(t *testing.T) {
	tp := &mockTunnelProvider{}
	renderer := &mockManifestRenderer{result: "---\napiVersion: v1\nkind: Namespace"}
	uc := newTestLinkUseCase(t, tp, renderer)

	manifest, err := uc.GenerateAgentManifest(t.Context(), "my-cluster", "admin@example.com", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if manifest != "---\napiVersion: v1\nkind: Namespace" {
		t.Errorf("unexpected manifest: %q", manifest)
	}
}

// isErrInvalidInput checks if err is *ErrInvalidInput using the
// standard errors.As mechanism.
func isErrInvalidInput(err error, target **ErrInvalidInput) bool {
	return errors.As(err, target)
}
