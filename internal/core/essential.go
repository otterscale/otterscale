package core

import (
	"context"
	"errors"
	"fmt"
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

type EssentialUseCase struct {
	conf           *config.Config
	scope          ScopeRepo
	facility       FacilityRepo
	facilityOffers FacilityOffersRepo
	machine        MachineRepo
	subnet         SubnetRepo
	ipRange        IPRangeRepo
	server         ServerRepo
	client         ClientRepo
	kubeRepo       KubeCoreRepo
	kubeAppsRepo   KubeAppsRepo
	action         ActionRepo
}

// VGpuInfo represents virtual GPU information for a pod
type VGpuInfo struct {
	IsVGpu          bool
	VGpuBindTime    time.Time
	BindPhase       string
	PhysicalGpuUUID string
	VRamMib         string
	VCoresPercent   string
}

// PodInfo represents pod information with GPU details
type PodInfo struct {
	Name             string
	Namespace        string
	ModelName        string
	MachineName      string
	VGpuInfo         []VGpuInfo
	DeploymentName   string
	DeploymentLabels map[string]string
}

// NewEssentialUseCase creates an essential use case
func NewEssentialUseCase(conf *config.Config, scope ScopeRepo, facility FacilityRepo, facilityOffers FacilityOffersRepo, machine MachineRepo, subnet SubnetRepo, ipRange IPRangeRepo, server ServerRepo, client ClientRepo, kubeRepo KubeCoreRepo, kubeAppsRepo KubeAppsRepo, action ActionRepo) *EssentialUseCase {
	return &EssentialUseCase{
		conf:           conf,
		scope:          scope,
		facility:       facility,
		facilityOffers: facilityOffers,
		machine:        machine,
		subnet:         subnet,
		ipRange:        ipRange,
		server:         server,
		client:         client,
		kubeRepo:       kubeRepo,
		kubeAppsRepo:   kubeAppsRepo,
		action:         action,
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
	if err := CreateCeph(ctx, uc.server, uc.machine, uc.facility, uuid, machineID, prefix, cephConfigs); err != nil {
		return err
	}
	if err := CreateKubernetes(ctx, uc.server, uc.machine, uc.facility, uuid, machineID, prefix, kubeConfigs); err != nil {
		return err
	}
	if err := CreateCommon(ctx, uc.server, uc.machine, uc.facility, uc.facilityOffers, uc.conf, uuid, prefix, commonConfigs); err != nil {
		return err
	}
	return nil
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

func createEssential(ctx context.Context, serverRepo ServerRepo, machineRepo MachineRepo, facilityRepo FacilityRepo, uuid, machineID, prefix string, charms []EssentialCharm, configs map[string]string) error {
	var (
		directive string
		err       error
	)
	if machineID != "" {
		directive, err = getDirective(ctx, machineRepo, machineID)
		if err != nil {
			return err
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

// GetGpuRelationByMachine returns GPU relation information by machine name
func (uc *EssentialUseCase) GetGpuRelationByMachine(ctx context.Context, uuid, facility, machineName string) ([]PodInfo, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}

	// Find pods on the specified machine with vGPU
	label := "hami.io/vgpu-node=" + machineName
	pods, err := uc.kubeRepo.ListPodsByLabel(ctx, config, "", label)
	if err != nil {
		return nil, err
	}

	podInfos := make([]PodInfo, 0, len(pods))
	for i := range pods {
		pod := &pods[i]
		// Check if pod has vGPU annotations
		annotations := pod.GetAnnotations()
		if _, hasVGpu := annotations["hami.io/vgpu-devices-allocated"]; !hasVGpu {
			continue // Skip pods without vGPU
		}

		podInfo := PodInfo{
			Name:        pod.GetName(),
			Namespace:   pod.GetNamespace(),
			MachineName: machineName,
		}

		// Get machine name from labels if available
		if nodeLabel, ok := pod.GetLabels()["hami.io/vgpu-node"]; ok {
			podInfo.MachineName = nodeLabel
		}

		// Find deployment for this pod
		deployment, err := uc.findDeploymentForPod(ctx, config, pod)
		if err == nil && deployment != nil {
			podInfo.DeploymentName = deployment.Name
			podInfo.DeploymentLabels = deployment.Labels

			// Get model name from deployment labels
			if modelName, ok := deployment.Labels["model-name"]; ok {
				podInfo.ModelName = modelName
			}
		}

		// Extract GPU information from annotations
		gpuInfos := extractGpuInfoFromPodAnnotations(annotations)
		if len(gpuInfos) > 0 {
			podInfo.VGpuInfo = gpuInfos
		}

		podInfos = append(podInfos, podInfo)
	}

	return podInfos, nil
}

// GetGpuRelationByModel returns GPU relation information by model name
func (uc *EssentialUseCase) GetGpuRelationByModel(ctx context.Context, uuid, facility, modelName string) ([]PodInfo, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}

	// Find deployments with the specified model
	label := "model-name=" + modelName
	deployments, err := uc.kubeAppsRepo.ListDeploymentsByLabel(ctx, config, "", label)
	if err != nil {
		return nil, err
	}

	podInfos := []PodInfo{}
	for i := range deployments {
		deployment := &deployments[i]
		// Get pods for this deployment
		selector := deployment.Spec.Selector
		if selector == nil || selector.MatchLabels == nil {
			continue
		}

		// Build label selector from deployment's selector
		podLabelSelector := ""
		for key, value := range selector.MatchLabels {
			if podLabelSelector != "" {
				podLabelSelector += ","
			}
			podLabelSelector += key + "=" + value
		}

		pods, err := uc.kubeRepo.ListPodsByLabel(ctx, config, deployment.Namespace, podLabelSelector)
		if err != nil {
			continue
		}

		for j := range pods {
			pod := &pods[j]
			podInfo := PodInfo{
				Name:             pod.GetName(),
				Namespace:        pod.GetNamespace(),
				ModelName:        modelName,
				DeploymentName:   deployment.Name,
				DeploymentLabels: deployment.Labels,
			}

			// Get machine name from labels
			if machineName, ok := pod.GetLabels()["hami.io/vgpu-node"]; ok {
				podInfo.MachineName = machineName
			}

			// Extract GPU information from annotations
			gpuInfos := extractGpuInfoFromPodAnnotations(pod.GetAnnotations())
			if len(gpuInfos) > 0 {
				podInfo.VGpuInfo = gpuInfos
			}

			podInfos = append(podInfos, podInfo)
		}
	}

	return podInfos, nil
}

// extractGpuInfoFromPodAnnotations extracts GPU information from pod annotations
func extractGpuInfoFromPodAnnotations(annotations map[string]string) []VGpuInfo {
	vgpuDevicesAllocated, hasVGpu := annotations["hami.io/vgpu-devices-allocated"]
	if !hasVGpu || vgpuDevicesAllocated == "" {
		return nil
	}

	vgpuInfos := []VGpuInfo{}
	devices := strings.Split(vgpuDevicesAllocated, ":")

	const minDevicePartsCount = 4
	for _, device := range devices {
		if device == "" {
			continue
		}

		parts := strings.Split(device, ",")
		if len(parts) < minDevicePartsCount {
			continue
		}

		// Extract UUID, VRAM, and cores from the device string
		gpuUUID := parts[0]
		vram := parts[2]

		// Extract cores (removing trailing colon)
		coresStr := parts[3]
		cores := strings.TrimSuffix(coresStr, ":")

		// Parse bind time
		var bindTime time.Time
		if bindTimeStr, ok := annotations["hami.io/bind-time"]; ok && bindTimeStr != "" {
			if timestamp, err := strconv.ParseInt(bindTimeStr, 10, 64); err == nil {
				bindTime = time.Unix(timestamp, 0)
			}
		}

		vgpuInfo := VGpuInfo{
			IsVGpu:          true,
			VGpuBindTime:    bindTime,
			BindPhase:       annotations["hami.io/bind-phase"],
			PhysicalGpuUUID: gpuUUID,
			VRamMib:         vram,
			VCoresPercent:   cores,
		}

		vgpuInfos = append(vgpuInfos, vgpuInfo)
	}

	return vgpuInfos
}

// findDeploymentForPod finds the deployment that owns the given pod
func (uc *EssentialUseCase) findDeploymentForPod(ctx context.Context, config *rest.Config, pod *Pod) (*Deployment, error) {
	// Get the owner references from the pod
	ownerRefs := pod.OwnerReferences
	for _, ownerRef := range ownerRefs {
		if ownerRef.Kind == "ReplicaSet" {
			// Find the ReplicaSet
			replicaSetName := ownerRef.Name

			// Get all deployments in the namespace and check if any owns this ReplicaSet
			deployments, err := uc.kubeAppsRepo.ListDeployments(ctx, config, pod.Namespace)
			if err != nil {
				return nil, err
			}

			for i := range deployments {
				deployment := &deployments[i]
				// Check if this deployment's name matches the ReplicaSet prefix
				deploymentName := deployment.Name
				if strings.HasPrefix(replicaSetName, deploymentName+"-") {
					return deployment, nil
				}
			}
		}
	}
	return nil, nil
}
