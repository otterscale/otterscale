// Package otterscale implements core.TunnelConsumer by calling the
// otterscale link gRPC service (via ConnectRPC) to register an agent
// and obtain mTLS tunnel credentials.
package otterscale

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	pb "github.com/otterscale/api/link/v1"
	"github.com/otterscale/otterscale/internal/core"
	"github.com/otterscale/otterscale/internal/pki"
)

// linkRegistrar implements core.TunnelConsumer by generating a fresh
// CSR on every registration, calling the remote link service to have
// it signed, and returning the resulting mTLS materials.
type linkRegistrar struct {
	agentID      string
	agentVersion string // agent binary version, sent during registration
	client       *http.Client
}

// NewLinkRegistrar returns a TunnelConsumer that registers agents
// against the otterscale link API using CSR-based mTLS enrolment.
// A fresh ECDSA P-256 key pair and CSR are generated on every
// Register call to ensure forward secrecy — a compromised key from a
// previous session cannot decrypt traffic from a new session.
func NewLinkRegistrar(version core.Version) (core.TunnelConsumer, error) {
	agentID, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("failed to get hostname: %w", err)
	}

	return &linkRegistrar{
		agentID:      agentID,
		agentVersion: string(version),
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

var _ core.TunnelConsumer = (*linkRegistrar)(nil)

// Register generates a fresh ECDSA key pair and CSR, then calls the
// link service's Register RPC. The server signs the CSR with its
// internal CA and returns the signed certificate, CA certificate,
// tunnel endpoint, and the server's own version. A new key pair is
// generated on every call to provide forward secrecy. The private
// key is returned inside the Registration to guarantee the cert/key
// pair is always consistent (no TOCTOU race).
func (f *linkRegistrar) Register(ctx context.Context, serverURL, cluster string) (core.Registration, error) {
	key, keyPEM, err := pki.GenerateKey()
	if err != nil {
		return core.Registration{}, fmt.Errorf("generate key pair: %w", err)
	}

	csrPEM, err := pki.GenerateCSR(key, f.agentID)
	if err != nil {
		return core.Registration{}, fmt.Errorf("generate CSR: %w", err)
	}

	client := pb.NewLinkServiceClient(f.client, serverURL)
	req := &pb.RegisterRequest{}
	req.SetCluster(cluster)
	req.SetAgentId(f.agentID)
	req.SetCsr(csrPEM)
	req.SetAgentVersion(f.agentVersion)

	resp, err := client.Register(ctx, req)
	if err != nil {
		return core.Registration{}, err
	}

	return core.Registration{
		Endpoint:      resp.GetEndpoint(),
		Certificate:   resp.GetCertificate(),
		CACertificate: resp.GetCaCertificate(),
		PrivateKeyPEM: keyPEM,
		AgentID:       f.agentID,
		ServerVersion: resp.GetServerVersion(),
	}, nil
}
