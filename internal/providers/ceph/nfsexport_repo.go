package ceph

import (
	"context"
	"strings"

	"github.com/otterscale/otterscale/internal/core/storage/file"
)

type nfsExportRepo struct {
	ceph *Ceph
}

func NewNFSExportRepo(ceph *Ceph) file.NFSExportRepo {
	return &nfsExportRepo{
		ceph: ceph,
	}
}

var _ file.NFSExportRepo = (*nfsExportRepo)(nil)

func (r *nfsExportRepo) List(_ context.Context, scope, clusterID string) ([]file.NFSExport, error) {
	conn, err := r.ceph.connection(scope)
	if err != nil {
		return nil, err
	}

	exports, err := listDetailedNFSExports(conn, clusterID)
	if err != nil {
		return nil, err
	}

	out := make([]file.NFSExport, 0, len(exports))
	for i := range exports {
		if e := r.toDomain(&exports[i]); e != nil {
			out = append(out, *e)
		}
	}
	return out, nil
}

func (r *nfsExportRepo) Apply(ctx context.Context, scope, clusterID string, e *file.NFSExport) error {
	conn, err := r.ceph.connection(scope)
	if err != nil {
		return err
	}

	if e == nil {
		return nil
	}

	ce := r.toWire(clusterID, e)
	return applyNFSExport(conn, clusterID, ce)
}

func (r *nfsExportRepo) Get(_ context.Context, scope, clusterID, pseudoPath string) (*file.NFSExport, error) {
	conn, err := r.ceph.connection(scope)
	if err != nil {
		return nil, err
	}

	info, err := getNFSExport(conn, clusterID, pseudoPath)
	if err != nil {
		return nil, err
	}

	return r.toDomain(info), nil
}

func (r *nfsExportRepo) Delete(_ context.Context, scope, clusterID, pseudoPath string) error {
	conn, err := r.ceph.connection(scope)
	if err != nil {
		return err
	}
	return removeNFSExport(conn, clusterID, pseudoPath)
}

func (r *nfsExportRepo) toDomain(info *cephNFSExport) *file.NFSExport {
	if info == nil {
		return nil
	}

	outProtocols := make([]file.NFSProtocol, 0, len(info.Protocols))
	for _, p := range info.Protocols {
		switch uint32(p) {
		case uint32(file.NFSProtocolV3):
			outProtocols = append(outProtocols, file.NFSProtocolV3)
		case uint32(file.NFSProtocolV4):
			outProtocols = append(outProtocols, file.NFSProtocolV4)
		}
	}

	outTransports := make([]file.NFSTransport, 0, len(info.Transports))
	for _, t := range info.Transports {
		tt := strings.ToUpper(strings.TrimSpace(t))
		switch file.NFSTransport(tt) {
		case file.NFSTransportTCP, file.NFSTransportUDP:
			outTransports = append(outTransports, file.NFSTransport(tt))
		}
	}

	outSecs := make([]file.Sectype, 0, len(info.SecType))
	for _, s := range info.SecType {
		ss := strings.ToLower(strings.TrimSpace(string(s)))
		switch file.Sectype(ss) {
		case file.SysSec, file.NoneSec, file.Krb5Sec, file.Krb5iSec, file.Krb5pSec:
			outSecs = append(outSecs, file.Sectype(ss))
		}
	}

	outClients := make([]file.Client, 0, len(info.Clients))
	for _, c := range info.Clients {
		addrs := make([]string, 0, len(c.Addresses))
		for _, a := range c.Addresses {
			a = strings.TrimSpace(a)
			if a != "" {
				addrs = append(addrs, a)
			}
		}
		if len(addrs) == 0 {
			continue
		}

		outClients = append(outClients, file.Client{
			Addresses:  addrs,
			AccessType: normalizeAccessType(c.AccessType),
			Squash:     normalizeSquashFromCeph(string(c.Squash)),
		})
	}

	return &file.NFSExport{
		ExportID:            info.ExportID,
		ClusterID:           info.ClusterID,
		FileSystemName:      strings.TrimSpace(info.FSAL.FileSystemName),
		GroupName:           "",
		SubvolumeName:       "",
		EnableSecurityLabel: info.SecurityLabel,
		Path:                info.Path,
		Protocols:           outProtocols,
		PseudoPath:          info.PseudoPath,
		AccessType:          normalizeAccessType(info.AccessType),
		Squash:              normalizeSquashFromCeph(string(info.Squash)),
		Transports:          outTransports,
		Clients:             outClients,
		Sectypes:            outSecs,
	}
}

func (r *nfsExportRepo) toWire(clusterID string, e *file.NFSExport) *cephNFSExport {
	ps := make([]int, 0, len(e.Protocols))
	for _, p := range e.Protocols {
		switch p {
		case file.NFSProtocolV3, file.NFSProtocolV4:
			ps = append(ps, int(p))
		}
	}

	ts := make([]string, 0, len(e.Transports))
	for _, t := range e.Transports {
		tt := strings.ToUpper(strings.TrimSpace(string(t)))
		if tt == "" {
			continue
		}
		switch file.NFSTransport(tt) {
		case file.NFSTransportTCP, file.NFSTransportUDP:
			ts = append(ts, tt)
		}
	}

	st := make([]string, 0, len(e.Sectypes))
	for _, s := range e.Sectypes {
		ss := strings.ToLower(strings.TrimSpace(string(s)))
		switch file.Sectype(ss) {
		case file.SysSec, file.NoneSec, file.Krb5Sec, file.Krb5iSec, file.Krb5pSec:
			st = append(st, ss)
		}
	}

	cs := make([]cephNFSClientInfo, 0, len(e.Clients))
	for _, c := range e.Clients {
		addrs := make([]string, 0, len(c.Addresses))
		for _, a := range c.Addresses {
			a = strings.TrimSpace(a)
			if a != "" {
				addrs = append(addrs, a)
			}
		}
		if len(addrs) == 0 {
			continue
		}

		cs = append(cs, cephNFSClientInfo{
			Addresses:  addrs,
			AccessType: string(normalizeAccessType(string(c.AccessType))),
			Squash:     toCephSquashString(c.Squash),
		})
	}

	return &cephNFSExport{
		ExportID:      e.ExportID,
		Path:          e.Path,
		ClusterID:     clusterID,
		PseudoPath:    e.PseudoPath,                 // json:"pseudo"
		AccessType:    string(e.AccessType),         // "RW"/"RO"/"NONE"
		Squash:        toCephSquashString(e.Squash), // None/Root/All/RootId
		SecurityLabel: e.EnableSecurityLabel,
		Protocols:     ps,
		Transports:    ts,
		FSAL: cephNFSFSALInfo{
			// 你目前 file.NFSExport 已經不帶 FSAL，就固定 CEPH + fs_name
			Name:           "CEPH",
			UserID:         "",
			FileSystemName: e.FileSystemName,
		},
		Clients: cs,
		SecType: st,
	}
}

// ---------- helpers ----------

func normalizeAccessType(s string) file.AccessType {
	ss := strings.ToUpper(strings.TrimSpace(s))
	switch file.AccessType(ss) {
	case file.AccessTypeRW, file.AccessTypeRO, file.AccessTypeNone:
		return file.AccessType(ss)
	default:
		return file.AccessTypeUnspecified
	}
}

// Ceph -> domain squash (None/Root/All/RootId) -> (none/root/all/rootid)
func normalizeSquashFromCeph(s string) file.Squash {
	ss := strings.ToLower(strings.TrimSpace(s))
	switch ss {
	case "none":
		return file.SquashNone
	case "root":
		return file.SquashRoot
	case "all":
		return file.SquashAll
	case "rootid", "root_id", "root-id":
		return file.SquashRootID
	}

	// 有些 Ceph 回來是 "None"/"Root"/"All"/"RootId"
	switch strings.TrimSpace(s) {
	case "None":
		return file.SquashNone
	case "Root":
		return file.SquashRoot
	case "All":
		return file.SquashAll
	case "RootId", "RootID":
		return file.SquashRootID
	default:
		return file.SquashUnspecified
	}
}

// domain squash (none/root/all/rootid) -> Ceph squash (None/Root/All/RootId)
func toCephSquash(s file.Squash) Squash {
	switch strings.ToLower(strings.TrimSpace(string(s))) {
	case "none":
		return NoneSquash
	case "root":
		return RootSquash
	case "all":
		return AllSquash
	case "rootid":
		return RootIDSquash
	default:
		return Unspecifiedquash
	}
}

func toCephSquashString(s file.Squash) string {
	switch strings.ToLower(strings.TrimSpace(string(s))) {
	case "none":
		return "None"
	case "root":
		return "Root"
	case "all":
		return "All"
	case "rootid":
		return "RootId"
	default:
		return ""
	}
}
