package core

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
)

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
		bistAnnotationKind:      bistAnnotationWarp,
		bistAnnotationFIO:       string(warp),
	}

	// spec
	var spec *JobSpec

	// s3 internal
	internal := target.Internal
	if internal != nil {
		switch internal.Type {
		case "ceph":
			// internal.Name
			spec = uc.internalS3JobSpec(internal.Endpoint, "", "", input)
		case "minio":
			spec = uc.internalS3JobSpec(internal.Endpoint, "", "", input)
		default:
			return nil, fmt.Errorf("unsupported kind %q", internal.Type)
		}
	}

	// s3 external
	external := target.External
	if external != nil {
		spec = uc.externalS3JobSpec(external, input)
	}

	// job
	job, err := uc.kubeBatch.CreateJob(ctx, config, bistNamespace, name, labels, annotations, spec)
	if err != nil {
		return nil, err
	}
	return uc.toBISTResult(ctx, config, job)
}

func (uc *BISTUseCase) internalS3JobSpec(endpoint, accessKey, secretKey string, input *WarpInput) *JobSpec {
	bistJobBackoffLimit := int32(2)
	env := []corev1.EnvVar{
		{Name: "BENCHMARK_TYPE", Value: "s3"},
		{Name: "BENCHMARK_ARGS_WARP_HOST", Value: endpoint},
		{Name: "BENCHMARK_ARGS_WARP_ACCESS_KEY", Value: accessKey},
		{Name: "BENCHMARK_ARGS_WARP_SECRET_KEY", Value: secretKey},
		{Name: "BENCHMARK_ARGS_WARP_ACTION", Value: input.Operation},
		{Name: "BENCHMARK_ARGS_WARP_DURATION", Value: input.Duration},
		{Name: "BENCHMARK_ARGS_WARP_CONCURRENT", Value: "2"},
		{Name: "BENCHMARK_ARGS_WARP_OBJ.SIZE", Value: input.ObjectSize},
		{Name: "BENCHMARK_ARGS_WARP_OBJECTS", Value: input.ObjectCount},
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

func (uc *BISTUseCase) externalS3JobSpec(target *WarpTargetExternal, input *WarpInput) *JobSpec {
	bistJobBackoffLimit := int32(2)
	env := []corev1.EnvVar{
		{Name: "BENCHMARK_TYPE", Value: "s3"},
		{Name: "BENCHMARK_ARGS_WARP_HOST", Value: target.Endpoint},
		{Name: "BENCHMARK_ARGS_WARP_ACCESS_KEY", Value: target.AccessKey},
		{Name: "BENCHMARK_ARGS_WARP_SECRET_KEY", Value: target.SecretKey},
		{Name: "BENCHMARK_ARGS_WARP_ACTION", Value: input.Operation},
		{Name: "BENCHMARK_ARGS_WARP_DURATION", Value: input.Duration},
		{Name: "BENCHMARK_ARGS_WARP_CONCURRENT", Value: "2"},
		{Name: "BENCHMARK_ARGS_WARP_OBJ.SIZE", Value: input.ObjectSize},
		{Name: "BENCHMARK_ARGS_WARP_OBJECTS", Value: input.ObjectCount},
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
