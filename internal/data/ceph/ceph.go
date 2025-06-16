package ceph

import (
	"sync"

	"github.com/ceph/go-ceph/rados"

	"github.com/openhdc/otterscale/internal/core"
)

type Ceph struct {
	connections sync.Map
}

func New() *Ceph {
	return &Ceph{}
}

func (m *Ceph) newConnection(config *core.StorageConfig) (*rados.Conn, error) {
	conn, err := rados.NewConn()
	if err != nil {
		return nil, err
	}

	if err := conn.SetConfigOption("fsid", config.FSID); err != nil {
		return nil, err
	}
	if err := conn.SetConfigOption("mon_host", config.MonHost); err != nil {
		return nil, err
	}
	if err := conn.SetConfigOption("key", config.Key); err != nil {
		return nil, err
	}

	if err := conn.Connect(); err != nil {
		return nil, err
	}
	return conn, nil
}

// TODO: Check connection with 'OpenIOContext'
func (m *Ceph) connection(config *core.StorageConfig) (*rados.Conn, error) {
	if v, ok := m.connections.Load(config.FSID); ok {
		conn := v.(*rados.Conn)
		return conn, nil
	}

	conn, err := m.newConnection(config)
	if err != nil {
		return nil, err
	}

	m.connections.Store(config.FSID, conn)

	return conn, nil
}
