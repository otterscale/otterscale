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
	jujuyaml "gopkg.in/yaml.v2"
	"k8s.io/client-go/rest"

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

// GPURelation domain models for core layer
type GPURelation struct {
	Machine *GPURelationMachine
	GPU     *GPURelationGPU
	Pod     *GPURelationPod
}

type GPURelationMachine struct {
	ID       string
	Hostname string
}

type GPURelationGPU struct {
	ID        string
	Vendor    string
	Product   string
	MachineID string
	VGPUs     []GPURelationVGPU
}

type GPURelationVGPU struct {
	PodName       string
	BindingPhase  string
	VramBytes     uint64
	VcoresPercent float32
	BoundAt       time.Time
}

type GPURelationPod struct {
	Name      string
	Namespace string
	ModelName string
	GPUIDs    []string
}

type gpoPodRelationTracker struct {
	processedPods     map[string]bool
	processedMachines map[string]bool
	processedGPUs     map[string]bool
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

func (uc *EssentialUseCase) ListGPURelationsByMachine(ctx context.Context, scopeUUID, facilityName, machineID string) ([]GPURelation, error) {
	machine, err := uc.machine.Get(ctx, machineID)
	if err != nil {
		return nil, err
	}
	if machine.Hostname == "" {
		return nil, err
	}

	config, err := kubeConfig(ctx, uc.facility, uc.action, scopeUUID, facilityName)
	if err != nil {
		return nil, err
	}

	node, err := uc.kubeCore.GetNode(ctx, config, machine.Hostname)
	if err != nil {
		return nil, err
	}

	nodeGpuMap, err := uc.getNodeGpuMap(node)
	if err != nil {
		return nil, err
	}

	pods, err := uc.kubeCore.ListPods(ctx, config, "")
	if err != nil {
		return nil, err
	}

	vgpuPods := filterVGpuPodsOnNode(ctx, pods, machine.Hostname, config, uc.kubeCore, uc.kubeApps)
	return uc.buildGPURelationsForMachine(machineID, machine.Hostname, nodeGpuMap, vgpuPods)
}

func (uc *EssentialUseCase) buildGPURelationsForMachine(machineID, hostname string, nodeGpuMap map[string]string, vgpuPods []EssentialVGpuPodInfo) ([]GPURelation, error) {
	relations := []GPURelation{{
		Machine: &GPURelationMachine{
			ID:       machineID,
			Hostname: hostname,
		},
	}}

	gpuToVGpus := make(map[string][]GPURelationVGPU)
	for _, vgpuPod := range vgpuPods {
		bindTime := vgpuPod.Pod.Annotations["hami.io/bind-time"]
		bindPhase := vgpuPod.Pod.Annotations["hami.io/bind-phase"]
		for i := range vgpuPod.VGpuInfos {
			vgpu := uc.createVGpuEntity(vgpuPod.Pod.Name, &vgpuPod.VGpuInfos[i], bindTime, bindPhase)
			gpuToVGpus[vgpuPod.VGpuInfos[i].GpuUUID] = append(gpuToVGpus[vgpuPod.VGpuInfos[i].GpuUUID], vgpu)
		}
	}
	relations = append(relations, uc.createGPURelationsWithVGpus(machineID, nodeGpuMap, gpuToVGpus)...)
	relations = append(relations, uc.createPodRelations(vgpuPods)...)

	return relations, nil
}

func (uc *EssentialUseCase) createGPURelationsWithVGpus(machineID string, nodeGpuMap map[string]string, gpuToVGpus map[string][]GPURelationVGPU) []GPURelation {
	var relations []GPURelation
	for gpuID, gpuVendorProduct := range nodeGpuMap {
		vendor, product := parseVendorProduct(gpuVendorProduct)
		gpu := GPURelationGPU{
			ID:        gpuID,
			Vendor:    vendor,
			Product:   product,
			MachineID: machineID,
		}
		if vgpus, exists := gpuToVGpus[gpuID]; exists {
			gpu.VGPUs = vgpus
		}
		relations = append(relations, GPURelation{GPU: &gpu})
	}
	return relations
}

func (uc *EssentialUseCase) createPodRelations(vgpuPods []EssentialVGpuPodInfo) []GPURelation {
	var relations []GPURelation
	for _, vgpuPod := range vgpuPods {
		var gpuIDs []string
		for _, vgpuAlloc := range vgpuPod.VGpuInfos {
			gpuIDs = append(gpuIDs, vgpuAlloc.GpuUUID)
		}
		relations = append(relations, GPURelation{
			Pod: &GPURelationPod{
				Name:      vgpuPod.Pod.Name,
				Namespace: vgpuPod.Pod.Namespace,
				ModelName: vgpuPod.ModelName,
				GPUIDs:    gpuIDs,
			},
		})
	}
	return relations
}

func (uc *EssentialUseCase) ListGPURelationsByModel(ctx context.Context, scopeUUID, facilityName, namespace, modelName string) ([]GPURelation, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, scopeUUID, facilityName)
	if err != nil {
		return nil, err
	}

	label := "model-name=" + modelName
	deployments, err := uc.kubeApps.ListDeploymentsByLabel(ctx, config, namespace, label)
	if err != nil {
		return nil, err
	}

	return uc.buildGPURelationsForModel(ctx, config, deployments, modelName)
}

func (uc *EssentialUseCase) createVGpuEntity(podName string, vgpuAlloc *EssentialVGpuAllocation, bindTime, bindPhase string) GPURelationVGPU {
	vgpu := GPURelationVGPU{
		PodName:      podName,
		BindingPhase: bindPhase,
	}
	if vramBytes, err := convertVramMibToBytes(vgpuAlloc.VramMib); err == nil {
		vgpu.VramBytes = vramBytes
	}
	if vcoresPercent, err := convertVcoresPercentToFloat(vgpuAlloc.VcoresPercent); err == nil {
		vgpu.VcoresPercent = vcoresPercent
	}
	if bindTime != "" {
		if bindTimeUnix, err := strconv.ParseInt(bindTime, 10, 64); err == nil {
			vgpu.BoundAt = time.Unix(bindTimeUnix, 0)
		}
	}
	return vgpu
}

func (uc *EssentialUseCase) buildGPURelationsForModel(ctx context.Context, config *rest.Config, deployments []Deployment, modelName string) ([]GPURelation, error) {
	var relations []GPURelation
	tracker := &gpoPodRelationTracker{
		processedPods:     make(map[string]bool),
		processedMachines: make(map[string]bool),
		processedGPUs:     make(map[string]bool),
	}

	for i := range deployments {
		podRelations, err := uc.processDeployment(ctx, config, &deployments[i], modelName, tracker)
		if err != nil {
			continue
		}
		relations = append(relations, podRelations...)
	}
	return relations, nil
}

func (uc *EssentialUseCase) processDeployment(ctx context.Context, config *rest.Config, deployment *Deployment, modelName string, tracker *gpoPodRelationTracker) ([]GPURelation, error) {
	if deployment.Spec.Selector == nil || deployment.Spec.Selector.MatchLabels == nil {
		return nil, nil
	}

	var selector string
	for key, value := range deployment.Spec.Selector.MatchLabels {
		if selector != "" {
			selector += ","
		}
		selector += key + "=" + value
	}

	pods, err := uc.kubeCore.ListPodsByLabel(ctx, config, deployment.Namespace, selector)
	if err != nil {
		return nil, err
	}

	var relations []GPURelation
	for i := range pods {
		podRelations, err := uc.processPod(ctx, config, &pods[i], modelName, tracker)
		if err != nil {
			continue
		}
		relations = append(relations, podRelations...)
	}
	return relations, nil
}

func (uc *EssentialUseCase) processPod(ctx context.Context, config *rest.Config, pod *Pod, modelName string, tracker *gpoPodRelationTracker) ([]GPURelation, error) {
	podKey := pod.Namespace + "/" + pod.Name
	if tracker.processedPods[podKey] {
		return nil, nil
	}
	tracker.processedPods[podKey] = true

	vgpuDevicesAllocated := pod.Annotations["hami.io/vgpu-devices-allocated"]
	if vgpuDevicesAllocated == "" {
		return nil, nil
	}

	vgpuAllocations := parseVGpuDevicesAllocated(vgpuDevicesAllocated)
	if len(vgpuAllocations) == 0 {
		return nil, nil
	}

	machineID, err := uc.getMachineIDFromNodeName(ctx, pod.Spec.NodeName)
	if err != nil {
		return nil, err
	}

	var relations []GPURelation
	if !tracker.processedMachines[machineID] {
		machine, err := uc.machine.Get(ctx, machineID)
		if err == nil {
			tracker.processedMachines[machineID] = true
			relations = append(relations, GPURelation{
				Machine: &GPURelationMachine{
					ID:       machineID,
					Hostname: machine.Hostname,
				},
			})
		}
	}

	node, err := uc.kubeCore.GetNode(ctx, config, pod.Spec.NodeName)
	if err != nil {
		return relations, nil
	}

	nodeGpuMap, err := uc.getNodeGpuMap(node)
	if err != nil {
		nodeGpuMap = make(map[string]string)
	}

	bindTime := pod.Annotations["hami.io/bind-time"]
	bindPhase := pod.Annotations["hami.io/bind-phase"]
	gpuIDs, gpuRelations := uc.processVGpuAllocations(pod.Name, vgpuAllocations, nodeGpuMap, machineID, bindTime, bindPhase, tracker)
	relations = append(relations, gpuRelations...)

	relations = append(relations, GPURelation{
		Pod: &GPURelationPod{
			Name:      pod.Name,
			Namespace: pod.Namespace,
			ModelName: modelName,
			GPUIDs:    gpuIDs,
		},
	})

	return relations, nil
}

func (uc *EssentialUseCase) processVGpuAllocations(podName string, vgpuAllocations []EssentialVGpuAllocation, nodeGpuMap map[string]string, machineID, bindTime, bindPhase string, tracker *gpoPodRelationTracker) ([]string, []GPURelation) {
	var gpuIDs []string
	var relations []GPURelation

	gpuToVGpus := make(map[string][]GPURelationVGPU)
	for i := range vgpuAllocations {
		gpuIDs = append(gpuIDs, vgpuAllocations[i].GpuUUID)
		vgpu := uc.createVGpuEntity(podName, &vgpuAllocations[i], bindTime, bindPhase)
		gpuToVGpus[vgpuAllocations[i].GpuUUID] = append(gpuToVGpus[vgpuAllocations[i].GpuUUID], vgpu)
	}

	for gpuID, vgpus := range gpuToVGpus {
		if !tracker.processedGPUs[gpuID] {
			if gpuVendorProduct, exists := nodeGpuMap[gpuID]; exists {
				vendor, product := parseVendorProduct(gpuVendorProduct)
				relations = append(relations, GPURelation{
					GPU: &GPURelationGPU{
						ID:        gpuID,
						Vendor:    vendor,
						Product:   product,
						MachineID: machineID,
						VGPUs:     vgpus,
					},
				})
				tracker.processedGPUs[gpuID] = true
			}
		}
	}

	return gpuIDs, relations
}

func filterVGpuPodsOnNode(ctx context.Context, pods []Pod, nodeName string, config *rest.Config, _ KubeCoreRepo, kubeApps KubeAppsRepo) []EssentialVGpuPodInfo {
	var vgpuPods []EssentialVGpuPodInfo

	for i := range pods {
		pod := &pods[i]

		if pod.Spec.NodeName != nodeName {
			continue
		}

		// Check if Pod has vGPU allocation annotation
		vgpuDevicesAllocated, vgpuDevicesAllocatedValues := pod.Annotations["hami.io/vgpu-devices-allocated"]
		if !vgpuDevicesAllocatedValues || vgpuDevicesAllocated == "" {
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
	// GPU-uuid,vendor,vram,vcores
	const (
		vgpuUIDIndex    = 0
		vgpuVendorIndex = 1
		vgpuVramIndex   = 2
		vgpuVcoresIndex = 3
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

		allocation := EssentialVGpuAllocation{
			GpuUUID:       parts[vgpuUIDIndex],    // GPU-c15ecdf3-444a-2d02-29e9-e978b2514335
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
	if len(pod.OwnerReferences) > 0 {
		for _, ownerRef := range pod.OwnerReferences {
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

	return ""
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
func (uc *EssentialUseCase) getNodeGpuMap(node *Node) (map[string]string, error) {
	const (
		hamiAnnotationKey = "hami.io/node-nvidia-register"
		gpuUID            = 0
		gpuVender         = 4
	)

	// Parse GPU annotations to get GPU UUID to vendor mapping
	gpuMap := make(map[string]string)

	annotationValue, exists := node.Annotations[hamiAnnotationKey]
	if !exists || annotationValue == "" {
		return gpuMap, nil
	}

	annotationValue = strings.Trim(annotationValue, "'\"")
	nodeNvidiaRegisters := strings.Split(annotationValue, ":")

	for i := 0; i < len(nodeNvidiaRegisters)-1; i++ {
		nodeNvidiaRegister := strings.TrimSpace(nodeNvidiaRegisters[i])
		if nodeNvidiaRegister == "" {
			continue
		}

		// Format: GPU-663aa370-535a-33b8-e01f-b325fb2025c7,10,24564,100,NVIDIA-NVIDIA GeForce RTX 4090,0,true,0,hami-core
		parts := strings.Split(nodeNvidiaRegister, ",")
		cardUID := parts[gpuUID]
		cardVendor := parts[gpuVender]
		gpuMap[cardUID] = cardVendor
	}

	return gpuMap, nil
}

// parseVendorProduct parses vendor and product from GPU vendor string
// Input: "NVIDIA-NVIDIA GeForce RTX 4090"
// Output: vendor="NVIDIA", product="NVIDIA GeForce RTX 4090"
func parseVendorProduct(vendorProduct string) (vendor, product string) {
	const (
		vendorIndex           = 0
		productIndex          = 1
		maxVendorProductParts = 2
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
	const bytesPerMiB = 1024 * 1024 // 1 MiB = 1024 * 1024 bytes
	mib, err := strconv.ParseUint(vramMib, 10, 64)
	if err != nil {
		return 0, err
	}
	return mib * bytesPerMiB, nil
}

// convertVcoresPercentToFloat converts vcores percent from string to float32
func convertVcoresPercentToFloat(vcoresPercent string) (float32, error) {
	percent, err := strconv.ParseFloat(vcoresPercent, 32)
	if err != nil {
		return 0.0, err
	}
	return float32(percent), nil
}
