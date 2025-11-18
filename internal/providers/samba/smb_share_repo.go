package samba

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/otterscale/otterscale/internal/core/storage/smb"
)

type smbShareRepo struct {
	samba *Samba
}

func NewSMBShareRepo(samba *Samba) smb.SMBShareRepo {
	return &smbShareRepo{
		samba: samba,
	}
}

var _ smb.SMBShareRepo = (*smbShareRepo)(nil)

func (r *smbShareRepo) List(ctx context.Context, scope, namespace, selector string) ([]smb.Share, error) {
	clientset, err := r.samba.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.ListOptions{
		LabelSelector: selector,
	}

	list, err := clientset.ApiV1alpha1().SmbShares(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func (r *smbShareRepo) Get(ctx context.Context, scope, namespace, name string) (*smb.Share, error) {
	clientset, err := r.samba.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.GetOptions{}

	return clientset.ApiV1alpha1().SmbShares(namespace).Get(ctx, name, opts)
}

func (r *smbShareRepo) Create(ctx context.Context, scope, namespace string, s *smb.Share) (*smb.Share, error) {
	clientset, err := r.samba.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.CreateOptions{}

	return clientset.ApiV1alpha1().SmbShares(namespace).Create(ctx, s, opts)
}

func (r *smbShareRepo) Update(ctx context.Context, scope, namespace string, s *smb.Share) (*smb.Share, error) {
	clientset, err := r.samba.clientset(scope)
	if err != nil {
		return nil, err
	}

	opts := metav1.UpdateOptions{}

	return clientset.ApiV1alpha1().SmbShares(namespace).Update(ctx, s, opts)
}

func (r *smbShareRepo) Delete(ctx context.Context, scope, namespace, name string) error {
	clientset, err := r.samba.clientset(scope)
	if err != nil {
		return err
	}

	opts := metav1.DeleteOptions{}

	return clientset.ApiV1alpha1().SmbShares(namespace).Delete(ctx, name, opts)
}
