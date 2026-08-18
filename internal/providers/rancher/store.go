// Package rancher provides Rancher Project discovery from the
// management cluster's informer cache.
package rancher

import (
	"cmp"
	"context"
	"log/slog"
	"slices"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"

	"github.com/otterscale/otterscale/internal/core"
	"github.com/otterscale/otterscale/internal/transport"
)

var projectGVR = schema.GroupVersionResource{
	Group:    "management.cattle.io",
	Version:  "v3",
	Resource: "projects",
}

// Store exposes valid Rancher Projects from a shared informer's Indexer.
type Store struct {
	informer cache.SharedIndexInformer
}

// NewUnavailableStore returns a Store that remains cache-not-ready.
// It is used when management-cluster client initialization fails so
// unrelated server features can continue running.
func NewUnavailableStore() *Store {
	return &Store{}
}

// NewStore constructs a cluster-wide Project informer using the server
// Pod's management-cluster ServiceAccount credentials.
func NewStore(config *rest.Config) (*Store, error) {
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(
		client,
		0*time.Second,
		metav1.NamespaceAll,
		nil,
	)
	informer := factory.ForResource(projectGVR).Informer()
	if err := informer.SetWatchErrorHandlerWithContext(func(ctx context.Context, _ *cache.Reflector, err error) {
		slog.ErrorContext(ctx, "Rancher Project informer error", "error", err)
	}); err != nil {
		return nil, err
	}
	return newStore(informer), nil
}

func newStore(informer cache.SharedIndexInformer) *Store {
	return &Store{informer: informer}
}

var (
	_ core.RancherProjectStore = (*Store)(nil)
	_ transport.Listener       = (*Store)(nil)
)

// Start runs the informer until the server lifecycle context is canceled.
func (s *Store) Start(ctx context.Context) error {
	if s.informer == nil {
		<-ctx.Done()
		return nil
	}
	s.informer.RunWithContext(ctx)
	return nil
}

// Stop is a no-op because the informer stops through its Start context.
func (s *Store) Stop(context.Context) error { return nil }

// ListProjects returns a stable snapshot from the informer store.
func (s *Store) ListProjects(context.Context) ([]core.RancherProject, error) {
	if s.informer == nil || !s.informer.HasSynced() {
		return nil, core.ErrRancherProjectCacheNotReady
	}

	items := s.informer.GetStore().List()
	projects := make([]core.RancherProject, 0, len(items))
	for _, item := range items {
		project, ok := projectFromObject(item)
		if ok {
			projects = append(projects, project)
		}
	}
	slices.SortFunc(projects, func(a, b core.RancherProject) int {
		return cmp.Or(cmp.Compare(a.DisplayName, b.DisplayName), cmp.Compare(a.ID, b.ID))
	})
	return projects, nil
}

// HasProject checks informer membership without issuing a live API GET.
func (s *Store) HasProject(_ context.Context, id string) (bool, error) {
	if s.informer == nil || !s.informer.HasSynced() {
		return false, core.ErrRancherProjectCacheNotReady
	}
	clusterID, projectName, err := core.ParseRancherProjectID(id)
	if err != nil {
		return false, err
	}
	item, exists, err := s.informer.GetIndexer().GetByKey(clusterID + "/" + projectName)
	if err != nil || !exists {
		return false, err
	}
	project, valid := projectFromObject(item)
	return valid && project.ID == id, nil
}

func projectFromObject(item any) (core.RancherProject, bool) {
	object, ok := item.(*unstructured.Unstructured)
	if !ok || object.GetName() == "" || object.GetNamespace() == "" || object.GetDeletionTimestamp() != nil {
		return core.RancherProject{}, false
	}
	if len(validation.IsDNS1123Label(object.GetNamespace())) > 0 || len(validation.IsDNS1123Subdomain(object.GetName())) > 0 {
		return core.RancherProject{}, false
	}
	clusterName, _, _ := unstructured.NestedString(object.Object, "spec", "clusterName")
	if clusterName != object.GetNamespace() {
		return core.RancherProject{}, false
	}
	displayName, _, _ := unstructured.NestedString(object.Object, "spec", "displayName")
	if displayName == "" {
		displayName = object.GetName()
	}
	return core.RancherProject{
		ID:          object.GetNamespace() + ":" + object.GetName(),
		DisplayName: displayName,
	}, true
}
