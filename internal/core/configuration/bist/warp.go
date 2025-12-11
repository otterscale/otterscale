package bist

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/otterscale/otterscale/internal/core/application/release"
	"github.com/otterscale/otterscale/internal/core/application/workload"
	"github.com/otterscale/otterscale/internal/core/scope"
	"github.com/otterscale/otterscale/internal/core/versions"
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
	Type  ObjectServiceType `json:"type"`
	Scope string            `json:"scope"`
	Name  string            `json:"name"`
	Host  string            `json:"host"`
}

type WarpTargetExternal struct {
	Host      string `json:"host"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
}

type WarpInput struct {
	Operation   WarpOperationType `json:"operation"`
	Duration    int64             `json:"duration"`
	ObjectSize  int64             `json:"object_size"`
	ObjectCount int64             `json:"object_count"`
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

func (uc *UseCase) CreateWarpResult(ctx context.Context, name, createdBy string, target WarpTarget, input *WarpInput) (*Result, error) {
	// namespace
	if err := uc.ensureNamespace(ctx, scope.ReservedName, bistNamespace); err != nil {
		return nil, err
	}

	// name
	if name == "" {
		name = fmt.Sprintf("bist-%s", shortID(strconv.FormatInt(time.Now().UnixNano(), 10)))
	}

	// annotations
	warp, _ := json.Marshal(Warp{
		Target: target,
		Input:  input,
	})

	// job
	job := &workload.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: bistNamespace,
			Labels: map[string]string{
				release.TypeLabel: "bist",
				kindLabel:         kindWarp,
				createdByLabel:    createdBy,
			},
			Annotations: map[string]string{
				warpAnnotation: string(warp),
			},
		},
	}

	// s3 internal
	internal := target.Internal
	if internal != nil {
		switch internal.Type {
		case ObjectServiceTypeCeph:
			job.Spec = uc.warpCephObjectGatewayJobSpec(internal, input)

		case ObjectServiceTypeMinIO:
			spec, err := uc.warpMinIOJobSpec(ctx, internal, input)
			if err != nil {
				return nil, err
			}
			job.Spec = spec

		default:
			return nil, fmt.Errorf("unsupported internal kind %q", internal.Type)
		}
	}

	// s3 external
	external := target.External
	if external != nil {
		job.Spec = uc.warpJobSpec(external, input)
	}

	// job
	job, err := uc.job.Create(ctx, scope.ReservedName, bistNamespace, job)
	if err != nil {
		return nil, err
	}

	return uc.toResult(ctx, job)
}

func (uc *UseCase) warpCephObjectGatewayJobSpec(target *WarpTargetInternal, input *WarpInput) batchv1.JobSpec {
	accessKey, secretKey := uc.bucket.Key(target.Scope)

	external := &WarpTargetExternal{
		Host:      target.Host, // without protocol prefix
		AccessKey: accessKey,
		SecretKey: secretKey,
	}

	return uc.warpJobSpec(external, input)
}

func (uc *UseCase) warpMinIOJobSpec(ctx context.Context, target *WarpTargetInternal, input *WarpInput) (batchv1.JobSpec, error) {
	tmp := strings.Split(target.Name, ".")
	if len(tmp) != 2 {
		return batchv1.JobSpec{}, fmt.Errorf("invalid name %q", target.Name)
	}

	secret, err := uc.secret.Get(ctx, target.Scope, tmp[0], tmp[1])
	if err != nil {
		return batchv1.JobSpec{}, err
	}

	external := &WarpTargetExternal{
		Host:      target.Host,
		AccessKey: string(secret.Data["root-user"]),
		SecretKey: string(secret.Data["root-password"]),
	}

	return uc.warpJobSpec(external, input), nil
}

func (uc *UseCase) warpJobSpec(target *WarpTargetExternal, input *WarpInput) batchv1.JobSpec {
	bistJobBackoffLimit := int32(2)

	env := []corev1.EnvVar{
		{Name: "BENCHMARK_TYPE", Value: "s3"},
		{Name: "BENCHMARK_ARGS_WARP_HOST", Value: target.Host},
		{Name: "BENCHMARK_ARGS_WARP_ACCESS_KEY", Value: target.AccessKey},
		{Name: "BENCHMARK_ARGS_WARP_SECRET_KEY", Value: target.SecretKey},
		{Name: "BENCHMARK_ARGS_WARP_ACTION", Value: input.Operation.String()},
		{Name: "BENCHMARK_ARGS_WARP_DURATION", Value: strconv.FormatInt(input.Duration, 10) + "s"}, // with unit
		{Name: "BENCHMARK_ARGS_WARP_CONCURRENT", Value: "2"},
		{Name: "BENCHMARK_ARGS_WARP_OBJ.SIZE", Value: strconv.FormatInt(input.ObjectSize, 10)},
	}

	if !strings.EqualFold(input.Operation.String(), http.MethodPut) {
		env = append(env, corev1.EnvVar{Name: "BENCHMARK_ARGS_WARP_OBJECTS", Value: strconv.FormatInt(input.ObjectCount, 10)})
	}

	return batchv1.JobSpec{
		BackoffLimit: &bistJobBackoffLimit,
		Template: corev1.PodTemplateSpec{
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:            "bist-container",
						Image:           fmt.Sprintf("ghcr.io/otterscale/built-in-self-test/bist-s3:v%s", versions.Bist),
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

func (uc *UseCase) unmarshalWarpOutput(ctx context.Context, job *workload.Job, val **WarpOutput) error {
	logs, err := uc.fetchLogs(ctx, job)
	if err != nil {
		return err
	}

	data, err := json.Marshal(logs)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &val)
}

func (uc *UseCase) toWarp(ctx context.Context, job *workload.Job) (*Warp, error) {
	warp := &Warp{}

	// target & input
	if err := json.Unmarshal([]byte(job.Annotations[warpAnnotation]), warp); err != nil {
		return nil, err
	}

	// output
	if err := uc.unmarshalWarpOutput(ctx, job, &warp.Output); err != nil {
		return nil, err
	}

	return warp, nil
}
