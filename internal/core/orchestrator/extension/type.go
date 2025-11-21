//go:generate go run -mod=mod golang.org/x/tools/cmd/stringer -type=Type -output=type_string.go

package extension

type Type int32

const (
	TypeGeneral Type = iota
	TypeRegistry
	TypeModel
	TypeInstance
	TypeStorage
)
