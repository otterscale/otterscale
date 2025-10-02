package core

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"
	"time"

	"connectrpc.com/connect"
	"github.com/canonical/gomaasclient/entity/node"
	"github.com/juju/juju/core/instance"
	jujustatus "github.com/juju/juju/core/status"
	"github.com/juju/juju/rpc/params"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/timestamppb"
	jujuyaml "gopkg.in/yaml.v2"
	"k8s.io/client-go/rest"

	pb "github.com/otterscale/otterscale/api/essential/v1"
	"github.com/otterscale/otterscale/internal/config"
)

type Essential struct {
	Type      int32
	Name      string
	ScopeUUID string
	ScopeName string
	Units     []EssentialUnit
}

type EssentialUnit struct {
	Name      string
	Directive string
}

type EssentialStatus struct {
	Level   int32
	Message string
	Details string
}

type EssentialCharm struct {
	Name        string
	Channel     string
	LXD         bool
	Machine     bool
	Subordinate bool
}

type EssentialVGpuPodInfo struct {
	Pod       *Pod
	ModelName string
	VGpuInfos []EssentialVGpuAllocation
}

type EssentialVGpuAllocation struct {
	GpuUUID       string
	Vendor        string
	VramMib       string
	VcoresPercent string
	BindTime      string
	BindPhase     string
}

type EssentialUseCase struct {
	conf           *config.Config
	kubeCore       KubeCoreRepo
	kubeApps       KubeAppsRepo
	action         ActionRepo
	scope          ScopeRepo
	facility       FacilityRepo
	facilityOffers FacilityOffersRepo
	machine        MachineRepo
	subnet         SubnetRepo
	ipRange        IPRangeRepo
	server         ServerRepo
	client         ClientRepo
	tag            TagRepo
}

func NewEssentialUseCase(conf *config.Config, kubeCore KubeCoreRepo, kubeApps KubeAppsRepo, action ActionRepo, scope ScopeRepo, facility FacilityRepo, facilityOffers FacilityOffersRepo, machine MachineRepo, subnet SubnetRepo, ipRange IPRangeRepo, server ServerRepo, client ClientRepo, tag TagRepo) *EssentialUseCase {
	return &EssentialUseCase{
		conf:           conf,
		kubeCore:       kubeCore,
		kubeApps:       kubeApps,
		action:         action,
		scope:          scope,
		facility:       facility,
		facilityOffers: facilityOffers,
		machine:        machine,
		subnet:         subnet,
		ipRange:        ipRange,
		server:         server,
		client:         client,
		tag:            tag,
	}
}

func (uc *EssentialUseCase) IsMachineDeployed(ctx context.Context, uuid string) (message string, ok bool, err error) {
	machines, err := uc.machine.List(ctx)
	if err != nil {
		return "", false, err
	}
	scopeMachines := []Machine{}
	for i := range machines {
		scopeUUID, err := getJujuModelUUID(machines[i].WorkloadAnnotations)
		if err != nil {
			continue
		}
		if scopeUUID == uuid {
			scopeMachines = append(scopeMachines, machines[i])
		}
	}
	for i := range scopeMachines {
		if scopeMachines[i].Status == node.StatusDeployed {
			return "", true, err
		}
	}
	return uc.getMachineStatusMessage(scopeMachines), false, nil
}

func (uc *EssentialUseCase) ListStatuses(ctx context.Context, uuid string) ([]EssentialStatus, error) {
	s, err := uc.client.Status(ctx, uuid, []string{"application", "*"})
	if err != nil {
		return nil, err
	}

	charms := []EssentialCharm{}
	charms = append(charms, kubernetesCharms...)
	charms = append(charms, cephCharms...)
	charms = append(charms, commonCharms...)

	statuses := []EssentialStatus{}
	for name := range s.Applications {
		ok := isEssentialCharm(s.Applications, name, charms)
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

func (uc *EssentialUseCase) ListEssentials(ctx context.Context, esType int32, uuid string) ([]Essential, error) {
	eg, egctx := errgroup.WithContext(ctx)
	result := make([][]Essential, 2)
	if esType == 0 || esType == 1 {
		eg.Go(func() error {
			v, err := listKuberneteses(egctx, uc.scope, uc.client, uuid)
			if err == nil {
				result[0] = v
			}
			return err
		})
	}
	if esType == 0 || esType == 2 {
		eg.Go(func() error {
			v, err := listCephs(egctx, uc.scope, uc.client, uuid)
			if err == nil {
				result[1] = v
			}
			return err
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return append(result[0], result[1]...), nil
}

func (uc *EssentialUseCase) CreateSingleNode(ctx context.Context, uuid, machineID, prefix string, userVirtualIPs []string, userCalicoCIDR string, userOSDDevices []string) error {
	// validate
	if err := uc.validateMachineStatus(ctx, uuid, machineID); err != nil {
		return err
	}

	// check
	osdDevices := strings.Join(userOSDDevices, " ")
	if osdDevices == "" {
		return connect.NewError(connect.CodeInvalidArgument, errors.New("no OSD devices provided"))
	}

	// default
	kubeVIPs := strings.Join(userVirtualIPs, " ")
	if kubeVIPs == "" {
		ip, err := GetAndReserveIP(ctx, uc.machine, uc.subnet, uc.ipRange, machineID, fmt.Sprintf("Kubernetes Load Balancer IP for %s", prefix))
		if err != nil {
			return err
		}
		kubeVIPs = ip.String()
	}

	cidr := userCalicoCIDR
	if cidr == "" {
		cidr = "198.19.0.0/16"
	}

	// config
	kubeConfigs, err := newKubernetesConfigs(prefix, kubeVIPs, cidr)
	if err != nil {
		return err
	}

	nfsVIP, err := GetAndReserveIP(ctx, uc.machine, uc.subnet, uc.ipRange, machineID, fmt.Sprintf("Ceph NFS IP for %s", prefix))
	if err != nil {
		return err
	}

	cephConfigs, err := newCephConfigs(prefix, osdDevices, nfsVIP.String())
	if err != nil {
		return err
	}

	commonConfigs, err := newCommonConfigs(prefix)
	if err != nil {
		return err
	}

	// create
	if err := CreateCeph(ctx, uc.server, uc.machine, uc.facility, uc.tag, uuid, machineID, prefix, cephConfigs); err != nil {
		return err
	}
	if err := CreateKubernetes(ctx, uc.server, uc.machine, uc.facility, uc.tag, uuid, machineID, prefix, kubeConfigs); err != nil {
		return err
	}
	if err := CreateCommon(ctx, uc.server, uc.machine, uc.facility, uc.facilityOffers, uc.conf, uuid, prefix, commonConfigs); err != nil {
		return err
	}
	return nil
}

func (uc *EssentialUseCase) ListKubernetesNodeLabels(ctx context.Context, uuid, facility, hostname string, all bool) (map[string]string, error) {
	const labelMinParts = 2 // Labels must have at least 2 parts (domain/key)

	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}

	node, err := uc.kubeCore.GetNode(ctx, config, hostname)
	if err != nil {
		return nil, err
	}

	if !all {
		maps.DeleteFunc(node.Labels, func(k, _ string) bool {
			parts := strings.Split(k, "/")
			return len(parts) < labelMinParts || !strings.HasSuffix(parts[0], LabelDomain)
		})
	}
	return node.Labels, nil
}

func (uc *EssentialUseCase) UpdateKubernetesNodeLabels(ctx context.Context, uuid, facility, hostname string, labels map[string]string) (map[string]string, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	node, err := uc.kubeCore.GetNode(ctx, config, hostname)
	if err != nil {
		return nil, err
	}
	if node.Labels == nil {
		node.Labels = map[string]string{}
	}
	for k, v := range labels {
		if v == "" {
			delete(node.Labels, k)
		} else {
			node.Labels[k] = v
		}
	}
	updatedNode, err := uc.kubeCore.UpdateNode(ctx, config, node)
	if err != nil {
		return nil, err
	}
	return updatedNode.Labels, nil
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
	message := "machine not found"
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

func (uc *EssentialUseCase) validateMachineStatus(ctx context.Context, uuid, machineID string) error {
	// maas
	machine, err := uc.machine.Get(ctx, machineID)
	if err != nil {
		return err
	}
	if machine.Status != node.StatusDeployed {
		return connect.NewError(connect.CodeInvalidArgument, errors.New("machine is not deployed"))
	}

	// juju
	id, err := getJujuMachineID(machine.WorkloadAnnotations)
	if err != nil {
		return err
	}
	status, err := uc.client.Status(ctx, uuid, []string{"machine", id})
	if err != nil {
		return err
	}
	m, ok := status.Machines[id]
	if !ok {
		return connect.NewError(connect.CodeInvalidArgument, errors.New("machine is not found"))
	}
	if m.AgentStatus.Status != jujustatus.Started.String() {
		return connect.NewError(connect.CodeInvalidArgument, errors.New("machine is not started"))
	}
	return nil
}

func NewCharmConfigs(prefix string, configs map[string]map[string]any) (map[string]string, error) {
	result := make(map[string]string)
	for name, config := range configs {
		key := toEssentialName(prefix, name)
		value, err := jujuyaml.Marshal(map[string]any{key: config})
		if err != nil {
			return nil, err
		}
		result["ch:"+name] = string(value)
	}
	return result, nil
}

// ch:amd64/kubernetes-control-plane-567 -> kubernetes-control-plane
func formatAppCharm(name string) (string, bool) {
	t := strings.Split(name, "/")
	if len(t) < 2 {
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

func isEssentialCharm(statusMap map[string]params.ApplicationStatus, name string, charms []EssentialCharm) bool {
	appCharm, ok := formatAppCharm(statusMap[name].Charm)
	if !ok {
		return false
	}
	for _, charm := range charms {
		essCharm := formatEssentialCharm(charm.Name)
		if appCharm == essCharm {
			return true
		}
	}
	return false
}

func listEssentials(ctx context.Context, scopeRepo ScopeRepo, clientRepo ClientRepo, charmName string, essentialType int32, scopeUUID string) ([]Essential, error) {
	scopes, err := scopeRepo.List(ctx)
	if err != nil {
		return nil, err
	}
	scopes = slices.DeleteFunc(scopes, func(s Scope) bool {
		return !strings.Contains(s.UUID, scopeUUID) || s.Status.Status != jujustatus.Available
	})

	eg, egctx := errgroup.WithContext(ctx)
	result := make([][]Essential, len(scopes))
	for i := range scopes {
		eg.Go(func() error {
			s, err := clientRepo.Status(egctx, scopes[i].UUID, []string{"application", "*"})
			if err != nil {
				return err
			}
			for name := range s.Applications {
				if !strings.Contains(s.Applications[name].Charm, charmName) {
					continue
				}
				units := []EssentialUnit{}
				for uname := range s.Applications[name].Units {
					units = append(units, EssentialUnit{
						Name:      uname,
						Directive: s.Applications[name].Units[uname].Machine,
					})
				}
				result[i] = append(result[i], Essential{
					Type:      essentialType,
					Name:      name,
					ScopeUUID: scopes[i].UUID,
					ScopeName: scopes[i].Name,
					Units:     units,
				})
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	ret := []Essential{}
	for i := range result {
		ret = append(ret, result[i]...)
	}
	slices.SortFunc(ret, func(e1, e2 Essential) int {
		return strings.Compare(e1.Name, e2.Name)
	})
	return ret, nil
}

func createEssential(ctx context.Context, serverRepo ServerRepo, machineRepo MachineRepo, facilityRepo FacilityRepo, tagRepo TagRepo, uuid, machineID, prefix string, charms []EssentialCharm, configs map[string]string, tags []string) error {
	var (
		directive string
		err       error
	)
	if machineID != "" {
		directive, err = getDirective(ctx, machineRepo, machineID)
		if err != nil {
			return err
		}
		for _, tag := range tags {
			_, _ = tagRepo.Create(ctx, tag, BuiltInMachineTagComment)
			if err := tagRepo.AddMachines(ctx, tag, []string{machineID}); err != nil {
				return err
			}
		}
	}

	base, err := defaultBase(ctx, serverRepo)
	if err != nil {
		return err
	}

	eg, egctx := errgroup.WithContext(ctx)
	for _, charm := range charms {
		eg.Go(func() error {
			name := toEssentialName(prefix, charm.Name)
			placements := []instance.Placement{}
			if directive != "" && !charm.Subordinate {
				placement := toPlacement(&MachinePlacement{LXD: charm.LXD, Machine: charm.Machine}, directive)
				placements = append(placements, *placement)
			}
			_, err := facilityRepo.Create(egctx, uuid, name, configs[charm.Name], charm.Name, charm.Channel, 0, 1, &base, placements, nil, true)
			return err
		})
	}
	return eg.Wait()
}

func createEssentialRelations(ctx context.Context, facilityRepo FacilityRepo, uuid string, endpointList [][]string) error {
	eg, egctx := errgroup.WithContext(ctx)
	for _, endpoints := range endpointList {
		eg.Go(func() error {
			_, err := facilityRepo.CreateRelation(egctx, uuid, endpoints)
			return err
		})
	}
	return eg.Wait()
}

func toEssentialName(prefix, charm string) string {
	if strings.HasPrefix(charm, "ch:") {
		return prefix + "-" + strings.Split(charm, ":")[1]
	}
	return prefix + "-" + charm
}

func toEndpointList(prefix string, relationList [][]string) [][]string {
	endpointList := [][]string{}
	for _, relations := range relationList {
		endpoints := []string{}
		for _, relation := range relations {
			endpoints = append(endpoints, toEssentialName(prefix, relation))
		}
		endpointList = append(endpointList, endpoints)
	}
	return endpointList
}

func getDirective(ctx context.Context, machineRepo MachineRepo, machineID string) (string, error) {
	machine, err := machineRepo.Get(ctx, machineID)
	if err != nil {
		return "", err
	}
	if machine.Status != node.StatusDeployed {
		return "", connect.NewError(connect.CodeInvalidArgument, errors.New("machine status is not deployed"))
	}
	return getJujuMachineID(machine.WorkloadAnnotations)
}

func (uc *EssentialUseCase) ListGPURelationsByMachine(ctx context.Context, scopeUUID, facilityName, machineID string) ([]*pb.GPURelation, error) {
	// Get Machine data and Hostname by machine_id
	machine, err := uc.machine.Get(ctx, machineID)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("machine not found: %w", err))
	}

	if machine.Hostname == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("machine hostname is empty"))
	}

	// Get Kubernetes configuration
	config, err := kubeConfig(ctx, uc.facility, uc.action, scopeUUID, facilityName)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get kube config: %w", err))
	}

	// Query Kubernetes Node by Hostname
	_, err = uc.kubeCore.GetNode(ctx, config, machine.Hostname)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("kubernetes node not found: %w", err))
	}

	// Parse GPU annotation data to get GPU information
	nodeGpuMap, err := uc.getNodeGpuMap(ctx, config, machine.Hostname)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get node GPU map: %w", err))
	}

	// Query all pods with vGPU configuration on this Node
	pods, err := uc.kubeCore.ListPods(ctx, config, "")
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to list pods: %w", err))
	}

	// Filter pods on this Node with vGPU configuration
	vgpuPods := filterVGpuPodsOnNode(ctx, pods, machine.Hostname, config, uc.kubeCore, uc.kubeApps)

	// Build GPU relations array
	return uc.buildGPURelationsForMachine(machineID, machine.Hostname, nodeGpuMap, vgpuPods)
}

// buildGPURelationsForMachine constructs the complete GPU relations array for a machine
func (uc *EssentialUseCase) buildGPURelationsForMachine(machineID, hostname string, nodeGpuMap map[string]string, vgpuPods []EssentialVGpuPodInfo) ([]*pb.GPURelation, error) {
	var gpuRelations []*pb.GPURelation

	// Add Machine relation
	machineRelation := uc.createMachineRelation(machineID, hostname)
	gpuRelations = append(gpuRelations, machineRelation)

	// Add GPU relations
	gpuRelations = append(gpuRelations, uc.createGPURelations(machineID, nodeGpuMap)...)

	// Add Pod and vGPU relations
	podVGpuRelations := uc.createPodAndVGpuRelations(vgpuPods)
	gpuRelations = append(gpuRelations, podVGpuRelations...)

	return gpuRelations, nil
}

// createMachineRelation creates a machine relation entity
func (uc *EssentialUseCase) createMachineRelation(machineID, hostname string) *pb.GPURelation {
	machineRelation := &pb.GPURelation{}
	machineEntity := &pb.GPURelation_Machine{}
	machineEntity.SetId(machineID)
	machineEntity.SetHostname(hostname)
	machineRelation.SetMachine(machineEntity)
	return machineRelation
}

// createGPURelations creates GPU relation entities for all GPUs on the node
func (uc *EssentialUseCase) createGPURelations(machineID string, nodeGpuMap map[string]string) []*pb.GPURelation {
	var gpuRelations []*pb.GPURelation

	for gpuID, gpuVendorProduct := range nodeGpuMap {
		// Parse vendor and product from the format "NVIDIA-NVIDIA GeForce RTX 4090"
		vendor, product := parseVendorProduct(gpuVendorProduct)

		gpuRelation := &pb.GPURelation{}
		gpuEntity := &pb.GPURelation_GPU{}
		gpuEntity.SetId(gpuID)
		gpuEntity.SetVendor(vendor)
		gpuEntity.SetProduct(product)
		gpuEntity.SetMachineId(machineID)
		gpuRelation.SetGpu(gpuEntity)
		gpuRelations = append(gpuRelations, gpuRelation)
	}

	return gpuRelations
}

// createPodAndVGpuRelations creates pod and vGPU relation entities for all vGPU pods
func (uc *EssentialUseCase) createPodAndVGpuRelations(vgpuPods []EssentialVGpuPodInfo) []*pb.GPURelation {
	var relations []*pb.GPURelation

	for _, vgpuPod := range vgpuPods {
		// Collect all GPU IDs used by this pod
		gpuIDs := uc.collectGpuIDsFromPod(vgpuPod)

		// Add Pod relation
		podRelation := uc.createPodRelation(vgpuPod, gpuIDs)
		relations = append(relations, podRelation)

		// Add vGPU relations for each vGPU allocation
		vgpuRelations := uc.createVGpuRelations(vgpuPod)
		relations = append(relations, vgpuRelations...)
	}

	return relations
}

// collectGpuIDsFromPod extracts all GPU IDs used by a pod
func (uc *EssentialUseCase) collectGpuIDsFromPod(vgpuPod EssentialVGpuPodInfo) []string {
	var gpuIDs []string
	for _, vgpuAlloc := range vgpuPod.VGpuInfos {
		gpuIDs = append(gpuIDs, vgpuAlloc.GpuUUID)
	}
	return gpuIDs
}

// createPodRelation creates a pod relation entity
func (uc *EssentialUseCase) createPodRelation(vgpuPod EssentialVGpuPodInfo, gpuIDs []string) *pb.GPURelation {
	podRelation := &pb.GPURelation{}
	podEntity := &pb.GPURelation_Pod{}
	podEntity.SetName(vgpuPod.Pod.Name)
	podEntity.SetNamespace(vgpuPod.Pod.Namespace)
	podEntity.SetModelName(vgpuPod.ModelName)
	podEntity.SetGpuId(gpuIDs)
	podRelation.SetPod(podEntity)
	return podRelation
}

// createVGpuRelations creates vGPU relation entities for all vGPU allocations in a pod
func (uc *EssentialUseCase) createVGpuRelations(vgpuPod EssentialVGpuPodInfo) []*pb.GPURelation {
	var vgpuRelations []*pb.GPURelation

	for _, vgpuAlloc := range vgpuPod.VGpuInfos {
		vgpuRelation := uc.createSingleVGpuRelation(vgpuPod.Pod.Name, vgpuAlloc)
		vgpuRelations = append(vgpuRelations, vgpuRelation)
	}

	return vgpuRelations
}

// createSingleVGpuRelation creates a single vGPU relation entity
func (uc *EssentialUseCase) createSingleVGpuRelation(podName string, vgpuAlloc EssentialVGpuAllocation) *pb.GPURelation {
	vgpuRelation := &pb.GPURelation{}

	// Convert VramMib string to bytes
	vramBytes, err := convertVramMibToBytes(vgpuAlloc.VramMib)
	if err != nil {
		vramBytes = 0 // Default to 0 if conversion fails
	}

	// Convert VcoresPercent string to float
	vcoresPercent, err := convertVcoresPercentToFloat(vgpuAlloc.VcoresPercent)
	if err != nil {
		vcoresPercent = 0.0 // Default to 0.0 if conversion fails
	}

	// Convert bind time to timestamp
	var boundAt *timestamppb.Timestamp
	if vgpuAlloc.BindTime != "" {
		if bindTimeUnix, err := strconv.ParseInt(vgpuAlloc.BindTime, 10, 64); err == nil {
			boundAt = timestamppb.New(time.Unix(bindTimeUnix, 0))
		}
	}

	vgpuEntity := &pb.GPURelationVGPU{}
	vgpuEntity.SetGpuId(vgpuAlloc.GpuUUID)
	vgpuEntity.SetPodName(podName)
	vgpuEntity.SetVramBytes(vramBytes)
	vgpuEntity.SetVcoresPercent(vcoresPercent)
	vgpuEntity.SetBindingPhase(vgpuAlloc.BindPhase)
	if boundAt != nil {
		vgpuEntity.SetBoundAt(boundAt)
	}
	vgpuRelation.SetVgpu(vgpuEntity)

	return vgpuRelation
}

func (uc *EssentialUseCase) ListGPURelationsByModel(ctx context.Context, scopeUUID, facilityName, namespace, modelName string) ([]*pb.GPURelation, error) {
	// Get Kubernetes configuration
	config, err := kubeConfig(ctx, uc.facility, uc.action, scopeUUID, facilityName)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get kube config: %w", err))
	}

	// Search for deployments with the specified model name
	label := "model-name=" + modelName
	deployments, err := uc.kubeApps.ListDeploymentsByLabel(ctx, config, namespace, label)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get deployments: %w", err))
	}

	return uc.buildGPURelationsForModel(ctx, config, deployments, modelName)
}

// buildGPURelationsForModel constructs GPU relations for all pods in deployments with the specified model
func (uc *EssentialUseCase) buildGPURelationsForModel(ctx context.Context, config *rest.Config, deployments []Deployment, modelName string) ([]*pb.GPURelation, error) {
	var gpuRelations []*pb.GPURelation
	tracker := newRelationTracker()

	for i := range deployments {
		relations, err := uc.processDeploymentForGPURelations(ctx, config, deployments[i], modelName, tracker)
		if err != nil {
			continue // Skip deployments with errors
		}
		gpuRelations = append(gpuRelations, relations...)
	}

	return gpuRelations, nil
}

// relationTracker tracks processed entities to avoid duplicates
type relationTracker struct {
	processedPods     map[string]bool
	processedMachines map[string]bool
	processedGPUs     map[string]bool
}

func newRelationTracker() *relationTracker {
	return &relationTracker{
		processedPods:     make(map[string]bool),
		processedMachines: make(map[string]bool),
		processedGPUs:     make(map[string]bool),
	}
}

// processDeploymentForGPURelations processes a single deployment and returns GPU relations
func (uc *EssentialUseCase) processDeploymentForGPURelations(ctx context.Context, config *rest.Config, deployment Deployment, modelName string, tracker *relationTracker) ([]*pb.GPURelation, error) {
	if deployment.Spec.Selector == nil || deployment.Spec.Selector.MatchLabels == nil {
		return nil, nil
	}

	// Build pod label selector from deployment selector
	podLabelSelector := uc.buildPodLabelSelector(deployment.Spec.Selector.MatchLabels)

	// Get pods for this deployment
	pods, err := uc.kubeCore.ListPodsByLabel(ctx, config, deployment.Namespace, podLabelSelector)
	if err != nil {
		return nil, err
	}

	var relations []*pb.GPURelation
	for i := range pods {
		podRelations, err := uc.processPodForGPURelations(ctx, config, pods[i], modelName, tracker)
		if err != nil {
			continue // Skip pods with errors
		}
		relations = append(relations, podRelations...)
	}

	return relations, nil
}

// buildPodLabelSelector builds a label selector string from match labels
func (uc *EssentialUseCase) buildPodLabelSelector(matchLabels map[string]string) string {
	var selector string
	for key, value := range matchLabels {
		if selector != "" {
			selector += ","
		}
		selector += key + "=" + value
	}
	return selector
}

// processPodForGPURelations processes a single pod and returns GPU relations
func (uc *EssentialUseCase) processPodForGPURelations(ctx context.Context, config *rest.Config, pod Pod, modelName string, tracker *relationTracker) ([]*pb.GPURelation, error) {
	// Skip if already processed
	podKey := pod.Namespace + "/" + pod.Name
	if tracker.processedPods[podKey] {
		return nil, nil
	}
	tracker.processedPods[podKey] = true

	// Check if pod has vGPU allocation
	vgpuDevicesAllocated, hasVGpuAllocation := pod.Annotations["hami.io/vgpu-devices-allocated"]
	if !hasVGpuAllocation || vgpuDevicesAllocated == "" {
		return nil, nil
	}

	// Parse vGPU allocation information
	vgpuAllocations := parseVGpuDevicesAllocated(vgpuDevicesAllocated)
	if len(vgpuAllocations) == 0 {
		return nil, nil
	}

	// Get machine information
	machineID, err := uc.getMachineIDFromNodeName(ctx, pod.Spec.NodeName)
	if err != nil {
		return nil, err
	}

	var relations []*pb.GPURelation

	// Add Machine relation if not already added
	if machineRelation := uc.addMachineRelationIfNeeded(ctx, machineID, tracker); machineRelation != nil {
		relations = append(relations, machineRelation)
	}

	// Get GPU information from node
	nodeGpuMap, err := uc.getNodeGpuMap(ctx, config, pod.Spec.NodeName)
	if err != nil {
		nodeGpuMap = make(map[string]string)
	}

	// Process vGPU allocations and create GPU relations
	gpuIDs, gpuRelations := uc.processVGpuAllocationsForGPUs(vgpuAllocations, nodeGpuMap, machineID, tracker)
	relations = append(relations, gpuRelations...)

	// Add Pod relation
	podRelation := uc.createPodRelationForModel(pod, modelName, gpuIDs)
	relations = append(relations, podRelation)

	// Add vGPU relations
	vgpuRelations := uc.createVGpuRelationsForModel(pod, vgpuAllocations)
	relations = append(relations, vgpuRelations...)

	return relations, nil
}

// addMachineRelationIfNeeded adds a machine relation if it hasn't been processed yet
func (uc *EssentialUseCase) addMachineRelationIfNeeded(ctx context.Context, machineID string, tracker *relationTracker) *pb.GPURelation {
	if tracker.processedMachines[machineID] {
		return nil
	}

	machine, err := uc.machine.Get(ctx, machineID)
	if err != nil {
		return nil
	}

	tracker.processedMachines[machineID] = true
	return uc.createMachineRelation(machineID, machine.Hostname)
}

// processVGpuAllocationsForGPUs processes vGPU allocations and creates GPU relations
func (uc *EssentialUseCase) processVGpuAllocationsForGPUs(vgpuAllocations []EssentialVGpuAllocation, nodeGpuMap map[string]string, machineID string, tracker *relationTracker) ([]string, []*pb.GPURelation) {
	var gpuIDs []string
	var gpuRelations []*pb.GPURelation

	for _, vgpuAlloc := range vgpuAllocations {
		gpuIDs = append(gpuIDs, vgpuAlloc.GpuUUID)

		// Add GPU relation if not already added
		if !tracker.processedGPUs[vgpuAlloc.GpuUUID] {
			if gpuVendorProduct, exists := nodeGpuMap[vgpuAlloc.GpuUUID]; exists {
				vendor, product := parseVendorProduct(gpuVendorProduct)

				gpuRelation := &pb.GPURelation{}
				gpuEntity := &pb.GPURelation_GPU{}
				gpuEntity.SetId(vgpuAlloc.GpuUUID)
				gpuEntity.SetVendor(vendor)
				gpuEntity.SetProduct(product)
				gpuEntity.SetMachineId(machineID)
				gpuRelation.SetGpu(gpuEntity)
				gpuRelations = append(gpuRelations, gpuRelation)
				tracker.processedGPUs[vgpuAlloc.GpuUUID] = true
			}
		}
	}

	return gpuIDs, gpuRelations
}

// createPodRelationForModel creates a pod relation entity for model-based query
func (uc *EssentialUseCase) createPodRelationForModel(pod Pod, modelName string, gpuIDs []string) *pb.GPURelation {
	podRelation := &pb.GPURelation{}
	podEntity := &pb.GPURelation_Pod{}
	podEntity.SetName(pod.Name)
	podEntity.SetNamespace(pod.Namespace)
	podEntity.SetModelName(modelName)
	podEntity.SetGpuId(gpuIDs)
	podRelation.SetPod(podEntity)
	return podRelation
}

// createVGpuRelationsForModel creates vGPU relation entities for model-based query
func (uc *EssentialUseCase) createVGpuRelationsForModel(pod Pod, vgpuAllocations []EssentialVGpuAllocation) []*pb.GPURelation {
	var vgpuRelations []*pb.GPURelation

	bindTime := pod.Annotations["hami.io/bind-time"]
	bindPhase := pod.Annotations["hami.io/bind-phase"]

	for _, vgpuAlloc := range vgpuAllocations {
		vgpuRelation := uc.createVGpuRelationForModel(pod.Name, vgpuAlloc, bindTime, bindPhase)
		vgpuRelations = append(vgpuRelations, vgpuRelation)
	}

	return vgpuRelations
}

// createVGpuRelationForModel creates a single vGPU relation entity for model-based query
func (uc *EssentialUseCase) createVGpuRelationForModel(podName string, vgpuAlloc EssentialVGpuAllocation, bindTime, bindPhase string) *pb.GPURelation {
	vgpuRelation := &pb.GPURelation{}

	// Convert VramMib string to bytes
	vramBytes, err := convertVramMibToBytes(vgpuAlloc.VramMib)
	if err != nil {
		vramBytes = 0
	}

	// Convert VcoresPercent string to float
	vcoresPercent, err := convertVcoresPercentToFloat(vgpuAlloc.VcoresPercent)
	if err != nil {
		vcoresPercent = 0.0
	}

	// Convert bind time to timestamp
	var boundAt *timestamppb.Timestamp
	if bindTime != "" {
		if bindTimeUnix, err := strconv.ParseInt(bindTime, 10, 64); err == nil {
			boundAt = timestamppb.New(time.Unix(bindTimeUnix, 0))
		}
	}

	vgpuEntity := &pb.GPURelationVGPU{}
	vgpuEntity.SetGpuId(vgpuAlloc.GpuUUID)
	vgpuEntity.SetPodName(podName)
	vgpuEntity.SetVramBytes(vramBytes)
	vgpuEntity.SetVcoresPercent(vcoresPercent)
	vgpuEntity.SetBindingPhase(bindPhase)
	if boundAt != nil {
		vgpuEntity.SetBoundAt(boundAt)
	}
	vgpuRelation.SetVgpu(vgpuEntity)

	return vgpuRelation
}

// filterVGpuPodsOnNode filters pods on the specified Node that have vGPU configuration
func filterVGpuPodsOnNode(ctx context.Context, pods []Pod, nodeName string, config *rest.Config, kubeRepo KubeCoreRepo, kubeApps KubeAppsRepo) []EssentialVGpuPodInfo {
	var vgpuPods []EssentialVGpuPodInfo

	for i := range pods {
		pod := &pods[i]

		// Check if Pod is on the specified Node
		if pod.Spec.NodeName != nodeName {
			continue
		}

		// Check if Pod has vGPU allocation annotation
		vgpuDevicesAllocated, hasVGpuAllocation := pod.Annotations["hami.io/vgpu-devices-allocated"]
		if !hasVGpuAllocation || vgpuDevicesAllocated == "" {
			continue
		}

		// Parse vGPU configuration
		vgpuAllocations := parseVGpuDevicesAllocated(vgpuDevicesAllocated)
		if len(vgpuAllocations) == 0 {
			continue
		}

		// Get additional Pod information
		bindTime := pod.Annotations["hami.io/bind-time"]
		bindPhase := pod.Annotations["hami.io/bind-phase"]

		// Update vGPU allocation bind information
		for j := range vgpuAllocations {
			vgpuAllocations[j].BindTime = bindTime
			vgpuAllocations[j].BindPhase = bindPhase
		}

		// Get model-name (from Deployment labels via Pod's owner reference)
		modelName := getModelNameFromPod(ctx, pod, config, kubeApps)

		vgpuPods = append(vgpuPods, EssentialVGpuPodInfo{
			Pod:       pod,
			ModelName: modelName,
			VGpuInfos: vgpuAllocations,
		})
	}

	return vgpuPods
}

// parseVGpuDevicesAllocated parses hami.io/vgpu-devices-allocated annotation
// Format: GPU-c15ecdf3-444a-2d02-29e9-e978b2514335,NVIDIA,3684,25:GPU-663aa370-535a-33b8-e01f-b325fb2025c7,NVIDIA,4684,35:;
func parseVGpuDevicesAllocated(annotation string) []EssentialVGpuAllocation {
	const (
		vgpuAllocationFieldCount = 4 // GPU-uuid,vendor,vram,vcores
		vgpuUUIDIndex            = 0
		vgpuVendorIndex          = 1
		vgpuVramIndex            = 2
		vgpuVcoresIndex          = 3
	)
	var allocations []EssentialVGpuAllocation

	// Remove trailing semicolon and split by colon
	annotation = strings.TrimSuffix(annotation, ";")
	entries := strings.Split(annotation, ":")

	for _, entry := range entries {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}

		// Split by comma: GPU-uuid,vendor,vram,vcores
		parts := strings.Split(entry, ",")
		if len(parts) != vgpuAllocationFieldCount {
			continue
		}

		allocation := EssentialVGpuAllocation{
			GpuUUID:       parts[vgpuUUIDIndex],   // GPU-c15ecdf3-444a-2d02-29e9-e978b2514335
			Vendor:        parts[vgpuVendorIndex], // NVIDIA
			VramMib:       parts[vgpuVramIndex],   // 3684
			VcoresPercent: parts[vgpuVcoresIndex], // 25
		}

		allocations = append(allocations, allocation)
	}

	return allocations
}

// getModelNameFromPod extracts model-name from Pod's Deployment labels
func getModelNameFromPod(ctx context.Context, pod *Pod, config *rest.Config, kubeApps KubeAppsRepo) string {
	// First try to get from Pod labels directly
	if modelName, exists := pod.Labels["model-name"]; exists {
		return modelName
	}

	// Try to get from Deployment labels via owner references
	if len(pod.OwnerReferences) > 0 {
		for _, ownerRef := range pod.OwnerReferences {
			// Check if owner is ReplicaSet
			if ownerRef.Kind == "ReplicaSet" {
				// Extract deployment name from ReplicaSet name
				// ReplicaSet naming pattern: {deployment-name}-{random-hash}
				replicaSetName := ownerRef.Name
				deploymentName := extractDeploymentNameFromReplicaSet(replicaSetName)

				if deploymentName != "" {
					// Get Deployment and check its labels
					deployment, err := kubeApps.GetDeployment(ctx, config, pod.Namespace, deploymentName)
					if err == nil && deployment != nil {
						if modelName, exists := deployment.Labels["model-name"]; exists {
							return modelName
						}
					}
				}
			}
		}
	}

	// Fallback: try to infer from Pod labels
	if appName, exists := pod.Labels["app"]; exists {
		return appName
	}

	// If not found, return Pod name as default
	return pod.Name
}

// extractDeploymentNameFromReplicaSet extracts deployment name from ReplicaSet name
// ReplicaSet naming pattern: {deployment-name}-{random-hash}
func extractDeploymentNameFromReplicaSet(replicaSetName string) string {
	const (
		replicaSetHashMinLength = 8  // Minimum length of ReplicaSet hash suffix
		replicaSetHashMaxLength = 10 // Maximum length of ReplicaSet hash suffix
	)

	// Find the last dash and remove the hash part
	lastDashIndex := strings.LastIndex(replicaSetName, "-")
	if lastDashIndex > 0 {
		// Check if the part after the last dash looks like a hash (alphanumeric, typically 8-10 chars)
		hashPart := replicaSetName[lastDashIndex+1:]
		if len(hashPart) >= replicaSetHashMinLength && len(hashPart) <= replicaSetHashMaxLength {
			// Assume it's a hash if it's the right length and contains only lowercase letters and numbers
			// ReplicaSet hashes are typically like: 54bb8b45bb, 85c77654c7, 7d8b49557f
			isHash := true
			hasNumber := false
			for _, r := range hashPart {
				if (r >= 'a' && r <= 'f') || (r >= '0' && r <= '9') {
					if r >= '0' && r <= '9' {
						hasNumber = true
					}
				} else {
					isHash = false
					break
				}
			}
			// Only consider it a hash if it contains at least one number and only hex characters
			if isHash && hasNumber {
				return replicaSetName[:lastDashIndex]
			}
		}
	}
	return replicaSetName
}

// getMachineIDFromNodeName gets machine_id from node_name using the same logic as GetGpuRelationByMachine
func (uc *EssentialUseCase) getMachineIDFromNodeName(ctx context.Context, nodeName string) (string, error) {
	// List all machines and find the one with matching hostname
	machines, err := uc.machine.List(ctx)
	if err != nil {
		return "", err
	}

	for _, machine := range machines {
		if machine.Hostname == nodeName {
			return machine.SystemID, nil
		}
	}

	return "", fmt.Errorf("machine not found for node: %s", nodeName)
}

// getNodeGpuMap creates a mapping from GPU UUID to GPU vendor for a specific node
func (uc *EssentialUseCase) getNodeGpuMap(ctx context.Context, config *rest.Config, nodeName string) (map[string]string, error) {
	const (
		hamiAnnotationKey        = "hami.io/node-nvidia-register"
		gpuEntryMinFieldCount    = 5 // Minimum fields: UUID,index,memory,count,vendor,...
		gpuEntryUUIDFieldIndex   = 0
		gpuEntryVendorFieldIndex = 4
	)

	// Get the node
	node, err := uc.kubeCore.GetNode(ctx, config, nodeName)
	if err != nil {
		return nil, err
	}

	// Parse GPU annotations to get GPU UUID to vendor mapping
	gpuMap := make(map[string]string)

	annotationValue, exists := node.Annotations[hamiAnnotationKey]
	if !exists || annotationValue == "" {
		return gpuMap, nil
	}

	// Remove leading and trailing quotes
	annotationValue = strings.Trim(annotationValue, "'\"")

	// Split different GPU cards by colon
	gpuEntries := strings.Split(annotationValue, ":")

	for i := 0; i < len(gpuEntries)-1; i++ { // Last element is usually empty, so -1
		entry := strings.TrimSpace(gpuEntries[i])
		if entry == "" {
			continue
		}

		// Parse single GPU entry
		// Format: GPU-663aa370-535a-33b8-e01f-b325fb2025c7,10,24564,100,NVIDIA-NVIDIA GeForce RTX 4090,0,true,0,hami-core
		parts := strings.Split(entry, ",")
		if len(parts) < gpuEntryMinFieldCount {
			continue // Skip incorrectly formatted entries
		}

		cardUUID := parts[gpuEntryUUIDFieldIndex]
		cardVendor := parts[gpuEntryVendorFieldIndex] // Extract full vendor name from parts[4]

		gpuMap[cardUUID] = cardVendor
	}

	return gpuMap, nil
}

// parseVendorProduct parses vendor and product from GPU vendor string
// Input: "NVIDIA-NVIDIA GeForce RTX 4090"
// Output: vendor="NVIDIA", product="NVIDIA GeForce RTX 4090"
func parseVendorProduct(vendorProduct string) (vendor, product string) {
	const (
		maxVendorProductParts = 2 // Split into at most 2 parts: vendor and product
		vendorIndex           = 0
		productIndex          = 1
	)

	parts := strings.SplitN(vendorProduct, "-", maxVendorProductParts)
	if len(parts) >= vendorIndex+1 {
		vendor = parts[vendorIndex]
	}
	if len(parts) >= productIndex+1 {
		product = parts[productIndex]
	} else {
		product = vendorProduct
	}
	return vendor, product
}

// convertVramMibToBytes converts VRAM from MiB string to bytes
func convertVramMibToBytes(vramMib string) (uint64, error) {
	mib, err := strconv.ParseUint(vramMib, 10, 64)
	if err != nil {
		return 0, err
	}
	// Convert MiB to bytes (1 MiB = 1024 * 1024 bytes)
	return mib * 1024 * 1024, nil
}

// convertVcoresPercentToFloat converts vcores percent from string to float32
func convertVcoresPercentToFloat(vcoresPercent string) (float32, error) {
	percent, err := strconv.ParseFloat(vcoresPercent, 32)
	if err != nil {
		return 0.0, err
	}
	return float32(percent), nil
}
