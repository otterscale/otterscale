package values

import (
	"strings"
	"testing"

	"github.com/Masterminds/semver/v3"
)

func TestProvideChartVersions(t *testing.T) {
	versions := ProvideChartVersions()

	for name, version := range map[string]string{
		"otterscale-agent-flux": versions.AgentFlux,
		"otterscale-agent":      versions.Agent,
		"flux":                  versions.Flux,
	} {
		if version == "" {
			t.Errorf("chart %s has no version", name)
			continue
		}
		if strings.TrimSpace(version) != version {
			t.Errorf("chart %s version %q is surrounded by whitespace", name, version)
		}
		if _, err := semver.NewConstraint(version); err != nil {
			t.Errorf("chart %s version %q is not a valid constraint: %v", name, version, err)
		}
	}
}
