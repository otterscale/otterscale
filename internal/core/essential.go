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
			return len(parts) < 2 || !strings.HasSuffix(parts[0], LabelDomain)
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
	var gpuRelations []*pb.GPURelation

	// Add Machine relation
	machineRelation := &pb.GPURelation{}
	machineEntity := &pb.GPURelation_Machine{}
	machineEntity.SetId(machineID)
	machineEntity.SetHostname(machine.Hostname)
	machineRelation.SetMachine(machineEntity)
	gpuRelations = append(gpuRelations, machineRelation)

	// Add GPU relations
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

	// Process vGPU pods and create Pod and vGPU relations
	for _, vgpuPod := range vgpuPods {
		// Collect all GPU IDs used by this pod
		var gpuIDs []string
		for _, vgpuAlloc := range vgpuPod.VGpuInfos {
			gpuIDs = append(gpuIDs, vgpuAlloc.GpuUUID)
		}

		// Add Pod relation
		podRelation := &pb.GPURelation{}
		podEntity := &pb.GPURelation_Pod{}
		podEntity.SetName(vgpuPod.Pod.Name)
		podEntity.SetNamespace(vgpuPod.Pod.Namespace)
		podEntity.SetModelName(vgpuPod.ModelName)
		podEntity.SetGpuId(gpuIDs)
		podRelation.SetPod(podEntity)
		gpuRelations = append(gpuRelations, podRelation)

		// Add vGPU relations for each vGPU allocation
		for _, vgpuAlloc := range vgpuPod.VGpuInfos {
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
			vgpuEntity.SetPodName(vgpuPod.Pod.Name)
			vgpuEntity.SetVramBytes(vramBytes)
			vgpuEntity.SetVcoresPercent(vcoresPercent)
			vgpuEntity.SetBindingPhase(vgpuAlloc.BindPhase)
			if boundAt != nil {
				vgpuEntity.SetBoundAt(boundAt)
			}
			vgpuRelation.SetVgpu(vgpuEntity)
			gpuRelations = append(gpuRelations, vgpuRelation)
		}
	}

	return gpuRelations, nil
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

	var gpuRelations []*pb.GPURelation
	processedPods := make(map[string]bool)     // Track processed pods to avoid duplicates
	processedMachines := make(map[string]bool) // Track processed machines to avoid duplicates
	processedGPUs := make(map[string]bool)     // Track processed GPUs to avoid duplicates

	for _, deployment := range deployments {
		selector := deployment.Spec.Selector
		if selector == nil || selector.MatchLabels == nil {
			continue
		}

		// Build pod label selector from deployment selector
		podLabelSelector := ""
		for key, value := range selector.MatchLabels {
			if podLabelSelector != "" {
				podLabelSelector += ","
			}
			podLabelSelector += key + "=" + value
		}

		// Get pods for this deployment
		pods, err := uc.kubeCore.ListPodsByLabel(ctx, config, deployment.Namespace, podLabelSelector)
		if err != nil {
			continue
		}

		for _, pod := range pods {
			// Skip if we've already processed this pod
			podKey := pod.Namespace + "/" + pod.Name
			if processedPods[podKey] {
				continue
			}
			processedPods[podKey] = true

			// Check if pod has vGPU allocation
			vgpuDevicesAllocated, hasVGpuAllocation := pod.Annotations["hami.io/vgpu-devices-allocated"]
			if !hasVGpuAllocation || vgpuDevicesAllocated == "" {
				continue
			}

			// Parse vGPU allocation information
			vgpuAllocations, err := parseVGpuDevicesAllocated(vgpuDevicesAllocated)
			if err != nil || len(vgpuAllocations) == 0 {
				continue
			}

			// Get machine information
			machineID, err := uc.getMachineIDFromNodeName(ctx, pod.Spec.NodeName)
			if err != nil {
				continue // Skip if we can't get machine ID
			}

			// Add Machine relation if not already added
			if !processedMachines[machineID] {
				machine, err := uc.machine.Get(ctx, machineID)
				if err == nil {
					machineRelation := &pb.GPURelation{}
					machineEntity := &pb.GPURelation_Machine{}
					machineEntity.SetId(machineID)
					machineEntity.SetHostname(machine.Hostname)
					machineRelation.SetMachine(machineEntity)
					gpuRelations = append(gpuRelations, machineRelation)
					processedMachines[machineID] = true
				}
			}

			// Get GPU information from node
			nodeGpuMap, err := uc.getNodeGpuMap(ctx, config, pod.Spec.NodeName)
			if err != nil {
				nodeGpuMap = make(map[string]string) // Empty map if we can't get node info
			}

			// Collect GPU IDs used by this pod for the Pod relation
			var gpuIDs []string
			for _, vgpuAlloc := range vgpuAllocations {
				gpuIDs = append(gpuIDs, vgpuAlloc.GpuUUID)

				// Add GPU relation if not already added
				if !processedGPUs[vgpuAlloc.GpuUUID] {
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
						processedGPUs[vgpuAlloc.GpuUUID] = true
					}
				}
			}

			// Add Pod relation
			podRelation := &pb.GPURelation{}
			podEntity := &pb.GPURelation_Pod{}
			podEntity.SetName(pod.Name)
			podEntity.SetNamespace(pod.Namespace)
			podEntity.SetModelName(modelName)
			podEntity.SetGpuId(gpuIDs)
			podRelation.SetPod(podEntity)
			gpuRelations = append(gpuRelations, podRelation)

			// Add vGPU relations
			bindTime := pod.Annotations["hami.io/bind-time"]
			bindPhase := pod.Annotations["hami.io/bind-phase"]

			for _, vgpuAlloc := range vgpuAllocations {
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
				if bindTime != "" {
					if bindTimeUnix, err := strconv.ParseInt(bindTime, 10, 64); err == nil {
						boundAt = timestamppb.New(time.Unix(bindTimeUnix, 0))
					}
				}

				vgpuEntity := &pb.GPURelationVGPU{}
				vgpuEntity.SetGpuId(vgpuAlloc.GpuUUID)
				vgpuEntity.SetPodName(pod.Name)
				vgpuEntity.SetVramBytes(vramBytes)
				vgpuEntity.SetVcoresPercent(vcoresPercent)
				vgpuEntity.SetBindingPhase(bindPhase)
				if boundAt != nil {
					vgpuEntity.SetBoundAt(boundAt)
				}
				vgpuRelation.SetVgpu(vgpuEntity)
				gpuRelations = append(gpuRelations, vgpuRelation)
			}
		}
	}

	return gpuRelations, nil
}

// filterVGpuPodsOnNode filters pods on the specified Node that have vGPU configuration
func filterVGpuPodsOnNode(ctx context.Context, pods []Pod, nodeName string, config *rest.Config, kubeRepo KubeCoreRepo, kubeAppsRepo KubeAppsRepo) []EssentialVGpuPodInfo {
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
		vgpuAllocations, err := parseVGpuDevicesAllocated(vgpuDevicesAllocated)
		if err != nil || len(vgpuAllocations) == 0 {
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
		modelName := getModelNameFromPod(ctx, pod, config, kubeRepo, kubeAppsRepo)

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
func parseVGpuDevicesAllocated(annotation string) ([]EssentialVGpuAllocation, error) {
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
		if len(parts) != 4 {
			continue
		}

		allocation := EssentialVGpuAllocation{
			GpuUUID:       parts[0], // GPU-c15ecdf3-444a-2d02-29e9-e978b2514335
			Vendor:        parts[1], // NVIDIA
			VramMib:       parts[2], // 3684
			VcoresPercent: parts[3], // 25
		}

		allocations = append(allocations, allocation)
	}

	return allocations, nil
}

// getModelNameFromPod extracts model-name from Pod's Deployment labels
func getModelNameFromPod(ctx context.Context, pod *Pod, config *rest.Config, kubeRepo KubeCoreRepo, kubeAppsRepo KubeAppsRepo) string {
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
					deployment, err := kubeAppsRepo.GetDeployment(ctx, config, pod.Namespace, deploymentName)
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
	// Find the last dash and remove the hash part
	lastDashIndex := strings.LastIndex(replicaSetName, "-")
	if lastDashIndex > 0 {
		// Check if the part after the last dash looks like a hash (alphanumeric, typically 8-10 chars)
		hashPart := replicaSetName[lastDashIndex+1:]
		if len(hashPart) >= 8 && len(hashPart) <= 10 {
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
	// Get the node
	node, err := uc.kubeCore.GetNode(ctx, config, nodeName)
	if err != nil {
		return nil, err
	}

	// Parse GPU annotations to get GPU UUID to vendor mapping
	gpuMap := make(map[string]string)

	const hamiAnnotationKey = "hami.io/node-nvidia-register"
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
		if len(parts) < 5 {
			continue // Skip incorrectly formatted entries
		}

		cardUUID := parts[0]
		cardVendor := parts[4] // Extract full vendor name from parts[4]

		gpuMap[cardUUID] = cardVendor
	}

	return gpuMap, nil
}

// parseVendorProduct parses vendor and product from GPU vendor string
// Input: "NVIDIA-NVIDIA GeForce RTX 4090"
// Output: vendor="NVIDIA", product="NVIDIA GeForce RTX 4090"
func parseVendorProduct(vendorProduct string) (vendor, product string) {
	parts := strings.SplitN(vendorProduct, "-", 2)
	if len(parts) >= 1 {
		vendor = parts[0]
	}
	if len(parts) >= 2 {
		product = parts[1]
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
