package workload

import (
	"github.com/goccy/go-yaml"

	"github.com/openhdc/openhdc/internal/workload/spec"
)

type Workload struct {
	Kind Kind `yaml:"kind"`
	Spec any  `yaml:"spec"`
}

func (w *Workload) Validate() error {
	return nil
}

func (w *Workload) Source() (*spec.Source, error) {
	data, err := yaml.Marshal(w.Spec)
	if err != nil {
		return nil, err
	}
	v := spec.Source{}
	if err := yaml.Unmarshal(data, &v); err != nil {
		return nil, err
	}
	return &v, nil
}

func (w *Workload) Destination() (*spec.Destination, error) {
	data, err := yaml.Marshal(w.Spec)
	if err != nil {
		return nil, err
	}
	v := spec.Destination{}
	if err := yaml.Unmarshal(data, &v); err != nil {
		return nil, err
	}
	return &v, nil
}

func (w *Workload) Transformer() (*spec.Transformer, error) {
	data, err := yaml.Marshal(w.Spec)
	if err != nil {
		return nil, err
	}
	v := spec.Transformer{}
	if err := yaml.Unmarshal(data, &v); err != nil {
		return nil, err
	}
	return &v, nil
}
