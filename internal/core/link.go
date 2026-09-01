// Package core defines the domain interfaces and use-case logic for the
// otterscale agent. Infrastructure adapters (chisel, kubernetes, otterscale)
// implement the interfaces declared here.
package core

import (
	"context"
	"fmt"
	"regexp"
)

// maxClusterNameLength matches the Kubernetes label value length limit.
const maxClusterNameLength = 63

// reClusterName matches a valid Kubernetes label value. Restricting the
// character set is also what prevents YAML injection through cluster names
// carrying quotes or newlines.
var reClusterName = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]*[a-z0-9])?$`)

// ValidateClusterName returns an *ErrInvalidInput on failure.
func ValidateClusterName(cluster string) error {
	if cluster == "" {
		return &ErrInvalidInput{Field: fieldCluster, Message: msgMustNotBeEmpty}
	}
	if len(cluster) > maxClusterNameLength {
		return &ErrInvalidInput{
			Field:   fieldCluster,
			Message: fmt.Sprintf("must not exceed %d characters", maxClusterNameLength),
		}
	}
	if !reClusterName.MatchString(cluster) {
		return &ErrInvalidInput{
			Field:   fieldCluster,
			Message: fmt.Sprintf("must match [a-z0-9]([a-z0-9-]*[a-z0-9])?, got %q", cluster),
		}
	}
	return nil
}

// TunnelProvider is the server-side abstraction for reverse tunnels: it
// allocates a unique endpoint per cluster, signs agent CSRs, and provisions a
// tunnel user for each connecting agent.
type TunnelProvider interface {
	// CACertPEM lets agents verify the tunnel server and the server configure mTLS.
	CACertPEM() []byte
	ListLinks() map[string]Link
	// RegisterLink validates and signs the agent's CSR, creates a tunnel user,
	// and returns the grant the agent needs to connect.
	RegisterLink(ctx context.Context, cluster, agentID, agentVersion string, csrPEM []byte) (TunnelGrant, error)
	// ResolveAddress returns the HTTP base URL for the given cluster.
	ResolveAddress(ctx context.Context, cluster string) (string, error)
}

// TunnelGrant is what a TunnelProvider issues for one registration: where to
// connect, the certificate that authenticates the agent at the TLS layer, and
// the credential it presents to the tunnel itself.
type TunnelGrant struct {
	// Endpoint is the allocated tunnel endpoint (host:port).
	Endpoint string
	// Certificate is signed from the agent's CSR, PEM-encoded.
	Certificate []byte
	// User is assigned by the provider: nothing the agent sent may decide it.
	// Agents identify themselves by hostname, and two clusters running the same
	// deployment report the same one, so a name taken from the request would let
	// one cluster's registration silently replace another's credentials.
	User string
	// Password is generated fresh per registration rather than derived from the
	// certificate — a value the agent already holds must not also be the secret
	// that authenticates it.
	Password string
}

// TunnelConsumer is the agent-side abstraction for registering with the link
// server and obtaining tunnel credentials via CSR/mTLS.
type TunnelConsumer interface {
	Register(ctx context.Context, serverURL, cluster string) (Registration, error)
}

// Registration holds what the link server returns after a successful
// CSR-based registration.
type Registration struct {
	// Endpoint is the tunnel endpoint the agent should connect to.
	Endpoint string
	// Certificate is the X.509 certificate signed by the server's CA, used for
	// mTLS client authentication.
	Certificate []byte
	// CACertificate verifies the tunnel server's identity.
	CACertificate []byte
	// PrivateKeyPEM corresponds to the CSR sent during this registration. It is
	// returned alongside the certificate so the pair is always consistent, with
	// no TOCTOU race against a separate key fetch.
	PrivateKeyPEM []byte
	// TunnelUser and TunnelPassword are the credential the agent presents to the
	// tunnel server. The server issues them, so an agent never has to reproduce
	// a value the server computed and the two sides cannot disagree.
	TunnelUser     string
	TunnelPassword string
	// ServerVersion is reported back to agents for diagnostics.
	ServerVersion string
}

// Link is the per-cluster tunnel state.
type Link struct {
	Host         string // unique 127.x.x.x loopback address
	User         string // chisel user name
	AgentVersion string // agent binary version
}

// LinkUseCase orchestrates cluster registration on the server side, delegating
// CSR signing and tunnel setup to the TunnelProvider.
type LinkUseCase struct {
	tunnel    TunnelProvider
	version   Version
	join *JoinAuthority
}

// NewLinkUseCase takes the server binary version, included in registration
// responses, and the join that authorizes agents to claim a cluster.
func NewLinkUseCase(tunnel TunnelProvider, version Version, join *JoinAuthority) *LinkUseCase {
	return &LinkUseCase{
		tunnel:    tunnel,
		version:   version,
		join: join,
	}
}

// RegistrationRequest is what an agent presents when it registers.
type RegistrationRequest struct {
	// Cluster is the name the agent claims.
	Cluster string
	// AgentID identifies the agent process (its hostname).
	AgentID string
	// AgentVersion is kept for diagnostics.
	AgentVersion string
	// JoinToken authorizes claiming Cluster.
	JoinToken string
	// CSRPEM is the PKCS#10 request to be signed.
	CSRPEM []byte
}

func (uc *LinkUseCase) ListLinks(_ context.Context) map[string]Link {
	return uc.tunnel.ListLinks()
}

// RegisterCluster authorizes the request, then has the tunnel provider sign the
// agent's CSR.
//
// The join token is checked before anything else, because registering a
// cluster replaces whatever was registered under that name: an unauthorized
// request must not be able to disturb the agent currently serving it.
func (uc *LinkUseCase) RegisterCluster(ctx context.Context, req *RegistrationRequest) (Registration, error) {
	if err := ValidateClusterName(req.Cluster); err != nil {
		return Registration{}, err
	}
	if err := uc.join.Verify(req.Cluster, req.JoinToken); err != nil {
		return Registration{}, err
	}
	if req.AgentID == "" {
		return Registration{}, &ErrInvalidInput{Field: "agent_id", Message: msgMustNotBeEmpty}
	}
	if len(req.CSRPEM) == 0 {
		return Registration{}, &ErrInvalidInput{Field: "csr", Message: msgMustNotBeEmpty}
	}

	grant, err := uc.tunnel.RegisterLink(ctx, req.Cluster, req.AgentID, req.AgentVersion, req.CSRPEM)
	if err != nil {
		return Registration{}, err
	}
	return Registration{
		Endpoint:       grant.Endpoint,
		Certificate:    grant.Certificate,
		CACertificate:  uc.tunnel.CACertPEM(),
		TunnelUser:     grant.User,
		TunnelPassword: grant.Password,
		ServerVersion:  string(uc.version),
	}, nil
}
