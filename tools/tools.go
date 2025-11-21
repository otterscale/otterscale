//go:build tools

// This package imports things required by wire, to force `go mod` to see them as dependencies
package tools

import (
	_ "github.com/google/subcommands"
	_ "golang.org/x/tools/cmd/stringer"
)
