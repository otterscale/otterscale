// Package core defines the domain interfaces and use-case logic for
// the otterscale agent. Infrastructure adapters (chisel, kubernetes,
// otterscale) implement the interfaces declared here.
package core

import (
	"context"
	"fmt"
	"regexp"
)

// maxClusterNameLength is the maximum allowed length for a cluster
// name. This matches the Kubernetes label value length limit.
const maxClusterNameLength = 63

// reClusterName matches a valid Kubernetes label value: lowercase
// alphanumeric characters or hyphens, must start and end with an
// alphanumeric character. This prevents YAML injection via cluster
// names that contain quotes, newlines, or other special characters.
var reClusterName = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]*[a-z0-9])?$`)

// ValidateClusterName checks that the given cluster name is non-empty,
// within the Kubernetes label value length limit, and matches the
// allowed character pattern. It returns an *ErrInvalidInput on failure.
func ValidateClusterName(cluster string) error {
	if cluster == "" {
		return &ErrInvalidInput{Field: "cluster", Message: "must not be empty"}
	}
	if len(cluster) > maxClusterNameLength {
		return &ErrInvalidInput{
			Field:   "cluster",
			Message: fmt.Sprintf("must not exceed %d characters", maxClusterNameLength),
		}
	}
	if !reClusterName.MatchString(cluster) {
		return &ErrInvalidInput{
			Field:   "cluster",
			Message: fmt.Sprintf("must match [a-z0-9]([a-z0-9-]*[a-z0-9])?, got %q", cluster),
		}
	}
	return nil
}

// TunnelProvider is the server-side abstraction for managing reverse
// tunnels. It allocates unique endpoints per cluster, signs agent
// CSRs, and provisions tunnel users for each connecting agent.
type TunnelProvider interface {
	// CACertPEM returns the PEM-encoded CA certificate so that
	// agents can verify the tunnel server and the server can
	// configure mTLS.
	CACertPEM() []byte
	// ListLinks returns the names of all registered clusters.
	ListLinks() map[string]Link
	// RegisterLink validates and signs the agent's CSR, creates
	// a tunnel user, and returns the allocated endpoint together
	// with the PEM-encoded signed certificate.
	RegisterLink(ctx context.Context, cluster, agentID, agentVersion string, csrPEM []byte) (endpoint string, certPEM []byte, err error)
	// ResolveAddress returns the HTTP base URL for the given cluster.
	ResolveAddress(ctx context.Context, cluster string) (string, error)
}

// TunnelConsumer is the agent-side abstraction for registering with
// the link server and obtaining tunnel credentials via CSR/mTLS.
type TunnelConsumer interface {
	// Register calls the link API with a CSR and returns the
	// signed certificate, CA certificate, tunnel endpoint, and the
	// private key that corresponds to the CSR. Returning the key
	// alongside the certificate eliminates the TOCTOU race that
	// would occur if callers had to fetch the key separately.
	Register(ctx context.Context, serverURL, cluster string) (Registration, error)
}

// Registration holds the credentials and connection details returned
// by the link server after a successful CSR-based registration.
type Registration struct {
	// Endpoint is the tunnel endpoint the agent should connect to.
	Endpoint string
	// Certificate is the PEM-encoded X.509 certificate signed by
	// the server's CA, used for mTLS client authentication.
	Certificate []byte
	// CACertificate is the PEM-encoded CA certificate used to
	// verify the tunnel server's identity.
	CACertificate []byte
	// PrivateKeyPEM is the PEM-encoded ECDSA private key that
	// corresponds to the CSR sent during this registration.
	// Returned alongside the certificate to ensure the key/cert
	// pair is always consistent (no TOCTOU race).
	PrivateKeyPEM []byte
	// AgentID is the identifier of the agent that registered. It is
	// set by the TunnelConsumer so that callers can derive auth
	// credentials without re-querying the hostname.
	AgentID string
	// ServerVersion is the version of the server binary, reported
	// back to agents for diagnostics.
	ServerVersion string
}

// Cluster holds the per-cluster tunnel state: the allocated
// loopback host and the chisel user name.
type Link struct {
	Host         string // unique 127.x.x.x loopback address
	User         string // chisel user name
	AgentVersion string // agent binary version
}

// LinkUseCase orchestrates cluster registration on the server side.
// It delegates CSR signing and tunnel setup to the TunnelProvider.
type LinkUseCase struct {
	tunnel    TunnelProvider
	version   Version
	enrolment *Enrolment
}

// NewLinkUseCase returns a LinkUseCase backed by the given
// TunnelProvider. version is the server binary version, included in
// registration responses. enrolment authorizes the agents allowed to
// claim a cluster.
func NewLinkUseCase(tunnel TunnelProvider, version Version, enrolment *Enrolment) *LinkUseCase {
	return &LinkUseCase{
		tunnel:    tunnel,
		version:   version,
		enrolment: enrolment,
	}
}

// RegistrationRequest is what an agent presents when it registers.
type RegistrationRequest struct {
	// Cluster is the name the agent claims.
	Cluster string
	// AgentID identifies the agent process (its hostname).
	AgentID string
	// AgentVersion is the agent binary version, kept for diagnostics.
	AgentVersion string
	// EnrolmentToken authorizes claiming Cluster.
	EnrolmentToken string
	// CSRPEM is the PEM-encoded PKCS#10 request to be signed.
	CSRPEM []byte
}

// ListLinks returns the names of all currently registered clusters.
func (uc *LinkUseCase) ListLinks(_ context.Context) map[string]Link {
	return uc.tunnel.ListLinks()
}

// RegisterCluster authorizes the request, forwards the agent's CSR to
// the tunnel provider for signing, and returns the signed certificate,
// CA certificate, tunnel endpoint, and the server's version.
//
// The enrolment token is checked before anything else happens, because
// registering a cluster replaces whatever was registered under that
// name: an unauthorized request must not be able to disturb the agent
// currently serving it.
func (uc *LinkUseCase) RegisterCluster(ctx context.Context, req *RegistrationRequest) (Registration, error) {
	if err := ValidateClusterName(req.Cluster); err != nil {
		return Registration{}, err
	}
	if err := uc.enrolment.Verify(req.Cluster, req.EnrolmentToken); err != nil {
		return Registration{}, err
	}
	if req.AgentID == "" {
		return Registration{}, &ErrInvalidInput{Field: "agent_id", Message: "must not be empty"}
	}
	if len(req.CSRPEM) == 0 {
		return Registration{}, &ErrInvalidInput{Field: "csr", Message: "must not be empty"}
	}

	endpoint, certPEM, err := uc.tunnel.RegisterLink(ctx, req.Cluster, req.AgentID, req.AgentVersion, req.CSRPEM)
	if err != nil {
		return Registration{}, err
	}
	return Registration{
		Endpoint:      endpoint,
		Certificate:   certPEM,
		CACertificate: uc.tunnel.CACertPEM(),
		ServerVersion: string(uc.version),
	}, nil
}
