package core

import (
	"encoding/base64"
	"strings"
	"testing"
)

func newTestJoinAuthority(t *testing.T, secret string) *JoinAuthority {
	t.Helper()

	e, err := NewJoinAuthority(secret)
	if err != nil {
		t.Fatalf("NewJoinAuthority: %v", err)
	}
	return e
}

func TestNewJoinAuthority_RequiresSecret(t *testing.T) {
	t.Parallel()

	if _, err := NewJoinAuthority(""); err == nil {
		t.Error("expected an error for an empty secret, got nil")
	}
}

func TestJoinAuthority_TokenIsStable(t *testing.T) {
	t.Parallel()

	e := newTestJoinAuthority(t, "root-secret")

	first := e.Token("prod")
	if first != e.Token("prod") {
		t.Error("the same cluster produced two different tokens")
	}
	if first != newTestJoinAuthority(t, "root-secret").Token("prod") {
		t.Error("a second JoinAuthority with the same secret produced a different token")
	}

	// 32 raw bytes, unpadded base64url.
	if decoded, err := base64.RawURLEncoding.DecodeString(first); err != nil || len(decoded) != 32 {
		t.Errorf("token %q does not decode to 32 raw bytes (err=%v)", first, err)
	}
}

// TestJoinAuthority_TokenIsPerCluster is the property the whole scheme
// rests on: an agent holding one cluster's token cannot register under
// another cluster's name.
func TestJoinAuthority_TokenIsPerCluster(t *testing.T) {
	t.Parallel()

	e := newTestJoinAuthority(t, "root-secret")

	if e.Token("prod") == e.Token("staging") {
		t.Fatal("two clusters share a token")
	}
	if err := e.Verify("staging", e.Token("prod")); err == nil {
		t.Error("prod's token was accepted for staging")
	}
}

func TestJoinAuthority_TokenDependsOnSecret(t *testing.T) {
	t.Parallel()

	token := newTestJoinAuthority(t, "root-secret").Token("prod")

	// Rotating the secret must invalidate every token issued before.
	if err := newTestJoinAuthority(t, "rotated-secret").Verify("prod", token); err == nil {
		t.Error("a token issued under the previous secret is still accepted")
	}
}

func TestJoinAuthority_Verify(t *testing.T) {
	t.Parallel()

	e := newTestJoinAuthority(t, "root-secret")
	valid := e.Token("prod")

	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{name: "the cluster's own token", token: valid},
		{name: "empty", token: "", wantErr: true},
		{name: "not base64", token: "not a token", wantErr: true},
		{name: "right encoding, wrong length", token: base64.RawURLEncoding.EncodeToString([]byte("short")), wantErr: true},
		{name: "one character changed", token: flipFirstChar(valid), wantErr: true},
		{name: "correct token with padding", token: valid + "=", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := e.Verify("prod", tt.token)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected an error, got nil")
				}
				code, ok := DomainErrorCode(err)
				if !ok || code != ErrorCodeUnauthenticated {
					t.Errorf("code = %v (domain=%v), want ErrorCodeUnauthenticated", code, ok)
				}
				// The message must not distinguish between failure
				// modes, or it becomes an oracle.
				if got := err.Error(); !strings.Contains(got, "invalid join token") {
					t.Errorf("message = %q, want the uniform rejection message", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

// flipFirstChar returns s with its first character changed.
func flipFirstChar(s string) string {
	if s == "" {
		return s
	}
	replacement := "A"
	if s[:1] == replacement {
		replacement = "B"
	}
	return replacement + s[1:]
}
