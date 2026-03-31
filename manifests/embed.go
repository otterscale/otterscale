// Package manifests embeds static Kubernetes YAML files that the
// agent applies during its Layer 0 bootstrap phase.
package manifests

import "embed"

//go:embed bootstrap/base/*.yaml
var Base embed.FS

//go:embed bootstrap/platform/*.yaml
var Platform embed.FS
