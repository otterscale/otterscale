package smb

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"

	"github.com/otterscale/otterscale/internal/core/application/config"
	"github.com/otterscale/otterscale/internal/core/application/persistent"
	"github.com/otterscale/otterscale/internal/core/application/service"
)

func mergeUsers(existing, desired []userEntry) []User {
	userMap := make(map[string]userEntry)

	// Add existing users to map
	for _, u := range existing {
		userMap[u.Name] = u
	}

	// Update or add desired users
	for _, u := range desired {
		if existing, found := userMap[u.Name]; found {
			// Keep existing password if new one is empty
			if u.Password != "" {
				userMap[u.Name] = u
			} else {
				userMap[u.Name] = existing
			}
		} else {
			// Add new user
			userMap[u.Name] = u
		}
	}

	// Remove users not in desired list
	maps.DeleteFunc(userMap, func(k string, _ userEntry) bool {
		for _, u := range desired {
			if u.Name == k {
				return false
			}
		}
		return true
	})

	// Convert map back to slice
	result := make([]User, 0, len(userMap))
	for _, u := range userMap {
		result = append(result, User{
			Username: u.Name,
			Password: u.Password,
		})
	}

	return result
}

func (uc *UseCase) upsertUsersSecret(ctx context.Context, scope, namespace, name string, secret *config.Secret) error {
	existing, err := uc.secret.Get(ctx, scope, namespace, name)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			_, err = uc.secret.Create(ctx, scope, namespace, secret)
			return err
		}
		return err
	}

	mergedUsers, err := uc.mergeUsersFromSecrets(existing, secret)
	if err != nil {
		return err
	}

	newSecret := uc.buildUsersSecret(namespace, name, mergedUsers)
	newSecret.ObjectMeta = existing.ObjectMeta

	_, err = uc.secret.Update(ctx, scope, namespace, newSecret)
	return err
}

func (uc *UseCase) mergeUsersFromSecrets(existing, desired *config.Secret) ([]User, error) {
	existingUsers, err := uc.extractUserEntries(existing)
	if err != nil {
		return nil, fmt.Errorf("extract existing users: %w", err)
	}

	desiredUsers, err := uc.extractUserEntries(desired)
	if err != nil {
		return nil, fmt.Errorf("extract desired users: %w", err)
	}

	return mergeUsers(existingUsers, desiredUsers), nil
}

func (uc *UseCase) extractUserEntries(secret *config.Secret) ([]userEntry, error) {
	usersData, ok := secret.Data["users"]
	if !ok {
		return nil, fmt.Errorf("users key not found in secret data")
	}

	var config sambaContainerConfig
	if err := json.Unmarshal(usersData, &config); err != nil {
		return nil, err
	}

	return config.Users["all_entries"], nil
}

func (uc *UseCase) extractJoinSource(secret *config.Secret) (*User, error) {
	joinData, ok := secret.Data["join"]
	if !ok {
		return nil, fmt.Errorf("join key not found in secret data")
	}

	var user User
	if err := json.Unmarshal(joinData, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (uc *UseCase) upsertJoinSecret(ctx context.Context, scope, namespace, name string, secret *config.Secret) error {
	existing, err := uc.secret.Get(ctx, scope, namespace, name)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			_, err = uc.secret.Create(ctx, scope, namespace, secret)
			return err
		}
		return err
	}

	joinSource, err := uc.extractJoinSource(existing)
	if err != nil {
		return err
	}

	// user password is not updated if empty
	if joinSource.Password == "" {
		return nil
	}

	secret.ObjectMeta = existing.ObjectMeta

	_, err = uc.secret.Update(ctx, scope, namespace, secret)
	return err
}

func (uc *UseCase) updateSecurityConfig(ctx context.Context, scope, namespace, name string, securityConfig *SecurityConfig) error {
	existing, err := uc.smbSecurityConfig.Get(ctx, scope, namespace, name)
	if err != nil {
		return err
	}

	securityConfig.ObjectMeta = existing.ObjectMeta

	_, err = uc.smbSecurityConfig.Update(ctx, scope, namespace, securityConfig)
	return err
}

func (uc *UseCase) updateCommonConfig(ctx context.Context, scope, namespace, name string, commonConfig *CommonConfig) error {
	existing, err := uc.smbCommonConfig.Get(ctx, scope, namespace, name)
	if err != nil {
		return err
	}

	commonConfig.ObjectMeta = existing.ObjectMeta

	_, err = uc.smbCommonConfig.Update(ctx, scope, namespace, commonConfig)
	return err
}

func (uc *UseCase) updateShare(ctx context.Context, scope, namespace, name string, share *Share) error {
	existing, err := uc.smbShare.Get(ctx, scope, namespace, name)
	if err != nil {
		return err
	}

	share.ObjectMeta = existing.ObjectMeta

	_, err = uc.smbShare.Update(ctx, scope, namespace, share)
	return err
}

func (uc *UseCase) updatePersistentVolumeClaim(ctx context.Context, scope, namespace string, pvc *persistent.PersistentVolumeClaim) error {
	existing, err := uc.persistentVolumeClaim.Get(ctx, scope, namespace, pvc.Name)
	if err != nil {
		return err
	}

	pvc.ObjectMeta = existing.ObjectMeta

	_, err = uc.persistentVolumeClaim.Update(ctx, scope, namespace, pvc)
	return err
}

func (uc *UseCase) updateService(ctx context.Context, scope, namespace, name string, service *service.Service) error {
	existing, err := uc.service.Get(ctx, scope, namespace, name)
	if err != nil {
		return err
	}

	service.ObjectMeta = existing.ObjectMeta

	_, err = uc.service.Update(ctx, scope, namespace, service)
	return err
}
