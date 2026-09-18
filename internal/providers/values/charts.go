package values

import "github.com/otterscale/otterscale/internal/core"

// The chart versions this build's rendered values are written for. Helm and
// Flux both accept a semver range here, so "1.1.x" is a legitimate value.
var chartVersions = core.ChartVersions{
	AgentFlux: "1.1.x", // renovate: datasource=helm depName=otterscale-agent-flux registryUrl=https://otterscale.github.io/helm-charts
	Agent:     "1.1.x", // renovate: datasource=helm depName=otterscale-agent registryUrl=https://otterscale.github.io/helm-charts
	Flux:      "1.0.x", // renovate: datasource=helm depName=flux registryUrl=https://otterscale.github.io/helm-charts
}

func ProvideChartVersions() core.ChartVersions {
	return chartVersions
}
