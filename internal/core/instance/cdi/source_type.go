//go:generate go run -mod=mod golang.org/x/tools/cmd/stringer -type=SourceType -output=source_type_string.go

package cdi

type SourceType int32

const (
	SourceTypeBlankImage SourceType = iota
	SourceTypeHTTPURL
	SourceTypeExistingPersistentVolumeClaim
)
