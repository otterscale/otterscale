package file

import (
	"context"
	"fmt"
	"sort"
	"strings"
)

type AccessType string
type Squash string
type Sectype string
type NFSProtocol uint32
type NFSTransport string

const (
	NFSProtocolV3 NFSProtocol = 3
	NFSProtocolV4 NFSProtocol = 4
)

const (
	AccessTypeUnspecified AccessType = ""
	AccessTypeRW          AccessType = "RW"
	AccessTypeRO          AccessType = "RO"
	AccessTypeNone        AccessType = "NONE"
)

const (
	SquashUnspecified Squash = ""
	SquashNone        Squash = "none"
	SquashRoot        Squash = "root"
	SquashAll         Squash = "all"
	SquashRootID      Squash = "rootid"
)

const (
	NFSTransportTCP NFSTransport = "TCP"
	NFSTransportUDP NFSTransport = "UDP"
)

const (
	SysSec   Sectype = "sys"
	NoneSec  Sectype = "none"
	Krb5Sec  Sectype = "krb5"
	Krb5iSec Sectype = "krb5i"
	Krb5pSec Sectype = "krb5p"
)

type FSAL struct {
	Name   string
	UserID string
	FSName string
}

type Client struct {
	Addresses  []string
	AccessType AccessType
	Squash     Squash
}

type NFSExport struct {
	ExportID            uint64
	ClusterID           string
	FileSystemName      string
	GroupName           string
	SubvolumeName       string
	EnableSecurityLabel bool
	Path                string
	Protocols           []NFSProtocol
	PseudoPath          string
	AccessType          AccessType
	Squash              Squash
	Transports          []NFSTransport
	Clients             []Client
	Sectypes            []Sectype
}

type NFSExportSpec struct {
	PseudoPath string

	EnableSecurityLabel bool
	Protocols           []NFSProtocol
	AccessType          AccessType
	Squash              Squash
	Transports          []NFSTransport
	Clients             []Client
	Sectypes            []Sectype
}

// Note: Ceph create and update operations only return error status.

type NFSExportRepo interface {
	List(ctx context.Context, scope, cluster string) ([]NFSExport, error)
	Get(ctx context.Context, scope, cluster, pseudo string) (*NFSExport, error)
	Apply(ctx context.Context, scope, cluster string, export *NFSExport) error
	Delete(ctx context.Context, scope, cluster, pseudo string) error
}

func (uc *UseCase) ListNFSExports(ctx context.Context, scope, cluster string) ([]NFSExport, error) {
	return uc.nfsexport.List(ctx, scope, cluster)
}

func (uc *UseCase) ApplyNFSExport(ctx context.Context, scope, cluster, filesystem, group, subvolume string, spec NFSExportSpec) (*NFSExport, error) {
	if strings.TrimSpace(cluster) == "" {
		return nil, fmt.Errorf("cluster is required")
	}
	if strings.TrimSpace(filesystem) == "" {
		return nil, fmt.Errorf("filesystem is required")
	}
	if strings.TrimSpace(subvolume) == "" {
		return nil, fmt.Errorf("subvolume is required")
	}
	spec.PseudoPath = strings.TrimSpace(spec.PseudoPath)
	if spec.PseudoPath == "" {
		return nil, fmt.Errorf("pseudo_path is required")
	}

	sv, err := uc.subvolume.Get(ctx, scope, filesystem, group, subvolume)
	if err != nil {
		return nil, err
	}
	path := sv.Info.Path

	spec.Protocols = normalizeProtocols(spec.Protocols)
	if len(spec.Protocols) == 0 {
		spec.Protocols = []NFSProtocol{NFSProtocolV4}
	}

	spec.Transports = normalizeTransports(spec.Transports)
	if len(spec.Transports) == 0 {
		spec.Transports = []NFSTransport{NFSTransportTCP}
	}

	spec.Sectypes = normalizeSectypes(spec.Sectypes)
	spec.Clients = normalizeClients(spec.Clients)

	e := &NFSExport{
		ClusterID:      cluster,
		FileSystemName: filesystem,
		GroupName:      group,
		SubvolumeName:  subvolume,

		EnableSecurityLabel: spec.EnableSecurityLabel,
		Path:                path,
		Protocols:           spec.Protocols,
		PseudoPath:          spec.PseudoPath,
		AccessType:          spec.AccessType,
		Squash:              spec.Squash,
		Transports:          spec.Transports,
		Clients:             spec.Clients,
		Sectypes:            spec.Sectypes,
	}

	if err := uc.nfsexport.Apply(ctx, scope, cluster, e); err != nil {
		return nil, err
	}
	return uc.nfsexport.Get(ctx, scope, cluster, e.PseudoPath)
}

func (uc *UseCase) GetNFSExport(ctx context.Context, scope, cluster, pseudo string) (*NFSExport, error) {
	return uc.nfsexport.Get(ctx, scope, cluster, pseudo)
}

func (uc *UseCase) DeleteNFSExport(ctx context.Context, scope, cluster, pseudo string) error {
	return uc.nfsexport.Delete(ctx, scope, cluster, pseudo)
}

func normalizeProtocols(in []NFSProtocol) []NFSProtocol {
	if len(in) == 0 {
		return nil
	}
	m := make(map[NFSProtocol]struct{}, len(in))
	out := make([]NFSProtocol, 0, len(in))
	for _, p := range in {
		if p != NFSProtocolV3 && p != NFSProtocolV4 {
			continue
		}
		if _, ok := m[p]; ok {
			continue
		}
		m[p] = struct{}{}
		out = append(out, p)
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

func normalizeTransports(in []NFSTransport) []NFSTransport {
	if len(in) == 0 {
		return nil
	}
	m := make(map[NFSTransport]struct{}, len(in))
	out := make([]NFSTransport, 0, len(in))
	for _, t := range in {
		tt := NFSTransport(strings.ToUpper(strings.TrimSpace(string(t))))
		if tt != NFSTransportTCP && tt != NFSTransportUDP {
			continue
		}
		if _, ok := m[tt]; ok {
			continue
		}
		m[tt] = struct{}{}
		out = append(out, tt)
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

func normalizeSectypes(in []Sectype) []Sectype {
	if len(in) == 0 {
		return nil
	}
	m := make(map[Sectype]struct{}, len(in))
	out := make([]Sectype, 0, len(in))
	for _, s := range in {
		ss := Sectype(strings.ToLower(strings.TrimSpace(string(s))))
		switch ss {
		case SysSec, NoneSec, Krb5Sec, Krb5iSec, Krb5pSec:
			// ok
		default:
			continue
		}
		if _, ok := m[ss]; ok {
			continue
		}
		m[ss] = struct{}{}
		out = append(out, ss)
	}
	// keep stable order (optional)
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

func normalizeClients(in []Client) []Client {
	if len(in) == 0 {
		return nil
	}
	out := make([]Client, 0, len(in))
	for _, c := range in {
		addrs := make([]string, 0, len(c.Addresses))
		for _, a := range c.Addresses {
			a = strings.TrimSpace(a)
			if a == "" {
				continue
			}
			addrs = append(addrs, a)
		}
		if len(addrs) == 0 {
			// skip empty client rule
			continue
		}
		c.Addresses = addrs
		out = append(out, c)
	}
	return out
}
