package smb

import (
	"context"
	"fmt"

	"github.com/samba-in-kubernetes/samba-operator/api/v1alpha1"
	"golang.org/x/sync/errgroup"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/otterscale/otterscale/internal/core/application/release"
	"github.com/otterscale/otterscale/internal/core/application/workload"
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
	Deployment     *workload.Deployment
	Pods           []workload.Pod
}

type User struct {
	Username string
	Password string
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

	deployment workload.DeploymentRepo
	pod        workload.PodRepo
}

func NewUseCase(smbCommonConfig SMBCommonConfigRepo, smbShare SMBShareRepo, smbSecurityConfig SMBSecurityConfigRepo, deployment workload.DeploymentRepo, pod workload.PodRepo) *UseCase {
	return &UseCase{
		smbCommonConfig:   smbCommonConfig,
		smbShare:          smbShare,
		smbSecurityConfig: smbSecurityConfig,
		deployment:        deployment,
		pod:               pod,
	}
}

func (uc *UseCase) ListSMBShares(ctx context.Context, scope, namespace string) ([]ShareData, error) {
	var (
		commonConfigMap   map[string]*CommonConfig
		securityConfigMap map[string]*SecurityConfig
		deploymentMap     map[string]*workload.Deployment
		allPods           []workload.Pod
		allShares         []Share
	)

	selector := release.TypeLabel + "=" + "samba"
	eg, egctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		commonConfigs, err := uc.smbCommonConfig.List(egctx, scope, namespace, "")
		if err == nil {
			commonConfigMap = map[string]*CommonConfig{}

			for i := range commonConfigs {
				cc := commonConfigs[i]
				commonConfigMap[cc.Name] = &cc
			}
		}

		return err
	})

	eg.Go(func() error {
		securityConfigs, err := uc.smbSecurityConfig.List(egctx, scope, namespace, "")
		if err == nil {
			securityConfigMap = map[string]*SecurityConfig{}

			for i := range securityConfigs {
				sc := securityConfigs[i]
				securityConfigMap[sc.Name] = &sc
			}
		}

		return err
	})

	eg.Go(func() error {
		deployments, err := uc.deployment.List(egctx, scope, namespace, selector)
		if err == nil {
			deploymentMap := map[string]*workload.Deployment{}

			for i := range deployments {
				d := deployments[i]
				deploymentMap[d.Name] = &d
			}
		}

		return err
	})

	eg.Go(func() error {
		pods, err := uc.pod.List(egctx, scope, namespace, selector)
		if err == nil {
			allPods = pods
		}
		return err
	})

	eg.Go(func() error {
		shares, err := uc.smbShare.List(egctx, scope, namespace, "")
		if err == nil {
			allShares = shares
		}
		return err
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	ret := []ShareData{}

	for i := range allShares {
		share := allShares[i]

		deployment, ok := deploymentMap[share.Name]
		if !ok {
			continue
		}

		selector, err := v1.LabelSelectorAsSelector(deployment.Spec.Selector)
		if err != nil {
			return nil, fmt.Errorf("failed to create selector: %w", err)
		}

		ret = append(ret, ShareData{
			Share:          &share,
			CommonConfig:   commonConfigMap[share.Spec.CommonConfig],
			SecurityConfig: securityConfigMap[share.Spec.SecurityConfig],
			Deployment:     deployment,
			Pods:           uc.filterPods(selector, namespace, allPods),
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

func (uc *UseCase) filterPods(selector labels.Selector, namespace string, pods []workload.Pod) []workload.Pod {
	ret := []workload.Pod{}

	for i := range pods {
		if pods[i].Namespace == namespace && selector.Matches(labels.Set(pods[i].Labels)) {
			ret = append(ret, pods[i])
		}
	}

	return ret
}
