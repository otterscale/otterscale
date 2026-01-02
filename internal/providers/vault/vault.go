package vault

import (
	"fmt"
	"os"

	"github.com/hashicorp/vault/api"

	"github.com/otterscale/otterscale/internal/config"
)

type Vault struct {
	conf *config.Config

	client  *api.Client
	kvMount string
}

func New(conf *config.Config) (*Vault, error) {
	vaultConfig := api.DefaultConfig()
	vaultConfig.Address = conf.VaultAddress()

	client, err := api.NewClient(vaultConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault client: %w", err)
	}

	if err := authenticate(client, conf); err != nil {
		return nil, fmt.Errorf("failed to authenticate with vault: %w", err)
	}

	return &Vault{
		conf:    conf,
		client:  client,
		kvMount: conf.VaultKVMount(),
	}, nil
}

func authenticate(client *api.Client, conf *config.Config) error {
	// Token authentication
	if token := conf.VaultToken(); token != "" {
		client.SetToken(token)
		return nil
	}

	// Kubernetes authentication (for running inside a pod)
	if role := conf.VaultKubernetesRole(); role != "" {
		return authenticateKubernetes(client, conf)
	}

	return fmt.Errorf("no vault authentication method configured")
}

func authenticateKubernetes(client *api.Client, conf *config.Config) error {
	jwt, err := readServiceAccountToken()
	if err != nil {
		return err
	}

	loginData := map[string]any{
		"jwt":  jwt,
		"role": conf.VaultKubernetesRole(),
	}

	authPath := conf.VaultKubernetesAuthPath()
	secret, err := client.Logical().Write(fmt.Sprintf("auth/%s/login", authPath), loginData)
	if err != nil {
		return err
	}

	client.SetToken(secret.Auth.ClientToken)
	return nil
}

func readServiceAccountToken() (string, error) {
	tokenPath := "/var/run/secrets/kubernetes.io/serviceaccount/token"
	token, err := os.ReadFile(tokenPath)
	if err != nil {
		return "", fmt.Errorf("failed to read service account token: %w", err)
	}
	return string(token), nil
}

func (m *Vault) Client() *api.Client {
	return m.client
}

func (m *Vault) KVMount() string {
	return m.kvMount
}
