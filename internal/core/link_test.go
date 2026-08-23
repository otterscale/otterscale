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

func newTestLinkUseCase(t *testing.T, tp TunnelProvider) *LinkUseCase {
	t.Helper()
	return NewLinkUseCase(tp, "v1.0.0")
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
	uc := newTestLinkUseCase(t, tp)

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

// isErrInvalidInput checks if err is *ErrInvalidInput using the
// standard errors.As mechanism.
func isErrInvalidInput(err error, target **ErrInvalidInput) bool {
	return errors.As(err, target)
}
