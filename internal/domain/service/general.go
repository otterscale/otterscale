package service

import (
	"context"
	"slices"
	"strings"

	"golang.org/x/sync/errgroup"

	"github.com/openhdc/openhdc/internal/domain/model"
)

var (
	kubernetesFacilityList = []struct {
		app string
		lxd bool
	}{
		{app: "calico", lxd: true},
		{app: "containerd", lxd: true},
		{app: "easyrsa", lxd: true},
		{app: "etcd", lxd: true},
		{app: "keepalived", lxd: true},
		{app: "kubeapi-load-balancer", lxd: true},
		{app: "kubernetes-control-plane", lxd: true},
		{app: "kubernetes-worker", lxd: false},
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

func (s *NexusService) CreateCeph(ctx context.Context) (*model.FacilityInfo, error) {
	return nil, nil
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

func (s *NexusService) CreateKubernetes(ctx context.Context) (*model.FacilityInfo, error) {
	// lxd: easyrsa, etcd, lb, cp, containerd, calico, keepalived
	// bare machine: worker
	// helm: prometheus-stack
	prefix := "abc"
	prefix = prefix + "-"
	// placements, err := s.toPlacements(ctx, uuid, mps)
	// if err != nil {
	// 	return nil, err
	// }
	// constraint := toConstraint(mc)
	// if _, err := s.facility.Create(ctx, uuid, name, configYAML, charmName, channel, revision, number, placements, &constraint, trust); err != nil {
	// 	return nil, err
	// }
	return nil, nil
}

func (s *NexusService) AddKuberneteUnit(ctx context.Context, uuid, name string, number int, machine string, force bool) error {
	// s.GetApplication(ctx, uuid)
	if force {

	}
	// eg, ctx := errgroup.WithContext(ctx)
	// eg.Go(func() error {
	// 	_, err := s.facility.AddUnits(ctx, uuid, name)
	// 	return err
	// })
	// return eg.Wait()

	return nil
}

func (s *NexusService) appendPrefixToRelationList(prefix string, relationList [][]string) [][]string {
	endpointList := [][]string{}
	for _, relations := range relationList {
		endpoints := []string{}
		for _, relation := range relations {
			endpoints = append(endpoints, (prefix + "-" + relation))
		}
		endpointList = append(endpointList, endpoints)
	}
	return endpointList
}

func (s *NexusService) createRelations(ctx context.Context, uuid string, prefix string, relationList [][]string) error {
	endpointList := s.appendPrefixToRelationList(prefix, relationList)
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
