package core

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/url"
	"time"

	"golang.org/x/sync/errgroup"
	"k8s.io/client-go/rest"

	"github.com/openhdc/otterscale/internal/config"
)

const (
	bistNamespace = "bist"
)

const (
	minioLabel       = "app.kubernetes.io/name=minio"
	minioField       = "spec.type=NodePort"
	minioServiceName = "minio-api"
)

type BISTResult struct {
	UID          string
	Name         string
	Status       string
	CreatedBy    string
	StartTime    time.Time
	CompleteTime time.Time
	FIO          *FIO
	Warp         *Warp
}

type FIO struct {
	Target FIOTarget
	Input  *FIOInput
	Output *FIOOutput
}

type FIOTarget struct {
	Ceph *FIOTargetCeph
	NFS  *FIOTargetNFS
}

type FIOTargetCeph struct {
	ScopeUUID    string
	FacilityName string
}

type FIOTargetNFS struct {
	Endpoint string
	Path     string
}

type FIOInput struct {
	AccessMode string
	JobCount   int64
	RunTime    string
	BlockSize  string
	FileSize   string
	IODepth    int64
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
	Target WarpTarget
	Input  *WarpInput
	Output *WarpOutput
}

type WarpTarget struct {
	Internal *WarpTargetInternal
	External *WarpTargetExternal
}

type WarpTargetInternal struct {
	Type     string
	Name     string
	Endpoint string
}

type WarpTargetExternal struct {
	Endpoint  string
	AccessKey string
	SecretKey string
}

type WarpInput struct {
	Operation  string
	Duration   string
	ObjectSize string
	ObjectNum  string
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
	return nil, nil
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
