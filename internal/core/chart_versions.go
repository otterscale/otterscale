package core

// ChartVersions pins the Helm charts one build installs. A struct rather than
// three strings so Wire can tell it apart, as Version does. Which versions
// those are is the providers layer's business.
type ChartVersions struct {
	// AgentFlux is the umbrella chart, reported to the caller rather than
	// rendered: the file cannot pin the chart that consumes it.
	AgentFlux string
	// Agent and Flux are the nested HelmReleases, pinned in the rendered values.
	Agent string
	Flux  string
}
