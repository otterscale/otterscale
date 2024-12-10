package spec

type Destination struct {
	Metadata `yaml:",inline"`
}

func (s *Destination) GetMetadata() Metadata {
	return s.Metadata
}

func (d *Destination) Validate() error {
	return nil
}
