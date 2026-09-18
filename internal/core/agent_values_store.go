package core

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"
)

// agentValuesTicketTTL is how long an issued URL is served. A person pastes
// these into a shell, so it is generous in human terms and short in credential
// ones.
const agentValuesTicketTTL = 1 * time.Hour

// agentValuesTicketIDBytes encodes to 22 URL-safe characters, and is far
// beyond guessable within the TTL.
const agentValuesTicketIDBytes = 16

// ErrAgentValuesTicketNotFound is what every lookup failure returns, so a
// caller cannot tell an id that was never issued from one that has expired.
var ErrAgentValuesTicketNotFound = &DomainError{
	Code:    ErrorCodeUnauthenticated,
	Message: "invalid or expired token",
}

// AgentValuesStore holds the parameters behind an issued URL until it expires.
//
// The URL carries an opaque id rather than the parameters, which keeps it
// short enough to paste and keeps the cluster and its addresses out of shell
// history, scrollback and proxy logs. Implementations live in the
// infrastructure layer.
type AgentValuesStore interface {
	// Put stores req under id until expiresAt. An absolute deadline rather
	// than a TTL, so the instant reported to the caller is the one the store
	// enforces, not a second computation of it.
	Put(ctx context.Context, id string, req *AgentValuesRequest, expiresAt time.Time) error
	// Get returns ErrAgentValuesTicketNotFound when there is no live entry.
	// The returned value is the stored one: read it, do not mutate it.
	Get(ctx context.Context, id string) (*AgentValuesRequest, error)
}

// newAgentValuesTicketID returns an unguessable, URL-safe id.
func newAgentValuesTicketID() (string, error) {
	buf := make([]byte, agentValuesTicketIDBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generate agent values ticket id: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}
