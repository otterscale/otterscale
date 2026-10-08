package config

import "testing"

// The agent's defaults name objects the prometheus-stack and otterscale-agent
// charts create; a mismatch is a silent 403 or x509 error at query time.
func TestAgentPrometheusDefaults(t *testing.T) {
	t.Setenv("OTTERSCALE_AGENT_PROXY_PROMETHEUS_URL", "")
	t.Setenv("OTTERSCALE_AGENT_PROXY_PROMETHEUS_CA_SECRET", "")

	conf, err := New()
	if err != nil {
		t.Fatal(err)
	}

	if got, want := conf.AgentProxyPrometheusURL(), "https://prometheus-stack-kube-prom-prometheus.monitoring.svc:9091"; got != want {
		t.Errorf("AgentProxyPrometheusURL() = %q, want %q", got, want)
	}
	if got, want := conf.AgentProxyPrometheusCASecret(), "monitoring/prometheus-stack-kube-prom-prometheus-tls"; got != want {
		t.Errorf("AgentProxyPrometheusCASecret() = %q, want %q", got, want)
	}
}
