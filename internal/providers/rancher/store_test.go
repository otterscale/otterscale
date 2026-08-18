package rancher

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic/dynamicinformer"
	dynamicfake "k8s.io/client-go/dynamic/fake"
	k8stesting "k8s.io/client-go/testing"
	"k8s.io/client-go/tools/cache"

	"github.com/otterscale/otterscale/internal/core"
)

//nolint:gocyclo,funlen // One lifecycle keeps informer state shared across all watch assertions.
func TestStoreInformerLifecycle(t *testing.T) {
	deleting := project("c-m-delete", "p-deleting", "Deleting")
	now := metav1.Now()
	deleting.SetDeletionTimestamp(&now)
	objects := []runtime.Object{
		project("c-m-prod", "p-zebra", "Zebra"),
		project("local", "p-alpha", "Alpha"),
		project("c-m-fallback", "p-fallback", ""),
		projectWithClusterName("c-m-wrong", "p-wrong", "other", "Wrong"),
		projectWithClusterName("bad.cluster", "p-invalid", "bad.cluster", "Invalid"),
		deleting,
	}
	client := dynamicfake.NewSimpleDynamicClientWithCustomListKinds(
		runtime.NewScheme(),
		map[schema.GroupVersionResource]string{projectGVR: "ProjectList"},
		objects...,
	)
	watches := make(chan watch.Interface, 2)
	client.PrependWatchReactor("projects", func(action k8stesting.Action) (bool, watch.Interface, error) {
		watcher, err := client.Tracker().Watch(action.GetResource(), action.GetNamespace())
		if err == nil {
			watches <- watcher
		}
		return true, watcher, err
	})
	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(client, 0, metav1.NamespaceAll, nil)
	informer := factory.ForResource(projectGVR).Informer()
	store := newStore(informer)

	if _, err := store.ListProjects(t.Context()); !errors.Is(err, core.ErrRancherProjectCacheNotReady) {
		t.Fatalf("ListProjects before sync error = %v", err)
	}
	if _, err := store.HasProject(t.Context(), "local:p-alpha"); !errors.Is(err, core.ErrRancherProjectCacheNotReady) {
		t.Fatalf("HasProject before sync error = %v", err)
	}

	runCtx, cancel := context.WithCancel(t.Context())
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- store.Start(runCtx) }()
	ctx, stopWaiting := context.WithTimeout(t.Context(), 5*time.Second)
	defer stopWaiting()
	if !cache.WaitForCacheSync(ctx.Done(), informer.HasSynced) {
		cancel()
		<-done
		t.Fatal("informer did not sync")
	}
	var firstWatch watch.Interface
	select {
	case firstWatch = <-watches:
	case <-ctx.Done():
		t.Fatal("informer watch did not start")
	}

	projects, err := store.ListProjects(ctx)
	if err != nil {
		t.Fatalf("ListProjects: %v", err)
	}
	want := []core.RancherProject{
		{ID: "local:p-alpha", DisplayName: "Alpha"},
		{ID: "c-m-prod:p-zebra", DisplayName: "Zebra"},
		{ID: "c-m-fallback:p-fallback", DisplayName: "p-fallback"},
	}
	if len(projects) != len(want) {
		t.Fatalf("projects = %#v, want %#v", projects, want)
	}
	for i := range want {
		if projects[i] != want[i] {
			t.Fatalf("projects[%d] = %#v, want %#v", i, projects[i], want[i])
		}
	}
	if has, err := store.HasProject(ctx, "local:p-alpha"); err != nil || !has {
		t.Fatalf("HasProject(local:p-alpha) = %v, %v", has, err)
	}
	if has, err := store.HasProject(ctx, "c-m-wrong:p-wrong"); err != nil || has {
		t.Fatalf("mismatched Project must be excluded: %v, %v", has, err)
	}

	added := project("local", "p-added", "Beta")
	resource := client.Resource(projectGVR).Namespace("local")
	if _, err := resource.Create(ctx, added, metav1.CreateOptions{}); err != nil {
		t.Fatalf("create Project: %v", err)
	}
	waitForProject(ctx, t, store, "local:p-added", true)

	updated, err := resource.Get(ctx, "p-added", metav1.GetOptions{})
	if err != nil {
		t.Fatalf("get Project: %v", err)
	}
	_ = unstructured.SetNestedField(updated.Object, "Aardvark", "spec", "displayName")
	if _, err := resource.Update(ctx, updated, metav1.UpdateOptions{}); err != nil {
		t.Fatalf("update Project: %v", err)
	}
	waitForDisplayName(ctx, t, store, "local:p-added", "Aardvark")

	if err := resource.Delete(ctx, "p-added", metav1.DeleteOptions{}); err != nil {
		t.Fatalf("delete Project: %v", err)
	}
	waitForProject(ctx, t, store, "local:p-added", false)

	// A warm cache remains readable while the informer reconnects after
	// an interrupted watch.
	firstWatch.Stop()
	select {
	case <-watches:
	case <-ctx.Done():
		t.Fatal("informer did not restart watch")
	}
	projects, err = store.ListProjects(ctx)
	if err != nil || len(projects) != len(want) {
		t.Fatalf("warm cache after watch stop = %#v, %v", projects, err)
	}
	cancel()
	<-done
}

func TestStoreRetriesInitialListFailure(t *testing.T) {
	client := dynamicfake.NewSimpleDynamicClientWithCustomListKinds(
		runtime.NewScheme(),
		map[schema.GroupVersionResource]string{projectGVR: "ProjectList"},
		project("local", "p-test", "Test"),
	)
	var attempts atomic.Int32
	client.PrependReactor("list", "projects", func(k8stesting.Action) (bool, runtime.Object, error) {
		if attempts.Add(1) == 1 {
			return true, nil, errors.New("initial list failed")
		}
		return false, nil, nil
	})
	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(client, 0, metav1.NamespaceAll, nil)
	store := newStore(factory.ForResource(projectGVR).Informer())

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	done := make(chan error, 1)
	go func() { done <- store.Start(ctx) }()
	if !cache.WaitForCacheSync(ctx.Done(), store.informer.HasSynced) {
		cancel()
		<-done
		t.Fatal("informer did not recover from initial list failure")
	}
	if attempts.Load() < 2 {
		t.Fatalf("list attempts = %d, want at least 2", attempts.Load())
	}
	if has, err := store.HasProject(ctx, "local:p-test"); err != nil || !has {
		t.Fatalf("HasProject after retry = %v, %v", has, err)
	}
	cancel()
	<-done
}

func TestUnavailableStore(t *testing.T) {
	store := NewUnavailableStore()
	if _, err := store.ListProjects(t.Context()); !errors.Is(err, core.ErrRancherProjectCacheNotReady) {
		t.Fatalf("ListProjects error = %v", err)
	}
	if _, err := store.HasProject(t.Context(), "local:p-test"); !errors.Is(err, core.ErrRancherProjectCacheNotReady) {
		t.Fatalf("HasProject error = %v", err)
	}

	ctx, cancel := context.WithCancel(t.Context())
	done := make(chan error, 1)
	go func() { done <- store.Start(ctx) }()
	select {
	case err := <-done:
		t.Fatalf("Start returned before cancellation: %v", err)
	default:
	}
	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Start returned error: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("Start did not stop after cancellation")
	}
}

func project(namespace, name, displayName string) *unstructured.Unstructured {
	return projectWithClusterName(namespace, name, namespace, displayName)
}

func projectWithClusterName(namespace, name, clusterName, displayName string) *unstructured.Unstructured {
	return &unstructured.Unstructured{Object: map[string]any{
		"apiVersion": "management.cattle.io/v3",
		"kind":       "Project",
		"metadata": map[string]any{
			"name":      name,
			"namespace": namespace,
		},
		"spec": map[string]any{
			"clusterName": clusterName,
			"displayName": displayName,
		},
	}}
}

func waitForProject(ctx context.Context, t *testing.T, store *Store, id string, want bool) {
	t.Helper()
	err := wait.PollUntilContextTimeout(ctx, 10*time.Millisecond, time.Second, true, func(context.Context) (bool, error) {
		has, err := store.HasProject(ctx, id)
		return err == nil && has == want, err
	})
	if err != nil {
		t.Fatalf("wait for Project %q=%v: %v", id, want, err)
	}
}

func waitForDisplayName(ctx context.Context, t *testing.T, store *Store, id, want string) {
	t.Helper()
	err := wait.PollUntilContextTimeout(ctx, 10*time.Millisecond, time.Second, true, func(context.Context) (bool, error) {
		projects, err := store.ListProjects(ctx)
		if err != nil {
			return false, err
		}
		for _, project := range projects {
			if project.ID == id {
				return project.DisplayName == want, nil
			}
		}
		return false, nil
	})
	if err != nil {
		t.Fatalf("wait for Project %q display name %q: %v", id, want, err)
	}
}
