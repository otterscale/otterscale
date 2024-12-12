package workload

import (
	"errors"
	"fmt"
	"os"

	"buf.build/go/protoyaml"
)

type Reader struct {
	Sources      []*Workload
	Destinations []*Workload
	Transformers []*Workload
}

func NewReader(paths []string) (*Reader, error) {
	r := &Reader{
		Sources:      []*Workload{},
		Destinations: []*Workload{},
		Transformers: []*Workload{},
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

func (r *Reader) readFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var w Workload
	if err := protoyaml.Unmarshal(data, &w); err != nil {
		return err
	}
	w.Internal = &Internal{
		FilePath: path,
	}
	switch w.Kind {
	case Kind_source:
		r.Sources = append(r.Sources, &w)
	case Kind_destination:
		r.Sources = append(r.Sources, &w)
	case Kind_transformer:
		r.Sources = append(r.Sources, &w)
	default:
		return fmt.Errorf("invalid kind %s", w.Kind)
	}
	return nil
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
