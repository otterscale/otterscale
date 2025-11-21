//go:generate go run -mod=mod golang.org/x/tools/cmd/stringer -type=ObjectServiceType -output=object_service_type_string.go -linecomment=true

package bist

type ObjectServiceType int32

const (
	ObjectServiceTypeUnspecified ObjectServiceType = iota // unspecified
	ObjectServiceTypeCeph                                 // ceph
	ObjectServiceTypeMinIO                                // minio
)
