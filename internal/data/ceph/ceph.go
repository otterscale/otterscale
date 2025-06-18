package ceph

import (
	"strconv"
	"sync"
	"time"

	"github.com/ceph/go-ceph/rados"

	"github.com/openhdc/otterscale/internal/config"
	"github.com/openhdc/otterscale/internal/core"
)

type Ceph struct {
	conf        *config.Config
	connections sync.Map
}

func New(conf *config.Config) *Ceph {
	return &Ceph{
		conf: conf,
	}
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

	radosTimeout := time.Second * 3
	if m.conf.Ceph.RadosTimeout > 0 {
		radosTimeout = m.conf.Ceph.RadosTimeout
	}
	timeout := strconv.FormatFloat(radosTimeout.Seconds(), 'f', -1, 64)
	if err = conn.SetConfigOption("rados_mon_op_timeout", timeout); err != nil {
		return nil, err
	}
	if err = conn.SetConfigOption("rados_osd_op_timeout", timeout); err != nil {
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
