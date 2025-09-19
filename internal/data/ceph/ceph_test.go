package ceph

import (
	"testing"

	"github.com/ceph/go-ceph/rados"
	"github.com/ceph/go-ceph/rgw/admin"
	"github.com/otterscale/otterscale/internal/config"
	"github.com/otterscale/otterscale/internal/core"
)

/* -------------------------------------------------------------------------- *
 * Helper – builds a minimal *core.StorageConfig suitable for the tests.
 * -------------------------------------------------------------------------- */
func newStorageConfig(fsid string, monHost string, endpoint string) *core.StorageConfig {
	return &core.StorageConfig{
		StorageCephConfig: &core.StorageCephConfig{
			FSID:    fsid,
			MONHost: monHost,
			Key:     "dummy-key",
		},
		StorageRGWConfig: &core.StorageRGWConfig{
			Endpoint:  endpoint,
			AccessKey: "dummy-ak",
			SecretKey: "dummy-sk",
		},
	}
}

/* -------------------------------------------------------------------------- *
 * TestNew – the constructor must return a non‑nil *Ceph instance.
 * -------------------------------------------------------------------------- */
func TestNew(t *testing.T) {
	c := New(&config.Config{})
	if c == nil {
		t.Fatalf("New returned nil")
	}
}

/* -------------------------------------------------------------------------- *
 * TestConnection_Cached – a value already stored in the connections map is
 * returned without invoking newConnection.
 * -------------------------------------------------------------------------- */
func TestConnection_Cached(t *testing.T) {
	// Create a Ceph instance with an empty config (no real connection needed).
	ceph := New(&config.Config{})

	// Insert a dummy *rados.Conn (nil is a valid value for the type assertion).
	const fsid = "cached-fsid"
	ceph.connections.Store(fsid, (*rados.Conn)(nil))

	// Retrieve it via the public method.
	conn, err := ceph.connection(newStorageConfig(fsid, "unused", "unused"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if conn != nil {
		t.Fatalf("expected cached nil connection, got non‑nil")
	}
}

/* -------------------------------------------------------------------------- *
 * TestConnection_NewError – when nothing is cached, newConnection is called.
 * The test uses an obviously invalid MONHost so that rados.Connect fails.
 * The method must return an error and must not cache the failed connection.
 * -------------------------------------------------------------------------- */
func TestConnection_NewError(t *testing.T) {
	ceph := New(&config.Config{})

	cfg := newStorageConfig("new-error-fsid", "", "unused") // empty MONHost → Connect will fail
	_, err := ceph.connection(cfg)
	if err == nil {
		t.Fatalf("expected error from connection, got nil")
	}
	// Ensure nothing was cached after the failure.
	if _, ok := ceph.connections.Load(cfg.FSID); ok {
		t.Fatalf("connection was cached despite error")
	}
}

/* -------------------------------------------------------------------------- *
 * TestClient_Cached – a value already stored in the clients map is returned
 * directly.
 * -------------------------------------------------------------------------- */
func TestClient_Cached(t *testing.T) {
	ceph := New(&config.Config{})

	const fsid = "cached-client"
	ceph.clients.Store(fsid, (*admin.API)(nil))

	client, err := ceph.client(newStorageConfig(fsid, "unused", "unused"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client != nil {
		t.Fatalf("expected cached nil client, got non‑nil")
	}
}

/* -------------------------------------------------------------------------- *
 * TestClient_Error – admin.New should fail when the endpoint is empty.
 * The method must surface the error and must not store a value in the cache.
 * -------------------------------------------------------------------------- */
func TestClient_Error(t *testing.T) {
	ceph := New(&config.Config{})

	cfg := newStorageConfig("error-fsid", "unused", "") // empty Endpoint → admin.New errors
	_, err := ceph.client(cfg)
	if err == nil {
		t.Fatalf("expected error from client creation, got nil")
	}
	// Verify that nothing was cached.
	if _, ok := ceph.clients.Load(cfg.FSID); ok {
		t.Fatalf("client was cached despite error")
	}
}

/* -------------------------------------------------------------------------- *
 * Compile‑time checks – ensure the concrete type still satisfies any
 * expectations (no methods are required here, just a sanity check).
 * -------------------------------------------------------------------------- */
var _ = (*Ceph)(nil)
