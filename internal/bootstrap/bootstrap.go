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
	"strings"

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

// harborURL, when non-empty, is the Harbor registry host used to
// register an additional OCI HelmRepository (oci://<host>/modules) on
// top of the embedded platform manifests.
//
// The method is idempotent and safe to call on every agent restart.
func (b *Bootstrapper) Run(ctx context.Context, harborURL string) error {
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

	// Harbor: register the OCI modules HelmRepository when a Harbor
	// registry host is configured. The Flux HelmRepository CRD is
	// already established by the base stage, so its GVR resolves here.
	if harborURL != "" {
		if err := b.applyManifest(ctx, harborModulesRepository(harborURL)); err != nil {
			return fmt.Errorf("harbor modules helm repository: %w", err)
		}
		b.log.Info("applied harbor modules HelmRepository", "host", harborHost(harborURL))
	}

	b.log.Info("layer 0 bootstrap completed successfully")
	return nil
}

// harborHost normalizes a configured Harbor value into a bare registry
// host suitable for an oci:// URL. It strips any http/https scheme and
// trailing slashes; the configured value may or may not carry a scheme.
func harborHost(harborURL string) string {
	host := harborURL
	host = strings.TrimPrefix(host, "https://")
	host = strings.TrimPrefix(host, "http://")
	return strings.TrimRight(host, "/")
}

// harborInsecure reports whether the OCI registry should be contacted
// over plaintext HTTP. A value with an explicit "http://" scheme is
// treated as insecure; "https://" or a bare host defaults to secure.
func harborInsecure(harborURL string) bool {
	return strings.HasPrefix(harborURL, "http://")
}

// harborModulesRepository renders a Flux OCI HelmRepository manifest
// pointing at oci://<host>/modules. The host is normalized via
// harborHost. When the configured value uses an "http://" scheme, the
// repository is marked insecure so Flux allows plaintext connections.
// The value originates from trusted operator configuration.
func harborModulesRepository(harborURL string) []byte {
	insecure := ""
	if harborInsecure(harborURL) {
		insecure = "\n  insecure: true"
	}
	return fmt.Appendf(nil, `apiVersion: source.toolkit.fluxcd.io/v1
kind: HelmRepository
metadata:
  name: oci-modules
  namespace: otterscale-system
  labels:
    tenant.otterscale.io/from-harbor: "true"
spec:
  type: oci
  interval: 6h
  provider: generic
  url: oci://%s/modules%s
`, harborHost(harborURL), insecure)
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
