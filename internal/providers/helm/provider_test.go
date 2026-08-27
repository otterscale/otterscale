package helm

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"helm.sh/helm/v4/pkg/helmpath"

	"github.com/otterscale/otterscale/internal/core"
)

// awaitTimeout bounds the cancellation assertions: the point of
// awaitWithContext is that it returns without waiting for fn, so a
// blocked return must fail the test rather than hang it.
const awaitTimeout = 2 * time.Second

func TestAwaitWithContext_ReturnsResult(t *testing.T) {
	t.Parallel()

	got, err := awaitWithContext(t.Context(), "doing a thing", func() (int, error) {
		return 42, nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 42 {
		t.Errorf("got %d, want 42", got)
	}
}

func TestAwaitWithContext_PropagatesError(t *testing.T) {
	t.Parallel()

	boom := errors.New("boom")
	if _, err := awaitWithContext(t.Context(), "doing a thing", func() (int, error) {
		return 0, boom
	}); !errors.Is(err, boom) {
		t.Fatalf("got %v, want %v", err, boom)
	}
}

// TestAwaitWithContext_GivesUp is the regression test for the dropped
// context: a fetch the Helm SDK cannot cancel must not pin the caller
// to it.
func TestAwaitWithContext_GivesUp(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		ctx      func(t *testing.T) context.Context
		wantCode core.ErrorCode
	}{
		{
			name: "caller cancels",
			ctx: func(t *testing.T) context.Context {
				t.Helper()
				ctx, cancel := context.WithCancel(t.Context())
				cancel()
				return ctx
			},
			wantCode: core.ErrorCodeCanceled,
		},
		{
			name: "deadline passes",
			ctx: func(t *testing.T) context.Context {
				t.Helper()
				ctx, cancel := context.WithTimeout(t.Context(), time.Millisecond)
				t.Cleanup(cancel)
				return ctx
			},
			wantCode: core.ErrorCodeDeadlineExceeded,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// fn outlives the call, exactly like a Helm fetch that
			// cannot be interrupted.
			release := make(chan struct{})
			t.Cleanup(func() { close(release) })

			type outcome struct{ err error }
			done := make(chan outcome, 1)
			go func() {
				_, err := awaitWithContext(tt.ctx(t), "fetching chart example/nginx", func() (int, error) {
					<-release
					return 0, nil
				})
				done <- outcome{err: err}
			}()

			select {
			case got := <-done:
				code, ok := core.DomainErrorCode(got.err)
				if !ok {
					t.Fatalf("got %v (%T), want a *core.DomainError", got.err, got.err)
				}
				if code != tt.wantCode {
					t.Errorf("code = %v, want %v", code, tt.wantCode)
				}
			case <-time.After(awaitTimeout):
				t.Fatalf("awaitWithContext did not return within %s", awaitTimeout)
			}
		})
	}
}

func TestSetHelmHome_CreatesWritableDirectories(t *testing.T) {
	// t.Setenv forbids parallel tests; it also restores the previous
	// values, so the process-wide change does not leak to other tests.
	t.Setenv(helmpath.CacheHomeEnvVar, "")
	t.Setenv(helmpath.ConfigHomeEnvVar, "")
	t.Setenv(helmpath.DataHomeEnvVar, "")

	if err := setHelmHome(); err != nil {
		t.Fatalf("setHelmHome: %v", err)
	}

	for _, sub := range []struct{ env, dir string }{
		{helmpath.CacheHomeEnvVar, "cache"},
		{helmpath.ConfigHomeEnvVar, "config"},
		{helmpath.DataHomeEnvVar, "data"},
	} {
		want := filepath.Join(helmHome(), sub.dir)
		if got := os.Getenv(sub.env); got != want {
			t.Errorf("%s = %q, want %q", sub.env, got, want)
		}
		info, err := os.Stat(want)
		if err != nil {
			t.Errorf("stat %s: %v", want, err)
			continue
		}
		if !info.IsDir() {
			t.Errorf("%s is not a directory", want)
		}
	}
}

func TestSetHelmHome_KeepsOperatorConfiguration(t *testing.T) {
	configured := t.TempDir()
	t.Setenv(helmpath.CacheHomeEnvVar, configured)
	t.Setenv(helmpath.ConfigHomeEnvVar, "")
	t.Setenv(helmpath.DataHomeEnvVar, "")

	if err := setHelmHome(); err != nil {
		t.Fatalf("setHelmHome: %v", err)
	}

	if got := os.Getenv(helmpath.CacheHomeEnvVar); got != configured {
		t.Errorf("%s = %q, want the operator's %q", helmpath.CacheHomeEnvVar, got, configured)
	}
	if got := os.Getenv(helmpath.DataHomeEnvVar); got == "" {
		t.Errorf("%s was left unset", helmpath.DataHomeEnvVar)
	}
}
