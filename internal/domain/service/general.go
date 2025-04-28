package service

import (
	"context"
	"slices"
	"strings"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/canonical/gomaasclient/entity/node"
	"github.com/juju/juju/core/instance"

	"github.com/openhdc/openhdc/internal/domain/model"
)

type generalFacility struct {
	charmName string
	lxd       bool
}

var (
	kubernetesFacilityList = []generalFacility{
		{charmName: "ch:calico", lxd: true},
		{charmName: "ch:containerd", lxd: true},
		{charmName: "ch:easyrsa", lxd: true},
		{charmName: "ch:etcd", lxd: true},
		{charmName: "ch:keepalived", lxd: true},
		{charmName: "ch:kubeapi-load-balancer", lxd: true},
		{charmName: "ch:kubernetes-control-plane", lxd: true},
		{charmName: "ch:kubernetes-worker", lxd: false},
	}

	kubernetesRelationList = [][]string{
		{"calico:cni", "kubernetes-control-plane:cni"},
		{"calico:etcd", "etcd:db"},
		{"containerd:containerd", "kubernetes-control-plane:container-runtime"},
		{"containerd:containerd", "kubernetes-worker:container-runtime"},
		{"etcd:certificates", "easyrsa:client"},
		{"keepalived:website", "kubeapi-load-balancer:apiserver"},
		{"kubeapi-load-balancer:certificates", "easyrsa:client"},
		{"kubernetes-control-plane:certificates", "easyrsa:client"},
		{"kubernetes-control-plane:etcd", "etcd:db"},
		{"kubernetes-control-plane:kube-control", "kubernetes-worker:kube-control"},
		{"kubernetes-control-plane:kube-api-endpoint", "kubeapi-load-balancer:apiserver"},
		{"kubernetes-control-plane:loadbalancer", "kubeapi-load-balancer:loadbalancer"},
		{"kubernetes-worker:certificates", "easyrsa:client"},
		{"kubernetes-worker:kube-api-endpoint", " kubeapi-load-balancer:website"},
	}
)

var (
	cephFacilityList = []generalFacility{
		{charmName: "ch:ceph-fs", lxd: true},
		{charmName: "ch:ceph-mon", lxd: true},
		{charmName: "ch:ceph-osd", lxd: false},
	}

	cephRelationList = [][]string{
		{"ceph-fs:ceph-mds", "ceph-mon:mds"},
		{"ceph-osd:mon", "ceph-mon:osd"},
	}
)

func (s *NexusService) VerifyEnvironment(ctx context.Context) ([]model.Error, error) {
	funcs := []func(context.Context) (*model.Error, error){}
	funcs = append(funcs, s.isCephExists, s.isKubernetesExists)

	eg, ctx := errgroup.WithContext(ctx)
	result := make([]model.Error, len(funcs))
	for i := range funcs {
		eg.Go(func() error {
			e, err := funcs[i](ctx)
			if err == nil && e != nil {
				result[i] = *e
			}
			return err
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	slices.SortFunc(result, func(e1, e2 model.Error) int {
		return strings.Compare(e1.Code, e2.Code)
	})
	return slices.DeleteFunc(result, func(e model.Error) bool { return e.Code == "" }), nil
}

func (s *NexusService) ListCephes(ctx context.Context, uuid string) ([]model.FacilityInfo, error) {
	fis, err := s.listFacilitiesAcrossScopes(ctx, charmNameCeph)
	if err != nil {
		return nil, err
	}
	filter := []model.FacilityInfo{}
	for i := range fis {
		if strings.Contains(fis[i].ScopeUUID, uuid) {
			filter = append(filter, fis[i])
		}
	}
	return filter, nil
}

// TODO: CONFIG
func (s *NexusService) CreateCeph(ctx context.Context, uuid, machineID, prefix string) (*model.FacilityInfo, error) {
	fi, err := s.createGeneralFacility(ctx, uuid, machineID, prefix, charmNameCeph, cephFacilityList)
	if err != nil {
		return nil, err
	}
	if err := s.createGeneralRelations(ctx, uuid, prefix, cephRelationList); err != nil {
		return nil, err
	}
	return fi, nil
}

func (s *NexusService) AddCephUnits(ctx context.Context, uuid, general string, number int, machineIDs []string) error {
	return s.addGeneralFacilityUnits(ctx, uuid, general, number, machineIDs, cephFacilityList)
}

func (s *NexusService) ListKuberneteses(ctx context.Context, uuid string) ([]model.FacilityInfo, error) {
	fis, err := s.listFacilitiesAcrossScopes(ctx, charmNameKubernetes)
	if err != nil {
		return nil, err
	}
	filter := []model.FacilityInfo{}
	for i := range fis {
		if strings.Contains(fis[i].ScopeUUID, uuid) {
			filter = append(filter, fis[i])
		}
	}
	return filter, nil
}

// TODO: CONFIG
func (s *NexusService) CreateKubernetes(ctx context.Context, uuid, machineID, prefix string) (*model.FacilityInfo, error) {
	fi, err := s.createGeneralFacility(ctx, uuid, machineID, prefix, charmNameKubernetes, kubernetesFacilityList)
	if err != nil {
		return nil, err
	}
	if err := s.createGeneralRelations(ctx, uuid, prefix, kubernetesRelationList); err != nil {
		return nil, err
	}
	return fi, nil
}

func (s *NexusService) AddKubernetesUnits(ctx context.Context, uuid, general string, number int, machineIDs []string, force bool) error {
	if !force {
		st, err := s.client.Status(ctx, uuid, []string{"application", general})
		if err != nil {
			return err
		}
		app, ok := st.Applications[general]
		if !ok {
			return status.Errorf(codes.NotFound, "kubernetes facility %q not found", general)
		}
		if len(app.Units) > 3 {
			return status.Errorf(codes.InvalidArgument, "cannot add more than 3 Kubernetes worker units without force flag")
		}
	}
	return s.addGeneralFacilityUnits(ctx, uuid, general, number, machineIDs, kubernetesFacilityList)
}

func (s *NexusService) createGeneralFacility(ctx context.Context, uuid, machineID, prefix, general string, facilityList []generalFacility) (*model.FacilityInfo, error) {
	m, err := s.machine.Get(ctx, machineID)
	if err != nil {
		return nil, err
	}
	if m.Status != node.StatusDeployed {
		return nil, status.Error(codes.InvalidArgument, "machine status is not deployed")
	}

	directive, err := getJujuMachineID(m.WorkloadAnnotations)
	if err != nil {
		return nil, err
	}

	base, err := s.imageBase(ctx)
	if err != nil {
		return nil, err
	}

	var facilityName string
	eg, ctx := errgroup.WithContext(ctx)
	for _, facility := range facilityList {
		name := toGeneralFacilityName(prefix, facility.charmName)

		if facility.charmName == general {
			facilityName = name
		}

		eg.Go(func() error {
			placements := []instance.Placement{
				{Scope: toPlacementScope(facility.lxd), Directive: directive},
			}
			_, err := s.facility.Create(ctx, uuid, name, "", facility.charmName, "", 0, 1, base, placements, nil, true)
			return err
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	scopeName, err := s.getScopeName(ctx, uuid)
	if err != nil {
		return nil, err
	}

	return &model.FacilityInfo{
		ScopeUUID:    uuid,
		ScopeName:    scopeName,
		FacilityName: facilityName,
	}, nil
}

func (s *NexusService) addGeneralFacilityUnits(ctx context.Context, uuid, general string, number int, machineIDs []string, facilityList []generalFacility) error {
	slices.Sort(machineIDs)
	directives := slices.Compact(machineIDs)
	if len(directives) != number {
		return status.Error(codes.InvalidArgument, "number of machines does not match requested number of units")
	}

	prefix := toGeneralFacilityPrefix(general)

	eg, ctx := errgroup.WithContext(ctx)
	for _, facility := range facilityList {
		name := toGeneralFacilityName(prefix, facility.charmName)
		lxd := facility.lxd

		eg.Go(func() error {
			placements := make([]instance.Placement, len(directives))
			for i, directive := range directives {
				placements[i] = instance.Placement{
					Scope:     toPlacementScope(lxd),
					Directive: directive,
				}
			}

			_, err := s.facility.AddUnits(ctx, uuid, name, number, placements)
			return err
		})
	}

	return eg.Wait()
}

func (s *NexusService) getScopeName(ctx context.Context, uuid string) (string, error) {
	scopes, err := s.scope.List(ctx)
	if err != nil {
		return "", err
	}

	for i := range scopes {
		if scopes[i].UUID == uuid {
			return scopes[i].Name, nil
		}
	}

	return "", nil
}

func appendPrefixToRelationList(prefix string, relationList [][]string) [][]string {
	endpointList := [][]string{}
	for _, relations := range relationList {
		endpoints := []string{}
		for _, relation := range relations {
			endpoints = append(endpoints, toGeneralFacilityName(prefix, relation))
		}
		endpointList = append(endpointList, endpoints)
	}
	return endpointList
}

func (s *NexusService) createGeneralRelations(ctx context.Context, uuid, prefix string, relationList [][]string) error {
	endpointList := appendPrefixToRelationList(prefix, relationList)
	eg, ctx := errgroup.WithContext(ctx)
	for _, endpoints := range endpointList {
		eg.Go(func() error {
			_, err := s.facility.CreateRelation(ctx, uuid, endpoints)
			return err
		})
	}
	return eg.Wait()
}

func (s *NexusService) isCephExists(ctx context.Context) (*model.Error, error) {
	cephes, err := s.ListCephes(ctx, "")
	if err != nil {
		return nil, err
	}
	if len(cephes) == 0 {
		return &model.ErrCephNotFound, nil
	}
	return nil, nil
}

func (s *NexusService) isKubernetesExists(ctx context.Context) (*model.Error, error) {
	kubernetes, err := s.ListKuberneteses(ctx, "")
	if err != nil {
		return nil, err
	}
	if len(kubernetes) == 0 {
		return &model.ErrKubernetesNotFound, nil
	}
	return nil, nil
}

func toPlacementScope(lxd bool) string {
	if lxd {
		return "lxd"
	}
	return instance.MachineScope
}

func toGeneralFacilityName(prefix, charmName string) string {
	if strings.HasPrefix(charmName, "ch:") {
		return prefix + "-" + strings.Split(charmName, ":")[1]
	}
	return prefix + "-" + charmName
}

func toGeneralFacilityPrefix(general string) string {
	return strings.Split(general, "-")[0]
}
