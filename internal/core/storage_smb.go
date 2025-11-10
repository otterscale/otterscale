package core

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
)

const (
	sambaConfigPrefix = `{
  "samba-container-config": "v0",
  "users": {
    "all_entries": [`
	sambaConfigSuffix = `
    ]
  }
}`
	sambaUserEntryTemplate = `
      {
        "name": %q,
        "password": %q
      }`
)

type SMBShare struct {
	Name         string
	Namespace    string
	Scope        string
	Facility     string
	SizeBytes    uint64
	Browseable   bool
	ReadOnly     bool
	GuestOk      bool
	MapToGuest   string
	SecurityMode string
	Realm        string
	ValidUsers   []string
	ADAuth       *ADAuth
	LocalAuth    *LocalAuth
	Status       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type SMBUser struct {
	Username string
	Password string
}

type ADAuth struct {
	Realm        string
	JoinUsername string
	JoinPassword string
}

type LocalAuth struct {
	Users []SMBUser
}

type SMBShareConfig struct {
	Browseable   bool
	ReadOnly     bool
	GuestOk      bool
	ValidUsers   string
	MapToGuest   string
	GuestAccount string
}

func (uc *StorageUseCase) ListSMBShares(ctx context.Context, scope, facility, namespace string) ([]SMBShare, error) {
	restConfig, err := kubeConfig(ctx, uc.facility, uc.action, scope, facility)
	if err != nil {
		return nil, err
	}

	return uc.kubeSMB.ListSMBShares(ctx, restConfig, namespace)
}

func (uc *StorageUseCase) GetSMBShare(ctx context.Context, scope, facility, namespace, name string) (*SMBShare, error) {
	restConfig, err := kubeConfig(ctx, uc.facility, uc.action, scope, facility)
	if err != nil {
		return nil, err
	}

	return uc.kubeSMB.GetSMBShare(ctx, restConfig, namespace, name)
}

func (uc *StorageUseCase) createUserSecret(ctx context.Context, restConfig *rest.Config, namespace, secretName string, users []SMBUser) error {
	secretData := sambaConfigPrefix

	for i, user := range users {
		if i > 0 {
			secretData += ","
		}
		secretData += fmt.Sprintf(sambaUserEntryTemplate, user.Username, user.Password)
	}
	secretData += sambaConfigSuffix

	data := map[string]string{
		"users": secretData,
	}
	_, err := uc.kubeCore.CreateSecret(ctx, restConfig, namespace, secretName, corev1.SecretTypeOpaque, data)
	return err
}

func (uc *StorageUseCase) createJoinSecret(ctx context.Context, restConfig *rest.Config, namespace, secretName string, adAuth *ADAuth) error {
	joinSecretData := fmt.Sprintf(`{"username": %q, "password": %q}`, adAuth.JoinUsername, adAuth.JoinPassword)
	data := map[string]string{
		"join": joinSecretData,
	}
	_, err := uc.kubeCore.CreateSecret(ctx, restConfig, namespace, secretName, corev1.SecretTypeOpaque, data)
	return err
}

type smbShareResources struct {
	hasCommonConfig   bool
	hasJoinSecret     bool
	hasUserSecret     bool
	hasSecurityConfig bool
}

func (r *smbShareResources) cleanup(ctx context.Context, uc *StorageUseCase, restConfig *rest.Config, namespace, commonConfigName, joinSecretName, userSecretName, securityConfigName string) {
	if r.hasUserSecret {
		_ = uc.kubeCore.DeleteSecret(ctx, restConfig, namespace, userSecretName)
	}
	if r.hasSecurityConfig {
		_ = uc.kubeSMB.DeleteSMBSecurityConfig(ctx, restConfig, namespace, securityConfigName)
	}
	if r.hasJoinSecret {
		_ = uc.kubeCore.DeleteSecret(ctx, restConfig, namespace, joinSecretName)
	}
	if r.hasCommonConfig {
		_ = uc.kubeSMB.DeleteSMBCommonConfig(ctx, restConfig, namespace, commonConfigName)
	}
}

func (uc *StorageUseCase) CreateSMBShare(ctx context.Context, scope, facility, namespace, name string, config *SMBShareConfig, sizeBytes uint64, securityMode string, adAuth *ADAuth, localAuth *LocalAuth) (*SMBShare, error) {
	restConfig, err := kubeConfig(ctx, uc.facility, uc.action, scope, facility)
	if err != nil {
		return nil, err
	}

	_, err = uc.kubeCore.GetNamespace(ctx, restConfig, namespace)
	if err != nil {
		if _, err := uc.kubeCore.CreateNamespace(ctx, restConfig, namespace); err != nil {
			return nil, fmt.Errorf("failed to create namespace: %w", err)
		}
	}

	commonConfigName := name + "-common"
	securityConfigName := name + "-security"
	userSecretName := name + "-users"
	joinSecretName := name + "-join"

	resources := &smbShareResources{}

	if err := uc.kubeSMB.CreateSMBCommonConfig(ctx, restConfig, namespace, commonConfigName, config); err != nil {
		return nil, fmt.Errorf("failed to create common config: %w", err)
	}
	resources.hasCommonConfig = true

	if adAuth != nil && adAuth.JoinUsername != "" && adAuth.JoinPassword != "" {
		if err := uc.createJoinSecret(ctx, restConfig, namespace, joinSecretName, adAuth); err != nil {
			resources.cleanup(ctx, uc, restConfig, namespace, commonConfigName, joinSecretName, userSecretName, securityConfigName)
			return nil, fmt.Errorf("failed to create join secret: %w", err)
		}
		resources.hasJoinSecret = true
	}

	if localAuth != nil && len(localAuth.Users) > 0 {
		if err := uc.createUserSecret(ctx, restConfig, namespace, userSecretName, localAuth.Users); err != nil {
			resources.cleanup(ctx, uc, restConfig, namespace, commonConfigName, joinSecretName, userSecretName, securityConfigName)
			return nil, fmt.Errorf("failed to create user secret: %w", err)
		}
		resources.hasUserSecret = true
	}

	realm := ""
	secretName := ""
	if adAuth != nil {
		realm = adAuth.Realm
		if resources.hasJoinSecret {
			secretName = joinSecretName
		}
	} else if resources.hasUserSecret {
		secretName = userSecretName
	}

	if err := uc.kubeSMB.CreateSMBSecurityConfig(ctx, restConfig, namespace, securityConfigName, securityMode, realm, secretName); err != nil {
		resources.cleanup(ctx, uc, restConfig, namespace, commonConfigName, joinSecretName, userSecretName, securityConfigName)
		return nil, fmt.Errorf("failed to create security config: %w", err)
	}
	resources.hasSecurityConfig = true

	if err := uc.kubeSMB.CreateSMBShare(ctx, restConfig, namespace, name, sizeBytes, config.Browseable, config.ReadOnly, config.GuestOk, config.ValidUsers, commonConfigName, securityConfigName, ""); err != nil {
		resources.cleanup(ctx, uc, restConfig, namespace, commonConfigName, joinSecretName, userSecretName, securityConfigName)
		return nil, fmt.Errorf("failed to create SMB share CRD: %w", err)
	}

	return uc.kubeSMB.GetSMBShare(ctx, restConfig, namespace, name)
}

func (uc *StorageUseCase) mergeAndUpdateUserSecret(ctx context.Context, restConfig *rest.Config, namespace, secretName string, newUsers []SMBUser) error {
	existingSecret, err := uc.kubeCore.GetSecret(ctx, restConfig, namespace, secretName)
	if err != nil {
		return fmt.Errorf("failed to get existing user secret: %w", err)
	}

	existingUsers := []SMBUser{}
	if existingData, ok := existingSecret.Data["users"]; ok && len(existingData) > 0 {
		type SambaUser struct {
			Name     string `json:"name"`
			Password string `json:"password"`
		}
		type SambaUsers struct {
			AllEntries []SambaUser `json:"all_entries"`
		}
		type SambaConfig struct {
			Users SambaUsers `json:"users"`
		}

		var config SambaConfig
		if err := json.Unmarshal(existingData, &config); err == nil {
			for _, u := range config.Users.AllEntries {
				existingUsers = append(existingUsers, SMBUser{
					Username: u.Name,
					Password: u.Password,
				})
			}
		}
	}

	// Append new users
	allUsers := existingUsers
	for _, newUser := range newUsers {
		found := false
		for i, existingUser := range allUsers {
			if existingUser.Username == newUser.Username {
				allUsers[i].Password = newUser.Password
				found = true
				break
			}
		}
		if !found {
			allUsers = append(allUsers, newUser)
		}
	}

	secretData := sambaConfigPrefix

	for i, user := range allUsers {
		if i > 0 {
			secretData += ","
		}
		secretData += fmt.Sprintf(sambaUserEntryTemplate, user.Username, user.Password)
	}
	secretData += sambaConfigSuffix

	data := map[string]string{
		"users": secretData,
	}
	_, err = uc.kubeCore.UpdateSecret(ctx, restConfig, namespace, secretName, data)
	return err
}

func (uc *StorageUseCase) updateJoinSecret(ctx context.Context, restConfig *rest.Config, namespace, secretName string, adAuth *ADAuth) error {
	joinSecretData := fmt.Sprintf(`{"username": %q, "password": %q}`, adAuth.JoinUsername, adAuth.JoinPassword)
	data := map[string]string{
		"join": joinSecretData,
	}
	_, err := uc.kubeCore.UpdateSecret(ctx, restConfig, namespace, secretName, data)
	return err
}

func (uc *StorageUseCase) restartSMBDeployment(ctx context.Context, restConfig *rest.Config, namespace, name string) error {
	deployment, err := uc.kubeApps.GetDeployment(ctx, restConfig, namespace, name)
	if err != nil {
		return fmt.Errorf("failed to get SMB share deployment for restart: %w", err)
	}

	// Trigger rolling update
	if deployment.Spec.Template.Annotations == nil {
		deployment.Spec.Template.Annotations = map[string]string{}
	}
	deployment.Spec.Template.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)

	_, err = uc.kubeApps.UpdateDeployment(ctx, restConfig, namespace, deployment)
	return err
}

func (uc *StorageUseCase) UpdateSMBShare(ctx context.Context, scope, facility, namespace, name string, config *SMBShareConfig, sizeBytes uint64, localAuth *LocalAuth, adAuth *ADAuth) (*SMBShare, error) {
	restConfig, err := kubeConfig(ctx, uc.facility, uc.action, scope, facility)
	if err != nil {
		return nil, err
	}

	commonConfigName := name + "-common"
	securityConfigName := name + "-security"
	userSecretName := name + "-users"
	joinSecretName := name + "-join"

	if err := uc.kubeSMB.UpdateSMBShare(ctx, restConfig, namespace, name, sizeBytes, config.Browseable, config.ReadOnly, config.GuestOk, config.ValidUsers); err != nil {
		return nil, fmt.Errorf("failed to update SMB share: %w", err)
	}

	if err := uc.kubeSMB.UpdateSMBCommonConfig(ctx, restConfig, namespace, commonConfigName, config); err != nil {
		return nil, fmt.Errorf("failed to update common config: %w", err)
	}

	restartNeeded := false
	if localAuth != nil && len(localAuth.Users) > 0 {
		if err := uc.mergeAndUpdateUserSecret(ctx, restConfig, namespace, userSecretName, localAuth.Users); err != nil {
			return nil, fmt.Errorf("failed to update user secret: %w", err)
		}
		restartNeeded = true
	}

	if adAuth != nil && adAuth.Realm != "" {
		// Update AD join secret
		if adAuth.JoinUsername != "" && adAuth.JoinPassword != "" {
			if err := uc.updateJoinSecret(ctx, restConfig, namespace, joinSecretName, adAuth); err != nil {
				return nil, fmt.Errorf("failed to update AD join secret: %w", err)
			}
			restartNeeded = true
		}

		if err := uc.kubeSMB.UpdateSMBSecurityConfig(ctx, restConfig, namespace, securityConfigName, "active-directory", adAuth.Realm, joinSecretName); err != nil {
			return nil, fmt.Errorf("failed to update security config: %w", err)
		}
	}

	if restartNeeded {
		if err := uc.restartSMBDeployment(ctx, restConfig, namespace, name); err != nil {
			return nil, fmt.Errorf("failed to restart SMB share deployment: %w", err)
		}
	}

	return uc.kubeSMB.GetSMBShare(ctx, restConfig, namespace, name)
}
