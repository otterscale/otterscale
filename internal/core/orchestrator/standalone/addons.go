package standalone

import "github.com/otterscale/otterscale/internal/core/versions"

type addons struct{}

func newAddons() base {
	return &addons{}
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
			"cephfs-enable":        "true",
			"default-storage":      "ceph-ext4",
			"image-registry":       "ghcr.io/otterscale",
			"provisioner-replicas": 1,
			"release":              "v" + versions.CephCSI,
		},
	}

	return buildConfig(charmName, configs)
}

func (a *addons) Relations() [][]string {
	return [][]string{
		{"ceph-csi", "ceph-mon"},
		{"ceph-csi", "kubernetes-control-plane"},
		{"grafana-agent:cos-agent", "ceph-mon:cos-agent"},
		{"grafana-agent:cos-agent", "kubeapi-load-balancer:cos-agent"},
		{"grafana-agent:cos-agent", "kubernetes-control-plane:cos-agent"},
	}
}

func (a *addons) Tags() []string {
	return []string{}
}
