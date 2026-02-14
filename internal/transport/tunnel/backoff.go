package tunnel

import (
	"context"
	"math/rand/v2"
	"strings"
	"time"
)

// isAuthErr detects authentication-related errors from chisel by
// inspecting the error message. This is necessary because chisel does
// not expose typed errors for auth failures.
func isAuthErr(err error) bool {
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "unable to authenticate") ||
		strings.Contains(msg, "authentication failed") ||
		strings.Contains(msg, "auth failed") ||
		strings.Contains(msg, "unauthorized") ||
		strings.Contains(msg, "invalid auth")
}

// sleepCtx blocks for d or until ctx is done.
// Returns true if the sleep completed (context still alive).
func sleepCtx(ctx context.Context, d time.Duration) bool {
	t := time.NewTimer(d)
	defer t.Stop()

	select {
	case <-t.C:
		return true
	case <-ctx.Done():
		return false
	}
}

// backoff implements simple exponential backoff capped at a maximum.
type backoff struct {
	base    time.Duration
	max     time.Duration
	current time.Duration
}

func newBackoff(base, max time.Duration) *backoff {
	return &backoff{base: base, max: max, current: base}
}

// Next returns a jittered delay based on the current backoff interval,
// then doubles the interval for the next call. Full jitter (uniform
// random between 0 and current) prevents thundering-herd effects when
// multiple agents reconnect simultaneously after a server restart.
func (b *backoff) Next() time.Duration {
	d := b.current
	// Full jitter: uniform random in [0, current].
	jittered := time.Duration(rand.Int64N(int64(d) + 1))
	if next := b.current * 2; next > b.max {
		b.current = b.max
	} else {
		b.current = next
	}
	return jittered
}

// Reset sets the delay back to the base value.
func (b *backoff) Reset() {
	b.current = b.base
}
