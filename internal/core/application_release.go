package core

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-faker/faker/v4"

	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/strvals"

	"k8s.io/client-go/rest"
	"sigs.k8s.io/yaml"
)

type Release = release.Release

type ReleaseRepo interface {
	List(config *rest.Config, namespace string) ([]release.Release, error)
	Get(restConfig *rest.Config, namespace, name string) (*release.Release, error)
	Install(config *rest.Config, namespace, name string, dryRun bool, chartRef string, labels, annotations map[string]string, values map[string]any) (*release.Release, error)
	Uninstall(config *rest.Config, namespace, name string, dryRun bool) (*release.Release, error)
	Upgrade(config *rest.Config, namespace, name string, dryRun bool, chartRef string, values map[string]any) (*release.Release, error)
	Rollback(config *rest.Config, namespace, name string, dryRun bool) error
	GetValues(config *rest.Config, namespace, name string) (map[string]any, error)
}

func (uc *ApplicationUseCase) ListReleases(ctx context.Context, uuid, facility string) ([]Release, error) {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.release.List(config, "")
}

func (uc *ApplicationUseCase) CreateRelease(ctx context.Context, uuid, facility, namespace, name string, dryRun bool, chartRef, valuesYAML string, valuesMap map[string]string) (*Release, error) {
	values, err := toReleaseValues(valuesYAML, valuesMap)
	if err != nil {
		return nil, err
	}

	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}

	// labels
	labels := map[string]string{
		ApplicationReleaseNameLabel: name,
	}
	if strings.Contains(chartRef, "llm-d-incubation/llm-d-modelservice") {
		if modelName, ok := values["modelArtifacts.name"].(string); ok {
			labels[ApplicationReleaseLLMDModelNameLabel] = modelName
		}
	}

	// annotations
	annotations := map[string]string{
		ApplicationReleaseChartRefAnnotation: chartRef,
	}
	return uc.release.Install(config, namespace, getReleaseName(name), dryRun, chartRef, labels, annotations, values)
}

func (uc *ApplicationUseCase) UpdateRelease(ctx context.Context, uuid, facility, namespace, name string, dryRun bool, chartRef, valuesYAML string) (*Release, error) {
	values, err := toReleaseValues(valuesYAML, nil)
	if err != nil {
		return nil, err
	}

	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return nil, err
	}
	return uc.release.Upgrade(config, namespace, name, dryRun, chartRef, values)
}

func (uc *ApplicationUseCase) DeleteRelease(ctx context.Context, uuid, facility, namespace, name string, dryRun bool) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}
	_, err = uc.release.Uninstall(config, namespace, name, dryRun)
	return err
}

func (uc *ApplicationUseCase) RollbackRelease(ctx context.Context, uuid, facility, namespace, name string, dryRun bool) error {
	config, err := kubeConfig(ctx, uc.facility, uc.action, uuid, facility)
	if err != nil {
		return err
	}
	return uc.release.Rollback(config, namespace, name, dryRun)
}

func getReleaseName(name string) string {
	if name != "" {
		return name
	}
	return strings.ToLower(faker.FirstName() + "-" + faker.Username())
}

func toReleaseValues(valuesYAML string, valuesMap map[string]string) (map[string]any, error) {
	// advanced
	values := map[string]any{}
	if err := yaml.Unmarshal([]byte(valuesYAML), &values); err != nil {
		return nil, err
	}

	// basic
	vals := []string{}
	for k, v := range valuesMap {
		if v != "" {
			vals = append(vals, fmt.Sprintf("%s=%s", k, v))
		}
	}
	if err := strvals.ParseInto(strings.Join(vals, ","), values); err != nil {
		return nil, err
	}

	return values, nil
}
