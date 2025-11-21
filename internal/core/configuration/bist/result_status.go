//go:generate go run -mod=mod golang.org/x/tools/cmd/stringer -type=ResultStatus -output=result_status_string.go

package bist

type ResultStatus int32

const (
	ResultStatusRunning ResultStatus = iota
	ResultStatusSucceeded
	ResultStatusFailed
)
