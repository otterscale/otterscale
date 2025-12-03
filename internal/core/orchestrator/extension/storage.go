package extension

import (
	"fmt"

	"github.com/otterscale/otterscale/internal/core/versions"
)

var sambaOperatorChartRef = fmt.Sprintf("https://github.com/otterscale/charts/releases/download/samba-operator-%[1]s/samba-operator-%[1]s.tgz", versions.SambaOperator)

var storageComponents = []component{
	{
		ID:          "samba-operator",
		DisplayName: "Samba",
		Description: "An operator for Samba as a service on PVCs in Kubernetes.",
		Logo:        "https://github.com/otterscale.png",
		Chart: &chartComponent{
			Name:      "samba-operator",
			Namespace: "samba-operator-system",
			Ref:       sambaOperatorChartRef,
			Version:   versions.SambaOperator,
		},
	},
}
