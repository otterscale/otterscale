package workload

import (
	"errors"
	"fmt"
	"os"
	"slices"

	"github.com/goccy/go-yaml"
)

type Reader struct {
	Sources      []*Source
	Destinations []*Destination
	Transformers []*Transformer
}

func (r *Reader) readFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var w Workload
	if err := yaml.Unmarshal(data, &w); err != nil {
		return err
	}
	switch w.Kind {
	case KindSource:
		source, err := w.Source()
		if err != nil {
			return err
		}
		if slices.ContainsFunc(r.Sources, func(s *Source) bool {
			return s.Name == source.Name
		}) {
			return fmt.Errorf("duplicate source name %s", source.Name)
		}
		if err := source.Validate(); err != nil {
			return fmt.Errorf("failed to validate source %s: %w", source.Name, err)
		}
		r.Sources = append(r.Sources, source)
		return nil

	case KindDestination:
		destination, err := w.Destination()
		if err != nil {
			return err
		}
		if slices.ContainsFunc(r.Destinations, func(d *Destination) bool {
			return d.Name == destination.Name
		}) {
			return fmt.Errorf("duplicate destination name %s", destination.Name)
		}
		if err := destination.Validate(); err != nil {
			return fmt.Errorf("failed to validate destination %s: %w", destination.Name, err)
		}
		r.Destinations = append(r.Destinations, destination)
		return nil

	case KindTransformer:
		transformer, err := w.Transformer()
		if err != nil {
			return err
		}
		if slices.ContainsFunc(r.Transformers, func(t *Transformer) bool {
			return t.Name == transformer.Name
		}) {
			return fmt.Errorf("duplicate transformer name %s", transformer.Name)
		}
		if err := transformer.Validate(); err != nil {
			return fmt.Errorf("failed to validate transformer %s: %w", transformer.Name, err)
		}
		r.Transformers = append(r.Transformers, transformer)
		return nil

	default:
		return fmt.Errorf("invalid kind %s", w.Kind)
	}
}

func (r *Reader) validate() error {
	if len(r.Sources) == 0 {
		return errors.New("require at least one source")
	}
	if len(r.Destinations) == 0 {
		return errors.New("require at least one destination")
	}
	return nil
}

func NewReader(paths []string) (*Reader, error) {
	r := &Reader{
		Sources:      []*Source{},
		Destinations: []*Destination{},
		Transformers: []*Transformer{},
	}
	for _, path := range paths {
		if err := r.readFile(path); err != nil {
			return nil, err
		}
	}
	if err := r.validate(); err != nil {
		return nil, err
	}
	return r, nil
}
