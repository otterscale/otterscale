package config

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"sigs.k8s.io/yaml"
)

type MAAS struct {
	URL     string `json:"url"`
	Key     string `json:"key"`
	Version string `json:"version"`
}

type Juju struct {
	Controller          string   `json:"controller"`
	ControllerAddresses []string `json:"controller_addresses"`
	Username            string   `json:"username"`
	Password            string   `json:"password"`
	CACert              string   `json:"ca_cert"`
	CloudName           string   `json:"cloud_name"`
	CloudRegion         string   `json:"cloud_region"`
	CharmhubAPIURL      string   `json:"charmhub_api_url"`
}

type MicroK8s struct {
	Config string `json:"config"`
}

type Ceph struct {
	RADOSTimeout time.Duration `json:"rados_timeout"`
}

type Vault struct {
	Address            string `json:"address"`
	Token              string `json:"token"`
	KVMount            string `json:"kv_mount"`
	KubernetesRole     string `json:"kubernetes_role"`
	KubernetesAuthPath string `json:"kubernetes_auth_path"`
}

type Schema struct {
	// System
	MAAS     MAAS     `json:"maas"`
	Juju     Juju     `json:"juju"`
	MicroK8s MicroK8s `json:"micro_k8s"`
	Vault    Vault    `json:"vault"`

	// User
	Ceph Ceph `json:"ceph"`
}

type Config struct {
	v *viper.Viper
}

func New() *Config {
	return &Config{
		v: viper.New(),
	}
}

func (c *Config) Load(path string) error {
	extension := filepath.Ext(path)
	if len(extension) < 2 {
		return fmt.Errorf("extension not found in filename: %q", path)
	}

	filename := filepath.Base(path)
	filenameOnly := filename[0 : len(filename)-len(extension)]

	c.v.AddConfigPath(filepath.Dir(path))
	c.v.SetConfigName(filenameOnly)
	c.v.SetConfigType(extension[1:]) // remove dot

	if err := c.v.ReadInConfig(); err != nil {
		return err
	}

	c.v.WatchConfig()

	c.v.OnConfigChange(func(in fsnotify.Event) {
		slog.Info("configuration file changed", "file", in.Name)
	})

	return nil
}

func (c *Config) MAASURL() string {
	return c.v.GetString("maas.url")
}

func (c *Config) MAASKey() string {
	return c.v.GetString("maas.key")
}

func (c *Config) MAASVersion() string {
	return c.v.GetString("maas.version")
}

func (c *Config) JujuController() string {
	return c.v.GetString("juju.controller")
}

func (c *Config) JujuControllerAddresses() []string {
	return c.v.GetStringSlice("juju.controller_addresses")
}

func (c *Config) JujuUsername() string {
	return c.v.GetString("juju.username")
}

func (c *Config) JujuPassword() string {
	return c.v.GetString("juju.password")
}

func (c *Config) JujuCACert() string {
	return c.v.GetString("juju.ca_cert")
}

func (c *Config) JujuCloudName() string {
	return c.v.GetString("juju.cloud_name")
}

func (c *Config) JujuCloudRegion() string {
	return c.v.GetString("juju.cloud_region")
}

func (c *Config) JujuCharmhubAPIURL() string {
	return c.v.GetString("juju.charmhub_api_url")
}

func (c *Config) MicroK8sConfig() string {
	return c.v.GetString("micro_k8s.config")
}

func (c *Config) CephRADOSTimeout() time.Duration {
	return c.v.GetDuration("ceph.rados_timeout")
}

func (c *Config) VaultAddress() string {
	return c.v.GetString("vault.address")
}

func (c *Config) VaultToken() string {
	return c.v.GetString("vault.token")
}

func (c *Config) VaultKVMount() string {
	if mount := c.v.GetString("vault.kv_mount"); mount != "" {
		return mount
	}
	return "secret"
}

func (c *Config) VaultKubernetesRole() string {
	return c.v.GetString("vault.kubernetes_role")
}

func (c *Config) VaultKubernetesAuthPath() string {
	if path := c.v.GetString("vault.kubernetes_auth_path"); path != "" {
		return path
	}
	return "kubernetes"
}

func PrintDefault() error {
	data, err := yaml.Marshal(defaultSchema())
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

func defaultSchema() *Schema {
	return &Schema{
		MAAS: MAAS{
			URL:     "http://localhost:5240/MAAS/",
			Key:     "::",
			Version: "2.0",
		},
		Juju: Juju{
			ControllerAddresses: []string{"localhost:17070"},
			Username:            "admin",
			Password:            "password",
			CACert: `-----BEGIN CERTIFICATE-----
...
-----END CERTIFICATE-----`,
			CloudName:      "maas-cloud",
			CloudRegion:    "default",
			CharmhubAPIURL: "https://api.charmhub.io",
		},
	}
}
