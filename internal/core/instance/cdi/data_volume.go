package cdi

import (
	"context"
	"encoding/json"

	"golang.org/x/sync/errgroup"
	"k8s.io/apimachinery/pkg/api/resource"
	cdiv1beta1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"

	"github.com/otterscale/otterscale/internal/core/application/storage"
)

type DataVolumeSourceType int64

const (
	DataVolumeSourceTypeBlank DataVolumeSourceType = iota
	DataVolumeSourceTypeHTTP
	DataVolumeSourceTypePVC
)

// DataVolume represents a KubeVirt DataVolume resource.
type DataVolume = cdiv1beta1.DataVolume

type DataVolumeStorage struct {
	*DataVolume
	*storage.Storage
}

type DataVolumeRepo interface {
	List(ctx context.Context, scope, namespace string, bootImage bool) ([]DataVolume, error)
	Get(ctx context.Context, scope, namespace, name string) (*DataVolume, error)
	Create(ctx context.Context, scope, namespace, name string, srcType DataVolumeSourceType, srcData string, size int64, bootImage bool) (*DataVolume, error)
	Delete(ctx context.Context, scope, namespace, name string) error
}

type DataVolumeUseCase struct {
	dataVolume DataVolumeRepo

	persistentVolumeClaim storage.PersistentVolumeClaimRepo
	storageClass          storage.StorageClassRepo
}

func NewDataVolumeUseCase(dataVolume DataVolumeRepo, persistentVolumeClaim storage.PersistentVolumeClaimRepo) *DataVolumeUseCase {
	return &DataVolumeUseCase{
		dataVolume:            dataVolume,
		persistentVolumeClaim: persistentVolumeClaim,
	}
}

func (uc *DataVolumeUseCase) ListDataVolumes(ctx context.Context, scope, namespace string, bootImage bool) ([]DataVolumeStorage, error) {
	var (
		dataVolumes            []DataVolume
		persistentVolumeClaims []storage.PersistentVolumeClaim
		storageClasses         []storage.StorageClass
	)

	eg, egctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		v, err := uc.dataVolume.List(ctx, scope, namespace, bootImage)
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

	storageClassMap := make(map[string]*storage.StorageClass)
	for i := range storageClasses {
		sc := storageClasses[i]
		storageClassMap[sc.Name] = &sc
	}

	pvcMap := make(map[string]*storage.PersistentVolumeClaim)
	for i := range persistentVolumeClaims {
		pvc := persistentVolumeClaims[i]
		if pvc.Namespace == namespace {
			pvcMap[pvc.Name] = &pvc
		}
	}

	ret := make([]DataVolumeStorage, len(dataVolumes))

	for i := range ret {
		dv := dataVolumes[i]

		pvc, found := pvcMap[dv.Name]
		if !found {
			continue
		}

		dvStorage := DataVolumeStorage{
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

func (uc *DataVolumeUseCase) GetDataVolume(ctx context.Context, scope, namespace, name string) (*DataVolumeStorage, error) {
	var (
		dataVolume            *DataVolume
		persistentVolumeClaim *storage.PersistentVolumeClaim
		storageClass          *storage.StorageClass
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

	return &DataVolumeStorage{
		DataVolume: dataVolume,
		Storage: &storage.Storage{
			PersistentVolumeClaim: persistentVolumeClaim,
			StorageClass:          storageClass,
		},
	}, nil
}

func (uc *DataVolumeUseCase) CreateDataVolume(ctx context.Context, scope, namespace, name string, srcType DataVolumeSourceType, srcData string, size int64, bootImage bool) (*DataVolumeStorage, error) {
	dataVolume, err := uc.dataVolume.Create(ctx, scope, namespace, name, srcType, srcData, size, bootImage)
	if err != nil {
		return nil, err
	}

	return &DataVolumeStorage{
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
