package ceph

type nfsExportRepo struct {
	ceph *Ceph
}

/*func NewNFSExportRepo(ceph *Ceph) file.NFSExportRepo {
	return &nfsExportRepo{
		ceph: ceph,
	}
}

var _ file.NFSExportRepo = (*nfsExportRepo)(nil)*/

/*func (r *nfsExportRepo) List(_ context.Context, scope, clusterID string) ([]file.NFSExport, error) {
	conn, err := r.ceph.connection(scope)
	if err != nil {
		return nil, err
	}

	exports, err := listDetailedExports(conn, clusterID)
	if err != nil {
		return nil, err
	}

	out := make([]file.NFSExport, 0, len(exports))
	for i := range exports {
		e := r.toNFSExport(&exports[i])
		if e == nil {
			continue
		}
		out = append(out, *e)
	}

	return out, nil
}

func (r *nfsExportRepo) Create(
	_ context.Context,
	scope, clusterID, filesystem string,
	enableSecurityLabel bool,
	group, subvolume string,
	protocolNFSv3, protocolNFSv4 bool,
	pseudoPath string,
	accessType file.AccessType,
	sq file.Squash,
	transportUDP, transportTCP bool,
	clients []file.Client,
) error {
	conn, err := r.ceph.connection(scope)
	if err != nil {
		return err
	}

	spec := r.toCephFSExportSpec(
		clusterID,
		filesystem,
		pseudoPath,
		accessType,
		sq,
		protocolNFSv3,
		protocolNFSv4,
		transportUDP,
		transportTCP,
		enableSecurityLabel,
		clients,
	)

	_, err = createCephFSExport(conn, spec)
	return err
}

func (r *nfsExportRepo) Get(_ context.Context, scope, clusterID, pseudoPath string) (*file.NFSExport, error) {
	conn, err := r.ceph.connection(scope)
	if err != nil {
		return nil, err
	}

	info, err := exportInfo(conn, clusterID, pseudoPath)
	if err != nil {
		return nil, err
	}

	return r.toNFSExport(info), nil
}

func (r *nfsExportRepo) Update(
	_ context.Context,
	scope, clusterID string,
	enableSecurityLabel bool,
	protocolNFSv3, protocolNFSv4 bool,
	pseudoPath string,
	accessType file.AccessType,
	sq file.Squash,
	transportUDP, transportTCP bool,
	clients []file.Client,
) error {
	conn, err := r.ceph.connection(scope)
	if err != nil {
		return err
	}

	oldInfo, err := exportInfo(conn, clusterID, pseudoPath)
	if err != nil {
		return err
	}

	if err := removeExport(conn, clusterID, pseudoPath); err != nil {
		return err
	}

	spec := r.toCephFSExportSpec(
		clusterID,
		oldInfo.FSAL.FileSystemName,
		pseudoPath,
		accessType,
		sq,
		protocolNFSv3,
		protocolNFSv4,
		transportUDP,
		transportTCP,
		enableSecurityLabel,
		clients,
	)

	_, err = createCephFSExport(conn, spec)
	return err
}

func (r *nfsExportRepo) Delete(_ context.Context, scope, clusterID, pseudoPath string) error {
	conn, err := r.ceph.connection(scope)
	if err != nil {
		return err
	}
	return removeExport(conn, clusterID, pseudoPath)
}

func (r *nfsExportRepo) toNFSExport(info *ExportInfo) *file.NFSExport {
	if info == nil {
		return nil
	}

	v3, v4 := false, false
	for _, p := range info.Protocols {
		if p == 3 {
			v3 = true
		} else if p == 4 {
			v4 = true
		}
	}

	tcp, udp := false, false
	for _, t := range info.Transports {
		switch strings.ToLower(t) {
		case "tcp":
			tcp = true
		case "udp":
			udp = true
		}
	}

	outClients := make([]file.Client, 0, len(info.Clients))
	for _, c := range info.Clients {
		outClients = append(outClients, file.Client{
			Addresses:  append([]string(nil), c.Addresses...),
			AccessType: file.AccessType(c.AccessType),
			Squash:     file.Squash(c.Squash),
		})
	}

	return &file.NFSExport{
		ClusterID:            info.ClusterID,
		FileSystemName:       info.FSAL.FileSystemName,
		EnableSecurityLabel:  info.SecurityLabel,
		Path:                 info.Path,
		PseudoPath:           info.PseudoPath,
		ProtocolNFSV3:        v3,
		ProtocolNFSV4:        v4,
		TransportProtocolTCP: tcp,
		TransportProtocolUDP: udp,
		AccessType:           file.AccessType(info.AccessType),
		Squash:               file.Squash(info.Squash),
		Clients:              outClients,
	}
}

func (r *nfsExportRepo) toCephFSExportSpec(
	clusterID, filesystem, pseudoPath string,
	accessType file.AccessType,
	sq file.Squash,
	_ bool, _ bool,
	_ bool, _ bool,
	_ bool,
	clients []file.Client,
) CephFSExportSpec {

	clientAddrs := make([]string, 0, len(clients))
	for _, c := range clients {
		clientAddrs = append(clientAddrs, c.Addresses...)
	}

	readOnly := accessType == file.AccessTypeRO

	sec := []SecType{}
	return CephFSExportSpec{
		FileSystemName: filesystem,
		ClusterID:      clusterID,
		PseudoPath:     pseudoPath,
		ReadOnly:       readOnly,
		ClientAddr:     clientAddrs,
		Squash:         SquashMode(sq),
		SecType:        sec,
	}
}*/
