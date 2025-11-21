//go:generate go run -mod=mod golang.org/x/tools/cmd/stringer -type=SecurityMode -output=security_mode_string.go -linecomment=true

package smb

type SecurityMode int32

const (
	SecurityModeUser            SecurityMode = iota // user
	SecurityModeActiveDirectory                     // active-directory
)
