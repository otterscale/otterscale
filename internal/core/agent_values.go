package core

import (
	"cmp"
	"context"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
)

// robotSecretContext namespaces the Harbor robot secret derivation, as
// joinTokenContext does for join tokens.
const robotSecretContext = "otterscale-harbor-robot:" //nolint:gosec // a derivation namespace, not a credential

// defaultNodePortRange is the Kubernetes default. Unlike the other
// cluster-info fields this is a real default, not an example, so filling it in
// is safe.
const defaultNodePortRange = "30000-32767"

// TrustedCASecretName is the Secret the agent chart expects the CA in. Its
// tenant operator's rendered manifest mounts this exact name, so it is a
// constant rather than a setting.
const TrustedCASecretName = "otterscale-ca"

// DefaultTrustedCAKey is the key inside that Secret, as the documented
// `kubectl create secret generic otterscale-ca --from-file=ca.crt` produces.
const DefaultTrustedCAKey = "ca.crt"

// Harbor project names the module and operator charts are published under.
const (
	modulesProject   = "modules"
	operatorsProject = "operators"
)

// reNodePortRange matches "<low>-<high>", as the agent chart's own check does.
var reNodePortRange = regexp.MustCompile(`^\d+-\d+$`)

// Field names for the cluster-info validation errors.
const (
	fieldExternalAddress = "external_address"
	fieldNodePortRange   = "node_port_range"
	fieldInferenceURL    = "inference_url"
	fieldClusterInfo     = "cluster_info"
)

// AgentClusterInfo is how users reach workloads on a joining cluster. The
// control plane cannot discover any of it, so a caller supplies it.
type AgentClusterInfo struct {
	// ExternalAddress is bare: no scheme, no port.
	ExternalAddress string
	// NodePortRange is "<low>-<high>"; empty means defaultNodePortRange.
	NodePortRange string
	// InferenceURL is absolute http(s), or empty.
	InferenceURL string
}

// AgentValuesRequest is what a caller asks for. ClusterInfo is required: the
// chart's defaults for those fields are example addresses, so a default here
// would publish them.
type AgentValuesRequest struct {
	Cluster     string
	Subject     string
	ExtraUsers  []string
	ClusterInfo *AgentClusterInfo
}

// AgentValuesConfig is the deployment-wide half, resolved from configuration
// once at startup. Only the umbrella chart knows these.
type AgentValuesConfig struct {
	// ExternalURL is this API's externally reachable URL, including any path
	// prefix the gateway rewrites. It is both what agents register against and
	// the base of the URL that serves them.
	ExternalURL string
	// TunnelServerURL is the tunnel listener agents dial.
	TunnelServerURL string
	// HarborURL is the registry joining clusters are pointed at. Empty means
	// no Harbor is configured, which makes issuing values a precondition
	// failure rather than a degraded render.
	HarborURL string
	// TrustedCASecret and TrustedCAKey are empty together when this server's
	// certificate chains to a public CA and an agent needs no CA at all.
	TrustedCASecret string
	TrustedCAKey    string
}

// HarborRobotCredentials is the identity a joining cluster's tenant operator
// authenticates to Harbor as.
type HarborRobotCredentials struct {
	Name   string
	Secret string
}

// HarborSettings is the resolved Harbor half of a render.
type HarborSettings struct {
	URL              string
	ModulesRepoURL   string
	OperatorsRepoURL string
	Robot            HarborRobotCredentials
}

// AgentValues is the fully resolved answer, ready to be turned into YAML.
type AgentValues struct {
	Cluster           string
	JoinToken         string
	ServerURL         string
	TunnelServerURL   string
	ClusterAdminUsers []string
	ClusterInfo       AgentClusterInfo
	Harbor            HarborSettings
	// TrustedCASecret and TrustedCAKey are empty together, and then the render
	// carries no CA plumbing.
	TrustedCASecret string
	TrustedCAKey    string
}

// AgentValuesResult is what a caller gets back: the file, and a URL serving
// the identical bytes until ExpiresAt.
type AgentValuesResult struct {
	YAML      string
	URL       string
	ExpiresAt time.Time
}

// AgentValuesRenderer turns resolved values into the YAML a Helm install
// consumes. The file's shape is an infrastructure detail, so it lives outside
// this package.
type AgentValuesRenderer interface {
	Render(values *AgentValues) (string, error)
}

// HarborClient provisions the per-cluster robot account.
type HarborClient interface {
	// EnsureRobotAccount creates the robot if absent and otherwise resets its
	// secret to the one given, so repeated calls converge on the same
	// credential instead of revoking the one in use.
	EnsureRobotAccount(ctx context.Context, cluster, secret string) (HarborRobotCredentials, error)
}

// AgentValuesUseCase assembles the Helm override values that install an agent
// on a joining cluster. Kept apart from LinkUseCase so registration does not
// also own rendering and registry provisioning.
type AgentValuesUseCase struct {
	cfg      AgentValuesConfig
	join     *JoinAuthority
	tickets  AgentValuesStore
	renderer AgentValuesRenderer
	harbor   HarborClient
}

func NewAgentValuesUseCase(
	cfg *AgentValuesConfig,
	join *JoinAuthority,
	tickets AgentValuesStore,
	renderer AgentValuesRenderer,
	harbor HarborClient,
) *AgentValuesUseCase {
	return &AgentValuesUseCase{
		cfg:      *cfg,
		join:     join,
		tickets:  tickets,
		renderer: renderer,
		harbor:   harbor,
	}
}

// Issue renders the values and stores the request behind an opaque id, so the
// same file can be fetched from a URL.
//
// Validation runs before the Harbor call, so a malformed request cannot touch
// the robot account of a running cluster; the ticket is stored last, so no URL
// exists for a request that failed.
func (uc *AgentValuesUseCase) Issue(ctx context.Context, req *AgentValuesRequest) (AgentValuesResult, error) {
	resolved, err := uc.resolve(req)
	if err != nil {
		return AgentValuesResult{}, err
	}

	yaml, err := uc.render(ctx, resolved, uc.ensureRobot)
	if err != nil {
		return AgentValuesResult{}, err
	}

	id, err := newAgentValuesTicketID()
	if err != nil {
		return AgentValuesResult{}, err
	}
	expiresAt := time.Now().Add(agentValuesTicketTTL)
	if err := uc.tickets.Put(ctx, id, resolved, expiresAt); err != nil {
		return AgentValuesResult{}, err
	}

	return AgentValuesResult{
		YAML:      yaml,
		URL:       strings.TrimRight(uc.cfg.ExternalURL, "/") + "/link/values/" + id,
		ExpiresAt: expiresAt,
	}, nil
}

// RenderFromTicket re-renders the file a URL stands for.
//
// It never calls Harbor: the robot secret is derived, so this path writes
// nothing, needs no Harbor reachability, and cannot disagree with what Issue
// already returned.
func (uc *AgentValuesUseCase) RenderFromTicket(ctx context.Context, id string) (string, error) {
	resolved, err := uc.tickets.Get(ctx, id)
	if err != nil {
		return "", err
	}

	// Re-validated here, not because a caller can reach it, but to keep the
	// guarantee that nothing unvalidated reaches the renderer local to this
	// path.
	if err := uc.validate(resolved); err != nil {
		return "", err
	}

	return uc.render(ctx, resolved, uc.deriveRobot)
}

// robotProvider resolves the Harbor robot credentials for a cluster. Issue
// provisions them; RenderFromTicket only derives them.
type robotProvider func(ctx context.Context, cluster string) (HarborRobotCredentials, error)

func (uc *AgentValuesUseCase) render(
	ctx context.Context,
	req *AgentValuesRequest,
	robot robotProvider,
) (string, error) {
	if uc.cfg.ExternalURL == "" {
		return "", &DomainError{
			Code: ErrorCodeFailedPrecondition,
			Message: "the externally reachable server URL is not configured; " +
				"set --external-url, OTTERSCALE_SERVER_EXTERNAL_URL or server.external_url " +
				"to the URL agents use to reach this API, including the gateway's path prefix",
		}
	}
	if uc.cfg.TunnelServerURL == "" {
		return "", &DomainError{
			Code: ErrorCodeFailedPrecondition,
			Message: "the externally reachable tunnel URL is not configured; " +
				"set --external-tunnel-url, OTTERSCALE_SERVER_EXTERNAL_TUNNEL_URL " +
				"or server.external_tunnel_url to the URL agents dial for the tunnel",
		}
	}
	if uc.cfg.HarborURL == "" {
		return "", &DomainError{
			Code: ErrorCodeFailedPrecondition,
			Message: "Harbor is not configured, and a joining cluster cannot run its " +
				"tenant operator without it; set --harbor-url, OTTERSCALE_SERVER_HARBOR_URL " +
				"or server.harbor_url, together with the Harbor admin password",
		}
	}

	harbor, err := uc.harborSettings(ctx, req.Cluster, robot)
	if err != nil {
		return "", err
	}

	return uc.renderer.Render(&AgentValues{
		Cluster:           req.Cluster,
		JoinToken:         uc.join.Token(req.Cluster),
		ServerURL:         uc.cfg.ExternalURL,
		TunnelServerURL:   uc.cfg.TunnelServerURL,
		ClusterAdminUsers: clusterAdminUsers(req.Subject, req.ExtraUsers),
		ClusterInfo:       *req.ClusterInfo,
		Harbor:            harbor,
		TrustedCASecret:   uc.cfg.TrustedCASecret,
		TrustedCAKey:      uc.cfg.TrustedCAKey,
	})
}

func (uc *AgentValuesUseCase) harborSettings(
	ctx context.Context,
	cluster string,
	robot robotProvider,
) (HarborSettings, error) {
	modules, err := ociRepoURL(uc.cfg.HarborURL, modulesProject)
	if err != nil {
		return HarborSettings{}, err
	}
	operators, err := ociRepoURL(uc.cfg.HarborURL, operatorsProject)
	if err != nil {
		return HarborSettings{}, err
	}

	creds, err := robot(ctx, cluster)
	if err != nil {
		return HarborSettings{}, err
	}

	return HarborSettings{
		URL:              uc.cfg.HarborURL,
		ModulesRepoURL:   modules,
		OperatorsRepoURL: operators,
		Robot:            creds,
	}, nil
}

func (uc *AgentValuesUseCase) ensureRobot(ctx context.Context, cluster string) (HarborRobotCredentials, error) {
	return uc.harbor.EnsureRobotAccount(ctx, cluster, uc.robotSecret(cluster))
}

// deriveRobot reproduces what ensureRobot stored, without contacting Harbor.
func (uc *AgentValuesUseCase) deriveRobot(_ context.Context, cluster string) (HarborRobotCredentials, error) {
	return HarborRobotCredentials{
		Name:   HarborRobotName(cluster),
		Secret: uc.robotSecret(cluster),
	}, nil
}

// HarborRobotName is the account name Harbor assigns a system-level robot.
// Defined once: the render has to name the same account whether it went
// through Harbor or derived it, and two spellings would diverge silently.
func HarborRobotName(cluster string) string {
	return "robot$" + cluster
}

// robotSecret derives the Harbor robot password from the join secret, rather
// than letting Harbor generate one: a generated secret can only be obtained by
// replacing the account, which revokes the credential the cluster is already
// using, and the URL may be fetched repeatedly.
//
// The leading "Ot" and trailing "0" satisfy Harbor's complexity policy
// whatever the base64 happens to contain.
func (uc *AgentValuesUseCase) robotSecret(cluster string) string {
	derived := uc.join.Derive(robotSecretContext, cluster)
	return "Ot" + base64.RawURLEncoding.EncodeToString(derived) + "0"
}

// resolve validates a request and fills in the one field with a real default,
// returning a copy so a caller cannot mutate what was stored.
func (uc *AgentValuesUseCase) resolve(req *AgentValuesRequest) (*AgentValuesRequest, error) {
	if req.ClusterInfo == nil {
		return nil, &ErrInvalidInput{
			Field: fieldClusterInfo,
			Message: "must be supplied; the chart's defaults for these fields are example " +
				"addresses, so leaving them out would publish the examples instead",
		}
	}

	resolved := &AgentValuesRequest{
		Cluster:    req.Cluster,
		Subject:    req.Subject,
		ExtraUsers: slices.Clone(req.ExtraUsers),
		ClusterInfo: &AgentClusterInfo{
			ExternalAddress: req.ClusterInfo.ExternalAddress,
			NodePortRange:   cmp.Or(req.ClusterInfo.NodePortRange, defaultNodePortRange),
			InferenceURL:    strings.TrimRight(req.ClusterInfo.InferenceURL, "/"),
		},
	}
	if err := uc.validate(resolved); err != nil {
		return nil, err
	}
	return resolved, nil
}

func (uc *AgentValuesUseCase) validate(req *AgentValuesRequest) error {
	if err := ValidateClusterName(req.Cluster); err != nil {
		return err
	}
	if req.Subject == "" {
		return &ErrInvalidInput{Field: "subject", Message: msgMustNotBeEmpty}
	}
	if req.ClusterInfo == nil {
		return &ErrInvalidInput{Field: fieldClusterInfo, Message: msgMustNotBeEmpty}
	}
	return validateClusterInfo(req.ClusterInfo)
}

// validateClusterInfo mirrors the agent chart's validate.yaml, so a bad value
// is reported here rather than by `helm install` minutes later.
func validateClusterInfo(info *AgentClusterInfo) error {
	address := info.ExternalAddress
	switch {
	case address == "":
		return &ErrInvalidInput{Field: fieldExternalAddress, Message: msgMustNotBeEmpty}
	case strings.Contains(address, "://"):
		return &ErrInvalidInput{
			Field:   fieldExternalAddress,
			Message: fmt.Sprintf("must be a bare address, got %q: the scheme is the dashboard's to choose", address),
		}
	case strings.Count(address, ":") == 1:
		return &ErrInvalidInput{
			Field:   fieldExternalAddress,
			Message: fmt.Sprintf("must not carry a port, got %q: the port comes from the node port range", address),
		}
	}

	if err := validateNodePortRange(info.NodePortRange); err != nil {
		return err
	}

	if url := info.InferenceURL; url != "" &&
		!strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return &ErrInvalidInput{
			Field:   fieldInferenceURL,
			Message: fmt.Sprintf("must be an absolute http or https URL, got %q", url),
		}
	}

	return nil
}

func validateNodePortRange(portRange string) error {
	if !reNodePortRange.MatchString(portRange) {
		return &ErrInvalidInput{
			Field:   fieldNodePortRange,
			Message: fmt.Sprintf("must look like %q, got %q", defaultNodePortRange, portRange),
		}
	}

	low, high, _ := strings.Cut(portRange, "-")
	lowPort, err := strconv.Atoi(low)
	if err != nil {
		return &ErrInvalidInput{Field: fieldNodePortRange, Message: fmt.Sprintf("lower bound %q is not a number", low)}
	}
	highPort, err := strconv.Atoi(high)
	if err != nil {
		return &ErrInvalidInput{Field: fieldNodePortRange, Message: fmt.Sprintf("upper bound %q is not a number", high)}
	}
	if lowPort >= highPort {
		return &ErrInvalidInput{
			Field:   fieldNodePortRange,
			Message: fmt.Sprintf("is inverted: %q", portRange),
		}
	}
	return nil
}

// clusterAdminUsers puts the caller first, dropping duplicates so a caller
// naming themselves does not appear twice.
func clusterAdminUsers(subject string, extraUsers []string) []string {
	users := make([]string, 0, len(extraUsers)+1)
	users = append(users, subject)
	for _, user := range extraUsers {
		if user != "" && !slices.Contains(users, user) {
			users = append(users, user)
		}
	}
	return users
}

// ociRepoURL turns the Harbor URL into the OCI URL of one of its projects.
// Only the host survives, so a domain and a bare address work alike.
func ociRepoURL(harborURL, project string) (string, error) {
	host, err := harborHost(harborURL)
	if err != nil {
		return "", err
	}
	return "oci://" + path.Join(host, project), nil
}

// harborHost extracts the host[:port] from the configured Harbor URL.
//
// Parsed rather than split on "://", so that any userinfo lands in u.User and
// stays out of the repository URL: the rendered values are pasted around and
// end up in a HelmRelease, and registry credentials must not travel with them.
//
// A scheme-less value is accepted rather than rejected, since url.Parse reads
// "harbor.example.com:8443" as a scheme plus an opaque part; prefixing "//"
// makes it parse as an authority instead.
func harborHost(harborURL string) (string, error) {
	trimmed := strings.TrimSpace(harborURL)

	parsed, err := url.Parse(trimmed)
	if err == nil && parsed.Host == "" && !strings.Contains(trimmed, "://") {
		parsed, err = url.Parse("//" + trimmed)
	}
	if err != nil {
		return "", &DomainError{
			Code:    ErrorCodeFailedPrecondition,
			Message: fmt.Sprintf("the configured Harbor URL %q is not a URL", harborURL),
			Cause:   err,
		}
	}
	if parsed.Host == "" {
		return "", &DomainError{
			Code: ErrorCodeFailedPrecondition,
			Message: fmt.Sprintf(
				"the configured Harbor URL %q has no host; expected a form like https://harbor.example.com:8443",
				harborURL,
			),
		}
	}
	return parsed.Host, nil
}
