package core

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"connectrpc.com/connect"
	"github.com/canonical/gomaasclient/entity/node"
	"github.com/juju/juju/core/instance"
	jujustatus "github.com/juju/juju/core/status"
	"github.com/juju/juju/rpc/params"
	"golang.org/x/sync/errgroup"
	jujuyaml "gopkg.in/yaml.v2"

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

type Capability struct {
	Category    string
	Name        string
	Description string
	Features    []string
	Available   bool
}

type CapabilitiesResponse struct {
	PlatformName        string
	PlatformDescription string
	Capabilities        []Capability
	UseCases            []string
	DocumentationURL    string
	Version             string
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
}

func NewEssentialUseCase(conf *config.Config, scope ScopeRepo, facility FacilityRepo, facilityOffers FacilityOffersRepo, machine MachineRepo, subnet SubnetRepo, ipRange IPRangeRepo, server ServerRepo, client ClientRepo) *EssentialUseCase {
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

func (uc *EssentialUseCase) GetCapabilities(ctx context.Context, language string) (*CapabilitiesResponse, error) {
	// Determine language for localization
	isZh := strings.HasPrefix(strings.ToLower(language), "zh")
	
	capabilities := []Capability{
		{
			Category:    ternary(isZh, "虛擬化管理", "Virtualization Management"),
			Name:        ternary(isZh, "虛擬機生命週期管理", "VM Lifecycle Management"),
			Description: ternary(isZh, "創建、啟動、停止、暫停、遷移虛擬機", "Create, start, stop, pause, migrate virtual machines"),
			Features:    ternary(isZh, []string{"KVM/QEMU集成", "GPU直通", "熱遷移", "快照管理"}, []string{"KVM/QEMU Integration", "GPU Passthrough", "Live Migration", "Snapshot Management"}),
			Available:   true,
		},
		{
			Category:    ternary(isZh, "容器編排", "Container Orchestration"),
			Name:        ternary(isZh, "Kubernetes原生支援", "Kubernetes Native Support"),
			Description: ternary(isZh, "部署和管理容器化應用程序", "Deploy and manage containerized applications"),
			Features:    ternary(isZh, []string{"Juju Charm部署", "工作負載管理", "服務網格", "自動擴展"}, []string{"Juju Charm Deployment", "Workload Management", "Service Mesh", "Auto Scaling"}),
			Available:   true,
		},
		{
			Category:    ternary(isZh, "存儲服務", "Storage Services"),
			Name:        ternary(isZh, "分佈式存儲", "Distributed Storage"),
			Description: ternary(isZh, "基於Ceph的可擴展塊、對象和文件存儲", "Ceph-based scalable block, object, and file storage"),
			Features:    ternary(isZh, []string{"S3兼容對象存儲", "高性能塊存儲", "POSIX文件系統", "備份與恢復"}, []string{"S3-Compatible Object Storage", "High-Performance Block Storage", "POSIX File Systems", "Backup & Recovery"}),
			Available:   true,
		},
		{
			Category:    ternary(isZh, "網絡", "Networking"),
			Name:        ternary(isZh, "軟件定義網絡", "Software-Defined Networking"),
			Description: ternary(isZh, "虛擬網絡、子網和路由", "Virtual networks, subnets, and routing"),
			Features:    ternary(isZh, []string{"負載均衡", "防火牆管理", "VPN集成", "網絡隔離"}, []string{"Load Balancing", "Firewall Management", "VPN Integration", "Network Isolation"}),
			Available:   true,
		},
		{
			Category:    ternary(isZh, "基礎設施管理", "Infrastructure Management"),
			Name:        ternary(isZh, "裸機配置", "Bare Metal Provisioning"),
			Description: ternary(isZh, "MAAS集成進行物理服務器管理", "MAAS integration for physical server management"),
			Features:    ternary(isZh, []string{"資源分配", "高可用性", "自動故障轉移", "水平擴展"}, []string{"Resource Allocation", "High Availability", "Automatic Failover", "Horizontal Scaling"}),
			Available:   true,
		},
		{
			Category:    ternary(isZh, "監控與可觀測性", "Monitoring & Observability"),
			Name:        ternary(isZh, "全面監控", "Comprehensive Monitoring"),
			Description: ternary(isZh, "基於Prometheus的監控和Grafana可視化", "Prometheus-based monitoring with Grafana visualization"),
			Features:    ternary(isZh, []string{"指標收集", "告警系統", "日誌聚合", "分佈式追蹤"}, []string{"Metrics Collection", "Alerting System", "Log Aggregation", "Distributed Tracing"}),
			Available:   true,
		},
		{
			Category:    ternary(isZh, "安全與訪問控制", "Security & Access Control"),
			Name:        ternary(isZh, "企業級安全", "Enterprise Security"),
			Description: ternary(isZh, "基於角色的訪問控制和企業認證", "Role-based access control and enterprise authentication"),
			Features:    ternary(isZh, []string{"RBAC", "LDAP/AD集成", "單點登錄", "數據加密", "審計日誌"}, []string{"RBAC", "LDAP/AD Integration", "Single Sign-On", "Data Encryption", "Audit Logging"}),
			Available:   true,
		},
		{
			Category:    ternary(isZh, "應用市場", "Application Marketplace"),
			Name:        ternary(isZh, "精選應用程序", "Curated Applications"),
			Description: ternary(isZh, "預配置的應用程序，可立即部署", "Pre-configured applications ready for deployment"),
			Features:    ternary(isZh, []string{"Charm商店", "自定義應用程序", "應用程序生命週期", "一鍵部署"}, []string{"Charm Store", "Custom Applications", "Application Lifecycle", "One-Click Deployment"}),
			Available:   true,
		},
		{
			Category:    ternary(isZh, "API與集成", "API & Integration"),
			Name:        ternary(isZh, "全面API支援", "Comprehensive API Support"),
			Description: ternary(isZh, "REST API和gRPC服務覆蓋所有平台操作", "REST API and gRPC services for all platform operations"),
			Features:    ternary(isZh, []string{"RESTful API", "gRPC服務", "CLI工具", "Webhook支持", "Terraform提供商"}, []string{"RESTful API", "gRPC Services", "CLI Tools", "Webhook Support", "Terraform Provider"}),
			Available:   true,
		},
	}

	useCases := ternary(isZh,
		[]string{
			"企業數據中心 - 多租戶基礎設施和資源優化",
			"開發與測試 - CI/CD集成和環境配置",
			"邊緣計算 - 分佈式部署和本地處理",
			"雲遷移 - 混合雲和工作負載遷移",
		},
		[]string{
			"Enterprise Data Centers - Multi-tenant infrastructure and resource optimization",
			"Development & Testing - CI/CD integration and environment provisioning",
			"Edge Computing - Distributed deployment and local processing",
			"Cloud Migration - Hybrid cloud and workload migration",
		},
	)

	return &CapabilitiesResponse{
		PlatformName:        "OtterScale",
		PlatformDescription: ternary(isZh, "統一基礎設施，賦能創新 - 超融合基礎設施平台", "Unifying Infrastructure, Empowering Innovation - Hyper-Converged Infrastructure Platform"),
		Capabilities:        capabilities,
		UseCases:            useCases,
		DocumentationURL:    "https://otterscale.github.io",
		Version:             "v0.6.0",
	}, nil
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

// Helper function for ternary operations (Go doesn't have ternary operator)
func ternary[T any](condition bool, trueVal T, falseVal T) T {
	if condition {
		return trueVal
	}
	return falseVal
}
