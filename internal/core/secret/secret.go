package secret

import (
	"context"
	"time"
)

// SecretType defines the type of secret.
type SecretType string

const (
	SecretTypeKubeConfig SecretType = "kubeconfig"
	SecretTypeTLSCert    SecretType = "tls-cert"
	SecretTypeCredential SecretType = "credential"
	SecretTypeAPIKey     SecretType = "api-key"
	SecretTypeGeneric    SecretType = "generic"
)

// Secret represents a secret stored in a secure vault.
type Secret struct {
	ID        string
	Scope     string
	Type      SecretType
	Name      string
	Data      map[string][]byte
	Metadata  map[string]string
	Version   int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Repository defines the interface for secret storage operations.
type Repository interface {
	Get(ctx context.Context, scope string, secretType SecretType, name string) (*Secret, error)
	Store(ctx context.Context, secret *Secret) error
	Delete(ctx context.Context, scope string, secretType SecretType, name string) error
	List(ctx context.Context, scope string, secretType SecretType) ([]*Secret, error)
	Exists(ctx context.Context, scope string, secretType SecretType, name string) (bool, error)
}
