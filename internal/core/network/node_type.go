//go:generate go run -mod=mod golang.org/x/tools/cmd/stringer -type=NodeType -output=node_type_string.go -linecomment=true

package network

type NodeType int32

const (
	NodeTypeMachine              NodeType = iota // Machine
	NodeTypeDevice                               // Device
	NodeTypeRackController                       // Rack Controller
	NodeTypeRegionRackController                 // Region & Rack Controller
	NodeTypeRegionController                     // Region Controller
)
