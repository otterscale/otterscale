package extension

import (
	"context"
	"fmt"
)

type component struct {
	ID          string
	DisplayName string
	Description string
	Logo        string

	Chart *chartComponent
	CRD   *crdComponent

	Dependencies []string
	PostFunc     func(ctx context.Context, scope string) error
}

type chartComponent struct {
	Name      string
	Namespace string
	Ref       string
	Version   string
	ValuesMap map[string]string
}

type crdComponent struct {
	Ref               string
	Version           string
	VersionAnnotation string
}

func (uc *UseCase) whichComponent(id string) (*component, error) {
	components := []component{}
	components = append(components, metricsComponents...)
	components = append(components, serviceMeshComponents...)
	components = append(components, registryComponents...)
	components = append(components, modelComponents...)
	components = append(components, instanceComponents...)
	components = append(components, storageComponents...)

	for i := range components {
		if components[i].ID == id {
			return &components[i], nil
		}
	}

	return nil, fmt.Errorf("extension component %s not found", id)
}
