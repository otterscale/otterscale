package core

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
)

// joinTokenContext namespaces the derivation so a token can never
// be mistaken for, or reused as, some other value derived from the same
// secret.
const joinTokenContext = "otterscale-join:" //nolint:gosec // a domain separator, not a credential

// JoinToken is what an agent presents when registering. It is a distinct
// type so Wire can tell it apart from other strings.
type JoinToken string

// JoinAuthority issues and verifies the tokens that authorize an agent to
// register a cluster.
//
// Tokens are derived rather than stored: the server holds one root
// secret and computes the expected token for whichever cluster a
// request claims, so onboarding a cluster needs no server-side change.
// Deriving per cluster also contains a leak — a token authorizes the
// one cluster it names, and an agent holding it cannot compute a token
// for any other cluster.
//
// Tokens do not expire and cannot be revoked individually; rotating the
// root secret invalidates all of them at once.
type JoinAuthority struct {
	secret []byte
}

// NewJoinAuthority requires a secret: without one the registration endpoint, which
// is reachable without authentication, would accept any caller.
func NewJoinAuthority(secret string) (*JoinAuthority, error) {
	if secret == "" {
		return nil, errors.New("core: join secret must not be empty")
	}
	return &JoinAuthority{secret: []byte(secret)}, nil
}

func (e *JoinAuthority) Token(cluster string) string {
	return base64.RawURLEncoding.EncodeToString(e.expected(cluster))
}

// Verify reports whether token authorizes registering cluster.
func (e *JoinAuthority) Verify(cluster, token string) error {
	invalid := &DomainError{
		Code: ErrorCodeUnauthenticated,
		// Deliberately uniform: the caller learns that the token was
		// not accepted, not whether the cluster is already registered
		// or how far the token was from correct.
		Message: "invalid join token",
	}

	got, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil || len(got) != sha256.Size {
		return invalid
	}

	if !hmac.Equal(got, e.expected(cluster)) {
		return invalid
	}
	return nil
}

// expected computes the raw MAC for a cluster. The context ends in a
// colon and ValidateClusterName rejects colons, so no cluster name can
// be confused with another by shifting the boundary.
func (e *JoinAuthority) expected(cluster string) []byte {
	mac := hmac.New(sha256.New, e.secret)
	// hash.Hash.Write is documented never to return an error.
	_, _ = mac.Write([]byte(joinTokenContext))
	_, _ = mac.Write([]byte(cluster))
	return mac.Sum(nil)
}
