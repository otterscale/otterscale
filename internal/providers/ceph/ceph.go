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
	"github.com/otterscale/otterscale/internal/providers/kubernetes"
	"gopkg.in/ini.v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sclient "k8s.io/client-go/kubernetes"

	"github.com/otterscale/otterscale/internal/config"
)

const (
	cephNamespace          = "rook-ceph"
	cephMonSecret          = "rook-ceph-mon"
	cephMonHostSecret      = "rook-ceph-config"
	cephAdminKeyringSecret = "rook-ceph-admin-keyring"
	cephObjectUserSecret   = "rook-ceph-object-user-ceph-objectstore-otterscale"
)

type Ceph struct {
	conf       *config.Config
	kubernetes *kubernetes.Kubernetes

	connections sync.Map
	clients     sync.Map
}

type connectionConfig struct {
	ID   string
	Host string
	Key  string
}

type clientConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
}

func New(conf *config.Config, kubernetes *kubernetes.Kubernetes) *Ceph {
	return &Ceph{
		conf:       conf,
		kubernetes: kubernetes,
	}
}

func (m *Ceph) getConnectionConfig(scope string) (connectionConfig, error) {
	clientset, err := m.kubernetes.Clientset(scope)
	if err != nil {
		return connectionConfig{}, err
	}

	// Create a single context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// Retrieve fsid
	fsid, err := getSecretData(ctx, clientset, cephNamespace, cephMonSecret, "fsid")
	if err != nil {
		return connectionConfig{}, err
	}

	// Retrieve mon_host
	monHost, err := getSecretData(ctx, clientset, cephNamespace, cephMonHostSecret, "mon_host")
	if err != nil {
		return connectionConfig{}, err
	}

	// Retrieve key
	key, err := getSecretData(ctx, clientset, cephNamespace, cephAdminKeyringSecret, "keyring")
	if err != nil {
		return connectionConfig{}, err
	}

	return connectionConfig{
		ID:   fsid,
		Host: monHost,
		Key:  key,
	}, nil
}

// Helper function to parse keyring data
func parseKeyring(keyring string) (string, error) {
	file, err := ini.Load([]byte(keyring))
	if err != nil {
		return "", err
	}

	key := file.Section("client.admin").Key("key").String()
	if key == "" {
		return "", errors.New("client.admin key not found in keyring")
	}

	return key, nil
}

// Helper function to get secret data
func getSecretData(ctx context.Context, clientset *k8sclient.Clientset, namespace, secretName, key string) (string, error) {
	secret, err := clientset.CoreV1().Secrets(namespace).Get(ctx, secretName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	data, ok := secret.Data[key]
	if !ok {
		return "", errors.New("key not found in secret")
	}

	// If the key is "keyring", parse it to extract the client.admin key
	if key == "keyring" {
		return parseKeyring(string(data))
	}

	return string(data), nil
}

func (m *Ceph) newConnection(c connectionConfig) (*rados.Conn, error) {
	conn, err := rados.NewConn()
	if err != nil {
		return nil, err
	}

	if err := conn.SetConfigOption("fsid", c.ID); err != nil {
		return nil, err
	}

	if err := conn.SetConfigOption("mon_host", c.Host); err != nil {
		return nil, err
	}

	if err := conn.SetConfigOption("key", c.Key); err != nil {
		return nil, err
	}

	radosTimeout := m.conf.CephRADOSTimeout()

	if radosTimeout <= 0 {
		radosTimeout = time.Second * 3
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
func (m *Ceph) connection(scope string) (*rados.Conn, error) {
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

func (m *Ceph) getClientConfig(scope string) (clientConfig, error) {
	clientset, err := m.kubernetes.Clientset(scope)
	if err != nil {
		return clientConfig{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	secret, err := clientset.CoreV1().Secrets(cephNamespace).Get(ctx, cephObjectUserSecret, metav1.GetOptions{})
	if err != nil {
		return clientConfig{}, err
	}

	getField := func(field string) (string, error) {
		data, ok := secret.Data[field]
		if !ok {
			return "", errors.New(field + " not found in secret")
		}
		return string(data), nil
	}

	accessKey, err := getField("AccessKey")
	if err != nil {
		return clientConfig{}, err
	}
	endpoint, err := getField("Endpoint")
	if err != nil {
		return clientConfig{}, err
	}
	secretKey, err := getField("SecretKey")
	if err != nil {
		return clientConfig{}, err
	}

	return clientConfig{
		Endpoint:  endpoint,
		AccessKey: accessKey,
		SecretKey: secretKey,
	}, nil
}

func (m *Ceph) client(scope string) (*admin.API, error) {
	if v, ok := m.clients.Load(scope); ok {
		client := v.(*admin.API)
		return client, nil
	}

	config, err := m.getClientConfig(scope)
	if err != nil {
		return nil, err
	}

	client, err := admin.New(config.Endpoint, config.AccessKey, config.SecretKey, nil)
	if err != nil {
		return nil, err
	}

	m.clients.Store(scope, client)

	return client, nil
}
