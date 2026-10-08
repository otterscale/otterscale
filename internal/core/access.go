package core

import "context"

// AccessAttributes is one Kubernetes authorization question — the
// resourceAttributes of a SelfSubjectAccessReview.
type AccessAttributes struct {
	Namespace   string
	Verb        string
	Group       string
	Resource    string
	Subresource string
	Name        string
}

// AccessReviewer answers whether the caller in ctx may do attrs on cluster, as
// that cluster's own RBAC decides. An error is no answer: callers must treat
// it as neither allowed nor denied.
type AccessReviewer interface {
	Allowed(ctx context.Context, cluster string, attrs AccessAttributes) (bool, error)
}

// ClusterMetricsAccess is "may read every metric on this cluster": get on the
// prometheuses/api subresource of the Prometheus that prometheus-stack runs.
// It is the attribute that cluster's kube-rbac-proxy checks for the agent, and
// the one OpenShift's Thanos Querier checks on openshift-monitoring/k8s.
// cluster-admin has it; the otterscale-agent chart's
// otterscale:metrics-cluster-view ClusterRole grants it to anyone else.
var ClusterMetricsAccess = AccessAttributes{
	Namespace:   "monitoring",
	Verb:        "get",
	Group:       "monitoring.coreos.com",
	Resource:    "prometheuses",
	Subresource: "api",
	Name:        "prometheus-stack-kube-prom-prometheus",
}

// NamespaceMetricsAccess is "may read this namespace's metrics": get pods in
// it, which every Workspace role (admin, edit, view) includes. OpenShift asks
// for pods.metrics.k8s.io instead, which needs metrics-server; we do not ship
// it.
func NamespaceMetricsAccess(namespace string) AccessAttributes {
	return AccessAttributes{Namespace: namespace, Verb: "get", Resource: "pods"}
}
