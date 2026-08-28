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

	linkv1 "github.com/otterscale/otterscale/api/link/v1"

	"github.com/otterscale/otterscale/internal/core"
	"github.com/otterscale/otterscale/internal/pki"
)

// linkRegistrar generates a fresh CSR per registration, has the remote link
// service sign it, and returns the resulting mTLS materials.
type linkRegistrar struct {
	agentID      string
	agentVersion string // agent binary version, sent during registration
	token        string // enrolment token authorizing this cluster
	client       *http.Client
}

// NewLinkRegistrar returns a TunnelConsumer that registers agents against the
// otterscale link API over CSR-based mTLS.
func NewLinkRegistrar(version core.Version, token core.EnrolmentToken) (core.TunnelConsumer, error) {
	agentID, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("failed to get hostname: %w", err)
	}

	return &linkRegistrar{
		agentID:      agentID,
		agentVersion: string(version),
		token:        string(token),
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}, nil
}

var _ core.TunnelConsumer = (*linkRegistrar)(nil)

// Register has the server sign a fresh CSR and returns the certificate, CA
// certificate, tunnel endpoint, and server version.
//
// A new ECDSA key pair per call gives forward secrecy: a key compromised in an
// earlier session cannot decrypt a new one. The private key travels inside the
// Registration, so the cert/key pair is always consistent.
func (f *linkRegistrar) Register(ctx context.Context, serverURL, cluster string) (core.Registration, error) {
	key, keyPEM, err := pki.GenerateKey()
	if err != nil {
		return core.Registration{}, fmt.Errorf("generate key pair: %w", err)
	}

	csrPEM, err := pki.GenerateCSR(key, f.agentID)
	if err != nil {
		return core.Registration{}, fmt.Errorf("generate CSR: %w", err)
	}

	client := linkv1.NewLinkServiceClient(f.client, serverURL)
	req := &linkv1.RegisterRequest{}
	req.SetCluster(cluster)
	req.SetAgentId(f.agentID)
	req.SetCsr(csrPEM)
	req.SetAgentVersion(f.agentVersion)
	req.SetEnrolmentToken(f.token)

	resp, err := client.Register(ctx, req)
	if err != nil {
		return core.Registration{}, err
	}

	return core.Registration{
		Endpoint:       resp.GetEndpoint(),
		Certificate:    resp.GetCertificate(),
		CACertificate:  resp.GetCaCertificate(),
		PrivateKeyPEM:  keyPEM,
		TunnelUser:     resp.GetTunnelUser(),
		TunnelPassword: resp.GetTunnelPassword(),
		ServerVersion:  resp.GetServerVersion(),
	}, nil
}
