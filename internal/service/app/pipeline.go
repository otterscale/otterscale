package app

import (
	"net/http"
	"strings"

	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"

	"github.com/openhdc/openhdc/internal/service/domain/service"
)

type PipelineApp struct {
	svc *service.KubeService
}

func NewPipelineApp(svc *service.KubeService) *PipelineApp {
	return &PipelineApp{
		svc: svc,
	}
}

func (a *PipelineApp) Bind(app core.App) {
	app.OnServe().BindFunc(a.onServe)
	app.OnRecordCreateRequest("pipelines").BindFunc(a.onRecordCreateRequest)
	app.OnRecordUpdateRequest("pipelines").BindFunc(a.onRecordUpdateRequest)
	app.OnRecordDeleteRequest("pipelines").BindFunc(a.onRecordDeleteRequest)
}

func (a *PipelineApp) onServe(se *core.ServeEvent) error {
	g := se.Router.Group("/api/hdc/pipeline")

	g.GET("/cron-job/{name}", func(e *core.RequestEvent) error {
		cj, err := a.svc.GetCronJob(e.Request.Context(), e.Request.PathValue("name"))
		if err != nil {
			return err
		}
		return e.JSON(http.StatusOK, cj)
	})

	g.GET("/cron-job/{name}/jobs", func(e *core.RequestEvent) error {
		cj, err := a.svc.ListJobsFromCronJob(e.Request.Context(), e.Request.PathValue("name"))
		if err != nil {
			return err
		}
		return e.JSON(http.StatusOK, cj)
	})

	g.POST("/cron-job/{name}/job", func(e *core.RequestEvent) error {
		j, err := a.svc.CreateJobFromCronJob(e.Request.Context(), e.Request.PathValue("name"), e.Auth.FieldsData()["name"].(string))
		if err != nil {
			return err
		}
		return e.JSON(http.StatusCreated, j)
	}).Bind(apis.RequireAuth())

	return se.Next()
}

func (a *PipelineApp) onRecordCreateRequest(e *core.RecordRequestEvent) error {
	ctx := e.Request.Context()
	name := e.Record.GetString("name")
	image := e.Record.GetString("image")
	schedule := e.Record.GetString("schedule")

	if strings.ToLower(name) != name {
		return e.BadRequestError("Name must be lowercase", nil)
	}
	if image == "" {
		return e.BadRequestError("Image must be provided", nil)
	}
	if schedule == "" {
		return e.BadRequestError("Schedule must be provided", nil)
	}

	if _, err := a.svc.CreateCronJob(ctx, name, image, schedule); err != nil {
		return e.InternalServerError("Failed to create cron job", err)
	}

	return e.Next()
}

func (a *PipelineApp) onRecordUpdateRequest(e *core.RecordRequestEvent) error {
	ctx := e.Request.Context()
	name := e.Record.GetString("name")
	image := e.Record.GetString("image")
	schedule := e.Record.GetString("schedule")

	if strings.ToLower(name) != name {
		return e.BadRequestError("Name must be lowercase", nil)
	}
	if image == "" {
		return e.BadRequestError("Image must be provided", nil)
	}
	if schedule == "" {
		return e.BadRequestError("Schedule must be provided", nil)
	}

	if _, err := a.svc.UpdateCronJob(ctx, name, image, schedule); err != nil {
		return e.InternalServerError("Failed to update cron job", err)
	}

	return e.Next()
}

func (a *PipelineApp) onRecordDeleteRequest(e *core.RecordRequestEvent) error {
	ctx := e.Request.Context()
	name := e.Record.GetString("name")

	if strings.ToLower(name) != name {
		return e.BadRequestError("Name must be lowercase", nil)
	}
	if err := a.svc.DeleteCronJob(ctx, name); err != nil {
		return e.InternalServerError("Failed to delete cron job", err)
	}

	return e.Next()
}
