package app

import (
	"net/http"

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
	app.OnRecordAfterCreateSuccess("pipelines").BindFunc(a.onRecordAfterCreateSuccess)
}

// TODO ADD AUTH
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
		data := struct {
			CreatedBy string `json:"createdBy" form:"createdBy"`
		}{}
		if err := e.BindBody(&data); err != nil {
			return e.BadRequestError("Failed to read request data", err)
		}
		if err := a.svc.CreateJobFromCronJob(e.Request.Context(), e.Request.PathValue("name"), data.CreatedBy); err != nil {
			return err
		}
		return e.JSON(http.StatusCreated, empty{})
	})
	//.Bind(apis.RequireAuth())

	return se.Next()
}

func (a *PipelineApp) onRecordAfterCreateSuccess(e *core.RecordEvent) error {
	return a.svc.CreateCronJob(e.Context, "hello", "busybox:1.28", "* * * * *")
}
