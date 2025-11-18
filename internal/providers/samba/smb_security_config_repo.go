package samba

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/otterscale/otterscale/internal/core/storage/smb"
)

type smbSecurityConfigRepo struct {
	samba *Samba
}

func NewSMBSecurityConfigRepo(samba *Samba) smb.SMBSecurityConfigRepo {
	return &smbSecurityConfigRepo{
		samba: samba,
	}
}

var _ smb.SMBSecurityConfigRepo = (*smbSecurityConfigRepo)(nil)

func (r *smbSecurityConfigRepo) List(ctx context.Context, scope, namespace, selector string) ([]smb.SecurityConfig, error) {
	clientset, err := r.samba.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.ListOptions{
		LabelSelector: selector,
	}

	list, err := clientset.ApiV1alpha1().SmbSecurityConfigs(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func (r *smbSecurityConfigRepo) Get(ctx context.Context, scope, namespace, name string) (*smb.SecurityConfig, error) {
	clientset, err := r.samba.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.GetOptions{}

	return clientset.ApiV1alpha1().SmbSecurityConfigs(namespace).Get(ctx, name, opts)
}

func (r *smbSecurityConfigRepo) Create(ctx context.Context, scope, namespace string, sc *smb.SecurityConfig) (*smb.SecurityConfig, error) {
	clientset, err := r.samba.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.CreateOptions{}

	return clientset.ApiV1alpha1().SmbSecurityConfigs(namespace).Create(ctx, sc, opts)
}

func (r *smbSecurityConfigRepo) Update(ctx context.Context, scope, namespace string, sc *smb.SecurityConfig) (*smb.SecurityConfig, error) {
	clientset, err := r.samba.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.UpdateOptions{}

	return clientset.ApiV1alpha1().SmbSecurityConfigs(namespace).Update(ctx, sc, opts)
}

func (r *smbSecurityConfigRepo) Delete(ctx context.Context, scope, namespace, name string) error {
	clientset, err := r.samba.clientset(scope)
	if err != nil {
		return err
	}

	opts := metav1.DeleteOptions{}

	return clientset.ApiV1alpha1().SmbSecurityConfigs(namespace).Delete(ctx, name, opts)
}
