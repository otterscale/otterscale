package spec

type Spec interface {
	GetMetadata() Metadata
	Validate() error
}
