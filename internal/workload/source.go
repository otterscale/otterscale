package workload

type Source struct {
	Metadata `yaml:",inline"`
}

func (s *Source) Validate() error {
	return nil
}
