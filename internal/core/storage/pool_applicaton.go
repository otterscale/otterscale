//go:generate go run -mod=mod golang.org/x/tools/cmd/stringer -type=PoolApplication -output=pool_application_string.go -linecomment=true

package storage

type PoolApplication int32

const (
	PoolApplicationUnspecified PoolApplication = iota
	PoolApplicationBlock                       // rbd
	PoolApplicationFile                        // cephfs
	PoolApplicationObject                      // rgw
)
