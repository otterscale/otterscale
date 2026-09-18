package handler

import (
	"context"

	"github.com/otterscale/otterscale/internal/core"
)

// AgentValuesHandler serves the agent values behind a URL. Separate from
// LinkService because its caller is a plain HTTP route; keeping it here lets
// the serving layer mount that route without importing core.
type AgentValuesHandler struct {
	values *core.AgentValuesUseCase
}

func NewAgentValuesHandler(values *core.AgentValuesUseCase) *AgentValuesHandler {
	return &AgentValuesHandler{values: values}
}

// Render returns the values file the id stands for. An unknown id and an
// expired one come back as the same error, so neither can be probed for.
func (h *AgentValuesHandler) Render(ctx context.Context, id string) (string, error) {
	return h.values.RenderFromTicket(ctx, id)
}
