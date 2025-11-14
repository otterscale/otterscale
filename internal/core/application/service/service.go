package service

import (
	"context"

	v1 "k8s.io/api/core/v1"
)

type (
	// Service represents a Kubernetes Service resource.
	Service = v1.Service

	// Port represents a Kubernetes ServicePort resource.
	Port = v1.ServicePort
)

type ServiceRepo interface {
	List(ctx context.Context, scope, namespace, selector string) ([]Service, error)
	Get(ctx context.Context, scope, namespace, name string) (*Service, error)
	Update(ctx context.Context, scope, namespace string, s *Service) (*Service, error)
	Create(ctx context.Context, scope, namespace string, s *Service) (*Service, error)
	Delete(ctx context.Context, scope, namespace, name string) error
	Host(scope string) string
}

type ServiceUseCase struct {
	service ServiceRepo
}

func NewServiceUseCase(service ServiceRepo) *ServiceUseCase {
	return &ServiceUseCase{
		service: service,
	}
}
