package handler

import (
	"context"
	"testing"

	"connectrpc.com/connect"

	pb "github.com/otterscale/otterscale/api/link/v1"

	"github.com/otterscale/otterscale/internal/core"
)

// stubTunnelProvider satisfies core.TunnelProvider without doing anything.
// IssueJoinToken never reaches the tunnel, so no call here should ever happen.
type stubTunnelProvider struct{}

func (stubTunnelProvider) CACertPEM() []byte { return nil }

func (stubTunnelProvider) ListLinks() map[string]core.Link { return map[string]core.Link{} }

func (stubTunnelProvider) ResolveAddress(context.Context, string) (string, error) {
	return "", nil
}

func (stubTunnelProvider) RegisterLink(
	context.Context, string, string, string, []byte,
) (core.TunnelGrant, error) {
	return core.TunnelGrant{}, nil
}

func newTestLinkService(t *testing.T) *LinkService {
	t.Helper()
	join, err := core.NewJoinAuthority("test-root-secret")
	if err != nil {
		t.Fatalf("NewJoinAuthority() error = %v", err)
	}
	return NewLinkService(core.NewLinkUseCase(stubTunnelProvider{}, core.Version("v1.0.0"), join))
}

func issueJoinTokenRequest(cluster string) *pb.IssueJoinTokenRequest {
	req := &pb.IssueJoinTokenRequest{}
	req.SetCluster(cluster)
	return req
}

// TestLinkService_IssueJoinToken_RequiresAdmin is the regression test for the
// gate on this procedure: the token it returns authorizes claiming a cluster,
// so an ordinary authenticated user must not be able to mint one.
func TestLinkService_IssueJoinToken_RequiresAdmin(t *testing.T) {
	tests := []struct {
		name string
		ctx  func(t *testing.T) context.Context
		want connect.Code
	}{
		{
			name: "no authenticated caller",
			ctx: func(t *testing.T) context.Context {
				t.Helper()
				return t.Context()
			},
			want: connect.CodeUnauthenticated,
		},
		{
			name: "authenticated but not an admin",
			ctx: func(t *testing.T) context.Context {
				t.Helper()
				return core.WithUserInfo(t.Context(), core.UserInfo{
					Subject: "user-1",
					Groups:  []string{"system:authenticated", "oidc:developer"},
				})
			},
			want: connect.CodePermissionDenied,
		},
		{
			// "admin" unprefixed is what a Kubernetes-native group looks like.
			// The OIDC middleware prefixes every claim it forwards, so matching an
			// unprefixed name would honour a group it never issued.
			name: "admin without the oidc: prefix",
			ctx: func(t *testing.T) context.Context {
				t.Helper()
				return core.WithUserInfo(t.Context(), core.UserInfo{
					Subject: "user-1",
					Groups:  []string{"system:authenticated", "admin"},
				})
			},
			want: connect.CodePermissionDenied,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newTestLinkService(t)

			resp, err := s.IssueJoinToken(tt.ctx(t), issueJoinTokenRequest("my-cluster"))
			if err == nil {
				t.Fatal("expected an error, got nil")
			}
			if got := connect.CodeOf(err); got != tt.want {
				t.Errorf("code = %v, want %v", got, tt.want)
			}
			if resp != nil {
				t.Error("a refused request must return no response")
			}
		})
	}
}

func TestLinkService_IssueJoinToken_AllowsAdmin(t *testing.T) {
	s := newTestLinkService(t)
	ctx := core.WithUserInfo(t.Context(), core.UserInfo{
		Subject: "admin-1",
		Groups:  []string{"system:authenticated", "oidc:admin"},
	})

	resp, err := s.IssueJoinToken(ctx, issueJoinTokenRequest("my-cluster"))
	if err != nil {
		t.Fatalf("IssueJoinToken() error = %v", err)
	}
	if resp.GetJoinToken() == "" {
		t.Error("join token is empty")
	}
}
