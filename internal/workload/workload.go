package workload

import (
	"github.com/goccy/go-yaml"
)

type Workload struct {
	Kind Kind `yaml:"kind"`
	Spec any  `yaml:"spec"`
}

func (w *Workload) Validate() error {
	return nil
}

func (w *Workload) Source() (*Source, error) {
	data, err := yaml.Marshal(w.Spec)
	if err != nil {
		return nil, err
	}
	v := Source{}
	if err := yaml.Unmarshal(data, &v); err != nil {
		return nil, err
	}
	return &v, nil
}

func (w *Workload) Destination() (*Destination, error) {
	data, err := yaml.Marshal(w.Spec)
	if err != nil {
		return nil, err
	}
	v := Destination{}
	if err := yaml.Unmarshal(data, &v); err != nil {
		return nil, err
	}
	return &v, nil
}

func (w *Workload) Transformer() (*Transformer, error) {
	data, err := yaml.Marshal(w.Spec)
	if err != nil {
		return nil, err
	}
	v := Transformer{}
	if err := yaml.Unmarshal(data, &v); err != nil {
		return nil, err
	}
	return &v, nil
}
