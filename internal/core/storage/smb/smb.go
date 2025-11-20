package smb

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/go-ldap/ldap/v3"
	"github.com/samba-in-kubernetes/samba-operator/api/v1alpha1"
	"golang.org/x/sync/errgroup"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/otterscale/otterscale/internal/core/application/config"
	"github.com/otterscale/otterscale/internal/core/application/persistent"
	"github.com/otterscale/otterscale/internal/core/application/service"
	"github.com/otterscale/otterscale/internal/core/application/workload"
)

const (
	servicePollingTimeout  = 5 * time.Minute
	servicePollingInterval = 5 * time.Second
)

const (
	GuestOKkey    = "guest ok"
	MapToGuestKey = "map to guest"
	ValidUsersKey = "valid users"
)

const (
	ldapDefaultPort   uint16 = 389
	EntityTypeUnknown int    = iota
	EntityTypeUser
	EntityTypeGroup
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

type ADValidateResult struct {
	Valid      bool
	EntityType int
	Message    string
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
func (uc *UseCase) ListSMBShares(ctx context.Context, scope, namespace string) (data []ShareData, hostname string, err error) {
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

func (uc *UseCase) CreateSMBShare(ctx context.Context, scope, namespace, name string, sizeBytes uint64, port int32, browsable, readOnly, guestOK bool, validUsers []string, mapToGuest, securityMode string, localUsers []User, realm string, joinSource *User) (data *ShareData, hostname string, err error) {
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

	// Wait for service to be created automatically and update NodePort in background
	if port > 0 {
		go uc.waitAndUpdateServiceNodePort(context.Background(), scope, namespace, names.Share, port)
	}

	url, err := uc.service.URL(scope)
	if err != nil {
		return nil, "", err
	}

	return &ShareData{
		Share: newShare,
	}, url.Hostname(), nil
}

func (uc *UseCase) UpdateSMBShare(ctx context.Context, scope, namespace, name string, sizeBytes uint64, port int32, browsable, readOnly, guestOK bool, validUsers []string, mapToGuest, securityMode string, localUsers []User, realm string, joinSource *User) (data *ShareData, hostname string, err error) {
	names := uc.newNames(name)

	// Update secrets based on security mode
	if securityMode == "active-directory" {
		if joinSource != nil && joinSource.Username != "" && joinSource.Password != "" {
			joinSecret := uc.buildJoinSecret(namespace, names.JoinSecret, joinSource)
			if err := uc.upsertSecret(ctx, scope, namespace, names.JoinSecret, joinSecret); err != nil {
				return nil, "", err
			}
		}
		_ = uc.secret.Delete(ctx, scope, namespace, names.UsersSecret)
	} else {
		if len(localUsers) > 0 {
			usersSecret := uc.buildUsersSecret(namespace, names.UsersSecret, localUsers)
			if err := uc.upsertSecret(ctx, scope, namespace, names.UsersSecret, usersSecret); err != nil {
				return nil, "", err
			}
		}
		_ = uc.secret.Delete(ctx, scope, namespace, names.JoinSecret)
	}

	// Update security config
	securityConfig := uc.buildSecurityConfig(namespace, names.SecurityConfig, securityMode, names.UsersSecret, realm, names.JoinSecret)
	updatedSecurityConfig, err := uc.upsertSecurityConfig(ctx, scope, namespace, securityConfig)
	if err != nil {
		return nil, "", err
	}

	// Update common config
	commonConfig := uc.buildCommonConfig(namespace, names.CommonConfig, mapToGuest)
	updatedCommonConfig, err := uc.upsertCommonConfig(ctx, scope, namespace, commonConfig)
	if err != nil {
		return nil, "", err
	}

	// Update share
	shareConfig := uc.buildShareConfig(guestOK, validUsers)
	share := uc.buildShare(namespace, names.Share, browsable, readOnly, sizeBytes, names.SecurityConfig, names.CommonConfig, shareConfig)
	newShare, err := uc.upsertShare(ctx, scope, namespace, share)
	if err != nil {
		return nil, "", err
	}

	// Update PVC if size changed
	if sizeBytes > 0 {
		pvcName := names.Share + "-pvc"
		pvc := uc.buildPVC(namespace, pvcName, sizeBytes)
		if err := uc.upsertPVC(ctx, scope, namespace, pvc); err != nil {
			return nil, "", fmt.Errorf("failed to update PVC: %w", err)
		}
	}

	// Update service with NodePort
	var svc *service.Service
	if port > 0 {
		updatedSvc, err := uc.setServiceNodePort(ctx, scope, namespace, names.Share, port)
		if err != nil {
			return nil, "", err
		}
		svc = updatedSvc
	} else {
		retrievedSvc, err := uc.service.Get(ctx, scope, namespace, names.Share)
		if err != nil && !k8serrors.IsNotFound(err) {
			return nil, "", fmt.Errorf("failed to get service: %w", err)
		}
		svc = retrievedSvc
	}

	url, err := uc.service.URL(scope)
	if err != nil {
		return nil, "", err
	}

	return &ShareData{
		Share:          newShare,
		CommonConfig:   updatedCommonConfig,
		SecurityConfig: updatedSecurityConfig,
		LocalUsers:     localUsers,
		JoinSource:     joinSource,
		Service:        svc,
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

func (uc *UseCase) newNames(name string) names {
	return names{
		UsersSecret:    name + "-users",
		JoinSecret:     name + "-join",
		SecurityConfig: name + "-security",
		CommonConfig:   name + "-common",
		Share:          name + "-share",
	}
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
			Labels: map[string]string{
				"app.kubernetes.io/name": "samba",
			},
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
			Labels: map[string]string{
				"app.kubernetes.io/name": "samba",
			},
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

func (uc *UseCase) buildPVC(namespace, name string, sizeBytes uint64) *persistent.PersistentVolumeClaim {
	return &persistent.PersistentVolumeClaim{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: *resource.NewQuantity(int64(sizeBytes), resource.BinarySI), //nolint:gosec // ignore
				},
			},
		},
	}
}

func (uc *UseCase) buildService(namespace, name string, nodePort int32) *service.Service {
	return &service.Service{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeNodePort,
			Ports: []corev1.ServicePort{
				{
					Name:     "smb",
					Port:     445,
					NodePort: nodePort,
					Protocol: corev1.ProtocolTCP,
				},
			},
		},
	}
}

func (uc *UseCase) upsertSecret(ctx context.Context, scope, namespace, name string, secret *config.Secret) error {
	existingSecret, err := uc.secret.Get(ctx, scope, namespace, name)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			_, err = uc.secret.Create(ctx, scope, namespace, secret)
			return err
		}
		return err
	}

	// Update existing secret with new data
	existingSecret.StringData = secret.StringData
	existingSecret.Data = secret.Data
	if existingSecret.Labels == nil {
		existingSecret.Labels = map[string]string{}
	}
	for k, v := range secret.Labels {
		existingSecret.Labels[k] = v
	}

	_, err = uc.secret.Update(ctx, scope, namespace, existingSecret)
	return err
}

func (uc *UseCase) upsertSecurityConfig(ctx context.Context, scope, namespace string, securityConfig *SecurityConfig) (*SecurityConfig, error) {
	existing, err := uc.smbSecurityConfig.Get(ctx, scope, namespace, securityConfig.ObjectMeta.Name) //nolint:staticcheck // ObjectMeta is embedded
	if err != nil {
		return nil, err
	}

	// Update specific fields while preserving metadata
	existing.Spec.Mode = securityConfig.Spec.Mode
	existing.Spec.Users = securityConfig.Spec.Users
	existing.Spec.DNS = securityConfig.Spec.DNS
	existing.Spec.Realm = securityConfig.Spec.Realm
	existing.Spec.JoinSources = securityConfig.Spec.JoinSources

	return uc.smbSecurityConfig.Update(ctx, scope, namespace, existing)
}

func (uc *UseCase) upsertCommonConfig(ctx context.Context, scope, namespace string, commonConfig *CommonConfig) (*CommonConfig, error) {
	existing, err := uc.smbCommonConfig.Get(ctx, scope, namespace, commonConfig.ObjectMeta.Name) //nolint:staticcheck // ObjectMeta is embedded
	if err != nil {
		return nil, err
	}

	// Update specific fields while preserving metadata
	existing.Spec.CustomGlobalConfig = commonConfig.Spec.CustomGlobalConfig

	return uc.smbCommonConfig.Update(ctx, scope, namespace, existing)
}

func (uc *UseCase) upsertShare(ctx context.Context, scope, namespace string, share *Share) (*Share, error) {
	existing, err := uc.smbShare.Get(ctx, scope, namespace, share.ObjectMeta.Name) //nolint:staticcheck // ObjectMeta is embedded
	if err != nil {
		return nil, err
	}

	// Update specific fields while preserving metadata and config references
	existing.Spec.Browseable = share.Spec.Browseable
	existing.Spec.ReadOnly = share.Spec.ReadOnly
	existing.Spec.CustomShareConfig = share.Spec.CustomShareConfig
	if share.Spec.Storage.Pvc != nil && share.Spec.Storage.Pvc.Spec != nil {
		if existing.Spec.Storage.Pvc == nil {
			existing.Spec.Storage.Pvc = &v1alpha1.SmbSharePvcSpec{}
		}
		if existing.Spec.Storage.Pvc.Spec == nil {
			existing.Spec.Storage.Pvc.Spec = &corev1.PersistentVolumeClaimSpec{}
		}
		existing.Spec.Storage.Pvc.Spec.Resources = share.Spec.Storage.Pvc.Spec.Resources
	}

	return uc.smbShare.Update(ctx, scope, namespace, existing)
}

func (uc *UseCase) upsertPVC(ctx context.Context, scope, namespace string, pvc *persistent.PersistentVolumeClaim) error {
	existing, err := uc.persistentVolumeClaim.Get(ctx, scope, namespace, pvc.Name)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
		return err
	}

	existing.Spec.Resources = pvc.Spec.Resources

	_, err = uc.persistentVolumeClaim.Update(ctx, scope, namespace, existing)
	return err
}

func (uc *UseCase) upsertService(ctx context.Context, scope, namespace string, svc *service.Service) (*service.Service, error) {
	existing, err := uc.service.Get(ctx, scope, namespace, svc.Name)
	if err != nil {
		return nil, err
	}

	// Update spec while preserving metadata
	existing.Spec = svc.Spec

	return uc.service.Update(ctx, scope, namespace, existing)
}

func (uc *UseCase) setServiceNodePort(ctx context.Context, scope, namespace, name string, nodePort int32) (*service.Service, error) {
	svc := uc.buildService(namespace, name, nodePort)
	return uc.upsertService(ctx, scope, namespace, svc)
}

func (uc *UseCase) parseSearchUsername(searchUsername string) (string, *ADValidateResult) {
	result := &ADValidateResult{EntityType: EntityTypeUnknown}
	actualSearchName := searchUsername

	if strings.HasPrefix(searchUsername, "@\"") {
		// Must end with " if starts with @"
		if !strings.HasSuffix(searchUsername, "\"") {
			result.Valid = false
			result.Message = "invalid format: missing closing quote"
			return "", result
		}

		// Remove @" prefix and " suffix
		actualSearchName = strings.TrimPrefix(searchUsername, "@\"")
		actualSearchName = strings.TrimSuffix(actualSearchName, "\"")

		// Reject invalid format with forward slash
		if strings.Contains(actualSearchName, "/") {
			result.Valid = false
			result.Message = "invalid format: use backslash (\\) instead of forward slash (/)"
			return "", result
		}

		// If format is "DOMAIN\NAME", extract the NAME part
		if idx := strings.Index(actualSearchName, "\\"); idx != -1 {
			actualSearchName = actualSearchName[idx+1:]
		}
	}

	actualSearchName = strings.TrimSpace(actualSearchName)
	if actualSearchName == "" {
		result.Valid = false
		result.Message = "invalid format: empty username or group name"
		return "", result
	}

	return actualSearchName, nil
}

func (uc *UseCase) resolveLDAPServer(realm string) (serverName string, port uint16, err error) {
	resolver := &net.Resolver{}
	_, srvs, err := resolver.LookupSRV(context.Background(), "ldap", "tcp", realm)
	if err == nil && len(srvs) > 0 {
		return strings.TrimSuffix(srvs[0].Target, "."), srvs[0].Port, nil
	}

	if addrs, err := resolver.LookupHost(context.Background(), realm); err == nil && len(addrs) > 0 {
		return realm, ldapDefaultPort, nil
	}

	return "", 0, fmt.Errorf("failed to lookup LDAP server: unable to resolve %s", realm)
}

func (uc *UseCase) connectLDAP(serverName string, port uint16, useTLS bool) (*ldap.Conn, error) {
	ldapURL := fmt.Sprintf("%s:%d", serverName, port)

	if useTLS {
		conn, err := ldap.DialURL(fmt.Sprintf("ldaps://%s", ldapURL),
			ldap.DialWithTLSConfig(&tls.Config{
				ServerName:         serverName,
				InsecureSkipVerify: false,
				MinVersion:         tls.VersionTLS12,
			}))
		if err != nil {
			return nil, fmt.Errorf("failed to connect to LDAP server: %w", err)
		}
		return conn, nil
	}

	conn, err := ldap.DialURL(fmt.Sprintf("ldap://%s", ldapURL))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to LDAP server: %w", err)
	}
	return conn, nil
}

func (uc *UseCase) ADValidate(_ context.Context, realm, username, password, searchUsername string, useTLS bool) (*ADValidateResult, error) {
	result := &ADValidateResult{EntityType: EntityTypeUnknown}
	searchUsername = strings.TrimSpace(searchUsername)

	// Parse and validate searchUsername format
	actualSearchName, parseResult := uc.parseSearchUsername(searchUsername)
	if parseResult != nil {
		return parseResult, nil
	}

	// Resolve LDAP server
	serverName, port, err := uc.resolveLDAPServer(realm)
	if err != nil {
		result.Message = "LDAP server not found"
		return result, err
	}

	// Connect to LDAP server
	conn, err := uc.connectLDAP(serverName, port, useTLS)
	if err != nil {
		result.Message = "failed to connect"
		return result, err
	}
	defer conn.Close()

	// Format username for binding
	if !strings.Contains(username, "@") && !strings.Contains(username, "=") {
		username = fmt.Sprintf("%s@%s", username, strings.ToUpper(realm))
	}

	// Bind to LDAP
	if err = conn.Bind(username, password); err != nil {
		result.Message = "authentication failed"
		return result, fmt.Errorf("LDAP bind failed")
	}

	// Search for user or group
	sr, err := conn.Search(ldap.NewSearchRequest(
		"DC="+strings.Join(strings.Split(strings.ToLower(realm), "."), ",DC="),
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(sAMAccountName=%s)", ldap.EscapeFilter(actualSearchName)),
		[]string{"objectClass"}, nil,
	))
	if err != nil {
		result.Message = "search failed"
		return result, fmt.Errorf("LDAP search failed")
	}

	if len(sr.Entries) == 0 {
		result.Message = "user or group not found"
		return result, nil
	}

	result.Valid = true
	result.EntityType = determineEntityType(sr.Entries[0].GetAttributeValues("objectClass"))
	result.Message = "validation successful"

	return result, nil
}

func determineEntityType(objectClasses []string) int {
	for _, c := range objectClasses {
		if cl := strings.ToLower(c); cl == "group" {
			return EntityTypeGroup
		}
	}
	for _, c := range objectClasses {
		if cl := strings.ToLower(c); cl == "user" || cl == "person" || cl == "inetorgperson" {
			return EntityTypeUser
		}
	}
	return EntityTypeUnknown
}

func (uc *UseCase) waitAndUpdateServiceNodePort(parentCtx context.Context, scope, namespace, name string, port int32) {
	ctx, cancel := context.WithTimeout(parentCtx, servicePollingTimeout)
	defer cancel()

	ticker := time.NewTicker(servicePollingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			_, err := uc.service.Get(ctx, scope, namespace, name)
			if err == nil {
				_, _ = uc.setServiceNodePort(context.Background(), scope, namespace, name, port)
				return
			}
		}
	}
}
