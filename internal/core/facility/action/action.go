package action

import (
	"context"

	"github.com/juju/juju/api/client/action"
)

type Action struct {
	Name string
	Spec *action.ActionSpec
}

type ActionRepo interface {
	List(ctx context.Context, scope, appName string) ([]Action, error)
	Run(ctx context.Context, scope, appName, actionName string, params map[string]any) (map[string]any, error)
	Execute(ctx context.Context, scope, appName, command string) (map[string]any, error)
}

type ActionUseCase struct {
	action ActionRepo
}

func (uc *ActionUseCase) ListActions(ctx context.Context, scope, appName string) ([]Action, error) {
	return uc.action.List(ctx, scope, appName)
}
