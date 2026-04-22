package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"connectrpc.com/authn"
	"github.com/coreos/go-oidc/v3/oidc"

	"github.com/otterscale/otterscale/internal/core"
)

// oidcGroupClaims holds the custom claims extracted from an OIDC ID
// token. "groups" is a standard OIDC claim; "resource_access" is the
// Keycloak-specific container for client-scoped roles, keyed by client
// id. Both sources are merged and prefixed with "oidc:" by the
// middleware so callers can treat roles and groups uniformly.
type oidcGroupClaims struct {
	Groups         []string                      `json:"groups"`
	ResourceAccess map[string]oidcResourceAccess `json:"resource_access"`
}

type oidcResourceAccess struct {
	Roles []string `json:"roles"`
}

// NewOIDC creates a ConnectRPC authentication middleware that verifies
// incoming Bearer tokens against the given OIDC issuer and client ID.
//
// On success, the authenticated user's subject and groups are stored
// in the request context as core.UserInfo. OIDC groups are prefixed
// with "oidc:" to keep them separate from Kubernetes-native groups and
// avoid unintended privilege escalation via name collisions. The
// "system:authenticated" group is always included.
func NewOIDC(issuer, clientID string) (*authn.Middleware, error) {
	const oidcDiscoveryTimeout = 10 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), oidcDiscoveryTimeout)
	defer cancel()

	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		return nil, fmt.Errorf("failed to init oidc provider: %w", err)
	}

	verifier := provider.Verifier(&oidc.Config{
		ClientID: clientID,
	})

	authenticate := func(ctx context.Context, r *http.Request) (any, error) {
		token, found := authn.BearerToken(r)
		if !found || token == "" {
			return nil, authn.Errorf("missing or invalid bearer token")
		}

		idToken, err := verifier.Verify(ctx, token)
		if err != nil {
			return nil, authn.Errorf("invalid token: %s", err)
		}

		var claims oidcGroupClaims
		if err := idToken.Claims(&claims); err != nil {
			return nil, authn.Errorf("parse token claims: %s", err)
		}

		// Prefix with "oidc:" to avoid collisions with Kubernetes
		// built-in groups (e.g. "system:masters"). Roles attached to
		// the verifier's client id in resource_access are merged with
		// top-level groups so a Keycloak client role named "admin"
		// and a realm group named "admin" both land on "oidc:admin".
		clientRoles := claims.ResourceAccess[clientID].Roles
		groups := make([]string, 0, 1+len(claims.Groups)+len(clientRoles))
		groups = append(groups, "system:authenticated")
		for _, names := range [][]string{claims.Groups, clientRoles} {
			for _, n := range names {
				groups = append(groups, "oidc:"+n)
			}
		}

		return core.UserInfo{
			Subject: idToken.Subject,
			Groups:  groups,
		}, nil
	}

	return authn.NewMiddleware(authenticate), nil
}
