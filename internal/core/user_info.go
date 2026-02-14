package core

import "context"

// UserInfo holds the authenticated user's identity and group memberships.
type UserInfo struct {
	Subject string
	Groups  []string
}

// userInfoKey is the context key for UserInfo. Using an unexported
// struct type prevents collisions with other packages.
type userInfoKey struct{}

// WithUserInfo returns a derived context that carries the given UserInfo.
// This is used by the authentication middleware to store the authenticated
// user's identity so that infrastructure adapters can retrieve it without
// depending on transport-specific context conventions.
func WithUserInfo(ctx context.Context, u UserInfo) context.Context {
	return context.WithValue(ctx, userInfoKey{}, u)
}

// UserInfoFromContext extracts the UserInfo stored by WithUserInfo.
// Returns false if the context does not carry a UserInfo value.
func UserInfoFromContext(ctx context.Context) (UserInfo, bool) {
	u, ok := ctx.Value(userInfoKey{}).(UserInfo)
	return u, ok
}
