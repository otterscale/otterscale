// Package bootstrap provides the Layer 0 bootstrap process for the
// otterscale agent.
//
// All operations are idempotent: re-running bootstrap on a cluster
// that already has the resources installed is a safe no-op (or a
// controlled version bump).
package bootstrap

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"sort"

	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"

	"github.com/otterscale/otterscale/manifests"
)

// fieldManager is the SSA field manager identifier used for all
// bootstrap-applied resources. This allows kubectl and other tools to
// see which fields are owned by the agent's bootstrap process.
const fieldManager = "otterscale-agent"

// certManagerWebhookNamespace is the namespace where cert-manager
// installs its webhook Deployment.
const certManagerWebhookNamespace = "cert-manager"

// certManagerWebhookName is the Deployment name we wait for before
// proceeding to stage 2.
const certManagerWebhookName = "cert-manager-webhook"

// Bootstrapper applies embedded infrastructure manifests to the local
// Kubernetes cluster. It is injected into the Agent via Wire and
// called during agent startup.
type Bootstrapper struct {
	dynamic dynamic.Interface
	disc    discovery.DiscoveryInterface
	log     *slog.Logger
}

// New creates a Bootstrapper from the given rest.Config. The config
// is typically an in-cluster config provided by Wire. New creates the
// dynamic and discovery clients internally — only the config is
// injected, keeping the Wire graph minimal.
func New(cfg *rest.Config) (*Bootstrapper, error) {
	dyn, err := dynamic.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("create dynamic client: %w", err)
	}

	disc, err := discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("create discovery client: %w", err)
	}

	return &Bootstrapper{
		dynamic: dyn,
		disc:    disc,
		log:     slog.Default().With("component", "bootstrap"),
	}, nil
}

// The method is idempotent and safe to call on every agent restart.
func (b *Bootstrapper) Run(ctx context.Context) error {
	b.log.Info("starting Layer 0 bootstrap")

	// Base: cert-manager + CRDs + FluxCD
	if err := b.applyStage(ctx, manifests.Base, "bootstrap/base"); err != nil {
		return fmt.Errorf("base: %w", err)
	}

	// Wait for cert-manager-webhook Deployment to be Available.
	if err := b.waitForDeployment(ctx, certManagerWebhookNamespace, certManagerWebhookName); err != nil {
		return fmt.Errorf("wait for cert-manager webhook: %w", err)
	}
	b.log.Info("cert-manager webhook is available")

	// Platform: tenant-operator + HelmRepository
	if err := b.applyStage(ctx, manifests.Platform, "bootstrap/platform"); err != nil {
		return fmt.Errorf("platform: %w", err)
	}

	b.log.Info("layer 0 bootstrap completed successfully")
	return nil
}

// applyStage reads every embedded YAML manifest from the given
// embed.FS directory and applies it to the cluster. Files are
// processed in lexicographic order so that ordering can be controlled
// via file-name prefixes if needed.
func (b *Bootstrapper) applyStage(ctx context.Context, fsys embed.FS, dir string) error {
	entries, err := fsys.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read embedded manifests directory %s: %w", dir, err)
	}

	// Sort entries explicitly (embed.FS returns sorted results per
	// the spec, but being explicit costs nothing).
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		data, err := fs.ReadFile(fsys, dir+"/"+name)
		if err != nil {
			return fmt.Errorf("read manifest %s: %w", name, err)
		}

		b.log.Info("applying manifest", "file", name)
		if err := b.applyManifest(ctx, data); err != nil {
			return fmt.Errorf("apply manifest %s: %w", name, err)
		}
	}

	return nil
}
