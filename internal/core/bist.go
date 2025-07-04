package core

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	pb "github.com/openhdc/otterscale/api/bist/v1"
	"golang.org/x/sync/errgroup"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"

	kuberr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openhdc/otterscale/internal/config"
	"k8s.io/client-go/rest"
)

const (
	BIST_NAME      = "bist"
	BIST_NAMESPACE = "bist"
	BIST_LABEL     = "app.otterscale.io/app=bist"
	BIST_JOB_RETRY = 3

	BIST_PHASE_COMPLETE = "COMPLETE"
	BIST_PHASE_RUNNING  = "RUNNING"
	BIST_PHASE_FAIL     = "FAILED"
	BIST_PHASE_CREATING = "CREATING"
)

var (
	BIST_JOB_LABEL = map[string]string{
		"app.otterscale.io/app": "bist",
	}
	BIST_IMAGE_LIST = map[string]string{
		"block": "docker.io/otterscale/bist-block:v1",
		"nfs":   "docker.io/otterscale/bist-nfs:v1",
		"s3":    "docker.io/otterscale/bist-s3:v2",
	}
	BIST_WARP_ENV = map[string]string{
		"BENCHMARK_TYPE":                 "s3",
		"BENCHMARK_ARGS_WARP_HOST":       "",
		"BENCHMARK_ARGS_WARP_ACCESS_KEY": "",
		"BENCHMARK_ARGS_WARP_SECRET_KEY": "",
		"BENCHMARK_ARGS_WARP_ACTION":     "",
		"BENCHMARK_ARGS_WARP_DURATION":   "",
		"BENCHMARK_ARGS_WARP_CONCURRENT": "2",
		"BENCHMARK_ARGS_WARP_OBJ.SIZE":   "",
	}

	BIST_FIO_ENV = map[string]string{
		"BENCHMARK_TYPE":                                "",
		"BENCHMARK_ARGS_FIO_DIRECT":                     "1",
		"BENCHMARK_ARGS_FIO_FILESIZE":                   "",
		"BENCHMARK_ARGS_FIO_IODEPTH":                    "",
		"BENCHMARK_ARGS_FIO_NUMJOBS":                    "",
		"BENCHMARK_ARGS_FIO_GROUP_REPORTING":            "True",
		"BENCHMARK_ARGS_FIO_RWMIXREAD":                  "100",
		"BENCHMARK_ARGS_FIO_NORANDOMMAP":                "True",
		"BENCHMARK_ARGS_FIO_BS":                         "",
		"BENCHMARK_ARGS_FIO_BUFFER_COMPRESS_PERCENTAGE": "0",
		"BENCHMARK_ARGS_FIO_RW":                         "",
		"BENCHMARK_ARGS_FIO_STARTDELAY":                 "5",
		"BENCHMARK_ARGS_FIO_TIME_BASED":                 "True",
		"BENCHMARK_ARGS_FIO_EXITALL_ON_ERROR":           "True",
		"BENCHMARK_ARGS_FIO_CREATE_SERIALIZE":           "0",
		"BENCHMARK_ARGS_FIO_RUNTIME":                    "",
	}
)

type BISTResult struct {
	Type         string
	Name         string
	Status       string
	StartTime    *metav1.Time
	CompleteTime *metav1.Time
	Args         string
	Logs         []string
	FIO          *BISTFIO
	Warp         *BISTWarp
}

type BISTFIO struct {
	AccessMode       string
	StorageClassName string
	NFSEndpoint      string
	NFSPath          string
	JobCount         uint64
	RunTime          string
	BlockSize        string
	FileSize         string
	IODepth          uint64
}

type BISTWarp struct {
	Operation  string `json:"operation"`
	Endpoint   string `json:"endpoint"`
	AccessKey  string `json:"access_key"`
	SecretKey  string `json:"secret_key"`
	Duration   string `json:"duration"`
	ObjectSize string `json:"object_size.SIZE"`
}

type BISTBlock struct {
	FacilityName     string
	StorageClassName string
}

type BISTS3 struct {
	Type     string
	Name     string
	Endpoint string
}

type BISTUseCase struct {
	action      ActionRepo
	scope       ScopeRepo
	kubeStorage KubeStorageRepo
	client      ClientRepo
	facility    FacilityRepo
	kubeCore    KubeCoreRepo
	kubeBatch   KubeBatchRepo
	conf        *config.Config
	configs     sync.Map
}

func NewBISTUseCase(action ActionRepo, scope ScopeRepo, kubeBatch KubeBatchRepo, kube KubeAppsRepo, kubeCore KubeCoreRepo, kubeStorage KubeStorageRepo, facility FacilityRepo, conf *config.Config, client ClientRepo) *BISTUseCase {
	return &BISTUseCase{
		action:      action,
		client:      client,
		scope:       scope,
		kubeCore:    kubeCore,
		kubeBatch:   kubeBatch,
		kubeStorage: kubeStorage,
		facility:    facility,
		conf:        conf,
	}
}

func (uc *BISTUseCase) ListResults(ctx context.Context) ([]BISTResult, error) {
	config, err := uc.newConfig(ctx)
	if err != nil {
		return nil, err
	}

	jobs, err := uc.kubeBatch.ListJobsByLabel(ctx, config, BIST_NAMESPACE, BIST_LABEL)
	if err != nil {
		return nil, err
	}

	bists := []BISTResult{}
	for _, job := range jobs {
		bist, err := uc.toBISTResult(ctx, job.Name, &job, config)
		if err != nil {
			return nil, err
		}
		bists = append(bists, *bist)
	}

	return bists, nil
}

func (uc *BISTUseCase) CreateResult(ctx context.Context, target, name string, fio *BISTFIO, warp *BISTWarp) (*BISTResult, error) {
	config, err := uc.newConfig(ctx)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, fmt.Errorf("invalid config")
	}

	_, err = uc.kubeCore.GetNamespace(ctx, config, BIST_NAMESPACE)
	if err != nil {
		if kuberr.IsNotFound(err) {
			// Create namespace if not exist
			namespace := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: BIST_NAMESPACE,
				},
			}
			_, err = uc.kubeCore.CreateNamespace(ctx, config, namespace)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	var job *batchv1.Job
	if target == "block" {
		job = toBISTJob(name, BIST_NAMESPACE, target, BIST_JOB_RETRY, fio)
	} else if target == "nfs" {
		job = toBISTJob(name, BIST_NAMESPACE, target, BIST_JOB_RETRY, fio)
	} else if target == "s3" {
		job = toBISTJob(name, BIST_NAMESPACE, target, BIST_JOB_RETRY, warp)
	}

	if job == nil {
		return nil, fmt.Errorf("Invalid bist job specified")
	}

	job, err = uc.kubeBatch.CreateJob(ctx, config, job)
	if err != nil {
		return nil, err
	}

	bist, err := uc.toBISTResult(ctx, job.Name, job, config)

	return bist, nil
}

func (uc *BISTUseCase) DeleteResult(ctx context.Context, name string) error {
	config, err := uc.newConfig(ctx)
	if err != nil {
		return err
	}
	return uc.kubeBatch.DeleteJob(ctx, config, BIST_NAMESPACE, name)
}

func (uc *BISTUseCase) ListBlocks(ctx context.Context) ([]BISTBlock, error) {

	config, err := uc.newConfig(ctx)
	if err != nil {
		return nil, err
	}
	sc, err := uc.kubeStorage.ListStorageClassesByLabel(ctx, config, "juju.io/manifest=rbd")

	return toBISTBlocks(sc), nil
}

func (uc *BISTUseCase) ListS3s(ctx context.Context, uuid string) ([]BISTS3, error) {
	//
	scopes, err := uc.scope.List(ctx)
	if err != nil {
		return nil, err
	}
	scopes = slices.DeleteFunc(scopes, func(s Scope) bool {
		return !strings.Contains(s.UUID, uuid)
	})

	minioEndpoints := make([]BISTS3, 0, len(scopes))
	eg, ctx := errgroup.WithContext(ctx)
	for i := range scopes {
		eg.Go(func() error {
			s, err := uc.client.Status(ctx, scopes[i].UUID, []string{"application", "*"})
			if err != nil {
				return err
			}
			for name := range s.Applications {
				if !strings.Contains(s.Applications[name].Charm, "kubernetes-control-plane") {
					continue
				}
				units := []EssentialUnit{}
				for uname := range s.Applications[name].Units {
					units = append(units, EssentialUnit{
						Name:      uname,
						Directive: s.Applications[name].Units[uname].Machine,
					})
				}

				for _, unit := range units {
					config, err := uc.config(ctx, uuid, removeLastSlashAndAfter(unit.Name))
					if config == nil {
						fmt.Println("nil config")
						continue
					}
					label := "app.otterscale.io/release-name=minio"
					field := "spec.type=NodePort"
					svcs, err := uc.kubeCore.ListServicesByOptions(ctx, config, "", label, field)
					if err != nil {
						return err
					}
					for _, svc := range svcs {
						for _, port := range svc.Spec.Ports {
							if port.Name == "minio-api" {
								minioEndpoints = append(minioEndpoints, toBISTS3("minio", svc.GetName()+"/"+svc.GetNamespace(), fmt.Sprintf("%s:%d", config.Host, port.NodePort)))
							}
						}
					}
				}
				break
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	cephEndpoints := make([]BISTS3, 0, len(scopes))

	for i := range scopes {
		eg.Go(func() error {
			s, err := uc.client.Status(ctx, scopes[i].UUID, []string{"application", "*"})
			if err != nil {
				return err
			}
			for name := range s.Applications {
				if !strings.Contains(s.Applications[name].Charm, "ceph-mon") {
					continue
				}
				units := []EssentialUnit{}
				for uname := range s.Applications[name].Units {
					units = append(units, EssentialUnit{
						Name:      uname,
						Directive: s.Applications[name].Units[uname].Machine,
					})
				}
				for _, unit := range units {
					leader, err := uc.facility.GetLeader(ctx, uuid, removeLastSlashAndAfter(unit.Name))
					if err != nil {
						return err
					}

					info, err := uc.facility.GetUnitInfo(ctx, uuid, leader)
					if err != nil {
						return err
					}

					cephEndpoints = append(cephEndpoints, toBISTS3("ceph", unit.Name, fmt.Sprintf("http://%s", info.PublicAddress)))
				}
				break
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	s3Endpoints := append(cephEndpoints, minioEndpoints...)

	// s3 endpoints
	return s3Endpoints, nil
}

func (uc *BISTUseCase) newConfig(ctx context.Context) (*rest.Config, error) {

	decodedBytes, _ := base64.StdEncoding.DecodeString(uc.conf.MicroK8s.Token)

	return &rest.Config{
		Host:        uc.conf.MicroK8s.Host,
		BearerToken: string(decodedBytes),
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true,
		},
	}, nil
}

func (uc *BISTUseCase) config(ctx context.Context, uuid, name string) (*rest.Config, error) {
	key := uuid + "/" + name

	if v, ok := uc.configs.Load(key); ok {
		return v.(*rest.Config), nil
	}

	config, err := uc.newKubeConfig(ctx, uuid, name)
	if err != nil {
		return nil, err
	}

	uc.configs.Store(key, config)

	return config, nil
}

func (uc *BISTUseCase) newKubeConfig(ctx context.Context, uuid, name string) (*rest.Config, error) {
	// kubernetes-control-plane
	leader, err := uc.facility.GetLeader(ctx, uuid, name)
	if err != nil {
		return nil, err
	}
	unitInfo, err := uc.facility.GetUnitInfo(ctx, uuid, leader)
	if err != nil {
		return nil, err
	}
	kubeControl, err := extractWorkerUnitName(unitInfo)
	if err != nil {
		return nil, err
	}

	// [TODO] Replace the last "/" and the char after
	kubeControl = removeLastSlashAndAfter(kubeControl)

	// kubernetes-worker
	leader, err = uc.facility.GetLeader(ctx, uuid, kubeControl)
	if err != nil {
		return nil, err
	}
	unitInfo, err = uc.facility.GetUnitInfo(ctx, uuid, leader)
	if err != nil {
		return nil, err
	}

	// config
	endpoint, err := extractEndpoint(unitInfo)
	if err != nil {
		return nil, err
	}
	clientToken, err := extractClientToken(unitInfo)
	if err != nil {
		return nil, err
	}
	return &rest.Config{
		Host:        endpoint,
		BearerToken: clientToken,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true,
		},
	}, nil
}

func removeLastSlashAndAfter(s string) string {
	lastSlashIndex := strings.LastIndex(s, "/")

	if lastSlashIndex == -1 {
		return s
	}
	return s[:lastSlashIndex]
}

func toBISTS3(s3type, name, endpoint string) BISTS3 {
	return BISTS3{
		Type:     s3type,
		Name:     name,
		Endpoint: endpoint,
	}
}

func toBISTBlocks(scs []StorageClass) []BISTBlock {
	ret := []BISTBlock{}
	for i := range scs {
		ret = append(ret, toBISTBlock(&scs[i]))
	}
	return ret
}

func toBISTBlock(sc *StorageClass) BISTBlock {

	fName := sc.GetLabels()["juju.io/application"]
	return BISTBlock{
		FacilityName:     fName,
		StorageClassName: sc.GetName(),
	}
}

func toBISTJob(name, namespace, target string, retry int32, BIST interface{}) *batchv1.Job {

	containerName := name + "-container"
	trueVal := true
	image := BIST_IMAGE_LIST[target]

	envJSON, err := json.Marshal(BIST)
	if err != nil {
		return nil
	}
	annotation := make(map[string]string)
	annotation["otterscale/bist-config"] = string(envJSON)

	BIST_JOB_LABEL["app.otterscale.io/bist"] = target

	jobEnv, _ := toJobEnv(target, BIST)

	//jobName := generateNameWithDateTime(name)
	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Labels:      BIST_JOB_LABEL,
			Annotations: annotation,
		},
		Spec: batchv1.JobSpec{
			BackoffLimit: &retry,
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            containerName,
							Image:           image,
							ImagePullPolicy: corev1.PullIfNotPresent,
							Command: []string{
								"./start.sh",
							},
							SecurityContext: &corev1.SecurityContext{
								Privileged: &trueVal,
							},
							Env:          jobEnv,
							VolumeMounts: generateVolumeMounts(target),
						},
					},
					Volumes:       generateVolume(target),
					RestartPolicy: corev1.RestartPolicyNever,
				},
			},
		},
	}
}

func generateVolume(target string) []corev1.Volume {

	switch target {
	case "block":
		return []corev1.Volume{
			{
				Name: "dev",
				VolumeSource: corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{
						Path: "/dev",
					},
				},
			},
			{
				Name: "modules",
				VolumeSource: corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{
						Path: "/lib/modules",
					},
				},
			},
			{
				Name: "run-udev",
				VolumeSource: corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{
						Path: "/run/udev",
					},
				},
			},
		}
	case "nfs":
		return []corev1.Volume{}
	case "s3":
		return []corev1.Volume{}
	}
	return []corev1.Volume{}
}

func generateVolumeMounts(target string) []corev1.VolumeMount {

	switch target {
	case "nfs":
		return []corev1.VolumeMount{}
	case "block":
		return []corev1.VolumeMount{
			{
				Name:      "dev",
				MountPath: "/dev",
			},
			{
				Name:      "modules",
				MountPath: "/lib/modules",
			},
			{
				Name:      "run-udev",
				MountPath: "/run/udev",
			},
		}
	case "s3":
		return []corev1.VolumeMount{}
	default:
		return []corev1.VolumeMount{}
	}
}

func generateNameWithDateTime(baseName string) string {
	now := time.Now()
	dateTimeStr := now.Format("20060102150405")
	return fmt.Sprintf("%s-%s", baseName, dateTimeStr)
}

func (uc *BISTUseCase) toBISTResult(ctx context.Context, name string, job *batchv1.Job, config *rest.Config) (*BISTResult, error) {

	var jobStatus string
	var ret BISTResult

	if job.Status.Succeeded == 1 && job.Status.CompletionTime != nil {
		jobStatus = BIST_PHASE_COMPLETE
	} else if job.Status.Failed == 1 {
		jobStatus = BIST_PHASE_FAIL
	} else if job.Status.Active == 1 && job.Status.Ready != nil {
		jobStatus = BIST_PHASE_RUNNING
	} else {
		jobStatus = BIST_PHASE_CREATING
	}

	//jsonBytes, err := json.Marshal(job.Spec.Template.Spec.Containers[0].Env)

	// [TODO] test result
	selector, err := metav1.LabelSelectorAsSelector(job.Spec.Selector)
	if err != nil {
		return nil, err
	}

	pods, err := uc.kubeCore.ListPodsByLabel(ctx, config, BIST_NAMESPACE, selector.String())
	if err != nil {
		return nil, err
	}

	podLog := make([]string, 0, len(pods))
	for _, pod := range pods {
		logs, err := uc.kubeCore.GetPodLogs(ctx, pod, config, BIST_NAMESPACE)
		if err != nil {
			return nil, err
		}
		podLog = append(podLog, logs)
	}

	ret.Name = name
	ret.Type = job.GetLabels()["app.otterscale.io/bist"]
	ret.Status = jobStatus
	ret.Args = job.GetAnnotations()["otterscale/bist-config"]
	ret.StartTime = job.Status.StartTime
	ret.CompleteTime = job.Status.CompletionTime
	ret.Logs = podLog

	target := toBISTType(job.GetLabels()["app.otterscale.io/bist"])
	if target == pb.TestResult_S3 {
		err := json.Unmarshal([]byte(job.GetAnnotations()["otterscale/bist-config"]), &ret.Warp)
		if err != nil {
			return nil, err
		}
	} else {
		err := json.Unmarshal([]byte(job.GetAnnotations()["otterscale/bist-config"]), &ret.FIO)
		if err != nil {
			return nil, err
		}
	}

	return &ret, nil
}

func toBISTType(s string) pb.TestResult_Type {
	switch s {
	case "s3":
		return pb.TestResult_S3
	case "block":
		return pb.TestResult_BLOCK
	case "nfs":
		return pb.TestResult_NFS
	}
	return pb.TestResult_UNSPECIFIED
}

func toJobEnv(target string, BIST interface{}) ([]corev1.EnvVar, error) {

	var jobEnv []corev1.EnvVar
	v := reflect.ValueOf(BIST)
	t := reflect.TypeOf(BIST)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	switch s := BIST.(type) {
	case *BISTWarp:
		BIST_WARP_ENV["BENCHMARK_ARGS_WARP_HOST"] = s.Endpoint
		BIST_WARP_ENV["BENCHMARK_ARGS_WARP_ACCESS_KEY"] = s.AccessKey
		BIST_WARP_ENV["BENCHMARK_ARGS_WARP_SECRET_KEY"] = s.SecretKey
		BIST_WARP_ENV["BENCHMARK_ARGS_WARP_ACTION"] = s.Operation
		BIST_WARP_ENV["BENCHMARK_ARGS_WARP_DURATION"] = s.Duration
		BIST_WARP_ENV["BENCHMARK_ARGS_WARP_OBJ.SIZE"] = s.ObjectSize
		jobEnv = make([]corev1.EnvVar, 0, len(BIST_WARP_ENV))
		for k, v := range BIST_WARP_ENV {
			jobEnv = append(jobEnv, corev1.EnvVar{
				Name:  k,
				Value: v,
			})
		}
	case *BISTFIO:
		BIST_FIO_ENV["BENCHMARK_TYPE"] = target
		BIST_FIO_ENV["BENCHMARK_ARGS_FIO_FILESIZE"] = s.FileSize
		BIST_FIO_ENV["BENCHMARK_ARGS_FIO_IODEPTH"] = strconv.FormatUint(s.IODepth, 10)
		BIST_FIO_ENV["BENCHMARK_ARGS_FIO_NUMJOBS"] = strconv.FormatUint(s.JobCount, 10)
		BIST_FIO_ENV["BENCHMARK_ARGS_FIO_BS"] = s.BlockSize
		BIST_FIO_ENV["BENCHMARK_ARGS_FIO_RW"] = s.AccessMode
		BIST_FIO_ENV["BENCHMARK_ARGS_FIO_RUNTIME"] = s.RunTime
		BIST_FIO_ENV["BENCHMARK_ARGS_NFS_ENDPOINT"] = s.NFSEndpoint
		BIST_FIO_ENV["BENCHMARK_ARGS_NFS_PATH"] = s.NFSPath
		BIST_FIO_ENV["BENCHMARK_ARGS_NFS_SC"] = s.StorageClassName
		jobEnv = make([]corev1.EnvVar, 0, len(BIST_FIO_ENV))
		for k, v := range BIST_FIO_ENV {
			jobEnv = append(jobEnv, corev1.EnvVar{
				Name:  k,
				Value: v,
			})
		}
	default:
		jobEnv = make([]corev1.EnvVar, 0, len(BIST_WARP_ENV))
		for k, v := range BIST_WARP_ENV {
			jobEnv = append(jobEnv, corev1.EnvVar{
				Name:  k,
				Value: v,
			})
		}
	}

	return jobEnv, nil
}
