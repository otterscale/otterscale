package tunnel

import (
	"context"
	"math/rand/v2"
	"strings"
	"time"
)

// isAuthErr matches on the message because chisel exposes no typed auth
// errors.
func isAuthErr(err error) bool {
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "unable to authenticate") ||
		strings.Contains(msg, "authentication failed") ||
		strings.Contains(msg, "auth failed") ||
		strings.Contains(msg, "unauthorized") ||
		strings.Contains(msg, "invalid auth")
}

// sleepCtx returns true if the full delay elapsed, false if ctx ended it.
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

type backoff struct {
	base    time.Duration
	max     time.Duration
	current time.Duration
}

func newBackoff(base, maxInterval time.Duration) *backoff {
	return &backoff{base: base, max: maxInterval, current: base}
}

// Next returns a jittered delay, then doubles the interval. Full jitter —
// uniform random in [0, current] — is what keeps agents from reconnecting in
// lockstep after a server restart.
func (b *backoff) Next() time.Duration {
	d := b.current
	jittered := time.Duration(rand.Int64N(int64(d) + 1)) //nolint:gosec // weak random is intentional for jitter
	if next := b.current * 2; next > b.max {
		b.current = b.max
	} else {
		b.current = next
	}
	return jittered
}

func (b *backoff) Reset() {
	b.current = b.base
}
