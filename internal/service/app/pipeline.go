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
}

func (a *PipelineApp) onServe(se *core.ServeEvent) error {
	g := se.Router.Group("/api/hdc/pipeline")

	g.GET("/{cluster}/{namespace}/cron-job/{name}", func(e *core.RequestEvent) error {
		ctx := e.Request.Context()
		cluster := e.Request.PathValue("cluster")
		namespace := e.Request.PathValue("namespace")
		name := e.Request.PathValue("name")
		cj, err := a.svc.GetCronJob(ctx, cluster, namespace, name)
		if err != nil {
			return err
		}
		return e.JSON(http.StatusOK, cj)
	})

	g.GET("/{cluster}/{namespace}/cron-job/{name}/jobs", func(e *core.RequestEvent) error {
		ctx := e.Request.Context()
		cluster := e.Request.PathValue("cluster")
		namespace := e.Request.PathValue("namespace")
		name := e.Request.PathValue("name")
		cj, err := a.svc.ListJobsFromCronJob(ctx, cluster, namespace, name)
		if err != nil {
			return err
		}
		return e.JSON(http.StatusOK, cj)
	})

	g.POST("/{cluster}/{namespace}/cron-job/{name}/job", func(e *core.RequestEvent) error {
		ctx := e.Request.Context()
		cluster := e.Request.PathValue("cluster")
		namespace := e.Request.PathValue("namespace")
		name := e.Request.PathValue("name")
		createdBy := e.Auth.FieldsData()["name"].(string)
		j, err := a.svc.CreateJobFromCronJob(ctx, cluster, namespace, name, createdBy)
		if err != nil {
			return err
		}
		return e.JSON(http.StatusCreated, j)
	}).Bind(apis.RequireAuth())

	return se.Next()
}

func (a *PipelineApp) onRecordCreateRequest(e *core.RecordRequestEvent) error {
	ctx := e.Request.Context()
	cluster := e.Record.GetString("cluster")
	namespace := e.Record.GetString("namespace")
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

	if _, err := a.svc.CreateCronJob(ctx, cluster, namespace, name, image, schedule); err != nil {
		return e.InternalServerError("Failed to create cron job", err)
	}

	return e.Next()
}

func (a *PipelineApp) onRecordUpdateRequest(e *core.RecordRequestEvent) error {
	ctx := e.Request.Context()
	cluster := e.Record.GetString("cluster")
	namespace := e.Record.GetString("namespace")
	name := e.Record.GetString("name")
	image := e.Record.GetString("image")
	schedule := e.Record.GetString("schedule")
	deleted := e.Record.GetBool("deleted")

	if strings.ToLower(name) != name {
		return e.BadRequestError("Name must be lowercase", nil)
	}
	if image == "" {
		return e.BadRequestError("Image must be provided", nil)
	}
	if schedule == "" {
		return e.BadRequestError("Schedule must be provided", nil)
	}

	if deleted {
		if err := a.svc.DeleteCronJob(ctx, cluster, namespace, name); err != nil {
			return e.InternalServerError("Failed to delete cron job", err)
		}

		return e.Next()
	}

	if _, err := a.svc.UpdateCronJob(ctx, cluster, namespace, name, image, schedule); err != nil {
		return e.InternalServerError("Failed to update cron job", err)
	}

	return e.Next()
}
