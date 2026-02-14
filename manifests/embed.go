// Package manifests embeds static Kubernetes YAML files that the
// agent applies during its Layer 0 bootstrap phase. Keeping the
// manifests in a top-level directory (rather than internal/) makes
// them easy to inspect and update without diving into Go packages.
package manifests

import "embed"

// Bootstrap holds the YAML manifests applied during Layer 0 bootstrap
// (FluxCD core components, otterscale-operator, Module CRD, etc.).
// Files are accessed via the "bootstrap/" prefix.
//
//go:embed bootstrap/*.yaml
var Bootstrap embed.FS
