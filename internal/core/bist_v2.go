package core

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"

	"github.com/openhdc/otterscale/internal/config"
)

const (
	bistKindFIO             = "fio"
	bistKindWarp            = "warp"
	bistNamespace           = "bist"
	bistLabel               = "bist.otterscale.io/name=bist"
	bistAnnotationCreatedBy = "bist.otterscale.io/created-by"
	bistAnnotationKind      = "bist.otterscale.io/kind"
	bistAnnotationFIO       = "bist.otterscale.io/fio"
	bistAnnotationWarp      = "bist.otterscale.io/warp"
)

const (
	minioLabel       = "app.kubernetes.io/name=minio"
	minioField       = "spec.type=NodePort"
	minioServiceName = "minio-api"
)

type BISTResult struct {
	UID            string
	Name           string
	Status         string
	CreatedBy      string
	StartTime      time.Time
	CompletionTime time.Time
	FIO            *FIO
	Warp           *Warp
}

type FIO struct {
	Target FIOTarget  `json:"target"`
	Input  *FIOInput  `json:"input,omitempty"`
	Output *FIOOutput `json:"output,omitempty"`
}

type FIOTarget struct {
	Ceph *FIOTargetCeph `json:"ceph,omitempty"`
	NFS  *FIOTargetNFS  `json:"nfs,omitempty"`
}

type FIOTargetCeph struct {
	ScopeUUID    string `json:"scope_uuid"`
	FacilityName string `json:"facility_name"`
}

type FIOTargetNFS struct {
	Endpoint string `json:"endpoint"`
	Path     string `json:"path"`
}

type FIOInput struct {
	AccessMode string `json:"access_mode"`
	JobCount   int64  `json:"job_count"`
	RunTime    string `json:"run_time"`
	BlockSize  string `json:"block_size"`
	FileSize   string `json:"file_size"`
	IODepth    int64  `json:"io_depth"`
}

type FIOOutput struct {
	Read  *FIOThroughput `json:"read"`
	Write *FIOThroughput `json:"write"`
	Trim  *FIOThroughput `json:"trim"`
}

type FIOThroughput struct {
	IOBytes   int64   `json:"io_bytes"`
	Bandwidth int64   `json:"bw"`
	IOPS      float64 `json:"iops"`
	TotalIOs  int64   `json:"total_ios"`
	Latency   struct {
		Min  int64   `json:"min"`
		Max  int64   `json:"max"`
		Mean float64 `json:"mean"`
	} `json:"lat_ns"`
}

type Warp struct {
	Target WarpTarget  `json:"target"`
	Input  *WarpInput  `json:"input,omitempty"`
	Output *WarpOutput `json:"output,omitempty"`
}

type WarpTarget struct {
	Internal *WarpTargetInternal `json:"internal,omitempty"`
	External *WarpTargetExternal `json:"external,omitempty"`
}

type WarpTargetInternal struct {
	Type     string `json:"type"`
	Name     string `json:"name"`
	Endpoint string `json:"endpoint"`
}

type WarpTargetExternal struct {
	Endpoint  string `json:"endpoint"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
}

type WarpInput struct {
	Operation  string `json:"operation"`
	Duration   string `json:"duration"`
	ObjectSize string `json:"object_size"`
	ObjectNum  string `json:"object_num"`
}

type WarpOutput struct {
	Type       string          `json:"type"`
	Operations []WarpOperation `json:"operations"`
}

type WarpOperation struct {
	Type       string `json:"type"`
	Throughput struct {
		Metrics struct {
			FastestBPS float64 `json:"fastest_bps"`
			MedianBPS  float64 `json:"median_bps"`
			SlowestBPS float64 `json:"slowest_bps"`
			FastestOPS float64 `json:"fastest_ops"`
			MedianOPS  float64 `json:"median_ops"`
			SlowestOPS float64 `json:"slowest_ops"`
		} `json:"segmented"`
		TotalBytes      float64 `json:"bytes"`
		TotalObjects    float64 `json:"objects"`
		TotalOperations int64   `json:"ops"`
	} `json:"throughput"`
}

type BISTUseCase struct {
	scope     ScopeRepo
	client    ClientRepo
	facility  FacilityRepo
	kubeBatch KubeBatchRepo
	kubeCore  KubeCoreRepo

	conf *config.Config
}

func NewBISTUseCase(scope ScopeRepo, client ClientRepo, facility FacilityRepo, kubeBatch KubeBatchRepo, kubeCore KubeCoreRepo, conf *config.Config) *BISTUseCase {
	return &BISTUseCase{
		scope:     scope,
		client:    client,
		facility:  facility,
		kubeBatch: kubeBatch,
		kubeCore:  kubeCore,
		conf:      conf,
	}
}

func (uc *BISTUseCase) ListResults(ctx context.Context) ([]BISTResult, error) {
	config, err := uc.newMicroK8sConfig()
	if err != nil {
		return nil, err
	}
	jobs, err := uc.kubeBatch.ListJobsByLabel(ctx, config, bistNamespace, bistLabel)
	if err != nil {
		return nil, err
	}
	results := []BISTResult{}
	for i := range jobs {
		bist, err := uc.toBISTResult(ctx, config, &jobs[i])
		if err != nil {
			return nil, err
		}
		results = append(results, *bist)
	}
	return results, nil
}

func (uc *BISTUseCase) CreateFIOResult(ctx context.Context, name, createdBy string, input *FIOInput, target *FIOTarget) (*BISTResult, error) {
	return nil, nil
}

func (uc *BISTUseCase) CreateWarpResult(ctx context.Context, name, createdBy string, input *WarpInput, target *WarpTarget) (*BISTResult, error) {
	return nil, nil
}

func (uc *BISTUseCase) DeleteResult(ctx context.Context, name string) error {
	config, err := uc.newMicroK8sConfig()
	if err != nil {
		return err
	}
	return uc.kubeBatch.DeleteJob(ctx, config, bistNamespace, name)
}

func (uc *BISTUseCase) ListInternalObjectServices(ctx context.Context, uuid string) ([]WarpTargetInternal, error) {
	var cephs, minios []WarpTargetInternal
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		svcs, err := uc.listCephObjectServices(ctx, uuid)
		if err != nil {
			return err
		}
		cephs = svcs
		return nil
	})
	eg.Go(func() error {
		svcs, err := uc.listMinIOs(ctx, uuid)
		if err != nil {
			return err
		}
		minios = svcs
		return nil
	})
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return append(cephs, minios...), nil
}

func (uc *BISTUseCase) listCephObjectServices(ctx context.Context, uuid string) ([]WarpTargetInternal, error) {
	cephs, err := listCephs(ctx, uc.scope, uc.client, uuid)
	if err != nil {
		return nil, err
	}
	services := []WarpTargetInternal{}
	for _, ceph := range cephs {
		leader, err := uc.facility.GetLeader(ctx, uuid, rgwName(ceph.Name))
		if err != nil {
			continue
		}
		info, err := uc.facility.GetUnitInfo(ctx, uuid, leader)
		if err != nil {
			continue
		}
		services = append(services, WarpTargetInternal{
			Type:     "ceph",
			Name:     ceph.Name,
			Endpoint: info.PublicAddress,
		})
	}
	return services, nil
}

func (uc *BISTUseCase) listMinIOs(ctx context.Context, uuid string) ([]WarpTargetInternal, error) {
	kubes, err := listKuberneteses(ctx, uc.scope, uc.client, uuid)
	if err != nil {
		return nil, err
	}
	services := []WarpTargetInternal{}
	for _, kube := range kubes {
		config, err := kubeConfig(ctx, uc.facility, uuid, kube.Name)
		if err != nil {
			continue
		}
		svcs, err := uc.kubeCore.ListServicesByOptions(ctx, config, "", minioLabel, minioField)
		if err != nil {
			continue
		}
		for i := range svcs {
			for _, port := range svcs[i].Spec.Ports {
				if port.Name != "minio-api" {
					continue
				}
				url, err := url.Parse(config.Host)
				if err != nil {
					continue
				}
				services = append(services, WarpTargetInternal{
					Type:     "minio",
					Name:     fmt.Sprintf("%s.%s", svcs[i].GetNamespace(), svcs[i].GetName()),
					Endpoint: fmt.Sprintf("%s:%d", url.Hostname(), port.NodePort),
				})
			}
		}
	}
	return services, nil
}

func (uc *BISTUseCase) newMicroK8sConfig() (*rest.Config, error) {
	token, err := base64.StdEncoding.DecodeString(uc.conf.MicroK8s.Token)
	if err != nil {
		return nil, err
	}
	return &rest.Config{
		Host:        uc.conf.MicroK8s.Host,
		BearerToken: string(token),
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true,
		},
	}, nil
}

func (uc *BISTUseCase) toBISTResult(ctx context.Context, config *rest.Config, job *Job) (*BISTResult, error) {
	annotations := job.GetAnnotations()
	kind, ok := annotations[bistAnnotationKind]
	if !ok {
		return nil, errors.New("kind of bist not found")
	}

	result := &BISTResult{
		UID:            string(job.UID),
		Name:           job.Name,
		Status:         uc.toBISTResultStatus(job),
		CreatedBy:      annotations[bistAnnotationCreatedBy],
		StartTime:      job.Status.StartTime.Time,
		CompletionTime: job.Status.CompletionTime.Time,
	}

	switch kind {
	case bistKindFIO:
		result.FIO, _ = uc.toFIO(ctx, config, job)
	case bistKindWarp:
		result.Warp, _ = uc.toWarp(ctx, config, job)
	}

	return result, nil
}

func (uc *BISTUseCase) getLogs(ctx context.Context, config *rest.Config, job *Job) (map[string]any, error) {
	if job.Status.CompletionTime == nil {
		return map[string]any{}, nil
	}
	selector, err := metav1.LabelSelectorAsSelector(job.Spec.Selector)
	if err != nil {
		return nil, err
	}

	pods, err := uc.kubeCore.ListPodsByLabel(ctx, config, bistNamespace, selector.String())
	if err != nil {
		return nil, err
	}

	for i := range pods {
		logs, err := uc.kubeCore.GetPodLogs(ctx, config, bistNamespace, pods[i].Name, pods[i].Spec.Containers[0].Name)
		if err != nil {
			continue
		}
		// warp result has redundant message
		for _, v := range []string{logs, uc.removeLastTwoLines(logs)} {
			var result map[string]any
			if err := json.Unmarshal([]byte(v), &result); err != nil {
				continue
			}
			return result, nil
		}
	}
	return map[string]any{}, nil
}

func (uc *BISTUseCase) toBISTResultStatus(job *Job) string {
	if job.Status.Succeeded > 0 {
		return "succeeded"
	}
	if job.Status.Failed > 0 {
		return "failed"
	}
	return "running"
}

func (uc *BISTUseCase) removeLastTwoLines(input string) string {
	lines := strings.Split(input, "\n")
	if len(lines) < 2 {
		return input
	}
	lines = lines[:len(lines)-2]
	return strings.Join(lines, "\n")
}

func (uc *BISTUseCase) toFIO(ctx context.Context, config *rest.Config, job *Job) (*FIO, error) {
	fio := &FIO{}
	// target & input
	if err := json.Unmarshal([]byte(job.Annotations[bistAnnotationFIO]), fio); err != nil {
		return nil, err
	}
	// output
	logs, err := uc.getLogs(ctx, config, job)
	if err != nil {
		return nil, err
	}
	v, ok := logs["jobs"].([]any)
	if ok {
		b, err := json.Marshal(v[0])
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(b, &fio.Output); err != nil {
			return nil, err
		}
	}
	return fio, nil
}

func (uc *BISTUseCase) toWarp(ctx context.Context, config *rest.Config, job *Job) (*Warp, error) {
	warp := &Warp{}
	// target & input
	if err := json.Unmarshal([]byte(job.Annotations[bistAnnotationWarp]), warp); err != nil {
		return nil, err
	}
	// output
	// logs, err := uc.getLogs(ctx, config, job)
	// if err != nil {
	// 	return nil, err
	// }
	// if err := json.Unmarshal([]byte(logs), &warp.Output); err != nil {
	// 	return nil, err
	// }

	// TODO

	return warp, nil
}
