package kube

import (
	"crypto/sha256"
	"os"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/registry"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type KubeMap map[string]*kubernetes.Clientset

func NewKubeMap() (KubeMap, error) {
	return map[string]*kubernetes.Clientset{}, nil
}

func (m KubeMap) exists(uuid, facility string) bool {
	key := key(uuid, facility)
	_, ok := m[key]
	return ok
}

func (m KubeMap) get(uuid, facility string) (*kubernetes.Clientset, error) {
	key := key(uuid, facility)
	if clientset, ok := m[key]; ok {
		return clientset, nil
	}
	return nil, status.Errorf(codes.NotFound, "kubernetes %q in scope %q not found", facility, uuid)
}

func (m KubeMap) set(uuid, facility string, config *rest.Config) error {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	key := key(uuid, facility)
	m[key] = clientset
	return nil
}

type HelmMap map[string]*genericclioptions.ConfigFlags

func NewHelmMap() (HelmMap, error) {
	return map[string]*genericclioptions.ConfigFlags{}, nil
}

func (m HelmMap) exists(uuid, facility string) bool {
	key := key(uuid, facility)
	_, ok := m[key]
	return ok
}

func (m HelmMap) get(uuid, facility, namespace string, rc *registry.Client) (*action.Configuration, error) {
	key := key(uuid, facility)
	if getter, ok := m[key]; ok {
		getter.Namespace = &namespace
		config := new(action.Configuration)
		if err := config.Init(getter, namespace, os.Getenv("HELM_DRIVER"), nil); err != nil {
			return nil, err
		}
		config.RegistryClient = rc
		return config, nil
	}
	return nil, status.Errorf(codes.NotFound, "helm with kubernetes %q in scope %q not found", facility, uuid)
}

func (m HelmMap) set(uuid, facility string, config *rest.Config) error {
	getter := genericclioptions.NewConfigFlags(true)
	getter.APIServer = &config.Host
	getter.BearerToken = &config.BearerToken
	getter.CAFile = &config.CAFile
	getter.CertFile = &config.TLSClientConfig.CertFile
	getter.KeyFile = &config.TLSClientConfig.KeyFile
	getter.Insecure = &config.Insecure

	key := key(uuid, facility)
	m[key] = getter
	return nil
}

func key(uuid, facility string) string {
	sha := sha256.New()
	sha.Write([]byte(uuid))
	sha.Write([]byte(facility))
	return string(sha.Sum(nil))
}
