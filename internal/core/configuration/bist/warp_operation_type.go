//go:generate go run -mod=mod golang.org/x/tools/cmd/stringer -type=WarpOperationType -output=warp_operation_type_string.go -linecomment=true

package bist

type WarpOperationType int32

const (
	WarpOperationTypeGet    WarpOperationType = iota // get
	WarpOperationTypePut                             // put
	WarpOperationTypeDelete                          // delete
	WarpOperationTypeList                            // list
	WarpOperationTypeStat                            // stat
	WarpOperationTypeMixed                           // mixed
)
