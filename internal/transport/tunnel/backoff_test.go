package tunnel

import (
	"testing"
	"time"
)

func TestBackoff_GrowsAndCaps(t *testing.T) {
	t.Parallel()

	const (
		base = 100 * time.Millisecond
		maxi = 400 * time.Millisecond
	)
	bo := newBackoff(base, maxi)

	// Next returns a jittered delay in [0, current] and doubles the
	// interval, so the bound is what can be asserted.
	for _, want := range []time.Duration{base, 2 * base, maxi, maxi} {
		if got := bo.Next(); got < 0 || got > want {
			t.Fatalf("Next() = %v, want a value in [0, %v]", got, want)
		}
	}

	if bo.current != maxi {
		t.Errorf("current = %v, want it capped at %v", bo.current, maxi)
	}
}

func TestBackoff_Reset(t *testing.T) {
	t.Parallel()

	const base = 100 * time.Millisecond
	bo := newBackoff(base, time.Second)

	bo.Next()
	bo.Next()
	bo.Reset()

	if bo.current != base {
		t.Errorf("current = %v, want %v after Reset", bo.current, base)
	}
}

func TestIsAuthErr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		msg  string
		want bool
	}{
		{msg: "Unable to authenticate", want: true},
		{msg: "server returned 401 Unauthorized", want: true},
		{msg: "dial tcp 127.0.0.1:8300: connect: connection refused", want: false},
		{msg: "EOF", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.msg, func(t *testing.T) {
			if got := isAuthErr(errString(tt.msg)); got != tt.want {
				t.Errorf("isAuthErr(%q) = %v, want %v", tt.msg, got, tt.want)
			}
		})
	}
}

// errString is a minimal error whose message is the string itself.
type errString string

func (e errString) Error() string { return string(e) }
