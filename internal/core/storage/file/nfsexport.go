package file

import (
	"context"
)

type AccessType string
type Squash string

const (
	AccessTypeUnspecified AccessType = ""
	AccessTypeRW          AccessType = "RW"
	AccessTypeRO          AccessType = "RO"
	AccessTypeMDOnly      AccessType = "MDONLY"
)

const (
	SquashUnspecified Squash = ""
	SquashNone        Squash = "None"
	SquashRoot        Squash = "Root"
	SquashAll         Squash = "All"
	SquashRootID      Squash = "RootId"
)

type Client struct {
	Addresses  []string
	AccessType AccessType
	Squash     Squash
}

type NFSExport struct {
	ClusterID            string
	FileSystemName       string
	EnableSecurityLabel  bool
	GroupName            string
	SubName              string
	Path                 string
	ProtocolNFSV3        bool
	ProtocolNFSV4        bool
	PseudoPath           string
	AccessType           AccessType
	Squash               Squash
	TransportProtocolUDP bool
	TransportProtocolTCP bool
	Clients              []Client
}

// Note: Ceph create and update operations only return error status.

type NFSExportRepo interface {
	List(ctx context.Context, scope, cluster string) ([]NFSExport, error)
	Get(ctx context.Context, scope, cluster, pseudo string) (*NFSExport, error)
	Create(ctx context.Context, scope, cluster, filesystem string, enable_security_label bool, group, subvolume string, protocol_nfs_v3, protocol_nfs_v4 bool, pseudo string, accessType AccessType, squash Squash, transport_protocol_udp, transport_protocol_tcp bool, clients []Client) error
	Update(ctx context.Context, scope, cluster string, enable_security_label bool, protocol_nfs_v3, protocol_nfs_v4 bool, pseudo string, accessType AccessType, squash Squash, transport_protocol_udp, transport_protocol_tcp bool, clients []Client) error
	Delete(ctx context.Context, scope, cluster, pseudo string) error
}

/*func (uc *UseCase) ListNFSExports(ctx context.Context, scope, cluster string) ([]NFSExport, error) {
	return uc.nfsExport.List(ctx, scope, cluster)
}

func (uc *UseCase) CreateNFSExport(ctx context.Context, scope, cluster, filesystem string, enable_security_label bool, group, subvolume string, protocol_nfs_v3, protocol_nfs_v4 bool, pseudo string, accessType AccessType, sq Squash, transport_protocol_udp, transport_protocol_tcp bool, cs []Client) (*NFSExport, error) {
	if err := uc.nfsExport.Create(ctx, scope, cluster, filesystem, enable_security_label, group, subvolume, protocol_nfs_v3, protocol_nfs_v4, pseudo, accessType, sq, transport_protocol_udp, transport_protocol_tcp, cs); err != nil {
		return nil, err
	}
	return uc.nfsExport.Get(ctx, scope, cluster, pseudo)
}

func (uc *UseCase) GetNFSExport(ctx context.Context, scope, cluster, pseudo string) (*NFSExport, error) {
	return uc.nfsExport.Get(ctx, scope, cluster, pseudo)
}

func (uc *UseCase) UpdateNFSExport(ctx context.Context, scope, cluster string, enable_security_label bool, protocol_nfs_v3, protocol_nfs_v4 bool, pseudo string, accessType AccessType, sq Squash, transport_protocol_udp, transport_protocol_tcp bool, cs []Client) (*NFSExport, error) {
	if err := uc.nfsExport.Update(ctx, scope, cluster, enable_security_label, protocol_nfs_v3, protocol_nfs_v4, pseudo, accessType, sq, transport_protocol_udp, transport_protocol_tcp, cs); err != nil {
		return nil, err
	}
	return uc.nfsExport.Get(ctx, scope, cluster, pseudo)
}

func (uc *UseCase) DeleteNFSExport(ctx context.Context, scope, cluster, pseudo string) error {
	return uc.nfsExport.Delete(ctx, scope, cluster, pseudo)
}*/
