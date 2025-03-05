package app

import "github.com/google/wire"

type empty struct{}

var ProviderSet = wire.NewSet(NewPipelineApp)
