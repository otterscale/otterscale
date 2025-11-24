package bist

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/otterscale/otterscale/internal/core/application/cluster"
	"github.com/otterscale/otterscale/internal/core/application/release"
	"github.com/otterscale/otterscale/internal/core/application/workload"
	"github.com/otterscale/otterscale/internal/core/scope"
	"github.com/otterscale/otterscale/internal/core/storage"
)

const (
	bistBlockPool  = "bist_pool"
	bistBlockImage = "bist_image"
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
	Scope string `json:"scope"`
}

type FIOTargetNFS struct {
	Host string `json:"host"`
	Path string `json:"path"`
}

type FIOInput struct {
	AccessMode FIOAccessMode `json:"access_mode"`
	JobCount   int64         `json:"job_count"`
	RunTime    int64         `json:"run_time"`
	BlockSize  int64         `json:"block_size"`
	FileSize   int64         `json:"file_size"`
	IODepth    int64         `json:"io_depth"`
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

func (uc *UseCase) CreateFIOResult(ctx context.Context, name, createdBy string, target FIOTarget, input *FIOInput) (*Result, error) {
	// namespace
	if err := uc.ensureNamespace(ctx, scope.ReservedName, bistNamespace); err != nil {
		return nil, err
	}

	// name
	if name == "" {
		name = fmt.Sprintf("bist-%s", shortID(strconv.FormatInt(time.Now().UnixNano(), 10)))
	}

	// annotations
	fio, _ := json.Marshal(FIO{
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
				kindLabel:         kindFIO,
				createdByLabel:    createdBy,
			},
			Annotations: map[string]string{
				fioAnnotation: string(fio),
			},
		},
	}

	// block
	block := target.Ceph
	if block != nil {
		if err := uc.ensurePool(ctx, block.Scope, bistBlockPool); err != nil {
			return nil, err
		}

		if err := uc.ensureImage(ctx, block.Scope, bistBlockPool, bistBlockImage); err != nil {
			return nil, err
		}

		configMapName := fmt.Sprintf("ceph-conf-%s", shortID(block.Scope))
		if err := uc.ensureConfigMap(ctx, scope.ReservedName, bistNamespace, configMapName); err != nil {
			return nil, err
		}

		job.Spec = uc.blockJobSpec(configMapName, input)
	}

	// nfs
	nfs := target.NFS
	if nfs != nil {
		job.Spec = uc.nfsJobSpec(nfs, input)
	}

	// job
	job, err := uc.job.Create(ctx, scope.ReservedName, bistNamespace, job)
	if err != nil {
		return nil, err
	}

	return uc.toResult(ctx, job)
}

func (uc *UseCase) blockJobSpec(configMapName string, input *FIOInput) batchv1.JobSpec {
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
		{Name: "BENCHMARK_ARGS_FIO_RW", Value: input.AccessMode.String()},
		{Name: "BENCHMARK_ARGS_FIO_STARTDELAY", Value: "5"},
		{Name: "BENCHMARK_ARGS_FIO_TIME_BASED", Value: "True"},
		{Name: "BENCHMARK_ARGS_FIO_EXITALL_ON_ERROR", Value: "True"},
		{Name: "BENCHMARK_ARGS_FIO_CREATE_SERIALIZE", Value: "0"},
		{Name: "BENCHMARK_ARGS_FIO_RUNTIME", Value: strconv.FormatInt(input.RunTime, 10)},
	}

	return batchv1.JobSpec{
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

func (uc *UseCase) nfsJobSpec(target *FIOTargetNFS, input *FIOInput) batchv1.JobSpec {
	bistJobBackoffLimit := int32(2)
	privileged := true
	env := []corev1.EnvVar{
		{Name: "BENCHMARK_TYPE", Value: "nfs"},
		{Name: "BENCHMARK_ARGS_NFS_ENDPOINT", Value: target.Host},
		{Name: "BENCHMARK_ARGS_NFS_PATH", Value: target.Path},
		{Name: "BENCHMARK_ARGS_FIO_DIRECT", Value: "1"},
		{Name: "BENCHMARK_ARGS_FIO_FILESIZE", Value: strconv.FormatInt(input.FileSize, 10)},
		{Name: "BENCHMARK_ARGS_FIO_IODEPTH", Value: strconv.FormatInt(input.IODepth, 10)},
		{Name: "BENCHMARK_ARGS_FIO_NUMJOBS", Value: strconv.FormatInt(input.JobCount, 10)},
		{Name: "BENCHMARK_ARGS_FIO_GROUP_REPORTING", Value: "True"},
		{Name: "BENCHMARK_ARGS_FIO_NORANDOMMAP", Value: "True"},
		{Name: "BENCHMARK_ARGS_FIO_BS", Value: strconv.FormatInt(input.BlockSize, 10)},
		{Name: "BENCHMARK_ARGS_FIO_BUFFER_COMPRESS_PERCENTAGE", Value: "0"},
		{Name: "BENCHMARK_ARGS_FIO_RW", Value: input.AccessMode.String()},
		{Name: "BENCHMARK_ARGS_FIO_STARTDELAY", Value: "5"},
		{Name: "BENCHMARK_ARGS_FIO_TIME_BASED", Value: "True"},
		{Name: "BENCHMARK_ARGS_FIO_EXITALL_ON_ERROR", Value: "True"},
		{Name: "BENCHMARK_ARGS_FIO_CREATE_SERIALIZE", Value: "0"},
		{Name: "BENCHMARK_ARGS_FIO_RUNTIME", Value: strconv.FormatInt(input.RunTime, 10)},
	}

	return batchv1.JobSpec{
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

func (uc *UseCase) ensureNamespace(ctx context.Context, scope, namespace string) error {
	_, err := uc.namespace.Get(ctx, scope, namespace)
	if apierrors.IsNotFound(err) {
		namespace := &cluster.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: bistNamespace,
			},
		}

		_, err := uc.namespace.Create(ctx, scope, namespace)
		return err
	}

	return err
}

func (uc *UseCase) ensureConfigMap(ctx context.Context, scope, namespace, name string) error {
	_, err := uc.configMap.Get(ctx, scope, namespace, name)
	if apierrors.IsNotFound(err) {
		host, id, key, err := uc.node.Config(scope)
		if err != nil {
			return err
		}

		configMap := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
			Data: map[string]string{
				"ceph.conf": fmt.Sprintf(`[global]\nmon host = %s\nfsid = %s\nkey = %s`, host, id, key),
			},
		}

		_, err = uc.configMap.Create(ctx, scope, namespace, configMap)
		return err
	}

	return err
}

func (uc *UseCase) ensurePool(ctx context.Context, scope, pool string) error {
	pools, err := uc.pool.List(ctx, scope, "")
	if err != nil {
		return err
	}

	for i := range pools {
		if pools[i].Name == pool {
			return nil
		}
	}

	if err := uc.pool.Create(ctx, scope, pool, storage.PoolTypeReplicated); err != nil {
		return err
	}

	return uc.pool.Enable(ctx, scope, pool, "rbd")
}

func (uc *UseCase) ensureImage(ctx context.Context, scope, pool, image string) error {
	images, err := uc.image.List(ctx, scope, pool)
	if err != nil {
		return err
	}

	for i := range images {
		if images[i].Name == image {
			return nil
		}
	}

	objectSizeBytes := 4194304
	stripeUnitBytes := uint64(4194304) //nolint:mnd // default 4MB
	stripeCount := uint64(1)
	size := uint64(10737418240) //nolint:mnd // default 10GB
	order := int(math.Round(math.Log2(float64(objectSizeBytes))))
	features := uint64(61) //nolint:mnd // Layering, ExclusiveLock, ObjectMap, FastDiff, DeepFlatten

	return uc.image.Create(ctx, scope, pool, image, order, stripeUnitBytes, stripeCount, size, features)
}

func (uc *UseCase) unmarshalFIOOutput(ctx context.Context, job *workload.Job, val **FIOOutput) error {
	logs, err := uc.fetchLogs(ctx, job)
	if err != nil {
		return err
	}

	jobs, ok := logs["jobs"].([]any)
	if ok {
		data, err := json.Marshal(jobs[0])
		if err != nil {
			return err
		}

		return json.Unmarshal(data, &val)
	}

	return nil
}

func (uc *UseCase) toFIO(ctx context.Context, job *workload.Job) (*FIO, error) {
	fio := &FIO{}

	// target & input
	if err := json.Unmarshal([]byte(job.Annotations[fioAnnotation]), fio); err != nil {
		return nil, err
	}

	// output
	if err := uc.unmarshalFIOOutput(ctx, job, &fio.Output); err != nil {
		return nil, err
	}

	if fio.Output != nil {
		if fio.Output.Read != nil && fio.Output.Read.IOBytes == 0 {
			fio.Output.Read = nil
		}

		if fio.Output.Write != nil && fio.Output.Write.IOBytes == 0 {
			fio.Output.Write = nil
		}

		if fio.Output.Trim != nil && fio.Output.Trim.IOBytes == 0 {
			fio.Output.Trim = nil
		}
	}

	return fio, nil
}
