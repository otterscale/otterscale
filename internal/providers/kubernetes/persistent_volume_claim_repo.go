package kubernetes

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/otterscale/otterscale/internal/core/application/persistent"
)

type persistentVolumeClaimRepo struct {
	kubernetes *Kubernetes
}

func NewPersistentVolumeClaimRepo(kubernetes *Kubernetes) persistent.PersistentVolumeClaimRepo {
	return &persistentVolumeClaimRepo{
		kubernetes: kubernetes,
	}
}

var _ persistent.PersistentVolumeClaimRepo = (*persistentVolumeClaimRepo)(nil)

func (r *persistentVolumeClaimRepo) List(ctx context.Context, scope, namespace, selector string) ([]persistent.PersistentVolumeClaim, error) {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.ListOptions{
		LabelSelector: selector,
	}

	list, err := clientset.CoreV1().PersistentVolumeClaims(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func (r *persistentVolumeClaimRepo) Get(ctx context.Context, scope, namespace, name string) (*persistent.PersistentVolumeClaim, error) {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.GetOptions{}

	return clientset.CoreV1().PersistentVolumeClaims(namespace).Get(ctx, name, opts)
}

func (r *persistentVolumeClaimRepo) Create(ctx context.Context, scope, namespace string, pvc *persistent.PersistentVolumeClaim) (*persistent.PersistentVolumeClaim, error) {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.CreateOptions{}

	return clientset.CoreV1().PersistentVolumeClaims(namespace).Create(ctx, pvc, opts)
}

func (r *persistentVolumeClaimRepo) Update(ctx context.Context, scope, namespace string, pvc *persistent.PersistentVolumeClaim) (*persistent.PersistentVolumeClaim, error) {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.UpdateOptions{}

	return clientset.CoreV1().PersistentVolumeClaims(namespace).Update(ctx, pvc, opts)
}

func (r *persistentVolumeClaimRepo) Patch(ctx context.Context, scope, namespace, name string, data []byte) (*persistent.PersistentVolumeClaim, error) {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.PatchOptions{}

	return clientset.CoreV1().PersistentVolumeClaims(namespace).Patch(ctx, name, types.JSONPatchType, data, opts)
}

func (r *persistentVolumeClaimRepo) Delete(ctx context.Context, scope, namespace, name string) error {
	clientset, err := r.kubernetes.clientset(scope)
	if err != nil {
		return err
	}

	opts := metav1.DeleteOptions{}

	return clientset.CoreV1().PersistentVolumeClaims(namespace).Delete(ctx, name, opts)
}
