package providers

import (
	"errors"
	"testing"

	"k8s.io/client-go/rest"

	"github.com/otterscale/otterscale/internal/core"
)

func TestProvideRancherStoreDegradesToCacheNotReady(t *testing.T) {
	tests := []struct {
		name    string
		provide func() (*rest.Config, error)
	}{
		{
			name: "Kubernetes config unavailable",
			provide: func() (*rest.Config, error) {
				return nil, errors.New("config unavailable")
			},
		},
		{
			name: "dynamic client initialization fails",
			provide: func() (*rest.Config, error) {
				return &rest.Config{
					Host: "https://127.0.0.1",
					TLSClientConfig: rest.TLSClientConfig{
						CAData: []byte("invalid CA data"),
					},
				}, nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := provideRancherStore(tt.provide)
			if _, err := store.ListProjects(t.Context()); !errors.Is(err, core.ErrRancherProjectCacheNotReady) {
				t.Fatalf("ListProjects error = %v", err)
			}
		})
	}
}
