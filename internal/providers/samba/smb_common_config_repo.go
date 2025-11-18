package samba

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/otterscale/otterscale/internal/core/storage/smb"
)

type smbCommonConfigRepo struct {
	samba *Samba
}

func NewSMBCommonConfigRepo(samba *Samba) smb.SMBCommonConfigRepo {
	return &smbCommonConfigRepo{
		samba: samba,
	}
}

var _ smb.SMBCommonConfigRepo = (*smbCommonConfigRepo)(nil)

func (r *smbCommonConfigRepo) List(ctx context.Context, scope, namespace, selector string) ([]smb.CommonConfig, error) {
	clientset, err := r.samba.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.ListOptions{
		LabelSelector: selector,
	}

	list, err := clientset.ApiV1alpha1().SmbCommonConfigs(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func (r *smbCommonConfigRepo) Get(ctx context.Context, scope, namespace, name string) (*smb.CommonConfig, error) {
	clientset, err := r.samba.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.GetOptions{}

	return clientset.ApiV1alpha1().SmbCommonConfigs(namespace).Get(ctx, name, opts)
}

func (r *smbCommonConfigRepo) Create(ctx context.Context, scope, namespace string, cc *smb.CommonConfig) (*smb.CommonConfig, error) {
	clientset, err := r.samba.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.CreateOptions{}

	return clientset.ApiV1alpha1().SmbCommonConfigs(namespace).Create(ctx, cc, opts)
}

func (r *smbCommonConfigRepo) Update(ctx context.Context, scope, namespace string, cc *smb.CommonConfig) (*smb.CommonConfig, error) {
	clientset, err := r.samba.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.UpdateOptions{}

	return clientset.ApiV1alpha1().SmbCommonConfigs(namespace).Update(ctx, cc, opts)
}

func (r *smbCommonConfigRepo) Delete(ctx context.Context, scope, namespace, name string) error {
	clientset, err := r.samba.clientset(scope)
	if err != nil {
		return err
	}

	opts := metav1.DeleteOptions{}

	return clientset.ApiV1alpha1().SmbCommonConfigs(namespace).Delete(ctx, name, opts)
}
