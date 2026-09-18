package values

import "github.com/otterscale/otterscale/internal/core"

// The chart versions this build's rendered values are written for. Helm and
// Flux both accept a semver range here, so "1.1.x" is a legitimate value.
var chartVersions = core.ChartVersions{
	// renovate: helm=otterscale-agent-flux registry=https://otterscale.github.io/helm-charts
	AgentFlux: "1.1.x",
	// renovate: helm=otterscale-agent registry=https://otterscale.github.io/helm-charts
	Agent: "1.1.x",
	// renovate: helm=flux registry=https://otterscale.github.io/helm-charts
	Flux: "1.0.x",
}

func ProvideChartVersions() core.ChartVersions {
	return chartVersions
}
