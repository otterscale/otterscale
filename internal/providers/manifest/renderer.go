// Package manifest provides the ManifestRenderer implementation that
// generates Kubernetes agent installation manifests from Go templates.
// The template and all rendering details are encapsulated here,
// keeping the domain layer (core) free of infrastructure concerns.
package manifest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"

	"github.com/otterscale/otterscale/internal/core"
)

// Renderer implements core.ManifestRenderer by executing a Go
// text/template that produces multi-document YAML.
type Renderer struct{}

// Verify at compile time that Renderer satisfies core.ManifestRenderer.
var _ core.ManifestRenderer = (*Renderer)(nil)

// NewRenderer returns a new manifest Renderer.
func NewRenderer() *Renderer {
	return &Renderer{}
}

// RenderAgentManifest produces a multi-document YAML manifest for
// installing the otterscale agent on a target Kubernetes cluster.
// The manifest includes a Namespace, ServiceAccount,
// ClusterRoleBinding (binding userName to cluster-admin), and a
// Deployment that runs the agent with the correct server/tunnel URLs.
func (r *Renderer) RenderAgentManifest(params *core.ManifestParams) (string, error) {
	data := agentManifestData{
		Cluster:   params.Cluster,
		UserName:  params.UserName,
		Image:     params.Image,
		ServerURL: params.ServerURL,
		TunnelURL: params.TunnelURL,
	}
	if params.HarborCreds != nil {
		data.HarborURL = params.HarborURL
		data.HarborRobotName = params.HarborCreds.Name
		data.HarborRobotSecret = params.HarborCreds.Secret
	}

	var buf bytes.Buffer
	if err := agentManifestTmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("render agent manifest: %w", err)
	}
	return buf.String(), nil
}

// agentManifestData holds the template parameters for agent manifest
// generation.
type agentManifestData struct {
	Cluster           string
	UserName          string
	Image             string
	ServerURL         string
	TunnelURL         string
	HarborURL         string
	HarborRobotName   string
	HarborRobotSecret string
}

// yamlQuote produces a JSON-encoded string (with surrounding quotes)
// that is safe to embed in a YAML double-quoted scalar. JSON string
// escaping is a strict subset of YAML double-quoted string escaping,
// so the result is always valid YAML regardless of the input content.
func yamlQuote(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}

// agentManifestTmpl is the parsed Go template for generating agent
// installation manifests. The "yamlQuote" function produces a
// JSON-encoded string that is safe for YAML double-quoted contexts.
var agentManifestTmpl = template.Must(
	template.New("agent-manifest").
		Funcs(template.FuncMap{"yamlQuote": yamlQuote}).
		Parse(agentManifestYAML),
)

const agentManifestYAML = `---
apiVersion: v1
kind: Namespace
metadata:
  name: otterscale-system
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: otterscale-agent
  namespace: otterscale-system
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: otterscale-agent
rules:
  # The agent proxies authenticated user requests to the local
  # kube-apiserver using impersonation headers. It must be allowed
  # to impersonate any user and group so that RBAC on the target
  # cluster enforces the actual caller's permissions.
  - apiGroups: [""]
    resources: ["users", "groups"]
    verbs: ["impersonate"]
  # Bootstrap: core resources required by FluxCD and Module CRD.
  - apiGroups: [""]
    resources: ["namespaces", "serviceaccounts", "services", "configmaps", "secrets", "resourcequotas"]
    verbs: ["get", "create", "patch"]
  # Bootstrap: workloads (FluxCD controllers).
  - apiGroups: ["apps"]
    resources: ["deployments"]
    verbs: ["get", "create", "patch"]
  # Bootstrap: RBAC for FluxCD and operator components.
  - apiGroups: ["rbac.authorization.k8s.io"]
    resources: ["clusterroles", "clusterrolebindings", "roles", "rolebindings"]
    verbs: ["get", "create", "patch", "bind", "escalate"]
  # Bootstrap: CRDs for FluxCD and Module.
  - apiGroups: ["apiextensions.k8s.io"]
    resources: ["customresourcedefinitions"]
    verbs: ["get", "create", "patch"]
  # Bootstrap: NetworkPolicy (FluxCD hardening).
  - apiGroups: ["networking.k8s.io"]
    resources: ["networkpolicies"]
    verbs: ["get", "create", "patch"]
  # Bootstrap: Admission webhooks (cert-manager + tenant-operator).
  - apiGroups: ["admissionregistration.k8s.io"]
    resources: ["mutatingwebhookconfigurations", "validatingwebhookconfigurations"]
    verbs: ["get", "create", "patch"]
  # Bootstrap: cert-manager resources (tenant-operator webhook TLS).
  - apiGroups: ["cert-manager.io"]
    resources: ["certificates", "issuers"]
    verbs: ["get", "create", "patch"]
  # Bootstrap: Module operator ModuleTemplate CRs.
  - apiGroups: ["module.otterscale.io"]
    resources: ["moduletemplates"]
    verbs: ["get", "create", "patch"]
  # Bootstrap: FluxCD source resources (GitRepository, HelmRepository).
  - apiGroups: ["source.toolkit.fluxcd.io"]
    resources: ["gitrepositories", "helmrepositories"]
    verbs: ["get", "create", "patch"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: otterscale-agent
subjects:
  - kind: ServiceAccount
    name: otterscale-agent
    namespace: otterscale-system
roleRef:
  kind: ClusterRole
  name: otterscale-agent
  apiGroup: rbac.authorization.k8s.io
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: otterscale-agent
  namespace: otterscale-system
rules:
  # The agent self-updates by patching its own Deployment image when
  # the server advertises a newer version.
  - apiGroups: ["apps"]
    resources: ["deployments"]
    resourceNames: ["otterscale-agent"]
    verbs: ["get", "patch"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: otterscale-agent
  namespace: otterscale-system
subjects:
  - kind: ServiceAccount
    name: otterscale-agent
    namespace: otterscale-system
roleRef:
  kind: Role
  name: otterscale-agent
  apiGroup: rbac.authorization.k8s.io
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: otterscale-cluster-admin
subjects:
  - kind: User
    name: {{ yamlQuote .UserName }}
    apiGroup: rbac.authorization.k8s.io
roleRef:
  kind: ClusterRole
  name: cluster-admin
  apiGroup: rbac.authorization.k8s.io
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: otterscale-node-reader
rules:
  - apiGroups: [""]
    resources: ["nodes"]
    verbs: ["get", "list", "watch"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: otterscale-node-reader
subjects:
  - kind: Group
    name: system:authenticated
    apiGroup: rbac.authorization.k8s.io
roleRef:
  kind: ClusterRole
  name: otterscale-node-reader
  apiGroup: rbac.authorization.k8s.io
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: otterscale-agent
  namespace: otterscale-system
spec:
  replicas: 1
  selector:
    matchLabels:
      app: otterscale-agent
  template:
    metadata:
      labels:
        app: otterscale-agent
    spec:
      serviceAccountName: otterscale-agent
      containers:
        - name: otterscale
          image: {{ yamlQuote .Image }}
          args:
            - agent
          env:
            - name: OTTERSCALE_AGENT_SERVER_URL
              value: {{ yamlQuote .ServerURL }}
            - name: OTTERSCALE_AGENT_TUNNEL_SERVER_URL
              value: {{ yamlQuote .TunnelURL }}
            - name: OTTERSCALE_AGENT_CLUSTER
              value: {{ yamlQuote .Cluster }}
{{- if .HarborRobotName }}
---
apiVersion: v1
kind: Secret
metadata:
  name: otterscale-harbor-robot
  namespace: otterscale-system
type: Opaque
stringData:
  HARBOR_URL: {{ yamlQuote .HarborURL }}
  HARBOR_ROBOT_NAME: {{ yamlQuote .HarborRobotName }}
  HARBOR_ROBOT_SECRET: {{ yamlQuote .HarborRobotSecret }}
{{- end }}
`
