package kube

import (
	"crypto/sha256"
	"log"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/registry"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type KubeMap struct {
	*sync.Map
}

func NewKubeMap() (KubeMap, error) {
	return KubeMap{&sync.Map{}}, nil
}

func (m KubeMap) exists(uuid, facility string) bool {
	key := key(uuid, facility)
	_, ok := m.Load(key)
	return ok
}

func (m KubeMap) get(uuid, facility string) (*kubernetes.Clientset, error) {
	key := key(uuid, facility)
	if v, ok := m.Load(key); ok {
		return v.(*kubernetes.Clientset), nil
	}
	return nil, status.Errorf(codes.NotFound, "kubernetes %q in scope %q not found", facility, uuid)
}

func (m KubeMap) set(uuid, facility string, config *rest.Config) error {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	key := key(uuid, facility)
	m.Store(key, clientset)
	return nil
}

type HelmMap struct {
	*sync.Map
}

func NewHelmMap() (HelmMap, error) {
	return HelmMap{&sync.Map{}}, nil
}

func (m HelmMap) exists(uuid, facility string) bool {
	key := key(uuid, facility)
	_, ok := m.Load(key)
	return ok
}

func (m HelmMap) get(uuid, facility, namespace string, rc *registry.Client) (*action.Configuration, error) {
	key := key(uuid, facility)
	if v, ok := m.Load(key); ok {
		restConfig := v.(*rest.Config)
		getter := genericclioptions.NewConfigFlags(true)
		getter.APIServer = &restConfig.Host
		getter.BearerToken = &restConfig.BearerToken
		getter.CAFile = &restConfig.CAFile
		getter.CertFile = &restConfig.TLSClientConfig.CertFile
		getter.KeyFile = &restConfig.TLSClientConfig.KeyFile
		getter.Insecure = &restConfig.Insecure
		getter.Namespace = &namespace

		config := new(action.Configuration)
		if err := config.Init(getter, namespace, "", log.Printf); err != nil {
			return nil, err
		}
		config.RegistryClient = rc
		return config, nil
	}
	return nil, status.Errorf(codes.NotFound, "helm with kubernetes %q in scope %q not found", facility, uuid)
}

func (m HelmMap) set(uuid, facility string, config *rest.Config) error {
	key := key(uuid, facility)
	m.Store(key, config)
	return nil
}

func key(uuid, facility string) string {
	sha := sha256.New()
	sha.Write([]byte(uuid))
	sha.Write([]byte(facility))
	return string(sha.Sum(nil))
}
