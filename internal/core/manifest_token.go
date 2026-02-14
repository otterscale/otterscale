package core

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

// manifestTokenTTL is the validity period of HMAC-signed manifest
// tokens. After this duration the token expires and a new one must
// be issued via the GetAgentManifest RPC.
const manifestTokenTTL = 1 * time.Hour

// errInvalidToken is the generic error returned for all token
// verification failures. Using a single message prevents attackers
// from inferring the verification stage that failed (e.g. decode vs
// signature vs expiry).
var errInvalidToken = errors.New("invalid or expired token")

// manifestTokenClaims is the JSON payload embedded in manifest tokens.
type manifestTokenClaims struct {
	Sub     string `json:"sub"`
	Cluster string `json:"cluster"`
	Iat     int64  `json:"iat"`
	Exp     int64  `json:"exp"`
}

// ManifestTokenIssuer signs and verifies HMAC-based manifest tokens.
// It is extracted from FleetUseCase to isolate token management as a
// single responsibility, making it easier to swap token formats (e.g.
// JWT, opaque) in the future without modifying the fleet orchestration
// logic.
//
// The now function is injected to decouple from wall-clock time,
// enabling deterministic tests without time.Sleep or reflect hacks.
type ManifestTokenIssuer struct {
	hmacKey []byte
	now     func() time.Time
}

// NewManifestTokenIssuer returns a ManifestTokenIssuer backed by the
// given HMAC key. The key must be non-empty.
func NewManifestTokenIssuer(hmacKey []byte) (*ManifestTokenIssuer, error) {
	if len(hmacKey) == 0 {
		return nil, fmt.Errorf("manifest token issuer: HMAC key is required")
	}
	return &ManifestTokenIssuer{hmacKey: hmacKey, now: time.Now}, nil
}

// Issue creates a signed token containing the user identity, cluster
// name, issued-at, and expiry timestamps.
func (i *ManifestTokenIssuer) Issue(cluster, userName string) (string, error) {
	now := i.now()
	claims := manifestTokenClaims{
		Sub:     userName,
		Cluster: cluster,
		Iat:     now.Unix(),
		Exp:     now.Add(manifestTokenTTL).Unix(),
	}

	payload, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("marshal token claims: %w", err)
	}

	mac := hmac.New(sha256.New, i.hmacKey)
	mac.Write(payload)
	sig := mac.Sum(nil)

	return base64.RawURLEncoding.EncodeToString(payload) + "." +
		base64.RawURLEncoding.EncodeToString(sig), nil
}

// Verify validates the HMAC signature and expiry of a manifest token
// and returns the embedded cluster name and user identity. All
// verification failures return a generic error to avoid leaking which
// stage failed; detailed reasons are available via VerifyDetailed.
func (i *ManifestTokenIssuer) Verify(token string) (cluster, userName string, err error) {
	cluster, userName, err = i.verifyDetailed(token)
	if err != nil {
		return "", "", errInvalidToken
	}
	return cluster, userName, nil
}

// verifyDetailed performs the actual token verification with detailed
// error messages for logging. The public Verify method wraps failures
// into a generic error before returning to the caller.
func (i *ManifestTokenIssuer) verifyDetailed(token string) (cluster, userName string, err error) {
	parts := strings.SplitN(token, ".", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("malformed token")
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return "", "", fmt.Errorf("decode payload: %w", err)
	}

	sig, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", "", fmt.Errorf("decode signature: %w", err)
	}

	// Verify HMAC before trusting any payload content.
	mac := hmac.New(sha256.New, i.hmacKey)
	mac.Write(payloadBytes)
	if !hmac.Equal(sig, mac.Sum(nil)) {
		return "", "", fmt.Errorf("invalid token signature")
	}

	var claims manifestTokenClaims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return "", "", fmt.Errorf("parse token claims: %w", err)
	}

	now := i.now().Unix()

	if now > claims.Exp {
		return "", "", fmt.Errorf("token expired")
	}

	// Sanity-check iat: reject tokens that claim to be issued in
	// the future (clock skew allowance: 5 minutes) or that are
	// older than the maximum token TTL plus a small buffer. This
	// limits the replay window for leaked tokens.
	const clockSkew = 5 * 60 // 5 minutes in seconds
	maxAge := int64(manifestTokenTTL.Seconds()) + clockSkew
	if claims.Iat > now+clockSkew {
		return "", "", fmt.Errorf("token issued in the future")
	}
	if now-claims.Iat > maxAge {
		return "", "", fmt.Errorf("token too old")
	}

	return claims.Cluster, claims.Sub, nil
}
