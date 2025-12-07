//go:generate go run -mod=mod golang.org/x/tools/cmd/stringer -type=Type -output=type_string.go

package workload

type Type int32

const (
	TypeUnknown Type = iota
	TypeDeployment
	TypeStatefulSet
	TypeDaemonSet
)
