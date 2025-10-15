package core

import (
	"context"
	"fmt"
	"strings"
)

// CharmConfig defines the structure for charm configurations
type CharmConfig struct {
	Charms    []EssentialCharm
	Relations [][]string
}

var (
	kubeCharmConfig = CharmConfig{
		Charms: []EssentialCharm{
			{Name: "ch:kubernetes-control-plane", Channel: "1.33/stable", Machine: true},
			{Name: "ch:etcd", Channel: "1.33/stable", LXD: true},
			{Name: "ch:easyrsa", Channel: "1.33/stable", LXD: true},
			{Name: "ch:kubeapi-load-balancer", Channel: "1.33/stable", LXD: true},
			{Name: "ch:calico", Channel: "1.33/stable", Subordinate: true},
			{Name: "ch:containerd", Channel: "1.33/stable", Subordinate: true},
			{Name: "ch:keepalived", Channel: "1.33/stable", Subordinate: true},
		},
		Relations: [][]string{
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
		},
	}

	cephCharmConfig = CharmConfig{
		Charms: []EssentialCharm{
			{Name: "ch:ceph-mon", Channel: "squid/stable", LXD: true},
			{Name: "ch:ceph-osd", Channel: "squid/stable", Machine: true},
			{Name: "ch:ceph-fs", Channel: "squid/stable", LXD: true},
			{Name: "ch:ceph-radosgw", Channel: "squid/stable", LXD: true},
			{Name: "ch:ceph-nfs", Channel: "squid/stable", LXD: true},
			{Name: "ch:hacluster", Channel: "2.4/stable", Subordinate: true},
		},
		Relations: [][]string{
			{"ceph-mon:client", "ceph-nfs:ceph-client"},
			{"ceph-mon:mds", "ceph-fs:ceph-mds"},
			{"ceph-mon:osd", "ceph-osd:mon"},
			{"ceph-mon:radosgw", "ceph-radosgw:mon"},
			{"hacluster:ha", "ceph-nfs:ha"},
		},
	}

	addonCharmConfig = CharmConfig{
		Charms: []EssentialCharm{
			{Name: "ch:ceph-csi", Channel: "1.33/stable", Subordinate: true},
			{Name: "ch:grafana-agent", Subordinate: true},
		},
		Relations: [][]string{
			{"ceph-csi", "ceph-mon"},
			{"ceph-csi", "kubernetes-control-plane"},
			{"grafana-agent:cos-agent", "ceph-mon:cos-agent"},
			{"grafana-agent:cos-agent", "kubeapi-load-balancer:cos-agent"},
			{"grafana-agent:cos-agent", "kubernetes-control-plane:cos-agent"},
		},
	}
)

// nodeParams encapsulates node creation parameters
type nodeParams struct {
	scope      string
	machineID  string
	prefix     string
	virtualIPs []string
	calicoCIDR string
	osdDevices []string
}

// CreateNode creates a new node with the specified configuration.
func (uc *OrchestratorUseCase) CreateNode(ctx context.Context, scope, machineID, prefix string, userVirtualIPs []string, userCalicoCIDR string, userOSDDevices []string) error {
	params := &nodeParams{
		scope:      scope,
		machineID:  machineID,
		prefix:     prefix,
		virtualIPs: userVirtualIPs,
		calicoCIDR: userCalicoCIDR,
		osdDevices: userOSDDevices,
	}
	return uc.createNodeWithParams(ctx, params)
}

func (uc *OrchestratorUseCase) createNodeWithParams(ctx context.Context, params *nodeParams) error {
	if err := uc.validateMachineStatus(ctx, params.scope, params.machineID); err != nil {
		return err
	}

	osdDevices, err := uc.validateOSDDevices(params.osdDevices)
	if err != nil {
		return err
	}

	kubeVIPs, err := uc.resolveKubeVIPs(ctx, params.machineID, params.prefix, params.virtualIPs)
	if err != nil {
		return err
	}

	calicoCIDR := uc.resolveCalicoCIDR(params.calicoCIDR)

	nfsVIP, err := uc.reserveIP(ctx, params.machineID, fmt.Sprintf("Ceph NFS IP for %s", params.prefix))
	if err != nil {
		return err
	}

	return uc.createNode(ctx, params.scope, params.machineID, params.prefix, kubeVIPs, calicoCIDR, osdDevices, nfsVIP.String())
}

func (uc *OrchestratorUseCase) createNode(ctx context.Context, scope, machineID, prefix, kubeVIPs, calicoCIDR, osdDevices, nfsVIP string) error {
	if err := uc.deployKubernetes(ctx, scope, machineID, prefix, kubeVIPs, calicoCIDR); err != nil {
		return err
	}

	if err := uc.deployCeph(ctx, scope, machineID, prefix, osdDevices, nfsVIP); err != nil {
		return err
	}

	if err := uc.deployAddons(ctx, scope, prefix); err != nil {
		return err
	}

	return nil
}

func (uc *OrchestratorUseCase) deployKubernetes(ctx context.Context, scope, machineID, prefix, kubeVIPs, calicoCIDR string) error {
	configs, err := uc.kubernetesConfigs(prefix, kubeVIPs, calicoCIDR)
	if err != nil {
		return err
	}

	tags := []string{Kubernetes, KubernetesControlPlane, KubernetesWorker}
	if err := uc.createEssential(ctx, scope, machineID, prefix, kubeCharmConfig.Charms, configs, tags); err != nil {
		return err
	}

	endpointList := toEssentialEndpointList(prefix, kubeCharmConfig.Relations)
	return uc.createEssentialRelations(ctx, scope, endpointList)
}

func (uc *OrchestratorUseCase) deployCeph(ctx context.Context, scope, machineID, prefix, osdDevices, nfsVIP string) error {
	configs, err := uc.cephConfigs(prefix, osdDevices, nfsVIP)
	if err != nil {
		return err
	}

	tags := []string{Ceph, CephMon, CephOSD}
	if err := uc.createEssential(ctx, scope, machineID, prefix, cephCharmConfig.Charms, configs, tags); err != nil {
		return err
	}

	endpointList := toEssentialEndpointList(prefix, cephCharmConfig.Relations)
	return uc.createEssentialRelations(ctx, scope, endpointList)
}

func (uc *OrchestratorUseCase) deployAddons(ctx context.Context, scope, prefix string) error {
	configs, err := uc.addonConfigs(prefix)
	if err != nil {
		return err
	}

	if err := uc.createEssential(ctx, scope, "", prefix, addonCharmConfig.Charms, configs, nil); err != nil {
		return err
	}

	if err := uc.createCOS(ctx, scope, prefix); err != nil {
		return err
	}

	endpointList := toEssentialEndpointList(prefix, addonCharmConfig.Relations)
	return uc.createEssentialRelations(ctx, scope, endpointList)
}

func (uc *OrchestratorUseCase) kubernetesConfigs(prefix, vips, cidr string) (map[string]string, error) {
	configs := map[string]map[string]any{
		"kubernetes-control-plane": {
			"allow-privileged": "true",
			"loadbalancer-ips": vips,
		},
		"kubeapi-load-balancer": {
			"loadbalancer-ips": vips,
		},
		"calico": {
			"ignore-loose-rpf": "true",
			"cidr":             cidr,
		},
		"containerd": {
			"gpu_driver": "none",
		},
		"keepalived": {
			"virtual_ip": strings.Fields(vips)[0],
		},
	}
	return toEssentialConfigs(prefix, configs)
}

func (uc *OrchestratorUseCase) cephConfigs(prefix, osdDevices, vip string) (map[string]string, error) {
	configs := map[string]map[string]any{
		"ceph-mon": {
			"monitor-count":       1,
			"expected-osd-count":  1,
			"config-flags":        `{ "global": {"osd_pool_default_size": 1, "osd_pool_default_min_size": 1, "mon_allow_pool_size_one": true} }`,
			"enable-perf-metrics": true,
		},
		"ceph-osd": {
			"osd-devices": osdDevices,
		},
		"ceph-fs": {
			"ceph-osd-replication-count": 2,
		},
		"ceph-nfs": {
			"vip": vip,
		},
		"ceph-radosgw": {
			"ceph-osd-replication-count": 2,
		},
		"hacluster": {
			"cluster_count": 1,
		},
	}
	return toEssentialConfigs(prefix, configs)
}

func (uc *OrchestratorUseCase) addonConfigs(prefix string) (map[string]string, error) {
	configs := map[string]map[string]any{
		"ceph-csi": {
			"default-storage":      "ceph-ext4",
			"cephfs-enable":        "true",
			"provisioner-replicas": 1,
		},
	}
	return toEssentialConfigs(prefix, configs)
}
