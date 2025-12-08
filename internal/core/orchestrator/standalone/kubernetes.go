package standalone

import (
	"strings"

	"github.com/otterscale/otterscale/internal/core/machine/tag"
	"github.com/otterscale/otterscale/internal/core/versions"
)

const defaultCalicoCIDR = "198.19.0.0/16"

type kubernetes struct {
	VirtualIPs []string
	CalicoCIDR string
}

func newKubernetes(virtualIPs []string, calicoCIDR string) base {
	return &kubernetes{
		VirtualIPs: virtualIPs,
		CalicoCIDR: calicoCIDR,
	}
}

func (k *kubernetes) Charms() []charm {
	return []charm{
		{Name: "ch:kubernetes-control-plane", Channel: versions.Kubernetes, PlacementScope: "#"},
		{Name: "ch:etcd", Channel: versions.Kubernetes, PlacementScope: "lxd"},
		{Name: "ch:easyrsa", Channel: versions.Kubernetes, PlacementScope: "lxd"},
		{Name: "ch:kubeapi-load-balancer", Channel: versions.Kubernetes, PlacementScope: "lxd"},
		{Name: "ch:calico", Channel: versions.Kubernetes, Subordinate: true},
		{Name: "ch:containerd", Channel: versions.Kubernetes, Subordinate: true},
		{Name: "ch:keepalived", Channel: versions.Kubernetes, Subordinate: true},
	}
}

func (k *kubernetes) Config(charmName string) (string, error) {
	configs := map[string]map[string]any{
		"kubernetes-control-plane": {
			"allow-privileged":     "true",
			"api-extra-args":       "service-node-port-range=1-65535",
			"loadbalancer-ips":     strings.Join(k.VirtualIPs, " "),
			"register-with-taints": "",
		},
		"kubeapi-load-balancer": {
			"loadbalancer-ips": strings.Join(k.VirtualIPs, " "),
		},
		"calico": {
			"cidr":             k.CalicoCIDR,
			"ignore-loose-rpf": "true",
		},
		"containerd": {
			"gpu_driver": "none",
		},
		"keepalived": {
			"virtual_ip": k.VirtualIPs[0],
		},
	}

	return buildConfig(charmName, configs)
}

func (k *kubernetes) Relations() [][]string {
	return [][]string{
		{"calico:cni", "kubernetes-control-plane:cni"},
		{"calico:etcd", "etcd:db"},
		{"easyrsa:client", "etcd:certificates"},
		{"easyrsa:client", "kubernetes-control-plane:certificates"},
		{"easyrsa:client", "kubeapi-load-balancer:certificates"},
		{"etcd:db", "kubernetes-control-plane:etcd"},
		{"kubernetes-control-plane:loadbalancer-external", "kubeapi-load-balancer:lb-consumers"},
		{"kubernetes-control-plane:loadbalancer-internal", "kubeapi-load-balancer:lb-consumers"},
		{"keepalived:juju-info", "kubeapi-load-balancer:juju-info"},
		{"keepalived:website", "kubeapi-load-balancer:apiserver"},
		{"containerd:containerd", "kubernetes-control-plane:container-runtime"},
	}
}

func (k *kubernetes) Tags() []string {
	return []string{
		tag.Kubernetes,
		tag.KubernetesControlPlane,
		tag.KubernetesWorker,
	}
}
