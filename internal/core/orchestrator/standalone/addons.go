package standalone

import "github.com/otterscale/otterscale/internal/core/versions"

type addons struct {
	Scope string
}

func newAddons(scope string) base {
	return &addons{
		Scope: scope,
	}
}

func (a *addons) Charms() []charm {
	return []charm{
		{Name: "ch:ceph-csi", Channel: versions.Kubernetes, Subordinate: true},
		{Name: "ch:grafana-agent", Subordinate: true},
	}
}

func (a *addons) Config(charmName string) (string, error) {
	configs := map[string]map[string]any{
		"ceph-csi": {
			"default-storage":      "ceph-ext4",
			"cephfs-enable":        "true",
			"provisioner-replicas": 1,
		},
	}

	return buildConfig(a.Scope, charmName, configs)
}

func (a *addons) Relations() [][]string {
	relations := [][]string{
		{"ceph-csi", "ceph-mon"},
		{"ceph-csi", "kubernetes-control-plane"},
		{"grafana-agent:cos-agent", "ceph-mon:cos-agent"},
		{"grafana-agent:cos-agent", "kubeapi-load-balancer:cos-agent"},
		{"grafana-agent:cos-agent", "kubernetes-control-plane:cos-agent"},
	}

	return buildRelations(a.Scope, relations)
}

func (a *addons) Tags() []string {
	return []string{}
}
