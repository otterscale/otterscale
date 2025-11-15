package standalone

import (
	"strings"

	"github.com/otterscale/otterscale/internal/core/machine/tag"
)

const defaultCalicoCIDR = "198.19.0.0/16"

type kubernetes struct {
	Scope      string
	VirtualIPs []string
	CalicoCIDR string
}

func newKubernetes(scope string, virtualIPs []string, calicoCIDR string) base {
	return &kubernetes{
		Scope:      scope,
		VirtualIPs: virtualIPs,
		CalicoCIDR: calicoCIDR,
	}
}

func (k *kubernetes) Charms() []charm {
	return []charm{
		{Name: "ch:kubernetes-control-plane", Channel: "1.33/stable", PlacementScope: "#"},
		{Name: "ch:etcd", Channel: "1.33/stable", PlacementScope: "lxd"},
		{Name: "ch:easyrsa", Channel: "1.33/stable", PlacementScope: "lxd"},
		{Name: "ch:kubeapi-load-balancer", Channel: "1.33/stable", PlacementScope: "lxd"},
		{Name: "ch:calico", Channel: "1.33/stable", Subordinate: true},
		{Name: "ch:containerd", Channel: "1.33/stable", Subordinate: true},
		{Name: "ch:keepalived", Channel: "1.33/stable", Subordinate: true},
	}
}

func (k *kubernetes) Configs() (string, error) {
	configs := map[string]map[string]any{
		"kubernetes-control-plane": {
			"register-with-taints": "",
			"allow-privileged":     "true",
			"loadbalancer-ips":     strings.Join(k.VirtualIPs, " "),
		},
		"kubeapi-load-balancer": {
			"loadbalancer-ips": strings.Join(k.VirtualIPs, " "),
		},
		"calico": {
			"ignore-loose-rpf": "true",
			"cidr":             k.CalicoCIDR,
		},
		"containerd": {
			"gpu_driver": "none",
		},
		"keepalived": {
			"virtual_ip": k.VirtualIPs[0],
		},
	}

	return buildConfigs(k.Scope, configs)
}

func (k *kubernetes) Relations() [][]string {
	relations := [][]string{
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

	return buildRelations(k.Scope, relations)
}

func (k *kubernetes) Tags() []string {
	return []string{
		tag.Kubernetes,
		tag.KubernetesControlPlane,
		tag.KubernetesWorker,
	}
}
