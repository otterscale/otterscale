package service

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/juju/juju/api/client/application"

	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/openhdc/openhdc/internal/domain/model"
)

type KubeClient interface {
	Get(cluster string) (*kubernetes.Clientset, error)
	Add(cluster string, cfg *rest.Config) error
}

type KubeDeployment interface {
	List(ctx context.Context, cluster, namespace string) (*appv1.DeploymentList, error)
}

type KubePersistentVolumeClaim interface {
	List(ctx context.Context, cluster, namespace string) (*corev1.PersistentVolumeClaimList, error)
}

type KubePod interface {
	List(ctx context.Context, cluster, namespace string) (*corev1.PodList, error)
}

type KubeSVC interface {
	List(ctx context.Context, cluster, namespace string) (*corev1.ServiceList, error)
}

type KubeService struct {
	client                KubeClient
	cronJob               KubeCronJob
	deployment            KubeDeployment
	job                   KubeJob
	namespace             KubeNamespace
	persistentVolumeClaim KubePersistentVolumeClaim
	pod                   KubePod
	service               KubeSVC
	application           JujuApplication
}

func NewKubeService(
	client KubeClient,
	cronJob KubeCronJob,
	deployment KubeDeployment,
	job KubeJob,
	namespace KubeNamespace,
	persistentVolumeClaim KubePersistentVolumeClaim,
	pod KubePod,
	service KubeSVC,
	application JujuApplication,
) *KubeService {
	return &KubeService{
		client:                client,
		cronJob:               cronJob,
		deployment:            deployment,
		job:                   job,
		namespace:             namespace,
		persistentVolumeClaim: persistentVolumeClaim,
		pod:                   pod,
		service:               service,
		application:           application,
	}
}

func (s *KubeService) ensureClient(ctx context.Context, uuid, cluster string) error {
	// Check if client already exists
	if _, err := s.client.Get(cluster); err == nil {
		return nil // Client already exists
	}

	// Create new client config
	cfg, err := s.newConfig(ctx, uuid, cluster)
	if err != nil {
		return err
	}

	return s.client.Add(cluster, cfg)
}

func (s *KubeService) newConfig(ctx context.Context, uuid, name string) (*rest.Config, error) {
	unit, err := s.application.GetLeader(ctx, uuid, name)
	if err != nil {
		return nil, err
	}

	unitInfo, err := s.application.GetUnitInfo(ctx, uuid, unit)
	if err != nil {
		return nil, err
	}

	endpoint, err := extractEndpointFromUnitInfo(unitInfo)
	if err != nil {
		return nil, err
	}

	clientToken, err := extractClientTokenFromUnitInfo(unitInfo)
	if err != nil {
		return nil, err
	}

	return &rest.Config{
		Host: endpoint,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true,
		},
		BearerToken: clientToken,
	}, nil
}

func extractEndpointFromUnitInfo(unitInfo *application.UnitInfo) (string, error) {
	endpoints := []string{}

	// Extract endpoints from relation data
	for _, erd := range unitInfo.RelationData {
		for _, rd := range erd.UnitRelationData {
			if endpointData, ok := rd.UnitData["api-endpoints"]; ok && endpointData != nil {
				if endpointStr, ok := endpointData.(string); ok {
					if err := json.Unmarshal([]byte(endpointStr), &endpoints); err != nil {
						return "", err
					}
				}
			}
		}
	}

	if len(endpoints) > 0 {
		return endpoints[0], nil
	}

	return "", errors.New("endpoint not found")
}

func extractClientTokenFromUnitInfo(unitInfo *application.UnitInfo) (string, error) {
	credentials := make(map[string]model.ControlPlaneCredential)

	// Extract credentials from relation data
	for _, erd := range unitInfo.RelationData {
		for _, rd := range erd.UnitRelationData {
			if credsData, ok := rd.UnitData["creds"]; ok && credsData != nil {
				if credsStr, ok := credsData.(string); ok {
					if err := json.Unmarshal([]byte(credsStr), &credentials); err != nil {
						return "", err
					}
				}
			}
		}
	}

	for _, cred := range credentials {
		return cred.ClientToken, nil
	}

	return "", errors.New("token not found")
}
