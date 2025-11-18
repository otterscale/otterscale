package smb

import (
	"context"

	"github.com/samba-in-kubernetes/samba-operator/api/v1alpha1"
)

const (
	GuestOKkey    = "guest ok"
	MapToGuestKey = "map to guest"
	ValidUsersKey = "valid users"
)

type (
	// CommonConfig represents a Samba-Operator SmbCommonConfig resource.
	CommonConfig = v1alpha1.SmbCommonConfig

	// SecurityConfig represents a Samba-Operator SecurityConfig resource.
	SecurityConfig = v1alpha1.SmbSecurityConfig

	// Share represents a Samba-Operator SmbShare resource.
	Share = v1alpha1.SmbShare

	// UserSpec represents a Samba-Operator SmbSecurityUsersSpec resource.
	UserSpec = v1alpha1.SmbSecurityUsersSpec

	// JoinSpec represents a Samba-Operator SmbSecurityJoinSpec resource.
	JoinSpec = v1alpha1.SmbSecurityJoinSpec

	// UserJoinSpec represents a Samba-Operator SmbSecurityUserJoinSpec resource.
	UserJoinSpec = v1alpha1.SmbSecurityUserJoinSpec
)

type ShareData struct {
	Share          *Share
	CommonConfig   *CommonConfig
	SecurityConfig *SecurityConfig
}

type User struct {
	Secret string
	Key    string
}

//nolint:revive // allows this exported interface name for specific domain clarity.
type SMBCommonConfigRepo interface {
	List(ctx context.Context, scope, namespace, selector string) ([]CommonConfig, error)
	Get(ctx context.Context, scope, namespace, name string) (*CommonConfig, error)
	Update(ctx context.Context, scope, namespace string, cc *CommonConfig) (*CommonConfig, error)
	Create(ctx context.Context, scope, namespace string, cc *CommonConfig) (*CommonConfig, error)
	Delete(ctx context.Context, scope, namespace, name string) error
}

//nolint:revive // allows this exported interface name for specific domain clarity.
type SMBShareRepo interface {
	List(ctx context.Context, scope, namespace, selector string) ([]Share, error)
	Get(ctx context.Context, scope, namespace, name string) (*Share, error)
	Update(ctx context.Context, scope, namespace string, s *Share) (*Share, error)
	Create(ctx context.Context, scope, namespace string, s *Share) (*Share, error)
	Delete(ctx context.Context, scope, namespace, name string) error
}

//nolint:revive // allows this exported interface name for specific domain clarity.
type SMBSecurityConfigRepo interface {
	List(ctx context.Context, scope, namespace, selector string) ([]SecurityConfig, error)
	Get(ctx context.Context, scope, namespace, name string) (*SecurityConfig, error)
	Update(ctx context.Context, scope, namespace string, sc *SecurityConfig) (*SecurityConfig, error)
	Create(ctx context.Context, scope, namespace string, sc *SecurityConfig) (*SecurityConfig, error)
	Delete(ctx context.Context, scope, namespace, name string) error
}

type UseCase struct {
	smbCommonConfig   SMBCommonConfigRepo
	smbShare          SMBShareRepo
	smbSecurityConfig SMBSecurityConfigRepo
}

func NewUseCase(smbCommonConfig SMBCommonConfigRepo, smbShare SMBShareRepo, smbSecurityConfig SMBSecurityConfigRepo) *UseCase {
	return &UseCase{
		smbCommonConfig:   smbCommonConfig,
		smbShare:          smbShare,
		smbSecurityConfig: smbSecurityConfig,
	}
}

func (uc *UseCase) ListSMBShares(ctx context.Context, scope, namespace string) ([]ShareData, error) {
	commonConfigs, err := uc.smbCommonConfig.List(ctx, scope, namespace, "")
	if err != nil {
		return nil, err
	}

	commonConfigMap := map[string]*CommonConfig{}
	for i := range commonConfigs {
		cc := commonConfigs[i]
		commonConfigMap[cc.Name] = &cc
	}

	securityConfigs, err := uc.smbSecurityConfig.List(ctx, scope, namespace, "")
	if err != nil {
		return nil, err
	}

	securityConfigMap := map[string]*SecurityConfig{}
	for i := range securityConfigs {
		sc := securityConfigs[i]
		securityConfigMap[sc.Name] = &sc
	}

	shares, err := uc.smbShare.List(ctx, scope, namespace, "")
	if err != nil {
		return nil, err
	}

	ret := []ShareData{}

	for i := range shares {
		ret = append(ret, ShareData{
			Share:          &shares[i],
			CommonConfig:   commonConfigMap[shares[i].Spec.CommonConfig],
			SecurityConfig: securityConfigMap[shares[i].Spec.SecurityConfig],
		})
	}

	return ret, nil
}

//nolint:revive // WIP
func (uc *UseCase) CreateSMBShare(ctx context.Context, scope, namespace, name string, sizeBytes uint64, browsable, readOnly, guestOK bool, validUsers []string, mapToGuest, securityMode string, localUser *User, realm string, joinSources []User) (*ShareData, error) {
	// Implementation goes here
	return nil, nil
}

//nolint:revive // WIP
func (uc *UseCase) UpdateSMBShare(ctx context.Context, scope, namespace, name string, sizeBytes uint64, browsable, readOnly, guestOK bool, validUsers []string, mapToGuest, securityMode string, localUser *User, realm string, joinSources []User) (*ShareData, error) {
	// Implementation goes here
	return nil, nil
}
