package service

import (
	"github.com/google/wire"
	"github.com/pocketbase/pocketbase/core"

	"github.com/openhdc/openhdc/internal/service/app"
)

var ProviderSet = wire.NewSet(NewPocketBase)

func NewPocketBase(ua *app.UserApp) []func(se *core.ServeEvent) error {
	pb := []func(se *core.ServeEvent) error{
		ua.Bind(),
	}
	return pb
}

// func bindRecordCrudApi(app core.App, rg *router.RouterGroup[*core.RequestEvent]) {
