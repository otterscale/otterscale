package cdi

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/otterscale/otterscale/internal/core/application/persistent"
	"golang.org/x/sync/errgroup"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	cdiv1beta1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
)

const DataVolumeBootImageLabel = "otterscale.com/data-volume.boot-image"

type DataVolumeSourceType int64

const (
	DataVolumeSourceTypeBlank DataVolumeSourceType = iota
	DataVolumeSourceTypeHTTP
	DataVolumeSourceTypePVC
)

// DataVolume represents a KubeVirt DataVolume resource.
type DataVolume = cdiv1beta1.DataVolume

type DataVolumePersistent struct {
	*DataVolume
	*persistent.Persistent
}

type DataVolumeRepo interface {
	List(ctx context.Context, scope, namespace, selector string) ([]DataVolume, error)
	Get(ctx context.Context, scope, namespace, name string) (*DataVolume, error)
	Create(ctx context.Context, scope, namespace string, dv *DataVolume) (*DataVolume, error)
	Delete(ctx context.Context, scope, namespace, name string) error
}

type DataVolumeUseCase struct {
	dataVolume DataVolumeRepo

	persistentVolumeClaim persistent.PersistentVolumeClaimRepo
	storageClass          persistent.StorageClassRepo
}

func NewDataVolumeUseCase(dataVolume DataVolumeRepo, persistentVolumeClaim persistent.PersistentVolumeClaimRepo) *DataVolumeUseCase {
	return &DataVolumeUseCase{
		dataVolume:            dataVolume,
		persistentVolumeClaim: persistentVolumeClaim,
	}
}

func (uc *DataVolumeUseCase) ListDataVolumes(ctx context.Context, scope, namespace string, bootImage bool) ([]DataVolumePersistent, error) {
	var (
		dataVolumes            []DataVolume
		persistentVolumeClaims []persistent.PersistentVolumeClaim
		storageClasses         []persistent.StorageClass
	)

	eg, egctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		selector := DataVolumeBootImageLabel + "=" + strconv.FormatBool(bootImage)

		v, err := uc.dataVolume.List(ctx, scope, namespace, selector)
		if err == nil {
			dataVolumes = v
		}
		return err
	})

	eg.Go(func() error {
		v, err := uc.persistentVolumeClaim.List(egctx, scope, namespace, "")
		if err == nil {
			persistentVolumeClaims = v
		}
		return err
	})

	eg.Go(func() error {
		v, err := uc.storageClass.List(egctx, scope, "")
		if err == nil {
			storageClasses = v
		}
		return err
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	storageClassMap := map[string]*persistent.StorageClass{}
	for i := range storageClasses {
		sc := storageClasses[i]
		storageClassMap[sc.Name] = &sc
	}

	pvcMap := map[string]*persistent.PersistentVolumeClaim{}
	for i := range persistentVolumeClaims {
		pvc := persistentVolumeClaims[i]
		if pvc.Namespace == namespace {
			pvcMap[pvc.Name] = &pvc
		}
	}

	ret := make([]DataVolumePersistent, len(dataVolumes))

	for i := range ret {
		dv := dataVolumes[i]

		pvc, found := pvcMap[dv.Name]
		if !found {
			continue
		}

		dvStorage := DataVolumePersistent{
			DataVolume: &dv,
		}

		scName := pvc.Spec.StorageClassName
		if scName != nil && *scName != "" {
			sc, found := storageClassMap[*scName]
			if found {
				dvStorage.StorageClass = sc
			}
		}

		ret[i] = dvStorage
	}

	return ret, nil
}

func (uc *DataVolumeUseCase) GetDataVolume(ctx context.Context, scope, namespace, name string) (*DataVolumePersistent, error) {
	var (
		dataVolume            *DataVolume
		persistentVolumeClaim *persistent.PersistentVolumeClaim
		storageClass          *persistent.StorageClass
	)

	eg, egctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		v, err := uc.dataVolume.Get(egctx, scope, namespace, name)
		if err == nil {
			dataVolume = v
		}
		return err
	})

	eg.Go(func() error {
		v, err := uc.persistentVolumeClaim.Get(egctx, scope, namespace, name)
		if err == nil {
			scName := persistentVolumeClaim.Spec.StorageClassName
			if scName != nil {
				sc, err := uc.storageClass.Get(egctx, scope, *scName)
				if err == nil {
					persistentVolumeClaim = v
					storageClass = sc
				}
				return err
			}
		}
		return err
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return &DataVolumePersistent{
		DataVolume: dataVolume,
		Persistent: &persistent.Persistent{
			PersistentVolumeClaim: persistentVolumeClaim,
			StorageClass:          storageClass,
		},
	}, nil
}

func (uc *DataVolumeUseCase) CreateDataVolume(ctx context.Context, scope, namespace, name string, srcType DataVolumeSourceType, srcData string, size int64, bootImage bool) (*DataVolumePersistent, error) {
	dataVolume, err := uc.dataVolume.Create(ctx, scope, namespace, uc.buildDataVolume(namespace, name, srcType, srcData, size, bootImage))
	if err != nil {
		return nil, err
	}

	return &DataVolumePersistent{
		DataVolume: dataVolume,
	}, nil
}

func (uc *DataVolumeUseCase) DeleteDataVolume(ctx context.Context, scope, namespace, name string) error {
	return uc.dataVolume.Delete(ctx, scope, namespace, name)
}

func (uc *DataVolumeUseCase) ExtendDataVolume(ctx context.Context, scope, namespace, name string, newSize int64) error {
	data, err := json.Marshal([]map[string]any{
		{
			"op":    "replace",
			"path":  "/spec/resources/requests/storage",
			"value": resource.NewQuantity(newSize, resource.BinarySI).String(),
		},
	})
	if err != nil {
		return err
	}

	_, err = uc.persistentVolumeClaim.Patch(ctx, scope, namespace, name, data)
	return err
}

func (uc *DataVolumeUseCase) buildDataVolume(namespace, name string, srcType DataVolumeSourceType, srcData string, size int64, bootImage bool) *DataVolume {
	var (
		source  *cdiv1beta1.DataVolumeSource
		storage *cdiv1beta1.StorageSpec
		pvc     *corev1.PersistentVolumeClaimSpec
	)

	switch srcType {
	case DataVolumeSourceTypeBlank:
		source = &cdiv1beta1.DataVolumeSource{
			Blank: &cdiv1beta1.DataVolumeBlankImage{},
		}

		storage = &cdiv1beta1.StorageSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: *resource.NewQuantity(size, resource.BinarySI),
				},
			},
		}

	case DataVolumeSourceTypeHTTP:
		source = &cdiv1beta1.DataVolumeSource{
			HTTP: &cdiv1beta1.DataVolumeSourceHTTP{URL: srcData},
		}

		storage = &cdiv1beta1.StorageSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: *resource.NewQuantity(size, resource.BinarySI),
				},
			},
		}

	case DataVolumeSourceTypePVC:
		source = &cdiv1beta1.DataVolumeSource{
			PVC: &cdiv1beta1.DataVolumeSourcePVC{
				Namespace: namespace,
				Name:      srcData,
			},
		}

		pvc = &corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: *resource.NewQuantity(size, resource.BinarySI),
				},
			},
		}
	}

	return &cdiv1beta1.DataVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				DataVolumeBootImageLabel: strconv.FormatBool(bootImage),
			},
		},
		Spec: cdiv1beta1.DataVolumeSpec{
			Source:  source,
			Storage: storage,
			PVC:     pvc,
		},
	}
}
