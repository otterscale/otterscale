package spec

type Source struct {
	Metadata `yaml:",inline"`
}

func (s *Source) GetMetadata() Metadata {
	return s.Metadata
}

func (s *Source) Validate() error {
	return nil
}
