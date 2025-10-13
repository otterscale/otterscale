package app

import (
	"github.com/otterscale/otterscale/api/model/v1/pbconnect"
	"github.com/otterscale/otterscale/internal/core"
)

type ModelService struct {
	pbconnect.UnimplementedModelServiceHandler

	model *core.ModelUseCase
}

func NewModelService(model *core.ModelUseCase) *ModelService {
	return &ModelService{model: model}
}

var _ pbconnect.ModelServiceHandler = (*ModelService)(nil)
