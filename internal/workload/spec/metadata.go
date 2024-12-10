package spec

type Metadata struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	Path    string `yaml:"path"`
}

func (m *Metadata) Validate() error {
	return nil
}
