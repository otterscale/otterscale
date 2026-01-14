package vault

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/otterscale/otterscale/internal/core/secret"
)

type secretRepo struct {
	vault *Vault
}

func NewSecretRepo(vault *Vault) secret.Repository {
	return &secretRepo{
		vault: vault,
	}
}

var _ secret.Repository = (*secretRepo)(nil)

func (r *secretRepo) Get(ctx context.Context, scope string, secretType secret.SecretType, name string) (*secret.Secret, error) {
	path := r.buildPath(scope, secretType, name)

	vaultSecret, err := r.vault.client.KVv2(r.vault.kvMount).Get(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to get secret from vault: %w", err)
	}

	if vaultSecret == nil || vaultSecret.Data == nil {
		return nil, fmt.Errorf("secret not found: %s/%s/%s", scope, secretType, name)
	}

	return r.toSecret(scope, secretType, name, vaultSecret.Data, vaultSecret.VersionMetadata.Version)
}

func (r *secretRepo) Store(ctx context.Context, s *secret.Secret) error {
	path := r.buildPath(s.Scope, s.Type, s.Name)

	// Convert binary data to base64
	data := make(map[string]any)
	for k, v := range s.Data {
		data[k] = base64.StdEncoding.EncodeToString(v)
	}

	// Add metadata
	data["_metadata"] = s.Metadata
	data["_updated_at"] = time.Now().UTC().Format(time.RFC3339)

	_, err := r.vault.client.KVv2(r.vault.kvMount).Put(ctx, path, data)
	if err != nil {
		return fmt.Errorf("failed to store secret in vault: %w", err)
	}

	return nil
}

func (r *secretRepo) Delete(ctx context.Context, scope string, secretType secret.SecretType, name string) error {
	path := r.buildPath(scope, secretType, name)

	err := r.vault.client.KVv2(r.vault.kvMount).DeleteMetadata(ctx, path)
	if err != nil {
		return fmt.Errorf("failed to delete secret from vault: %w", err)
	}

	return nil
}

func (r *secretRepo) List(ctx context.Context, scope string, secretType secret.SecretType) ([]*secret.Secret, error) {
	path := fmt.Sprintf("scopes/%s/%s", scope, secretType)

	list, err := r.vault.client.Logical().ListWithContext(ctx, fmt.Sprintf("%s/metadata/%s", r.vault.kvMount, path))
	if err != nil {
		return nil, fmt.Errorf("failed to list secrets from vault: %w", err)
	}

	if list == nil || list.Data == nil {
		return []*secret.Secret{}, nil
	}

	keys, ok := list.Data["keys"].([]any)
	if !ok {
		return []*secret.Secret{}, nil
	}

	secrets := make([]*secret.Secret, 0, len(keys))
	for _, key := range keys {
		name := key.(string)
		s, err := r.Get(ctx, scope, secretType, name)
		if err != nil {
			continue // skip errors
		}
		secrets = append(secrets, s)
	}

	return secrets, nil
}

func (r *secretRepo) Exists(ctx context.Context, scope string, secretType secret.SecretType, name string) (bool, error) {
	_, err := r.Get(ctx, scope, secretType, name)
	if err != nil {
		// Only return false for "not found" errors; propagate other errors
		if strings.Contains(err.Error(), "secret not found") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *secretRepo) buildPath(scope string, secretType secret.SecretType, name string) string {
	return fmt.Sprintf("scopes/%s/%s/%s", scope, secretType, name)
}

func (r *secretRepo) toSecret(scope string, secretType secret.SecretType, name string, data map[string]any, version int) (*secret.Secret, error) {
	s := &secret.Secret{
		Scope:   scope,
		Type:    secretType,
		Name:    name,
		Data:    make(map[string][]byte),
		Version: version,
	}

	for k, v := range data {
		if k == "_metadata" {
			if meta, ok := v.(map[string]any); ok {
				s.Metadata = make(map[string]string)
				for mk, mv := range meta {
					s.Metadata[mk] = fmt.Sprint(mv)
				}
			}
			continue
		}
		if k == "_updated_at" {
			if t, err := time.Parse(time.RFC3339, v.(string)); err == nil {
				s.UpdatedAt = t
			}
			continue
		}

		// Decode base64 data
		if str, ok := v.(string); ok {
			decoded, err := base64.StdEncoding.DecodeString(str)
			if err != nil {
				s.Data[k] = []byte(str) // fallback to raw string
			} else {
				s.Data[k] = decoded
			}
		}
	}

	return s, nil
}
