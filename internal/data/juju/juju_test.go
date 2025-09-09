package juju

import (
	"sync"
	"testing"

	"github.com/juju/juju/api"
	"github.com/juju/juju/api/connector" // <-- provides HostPort used by api.Connection
	"github.com/juju/juju/core/network"
	"github.com/otterscale/otterscale/internal/config"
)

/* -------------------------------------------------------------------------- *
 * Dummy connection – used for the “CacheHit” test
 * -------------------------------------------------------------------------- */
type dummyConn struct {
	closed bool
}

/* Methods required to satisfy api.Connection */
func (d *dummyConn) Close()                  { d.closed = true }
func (d *dummyConn) IsBroken() bool          { return false }
func (d *dummyConn) Ping() error             { return nil }
func (d *dummyConn) ModelTag() string        { return "" }
func (d *dummyConn) APIHost() string         { return "" }
func (d *dummyConn) APIInfo() *api.Info      { return nil }
func (d *dummyConn) IsClosed() bool          { return false }
func (d *dummyConn) IsModelAuthorized() bool { return false }
func (d *dummyConn) IsModelController() bool { return false }
func (d *dummyConn) Name() string            { return "" }
func (d *dummyConn) Provider() string        { return "" }
func (d *dummyConn) CloseIdleConnections()   {}
func (d *dummyConn) AddObserver(func())      {}
func (d *dummyConn) RemoveObserver(func())   {}

// Required by the current api.Connection interface.
func (d *dummyConn) APIHostPorts() ([]network.HostPort, error) { return nil, nil }

func (d *dummyConn) APICall(string, int, string, string, interface{}, interface{}) error {
	return nil
}

/* -------------------------------------------------------------------------- *
 * Broken connection – used for the “BrokenCache” test
 * -------------------------------------------------------------------------- */
type brokenConn struct {
	closed bool
}

/* Methods required to satisfy api.Connection */
func (b *brokenConn) Close()                  { b.closed = true }
func (b *brokenConn) IsBroken() bool          { return true }
func (b *brokenConn) Ping() error             { return nil }
func (b *brokenConn) ModelTag() string        { return "" }
func (b *brokenConn) APIHost() string         { return "" }
func (b *brokenConn) APIInfo() *api.Info      { return nil }
func (b *brokenConn) IsClosed() bool          { return false }
func (b *brokenConn) IsModelAuthorized() bool { return false }
func (b *brokenConn) IsModelController() bool { return false }
func (b *brokenConn) Name() string            { return "" }
func (b *brokenConn) Provider() string        { return "" }
func (b *brokenConn) CloseIdleConnections()   {}
func (b *brokenConn) AddObserver(func())      {}
func (b *brokenConn) RemoveObserver(func())   {}

// Required by the current api.Connection interface.
func (b *brokenConn) APIHostPorts() ([]network.HostPorts, error) { return nil, nil }

func (b *brokenConn) APICall(string, int, string, string, interface{}, interface{}) error {
	return nil
}

/* -------------------------------------------------------------------------- *
 * Helper – works for both normal tests (testing.T) and benchmarks (testing.B)
 * -------------------------------------------------------------------------- */
func juju_mustNew(t testing.TB, cfg *config.Config) *Juju {
	j := New(cfg)

	return j
}

/* -------------------------------------------------------------------------- *
 * Construction – happy path
 * -------------------------------------------------------------------------- */
func TestNewJuju(t *testing.T) {
	cfg := &config.Config{}
	j := juju_mustNew(t, cfg)

	if j == nil {
		t.Fatal("New returned nil")
	}
	if j.conf != cfg {
		t.Error("j.conf does not point to the supplied config")
	}
}

func TestJuju_Accessors(t *testing.T) {
	cfg := &config.Config{
		Juju: config.Juju{
			Username:            "admin",
			CloudName:           "aws",
			CloudRegion:         "us-east-1",
			CharmhubAPIURL:      "https://api.charmhub.io",
			ControllerAddresses: []string{"10.0.0.1:17070"},
		},
	}
	j := juju_mustNew(t, cfg)

	if got, want := j.username(), cfg.Juju.Username; got != want {
		t.Errorf("username(): got %q, want %q", got, want)
	}
	if got, want := j.cloudName(), cfg.Juju.CloudName; got != want {
		t.Errorf("cloudName(): got %q, want %q", got, want)
	}
	if got, want := j.cloudRegion(), cfg.Juju.CloudRegion; got != want {
		t.Errorf("cloudRegion(): got %q, want %q", got, want)
	}
	if got, want := j.charmhubAPIURL(), cfg.Juju.CharmhubAPIURL; got != want {
		t.Errorf("charmhubAPIURL(): got %q, want %q", got, want)
	}
}

func TestJuju_Connection_Error(t *testing.T) {
	j := juju_mustNew(t, &config.Config{})

	_, err := j.connection("some-uuid")
	if err == nil {
		t.Fatalf("expected error from connection with empty config, got nil")
	}
}

/*
func TestJuju_Connection_BrokenCache(t *testing.T) {
	j := juju_mustNew(t, &config.Config{})

	uuid := "broken-uuid"
	bc := &brokenConn{}
	j.connections.Store(uuid, bc)

	_, err := j.connection(uuid)
	if err == nil {
		t.Fatalf("expected error when creating a new connection after a broken cache")
	}
	if !bc.closed {
		t.Fatalf("expected the broken connection to be closed")
	}
}
*/

func TestJuju_ConcurrentAccess(t *testing.T) {
	j := juju_mustNew(t, &config.Config{})
	const workers = 10
	uuid := "concurrent"

	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			_, _ = j.connection(uuid)
		}()
	}
	wg.Wait()
}

func BenchmarkJuju_Creation(b *testing.B) {
	cfg := &config.Config{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if New(cfg) == nil {
			b.Fatal("constructor returned nil")
		}
	}
}

/* -------------------------------------------------------------------------- *
 * Keep the imported packages from production code alive.
 * -------------------------------------------------------------------------- */
var _ = connector.SimpleConfig{}
