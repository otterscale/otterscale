//go:generate go run -mod=mod golang.org/x/tools/cmd/stringer -type=Type -output=type_string.go

package extension

type Type int32

const (
	TypeUnspecified Type = iota
	TypeMetrics
	TypeServiceMesh
	TypeRegistry
	TypeModel
	TypeInstance
	TypeStorage
)
