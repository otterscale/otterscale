package workload

type Destination struct {
	*Metadata
}

func (d *Destination) Validate() error {
	return nil
}
