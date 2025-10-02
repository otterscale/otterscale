package kube

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"

	oscore "github.com/otterscale/otterscale/internal/core"
)

type apps struct {
	kube *Kube
}

func NewApps(kube *Kube) oscore.KubeAppsRepo {
	return &apps{
		kube: kube,
	}
}

var _ oscore.KubeAppsRepo = (*apps)(nil)

func (r *apps) ListDeployments(ctx context.Context, config *rest.Config, namespace string) ([]oscore.Deployment, error) {
	clientset, err := r.kube.clientset(config)
	if err != nil {
		return nil, err
	}

	opts := metav1.ListOptions{}
	list, err := clientset.AppsV1().Deployments(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (r *apps) ListDeploymentsByLabel(ctx context.Context, config *rest.Config, namespace, label string) ([]oscore.Deployment, error) {
	clientset, err := r.kube.clientset(config)
	if err != nil {
		return nil, err
	}

	opts := metav1.ListOptions{
		LabelSelector: label,
	}

	list, err := clientset.AppsV1().Deployments(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (r *apps) GetDeployment(ctx context.Context, config *rest.Config, namespace, name string) (*oscore.Deployment, error) {
	clientset, err := r.kube.clientset(config)
	if err != nil {
		return nil, err
	}

	opts := metav1.GetOptions{}
	return clientset.AppsV1().Deployments(namespace).Get(ctx, name, opts)
}

func (r *apps) UpdateDeployment(ctx context.Context, config *rest.Config, namespace string, deployment *oscore.Deployment) (*oscore.Deployment, error) {
	clientset, err := r.kube.clientset(config)
	if err != nil {
		return nil, err
	}

	opts := metav1.UpdateOptions{}
	return clientset.AppsV1().Deployments(namespace).Update(ctx, deployment, opts)
}

func (r *apps) ListStatefulSets(ctx context.Context, config *rest.Config, namespace string) ([]oscore.StatefulSet, error) {
	clientset, err := r.kube.clientset(config)
	if err != nil {
		return nil, err
	}

	opts := metav1.ListOptions{}
	list, err := clientset.AppsV1().StatefulSets(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (r *apps) GetStatefulSet(ctx context.Context, config *rest.Config, namespace, name string) (*oscore.StatefulSet, error) {
	clientset, err := r.kube.clientset(config)
	if err != nil {
		return nil, err
	}

	opts := metav1.GetOptions{}
	return clientset.AppsV1().StatefulSets(namespace).Get(ctx, name, opts)
}

func (r *apps) UpdateStatefulSet(ctx context.Context, config *rest.Config, namespace string, statefulSet *oscore.StatefulSet) (*oscore.StatefulSet, error) {
	clientset, err := r.kube.clientset(config)
	if err != nil {
		return nil, err
	}

	opts := metav1.UpdateOptions{}
	return clientset.AppsV1().StatefulSets(namespace).Update(ctx, statefulSet, opts)
}

func (r *apps) ListDaemonSets(ctx context.Context, config *rest.Config, namespace string) ([]oscore.DaemonSet, error) {
	clientset, err := r.kube.clientset(config)
	if err != nil {
		return nil, err
	}

	opts := metav1.ListOptions{}
	list, err := clientset.AppsV1().DaemonSets(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (r *apps) GetDaemonSet(ctx context.Context, config *rest.Config, namespace, name string) (*oscore.DaemonSet, error) {
	clientset, err := r.kube.clientset(config)
	if err != nil {
		return nil, err
	}

	opts := metav1.GetOptions{}
	return clientset.AppsV1().DaemonSets(namespace).Get(ctx, name, opts)
}

func (r *apps) UpdateDaemonSet(ctx context.Context, config *rest.Config, namespace string, daemonSet *oscore.DaemonSet) (*oscore.DaemonSet, error) {
	clientset, err := r.kube.clientset(config)
	if err != nil {
		return nil, err
	}

	opts := metav1.UpdateOptions{}
	return clientset.AppsV1().DaemonSets(namespace).Update(ctx, daemonSet, opts)
}
