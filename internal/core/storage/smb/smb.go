package smb

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/samba-in-kubernetes/samba-operator/api/v1alpha1"
	"golang.org/x/sync/errgroup"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/otterscale/otterscale/internal/core/application/config"
	"github.com/otterscale/otterscale/internal/core/application/service"
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
	LocalUsers     []User
	JoinSource     *User

	Deployment *workload.Deployment
	Pods       []workload.Pod
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}

type names struct {
	UsersSecret    string
	JoinSecret     string
	SecurityConfig string
	CommonConfig   string
	Share          string
}

type userEntry struct {
	Name     string `json:"name"`
	Password string `json:"password,omitempty"`
}

type sambaContainerConfig struct {
	SCCVersion string                 `json:"samba-container-config"`
	Users      map[string][]userEntry `json:"users,omitempty"`
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
	secret     config.SecretRepo
	service    service.ServiceRepo
}

func NewUseCase(smbCommonConfig SMBCommonConfigRepo, smbShare SMBShareRepo, smbSecurityConfig SMBSecurityConfigRepo, deployment workload.DeploymentRepo, pod workload.PodRepo, secret config.SecretRepo, service service.ServiceRepo) *UseCase {
	return &UseCase{
		smbCommonConfig:   smbCommonConfig,
		smbShare:          smbShare,
		smbSecurityConfig: smbSecurityConfig,
		deployment:        deployment,
		pod:               pod,
		secret:            secret,
		service:           service,
	}
}

func (uc *UseCase) ListSMBShares(ctx context.Context, scope, namespace string) (shareData []ShareData, host string, err error) {
	var (
		commonConfigMap   map[string]*CommonConfig
		securityConfigMap map[string]*SecurityConfig
		deploymentMap     map[string]*workload.Deployment
		allPods           []workload.Pod
		allShares         []Share
	)

	selector := "app.kubernetes.io/name=samba"
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
			deploymentMap = map[string]*workload.Deployment{}

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
		return nil, "", err
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
			return nil, "", fmt.Errorf("failed to create selector: %w", err)
		}

		ret = append(ret, ShareData{
			Share:          &share,
			CommonConfig:   commonConfigMap[share.Spec.CommonConfig],
			SecurityConfig: securityConfigMap[share.Spec.SecurityConfig],
			Deployment:     deployment,
			Pods:           uc.filterPods(selector, namespace, allPods),
		})
	}

	return ret, uc.service.Host(scope), nil
}

func (uc *UseCase) CreateSMBShare(ctx context.Context, scope, namespace, name string, sizeBytes uint64, browsable, readOnly, guestOK bool, validUsers []string, mapToGuest, securityMode string, localUsers []User, realm string, joinSource *User) (shareData *ShareData, err error) {
	names := uc.newNames(name)

	// Cleanup on error
	defer func() {
		if err != nil {
			_ = uc.secret.Delete(ctx, scope, namespace, names.UsersSecret)
			_ = uc.secret.Delete(ctx, scope, namespace, names.JoinSecret)
			_ = uc.smbSecurityConfig.Delete(ctx, scope, namespace, names.SecurityConfig)
			_ = uc.smbCommonConfig.Delete(ctx, scope, namespace, names.CommonConfig)
			_ = uc.smbShare.Delete(ctx, scope, namespace, names.Share)
		}
	}()

	// Create users secret
	if len(localUsers) > 0 {
		usersSecret := uc.buildUsersSecret(namespace, names.UsersSecret, localUsers)

		if _, err = uc.secret.Create(ctx, scope, namespace, usersSecret); err != nil {
			return nil, err
		}
	}

	// Create join secret
	if joinSource != nil {
		joinSecret := uc.buildJoinSecret(namespace, names.JoinSecret, joinSource)

		if _, err = uc.secret.Create(ctx, scope, namespace, joinSecret); err != nil {
			return nil, err
		}
	}

	// Create security config
	securityConfig := uc.buildSecurityConfig(namespace, names.SecurityConfig, securityMode, names.UsersSecret, realm, names.JoinSecret)

	_, err = uc.smbSecurityConfig.Create(ctx, scope, namespace, securityConfig)
	if err != nil {
		return nil, err
	}

	// Create common config
	commonConfig := uc.buildCommonConfig(namespace, names.CommonConfig, mapToGuest)

	_, err = uc.smbCommonConfig.Create(ctx, scope, namespace, commonConfig)
	if err != nil {
		return nil, err
	}

	// Create share
	shareConfig := uc.buildShareConfig(guestOK, validUsers)
	share := uc.buildShare(namespace, names.Share, browsable, readOnly, sizeBytes, names.SecurityConfig, names.CommonConfig, shareConfig)

	newShare, err := uc.smbShare.Create(ctx, scope, namespace, share)
	if err != nil {
		return nil, err
	}

	return &ShareData{
		Share: newShare,
	}, nil
}

func (uc *UseCase) UpdateSMBShare(ctx context.Context, scope, namespace, name string, sizeBytes uint64, browsable, readOnly, guestOK bool, validUsers []string, mapToGuest, securityMode string, localUsers []User, realm string, joinSource *User) (*ShareData, error) {
	names := uc.newNames(name)

	// Update users secret
	if len(localUsers) > 0 {
		usersSecret := uc.buildUsersSecret(namespace, names.UsersSecret, localUsers)

		if err := uc.upsertSecret(ctx, scope, namespace, usersSecret); err != nil {
			return nil, err
		}
	} else {
		_ = uc.secret.Delete(ctx, scope, namespace, names.UsersSecret)
	}

	// Create join secret
	if joinSource != nil {
		joinSecret := uc.buildJoinSecret(namespace, names.JoinSecret, joinSource)

		if err := uc.upsertSecret(ctx, scope, namespace, joinSecret); err != nil {
			return nil, err
		}
	} else {
		_ = uc.secret.Delete(ctx, scope, namespace, names.JoinSecret)
	}

	// Create security config
	securityConfig := uc.buildSecurityConfig(namespace, names.SecurityConfig, securityMode, names.UsersSecret, realm, names.JoinSecret)

	if _, err := uc.smbSecurityConfig.Update(ctx, scope, namespace, securityConfig); err != nil {
		return nil, err
	}

	// Create common config
	commonConfig := uc.buildCommonConfig(namespace, names.CommonConfig, mapToGuest)

	if _, err := uc.smbCommonConfig.Update(ctx, scope, namespace, commonConfig); err != nil {
		return nil, err
	}

	// Create share
	shareConfig := uc.buildShareConfig(guestOK, validUsers)
	share := uc.buildShare(namespace, names.Share, browsable, readOnly, sizeBytes, names.SecurityConfig, names.CommonConfig, shareConfig)

	newShare, err := uc.smbShare.Update(ctx, scope, namespace, share)
	if err != nil {
		return nil, err
	}

	return &ShareData{
		Share: newShare,
	}, nil
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

func (uc *UseCase) newNames(name string) names {
	return names{
		UsersSecret:    name + "-users",
		JoinSecret:     name + "-join",
		SecurityConfig: name + "-security",
		CommonConfig:   name + "-common",
		Share:          name + "-share",
	}
}

func (uc *UseCase) buildUsersSecret(namespace, name string, users []User) *config.Secret {
	userEntries := []userEntry{}

	for _, user := range users {
		userEntries = append(userEntries, userEntry{
			Name:     user.Username,
			Password: user.Password,
		})
	}

	data, _ := json.Marshal(&sambaContainerConfig{
		SCCVersion: "v0",
		Users: map[string][]userEntry{
			"all_entries": userEntries,
		},
	})

	return &config.Secret{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Type: corev1.SecretTypeOpaque,
		StringData: map[string]string{
			"users": string(data),
		},
	}
}

func (uc *UseCase) buildJoinSecret(namespace, name string, user *User) *config.Secret {
	data, _ := json.Marshal(&User{
		Username: user.Username,
		Password: user.Password,
	})

	return &config.Secret{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Type: corev1.SecretTypeOpaque,
		StringData: map[string]string{
			"join": string(data),
		},
	}
}

func (uc *UseCase) buildCommonConfig(namespace, name, mapToGuest string) *CommonConfig {
	return &CommonConfig{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1alpha1.SmbCommonConfigSpec{
			CustomGlobalConfig: &v1alpha1.SmbCommonConfigGlobalConfig{
				UseUnsafeCustomConfig: true,
				Configs: map[string]string{
					MapToGuestKey: mapToGuest,
				},
			},
		},
	}
}

func (uc *UseCase) buildSecurityConfig(namespace, name, securityMode, usersSecretName, realm, joinSecretName string) *SecurityConfig {
	return &SecurityConfig{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1alpha1.SmbSecurityConfigSpec{
			Mode: securityMode,
			Users: &v1alpha1.SmbSecurityUsersSpec{
				Secret: usersSecretName,
				Key:    "users",
			},
			DNS: &v1alpha1.SmbSecurityDNSSpec{
				Register: "cluster-ip",
			},
			Realm: realm,
			JoinSources: []v1alpha1.SmbSecurityJoinSpec{
				{
					UserJoin: &v1alpha1.SmbSecurityUserJoinSpec{
						Secret: joinSecretName,
						Key:    "join",
					},
				},
			},
		},
	}
}

func (uc *UseCase) buildShareConfig(guestOK bool, validUsers []string) *v1alpha1.SmbShareConfig {
	configs := map[string]string{}

	if guestOK {
		configs[GuestOKkey] = "yes"
	}

	if len(validUsers) > 0 {
		configs[ValidUsersKey] = strings.Join(validUsers, " ")
	}

	return &v1alpha1.SmbShareConfig{
		UseUnsafeCustomConfig: true,
		Configs:               configs,
	}
}

func (uc *UseCase) buildShare(namespace, name string, browsable, readOnly bool, sizeBytes uint64, securityConfigName, commonConfigName string, customShareConfig *v1alpha1.SmbShareConfig) *v1alpha1.SmbShare {
	storageClassName := "cephfs"

	return &Share{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1alpha1.SmbShareSpec{
			ShareName:  name,
			Browseable: browsable,
			ReadOnly:   readOnly,
			Storage: v1alpha1.SmbShareStorageSpec{
				Pvc: &v1alpha1.SmbSharePvcSpec{
					Spec: &corev1.PersistentVolumeClaimSpec{
						AccessModes: []corev1.PersistentVolumeAccessMode{
							corev1.ReadWriteMany,
						},
						Resources: corev1.VolumeResourceRequirements{
							Requests: corev1.ResourceList{
								"storage": *resource.NewQuantity(int64(sizeBytes), resource.BinarySI), //nolint:gosec // ignore
							},
						},
						StorageClassName: &storageClassName,
					},
				},
			},
			SecurityConfig:    securityConfigName,
			CommonConfig:      commonConfigName,
			CustomShareConfig: customShareConfig,
		},
	}
}

func (uc *UseCase) upsertSecret(ctx context.Context, scope, namespace string, secret *config.Secret) error {
	secret, err := uc.secret.Get(ctx, scope, namespace, secret.Name)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			_, err = uc.secret.Create(ctx, scope, namespace, secret)
			return err
		}
		return err
	}

	_, err = uc.secret.Update(ctx, scope, namespace, secret)
	return err
}
