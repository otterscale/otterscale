package workload

type Destination struct {
	Metadata `yaml:",inline"`
}

func (d *Destination) Validate() error {
	return nil
}
