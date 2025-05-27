package core

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/canonical/gomaasclient/entity/node"
	jujustatus "github.com/juju/juju/core/status"
	"github.com/juju/juju/rpc/params"
	"golang.org/x/sync/errgroup"
)

const (
	charmKubernetes = "kubernetes-control-plane"
	charmCeph       = "ceph-mon"
	charmCephCSI    = "ceph-csi"
)

var (
	kubernetesCharms = []essentialCharm{
		{name: "ch:calico", lxd: true, subordinate: true},
		{name: "ch:containerd", lxd: true, subordinate: true},
		{name: "ch:easyrsa", lxd: true},
		{name: "ch:etcd", lxd: true},
		{name: "ch:keepalived", lxd: true, subordinate: true},
		{name: "ch:kubeapi-load-balancer", lxd: true},
		{name: "ch:kubernetes-control-plane", lxd: true},
		{name: "ch:kubernetes-worker"},
	}

	kubernetesRelations = [][]string{
		{"calico:cni", "kubernetes-control-plane:cni"},
		{"calico:cni", "kubernetes-worker:cni"},
		{"calico:etcd", "etcd:db"},
		{"easyrsa:client", "etcd:certificates"},
		{"easyrsa:client", "kubernetes-control-plane:certificates"},
		{"easyrsa:client", "kubernetes-worker:certificates"},
		{"easyrsa:client", "kubeapi-load-balancer:certificates"},
		{"etcd:db", "kubernetes-control-plane:etcd"},
		{"kubernetes-control-plane:kube-control", "kubernetes-worker:kube-control"},
		{"kubernetes-control-plane:loadbalancer-external", "kubeapi-load-balancer:lb-consumers"},
		{"kubernetes-control-plane:loadbalancer-internal", "kubeapi-load-balancer:lb-consumers"},
		{"keepalived:juju-info", "kubeapi-load-balancer:juju-info"},
		{"keepalived:website", "kubeapi-load-balancer:apiserver"},
		{"containerd:containerd", "kubernetes-control-plane:container-runtime"},
		{"containerd:containerd", "kubernetes-worker:container-runtime"},
	}

	cephCharms = []essentialCharm{
		{name: "ch:ceph-fs", lxd: true},
		{name: "ch:ceph-mon", lxd: true},
		{name: "ch:ceph-osd", lxd: false},
	}

	cephRelations = [][]string{
		{"ceph-fs:ceph-mds", "ceph-mon:mds"},
		{"ceph-osd:mon", "ceph-mon:osd"},
	}

	cephCSICharms = []essentialCharm{
		{name: "ch:ceph-csi", lxd: true},
	}

	cephCSIRelations = [][]string{
		{"ceph-csi", "ceph-mon"},
		{"ceph-csi", "kubernetes-control-plane"},
	}
)

type essentialCharm struct {
	name        string
	lxd         bool
	subordinate bool
}

type Essential struct {
	Type      int32
	Name      string
	ScopeUUID string
	ScopeName string
	Units     []EssentialUnit
}

type EssentialUnit struct {
	Type      int32
	Name      string
	MachineID string
}

type EssentialStatus struct {
	Level   int32
	Message string
	Details string
}

type EssentialUseCase struct {
	machine MachineRepo
	scope   ScopeRepo
	client  ClientRepo
}

func NewEssentialUseCase(machine MachineRepo, scope ScopeRepo, client ClientRepo) *EssentialUseCase {
	return &EssentialUseCase{
		machine: machine,
		scope:   scope,
		client:  client,
	}
}

func (uc *EssentialUseCase) IsMachineDeployed(ctx context.Context, uuid string) (string, bool, error) {
	machines, err := uc.machine.List(ctx)
	if err != nil {
		return "", false, err
	}
	scopeMachines := []Machine{}
	for i := range machines {
		scopeUUID, _ := getJujuModelUUID(machines[i].WorkloadAnnotations)
		if scopeUUID == uuid {
			scopeMachines = append(scopeMachines, machines[i])
		}
	}
	for i := range scopeMachines {
		if scopeMachines[i].Status == node.StatusDeployed {
			return "", true, err
		}
	}
	return uc.getMachineStatusMessage(machines), false, nil
}

func (uc *EssentialUseCase) ListStatuses(ctx context.Context, uuid string) ([]EssentialStatus, error) {
	s, err := uc.client.Status(ctx, uuid, []string{"application", "*"})
	if err != nil {
		return nil, err
	}

	charms := []essentialCharm{}
	charms = append(charms, kubernetesCharms...)
	charms = append(charms, cephCharms...)
	charms = append(charms, cephCSICharms...)

	statuses := []EssentialStatus{}
	for name := range s.Applications {
		ok := isEssentialCharm(s.Applications[name], charms)
		if !ok {
			continue
		}

		status := s.Applications[name].Status
		level := int32(0) // info
		switch status.Status {
		case jujustatus.Maintenance.String():
			level = 1 // low
		case jujustatus.Unknown.String(), jujustatus.Waiting.String():
			level = 2 // medium
		case jujustatus.Blocked.String():
			level = 3 // high
		case jujustatus.Unset.String(), jujustatus.Terminated.String(), jujustatus.Active.String():
			continue
		}

		statuses = append(statuses, EssentialStatus{
			Level:   level,
			Message: fmt.Sprintf("[%s] %s", status.Status, name),
			Details: status.Info,
		})
	}
	return statuses, nil
}

func (uc *EssentialUseCase) ListEssentials(ctx context.Context, uuid string) ([]Essential, error) {
	eg, ctx := errgroup.WithContext(ctx)
	result := make([][]Essential, 2)
	eg.Go(func() error {
		v, err := listKuberneteses(ctx, uc.scope, uc.client)
		if err == nil {
			result[0] = v
		}
		return err
	})
	eg.Go(func() error {
		v, err := listCephs(ctx, uc.scope, uc.client)
		if err == nil {
			result[1] = v
		}
		return err
	})
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return append(result[0], result[1]...), nil
}

func (uc *EssentialUseCase) getMachineStatusMessage(machines []Machine) string {
	statuses := []node.Status{
		node.StatusDefault,
		node.StatusCommissioning,
		node.StatusFailedCommissioning,
		node.StatusTesting,
		node.StatusFailedTesting,
		node.StatusDeploying,
		node.StatusReady,
	}
	statusMessages := []string{
		"",
		"commissioning",
		"failed to commission",
		"testing",
		"failed to test",
		"deploying",
		"unknown",
	}
	statusIndex := 0
	message := ""
	for i := range machines {
		currentIndex := 0
		for j := range statuses {
			if machines[i].Status == statuses[j] {
				currentIndex = j
				break
			}
		}
		if statusIndex < currentIndex {
			statusIndex = currentIndex
			message = fmt.Sprintf("machine %q is %s", machines[i].FQDN, statusMessages[statusIndex])
		}
	}
	return message
}

// ch:amd64/kubernetes-control-plane-567 -> kubernetes-control-plane
func formatAppCharm(name string) (string, bool) {
	t := strings.Split(name, "/")
	if len(t) < 1 {
		return "", false
	}
	u := strings.Split(t[1], "-")
	_, err := strconv.Atoi(u[len(u)-1])
	if err != nil {
		return "", false
	}
	return strings.Join(u[:len(u)-1], "-"), true
}

// ch:kubernetes-control-plane -> kubernetes-control-plane
func formatEssentialCharm(name string) string {
	return strings.TrimPrefix(name, "ch:")
}

func isEssentialCharm(appStatus params.ApplicationStatus, charms []essentialCharm) bool {
	appCharm, ok := formatAppCharm(appStatus.Charm)
	if !ok {
		return false
	}
	for _, charm := range charms {
		essCharm := formatEssentialCharm(charm.name)
		if appCharm == essCharm {
			return true
		}
	}
	return false
}

// func (uc *NexusService) VerifyEnvironment(ctx context.Context, scopeUUID string) ([]model.Error, error) {
// 	funcs := []func(context.Context, string) ([]model.Error, error){
// 		s.isCephExists,
// 		s.isKubernetesExists,
// 		s.isDeployedMachineExists,
// 		s.listCephStatusMessage,
// 		s.listCephCSIStatusMessage,
// 		s.listKubernetesStatusMessage,
// 	}

// 	eg, ctx := errgroup.WithContext(ctx)
// 	result := make([][]model.Error, len(funcs))
// 	for i := range funcs {
// 		i := i // fixed on go 1.22
// 		eg.Go(func() error {
// 			es, err := funcs[i](ctx, scopeUUID)
// 			if err == nil && es != nil {
// 				result[i] = es
// 			}
// 			return err
// 		})
// 	}
// 	if err := eg.Wait(); err != nil {
// 		return nil, err
// 	}
// 	errs := []model.Error{}
// 	for i := range result {
// 		errs = append(errs, result[i]...)
// 	}
// 	slices.SortFunc(errs, func(e1, e2 model.Error) int {
// 		return strings.Compare(e1.Code, e2.Code)
// 	})
// 	return slices.DeleteFunc(errs, func(e model.Error) bool { return e.Code == "" }), nil
// }

// func (uc *NexusService) ListCephes(ctx context.Context, uuid string) ([]model.FacilityInfo, error) {
// 	return s.listGeneralFacilities(ctx, uuid, charmNameCeph)
// }

// func (uc *NexusService) CreateCeph(ctx context.Context, uuid, machineID, prefix string, userOSDDevices []string, development bool) (*model.FacilityInfo, error) {
// 	osdDevices := strings.Join(userOSDDevices, " ")
// 	if osdDevices == "" {
// 		return nil, status.Error(codes.InvalidArgument, "no OSD devices provided")
// 	}
// 	configs, err := getCephConfigs(prefix, osdDevices, development)
// 	if err != nil {
// 		return nil, err
// 	}
// 	fi, err := s.createGeneralFacility(ctx, uuid, machineID, prefix, charmNameCeph, cephFacilityList, configs)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if err := s.createGeneralRelations(ctx, uuid, toEndpointList(prefix, cephRelationList)); err != nil {
// 		return nil, err
// 	}
// 	return fi, nil
// }

// func (uc *NexusService) AddCephUnits(ctx context.Context, uuid, general string, number int, machineIDs []string) error {
// 	return s.addGeneralFacilityUnits(ctx, uuid, general, number, machineIDs, cephFacilityList)
// }

// func (uc *NexusService) ListKuberneteses(ctx context.Context, uuid string) ([]model.FacilityInfo, error) {
// 	return s.listGeneralFacilities(ctx, uuid, charmNameKubernetes)
// }

// func (uc *NexusService) getAndReserveIP(ctx context.Context, machineID, comment string) (net.IP, error) {
// 	machine, err := s.machine.Get(ctx, machineID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	links := machine.BootInterface.Links
// 	if len(links) == 0 {
// 		return nil, status.Error(codes.InvalidArgument, "machine has no network links")
// 	}
// 	subnet := &links[0].Subnet
// 	ip, err := s.getFreeIP(ctx, subnet)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if _, err := s.CreateIPRange(ctx, subnet.ID, ip.String(), ip.String(), comment); err != nil {
// 		return nil, err
// 	}
// 	return ip, nil
// }

// func (uc *NexusService) CreateKubernetes(ctx context.Context, uuid, machineID, prefix string, userVirtualIPs []string, userCalicoCIDR string) (*model.FacilityInfo, error) {
// 	vips := strings.Join(userVirtualIPs, " ")
// 	if vips == "" {
// 		ip, err := s.getAndReserveIP(ctx, machineID, fmt.Sprintf("kubernetes load balancer IP for %s", prefix))
// 		if err != nil {
// 			return nil, err
// 		}
// 		vips = ip.String()
// 	}

// 	cidr := userCalicoCIDR
// 	if cidr == "" {
// 		cidr = "198.19.0.0/16"
// 	}

// 	configs, err := getKubernetesConfigs(prefix, vips, cidr)
// 	if err != nil {
// 		return nil, err
// 	}
// 	fi, err := s.createGeneralFacility(ctx, uuid, machineID, prefix, charmNameKubernetes, kubernetesFacilityList, configs)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if err := s.createGeneralRelations(ctx, uuid, toEndpointList(prefix, kubernetesRelationList)); err != nil {
// 		return nil, err
// 	}
// 	return fi, nil
// }

// func (uc *NexusService) AddKubernetesUnits(ctx context.Context, uuid, general string, number int, machineIDs []string, force bool) error {
// 	if !force {
// 		st, err := s.client.Status(ctx, uuid, []string{"application", general})
// 		if err != nil {
// 			return err
// 		}
// 		app, ok := st.Applications[general]
// 		if !ok {
// 			return status.Errorf(codes.NotFound, "kubernetes facility %q not found", general)
// 		}
// 		if len(app.Units) > 3 {
// 			return status.Errorf(codes.InvalidArgument, "cannot add more than 3 Kubernetes worker units without force flag")
// 		}
// 	}
// 	return s.addGeneralFacilityUnits(ctx, uuid, general, number, machineIDs, kubernetesFacilityList)
// }

// func (uc *NexusService) SetCephCSI(ctx context.Context, kubernetes, ceph *model.FacilityInfo, prefix string, development bool) error {
// 	if kubernetes.ScopeUUID != ceph.ScopeUUID {
// 		return status.Error(codes.Unimplemented, "cross-model integration between facilities is not yet supported")
// 	}
// 	configs, err := getCephCSIConfigs(prefix, development)
// 	if err != nil {
// 		return err
// 	}
// 	if _, err := s.createGeneralFacility(ctx, kubernetes.ScopeUUID, "", prefix, charmNameCephCSI, cephCSIFacilityList, configs); err != nil {
// 		return err
// 	}
// 	return s.createGeneralRelations(ctx, kubernetes.ScopeUUID, toCephCSIEndpointList(kubernetes, ceph, prefix))
// }

// func (uc *NexusService) createGeneralFacility(ctx context.Context, uuid, machineID, prefix, general string, facilityList []generalFacility, configs map[string]string) (*model.FacilityInfo, error) {
// 	var directive string

// 	if machineID != "" {
// 		m, err := s.machine.Get(ctx, machineID)
// 		if err != nil {
// 			return nil, err
// 		}
// 		if m.Status != node.StatusDeployed {
// 			return nil, status.Error(codes.InvalidArgument, "machine status is not deployed")
// 		}

// 		machineID, err := getJujuMachineID(m.WorkloadAnnotations)
// 		if err != nil {
// 			return nil, err
// 		}
// 		directive = machineID
// 	}

// 	base, err := s.imageBase(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var facilityName string
// 	eg, ctx := errgroup.WithContext(ctx)
// 	for _, facility := range facilityList {
// 		facility := facility // fixed on go 1.22
// 		eg.Go(func() error {
// 			name := toGeneralFacilityName(prefix, facility.charmName)
// 			config := configs[facility.charmName]
// 			placements := []instance.Placement{}
// 			if directive != "" {
// 				placements = append(placements, instance.Placement{
// 					Scope:     toPlacementScope(facility.lxd),
// 					Directive: directive,
// 				})
// 			}
// 			_, err := s.facility.Create(ctx, uuid, name, config, facility.charmName, "", 0, 1, base, placements, nil, true)
// 			return err
// 		})
// 		if facility.charmName == "ch:"+general {
// 			facilityName = toGeneralFacilityName(prefix, facility.charmName)
// 		}
// 	}

// 	if err := eg.Wait(); err != nil {
// 		return nil, err
// 	}

// 	scopeName, err := s.getScopeName(ctx, uuid)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &model.FacilityInfo{
// 		ScopeUUID:    uuid,
// 		ScopeName:    scopeName,
// 		FacilityName: facilityName,
// 	}, nil
// }

// func (uc *NexusService) addGeneralFacilityUnits(ctx context.Context, uuid, general string, number int, machineIDs []string, facilityList []generalFacility) error {
// 	directives := []string{}
// 	for _, machineID := range machineIDs {
// 		m, err := s.machine.Get(ctx, machineID)
// 		if err != nil {
// 			return err
// 		}
// 		if m.Status != node.StatusDeployed {
// 			return status.Errorf(codes.InvalidArgument, "machine %q status is not deployed", machineID)
// 		}
// 		directive, err := getJujuMachineID(m.WorkloadAnnotations)
// 		if err != nil {
// 			return err
// 		}
// 		directives = append(directives, directive)
// 	}

// 	slices.Sort(directives)
// 	directives = slices.Compact(directives)

// 	if len(directives) != number {
// 		return status.Error(codes.InvalidArgument, "number of machines does not match requested number of units")
// 	}

// 	prefix := toGeneralFacilityPrefix(general)

// 	eg, ctx := errgroup.WithContext(ctx)
// 	for _, facility := range facilityList {
// 		facility := facility // fixed on go 1.22
// 		if facility.subordinate {
// 			continue
// 		}
// 		eg.Go(func() error {
// 			name := toGeneralFacilityName(prefix, facility.charmName)
// 			lxd := facility.lxd
// 			placements := make([]instance.Placement, len(directives))
// 			for i, directive := range directives {
// 				placements[i] = instance.Placement{
// 					Scope:     toPlacementScope(lxd),
// 					Directive: directive,
// 				}
// 			}
// 			_, err := s.facility.AddUnits(ctx, uuid, name, number, placements)
// 			return err
// 		})
// 	}

// 	return eg.Wait()
// }

// func (uc *NexusService) getScopeName(ctx context.Context, uuid string) (string, error) {
// 	scopes, err := s.scope.List(ctx)
// 	if err != nil {
// 		return "", err
// 	}

// 	for i := range scopes {
// 		if scopes[i].UUID == uuid {
// 			return scopes[i].Name, nil
// 		}
// 	}

// 	return "", nil
// }

// func toEndpointList(prefix string, relationList [][]string) [][]string {
// 	endpointList := [][]string{}
// 	for _, relations := range relationList {
// 		endpoints := []string{}
// 		for _, relation := range relations {
// 			endpoints = append(endpoints, toGeneralFacilityName(prefix, relation))
// 		}
// 		endpointList = append(endpointList, endpoints)
// 	}
// 	return endpointList
// }

// func toCephCSIEndpointList(kubernetes, ceph *model.FacilityInfo, prefix string) [][]string {
// 	endpointList := [][]string{}
// 	for _, relations := range cephCSIRelationList {
// 		endpoints := []string{}
// 		for _, relation := range relations {
// 			if relation == charmNameCephCSI {
// 				endpoints = append(endpoints, toGeneralFacilityName(prefix, relation))
// 			} else if relation == charmNameKubernetes {
// 				endpoints = append(endpoints, kubernetes.FacilityName)
// 			} else if relation == charmNameCeph {
// 				endpoints = append(endpoints, ceph.FacilityName)
// 			}
// 		}
// 		endpointList = append(endpointList, endpoints)
// 	}
// 	return endpointList
// }

// func (uc *NexusService) createGeneralRelations(ctx context.Context, uuid string, endpointList [][]string) error {
// 	eg, ctx := errgroup.WithContext(ctx)
// 	for _, endpoints := range endpointList {
// 		endpoints := endpoints // fixed on go 1.22
// 		eg.Go(func() error {
// 			_, err := s.facility.CreateRelation(ctx, uuid, endpoints)
// 			return err
// 		})
// 	}
// 	return eg.Wait()
// }

// func (uc *NexusService) isCephExists(ctx context.Context, scopeUUID string) ([]model.Error, error) {
// 	cephes, err := s.ListCephes(ctx, scopeUUID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if len(cephes) == 0 {
// 		return []model.Error{model.ErrCephNotFound}, nil
// 	}
// 	return nil, nil
// }

// func (uc *NexusService) isKubernetesExists(ctx context.Context, scopeUUID string) ([]model.Error, error) {
// 	kuberneteses, err := s.ListKuberneteses(ctx, scopeUUID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if len(kuberneteses) == 0 {
// 		return []model.Error{model.ErrKubernetesNotFound}, nil
// 	}
// 	return nil, nil
// }

// func (uc *NexusService) listCephStatusMessage(ctx context.Context, scopeUUID string) ([]model.Error, error) {
// 	return s.listStatusMessage(ctx, scopeUUID, cephFacilityList, model.ErrCephStatusMessageCode)
// }

// func (uc *NexusService) listCephCSIStatusMessage(ctx context.Context, scopeUUID string) ([]model.Error, error) {
// 	return s.listStatusMessage(ctx, scopeUUID, cephCSIFacilityList, model.ErrCephCSIStatusMessageCode)
// }

// func (uc *NexusService) listKubernetesStatusMessage(ctx context.Context, scopeUUID string) ([]model.Error, error) {
// 	return s.listStatusMessage(ctx, scopeUUID, kubernetesFacilityList, model.ErrKubernetesStatusMessageCode)
// }

// func (uc *NexusService) getReservedIPs(ctx context.Context, cidr string) ([]uint32, error) {
// 	_, ipNet, err := net.ParseCIDR(cidr)
// 	if err != nil {
// 		return nil, err
// 	}
// 	ipRanges, err := s.ipRange.List(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	record := []uint32{}
// 	for i := range ipRanges {
// 		if ipNet.Contains(ipRanges[i].StartIP) && ipNet.Contains(ipRanges[i].EndIP) {
// 			start := ipToUint32(ipRanges[i].StartIP)
// 			end := ipToUint32(ipRanges[i].EndIP)
// 			for i := start; i <= end; i++ {
// 				record = append(record, i)
// 			}
// 		}
// 	}
// 	return record, nil
// }

// func (uc *NexusService) getUsedIPs(ctx context.Context, subnetID int) ([]uint32, error) {
// 	ipas, err := s.subnet.GetIPAddresses(ctx, subnetID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	record := []uint32{}
// 	for i := range ipas {
// 		record = append(record, ipToUint32(ipas[i].IP))
// 	}
// 	return record, nil
// }

// func (uc *NexusService) getFreeIP(ctx context.Context, subnet *entity.Subnet) (net.IP, error) {
// 	skip := []uint32{}
// 	used, err := s.getUsedIPs(ctx, subnet.ID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	skip = append(skip, used...)

// 	reserved, err := s.getReservedIPs(ctx, subnet.CIDR)
// 	if err != nil {
// 		return nil, err
// 	}
// 	skip = append(skip, reserved...)

// 	_, ipNet, err := net.ParseCIDR(subnet.CIDR)
// 	if err != nil {
// 		return nil, err
// 	}

// 	ip := ipToUint32(ipNet.IP)
// 	mask := ipToUint32(net.IP(ipNet.Mask))
// 	network := ip & mask
// 	broadcast := network | ^mask

// 	next := false // get next to prevent time gap
// 	for i := network + 1; i < broadcast; i++ {
// 		if slices.Contains(skip, i) {
// 			continue
// 		}
// 		if next {
// 			return uint32ToIP(i), nil
// 		}
// 		next = true
// 	}

// 	return nil, status.Errorf(codes.ResourceExhausted, "no free IP found")
// }

// func toPlacementScope(lxd bool) string {
// 	if lxd {
// 		return "lxd"
// 	}
// 	return instance.MachineScope
// }

// func toGeneralFacilityName(prefix, charmName string) string {
// 	if strings.HasPrefix(charmName, "ch:") {
// 		return prefix + "-" + strings.Split(charmName, ":")[1]
// 	}
// 	return prefix + "-" + charmName
// }

// func toGeneralFacilityPrefix(general string) string {
// 	return strings.Split(general, "-")[0]
// }

// func getKubernetesConfigs(prefix, vips, cidr string) (map[string]string, error) {
// 	configs := map[string]map[string]any{
// 		"kubernetes-control-plane": {
// 			"allow-privileged": "true",
// 			"extra_sans":       vips,
// 			"loadbalancer-ips": vips,
// 		},
// 		"kubeapi-load-balancer": {
// 			"extra_sans":       vips,
// 			"loadbalancer-ips": vips,
// 		},
// 		"calico": {
// 			"ignore-loose-rpf": "true",
// 			"cidr":             cidr,
// 		},
// 		"containerd": {
// 			"gpu_driver": "none",
// 		},
// 		"keepalived": {
// 			"virtual_ip": strings.Split(vips, " ")[0],
// 		},
// 	}

// 	result := make(map[string]string)
// 	for name, config := range configs {
// 		key := toGeneralFacilityName(prefix, name)
// 		yamlData, err := yaml.Marshal(map[string]any{key: config})
// 		if err != nil {
// 			return nil, err
// 		}
// 		result["ch:"+name] = string(yamlData)
// 	}

// 	return result, nil
// }

// func getCephConfigs(prefix, osdDevices string, development bool) (map[string]string, error) {
// 	count := 2
// 	if development {
// 		count = 1
// 	}
// 	configs := map[string]map[string]any{
// 		"ceph-mon": {
// 			"monitor-count":      count,
// 			"expected-osd-count": count,
// 		},
// 		"ceph-osd": {
// 			"osd-devices": osdDevices,
// 		},
// 		"ceph-fs": {
// 			"ceph-osd-replication-count": count,
// 		},
// 	}
// 	if development {
// 		configs["ceph-mon"]["config-flags"] = `{ "global": {"osd_pool_default_size": 1, "osd_pool_default_min_size": 1, "mon_allow_pool_size_one": true} }`
// 	}

// 	result := make(map[string]string)
// 	for name, config := range configs {
// 		key := toGeneralFacilityName(prefix, name)
// 		yamlData, err := yaml.Marshal(map[string]any{key: config})
// 		if err != nil {
// 			return nil, err
// 		}
// 		result["ch:"+name] = string(yamlData)
// 	}

// 	return result, nil
// }

// func getCephCSIConfigs(prefix string, development bool) (map[string]string, error) {
// 	count := 3
// 	if development {
// 		count = 1
// 	}
// 	configs := map[string]map[string]any{
// 		"ceph-csi": {
// 			"default-storage":       "ceph-ext4",
// 			"provisioner-replicas": count,
// 		},
// 	}

// 	result := make(map[string]string)
// 	for name, config := range configs {
// 		key := toGeneralFacilityName(prefix, name)
// 		yamlData, err := yaml.Marshal(map[string]any{key: config})
// 		if err != nil {
// 			return nil, err
// 		}
// 		result["ch:"+name] = string(yamlData)
// 	}

// 	return result, nil
// }

// func ipToUint32(ip net.IP) uint32 {
// 	ip = ip.To4()
// 	return uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3])
// }

// func uint32ToIP(n uint32) net.IP {
// 	return net.IP{
// 		byte(n >> 24),
// 		byte(n >> 16),
// 		byte(n >> 8),
// 		byte(n),
// 	}
// }

func listKuberneteses(ctx context.Context, scopeRepo ScopeRepo, clientRepo ClientRepo) ([]Essential, error) {
	return listEssentials(ctx, scopeRepo, clientRepo, charmKubernetes, 1)
}

func listCephs(ctx context.Context, scopeRepo ScopeRepo, clientRepo ClientRepo) ([]Essential, error) {
	return listEssentials(ctx, scopeRepo, clientRepo, charmCeph, 2)
}

func listEssentials(ctx context.Context, scopeRepo ScopeRepo, clientRepo ClientRepo, charmName string, essentialType int32) ([]Essential, error) {
	scopes, err := scopeRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	eg, ctx := errgroup.WithContext(ctx)
	result := make([]Essential, len(scopes))
	for i := range scopes {
		i := i // fixed on go 1.22
		eg.Go(func() error {
			s, err := clientRepo.Status(ctx, scopes[i].UUID, []string{"application", "*"})
			if err != nil {
				return err
			}
			for name := range s.Applications {
				if !strings.Contains(s.Applications[name].Charm, charmCeph) {
					continue
				}
				result[i] = Essential{
					Type:      essentialType,
					Name:      name,
					ScopeUUID: scopes[i].UUID,
					ScopeName: scopes[i].Name,
				}
				break
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return result, nil
}
