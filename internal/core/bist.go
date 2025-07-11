package core

// import (
// 	"context"
// 	"sync"

// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// 	"github.com/openhdc/otterscale/internal/config"
// )

// const (
// 	BIST_NAME       = "bist"
// 	BIST_NAMESPACE  = "bist"
// 	BIST_LABEL      = "bist.otterscale.io/type"
// 	BIST_JOB_RETRY  = 2
// 	BIST_TYPE_BLOCK = "block"
// 	BIST_TYPE_NFS   = "nfs"
// 	BIST_TYPE_S3    = "s3"

// 	BIST_PHASE_COMPLETE = "COMPLETE"
// 	BIST_PHASE_RUNNING  = "RUNNING"
// 	BIST_PHASE_FAIL     = "FAILED"
// 	BIST_PHASE_CREATING = "CREATING"
// )

// var (
// 	BIST_IMAGE_LIST = map[string]string{
// 		"block": "docker.io/otterscale/bist-block:v1",
// 		"nfs":   "docker.io/otterscale/bist-nfs:v1",
// 		"s3":    "docker.io/otterscale/bist-s3:v3",
// 	}
// 	BIST_WARP_ENV = map[string]string{
// 		"BENCHMARK_TYPE":                 "s3",
// 		"BENCHMARK_ARGS_WARP_HOST":       "",
// 		"BENCHMARK_ARGS_WARP_ACCESS_KEY": "",
// 		"BENCHMARK_ARGS_WARP_SECRET_KEY": "",
// 		"BENCHMARK_ARGS_WARP_ACTION":     "",
// 		"BENCHMARK_ARGS_WARP_DURATION":   "",
// 		"BENCHMARK_ARGS_WARP_CONCURRENT": "2",
// 		"BENCHMARK_ARGS_WARP_OBJ.SIZE":   "",
// 		"BENCHMARK_ARGS_WARP_OBJECTS":    "",
// 	}

// 	BIST_FIO_ENV = map[string]string{
// 		"BENCHMARK_TYPE":                                "",
// 		"BENCHMARK_ARGS_BLOCK_POOL":                     "",
// 		"BENCHMARK_ARGS_BLOCK_IMAGE":                    "",
// 		"BENCHMARK_ARGS_FIO_DIRECT":                     "1",
// 		"BENCHMARK_ARGS_FIO_FILESIZE":                   "",
// 		"BENCHMARK_ARGS_FIO_IODEPTH":                    "",
// 		"BENCHMARK_ARGS_FIO_NUMJOBS":                    "",
// 		"BENCHMARK_ARGS_FIO_GROUP_REPORTING":            "True",
// 		"BENCHMARK_ARGS_FIO_NORANDOMMAP":                "True",
// 		"BENCHMARK_ARGS_FIO_BS":                         "",
// 		"BENCHMARK_ARGS_FIO_BUFFER_COMPRESS_PERCENTAGE": "0",
// 		"BENCHMARK_ARGS_FIO_RW":                         "",
// 		"BENCHMARK_ARGS_FIO_STARTDELAY":                 "5",
// 		"BENCHMARK_ARGS_FIO_TIME_BASED":                 "True",
// 		"BENCHMARK_ARGS_FIO_EXITALL_ON_ERROR":           "True",
// 		"BENCHMARK_ARGS_FIO_CREATE_SERIALIZE":           "0",
// 		"BENCHMARK_ARGS_FIO_RUNTIME":                    "",
// 	}
// )

// type LatNs struct {
// 	Min  int     `json:"min"`
// 	Max  int     `json:"max"`
// 	Mean float64 `json:"mean"`
// }

// type Operation struct {
// 	IoBytes  int     `json:"io_bytes"`
// 	Bw       int     `json:"bw"`
// 	Iops     float64 `json:"iops"`
// 	Runtime  int     `json:"runtime"`
// 	TotalIos int     `json:"total_ios"`
// 	LatNs    LatNs   `json:"lat_ns"`
// }

// type FIOJob struct {
// 	Read   Operation `json:"read"`
// 	Write  Operation `json:"write"`
// 	Trim   Operation `json:"trim"`
// 	UsrCPU float64   `json:"usr_cpu"`
// 	SysCPU float64   `json:"sys_cpu"`
// }

// type FIOResult struct {
// 	FioVersion string   `json:"fio version"`
// 	Timestamp  int      `json:"timestamp"`
// 	Time       string   `json:"time"`
// 	Jobs       []FIOJob `json:"jobs"`
// }

// type WarpResult struct {
// 	Type           string          `json:"type"`
// 	WarpOperations []WarpOperation `json:"operations"`
// }
// type BISTResult2 struct {
// 	Type         string
// 	Name         string
// 	Status       string
// 	StartTime    *metav1.Time
// 	CompleteTime *metav1.Time
// 	Args         string
// 	Logs         []string
// 	FIO          *BISTFIO
// 	Warp         *BISTWarp
// 	FIOResult    *FIOResult
// 	WarpResult   *WarpResult
// }

// // type BISTFIOTarget struct {
// // 	ScopeUUID    string
// // 	FacilityName string
// // 	NFS
// // }

// type BISTFIO struct {
// 	AccessMode       string
// 	StorageClassName string
// 	NFSEndpoint      string
// 	NFSPath          string
// 	JobCount         uint64
// 	RunTime          string
// 	BlockSize        string
// 	FileSize         string
// 	IODepth          uint64
// }

// type BISTWarp struct {
// 	Operation  string `json:"operation"`
// 	Endpoint   string `json:"endpoint"`
// 	AccessKey  string `json:"access_key"`
// 	SecretKey  string `json:"secret_key"`
// 	Duration   string `json:"duration"`
// 	ObjectSize string `json:"object_size.SIZE"`
// 	ObjectNum  string `json:"object_size.NUM"`
// }

// type BISTBlock struct {
// 	FacilityName     string
// 	StorageClassName string
// }

// type BISTObjectServiceList struct {
// 	Ceph   *BISTObjectService
// 	MinIOs []BISTObjectService
// }

// type BISTObjectService struct {
// 	Name     string
// 	Endpoint string
// }

// type BISTUseCase struct {
// 	action      ActionRepo
// 	scope       ScopeRepo
// 	kubeStorage KubeStorageRepo
// 	client      ClientRepo
// 	facility    FacilityRepo
// 	kubeCore    KubeCoreRepo
// 	kubeBatch   KubeBatchRepo
// 	conf        *config.Config
// 	configs     sync.Map
// }

// func NewBISTUseCase(action ActionRepo, scope ScopeRepo, kubeBatch KubeBatchRepo, kube KubeAppsRepo, kubeCore KubeCoreRepo, kubeStorage KubeStorageRepo, facility FacilityRepo, conf *config.Config, client ClientRepo) *BISTUseCase {
// 	return &BISTUseCase{
// 		action:      action,
// 		client:      client,
// 		scope:       scope,
// 		kubeCore:    kubeCore,
// 		kubeBatch:   kubeBatch,
// 		kubeStorage: kubeStorage,
// 		facility:    facility,
// 		conf:        conf,
// 	}
// }

// func (uc *BISTUseCase) ListResults(ctx context.Context) ([]BISTResult, error) {
// 	return nil, nil
// }

// func (uc *BISTUseCase) CreateResult(ctx context.Context, name, createdBy string, fio *BISTFIO, warp *BISTWarp) (*BISTResult, error) {
// 	return nil, nil
// }

// func (uc *BISTUseCase) DeleteResult(ctx context.Context, name string) error {
// 	return nil
// }

// func (uc *BISTUseCase) GetObjectServiceList(ctx context.Context, uuid string) (*BISTObjectServiceList, error) {
// 	return &BISTObjectServiceList{
// 		Ceph: &BISTObjectService{
// 			Name:     "aaa-ceph-mon",
// 			Endpoint: "1.2.3.4",
// 		},
// 		MinIOs: []BISTObjectService{
// 			{
// 				Name:     "minio-A",
// 				Endpoint: "1.2.3.4",
// 			},
// 			{
// 				Name:     "minio-B",
// 				Endpoint: "1.2.3.4",
// 			},
// 		},
// 	}, nil
// }

// // func (uc *BISTUseCase) ListResults(ctx context.Context) ([]BISTResult, error) {
// // 	config, err := uc.newConfig()
// // 	if err != nil {
// // 		return nil, err
// // 	}

// // 	jobs, err := uc.kubeBatch.ListJobsByLabel(ctx, config, BIST_NAMESPACE, BIST_LABEL)
// // 	if err != nil {
// // 		return nil, err
// // 	}

// // 	bists := []BISTResult{}
// // 	for _, job := range jobs {
// // 		bist, err := uc.toBISTResult(ctx, job.Name, &job, config)
// // 		if err != nil {
// // 			return nil, err
// // 		}
// // 		bists = append(bists, *bist)
// // 	}

// // 	return bists, nil
// // }

// // func (uc *BISTUseCase) CreateNamespace(ctx context.Context, config *rest.Config) error {

// // 	_, err := uc.kubeCore.GetNamespace(ctx, config, BIST_NAMESPACE)
// // 	if err != nil {
// // 		if kuberr.IsNotFound(err) {
// // 			namespace := &corev1.Namespace{
// // 				ObjectMeta: metav1.ObjectMeta{
// // 					Name: BIST_NAMESPACE,
// // 				},
// // 			}
// // 			_, err = uc.kubeCore.CreateNamespace(ctx, config, namespace)
// // 			if err != nil {
// // 				return err
// // 			}
// // 		} else {
// // 			return err
// // 		}
// // 	}
// // 	return nil
// // }

// // func (uc *BISTUseCase) CreateConfigMap(ctx context.Context, uuid, facilityName string, config *rest.Config) error {

// // 	cmName := fmt.Sprintf("ceph-conf-%s", generateShortName(uuid+facilityName))

// // 	_, err := uc.kubeCore.GetConfigMap(ctx, config, BIST_NAMESPACE, cmName)
// // 	if kuberr.IsNotFound(err) {

// // 		cephConfig, err := uc.newCephConfig(ctx, uuid, facilityName)

// // 		cm := uc.toBISTConfigMap(cmName, cephConfig.MONHost, cephConfig.FSID, cephConfig.Key)

// // 		cm, err = uc.kubeCore.CreateConfigMap(ctx, config, BIST_NAMESPACE, cm)
// // 		if err != nil {
// // 			return err
// // 		}
// // 	}

// // 	return nil
// // }

// // func (uc *BISTUseCase) toBISTConfigMap(name, host, fsid, key string) *corev1.ConfigMap {
// // 	return &corev1.ConfigMap{
// // 		ObjectMeta: metav1.ObjectMeta{
// // 			Name: name,
// // 		},
// // 		Data: map[string]string{
// // 			"ceph.conf": fmt.Sprintf(
// // 				`[global]
// // 				mon host = %s
// // 				fsid = %s
// // 				key = %s`, host, fsid, key),
// // 		},
// // 	}
// // }

// // func (uc *BISTUseCase) CreateResult(ctx context.Context, target, name, uuid, facility string, fio *BISTFIO, warp *BISTWarp) (*BISTResult, error) {
// // 	if target == "block" {
// // 		// [TODO] if exist
// // 		if !s.suc.IsPoolExists(ctx, uuid, facility, "otterscale_pool") {
// // 			_, err := s.suc.CreatePool(ctx, req.Msg.GetFio().GetScopeUuid(), req.Msg.GetFio().GetFacilityName(), "otterscale_pool", "replicated", false, 1, 0, 0, []string{"rbd"})
// // 			if err != nil {
// // 				return nil, err
// // 			}
// // 			_, err = s.suc.CreateImage(ctx, req.Msg.GetFio().GetScopeUuid(), req.Msg.GetFio().GetFacilityName(), "otterscale_pool", "otterscale_image", 4194304, 4194304, 1, 1073741824, true, true, true, true, true)
// // 			if err != nil {
// // 				return nil, err
// // 			}
// // 		}
// // 	}

// // 	config, err := uc.newConfig()
// // 	if err != nil {
// // 		return nil, err
// // 	}

// // 	err = uc.CreateNamespace(ctx, config)
// // 	if err != nil {
// // 		return nil, err
// // 	}

// // 	var job *batchv1.Job
// // 	switch target {
// // 	case "block":
// // 		err = uc.CreateConfigMap(ctx, uuid, facilityName, config)
// // 		if err != nil {
// // 			return nil, err
// // 		}

// // 		job, err = toBISTJob(name, BIST_NAMESPACE, target, BIST_JOB_RETRY, fio)
// // 	case "nfs":
// // 		job, err = toBISTJob(name, BIST_NAMESPACE, target, BIST_JOB_RETRY, fio)
// // 	case "s3":
// // 		job, err = toBISTJob(name, BIST_NAMESPACE, target, BIST_JOB_RETRY, warp)
// // 	}
// // 	if err != nil {
// // 		return nil, err
// // 	}

// // 	job, err = uc.kubeBatch.CreateJob(ctx, config, job)
// // 	if err != nil {
// // 		return nil, err
// // 	}

// // 	return uc.toBISTResult(ctx, job.Name, job, config)
// // }

// // func (uc *BISTUseCase) DeleteResult(ctx context.Context, name string) error {
// // 	config, err := uc.newConfig()
// // 	if err != nil {
// // 		return err
// // 	}
// // 	return uc.kubeBatch.DeleteJob(ctx, config, BIST_NAMESPACE, name)
// // }

// // func (uc *BISTUseCase) ListS3s(ctx context.Context, uuid string) ([]BISTS3, error) {
// // 	scopes, err := uc.scope.List(ctx)
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	scopes = slices.DeleteFunc(scopes, func(s Scope) bool {
// // 		return !strings.Contains(s.UUID, uuid)
// // 	})

// // 	minioEndpoints := make([]BISTS3, 0, len(scopes))
// // 	eg, ctx := errgroup.WithContext(ctx)
// // 	for i := range scopes {
// // 		eg.Go(func() error {
// // 			s, err := uc.client.Status(ctx, scopes[i].UUID, []string{"application", "*"})
// // 			if err != nil {
// // 				return err
// // 			}
// // 			for name := range s.Applications {
// // 				if !strings.Contains(s.Applications[name].Charm, "kubernetes-control-plane") {
// // 					continue
// // 				}
// // 				units := []EssentialUnit{}
// // 				for uname := range s.Applications[name].Units {
// // 					units = append(units, EssentialUnit{
// // 						Name:      uname,
// // 						Directive: s.Applications[name].Units[uname].Machine,
// // 					})
// // 				}

// // 				for _, unit := range units {
// // 					config, err := uc.config(ctx, uuid, removeLastSlashAndAfter(unit.Name))
// // 					if config == nil {
// // 						fmt.Println("nil config")
// // 						continue
// // 					}
// // 					label := "app.otterscale.io/release-name=minio"
// // 					field := "spec.type=NodePort"
// // 					svcs, err := uc.kubeCore.ListServicesByOptions(ctx, config, "", label, field)
// // 					if err != nil {
// // 						return err
// // 					}
// // 					for _, svc := range svcs {
// // 						for _, port := range svc.Spec.Ports {
// // 							if port.Name == "minio-api" {
// // 								minioEndpoints = append(minioEndpoints, toBISTS3("minio", svc.GetName()+"/"+svc.GetNamespace(), fmt.Sprintf("%s:%d", RemoveHTTPPrefix(config.Host), port.NodePort)))
// // 							}
// // 						}
// // 					}
// // 				}
// // 				break
// // 			}
// // 			return nil
// // 		})
// // 	}
// // 	if err := eg.Wait(); err != nil {
// // 		return nil, err
// // 	}

// // 	cephEndpoints := make([]BISTS3, 0, len(scopes))

// // 	for i := range scopes {
// // 		eg.Go(func() error {
// // 			s, err := uc.client.Status(ctx, scopes[i].UUID, []string{"application", "*"})
// // 			if err != nil {
// // 				return err
// // 			}
// // 			for name := range s.Applications {
// // 				if !strings.Contains(s.Applications[name].Charm, "ceph-mon") {
// // 					continue
// // 				}
// // 				units := []EssentialUnit{}
// // 				for uname := range s.Applications[name].Units {
// // 					units = append(units, EssentialUnit{
// // 						Name:      uname,
// // 						Directive: s.Applications[name].Units[uname].Machine,
// // 					})
// // 				}
// // 				for _, unit := range units {
// // 					info, err := uc.facility.GetUnitInfo(ctx, uuid, unit.Name)
// // 					if err != nil {
// // 						return err
// // 					}

// // 					cephEndpoints = append(cephEndpoints, toBISTS3("ceph", unit.Name, info.PublicAddress))
// // 				}
// // 				break
// // 			}
// // 			return nil
// // 		})
// // 	}
// // 	if err := eg.Wait(); err != nil {
// // 		return nil, err
// // 	}

// // 	s3Endpoints := append(cephEndpoints, minioEndpoints...)

// // 	return s3Endpoints, nil
// // }

// // func RemoveHTTPPrefix(input string) string {
// // 	if after, ok := strings.CutPrefix(input, "http://"); ok {
// // 		return after
// // 	} else if after, ok := strings.CutPrefix(input, "https://"); ok {
// // 		return after
// // 	}
// // 	return input
// // }

// // func (uc *BISTUseCase) newConfig() (*rest.Config, error) {
// // 	decodedBytes, err := base64.StdEncoding.DecodeString(uc.conf.MicroK8s.Token)
// // 	if err != nil {
// // 		return nil, err
// // 	}

// // 	return &rest.Config{
// // 		Host:        uc.conf.MicroK8s.Host,
// // 		BearerToken: string(decodedBytes),
// // 		TLSClientConfig: rest.TLSClientConfig{
// // 			Insecure: true,
// // 		},
// // 	}, nil
// // }

// // func (uc *BISTUseCase) config(ctx context.Context, uuid, name string) (*rest.Config, error) {
// // 	key := uuid + "/" + name

// // 	if v, ok := uc.configs.Load(key); ok {
// // 		return v.(*rest.Config), nil
// // 	}

// // 	config, err := uc.newKubeConfig(ctx, uuid, name)
// // 	if err != nil {
// // 		return nil, err
// // 	}

// // 	uc.configs.Store(key, config)

// // 	return config, nil
// // }

// // func (uc *BISTUseCase) newKubeConfig(ctx context.Context, uuid, name string) (*rest.Config, error) {
// // 	// kubernetes-control-plane
// // 	leader, err := uc.facility.GetLeader(ctx, uuid, name)
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	unitInfo, err := uc.facility.GetUnitInfo(ctx, uuid, leader)
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	kubeControl, err := extractWorkerUnitName(unitInfo)
// // 	if err != nil {
// // 		return nil, err
// // 	}

// // 	kubeControl = removeLastSlashAndAfter(kubeControl)

// // 	// kubernetes-worker
// // 	leader, err = uc.facility.GetLeader(ctx, uuid, kubeControl)
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	unitInfo, err = uc.facility.GetUnitInfo(ctx, uuid, leader)
// // 	if err != nil {
// // 		return nil, err
// // 	}

// // 	// config
// // 	endpoint, err := extractEndpoint(unitInfo)
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	clientToken, err := extractClientToken(unitInfo)
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	return &rest.Config{
// // 		Host:        endpoint,
// // 		BearerToken: clientToken,
// // 		TLSClientConfig: rest.TLSClientConfig{
// // 			Insecure: true,
// // 		},
// // 	}, nil
// // }

// // func (uc *BISTUseCase) newCephConfig(ctx context.Context, uuid, name string) (*StorageCephConfig, error) {
// // 	var (
// // 		leader string
// // 	)
// // 	eg, egctx := errgroup.WithContext(ctx)
// // 	eg.Go(func() error {
// // 		monLeader, err := uc.facility.GetLeader(egctx, uuid, name) // ceph-mon
// // 		if err != nil {
// // 			return err
// // 		}
// // 		leader = monLeader
// // 		return nil
// // 	})

// // 	if err := eg.Wait(); err != nil {
// // 		return nil, err
// // 	}

// // 	var cephConfig *StorageCephConfig

// // 	eg, egctx = errgroup.WithContext(ctx)
// // 	eg.Go(func() error {
// // 		result, err := uc.runAction(egctx, uuid, leader, cephConfigCommand)
// // 		if err != nil {
// // 			return err
// // 		}
// // 		config, err := uc.extractStorageCephConfig(result)
// // 		if err != nil {
// // 			return err
// // 		}
// // 		cephConfig = config
// // 		return nil
// // 	})

// // 	if err := eg.Wait(); err != nil {
// // 		return nil, err
// // 	}

// // 	return cephConfig, nil
// // }

// // func (uc *BISTUseCase) runAction(ctx context.Context, uuid, leader, command string) (*action.ActionResult, error) {
// // 	id, err := uc.action.RunCommand(ctx, uuid, leader, command)
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	return uc.waitForActionCompleted(ctx, uuid, id)
// // }

// // func (uc *BISTUseCase) waitForActionCompleted(ctx context.Context, uuid, id string) (*action.ActionResult, error) {
// // 	const tickInterval = time.Second
// // 	const timeoutDuration = 5 * time.Second

// // 	ticker := time.NewTicker(tickInterval)
// // 	defer ticker.Stop()

// // 	timeout := time.After(timeoutDuration)
// // 	for {
// // 		select {
// // 		case <-ticker.C:
// // 			result, err := uc.action.GetResult(ctx, uuid, id)
// // 			if err != nil {
// // 				return nil, err
// // 			}
// // 			if result.Status == "completed" { // state.ActionCompleted
// // 				return result, nil
// // 			}
// // 			continue

// // 		case <-timeout:
// // 			return nil, fmt.Errorf("timeout waiting for action %s to become completed", id)

// // 		case <-ctx.Done():
// // 			return nil, ctx.Err()
// // 		}
// // 	}
// // }

// // func (uc *BISTUseCase) extractStorageCephConfig(result *action.ActionResult) (*StorageCephConfig, error) {
// // 	stdout, ok := result.Output["stdout"]
// // 	if !ok {
// // 		return nil, errors.New("ceph config stdout not found")
// // 	}
// // 	file, err := ini.Load([]byte(stdout.(string)))
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	fsID := file.Section("global").Key("fsid").String()
// // 	if fsID == "" {
// // 		return nil, errors.New("ceph config fsid not found")
// // 	}
// // 	monHost := file.Section("global").Key("mon_host").String()
// // 	if monHost == "" {
// // 		return nil, errors.New("ceph config mon_host not found")
// // 	}
// // 	key := file.Section("client.admin").Key("key").String()
// // 	if key == "" {
// // 		return nil, errors.New("ceph config key not found")
// // 	}
// // 	return &StorageCephConfig{
// // 		FSID:    fsID,
// // 		MONHost: monHost,
// // 		Key:     key,
// // 	}, nil
// // }

// // func removeLastSlashAndAfter(s string) string {
// // 	lastSlashIndex := strings.LastIndex(s, "/")

// // 	if lastSlashIndex == -1 {
// // 		return s
// // 	}
// // 	return s[:lastSlashIndex]
// // }

// // func toBISTS3(s3type, name, endpoint string) BISTS3 {
// // 	return BISTS3{
// // 		Type:     s3type,
// // 		Name:     name,
// // 		Endpoint: endpoint,
// // 	}
// // }

// // func toBISTJob(name, namespace, target string, retry int32, BIST interface{}) (*batchv1.Job, error) {
// // 	trueVal := false
// // 	image := BIST_IMAGE_LIST[target]

// // 	if name == "" {
// // 		name = generateNameWithDateTime(BIST_NAME)
// // 	}

// // 	if target == BIST_TYPE_BLOCK || target == BIST_TYPE_NFS {
// // 		trueVal = true
// // 	}

// // 	envJSON, err := json.Marshal(BIST)
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	annotation := make(map[string]string)
// // 	annotation["bist.otterscale.io/config"] = string(envJSON)

// // 	BIST_JOB_LABEL := map[string]string{
// // 		//"bist.otterscale.io/create-by": user,
// // 		"bist.otterscale.io/type": target,
// // 	}

// // 	jobEnv, err := toJobEnv(target, BIST)
// // 	if err != nil {
// // 		return nil, err
// // 	}

// // 	return &batchv1.Job{
// // 		ObjectMeta: metav1.ObjectMeta{
// // 			Name:        name,
// // 			Namespace:   namespace,
// // 			Labels:      BIST_JOB_LABEL,
// // 			Annotations: annotation,
// // 		},
// // 		Spec: batchv1.JobSpec{
// // 			BackoffLimit: &retry,
// // 			Template: corev1.PodTemplateSpec{
// // 				Spec: corev1.PodSpec{
// // 					Containers: []corev1.Container{
// // 						{
// // 							Name:            "bist-container",
// // 							Image:           image,
// // 							ImagePullPolicy: corev1.PullIfNotPresent,
// // 							Command: []string{
// // 								"./start.sh",
// // 							},
// // 							SecurityContext: &corev1.SecurityContext{
// // 								Privileged: &trueVal,
// // 							},
// // 							Env:          jobEnv,
// // 							VolumeMounts: generateVolumeMounts(target),
// // 						},
// // 					},
// // 					Volumes:       generateVolume(target),
// // 					RestartPolicy: corev1.RestartPolicyNever,
// // 				},
// // 			},
// // 		},
// // 	}, nil
// // }

// // func generateVolume(target string) []corev1.Volume {
// // 	if target == BIST_TYPE_NFS {
// // 		return []corev1.Volume{
// // 			{
// // 				Name: "ceph-conf",
// // 				VolumeSource: corev1.VolumeSource{
// // 					ConfigMap: &corev1.ConfigMapVolumeSource{
// // 						LocalObjectReference: corev1.LocalObjectReference{
// // 							Name: "ceph-conf",
// // 						},
// // 						Items: []corev1.KeyToPath{
// // 							{
// // 								Key:  "ceph.conf",
// // 								Path: "ceph.conf",
// // 							},
// // 						},
// // 					},
// // 				},
// // 			},
// // 			{
// // 				Name: "dev",
// // 				VolumeSource: corev1.VolumeSource{
// // 					HostPath: &corev1.HostPathVolumeSource{
// // 						Path: "/dev",
// // 					},
// // 				},
// // 			},
// // 			{
// // 				Name: "modules",
// // 				VolumeSource: corev1.VolumeSource{
// // 					HostPath: &corev1.HostPathVolumeSource{
// // 						Path: "/lib/modules",
// // 					},
// // 				},
// // 			},
// // 			{
// // 				Name: "run-udev",
// // 				VolumeSource: corev1.VolumeSource{
// // 					HostPath: &corev1.HostPathVolumeSource{
// // 						Path: "/run/udev",
// // 					},
// // 				},
// // 			},
// // 		}
// // 	}

// // 	return []corev1.Volume{}
// // }

// // func generateVolumeMounts(target string) []corev1.VolumeMount {
// // 	if target == BIST_TYPE_NFS {
// // 		return []corev1.VolumeMount{
// // 			{
// // 				Name:      "ceph-conf",
// // 				MountPath: "/etc/ceph/ceph.conf",
// // 				SubPath:   "ceph.conf",
// // 			},
// // 			{
// // 				Name:      "dev",
// // 				MountPath: "/dev",
// // 			},
// // 			{
// // 				Name:      "modules",
// // 				MountPath: "/lib/modules",
// // 			},
// // 			{
// // 				Name:      "run-udev",
// // 				MountPath: "/run/udev",
// // 			},
// // 		}
// // 	}

// // 	return []corev1.VolumeMount{}
// // }

// // func generateNameWithDateTime(baseName string) string {
// // 	now := time.Now()
// // 	dateTimeStr := now.Format("20060102150405")
// // 	return fmt.Sprintf("%s-%s", baseName, dateTimeStr)
// // }

// // func generateShortName(baseName string) string {
// // 	hash := sha256.Sum256([]byte(baseName))
// // 	hashString := hex.EncodeToString(hash[:])
// // 	shortString := hashString[:8]

// // 	return shortString
// // }

// // // Warp result has redundant message
// // func RemoveLastTwoLines(input string) string {
// // 	lines := strings.Split(input, "\n")
// // 	if len(lines) < 3 {
// // 		return input
// // 	}
// // 	lines = lines[:len(lines)-2]
// // 	return strings.Join(lines, "\n")
// // }

// // func (uc *BISTUseCase) toBISTResult(ctx context.Context, name string, job *batchv1.Job, config *rest.Config) (*BISTResult, error) {

// // 	var jobStatus string
// // 	var ret BISTResult

// // 	if job.Status.Succeeded == 1 && job.Status.CompletionTime != nil {
// // 		jobStatus = BIST_PHASE_COMPLETE
// // 	} else if job.Status.Failed == 1 {
// // 		jobStatus = BIST_PHASE_FAIL
// // 	} else if job.Status.Active == 1 && job.Status.Ready != nil {
// // 		jobStatus = BIST_PHASE_RUNNING
// // 	} else {
// // 		jobStatus = BIST_PHASE_CREATING
// // 	}

// // 	selector, err := metav1.LabelSelectorAsSelector(job.Spec.Selector)
// // 	if err != nil {
// // 		return nil, err
// // 	}

// // 	pods, err := uc.kubeCore.ListPodsByLabel(ctx, config, BIST_NAMESPACE, selector.String())
// // 	if err != nil {
// // 		return nil, err
// // 	}

// // 	podLog := make([]string, 0, len(pods))
// // 	for _, pod := range pods {
// // 		//[TODO] go routine
// // 		if job.Status.CompletionTime != nil {
// // 			logs, err := uc.kubeCore.GetPodLogs(ctx, pod, config, BIST_NAMESPACE)
// // 			if err != nil {
// // 				return nil, err
// // 			}
// // 			podLog = append(podLog, logs)
// // 			json.Unmarshal([]byte(logs), &ret.FIOResult)
// // 			json.Unmarshal([]byte(RemoveLastTwoLines(logs)), &ret.WarpResult)
// // 			break
// // 		}
// // 	}

// // 	ret.Name = name
// // 	ret.Type = job.GetLabels()["bist.otterscale.io/type"]
// // 	ret.Status = jobStatus
// // 	ret.StartTime = job.Status.StartTime
// // 	ret.CompleteTime = job.Status.CompletionTime
// // 	ret.Logs = podLog

// // 	json.Unmarshal([]byte(job.GetAnnotations()["bist.otterscale.io/config"]), &ret.Warp)
// // 	json.Unmarshal([]byte(job.GetAnnotations()["bist.otterscale.io/config"]), &ret.FIO)

// // 	return &ret, nil
// // }

// // func toJobEnv(target string, BIST interface{}) ([]corev1.EnvVar, error) {

// // 	var jobEnv []corev1.EnvVar
// // 	v := reflect.ValueOf(BIST)
// // 	t := reflect.TypeOf(BIST)

// // 	if v.Kind() == reflect.Ptr {
// // 		v = v.Elem()
// // 		t = t.Elem()
// // 	}

// // 	switch s := BIST.(type) {
// // 	case *BISTWarp:
// // 		BIST_WARP_ENV["BENCHMARK_ARGS_WARP_HOST"] = s.Endpoint
// // 		BIST_WARP_ENV["BENCHMARK_ARGS_WARP_ACCESS_KEY"] = s.AccessKey
// // 		BIST_WARP_ENV["BENCHMARK_ARGS_WARP_SECRET_KEY"] = s.SecretKey
// // 		BIST_WARP_ENV["BENCHMARK_ARGS_WARP_ACTION"] = s.Operation
// // 		BIST_WARP_ENV["BENCHMARK_ARGS_WARP_DURATION"] = s.Duration
// // 		BIST_WARP_ENV["BENCHMARK_ARGS_WARP_OBJ.SIZE"] = s.ObjectSize
// // 		BIST_WARP_ENV["BENCHMARK_ARGS_WARP_OBJECTS"] = s.ObjectNum
// // 		jobEnv = make([]corev1.EnvVar, 0, len(BIST_WARP_ENV))
// // 		for k, v := range BIST_WARP_ENV {
// // 			jobEnv = append(jobEnv, corev1.EnvVar{
// // 				Name:  k,
// // 				Value: v,
// // 			})
// // 		}
// // 	case *BISTFIO:
// // 		BIST_FIO_ENV["BENCHMARK_TYPE"] = target
// // 		if target == BIST_TYPE_BLOCK {
// // 			BIST_FIO_ENV["BENCHMARK_ARGS_BLOCK_POOL"] = "otterscale_pool"
// // 			BIST_FIO_ENV["BENCHMARK_ARGS_BLOCK_IMAGE"] = "otterscale_image"
// // 		}
// // 		BIST_FIO_ENV["BENCHMARK_ARGS_FIO_FILESIZE"] = s.FileSize
// // 		BIST_FIO_ENV["BENCHMARK_ARGS_FIO_IODEPTH"] = strconv.FormatUint(s.IODepth, 10)
// // 		BIST_FIO_ENV["BENCHMARK_ARGS_FIO_NUMJOBS"] = strconv.FormatUint(s.JobCount, 10)
// // 		BIST_FIO_ENV["BENCHMARK_ARGS_FIO_BS"] = s.BlockSize
// // 		BIST_FIO_ENV["BENCHMARK_ARGS_FIO_RW"] = s.AccessMode
// // 		BIST_FIO_ENV["BENCHMARK_ARGS_FIO_RUNTIME"] = s.RunTime
// // 		BIST_FIO_ENV["BENCHMARK_ARGS_NFS_ENDPOINT"] = s.NFSEndpoint
// // 		BIST_FIO_ENV["BENCHMARK_ARGS_NFS_PATH"] = s.NFSPath
// // 		jobEnv = make([]corev1.EnvVar, 0, len(BIST_FIO_ENV))
// // 		for k, v := range BIST_FIO_ENV {
// // 			jobEnv = append(jobEnv, corev1.EnvVar{
// // 				Name:  k,
// // 				Value: v,
// // 			})
// // 		}
// // 	default:
// // 		return nil, errors.New("unsupported benchmark type")
// // 	}

// // 	return jobEnv, nil
// // }
