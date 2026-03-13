// Package harbor implements the core.HarborClient interface using the
// Harbor v2.0 REST API. It creates per-cluster system-level robot
// accounts with permissions for project, member, and robot management.
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
	"os"
	"strings"
	"sync"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/otterscale/otterscale/internal/core"
)

const (
	// harborSecretNamespace is the Kubernetes namespace where the
	// Harbor admin secret is stored.
	harborSecretNamespace = "otterscale-system"
	// harborSecretName is the name of the Kubernetes Secret
	// containing the Harbor admin password.
	harborSecretName = "otterscale-harbor-admin"
	// harborSecretKey is the data key within the Secret that holds
	// the admin password.
	harborSecretKey = "HARBOR_ADMIN_PASSWORD"
)

// Client implements core.HarborClient using the Harbor v2.0 REST API.
// The admin password is lazily read from a Kubernetes Secret on first
// use and cached for subsequent calls.
type Client struct {
	baseURL    string
	httpClient *http.Client

	mu       sync.Mutex
	password string // cached after first read from K8s secret
}

// Verify at compile time that Client satisfies core.HarborClient.
var _ core.HarborClient = (*Client)(nil)

// NewClient returns a Harbor API client. baseURL is the Harbor
// instance URL (e.g. "https://harbor.example.com").
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{},
	}
}

// EnsureRobotAccount creates (or re-creates) a system-level robot
// account for the given cluster name. On 409 Conflict the existing
// robot is deleted and a new one is created to obtain a fresh secret.
func (c *Client) EnsureRobotAccount(ctx context.Context, clusterName string) (*core.HarborRobotCredentials, error) {
	password, err := c.adminPassword(ctx)
	if err != nil {
		return nil, err
	}

	creds, err := c.createRobot(ctx, clusterName, password)
	if err == nil {
		return creds, nil
	}

	// If the robot already exists, delete it and retry.
	var conflictErr *errConflict
	if !errors.As(err, &conflictErr) {
		return nil, err
	}

	robotID, err := c.findRobotID(ctx, clusterName, password)
	if err != nil {
		return nil, fmt.Errorf("find existing robot: %w", err)
	}

	if err := c.deleteRobot(ctx, robotID, password); err != nil {
		return nil, fmt.Errorf("delete existing robot: %w", err)
	}

	return c.createRobot(ctx, clusterName, password)
}

// adminPassword returns the Harbor admin password, reading it from
// the Kubernetes Secret on first call and caching the result.
func (c *Client) adminPassword(ctx context.Context) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.password != "" {
		return c.password, nil
	}

	cfg, err := kubeConfig()
	if err != nil {
		return "", fmt.Errorf("harbor: load kubernetes config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return "", fmt.Errorf("harbor: create kubernetes client: %w", err)
	}

	secret, err := clientset.CoreV1().Secrets(harborSecretNamespace).Get(ctx, harborSecretName, metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("harbor: read secret %s/%s: %w", harborSecretNamespace, harborSecretName, err)
	}

	pw, ok := secret.Data[harborSecretKey]
	if !ok || len(pw) == 0 {
		return "", fmt.Errorf("harbor: secret %s/%s missing key %q", harborSecretNamespace, harborSecretName, harborSecretKey)
	}

	c.password = string(pw)
	return c.password, nil
}

// robotRequest is the JSON body for POST /api/v2.0/robots.
type robotRequest struct {
	Name        string            `json:"name"`
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

// robotResponse is the JSON response from a successful robot creation.
type robotResponse struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Secret string `json:"secret"`
}

// robotListItem represents a single robot in the list response.
type robotListItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// createRobot sends POST /api/v2.0/robots to create a system-level
// robot account with the required permissions.
func (c *Client) createRobot(ctx context.Context, clusterName, password string) (*core.HarborRobotCredentials, error) {
	body := robotRequest{
		Name:     clusterName,
		Duration: -1,
		Level:    "system",
		Permissions: []robotPermission{
			{
				Kind:      "system",
				Namespace: "/",
				Access: []robotAccess{
					{Resource: "project", Action: "list"},
					{Resource: "project", Action: "create"},
				},
			},
			{
				Kind:      "project",
				Namespace: "*",
				Access: []robotAccess{
					{Resource: "member", Action: "list"},
					{Resource: "member", Action: "read"},
					{Resource: "member", Action: "create"},
					{Resource: "member", Action: "update"},
					{Resource: "member", Action: "delete"},
					{Resource: "robot", Action: "list"},
					{Resource: "robot", Action: "read"},
					{Resource: "robot", Action: "create"},
					{Resource: "robot", Action: "delete"},
				},
			},
		},
	}

	data, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("harbor: marshal robot request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/v2.0/robots", bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("harbor: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("admin", password)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("harbor: create robot: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusConflict {
		return nil, &errConflict{cluster: clusterName}
	}

	if resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("harbor: create robot: unexpected status %d: %s", resp.StatusCode, string(respBody))
	}

	var result robotResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("harbor: decode robot response: %w", err)
	}

	return &core.HarborRobotCredentials{
		Name:   result.Name,
		Secret: result.Secret,
	}, nil
}

// findRobotID searches for a robot by cluster name and returns its ID.
func (c *Client) findRobotID(ctx context.Context, clusterName, password string) (int, error) {
	query := url.Values{}
	query.Set("q", "name=~"+clusterName)

	reqURL := c.baseURL + "/api/v2.0/robots?" + query.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return 0, fmt.Errorf("harbor: create request: %w", err)
	}
	req.SetBasicAuth("admin", password)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("harbor: list robots: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("harbor: list robots: unexpected status %d: %s", resp.StatusCode, string(respBody))
	}

	var robots []robotListItem
	if err := json.NewDecoder(resp.Body).Decode(&robots); err != nil {
		return 0, fmt.Errorf("harbor: decode robots list: %w", err)
	}

	// Match on the full robot name suffix (Harbor prefixes with "robot$").
	suffix := "$" + clusterName
	for _, r := range robots {
		if strings.HasSuffix(r.Name, suffix) {
			return r.ID, nil
		}
	}

	return 0, fmt.Errorf("harbor: robot for cluster %q not found", clusterName)
}

// deleteRobot sends DELETE /api/v2.0/robots/{id}.
func (c *Client) deleteRobot(ctx context.Context, robotID int, password string) error {
	reqURL := fmt.Sprintf("%s/api/v2.0/robots/%d", c.baseURL, robotID)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, reqURL, nil)
	if err != nil {
		return fmt.Errorf("harbor: create request: %w", err)
	}
	req.SetBasicAuth("admin", password)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("harbor: delete robot: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("harbor: delete robot: unexpected status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

// errConflict signals that a robot with the same name already exists.
type errConflict struct {
	cluster string
}

func (e *errConflict) Error() string {
	return fmt.Sprintf("harbor: robot for cluster %q already exists", e.cluster)
}

// kubeConfig returns a Kubernetes REST config, using in-cluster config in production and falling back to local kubeconfig for debugging.
func kubeConfig() (*rest.Config, error) {
	if os.Getenv("OTTERSCALE_DEBUG") != "" {
		return clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	}
	return rest.InClusterConfig()
}
