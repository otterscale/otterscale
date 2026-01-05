package vault

import (
	"fmt"
	"os"

	"github.com/hashicorp/vault/api"
)

type Vault struct {
	client  *api.Client
	kvMount string
}

func New() (*Vault, error) {
	vaultConfig := api.DefaultConfig()

	client, err := api.NewClient(vaultConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault client: %w", err)
	}

	kvMount := os.Getenv("VAULT_KV_MOUNT")
	if kvMount == "" {
		kvMount = "secret"
	}

	return &Vault{
		client:  client,
		kvMount: kvMount,
	}, nil
}

func (m *Vault) Client() *api.Client {
	return m.client
}

func (m *Vault) KVMount() string {
	return m.kvMount
}
