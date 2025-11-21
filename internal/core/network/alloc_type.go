//go:generate go run -mod=mod golang.org/x/tools/cmd/stringer -type=AllocType -output=alloc_type_string.go -linecomment=true

package network

type AllocType int32

const (
	AllocTypeAutomatic    AllocType = iota     // Automatic
	AllocTypeSticky                            // Sticky
	AllocTypeUserReserved AllocType = iota + 2 // User Reserved
	AllocTypeDHCP                              // DHCP
	AllocTypeDiscovered                        // Discovered
)
