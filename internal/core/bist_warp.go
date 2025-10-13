package core

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
)

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
	Type         string `json:"type"`
	ScopeUUID    string `json:"scope_uuid"`
	FacilityName string `json:"facility_name"`
	Name         string `json:"name"`
	Endpoint     string `json:"endpoint"`
}

type WarpTargetExternal struct {
	Endpoint  string `json:"endpoint"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
}

type WarpInput struct {
	Operation   string `json:"operation"`
	Duration    int64  `json:"duration"`
	ObjectSize  int64  `json:"object_size"`
	ObjectCount int64  `json:"object_count"`
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

func (uc *BISTUseCase) CreateWarpResult(ctx context.Context, name, createdBy string, target WarpTarget, input *WarpInput) (*BISTResult, error) {
	config, err := uc.newMicroK8sConfig()
	if err != nil {
		return nil, err
	}

	// namespace
	if err := uc.ensureNamespace(ctx, config); err != nil {
		return nil, err
	}

	// name
	if name == "" {
		name = fmt.Sprintf("bist-%s", generateHashedName(strconv.FormatInt(time.Now().UnixNano(), 10)))
	}

	// labels
	labelName := strings.Split(bistLabel, "=")
	labels := map[string]string{
		labelName[0]: labelName[1],
	}

	// annotations
	warp, _ := json.Marshal(Warp{
		Target: target,
		Input:  input,
	})
	annotations := map[string]string{
		bistAnnotationCreatedBy: createdBy,
		bistAnnotationKind:      bistKindWarp,
		bistAnnotationWarp:      string(warp),
	}

	// spec
	var spec *JobSpec

	// s3 internal
	internal := target.Internal
	if internal != nil {
		switch internal.Type {
		case "ceph":
			spec, err = uc.warpCephObjectGatewayJobSpec(ctx, internal, input)
			if err != nil {
				return nil, err
			}
		case "minio":
			spec, err = uc.warpMinIOJobSpec(ctx, internal, input)
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("unsupported internal kind %q", internal.Type)
		}
	}

	// s3 external
	external := target.External
	if external != nil {
		spec = uc.warpJobSpec(external, input)
	}

	// job
	job, err := uc.kubeBatch.CreateJob(ctx, config, bistNamespace, name, labels, annotations, spec)
	if err != nil {
		return nil, err
	}
	return uc.toBISTResult(ctx, config, job)
}

func (uc *BISTUseCase) warpCephObjectGatewayJobSpec(ctx context.Context, target *WarpTargetInternal, input *WarpInput) (*JobSpec, error) {
	sc, err := storageConfig(ctx, uc.facility, uc.action, target.ScopeUUID, target.FacilityName)
	if err != nil {
		return nil, err
	}
	return uc.warpJobSpec(&WarpTargetExternal{
		Endpoint:  target.Endpoint, // without protocol prefix
		AccessKey: sc.AccessKey,
		SecretKey: sc.SecretKey,
	}, input), nil
}

func (uc *BISTUseCase) warpMinIOJobSpec(ctx context.Context, target *WarpTargetInternal, input *WarpInput) (*JobSpec, error) {
	tmp := strings.Split(target.Name, ".")
	if len(tmp) != 2 {
		return nil, fmt.Errorf("invalid name %q", target.Name)
	}
	kc, err := kubeConfig(ctx, uc.facility, uc.action, target.ScopeUUID, target.FacilityName)
	if err != nil {
		return nil, err
	}
	secret, err := uc.kubeCore.GetSecret(ctx, kc, tmp[0], tmp[1])
	if err != nil {
		return nil, err
	}
	return uc.warpJobSpec(&WarpTargetExternal{
		Endpoint:  target.Endpoint,
		AccessKey: string(secret.Data["root-user"]),
		SecretKey: string(secret.Data["root-password"]),
	}, input), nil
}

func (uc *BISTUseCase) warpJobSpec(target *WarpTargetExternal, input *WarpInput) *JobSpec {
	bistJobBackoffLimit := int32(2)
	env := []corev1.EnvVar{
		{Name: "BENCHMARK_TYPE", Value: "s3"},
		{Name: "BENCHMARK_ARGS_WARP_HOST", Value: target.Endpoint},
		{Name: "BENCHMARK_ARGS_WARP_ACCESS_KEY", Value: target.AccessKey},
		{Name: "BENCHMARK_ARGS_WARP_SECRET_KEY", Value: target.SecretKey},
		{Name: "BENCHMARK_ARGS_WARP_ACTION", Value: input.Operation},
		{Name: "BENCHMARK_ARGS_WARP_DURATION", Value: strconv.FormatInt(input.Duration, 10) + "s"}, // with unit
		{Name: "BENCHMARK_ARGS_WARP_CONCURRENT", Value: "2"},
		{Name: "BENCHMARK_ARGS_WARP_OBJ.SIZE", Value: strconv.FormatInt(input.ObjectSize, 10)},
	}
	if input.Operation == http.MethodPut {
		env = append(env, corev1.EnvVar{Name: "BENCHMARK_ARGS_WARP_OBJECTS", Value: strconv.FormatInt(input.ObjectCount, 10)})
	}
	return &JobSpec{
		BackoffLimit: &bistJobBackoffLimit,
		Template: corev1.PodTemplateSpec{
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:            "bist-container",
						Image:           "docker.io/otterscale/bist-s3:v3",
						Command:         []string{"./start.sh"},
						Env:             env,
						ImagePullPolicy: corev1.PullIfNotPresent,
					},
				},
				RestartPolicy: corev1.RestartPolicyNever,
			},
		},
	}
}
