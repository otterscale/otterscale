package http

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"connectrpc.com/authn"
	"github.com/coreos/go-oidc/v3/oidc"

	"github.com/otterscale/otterscale/internal/core"
)

// oidcGroupClaims covers both group sources: "groups" is the standard OIDC
// claim, "resource_access" the Keycloak container for client-scoped roles. The
// middleware merges them under an "oidc:" prefix, so callers treat roles and
// groups uniformly.
type oidcGroupClaims struct {
	Groups         []string                      `json:"groups"`
	ResourceAccess map[string]oidcResourceAccess `json:"resource_access"`
}

type oidcResourceAccess struct {
	Roles []string `json:"roles"`
}

// NewOIDC verifies incoming Bearer tokens against the issuer and client ID,
// storing the subject and groups in the request context as core.UserInfo.
// Groups are prefixed with "oidc:" so a name collision with a Kubernetes-native
// group cannot escalate privileges. "system:authenticated" is always included.
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

		return newUserInfo(idToken.Subject, claims, clientID)
	}

	return authn.NewMiddleware(authenticate), nil
}

// newUserInfo maps verified claims onto the identity sent as Impersonate-User
// and Impersonate-Group.
func newUserInfo(subject string, claims oidcGroupClaims, clientID string) (core.UserInfo, error) {
	// The subject doubles as the RBAC subject and the Harbor identity, so it
	// cannot carry a prefix the way groups do. Reject the reserved namespace
	// instead: "system:serviceaccount:<ns>:<name>" would impersonate a service
	// account rather than a user.
	if subject == "" || strings.HasPrefix(subject, "system:") {
		return core.UserInfo{}, authn.Errorf("invalid subject")
	}

	// The "oidc:" prefix keeps these clear of built-in groups such as
	// "system:masters". Client roles from resource_access merge with
	// top-level groups, so a Keycloak client role "admin" and a realm
	// group "admin" both land on "oidc:admin".
	clientRoles := claims.ResourceAccess[clientID].Roles
	groups := make([]string, 0, 1+len(claims.Groups)+len(clientRoles))
	groups = append(groups, "system:authenticated")
	for _, names := range [][]string{claims.Groups, clientRoles} {
		for _, n := range names {
			groups = append(groups, "oidc:"+n)
		}
	}

	return core.UserInfo{
		Subject: subject,
		Groups:  groups,
	}, nil
}
