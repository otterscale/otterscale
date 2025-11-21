package storage

//go:generate go run -mod=mod golang.org/x/tools/cmd/stringer -type=PoolType -output=pool_type_string.go -linecomment=true

type PoolType int32

const (
	PoolTypeUnspecified PoolType = iota // unspecified
	PoolTypeErasure                     // erasure
	PoolTypeReplicated                  // replicated
)
