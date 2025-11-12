package ceph

import (
	"context"
	"encoding/json"
	"errors"
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/ceph/go-ceph/rados"
	"github.com/ceph/go-ceph/rgw/admin"
	"golang.org/x/sync/errgroup"
	"gopkg.in/ini.v1"

	"github.com/otterscale/otterscale/internal/config"
	"github.com/otterscale/otterscale/internal/providers/juju"
)

type Ceph struct {
	conf *config.Config
	juju *juju.Juju

	connections sync.Map
	clients     sync.Map
}

type connectionConfig struct {
	id   string
	host string
	key  string
}

type clientConfig struct {
	endpoint  string
	accessKey string
	secretKey string
}

func New(conf *config.Config, juju *juju.Juju) *Ceph {
	return &Ceph{
		conf: conf,
		juju: juju,
	}
}

func (m *Ceph) extractAsConnectionConfig(r map[string]any) (connectionConfig, error) {
	stdout, ok := r["stdout"]
	if !ok {
		return connectionConfig{}, errors.New("ceph config stdout not found")
	}

	file, err := ini.Load([]byte(stdout.(string)))
	if err != nil {
		return connectionConfig{}, err
	}

	id := file.Section("global").Key("fsid").String()
	if id == "" {
		return connectionConfig{}, errors.New("ceph config fsid not found")
	}

	host := file.Section("global").Key("mon_host").String()
	if host == "" {
		return connectionConfig{}, errors.New("ceph config mon_host not found")
	}

	key := file.Section("client.admin").Key("key").String()
	if key == "" {
		return connectionConfig{}, errors.New("ceph config key not found")
	}

	return connectionConfig{
		id:   id,
		host: host,
		key:  key,
	}, nil
}

func (m *Ceph) getConnectionConfig(scope string) (connectionConfig, error) {
	name := scope + "-ceph-mon"

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	result, err := m.juju.Execute(ctx, scope, name, "ceph config generate-minimal-conf && ceph auth get client.admin")
	if err != nil {
		return connectionConfig{}, err
	}

	return m.extractAsConnectionConfig(result)
}

func (m *Ceph) newConnection(c connectionConfig) (*rados.Conn, error) {
	conn, err := rados.NewConn()
	if err != nil {
		return nil, err
	}

	if err := conn.SetConfigOption("fsid", c.id); err != nil {
		return nil, err
	}
	if err := conn.SetConfigOption("mon_host", c.host); err != nil {
		return nil, err
	}
	if err := conn.SetConfigOption("key", c.key); err != nil {
		return nil, err
	}

	radosTimeout := time.Second * 3
	if m.conf.Ceph.RADOSTimeout > 0 {
		radosTimeout = m.conf.Ceph.RADOSTimeout
	}
	timeout := strconv.FormatFloat(radosTimeout.Seconds(), 'f', -1, 64)

	if err = conn.SetConfigOption("rados_mon_op_timeout", timeout); err != nil {
		return nil, err
	}
	if err = conn.SetConfigOption("rados_osd_op_timeout", timeout); err != nil {
		return nil, err
	}

	if err := conn.Connect(); err != nil {
		return nil, err
	}
	return conn, nil
}

// TODO: Check connection with 'OpenIOContext'
func (m *Ceph) Connection(scope string) (*rados.Conn, error) {
	if v, ok := m.connections.Load(scope); ok {
		conn := v.(*rados.Conn)
		return conn, nil
	}

	config, err := m.getConnectionConfig(scope)
	if err != nil {
		return nil, err
	}

	conn, err := m.newConnection(config)
	if err != nil {
		return nil, err
	}

	m.connections.Store(scope, conn)

	return conn, nil
}

func (m *Ceph) getRGWCommand(r map[string]any) (string, error) {
	stdout, ok := r["stdout"]
	if !ok {
		return "", errors.New("rgw list config stdout not found")
	}

	var users []string
	if err := json.Unmarshal([]byte(stdout.(string)), &users); err != nil {
		return "", err
	}

	if slices.Contains(users, "otterscale") {
		return "radosgw-admin user info --uid=otterscale --format=json", nil
	}
	return "radosgw-admin user create --system --uid=otterscale --display-name=OtterScale --format json", nil
}

func (m *Ceph) extractObjectKeys(r map[string]any) (accessKey, secretKey string, err error) {
	stdout, ok := r["stdout"]
	if !ok {
		return "", "", errors.New("rgw config stdout not found")
	}

	type Info struct {
		Keys []struct {
			AccessKey string `json:"access_key,omitempty"`
			SecretKey string `json:"secret_key,omitempty"`
		} `json:"keys,omitempty"`
	}
	var info Info

	if err := json.Unmarshal([]byte(stdout.(string)), &info); err != nil {
		return "", "", err
	}

	if len(info.Keys) == 0 {
		return "", "", errors.New("rgw config key not found")
	}

	accessKey = info.Keys[0].AccessKey
	if accessKey == "" {
		return "", "", errors.New("rgw config access key not found")
	}

	secretKey = info.Keys[0].SecretKey
	if secretKey == "" {
		return "", "", errors.New("rgw config secret key not found")
	}

	return accessKey, secretKey, nil
}

func (m *Ceph) getObjectKeys(ctx context.Context, scope, name string) (accessKey, secretKey string, err error) {
	result, err := m.juju.Execute(ctx, scope, name, "radosgw-admin user list")
	if err != nil {
		return "", "", err
	}

	cmd, err := m.getRGWCommand(result)
	if err != nil {
		return "", "", err
	}

	result, err = m.juju.Execute(ctx, scope, name, cmd)
	if err != nil {
		return "", "", err
	}

	return m.extractObjectKeys(result)
}

func (m *Ceph) getClientConfig(scope string) (clientConfig, error) {
	name := scope + "-ceph-mon"
	rgwName := scope + "-ceph-radosgw"

	config := clientConfig{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	eg, egctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		accessKey, secretKey, err := m.getObjectKeys(egctx, scope, name)
		if err != nil {
			return err
		}
		config.accessKey = accessKey
		config.secretKey = secretKey
		return nil
	})

	eg.Go(func() error {
		endpoint, err := m.juju.GetEndpoint(egctx, scope, rgwName)
		if err != nil {
			return err
		}
		config.endpoint = endpoint
		return nil
	})

	if err := eg.Wait(); err != nil {
		return clientConfig{}, err
	}

	return config, nil
}

func (m *Ceph) Client(scope string) (*admin.API, error) {
	if v, ok := m.clients.Load(scope); ok {
		client := v.(*admin.API)
		return client, nil
	}

	config, err := m.getClientConfig(scope)
	if err != nil {
		return nil, err
	}

	client, err := admin.New(config.endpoint, config.accessKey, config.secretKey, nil)
	if err != nil {
		return nil, err
	}

	m.clients.Store(scope, client)

	return client, nil
}
