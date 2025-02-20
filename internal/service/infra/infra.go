package infra

import (
	"github.com/google/wire"

	"github.com/openhdc/openhdc/internal/service/infra/repo"
)

var ProviderSet = wire.NewSet(repo.NewEntClient, repo.NewUserRepo)
