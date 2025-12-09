package extension

import (
	"context"
	"errors"
	"fmt"
)

func (uc *UseCase) getNodePort(ctx context.Context, scope, namespace, name, portName string) (int32, error) {
	svc, err := uc.service.Get(ctx, scope, namespace, name)
	if err != nil {
		return 0, err
	}

	ports := svc.Spec.Ports

	for i := range ports {
		if ports[i].Name == portName {
			return ports[i].NodePort, nil
		}
	}

	return 0, fmt.Errorf("service %q has no %q port defined", name, portName)
}

func (uc *UseCase) getConfig(ctx context.Context, scope, name string) (map[string]any, error) {
	return uc.facility.Config(ctx, scope, name)
}

func (uc *UseCase) setConfig(ctx context.Context, scope, name, key, value string) error {
	return uc.facility.Update(ctx, scope, name, map[string]string{key: value})
}

func getValue[T any](config map[string]any, key string) (T, error) {
	var zero T

	target, ok := config[key].(map[string]any)
	if !ok {
		return zero, fmt.Errorf("invalid type for %s field", key)
	}

	source, ok := target["source"].(string)
	if !ok {
		return zero, errors.New("invalid type for source field")
	}

	if source == "unset" {
		return zero, nil
	}

	value, ok := target["value"].(T)
	if !ok {
		return zero, fmt.Errorf("invalid type for %s.value field", key)
	}

	return value, nil
}
