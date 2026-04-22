// Package core defines the domain interfaces and use-case logic for
// the otterscale agent. Infrastructure adapters (chisel, kubernetes,
// otterscale) implement the interfaces declared here.
package core

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
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
	// ServerVersion is the version of the server binary. Agents
	// compare this against their own version to decide whether a
	// self-update is needed.
	ServerVersion string
}

// Cluster holds the per-cluster tunnel state: the allocated
// loopback host and the chisel user name.
type Link struct {
	Host         string // unique 127.x.x.x loopback address
	User         string // chisel user name
	AgentVersion string // agent binary version
}

// HarborRobotCredentials holds the name and secret for a Harbor
// robot account created via the Harbor v2.0 API.
type HarborRobotCredentials struct {
	// Name is the full robot account name (e.g. "robot$my-cluster").
	Name string
	// Secret is the robot account token/password.
	Secret string
}

// HarborClient creates per-cluster robot accounts in Harbor.
// Implementations live in the providers layer.
type HarborClient interface {
	// EnsureRobotAccount creates (or re-creates) a system-level
	// robot account for the given cluster name and returns its
	// credentials. On conflict, the existing robot is deleted and
	// re-created to obtain a fresh secret.
	EnsureRobotAccount(ctx context.Context, clusterName string) (*HarborRobotCredentials, error)
}

// AgentManifestConfig holds the external URLs and HMAC key needed to
// generate agent installation manifests and sign manifest tokens.
type AgentManifestConfig struct {
	// ServerURL is the externally reachable URL of the control-plane
	// server (e.g. "https://otterscale.example.com").
	ServerURL string
	// TunnelURL is the externally reachable URL of the tunnel server
	// (e.g. "https://tunnel.example.com:8300").
	TunnelURL string
	// HMACKey is a 32-byte key derived from the CA seed via HKDF.
	// It is used to sign and verify stateless manifest tokens.
	HMACKey []byte
	// HarborURL is the externally reachable Harbor registry URL.
	// Empty when Harbor integration is disabled.
	HarborURL string
}

// ManifestParams holds the parameters needed to render an agent
// installation manifest. It is defined in the core layer as a
// pure value object; the rendering logic lives in the providers layer.
type ManifestParams struct {
	Cluster   string
	UserName  string
	Image     string
	ServerURL string
	TunnelURL string
	// ExtraUsers are additional user identities bound to cluster-admin
	// via the otterscale-cluster-admin ClusterRoleBinding, in addition
	// to UserName. Already deduplicated and non-empty.
	ExtraUsers []string
	// HarborURL is the Harbor registry URL. Empty when Harbor
	// integration is disabled.
	HarborURL string
	// HarborCreds holds the per-cluster robot account credentials.
	// Nil when Harbor integration is disabled.
	HarborCreds *HarborRobotCredentials
}

// ManifestRenderer renders agent installation manifests from the given
// parameters. Implementations live in the providers layer and own the
// template and formatting details.
type ManifestRenderer interface {
	RenderAgentManifest(params *ManifestParams) (string, error)
}

// LinkUseCase orchestrates cluster registration on the server side.
// It delegates CSR signing and tunnel setup to the TunnelProvider,
// and token management to the ManifestTokenIssuer.
type LinkUseCase struct {
	tunnel      TunnelProvider
	version     Version
	manifestCfg AgentManifestConfig
	renderer    ManifestRenderer
	tokenIssuer *ManifestTokenIssuer
	harbor      HarborClient // nil when Harbor integration is disabled
}

// NewLinkUseCase returns a LinkUseCase backed by the given
// TunnelProvider. version is the server binary version, included in
// registration responses so agents can detect mismatches.
// manifestCfg provides the external URLs embedded in generated agent
// installation manifests. It returns an error if any required
// manifest configuration field is missing.
func NewLinkUseCase(tunnel TunnelProvider, version Version, manifestCfg AgentManifestConfig, renderer ManifestRenderer, harbor HarborClient) (*LinkUseCase, error) {
	if manifestCfg.ServerURL == "" {
		return nil, fmt.Errorf("manifest config: server URL is required")
	}
	if manifestCfg.TunnelURL == "" {
		return nil, fmt.Errorf("manifest config: tunnel URL is required")
	}
	tokenIssuer, err := NewManifestTokenIssuer(manifestCfg.HMACKey)
	if err != nil {
		return nil, err
	}
	return &LinkUseCase{
		tunnel:      tunnel,
		version:     version,
		manifestCfg: manifestCfg,
		renderer:    renderer,
		tokenIssuer: tokenIssuer,
		harbor:      harbor,
	}, nil
}

// ListLinks returns the names of all currently registered clusters.
func (uc *LinkUseCase) ListLinks(_ context.Context) map[string]Link {
	return uc.tunnel.ListLinks()
}

// RegisterCluster validates the inputs, forwards the agent's CSR to
// the tunnel provider for signing, and returns the signed certificate,
// CA certificate, tunnel endpoint, and the server's version.
func (uc *LinkUseCase) RegisterCluster(ctx context.Context, cluster, agentID, agentVersion string, csrPEM []byte) (Registration, error) {
	if err := ValidateClusterName(cluster); err != nil {
		return Registration{}, err
	}
	if agentID == "" {
		return Registration{}, &ErrInvalidInput{Field: "agent_id", Message: "must not be empty"}
	}
	if len(csrPEM) == 0 {
		return Registration{}, &ErrInvalidInput{Field: "csr", Message: "must not be empty"}
	}

	endpoint, certPEM, err := uc.tunnel.RegisterLink(ctx, cluster, agentID, agentVersion, csrPEM)
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

// IssueManifestURL generates an HMAC-signed token that encodes the
// cluster name, user identity, and extra users bound to cluster-admin,
// and returns a full URL that serves the agent manifest as raw YAML.
// The token is valid for manifestTokenTTL.
func (uc *LinkUseCase) IssueManifestURL(_ context.Context, cluster, userName string, extraUsers []string) (string, error) {
	token, err := uc.tokenIssuer.Issue(cluster, userName, extraUsers)
	if err != nil {
		return "", fmt.Errorf("issue manifest token: %w", err)
	}
	return strings.TrimRight(uc.manifestCfg.ServerURL, "/") + "/link/manifest/" + token, nil
}

// VerifyManifestToken validates the HMAC signature and expiry of a
// manifest token and returns the extracted claims. All verification
// failures return a generic error to avoid leaking which stage failed;
// detailed reasons are logged at debug level.
func (uc *LinkUseCase) VerifyManifestToken(_ context.Context, token string) (ManifestTokenClaims, error) {
	claims, err := uc.tokenIssuer.Verify(token)
	if err != nil {
		slog.Debug("manifest token verification failed", "error", err)
		return ManifestTokenClaims{}, err
	}
	return claims, nil
}

// GenerateAgentManifest produces a multi-document YAML manifest for
// installing the otterscale agent on a target Kubernetes cluster.
// The manifest includes a Namespace, ServiceAccount,
// ClusterRoleBinding (binding userName to cluster-admin), and a
// Deployment that runs the agent with the correct server/tunnel URLs.
func (uc *LinkUseCase) GenerateAgentManifest(ctx context.Context, cluster, userName string, extraUsers []string) (string, error) {
	if err := ValidateClusterName(cluster); err != nil {
		return "", err
	}
	if userName == "" {
		return "", &ErrInvalidInput{Field: "user_name", Message: "must not be empty"}
	}

	params := &ManifestParams{
		Cluster:    cluster,
		UserName:   userName,
		ExtraUsers: extraUsers,
		Image:      fmt.Sprintf("ghcr.io/otterscale/otterscale:%s", uc.version),
		ServerURL:  uc.manifestCfg.ServerURL,
		TunnelURL:  uc.manifestCfg.TunnelURL,
	}

	if uc.harbor != nil {
		creds, err := uc.harbor.EnsureRobotAccount(ctx, cluster)
		if err != nil {
			return "", fmt.Errorf("create harbor robot account: %w", err)
		}
		params.HarborURL = uc.manifestCfg.HarborURL
		params.HarborCreds = creds
	}

	return uc.renderer.RenderAgentManifest(params)
}
