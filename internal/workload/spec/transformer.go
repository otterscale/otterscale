package spec

type Transformer struct {
	Metadata `yaml:",inline"`
}

func (s *Transformer) GetMetadata() Metadata {
	return s.Metadata
}

func (t *Transformer) Validate() error {
	return nil
}
