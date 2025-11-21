//go:generate go run -mod=mod golang.org/x/tools/cmd/stringer -type=HealthStatus -output=health_status_string.go

package environment

type HealthStatus int32

const (
	HealthStatusUnknown HealthStatus = iota
	HealthStatusOK
	HealthStatusNotInstalled
)
