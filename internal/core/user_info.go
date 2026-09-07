package core

import (
	"context"
	"slices"
)

// UserInfo holds the authenticated user's identity and group memberships.
type UserInfo struct {
	Subject string
	Groups  []string
}

// adminGroup is the group membership required to perform privileged
// operations such as issuing a join token (which authorizes claiming a
// cluster, and thereby cluster-admin on it). The "oidc:" prefix matches
// the convention applied by the OIDC middleware to token group claims.
const adminGroup = "oidc:admin"

// IsAdmin reports whether the given group list contains the admin group.
func IsAdmin(groups []string) bool {
	return slices.Contains(groups, adminGroup)
}

// userInfoKey is unexported, so it cannot collide with other packages' keys.
type userInfoKey struct{}

// WithUserInfo is how the auth middleware publishes the caller's identity, so
// infrastructure adapters can read it without knowing the transport's context
// conventions.
func WithUserInfo(ctx context.Context, u UserInfo) context.Context {
	return context.WithValue(ctx, userInfoKey{}, u)
}

// UserInfoFromContext reports false when the context carries no UserInfo.
func UserInfoFromContext(ctx context.Context) (UserInfo, bool) {
	u, ok := ctx.Value(userInfoKey{}).(UserInfo)
	return u, ok
}
