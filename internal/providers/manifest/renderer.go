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
// The manifest includes downstream admission prerequisites, a Namespace,
// ServiceAccount, ClusterRoleBindings, and a Deployment that runs the agent
// with the correct server/tunnel URLs.
func (r *Renderer) RenderAgentManifest(params *core.ManifestParams) (string, error) {
	projectName, err := rancherProjectName(params.RancherProjectID)
	if err != nil {
		return "", fmt.Errorf("render agent manifest: %w", err)
	}

	data := agentManifestData{
		Cluster:            params.Cluster,
		ClusterAdminUsers:  append([]string{params.UserName}, params.ExtraUsers...),
		Image:              params.Image,
		ServerURL:          params.ServerURL,
		TunnelURL:          params.TunnelURL,
		RancherProjectID:   params.RancherProjectID,
		RancherProjectName: projectName,
		HarborURL:          params.HarborURL,
	}
	if params.HarborCreds != nil {
		data.HarborRobotName = params.HarborCreds.Name
		data.HarborRobotSecret = params.HarborCreds.Secret
	}

	var buf bytes.Buffer
	if err := agentManifestTmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("render agent manifest: %w", err)
	}
	return buf.String(), nil
}

// rancherProjectName returns the Project resource name used by Kubernetes
// RBAC. An empty full ID is the supported no-selection path.
func rancherProjectName(id string) (string, error) {
	if id == "" {
		return "", nil
	}
	_, projectName, err := core.ParseRancherProjectID(id)
	return projectName, err
}

// agentManifestData holds the template parameters for agent manifest
// generation.
type agentManifestData struct {
	Cluster            string
	ClusterAdminUsers  []string
	Image              string
	ServerURL          string
	TunnelURL          string
	RancherProjectID   string
	RancherProjectName string
	HarborURL          string
	HarborRobotName    string
	HarborRobotSecret  string
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
apiVersion: admissionregistration.k8s.io/v1
kind: ValidatingAdmissionPolicy
metadata:
  name: otterscale-workspace-rancher-project
  annotations:
    security.otterscale.io/rancher-project-id: {{ yamlQuote .RancherProjectID }}
spec:
  failurePolicy: Fail
  matchConstraints:
    resourceRules:
      - apiGroups: ["tenant.otterscale.io"]
        apiVersions: ["v1alpha1"]
        operations: ["CREATE", "UPDATE"]
        resources: ["workspaces", "workspaces/status"]
        scope: Cluster
  validations:
    - expression: >-{{ if .RancherProjectID }}
        request.operation != 'CREATE' ||
        (has(object.spec.rancherProjectID) &&
        object.spec.rancherProjectID == {{ yamlQuote .RancherProjectID }}){{ else }}
        request.operation != 'CREATE' ||
        !has(object.spec.rancherProjectID) ||
        object.spec.rancherProjectID == ''{{ end }}
      message: Workspace Rancher Project must match the managed Cluster configuration
    - expression: >-
        request.operation != 'UPDATE' ||
        ((!has(oldObject.spec.rancherProjectID) && !has(object.spec.rancherProjectID)) ||
        (has(oldObject.spec.rancherProjectID) && has(object.spec.rancherProjectID) &&
        oldObject.spec.rancherProjectID == object.spec.rancherProjectID))
      message: Workspace Rancher Project is immutable
---
apiVersion: admissionregistration.k8s.io/v1
kind: ValidatingAdmissionPolicyBinding
metadata:
  name: otterscale-workspace-rancher-project
spec:
  policyName: otterscale-workspace-rancher-project
  validationActions: [Deny]
---
apiVersion: admissionregistration.k8s.io/v1
kind: ValidatingAdmissionPolicy
metadata:
  name: otterscale-workspace-namespace-security
  annotations:
    security.otterscale.io/rancher-project-id: {{ yamlQuote .RancherProjectID }}
spec:
  failurePolicy: Fail
  matchConstraints:
    resourceRules:
      - apiGroups: [""]
        apiVersions: ["v1"]
        operations: ["CREATE", "UPDATE"]
        resources: ["namespaces"]
        scope: Cluster
  matchConditions:
    - name: tenant-operator-or-workspace-namespace
      expression: >-
        request.userInfo.username == 'system:serviceaccount:otterscale-system:tenant-operator-controller-manager' ||
        (has(object.metadata.ownerReferences) && object.metadata.ownerReferences.exists(r,
        r.apiVersion == 'tenant.otterscale.io/v1alpha1' && r.kind == 'Workspace' &&
        has(r.controller) && r.controller == true)) ||
        (request.operation == 'UPDATE' && has(oldObject.metadata.ownerReferences) &&
        oldObject.metadata.ownerReferences.exists(r,
        r.apiVersion == 'tenant.otterscale.io/v1alpha1' && r.kind == 'Workspace' &&
        has(r.controller) && r.controller == true))
  validations:
    - expression: >-
        request.operation != 'CREATE' ||
        request.userInfo.username == 'system:serviceaccount:otterscale-system:tenant-operator-controller-manager'
      message: Only tenant-operator may create an OtterScale Workspace Namespace
    - expression: >-
        request.userInfo.username != 'system:serviceaccount:otterscale-system:tenant-operator-controller-manager' ||
        (has(object.metadata.ownerReferences) && object.metadata.ownerReferences.exists(r,
        r.apiVersion == 'tenant.otterscale.io/v1alpha1' && r.kind == 'Workspace' &&
        has(r.controller) && r.controller == true))
      message: tenant-operator may only modify OtterScale Workspace Namespaces
    - expression: >-
        has(object.metadata.ownerReferences) &&
        object.metadata.ownerReferences.filter(r,
        r.apiVersion == 'tenant.otterscale.io/v1alpha1' && r.kind == 'Workspace' &&
        has(r.controller) && r.controller == true).size() == 1
      message: An OtterScale Workspace Namespace must have exactly one Workspace controller owner
    - expression: >-
        request.operation != 'UPDATE' ||
        (has(oldObject.metadata.ownerReferences) &&
        oldObject.metadata.ownerReferences.filter(r,
        r.apiVersion == 'tenant.otterscale.io/v1alpha1' && r.kind == 'Workspace' &&
        has(r.controller) && r.controller == true).size() == 1 &&
        (has(object.metadata.ownerReferences) && object.metadata.ownerReferences.exists(r,
        r.apiVersion == 'tenant.otterscale.io/v1alpha1' && r.kind == 'Workspace' &&
        has(r.controller) && r.controller == true &&
        r.name == oldObject.metadata.ownerReferences.filter(o,
        o.apiVersion == 'tenant.otterscale.io/v1alpha1' && o.kind == 'Workspace' &&
        has(o.controller) && o.controller == true)[0].name &&
        r.uid == oldObject.metadata.ownerReferences.filter(o,
        o.apiVersion == 'tenant.otterscale.io/v1alpha1' && o.kind == 'Workspace' &&
        has(o.controller) && o.controller == true)[0].uid)))
      message: The Workspace controller owner is immutable
    - expression: >-
        has(object.metadata.labels) &&
        object.metadata.labels['pod-security.kubernetes.io/enforce'] == 'baseline' &&
        object.metadata.labels['pod-security.kubernetes.io/warn'] == 'restricted' &&
        object.metadata.labels['pod-security.kubernetes.io/audit'] == 'restricted' &&
        !('pod-security.kubernetes.io/enforce-version' in object.metadata.labels) &&
        !('pod-security.kubernetes.io/warn-version' in object.metadata.labels) &&
        !('pod-security.kubernetes.io/audit-version' in object.metadata.labels)
      message: Workspace Namespace Pod Security labels must remain at the OtterScale baseline
    - expression: >-{{ if .RancherProjectID }}
        request.operation != 'CREATE' ||
        (has(object.metadata.annotations) &&
        object.metadata.annotations['field.cattle.io/projectId'] == {{ yamlQuote .RancherProjectID }}){{ else }}
        request.operation != 'CREATE' ||
        !has(object.metadata.annotations) ||
        !('field.cattle.io/projectId' in object.metadata.annotations){{ end }}
      message: Workspace Namespace Rancher Project must match the managed Cluster configuration
    - expression: >-{{ if .RancherProjectID }}
        request.operation != 'UPDATE' ||
        ((!has(oldObject.metadata.annotations) || !('field.cattle.io/projectId' in oldObject.metadata.annotations)) &&
        (!has(object.metadata.annotations) || !('field.cattle.io/projectId' in object.metadata.annotations))) ||
        (has(oldObject.metadata.annotations) && has(object.metadata.annotations) &&
        'field.cattle.io/projectId' in oldObject.metadata.annotations &&
        'field.cattle.io/projectId' in object.metadata.annotations &&
        oldObject.metadata.annotations['field.cattle.io/projectId'] == object.metadata.annotations['field.cattle.io/projectId']) ||
        (request.userInfo.username == 'system:serviceaccount:otterscale-system:tenant-operator-controller-manager' &&
        has(object.metadata.annotations) &&
        object.metadata.annotations['field.cattle.io/projectId'] == {{ yamlQuote .RancherProjectID }}){{ else }}
        request.operation != 'UPDATE' ||
        ((!has(oldObject.metadata.annotations) || !('field.cattle.io/projectId' in oldObject.metadata.annotations)) &&
        (!has(object.metadata.annotations) || !('field.cattle.io/projectId' in object.metadata.annotations))) ||
        (has(oldObject.metadata.annotations) && has(object.metadata.annotations) &&
        'field.cattle.io/projectId' in oldObject.metadata.annotations &&
        'field.cattle.io/projectId' in object.metadata.annotations &&
        oldObject.metadata.annotations['field.cattle.io/projectId'] == object.metadata.annotations['field.cattle.io/projectId']){{ end }}
      message: Workspace Namespace Rancher Project may not move to an unapproved Project
---
apiVersion: admissionregistration.k8s.io/v1
kind: ValidatingAdmissionPolicyBinding
metadata:
  name: otterscale-workspace-namespace-security
spec:
  policyName: otterscale-workspace-namespace-security
  validationActions: [Deny]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: otterscale-tenant-operator-rancher-webhook
  annotations:
    security.otterscale.io/rancher-project-id: {{ yamlQuote .RancherProjectID }}
rules:
  - apiGroups: ["management.cattle.io"]
    resources: ["projects"]{{ if .RancherProjectID }}
    resourceNames: [{{ yamlQuote .RancherProjectName }}]
    verbs: ["updatepsa", "manage-namespaces"]{{ else }}
    verbs: ["updatepsa"]{{ end }}
---
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
  # Rancher webhook access is deactivated before its fail-closed guards are
  # verified on every bootstrap run.
  - apiGroups: ["rbac.authorization.k8s.io"]
    resources: ["clusterrolebindings"]
    resourceNames: ["otterscale-tenant-operator-rancher-webhook"]
    verbs: ["delete"]
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
  # Read the VAP status/spec before activating tenant-operator's synthetic
  # Rancher webhook access.
  - apiGroups: ["admissionregistration.k8s.io"]
    resources: ["validatingadmissionpolicies", "validatingadmissionpolicybindings"]
    resourceNames: ["otterscale-workspace-rancher-project", "otterscale-workspace-namespace-security"]
    verbs: ["get"]
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
{{- range .ClusterAdminUsers }}
  - kind: User
    name: {{ yamlQuote . }}
    apiGroup: rbac.authorization.k8s.io
{{- end }}
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
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: otterscale-storageclass-reader
rules:
  - apiGroups: ["storage.k8s.io"]
    resources: ["storageclasses"]
    verbs: ["get", "list", "watch"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: otterscale-storageclass-reader
subjects:
  - kind: Group
    name: system:authenticated
    apiGroup: rbac.authorization.k8s.io
roleRef:
  kind: ClusterRole
  name: otterscale-storageclass-reader
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
              value: {{ yamlQuote .Cluster }}{{ if .RancherProjectID }}
            - name: OTTERSCALE_AGENT_RANCHER_PROJECT_ID
              value: {{ yamlQuote .RancherProjectID }}{{ end }}
{{- if .HarborURL }}
            - name: OTTERSCALE_AGENT_HARBOR_URL
              value: {{ yamlQuote .HarborURL }}
{{- end }}
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
