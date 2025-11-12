package storage

import (
	"github.com/google/wire"

	"github.com/otterscale/otterscale/internal/core/storage/block"
	"github.com/otterscale/otterscale/internal/core/storage/file"
	"github.com/otterscale/otterscale/internal/core/storage/node"
	"github.com/otterscale/otterscale/internal/core/storage/object"
	"github.com/otterscale/otterscale/internal/core/storage/pool"
)

var ProviderSet = wire.NewSet(
	block.NewBlockUseCase,
	file.NewFileUseCase,
	node.NewNodeUseCase,
	object.NewObjectUseCase,
	pool.NewPoolUseCase,
)
