package core

import (
	"context"
)

type BISTResult struct {
	Type string
	Name string
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
	Operation  string
	Endpoint   string
	AccessKey  string
	SecretKey  string
	Duration   string
	ObjectSize string
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
	kubeCore KubeCoreRepo
}

func NewBISTUseCase(kubeCore KubeCoreRepo) *BISTUseCase {
	return &BISTUseCase{
		kubeCore: kubeCore,
	}
}

func (uc *BISTUseCase) ListResults(ctx context.Context) ([]BISTResult, error) {
	return nil, nil
}

func (uc *BISTUseCase) CreateResult(ctx context.Context, target, name string, fio *BISTFIO, warp *BISTWarp) (*BISTResult, error) {
	return nil, nil
}

func (uc *BISTUseCase) DeleteResult(ctx context.Context, name string) error {
	return nil
}

func (uc *BISTUseCase) ListBlocks(ctx context.Context) ([]BISTBlock, error) {
	return nil, nil
}

func (uc *BISTUseCase) ListS3s(ctx context.Context, uuid string) ([]BISTS3, error) {
	return nil, nil
}
