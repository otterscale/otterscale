package kubernetes

import (
	"testing"

	"github.com/otterscale/otterscale/internal/core"
)

// TestImpersonationRequiresASubject guards the invariant the whole
// authorization model rests on. client-go attaches impersonation
// headers only when the config carries a user, uid, group or extra —
// so a UserInfo with no subject and no groups would produce a request
// with no identity at all, which the API server answers as the agent's
// own ServiceAccount rather than refusing.
func TestImpersonationRequiresASubject(t *testing.T) {
	tests := []struct {
		name    string
		user    core.UserInfo
		present bool // false marks "no UserInfo in the context at all"
	}{
		{name: "no user info", present: false},
		{name: "empty subject", user: core.UserInfo{}, present: true},
		{
			name:    "empty subject with groups",
			user:    core.UserInfo{Groups: []string{"system:authenticated"}},
			present: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			if tt.present {
				ctx = core.WithUserInfo(ctx, tt.user)
			}

			_, err := impersonation(ctx)
			if err == nil {
				t.Fatal("expected an error: a request with no subject would run as the agent's ServiceAccount")
			}
			code, ok := core.DomainErrorCode(err)
			if !ok || code != core.ErrorCodeUnauthenticated {
				t.Errorf("error = %v (code %v), want ErrorCodeUnauthenticated", err, code)
			}
		})
	}
}

// TestImpersonationCarriesSubjectAndGroups covers the accepted case.
func TestImpersonationCarriesSubjectAndGroups(t *testing.T) {
	ctx := core.WithUserInfo(t.Context(), core.UserInfo{
		Subject: "alice",
		Groups:  []string{"system:authenticated", "oidc:platform"},
	})

	got, err := impersonation(ctx)
	if err != nil {
		t.Fatalf("impersonation: %v", err)
	}
	if got.UserName != "alice" {
		t.Errorf("UserName = %q, want %q", got.UserName, "alice")
	}
	if len(got.Groups) != 2 {
		t.Errorf("Groups = %v, want both groups forwarded", got.Groups)
	}
}
