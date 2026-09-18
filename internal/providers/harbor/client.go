// Package harbor implements core.HarborClient against the Harbor v2.0 REST
// API, provisioning the system-level robot account a joining cluster's tenant
// operator authenticates as. The admin credential arrives as a function, so
// this package needs no cluster access.
package harbor

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/otterscale/otterscale/internal/core"
)

const (
	// robotsPath is the Harbor v2.0 robot collection.
	robotsPath = "/api/v2.0/robots"

	// adminUser is the account the admin password belongs to.
	adminUser = "admin"

	// robotNeverExpires is Harbor's sentinel for no expiry: the credential
	// lasts as long as the cluster's membership.
	robotNeverExpires = -1

	// robotLevel must be system-wide to provision projects for workspaces that
	// do not exist yet.
	robotLevel = "system"

	// requestTimeout bounds a single call, so an unresponsive Harbor does not
	// hold the RPC open until the transport's much longer timeout.
	requestTimeout = 15 * time.Second
)

// Client talks to one Harbor instance.
type Client struct {
	baseURL    string
	httpClient *http.Client
	// Called per request, not read once: the Secret the chart mounts it from
	// may appear after this process starts, and rotation must not need a
	// restart.
	password func() (string, error)
}

var _ core.HarborClient = (*Client)(nil)

// NewClient returns a client for the Harbor at baseURL, authenticating as the
// admin whose password the given function resolves.
func NewClient(baseURL string, password func() (string, error)) *Client {
	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{Timeout: requestTimeout},
		password:   password,
	}
}

// EnsureRobotAccount converges the robot on the given secret: it creates the
// account when absent, and otherwise resets the existing secret to the same
// value. It never deletes and re-creates, which is what obtaining a
// Harbor-generated secret would require, and would revoke the credential the
// cluster is already using.
func (c *Client) EnsureRobotAccount(ctx context.Context, cluster, secret string) (core.HarborRobotCredentials, error) {
	password, err := c.adminPassword()
	if err != nil {
		return core.HarborRobotCredentials{}, err
	}

	creds, err := c.createRobot(ctx, cluster, secret, password)
	if err == nil {
		return creds, nil
	}

	var conflict *errConflict
	if !errors.As(err, &conflict) {
		return core.HarborRobotCredentials{}, err
	}

	// The account is already there, so adopt it: find it and write the secret
	// we derived, which also repairs one that was changed out of band.
	id, err := c.findRobotID(ctx, cluster, password)
	if err != nil {
		return core.HarborRobotCredentials{}, err
	}
	if err := c.refreshRobotSecret(ctx, id, secret, password); err != nil {
		return core.HarborRobotCredentials{}, err
	}

	return core.HarborRobotCredentials{Name: robotName(cluster), Secret: secret}, nil
}

const msgPasswordUnavailable = "the Harbor admin password is not available; " +
	"set --harbor-admin-password, --harbor-admin-password-file or " +
	"OTTERSCALE_SERVER_HARBOR_ADMIN_PASSWORD, and check that the Secret behind it exists"

// adminPassword turns an unconfigured credential into an actionable
// precondition failure rather than an unauthorized call.
func (c *Client) adminPassword() (string, error) {
	password, err := c.password()
	if err != nil {
		return "", &core.DomainError{
			Code:    core.ErrorCodeFailedPrecondition,
			Message: msgPasswordUnavailable,
			Cause:   err,
		}
	}
	if password == "" {
		return "", &core.DomainError{
			Code:    core.ErrorCodeFailedPrecondition,
			Message: msgPasswordUnavailable,
		}
	}
	return password, nil
}

// robotName is what Harbor assigns a system-level robot.
func robotName(cluster string) string {
	return "robot$" + cluster
}

// robotRequest is the body of a robot create.
type robotRequest struct {
	Name        string            `json:"name"`
	Secret      string            `json:"secret"`
	Duration    int               `json:"duration"`
	Level       string            `json:"level"`
	Permissions []robotPermission `json:"permissions"`
}

type robotPermission struct {
	Kind      string        `json:"kind"`
	Namespace string        `json:"namespace"`
	Access    []robotAccess `json:"access"`
}

type robotAccess struct {
	Resource string `json:"resource"`
	Action   string `json:"action"`
}

type robotResponse struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Secret string `json:"secret"`
}

type robotListItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// robotSecretRequest is the body of a secret refresh.
type robotSecretRequest struct {
	Secret string `json:"secret"`
}

// Harbor permission resources and actions.
const (
	resourceProject    = "project"
	resourceMember     = "member"
	resourceRobot      = "robot"
	resourceRepository = "repository"

	actionList   = "list"
	actionRead   = "read"
	actionCreate = "create"
	actionUpdate = "update"
	actionDelete = "delete"
	actionPull   = "pull"
	actionPush   = "push"
)

// robotPermissions is what the tenant operator needs: enough to discover and
// create projects, and within each one to manage members, robots and images.
func robotPermissions() []robotPermission {
	return []robotPermission{
		{
			Kind:      "system",
			Namespace: "/",
			Access: []robotAccess{
				{Resource: resourceProject, Action: actionList},
				{Resource: resourceProject, Action: actionCreate},
				{Resource: resourceRobot, Action: actionList},
				{Resource: resourceRobot, Action: actionRead},
			},
		},
		{
			Kind:      resourceProject,
			Namespace: "*",
			Access: []robotAccess{
				{Resource: resourceProject, Action: actionRead},
				{Resource: resourceMember, Action: actionList},
				{Resource: resourceMember, Action: actionRead},
				{Resource: resourceMember, Action: actionCreate},
				{Resource: resourceMember, Action: actionUpdate},
				{Resource: resourceMember, Action: actionDelete},
				{Resource: resourceRobot, Action: actionList},
				{Resource: resourceRobot, Action: actionRead},
				{Resource: resourceRobot, Action: actionCreate},
				{Resource: resourceRobot, Action: actionDelete},
				{Resource: resourceRepository, Action: actionPull},
				{Resource: resourceRepository, Action: actionPush},
			},
		},
	}
}

// createRobot reports a 409 as errConflict, the signal to adopt the existing
// account.
func (c *Client) createRobot(ctx context.Context, cluster, secret, password string) (core.HarborRobotCredentials, error) {
	// The secret travels in the body by design: a caller-chosen one is what
	// makes issuing agent values idempotent.
	body, err := json.Marshal(robotRequest{ //nolint:gosec // G117: sending the secret is the operation
		Name:        cluster,
		Secret:      secret,
		Duration:    robotNeverExpires,
		Level:       robotLevel,
		Permissions: robotPermissions(),
	})
	if err != nil {
		return core.HarborRobotCredentials{}, fmt.Errorf("harbor: marshal robot request: %w", err)
	}

	resp, err := c.do(ctx, http.MethodPost, c.baseURL+robotsPath, password, bytes.NewReader(body))
	if err != nil {
		return core.HarborRobotCredentials{}, fmt.Errorf("harbor: create robot: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusConflict {
		return core.HarborRobotCredentials{}, &errConflict{cluster: cluster}
	}
	if resp.StatusCode != http.StatusCreated {
		return core.HarborRobotCredentials{}, statusError(resp, "create robot")
	}

	var result robotResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return core.HarborRobotCredentials{}, fmt.Errorf("harbor: decode robot response: %w", err)
	}

	// The name comes from Harbor, which owns the prefixing; the secret does
	// not, since some versions omit it from the response.
	return core.HarborRobotCredentials{Name: result.Name, Secret: secret}, nil
}

// findRobotID locates the robot for a cluster.
func (c *Client) findRobotID(ctx context.Context, cluster, password string) (int, error) {
	query := url.Values{}
	query.Set("q", "name=~"+cluster)

	resp, err := c.do(ctx, http.MethodGet, c.baseURL+robotsPath+"?"+query.Encode(), password, http.NoBody)
	if err != nil {
		return 0, fmt.Errorf("harbor: list robots: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, statusError(resp, "list robots")
	}

	var robots []robotListItem
	if err := json.NewDecoder(resp.Body).Decode(&robots); err != nil {
		return 0, fmt.Errorf("harbor: decode robots list: %w", err)
	}

	// The query is a substring match, so compare the whole name: "prod" must
	// not adopt "eu-prod"'s robot.
	for _, robot := range robots {
		if robot.Name == robotName(cluster) {
			return robot.ID, nil
		}
	}

	return 0, fmt.Errorf("harbor: robot for cluster %q not found", cluster)
}

// refreshRobotSecret sets an existing robot's secret.
func (c *Client) refreshRobotSecret(ctx context.Context, id int, secret, password string) error {
	body, err := json.Marshal(robotSecretRequest{Secret: secret}) //nolint:gosec // G117: sending the secret is the operation
	if err != nil {
		return fmt.Errorf("harbor: marshal robot secret request: %w", err)
	}

	robotURL := fmt.Sprintf("%s%s/%d", c.baseURL, robotsPath, id)
	resp, err := c.do(ctx, http.MethodPatch, robotURL, password, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("harbor: refresh robot secret: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return statusError(resp, "refresh robot secret")
	}
	return nil
}

func (c *Client) do(ctx context.Context, method, reqURL, password string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, reqURL, body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	if body != http.NoBody {
		req.Header.Set("Content-Type", "application/json")
	}
	req.SetBasicAuth(adminUser, password)

	return c.httpClient.Do(req) //nolint:wrapcheck // callers name the operation
}

// statusError includes the body, where Harbor puts the reason. A rejected
// credential becomes a precondition failure, because no retry will help.
func statusError(resp *http.Response, operation string) error {
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return &core.DomainError{
			Code: core.ErrorCodeFailedPrecondition,
			Message: "Harbor rejected the configured admin credential; " +
				"check --harbor-admin-password-file against the password in use",
			Cause: fmt.Errorf("harbor: %s: status %d: %s", operation, resp.StatusCode, body),
		}
	}

	return fmt.Errorf("harbor: %s: unexpected status %d: %s", operation, resp.StatusCode, body)
}

// errConflict means a robot with that name is already registered.
type errConflict struct {
	cluster string
}

func (e *errConflict) Error() string {
	return fmt.Sprintf("harbor: robot for cluster %q already exists", e.cluster)
}
