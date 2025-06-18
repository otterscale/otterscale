package config

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/goccy/go-yaml"
)

const filePerm600 os.FileMode = 0o600

type MAAS struct {
	URL     string `yaml:"url"`
	Key     string `yaml:"key"`
	Version string `yaml:"version"`
}

type Juju struct {
	Controller          string   `yaml:"controller"`
	ControllerAddresses []string `yaml:"controller_addresses"`
	Username            string   `yaml:"username"`
	Password            string   `yaml:"password"`
	CACert              string   `yaml:"ca_cert"`
	CloudName           string   `yaml:"cloud_name"`
	CloudRegion         string   `yaml:"cloud_region"`
	CharmhubAPIURL      string   `yaml:"charmhub_api_url"`
}

type Kube struct {
	HelmRepositoryURLs []string `yaml:"helm_repository_urls"`
}

type Ceph struct {
	RadosTimeout time.Duration `yaml:"rados_timeout"`
}

type Config struct {
	MAAS MAAS `yaml:"maas"`
	Juju Juju `yaml:"juju"`
	Kube Kube `yaml:"kube"`
	Ceph Ceph `yaml:"ceph"`

	watcher *fsnotify.Watcher
	path    string
}

func New() (*Config, func(), error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, nil, err
	}
	conf := &Config{
		watcher: watcher,
	}
	return conf, func() { conf.Close() }, nil
}

func (c *Config) Load(path string) error {
	c.path = path
	if err := c.scan(); err != nil {
		return err
	}
	if err := c.watch(); err != nil {
		return err
	}
	return nil
}

func (c *Config) Close() error {
	return c.watcher.Close()
}

func (c *Config) Override(config *Config) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	return os.WriteFile(c.path, data, filePerm600)
}

func (c *Config) scan() error {
	data, err := os.ReadFile(c.path)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(data, c); err != nil {
		return err
	}
	return nil
}

func (c *Config) watch() error {
	if err := c.watcher.Add(c.path); err != nil {
		return err
	}
	go func() {
		for {
			select {
			case event, ok := <-c.watcher.Events:
				if !ok {
					return
				}
				if event.Op == fsnotify.Rename {
					_, err := os.Stat(event.Name)
					if err == nil || os.IsExist(err) {
						slog.Debug("config event renamed", "name", event.Name)
						if err := c.watcher.Add(event.Name); err != nil {
							return
						}
					}
				}
				if event.Has(fsnotify.Write) {
					slog.Debug("config event modified", "name", event.Name)
				}
				if err := c.scan(); err != nil {
					slog.Error("config event error", "name", event.Name, "err", err)
				}
			case err, ok := <-c.watcher.Errors:
				if !ok {
					return
				}
				slog.Error("config watcher error", "err", err)
			}
		}
	}()
	return nil
}

func InitFile(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return createDefaultFile(path)
	}
	return fmt.Errorf("file already exists: %s", path)
}

func defaultConfig() *Config {
	return &Config{
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
		Kube: Kube{
			HelmRepositoryURLs: []string{"http://chartmuseum:8080"},
		},
	}
}

func createDefaultFile(path string) error {
	data, err := yaml.Marshal(defaultConfig())
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, filePerm600)
}
