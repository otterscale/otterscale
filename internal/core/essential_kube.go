package core

import (
	"context"
	"errors"
	"net"
	"slices"
	"strconv"
	"strings"

	"connectrpc.com/connect"
	"github.com/canonical/gomaasclient/entity"
)

const charmKubernetes = "kubernetes-control-plane"

var (
	kubernetesCharms = []EssentialCharm{
		{Name: "ch:kubernetes-control-plane", Machine: true},
		{Name: "ch:etcd", LXD: true},
		{Name: "ch:easyrsa", LXD: true},
		{Name: "ch:kubeapi-load-balancer", LXD: true},
		{Name: "ch:calico", Subordinate: true},
		{Name: "ch:containerd", Subordinate: true},
		{Name: "ch:keepalived", Subordinate: true},
	}

	kubernetesRelations = [][]string{
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
)

func CreateKubernetes(ctx context.Context, serverRepo ServerRepo, machineRepo MachineRepo, facilityRepo FacilityRepo, uuid, machineID, prefix string, configs map[string]string) error {
	if err := createEssential(ctx, serverRepo, machineRepo, facilityRepo, uuid, machineID, prefix, kubernetesCharms, configs); err != nil {
		return err
	}
	return createEssentialRelations(ctx, facilityRepo, uuid, toEndpointList(prefix, kubernetesRelations))
}

func GetKubernetesCharms() []EssentialCharm {
	return kubernetesCharms
}

func newKubernetesConfigs(prefix, vips, cidr string) (map[string]string, error) {
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
			"virtual_ip": strings.Split(vips, " ")[0],
		},
	}
	return NewCharmConfigs(prefix, configs)
}

func listKuberneteses(ctx context.Context, scopeRepo ScopeRepo, clientRepo ClientRepo, uuid string) ([]Essential, error) {
	return listEssentials(ctx, scopeRepo, clientRepo, charmKubernetes, 1, uuid)
}

func GetAndReserveIP(ctx context.Context, machineRepo MachineRepo, subnetRepo SubnetRepo, ipRangeRepo IPRangeRepo, machineID, comment string) (net.IP, error) {
	machine, err := machineRepo.Get(ctx, machineID)
	if err != nil {
		return nil, err
	}
	links := machine.BootInterface.Links
	if len(links) == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("machine has no network links"))
	}
	subnet := &links[0].Subnet
	ip, err := getFreeIP(ctx, subnetRepo, ipRangeRepo, subnet)
	if err != nil {
		return nil, err
	}
	if _, err := createIPRange(ctx, ipRangeRepo, subnet.ID, ip.String(), ip.String(), comment); err != nil {
		return nil, err
	}
	return ip, nil
}

func getFreeIP(ctx context.Context, subnetRepo SubnetRepo, ipRangeRepo IPRangeRepo, subnet *entity.Subnet) (net.IP, error) {
	skip := []uint32{}
	used, err := getUsedIPs(ctx, subnetRepo, subnet.ID)
	if err != nil {
		return nil, err
	}
	skip = append(skip, used...)

	reserved, err := getReservedIPs(ctx, ipRangeRepo, subnet.CIDR)
	if err != nil {
		return nil, err
	}
	skip = append(skip, reserved...)

	_, ipNet, err := net.ParseCIDR(subnet.CIDR)
	if err != nil {
		return nil, err
	}

	ip := ipToUint32(ipNet.IP)
	mask := ipToUint32(net.IP(ipNet.Mask))
	network := ip & mask
	broadcast := network | ^mask

	next := false // get next to prevent time gap
	for i := network + 1; i < broadcast; i++ {
		if slices.Contains(skip, i) {
			continue
		}
		if next {
			return uint32ToIP(i), nil
		}
		next = true
	}

	return nil, connect.NewError(connect.CodeResourceExhausted, errors.New("no free IP found"))
}

func getUsedIPs(ctx context.Context, subnetRepo SubnetRepo, subnetID int) ([]uint32, error) {
	ips, err := subnetRepo.GetIPAddresses(ctx, subnetID)
	if err != nil {
		return nil, err
	}
	record := []uint32{}
	for i := range ips {
		record = append(record, ipToUint32(ips[i].IP))
	}
	return record, nil
}

func getReservedIPs(ctx context.Context, ipRangeRepo IPRangeRepo, cidr string) ([]uint32, error) {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}
	ipRanges, err := ipRangeRepo.List(ctx)
	if err != nil {
		return nil, err
	}
	record := []uint32{}
	for i := range ipRanges {
		if ipNet.Contains(ipRanges[i].StartIP) && ipNet.Contains(ipRanges[i].EndIP) {
			start := ipToUint32(ipRanges[i].StartIP)
			end := ipToUint32(ipRanges[i].EndIP)
			for i := start; i <= end; i++ {
				record = append(record, i)
			}
		}
	}
	return record, nil
}

func createIPRange(ctx context.Context, ipRangeRepo IPRangeRepo, subnetID int, startIP, endIP, comment string) (*IPRange, error) {
	params := &entity.IPRangeParams{
		Type:    "reserved",
		Subnet:  strconv.Itoa(subnetID),
		StartIP: startIP,
		EndIP:   endIP,
		Comment: comment,
	}
	return ipRangeRepo.Create(ctx, params)
}

func ipToUint32(ip net.IP) uint32 {
	ip = ip.To4()
	return uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3])
}

func uint32ToIP(n uint32) net.IP {
	return net.IP{
		byte(n >> 24), //nolint:mnd // shift
		byte(n >> 16), //nolint:mnd // shift
		byte(n >> 8),  //nolint:mnd // shift
		byte(n),
	}
}
