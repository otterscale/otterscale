package cache

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/otterscale/otterscale/internal/core"
)

const testTicketTTL = time.Hour

func testRequest() *core.AgentValuesRequest {
	return &core.AgentValuesRequest{
		Cluster: "prod",
		Subject: "admin-1",
		ClusterInfo: &core.AgentClusterInfo{
			ExternalAddress: "192.0.2.1",
			NodePortRange:   "30000-32767",
		},
	}
}

// newTestStore lets the test drive the clock, so expiry needs no sleeping.
func newTestStore(now *time.Time) *AgentValuesStore {
	s := NewAgentValuesStore()
	s.now = func() time.Time { return *now }
	return s
}

func TestAgentValuesStore_RoundTrip(t *testing.T) {
	now := time.Date(2026, time.September, 18, 12, 0, 0, 0, time.UTC)
	s := newTestStore(&now)

	want := testRequest()
	if err := s.Put(t.Context(), "ticket-1", want, testTicketTTL); err != nil {
		t.Fatalf("Put() error = %v", err)
	}

	got, err := s.Get(t.Context(), "ticket-1")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if got != want {
		t.Errorf("Get() returned %+v, want the stored %+v", got, want)
	}
}

func TestAgentValuesStore_GetUnknown(t *testing.T) {
	now := time.Now()
	s := newTestStore(&now)

	_, err := s.Get(t.Context(), "never-issued")
	if !errors.Is(err, core.ErrAgentValuesTicketNotFound) {
		t.Errorf("error = %v, want %v", err, core.ErrAgentValuesTicketNotFound)
	}
}

// The clock enforces the TTL, not the sweep: an entry the loop has not reached
// must already be gone.
func TestAgentValuesStore_Expiry(t *testing.T) {
	now := time.Date(2026, time.September, 18, 12, 0, 0, 0, time.UTC)
	s := newTestStore(&now)

	if err := s.Put(t.Context(), "ticket-1", testRequest(), testTicketTTL); err != nil {
		t.Fatalf("Put() error = %v", err)
	}

	now = now.Add(testTicketTTL - time.Minute)
	if _, err := s.Get(t.Context(), "ticket-1"); err != nil {
		t.Errorf("Get() inside the TTL failed: %v", err)
	}

	now = now.Add(2 * time.Minute)
	_, err := s.Get(t.Context(), "ticket-1")
	if !errors.Is(err, core.ErrAgentValuesTicketNotFound) {
		t.Errorf("error after expiry = %v, want %v", err, core.ErrAgentValuesTicketNotFound)
	}
}

func TestAgentValuesStore_Capacity(t *testing.T) {
	now := time.Now()
	s := newTestStore(&now)

	for i := range maxAgentValuesTickets {
		id := fmt.Sprintf("ticket-%d", i)
		if err := s.Put(t.Context(), id, testRequest(), testTicketTTL); err != nil {
			t.Fatalf("Put() %d error = %v", i, err)
		}
	}

	err := s.Put(t.Context(), "one-too-many", testRequest(), testTicketTTL)
	if err == nil {
		t.Fatal("expected an error past the ceiling, got nil")
	}
	if code, ok := core.DomainErrorCode(err); !ok || code != core.ErrorCodeResourceExhausted {
		t.Errorf("code = %v (domain=%t), want %v", code, ok, core.ErrorCodeResourceExhausted)
	}
}

// Why the ceiling is not a permanent wall: URLs issued long ago must not block
// one now.
func TestAgentValuesStore_PutSweepsExpired(t *testing.T) {
	now := time.Date(2026, time.September, 18, 12, 0, 0, 0, time.UTC)
	s := newTestStore(&now)

	for i := range maxAgentValuesTickets {
		id := fmt.Sprintf("ticket-%d", i)
		if err := s.Put(t.Context(), id, testRequest(), testTicketTTL); err != nil {
			t.Fatalf("Put() %d error = %v", i, err)
		}
	}

	now = now.Add(testTicketTTL + time.Minute)
	if err := s.Put(t.Context(), "fresh", testRequest(), testTicketTTL); err != nil {
		t.Fatalf("Put() after the old tickets expired failed: %v", err)
	}
	if got := len(s.tickets); got != 1 {
		t.Errorf("stored tickets = %d, want 1: the expired ones were not swept", got)
	}
}

// A rejected fetch reclaims what it rejected, so the common path needs no
// sweeper behind it.
func TestAgentValuesStore_GetReclaimsExpired(t *testing.T) {
	now := time.Date(2026, time.September, 18, 12, 0, 0, 0, time.UTC)
	s := newTestStore(&now)

	if err := s.Put(t.Context(), "stale", testRequest(), testTicketTTL); err != nil {
		t.Fatalf("Put() error = %v", err)
	}
	now = now.Add(testTicketTTL + time.Minute)

	if _, err := s.Get(t.Context(), "stale"); !errors.Is(err, core.ErrAgentValuesTicketNotFound) {
		t.Fatalf("error = %v, want %v", err, core.ErrAgentValuesTicketNotFound)
	}
	if got := len(s.tickets); got != 0 {
		t.Errorf("stored tickets = %d, want 0: the expired one was not dropped", got)
	}
}
