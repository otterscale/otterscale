package connector

type Kind string

const (
	KindUnspecified Kind = ""
	KindSource      Kind = "source"
	KindDestination Kind = "destination"
)

type WriteKind int

const (
	WriteKindUnspecified WriteKind = iota
	WriteKindInsert
	WriteKindUpdate
	WriteKindUpsert
	WriteKindDelete
	WriteKindMigrate
)
