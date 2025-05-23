package config

import (
	"context"
	"log"
	"os"

	"buf.build/go/protoyaml"
	"github.com/fsnotify/fsnotify"
)

type Config struct {
	ConfigSet *ConfigSet

	path    string
	watcher *fsnotify.Watcher

	ctx    context.Context
	cancel context.CancelFunc
}

func New(path string) (*Config, error) {
	if err := initFile(path); err != nil {
		return nil, err
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	if err := watcher.Add(path); err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &Config{
		ConfigSet: &ConfigSet{},
		path:      path,
		watcher:   watcher,
		ctx:       ctx,
		cancel:    cancel,
	}, nil
}

func (c *Config) Load() error {
	if err := c.scan(); err != nil {
		return err
	}
	if err := c.watch(); err != nil {
		return err
	}
	return nil
}

func (c *Config) Close() error {
	c.cancel()
	return c.watcher.Close()
}

func (c *Config) watch() error {
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
						log.Println("config renamed:", event.Name)
						if err := c.watcher.Add(event.Name); err != nil {
							return
						}
					}
				}
				if event.Has(fsnotify.Write) {
					log.Println("config modified:", event.Name)
				}
				if err := c.scan(); err != nil {
					log.Println("config error:", err)
				}
			case err, ok := <-c.watcher.Errors:
				if !ok {
					return
				}
				log.Println("config error:", err)
			}
		}
	}()
	return nil
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

func initFile(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()
	}
	return nil
}
