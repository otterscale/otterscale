package workload

type Transformer struct {
	Metadata `yaml:",inline"`
}

func (t *Transformer) Validate() error {
	return nil
}
