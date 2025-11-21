package smb

//go:generate go run -mod=mod golang.org/x/tools/cmd/stringer -type=MapToGuest -output=map_to_guest_string.go -linecomment=true

type MapToGuest int32

const (
	MapToGuestNever       MapToGuest = iota // never
	MapToGuestBadUser                       // bad user
	MapToGuestBadPassword                   // bad password
)
