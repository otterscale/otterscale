package handler

import (
	"context"

	"github.com/otterscale/otterscale/internal/core"
)

// ManifestHandler provides token verification and manifest rendering
// for the raw HTTP manifest endpoint (kubectl apply -f <url>). It is
// separated from LinkService to keep the gRPC handler focused on
// ConnectRPC concerns and avoid coupling the transport layer to the
// handler layer for non-RPC operations.
type ManifestHandler struct {
	link *core.LinkUseCase
}

// NewManifestHandler returns a ManifestHandler backed by the given
// LinkUseCase.
func NewManifestHandler(link *core.LinkUseCase) *ManifestHandler {
	return &ManifestHandler{link: link}
}

// VerifyManifestToken validates an HMAC-signed manifest token and
// returns the embedded cluster name, user identity, and extra users
// bound to cluster-admin.
func (h *ManifestHandler) VerifyManifestToken(ctx context.Context, token string) (cluster, userName string, extraUsers []string, err error) {
	return h.link.VerifyManifestToken(ctx, token)
}

// RenderManifest generates the agent installation manifest for the
// given cluster, user, and additional cluster-admin users.
func (h *ManifestHandler) RenderManifest(ctx context.Context, cluster, userName string, extraUsers []string) (string, error) {
	return h.link.GenerateAgentManifest(ctx, cluster, userName, extraUsers)
}
