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
// returns the embedded cluster name and user identity.
func (h *ManifestHandler) VerifyManifestToken(ctx context.Context, token string) (cluster, userName string, err error) {
	return h.link.VerifyManifestToken(ctx, token)
}

// RenderManifest generates the agent installation manifest for the
// given cluster and user.
func (h *ManifestHandler) RenderManifest(ctx context.Context, cluster, userName string) (string, error) {
	return h.link.GenerateAgentManifest(ctx, cluster, userName)
}
