package config

import (
	"fmt"
	"log/slog"
	"os"

	"buf.build/go/protoyaml"
	"github.com/fsnotify/fsnotify"
)

type Config struct {
	ConfigSet *ConfigSet

	watcher *fsnotify.Watcher
	path    string
}

func New() (*Config, func(), error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, nil, err
	}
	conf := &Config{
		ConfigSet: &ConfigSet{},
		watcher:   watcher,
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

func (c *Config) Override(set *ConfigSet) error {
	data, err := protoyaml.Marshal(set)
	if err != nil {
		return err
	}
	return os.WriteFile(c.path, data, 0600)
}

func (c *Config) scan() error {
	data, err := os.ReadFile(c.path)
	if err != nil {
		return err
	}
	if err := protoyaml.Unmarshal(data, c.ConfigSet); err != nil {
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

func defaultConfigSet() *ConfigSet {
	maas := &MAAS{}
	maas.SetUrl("http://localhost:5240/MAAS/")
	maas.SetKey("::")
	maas.SetVersion("2.0")

	juju := &Juju{}
	juju.SetControllerAddresses([]string{"localhost:17070"})
	juju.SetUsername("admin")
	juju.SetPassword("password")
	juju.SetCaCert(`-----BEGIN CERTIFICATE-----
...
-----END CERTIFICATE-----`)
	juju.SetCloudName("maas-cloud")
	juju.SetCloudRegion("default")
	juju.SetCharmhubApiUrl("https://api.charmhub.io")

	kube := &Kube{}
	kube.SetHelmRepositoryUrls([]string{"http://chartmuseum:8080"})

	pb := &ConfigSet{}
	pb.SetMaas(maas)
	pb.SetJuju(juju)
	pb.SetKube(kube)

	return pb
}

func createDefaultFile(path string) error {
	data, err := protoyaml.Marshal(defaultConfigSet())
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}
