package cache

import (
	"context"
	"sync"
	"time"

	"github.com/otterscale/otterscale/internal/core"
)

// maxAgentValuesTickets bounds outstanding URLs. Issuing is admin-only and one
// per joining cluster, so this is far above real use.
const maxAgentValuesTickets = 1024

// AgentValuesStore keeps issued agent-values URLs in memory until they expire.
// Expiry needs no background sweeper: Get enforces it and Put clears what has
// lapsed, and the ceiling bounds what an idle server can hold.
//
// In memory rather than in Valkey, which the server does not otherwise talk
// to. Outstanding URLs therefore do not survive a restart, which is tolerable
// because the procedure also returns the file inline and re-issuing produces
// an identical one. core.AgentValuesStore is an interface so a Valkey-backed
// implementation can replace this without the domain layer changing.
type AgentValuesStore struct {
	now func() time.Time

	mu      sync.Mutex
	tickets map[string]agentValuesTicket
}

type agentValuesTicket struct {
	request   *core.AgentValuesRequest
	expiresAt time.Time
}

var _ core.AgentValuesStore = (*AgentValuesStore)(nil)

func NewAgentValuesStore() *AgentValuesStore {
	return &AgentValuesStore{
		now:     time.Now,
		tickets: make(map[string]agentValuesTicket),
	}
}

// Put sweeps expired entries first, so an old burst of short-lived URLs
// cannot hold a steady-state deployment at the ceiling.
func (s *AgentValuesStore) Put(
	_ context.Context, id string, req *core.AgentValuesRequest, expiresAt time.Time,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.evictExpiredLocked()

	if len(s.tickets) >= maxAgentValuesTickets {
		return &core.DomainError{
			Code: core.ErrorCodeResourceExhausted,
			Message: "too many agent values URLs are outstanding; " +
				"they expire on their own, so retry shortly",
		}
	}

	s.tickets[id] = agentValuesTicket{request: req, expiresAt: expiresAt}
	return nil
}

// Get enforces the TTL by the clock, and drops an expired entry on the way
// out, so nothing has to sweep behind it.
func (s *AgentValuesStore) Get(_ context.Context, id string) (*core.AgentValuesRequest, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ticket, ok := s.tickets[id]
	if !ok {
		return nil, core.ErrAgentValuesTicketNotFound
	}
	if s.now().After(ticket.expiresAt) {
		delete(s.tickets, id)
		return nil, core.ErrAgentValuesTicketNotFound
	}
	return ticket.request, nil
}

func (s *AgentValuesStore) evictExpiredLocked() {
	now := s.now()
	for id, ticket := range s.tickets {
		if now.After(ticket.expiresAt) {
			delete(s.tickets, id)
		}
	}
}
