package app

import (
	"fmt"
	"net/http"

	"github.com/pocketbase/pocketbase/core"

	"github.com/openhdc/openhdc/internal/service/domain/service"
)

type UserApp struct {
	svc *service.UserService
}

func NewUserApp(svc *service.UserService) *UserApp {
	return &UserApp{
		svc: svc,
	}
}

func (a *UserApp) BindOnRecordAfterCreateSuccess(app core.App, fn func(e *core.RecordEvent) error) {
	app.OnRecordAfterCreateSuccess("kubernetes_jobs").BindFunc(fn)
}

func (a *UserApp) Bind() func(se *core.ServeEvent) error {
	return func(se *core.ServeEvent) error {
		subGroup := se.Router.Group("/v1/api/users").Bind() // apis.RequireSuperuserAuth()
		subGroup.GET("/{id}", collectionView)
		// _, err := a.svc.Get(ctx, req.GetId())
		// if err != nil {
		// 	return nil, err
		// }
		return se.Next()
	}
}

func collectionView(e *core.RequestEvent) error {
	fmt.Println(e.Request.PathValue("id"))
	collection, err := e.App.FindCachedCollectionByNameOrId(e.Request.PathValue("id"))
	if err != nil || collection == nil {
		return e.NotFoundError("", err)
	}

	event := new(core.CollectionRequestEvent)
	event.RequestEvent = e
	event.Collection = collection

	return e.App.OnCollectionViewRequest().Trigger(event, func(e *core.CollectionRequestEvent) error {
		return e.JSON(http.StatusOK, e.Collection)
	})
}
