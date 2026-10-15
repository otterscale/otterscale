package config

import (
	"testing"

	"github.com/spf13/pflag"
)

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

// TestAgentPrometheusURLEmptyFlag: an empty flag disables the proxy, which an
// empty environment variable cannot (TestAgentPrometheusDefaults).
func TestAgentPrometheusURLEmptyFlag(t *testing.T) {
	t.Setenv("OTTERSCALE_AGENT_PROXY_PROMETHEUS_URL", "")

	conf, err := New()
	if err != nil {
		t.Fatal(err)
	}
	fs := pflag.NewFlagSet("agent", pflag.ContinueOnError)
	if err := conf.BindFlags(fs, AgentOptions); err != nil {
		t.Fatal(err)
	}
	if err := fs.Parse([]string{"--proxy-prometheus-url="}); err != nil {
		t.Fatal(err)
	}

	if got := conf.AgentProxyPrometheusURL(); got != "" {
		t.Errorf("AgentProxyPrometheusURL() = %q, want empty", got)
	}
}
