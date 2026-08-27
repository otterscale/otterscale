package core

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
)

// enrolmentTokenContext namespaces the derivation so a token can never
// be mistaken for, or reused as, some other value derived from the same
// secret.
const enrolmentTokenContext = "otterscale-enrolment:" //nolint:gosec // a domain separator, not a credential

// EnrolmentToken is the token an agent presents when registering. It
// is a distinct type so that Wire can tell it apart from other strings
// when injecting dependencies.
type EnrolmentToken string

// Enrolment issues and verifies the tokens that authorize an agent to
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
type Enrolment struct {
	secret []byte
}

// NewEnrolment returns an Enrolment backed by the given root secret.
// The secret is required: without one the registration endpoint, which
// is reachable without authentication, would accept any caller.
func NewEnrolment(secret string) (*Enrolment, error) {
	if secret == "" {
		return nil, errors.New("core: enrolment secret must not be empty")
	}
	return &Enrolment{secret: []byte(secret)}, nil
}

// Token returns the enrolment token for the given cluster.
func (e *Enrolment) Token(cluster string) string {
	return base64.RawURLEncoding.EncodeToString(e.expected(cluster))
}

// Verify reports whether token authorizes registering cluster.
func (e *Enrolment) Verify(cluster, token string) error {
	invalid := &DomainError{
		Code: ErrorCodeUnauthenticated,
		// Deliberately uniform: the caller learns that the token was
		// not accepted, not whether the cluster is already registered
		// or how far the token was from correct.
		Message: "invalid enrolment token",
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
func (e *Enrolment) expected(cluster string) []byte {
	mac := hmac.New(sha256.New, e.secret)
	// hash.Hash.Write is documented never to return an error.
	_, _ = mac.Write([]byte(enrolmentTokenContext))
	_, _ = mac.Write([]byte(cluster))
	return mac.Sum(nil)
}
