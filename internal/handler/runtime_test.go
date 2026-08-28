package handler

import (
	"context"
	"errors"
	"testing"

	"connectrpc.com/connect"

	"github.com/otterscale/otterscale/internal/core"
)

// TestSessionOutcome covers the decision a streaming handler makes once
// its session has finished. Before this existed, every ending was
// reported as success: a session that never started looked exactly like
// one that produced no output.
func TestSessionOutcome(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want connect.Code // 0 means "no error"
	}{
		{
			name: "clean end",
			err:  nil,
		},
		{
			name: "caller canceled",
			err:  context.Canceled,
		},
		{
			name: "caller deadline",
			err:  context.DeadlineExceeded,
		},
		{
			// The command ran and its output already reached the client.
			// Reporting the exit status as a transport failure would make
			// every unsuccessful shell command look like a server error.
			name: "command exited non-zero",
			err:  &core.ErrCommandExited{Code: 1},
		},
		{
			name: "container missing",
			err:  &core.DomainError{Code: core.ErrorCodeNotFound, Message: "no such container"},
			want: connect.CodeNotFound,
		},
		{
			name: "kubelet refused the forward",
			err:  &core.DomainError{Code: core.ErrorCodeFailedPrecondition, Message: "socat not found"},
			want: connect.CodeFailedPrecondition,
		},
		{
			name: "rbac denial",
			err:  &core.DomainError{Code: core.ErrorCodePermissionDenied, Message: "forbidden"},
			want: connect.CodePermissionDenied,
		},
		{
			name: "unrecognized failure",
			err:  errors.New("stream reset"),
			want: connect.CodeInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sessionOutcome("exec", "session-1", tt.err)

			if tt.want == 0 {
				if got != nil {
					t.Fatalf("sessionOutcome = %v, want nil", got)
				}
				return
			}

			if got == nil {
				t.Fatalf("sessionOutcome = nil, want code %v", tt.want)
			}
			if code := connect.CodeOf(got); code != tt.want {
				t.Fatalf("code = %v, want %v", code, tt.want)
			}
		})
	}
}

// TestSessionOutcomeUnwrapsWrappedCancellation checks that a
// cancellation buried under an adapter's wrapping is still recognized —
// the caller is gone either way, and there is nobody to report to.
func TestSessionOutcomeUnwrapsWrappedCancellation(t *testing.T) {
	wrapped := &core.DomainError{
		Code:    core.ErrorCodeCanceled,
		Message: "exec stream",
		Cause:   context.Canceled,
	}

	if got := sessionOutcome("exec", "session-1", wrapped); got != nil {
		t.Fatalf("sessionOutcome = %v, want nil", got)
	}
}
