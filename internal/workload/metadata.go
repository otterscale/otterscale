package workload

type Metadata struct {
	Name string
}

func (m *Metadata) Validate() error {
	return nil
}
