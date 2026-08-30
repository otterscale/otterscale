package http

import (
	"slices"
	"strings"
	"testing"
)

const testClientID = "otterscale"

func TestNewUserInfo_RejectsReservedSubjects(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		subject string
	}{
		{"service account", "system:serviceaccount:otterscale-system:flux-admin"},
		{"anonymous", "system:anonymous"},
		{"reserved prefix", "system:whatever"},
		{"empty", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if _, err := newUserInfo(tt.subject, oidcGroupClaims{}, testClientID); err == nil {
				t.Fatalf("newUserInfo(%q) = nil error, want rejection", tt.subject)
			}
		})
	}
}

func TestNewUserInfo_PrefixesEveryClaimedGroup(t *testing.T) {
	t.Parallel()

	claims := oidcGroupClaims{
		// Unprefixed, this claim would grant cluster-admin outright.
		Groups: []string{"system:masters", "platform"},
		ResourceAccess: map[string]oidcResourceAccess{
			testClientID: {Roles: []string{"admin"}},
			"other":      {Roles: []string{"must-not-appear"}},
		},
	}

	got, err := newUserInfo("alice@example.com", claims, testClientID)
	if err != nil {
		t.Fatalf("newUserInfo: %v", err)
	}

	if got.Subject != "alice@example.com" {
		t.Errorf("Subject = %q, want %q", got.Subject, "alice@example.com")
	}

	want := []string{"system:authenticated", "oidc:system:masters", "oidc:platform", "oidc:admin"}
	if !slices.Equal(got.Groups, want) {
		t.Fatalf("Groups = %v, want %v", got.Groups, want)
	}

	// system:authenticated is ours to add; no claim may land outside oidc:.
	for _, g := range got.Groups {
		if g == "system:authenticated" {
			continue
		}
		if !strings.HasPrefix(g, "oidc:") {
			t.Errorf("group %q escaped the oidc: namespace", g)
		}
	}
}
