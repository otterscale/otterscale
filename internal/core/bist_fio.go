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
	Scope    string `json:"scope"`
	Facility string `json:"facility"`
}

type FIOTargetNFS struct {
	Endpoint string `json:"endpoint"`
	Path     string `json:"path"`
}

type FIOInput struct {
	AccessMode string `json:"access_mode"`
	JobCount   int64  `json:"job_count"`
	RunTime    int64  `json:"run_time"`
	BlockSize  int64  `json:"block_size"`
	FileSize   int64  `json:"file_size"`
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

func (uc *BISTUseCase) CreateFIOResult(ctx context.Context, name, createdBy string, target FIOTarget, input *FIOInput) (*BISTResult, error) {
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
	fio, _ := json.Marshal(FIO{
		Target: target,
		Input:  input,
	})
	annotations := map[string]string{
		bistAnnotationCreatedBy: createdBy,
		bistAnnotationKind:      bistKindFIO,
		bistAnnotationFIO:       string(fio),
	}

	// spec
	var spec *JobSpec

	// block
	block := target.Ceph
	if block != nil {
		// pool & image
		if err := uc.ensurePool(ctx, block.Scope, block.Facility, bistBlockPool); err != nil {
			return nil, err
		}
		if err := uc.ensureImage(ctx, block.Scope, block.Facility, bistBlockPool, bistBlockImage); err != nil {
			return nil, err
		}
		// config map
		configMapName := fmt.Sprintf("ceph-conf-%s", generateHashedName(block.Scope+block.Facility))
		if err := uc.ensureConfigMap(ctx, config, block.Scope, block.Facility, configMapName); err != nil {
			return nil, err
		}
		spec = uc.blockJobSpec(configMapName, input)
	}

	// nfs
	nfs := target.NFS
	if nfs != nil {
		spec = uc.nfsJobSpec(nfs, input)
	}

	// job
	job, err := uc.kubeBatch.CreateJob(ctx, config, bistNamespace, name, labels, annotations, spec)
	if err != nil {
		return nil, err
	}
	return uc.toBISTResult(ctx, config, job)
}

func (uc *BISTUseCase) blockJobSpec(configMapName string, input *FIOInput) *JobSpec {
	bistJobBackoffLimit := int32(2)
	privileged := true
	env := []corev1.EnvVar{
		{Name: "BENCHMARK_TYPE", Value: "block"},
		{Name: "BENCHMARK_ARGS_BLOCK_POOL", Value: bistBlockPool},
		{Name: "BENCHMARK_ARGS_BLOCK_IMAGE", Value: bistBlockImage},
		{Name: "BENCHMARK_ARGS_FIO_DIRECT", Value: "1"},
		{Name: "BENCHMARK_ARGS_FIO_FILESIZE", Value: strconv.FormatInt(input.FileSize, 10)},
		{Name: "BENCHMARK_ARGS_FIO_IODEPTH", Value: strconv.FormatInt(input.IODepth, 10)},
		{Name: "BENCHMARK_ARGS_FIO_NUMJOBS", Value: strconv.FormatInt(input.JobCount, 10)},
		{Name: "BENCHMARK_ARGS_FIO_GROUP_REPORTING", Value: "True"},
		{Name: "BENCHMARK_ARGS_FIO_NORANDOMMAP", Value: "True"},
		{Name: "BENCHMARK_ARGS_FIO_BS", Value: strconv.FormatInt(input.BlockSize, 10)},
		{Name: "BENCHMARK_ARGS_FIO_BUFFER_COMPRESS_PERCENTAGE", Value: "0"},
		{Name: "BENCHMARK_ARGS_FIO_RW", Value: input.AccessMode},
		{Name: "BENCHMARK_ARGS_FIO_STARTDELAY", Value: "5"},
		{Name: "BENCHMARK_ARGS_FIO_TIME_BASED", Value: "True"},
		{Name: "BENCHMARK_ARGS_FIO_EXITALL_ON_ERROR", Value: "True"},
		{Name: "BENCHMARK_ARGS_FIO_CREATE_SERIALIZE", Value: "0"},
		{Name: "BENCHMARK_ARGS_FIO_RUNTIME", Value: strconv.FormatInt(input.RunTime, 10)},
	}
	return &JobSpec{
		BackoffLimit: &bistJobBackoffLimit,
		Template: corev1.PodTemplateSpec{
			Spec: corev1.PodSpec{
				Volumes: []corev1.Volume{
					{
						Name: "ceph-conf",
						VolumeSource: corev1.VolumeSource{
							ConfigMap: &corev1.ConfigMapVolumeSource{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: configMapName,
								},
								Items: []corev1.KeyToPath{{Key: "ceph.conf", Path: "ceph.conf"}},
							},
						},
					},
					{
						Name: "dev",
						VolumeSource: corev1.VolumeSource{
							HostPath: &corev1.HostPathVolumeSource{Path: "/dev"},
						},
					},
					{
						Name: "modules",
						VolumeSource: corev1.VolumeSource{
							HostPath: &corev1.HostPathVolumeSource{Path: "/lib/modules"},
						},
					},
					{
						Name: "run-udev",
						VolumeSource: corev1.VolumeSource{
							HostPath: &corev1.HostPathVolumeSource{Path: "/run/udev"},
						},
					},
					{
						Name: "sys",
						VolumeSource: corev1.VolumeSource{
							HostPath: &corev1.HostPathVolumeSource{Path: "/sys"},
						},
					},
				},
				Containers: []corev1.Container{
					{
						Name:    "bist-container",
						Image:   "docker.io/otterscale/bist-block:v1",
						Command: []string{"./start.sh"},
						Env:     env,
						VolumeMounts: []corev1.VolumeMount{
							{Name: "ceph-conf", MountPath: "/etc/ceph/ceph.conf", SubPath: "ceph.conf"},
							{Name: "dev", MountPath: "/dev"},
							{Name: "modules", MountPath: "/lib/modules"},
							{Name: "run-udev", MountPath: "/run/udev"},
							{Name: "sys", MountPath: "/sys"},
						},
						ImagePullPolicy: corev1.PullIfNotPresent,
						SecurityContext: &corev1.SecurityContext{Privileged: &privileged},
					},
				},
				RestartPolicy: corev1.RestartPolicyNever,
			},
		},
	}
}

func (uc *BISTUseCase) nfsJobSpec(target *FIOTargetNFS, input *FIOInput) *JobSpec {
	bistJobBackoffLimit := int32(2)
	privileged := true
	env := []corev1.EnvVar{
		{Name: "BENCHMARK_TYPE", Value: "nfs"},
		{Name: "BENCHMARK_ARGS_NFS_ENDPOINT", Value: target.Endpoint},
		{Name: "BENCHMARK_ARGS_NFS_PATH", Value: target.Path},
		{Name: "BENCHMARK_ARGS_FIO_DIRECT", Value: "1"},
		{Name: "BENCHMARK_ARGS_FIO_FILESIZE", Value: strconv.FormatInt(input.FileSize, 10)},
		{Name: "BENCHMARK_ARGS_FIO_IODEPTH", Value: strconv.FormatInt(input.IODepth, 10)},
		{Name: "BENCHMARK_ARGS_FIO_NUMJOBS", Value: strconv.FormatInt(input.JobCount, 10)},
		{Name: "BENCHMARK_ARGS_FIO_GROUP_REPORTING", Value: "True"},
		{Name: "BENCHMARK_ARGS_FIO_NORANDOMMAP", Value: "True"},
		{Name: "BENCHMARK_ARGS_FIO_BS", Value: strconv.FormatInt(input.BlockSize, 10)},
		{Name: "BENCHMARK_ARGS_FIO_BUFFER_COMPRESS_PERCENTAGE", Value: "0"},
		{Name: "BENCHMARK_ARGS_FIO_RW", Value: input.AccessMode},
		{Name: "BENCHMARK_ARGS_FIO_STARTDELAY", Value: "5"},
		{Name: "BENCHMARK_ARGS_FIO_TIME_BASED", Value: "True"},
		{Name: "BENCHMARK_ARGS_FIO_EXITALL_ON_ERROR", Value: "True"},
		{Name: "BENCHMARK_ARGS_FIO_CREATE_SERIALIZE", Value: "0"},
		{Name: "BENCHMARK_ARGS_FIO_RUNTIME", Value: strconv.FormatInt(input.RunTime, 10)},
	}
	return &JobSpec{
		BackoffLimit: &bistJobBackoffLimit,
		Template: corev1.PodTemplateSpec{
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:            "bist-container",
						Image:           "docker.io/otterscale/bist-nfs:v1",
						Command:         []string{"./start.sh"},
						Env:             env,
						ImagePullPolicy: corev1.PullIfNotPresent,
						SecurityContext: &corev1.SecurityContext{Privileged: &privileged},
					},
				},
				RestartPolicy: corev1.RestartPolicyNever,
			},
		},
	}
}
