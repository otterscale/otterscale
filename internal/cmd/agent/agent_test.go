package agent

import (
	"context"
	"encoding/pem"
	"testing"

	"github.com/otterscale/otterscale/internal/core"
)

type registrationTunnel struct {
	rancherProjectID string
}

func (t *registrationTunnel) Register(_ context.Context, _, _, rancherProjectID string) (core.Registration, error) {
	t.rancherProjectID = rancherProjectID
	return core.Registration{
		Endpoint:      "127.0.0.1:16598",
		Certificate:   pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: []byte("cert")}),
		CACertificate: []byte("ca"),
		PrivateKeyPEM: []byte("key"),
		AgentID:       "agent",
		ServerVersion: "v1.0.0",
	}, nil
}

type registrationUpdater struct{}

func (registrationUpdater) Patch(context.Context, string) error { return nil }

func TestRegisterPropagatesRancherProjectID(t *testing.T) {
	tunnel := &registrationTunnel{}
	agent := &Agent{tunnel: tunnel, version: "v1.0.0", updater: registrationUpdater{}}

	if _, err := agent.register("local:p-test")(t.Context(), "https://server.example.com", "cluster"); err != nil {
		t.Fatalf("register: %v", err)
	}
	if tunnel.rancherProjectID != "local:p-test" {
		t.Fatalf("Rancher Project ID = %q", tunnel.rancherProjectID)
	}

	if _, err := agent.register("")(t.Context(), "https://server.example.com", "cluster"); err != nil {
		t.Fatalf("legacy register: %v", err)
	}
	if tunnel.rancherProjectID != "" {
		t.Fatalf("legacy Rancher Project ID = %q", tunnel.rancherProjectID)
	}
}
