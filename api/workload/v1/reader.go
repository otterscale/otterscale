package workload

import (
	"errors"
	"fmt"
	"os"

	"buf.build/go/protoyaml"

	"github.com/openhdc/openhdc/api/property/v1"
)

type Reader struct {
	wls *Workloads
}

func NewReader(files []string) (*Reader, error) {
	r := &Reader{
		wls: &Workloads{},
	}
	for _, file := range files {
		if err := r.readFile(file); err != nil {
			return nil, err
		}
	}
	if err := r.validate(); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Reader) Workloads() *Workloads {
	return r.wls
}

func (r *Reader) readFile(file string) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	var w Workload
	if err := protoyaml.Unmarshal(data, &w); err != nil {
		return err
	}
	// set internal filepath
	w.SetInternal(Internal_builder{FilePath: &file}.Build())
	// append
	switch w.GetKind() {
	case property.WorkloadKind_source:
		r.wls.AppendSources(&w)
	case property.WorkloadKind_destination:
		r.wls.AppendDestinations(&w)
	case property.WorkloadKind_transformer:
		r.wls.AppendTransformers(&w)
	default:
		return fmt.Errorf("invalid kind %s", w.GetKind())
	}
	return nil
}

func (r *Reader) validate() error {
	if len(r.wls.Sources()) == 0 {
		return errors.New("require at least one source")
	}
	if len(r.wls.Destinations()) == 0 {
		return errors.New("require at least one destination")
	}
	return nil
}
