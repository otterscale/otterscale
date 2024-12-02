package workload

type Workload struct {
	Kind Kind `yaml:"kind"`
	Spec any  `yaml:"spec"`
}

func (w *Workload) Validate() error {
	return nil
}
