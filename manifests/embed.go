// Package manifests embeds static Kubernetes YAML files that the
// agent applies during its Layer 0 bootstrap phase.
package manifests

import "embed"

//go:embed bootstrap/stage1/*.yaml
var Stage1 embed.FS

//go:embed bootstrap/stage2/*.yaml
var Stage2 embed.FS
