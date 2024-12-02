package workload

type Source struct {
	*Metadata
}

func (s *Source) Validate() error {
	return nil
}
