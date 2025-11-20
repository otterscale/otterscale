package smb

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/samba-in-kubernetes/samba-operator/api/v1alpha1"
	"golang.org/x/sync/errgroup"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/otterscale/otterscale/internal/core/application/config"
	"github.com/otterscale/otterscale/internal/core/application/persistent"
	"github.com/otterscale/otterscale/internal/core/application/service"
	"github.com/otterscale/otterscale/internal/core/application/workload"
)

const namespace = "samba-operator-system"

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

	Service    *service.Service
	Deployment *workload.Deployment
	Pods       []workload.Pod
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}

type names struct {
	UsersSecret           string
	JoinSecret            string
	SecurityConfig        string
	CommonConfig          string
	Share                 string
	PersistentVolumeClaim string
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

	deployment            workload.DeploymentRepo
	pod                   workload.PodRepo
	secret                config.SecretRepo
	service               service.ServiceRepo
	persistentVolumeClaim persistent.PersistentVolumeClaimRepo
}

func NewUseCase(smbCommonConfig SMBCommonConfigRepo, smbShare SMBShareRepo, smbSecurityConfig SMBSecurityConfigRepo, deployment workload.DeploymentRepo, pod workload.PodRepo, secret config.SecretRepo, service service.ServiceRepo, persistentVolumeClaim persistent.PersistentVolumeClaimRepo) *UseCase {
	return &UseCase{
		smbCommonConfig:       smbCommonConfig,
		smbShare:              smbShare,
		smbSecurityConfig:     smbSecurityConfig,
		deployment:            deployment,
		pod:                   pod,
		secret:                secret,
		service:               service,
		persistentVolumeClaim: persistentVolumeClaim,
	}
}

//nolint:funlen // ignore
func (uc *UseCase) ListSMBShares(ctx context.Context, scope string) (data []ShareData, hostname string, err error) {
	var (
		commonConfigMap   map[string]*CommonConfig
		securityConfigMap map[string]*SecurityConfig
		deploymentMap     map[string]*workload.Deployment
		secretMap         map[string]*config.Secret
		allServices       []service.Service
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
				commonConfigMap[commonConfigs[i].Name] = &commonConfigs[i]
			}
		}

		return err
	})

	eg.Go(func() error {
		securityConfigs, err := uc.smbSecurityConfig.List(egctx, scope, namespace, "")
		if err == nil {
			securityConfigMap = map[string]*SecurityConfig{}

			for i := range securityConfigs {
				securityConfigMap[securityConfigs[i].Name] = &securityConfigs[i]
			}
		}

		return err
	})

	eg.Go(func() error {
		deployments, err := uc.deployment.List(egctx, scope, namespace, selector)
		if err == nil {
			deploymentMap = map[string]*workload.Deployment{}

			for i := range deployments {
				deploymentMap[deployments[i].Name] = &deployments[i]
			}
		}

		return err
	})

	eg.Go(func() error {
		secrets, err := uc.secret.List(egctx, scope, namespace, selector)
		if err == nil {
			secretMap = map[string]*config.Secret{}

			for i := range secrets {
				secretMap[secrets[i].Name] = &secrets[i]
			}
		}
		return err
	})

	eg.Go(func() error {
		services, err := uc.service.List(egctx, scope, namespace, selector)
		if err == nil {
			allServices = services
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

		names := uc.newNames(share.Name)
		localUsers, joinSource, err := uc.extractSecrets(secretMap, names.UsersSecret, names.JoinSecret)
		if err != nil {
			return nil, "", err
		}

		ret = append(ret, ShareData{
			Share:          &share,
			CommonConfig:   commonConfigMap[share.Spec.CommonConfig],
			SecurityConfig: securityConfigMap[share.Spec.SecurityConfig],
			LocalUsers:     localUsers,
			JoinSource:     joinSource,
			Service:        uc.filterService(deployment.Spec.Template.Labels, namespace, allServices),
			Deployment:     deployment,
			Pods:           uc.filterPods(selector, namespace, allPods),
		})
	}

	url, err := uc.service.URL(scope)
	if err != nil {
		return nil, "", err
	}

	return ret, url.Hostname(), nil
}

func (uc *UseCase) CreateSMBShare(ctx context.Context, scope, name string, sizeBytes uint64, port int32, browsable, readOnly, guestOK bool, validUsers []string, mapToGuest, securityMode string, localUsers []User, realm string, joinSource *User) (data *ShareData, hostname string, err error) {
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
			return nil, "", err
		}
	}

	// Create join secret
	if joinSource != nil && joinSource.Username != "" && joinSource.Password != "" {
		joinSecret := uc.buildJoinSecret(namespace, names.JoinSecret, joinSource)

		if _, err = uc.secret.Create(ctx, scope, namespace, joinSecret); err != nil {
			return nil, "", err
		}
	}

	// Create security config
	securityConfig := uc.buildSecurityConfig(namespace, names.SecurityConfig, securityMode, names.UsersSecret, realm, names.JoinSecret)

	_, err = uc.smbSecurityConfig.Create(ctx, scope, namespace, securityConfig)
	if err != nil {
		return nil, "", err
	}

	// Create common config
	commonConfig := uc.buildCommonConfig(namespace, names.CommonConfig, mapToGuest)

	_, err = uc.smbCommonConfig.Create(ctx, scope, namespace, commonConfig)
	if err != nil {
		return nil, "", err
	}

	// Create share
	shareConfig := uc.buildShareConfig(guestOK, validUsers)
	share := uc.buildShare(namespace, names.Share, browsable, readOnly, sizeBytes, names.SecurityConfig, names.CommonConfig, shareConfig)

	newShare, err := uc.smbShare.Create(ctx, scope, namespace, share)
	if err != nil {
		return nil, "", err
	}

	url, err := uc.service.URL(scope)
	if err != nil {
		return nil, "", err
	}

	// Wait for service to be created automatically and update NodePort in background
	go uc.waitAndUpdateServiceNodePort(context.Background(), scope, namespace, names.Share, port)

	return &ShareData{
		Share: newShare,
	}, url.Hostname(), nil
}

func (uc *UseCase) UpdateSMBShare(ctx context.Context, scope, name string, sizeBytes uint64, port int32, browsable, readOnly, guestOK bool, validUsers []string, mapToGuest, securityMode string, localUsers []User, realm string, joinSource *User) (data *ShareData, hostname string, err error) {
	names := uc.newNames(name)

	// Update users secret
	if len(localUsers) > 0 {
		usersSecret := uc.buildUsersSecret(namespace, names.UsersSecret, localUsers)

		if err := uc.upsertUsersSecret(ctx, scope, namespace, names.UsersSecret, usersSecret); err != nil {
			return nil, "", err
		}
	} else {
		_ = uc.secret.Delete(ctx, scope, namespace, names.UsersSecret)
	}

	// Update join secret
	if joinSource != nil && joinSource.Username != "" && joinSource.Password != "" {
		joinSecret := uc.buildJoinSecret(namespace, names.JoinSecret, joinSource)

		if err := uc.upsertJoinSecret(ctx, scope, namespace, names.JoinSecret, joinSecret); err != nil {
			return nil, "", err
		}
	} else {
		_ = uc.secret.Delete(ctx, scope, namespace, names.JoinSecret)
	}

	// Update security config
	securityConfig := uc.buildSecurityConfig(namespace, names.SecurityConfig, securityMode, names.UsersSecret, realm, names.JoinSecret)

	if err := uc.updateSecurityConfig(ctx, scope, namespace, names.SecurityConfig, securityConfig); err != nil {
		return nil, "", err
	}

	// Update common config
	commonConfig := uc.buildCommonConfig(namespace, names.CommonConfig, mapToGuest)

	if err := uc.updateCommonConfig(ctx, scope, namespace, names.CommonConfig, commonConfig); err != nil {
		return nil, "", err
	}

	// Update share
	shareConfig := uc.buildShareConfig(guestOK, validUsers)
	share := uc.buildShare(namespace, names.Share, browsable, readOnly, sizeBytes, names.SecurityConfig, names.CommonConfig, shareConfig)

	if err := uc.updateShare(ctx, scope, namespace, names.Share, share); err != nil {
		return nil, "", err
	}

	// Update persistent volume claim
	pvc := uc.buildPersistentVolumeClaim(namespace, names.PersistentVolumeClaim, sizeBytes)

	if err := uc.updatePersistentVolumeClaim(ctx, scope, namespace, pvc); err != nil {
		return nil, "", err
	}

	// Update service
	service := uc.buildService(namespace, names.Share, port)

	if err := uc.updateService(ctx, scope, namespace, names.Share, service); err != nil {
		return nil, "", err
	}

	// Get updated share
	newShare, err := uc.smbShare.Get(ctx, scope, namespace, names.Share)
	if err != nil {
		return nil, "", err
	}

	url, err := uc.service.URL(scope)
	if err != nil {
		return nil, "", err
	}

	return &ShareData{
		Share: newShare,
	}, url.Hostname(), nil
}

// only return the first matched service
func (uc *UseCase) filterService(podLabels map[string]string, namespace string, services []service.Service) *service.Service {
	for i := range services {
		selector := labels.SelectorFromSet(services[i].Spec.Selector)

		if services[i].Namespace == namespace && selector.Matches(labels.Set(podLabels)) {
			return &services[i]
		}
	}

	return nil
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

func (uc *UseCase) extractSecrets(secretMap map[string]*config.Secret, usersSecretName, joinSecretName string) (localUsers []User, joinSource *User, err error) {
	var config *sambaContainerConfig

	usersSecret, ok := secretMap[usersSecretName]
	if ok {
		if err := json.Unmarshal(usersSecret.Data["users"], &config); err != nil {
			return nil, nil, fmt.Errorf("failed to unmarshal join secret: %w", err)
		}

		if entries, ok := config.Users["all_entries"]; ok {
			for _, entry := range entries {
				localUsers = append(localUsers, User{
					Username: entry.Name,
					Password: entry.Password,
				})
			}
		}
	}

	joinSecret, ok := secretMap[joinSecretName]
	if ok {
		if err := json.Unmarshal(joinSecret.Data["join"], &joinSource); err != nil {
			return nil, nil, fmt.Errorf("failed to unmarshal join secret: %w", err)
		}
	}

	return localUsers, joinSource, nil
}

func (uc *UseCase) waitAndUpdateServiceNodePort(ctx context.Context, scope, namespace, name string, port int32) {
	const (
		timeout  = 5 * time.Minute
		interval = 5 * time.Second
	)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	newService := uc.buildService(namespace, name, port)

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			if err := uc.updateService(ctx, scope, namespace, name, newService); err != nil {
				if k8serrors.IsNotFound(err) {
					continue
				}
				slog.Error("failed to update service", "error", err)
				return
			}
		}
	}
}

func (uc *UseCase) newNames(name string) names {
	return names{
		UsersSecret:           name + "-users",
		JoinSecret:            name + "-join",
		SecurityConfig:        name + "-security",
		CommonConfig:          name + "-common",
		Share:                 name + "-share",
		PersistentVolumeClaim: name + "-pvc",
	}
}
