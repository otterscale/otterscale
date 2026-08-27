package core

import (
	"context"
	"errors"
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// watchRecorder implements ResourceRepo, capturing the options the use
// case passes to Watch. Only Watch is exercised here; the remaining
// methods satisfy the interface.
type watchRecorder struct {
	gotOpts WatchOptions
	watcher Watcher
	err     error
}

func (r *watchRecorder) Watch(_ context.Context, _ string, _ schema.GroupVersionResource, _ string, opts WatchOptions) (Watcher, error) {
	r.gotOpts = opts
	return r.watcher, r.err
}

func (r *watchRecorder) List(context.Context, string, schema.GroupVersionResource, string, ListOptions) (*unstructured.UnstructuredList, error) {
	return nil, nil
}

func (r *watchRecorder) Get(context.Context, string, schema.GroupVersionResource, string, string) (*unstructured.Unstructured, error) {
	return nil, nil
}

func (r *watchRecorder) Create(context.Context, string, schema.GroupVersionResource, string, []byte) (*unstructured.Unstructured, error) {
	return nil, nil
}

func (r *watchRecorder) Apply(context.Context, string, schema.GroupVersionResource, string, string, []byte, ApplyOptions) (*unstructured.Unstructured, error) {
	return nil, nil
}

func (r *watchRecorder) Update(context.Context, string, schema.GroupVersionResource, string, string, []byte, UpdateOptions) (*unstructured.Unstructured, error) {
	return nil, nil
}

func (r *watchRecorder) Delete(context.Context, string, schema.GroupVersionResource, string, string, DeleteOptions) error {
	return nil
}

func (r *watchRecorder) ListEvents(context.Context, string, string, ListOptions) (*unstructured.UnstructuredList, error) {
	return nil, nil
}

// noopWatcher is a Watcher that never produces an event.
type noopWatcher struct {
	ch chan WatchEvent
}

func newNoopWatcher() *noopWatcher { return &noopWatcher{ch: make(chan WatchEvent)} }

func (w *noopWatcher) ResultChan() <-chan WatchEvent { return w.ch }
func (w *noopWatcher) Stop()                         { close(w.ch) }

func watchTestIdentifier() *ResourceIdentifier {
	return &ResourceIdentifier{
		Cluster:   "prod",
		Group:     "",
		Version:   "v1",
		Resource:  "pods",
		Namespace: "default",
	}
}

// TestWatchResource_SendInitialEvents pins down when the use case asks
// the API server to stream the current state. The resumed-watch cases
// are the important ones: sendInitialEvents combined with a real
// resource version is rejected by the API server.
func TestWatchResource_SendInitialEvents(t *testing.T) {
	tests := []struct {
		name            string
		resourceVersion string
		watchList       bool
		watchListErr    error
		want            bool
	}{
		{
			name:      "fresh watch on a supporting cluster",
			watchList: true,
			want:      true,
		},
		{
			name:      "fresh watch on an older cluster",
			watchList: false,
			want:      false,
		},
		{
			name:            "explicit resource version 0 is still a fresh watch",
			resourceVersion: "0",
			watchList:       true,
			want:            true,
		},
		{
			name:            "resumed watch never asks for initial events",
			resourceVersion: "123456",
			watchList:       true,
			want:            false,
		},
		{
			name:         "unknown capability degrades to a plain watch",
			watchListErr: errors.New("version endpoint unavailable"),
			want:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			disco := &mockDiscovery{watchList: tt.watchList, watchListErr: tt.watchListErr}
			repo := &watchRecorder{watcher: newNoopWatcher()}
			uc := NewResourceUseCase(disco, repo, nil)

			watcher, err := uc.WatchResource(t.Context(), watchTestIdentifier(), WatchOptions{
				ResourceVersion: tt.resourceVersion,
			})
			if err != nil {
				t.Fatalf("WatchResource: %v", err)
			}
			defer watcher.Stop()

			if repo.gotOpts.SendInitialEvents != tt.want {
				t.Errorf("SendInitialEvents = %v, want %v", repo.gotOpts.SendInitialEvents, tt.want)
			}
			if repo.gotOpts.ResourceVersion != tt.resourceVersion {
				t.Errorf("ResourceVersion = %q, want %q", repo.gotOpts.ResourceVersion, tt.resourceVersion)
			}
		})
	}
}

// TestWatchResource_LookupError checks that an unknown GVR still fails
// the call rather than opening a watch.
func TestWatchResource_LookupError(t *testing.T) {
	disco := &mockDiscovery{lookupErr: errors.New("unable to recognize resource")}
	repo := &watchRecorder{watcher: newNoopWatcher()}
	uc := NewResourceUseCase(disco, repo, nil)

	if _, err := uc.WatchResource(t.Context(), watchTestIdentifier(), WatchOptions{}); err == nil {
		t.Fatal("expected an error for an unknown resource, got nil")
	}
}
