package connector

type Kind int

const (
	KindUnspecified Kind = iota
	KindSource
	KindDestination
)

type WriteKind int

const (
	WriteKindUnspecified Kind = iota
	WriteKindInsert
	WriteKindUpdate
	WriteKindUpsert
	WriteKindDelete
	WriteKindMigrate
)
