package extension

import (
	"fmt"

	"github.com/otterscale/otterscale/internal/core/versions"
)

var (
	sambaOperatorChartRef    = fmt.Sprintf("https://github.com/otterscale/charts/releases/download/samba-operator-%[1]s/samba-operator-%[1]s.tgz", versions.SambaOperator)
	rookCephOperatorChartRef = fmt.Sprintf("https://charts.rook.io/release/rook-ceph-v%s.tgz", versions.RookCephOperator)
	rookCephClusterChartRef  = fmt.Sprintf("https://charts.rook.io/release/rook-ceph-cluster-v%s.tgz", versions.RookCephCluster)
)

var storageComponents = []component{
	{
		ID:          "rook-ceph",
		DisplayName: "Rook Ceph Operator",
		Description: "Open-Source, Cloud-Native Storage Orchestrator for Kubernetes.",
		Logo:        "https://rook.io/images/rook-logo.svg",
		Chart: &chartComponent{
			Name:      "rook-ceph",
			Namespace: "rook-ceph",
			Ref:       rookCephOperatorChartRef,
			Version:   versions.RookCephOperator,
		},
	},
	{
		ID:          "rook-ceph-cluster",
		DisplayName: "Rook Ceph Cluster",
		Description: "Manages a Ceph cluster deployed in Kubernetes.",
		Logo:        "https://rook.io/images/rook-logo.svg",
		Chart: &chartComponent{
			Name:      "rook-ceph-cluster",
			Namespace: "rook-ceph",
			Ref:       rookCephClusterChartRef,
			Version:   versions.RookCephCluster,
			ValuesMap: map[string]string{
				"cephClusterSpec.network.provider":      "host",
				"cephObjectStores.spec[0].gateway.port": "8080",
			},
		},
	},
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
