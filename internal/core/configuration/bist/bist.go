package bist

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	conf "github.com/otterscale/otterscale/internal/config"
	"github.com/otterscale/otterscale/internal/core/application/cluster"
	"github.com/otterscale/otterscale/internal/core/application/config"
	"github.com/otterscale/otterscale/internal/core/application/release"
	"github.com/otterscale/otterscale/internal/core/application/service"
	"github.com/otterscale/otterscale/internal/core/application/workload"
	"github.com/otterscale/otterscale/internal/core/scope"
	"github.com/otterscale/otterscale/internal/core/storage"
	"github.com/otterscale/otterscale/internal/core/storage/block"
	"github.com/otterscale/otterscale/internal/core/storage/object"
)

const bistNamespace = "bist"

const (
	kindFIO  = "fio"
	kindWarp = "warp"
)

const (
	kindLabel      = "otterscale.com/bist.kind"
	createdByLabel = "otterscale.com/bist.created-by"
)

const (
	fioAnnotation  = "otterscale.com/bist.fio"
	warpAnnotation = "otterscale.com/bist.warp"
)

type Result struct {
	UID            string
	Name           string
	Status         ResultStatus
	CreatedBy      string
	StartTime      time.Time
	CompletionTime time.Time
	FIO            *FIO
	Warp           *Warp
}

type UseCase struct {
	conf *conf.Config

	bucket      object.BucketRepo
	configMap   config.ConfigMapRepo
	image       block.ImageRepo
	job         workload.JobRepo
	namespace   cluster.NamespaceRepo
	clusterNode cluster.NodeRepo
	storageNode storage.NodeRepo
	pod         workload.PodRepo
	pool        storage.PoolRepo
	secret      config.SecretRepo
	service     service.ServiceRepo
}

func NewUseCase(conf *conf.Config, bucket object.BucketRepo, configMap config.ConfigMapRepo, image block.ImageRepo, job workload.JobRepo, namespace cluster.NamespaceRepo, clusterNode cluster.NodeRepo, storageNode storage.NodeRepo, pod workload.PodRepo, pool storage.PoolRepo, secret config.SecretRepo, service service.ServiceRepo) *UseCase {
	return &UseCase{
		conf:        conf,
		bucket:      bucket,
		configMap:   configMap,
		image:       image,
		job:         job,
		namespace:   namespace,
		clusterNode: clusterNode,
		storageNode: storageNode,
		pod:         pod,
		pool:        pool,
		secret:      secret,
		service:     service,
	}
}

func (uc *UseCase) ListResults(ctx context.Context) ([]Result, error) {
	selector := release.TypeLabel + "=" + "bist"

	jobs, err := uc.job.List(ctx, scope.ReservedName, bistNamespace, selector)
	if err != nil {
		return nil, err
	}

	return uc.toResults(ctx, jobs)
}

func (uc *UseCase) DeleteResult(ctx context.Context, name string) error {
	return uc.job.Delete(ctx, scope.ReservedName, bistNamespace, name)
}

func (uc *UseCase) ListInternalObjectServices(ctx context.Context, scope string) ([]WarpTargetInternal, error) {
	targets, err := uc.listMinIOs(ctx, scope)
	if err != nil {
		return nil, err
	}

	url, err := url.Parse(uc.bucket.Endpoint(scope))
	if err != nil {
		return nil, err
	}

	targets = append(targets, WarpTargetInternal{
		Type:  ObjectServiceTypeCeph,
		Scope: scope,
		Name:  scope,
		Host:  url.Host,
	})

	return targets, nil
}

func (uc *UseCase) listMinIOs(ctx context.Context, scope string) ([]WarpTargetInternal, error) {
	internalIP, err := uc.clusterNode.InternalIP(ctx, scope)
	if err != nil {
		return nil, err
	}

	selector := "app.kubernetes.io/name" + "=" + "minio"

	services, err := uc.service.List(ctx, scope, "", selector)
	if err != nil {
		return nil, err
	}

	targets := []WarpTargetInternal{}

	for i := range services {
		if services[i].Spec.Type != "NodePort" {
			continue
		}

		for _, port := range services[i].Spec.Ports {
			if port.Name != "minio-api" {
				continue
			}

			targets = append(targets, WarpTargetInternal{
				Type:  ObjectServiceTypeMinIO,
				Scope: scope,
				Name:  fmt.Sprintf("%s.%s", services[i].GetNamespace(), services[i].GetName()),
				Host:  fmt.Sprintf("%s:%d", internalIP, port.NodePort),
			})
		}
	}

	return targets, nil
}

func (uc *UseCase) toResultStatus(job *workload.Job) ResultStatus {
	if job.Status.Succeeded > 0 {
		return ResultStatusSucceeded
	}
	if job.Status.Failed > 0 {
		return ResultStatusFailed
	}
	return ResultStatusRunning
}

func (uc *UseCase) toResult(ctx context.Context, job *workload.Job) (*Result, error) {
	labels := job.GetLabels()

	kind, ok := labels[kindLabel]
	if !ok {
		return nil, errors.New("kind of bist not found")
	}

	result := &Result{
		UID:       string(job.UID),
		Name:      job.Name,
		Status:    uc.toResultStatus(job),
		CreatedBy: labels[createdByLabel],
	}

	if job.Status.StartTime != nil {
		result.StartTime = job.Status.StartTime.Time
	}

	if job.Status.CompletionTime != nil {
		result.CompletionTime = job.Status.CompletionTime.Time
	}

	switch kind {
	case kindFIO:
		result.FIO, _ = uc.toFIO(ctx, job)

	case kindWarp:
		result.Warp, _ = uc.toWarp(ctx, job)
	}

	return result, nil
}

func (uc *UseCase) toResults(ctx context.Context, jobs []workload.Job) ([]Result, error) {
	results := []Result{}

	for i := range jobs {
		bist, err := uc.toResult(ctx, &jobs[i])
		if err != nil {
			return nil, err
		}

		results = append(results, *bist)
	}

	return results, nil
}

func (uc *UseCase) fetchLogs(ctx context.Context, job *workload.Job) (map[string]any, error) {
	if job.Status.CompletionTime == nil {
		return map[string]any{}, nil
	}

	selector, err := v1.LabelSelectorAsSelector(job.Spec.Selector)
	if err != nil {
		return nil, err
	}

	pods, err := uc.pod.List(ctx, scope.ReservedName, bistNamespace, selector.String())
	if err != nil {
		return nil, err
	}

	for i := range pods {
		logs, err := uc.pod.Stream(ctx, scope.ReservedName, bistNamespace, pods[i].Name, pods[i].Spec.Containers[0].Name, 0, false)
		if err != nil {
			continue
		}

		data, err := io.ReadAll(logs)
		if err != nil {
			continue
		}

		logs.Close()

		// warp result has redundant message
		for _, v := range []string{string(data), removeLastTwoLines(string(data))} {
			var result map[string]any
			if err := json.Unmarshal([]byte(v), &result); err != nil {
				continue
			}

			return result, nil
		}
	}

	return map[string]any{}, nil
}

func removeLastTwoLines(input string) string {
	lines := strings.Split(input, "\n")
	if len(lines) < 2 {
		return input
	}
	lines = lines[:len(lines)-2]
	return strings.Join(lines, "\n")
}

func shortID(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:4])
}
