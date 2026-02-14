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
// token. The "groups" claim is a standard OIDC claim supported by
// most providers (Keycloak, Dex, Auth0, etc.).
type oidcGroupClaims struct {
	Groups []string `json:"groups"`
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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

		groups := make([]string, 0, len(claims.Groups)+1)
		groups = append(groups, "system:authenticated")
		for _, g := range claims.Groups {
			// Prefix with "oidc:" to avoid collisions with
			// Kubernetes built-in groups (e.g. "system:masters").
			groups = append(groups, "oidc:"+g)
		}

		return core.UserInfo{
			Subject: idToken.Subject,
			Groups:  groups,
		}, nil
	}

	return authn.NewMiddleware(authenticate), nil
}
