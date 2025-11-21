//go:generate go run -mod=mod golang.org/x/tools/cmd/stringer -type=EntityType -output=entity_type_string.go

package smb

type EntityType int32

const (
	EntityTypeUnknown EntityType = iota
	EntityTypeUser
	EntityTypeGroup
)
