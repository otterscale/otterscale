package core

import "context"

// UserInfo holds the authenticated user's identity and group memberships.
type UserInfo struct {
	Subject string
	Groups  []string
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
