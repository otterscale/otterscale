package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"

	oscmd "github.com/otterscale/otterscale/internal/cmd"
	"github.com/otterscale/otterscale/internal/config"
)

var version = "devel"

func newCmd(conf *config.Config, mux *http.ServeMux) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "otterscale",
		Short:         "Open-source hyper-converged infrastructure platform / 開源超融合基礎設施平台",
		Long: "OtterScale is an open-source hyper-converged infrastructure platform that unifies compute, " +
			"storage, networking, and orchestration for scalable, enterprise-grade data center operations.\n\n" +
			"🌟 What can OtterScale do? / OtterScale可以做甚麼？\n" +
			"• 🖥️  Virtualization: KVM/QEMU VMs with live migration and GPU management\n" +
			"• 🐳 Container Orchestration: Native Kubernetes and Juju charm deployment\n" +
			"• 💾 Distributed Storage: Built-in Ceph clusters with automated backup\n" +
			"• 🌐 Software-Defined Networking: Virtual networks, load balancing, and security\n" +
			"• 📊 Monitoring: Integrated Prometheus and Grafana stack\n" +
			"• 🔐 Security: RBAC with LDAP/AD integration and SSO\n" +
			"• 🛒 Application Marketplace: Curated catalog of ready-to-deploy apps\n" +
			"• ⚡ High Availability: Multi-node deployment with automatic failover\n\n" +
			"For detailed capabilities, run: otterscale capabilities\n" +
			"詳細功能請運行：otterscale capabilities",
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	cmd.AddCommand(
		oscmd.NewInit(),
		oscmd.NewServe(conf, mux),
		oscmd.NewCapabilities(),
	)
	return cmd
}

func run() error {
	// options
	grpcHelper := true

	// wire cmd
	cmd, cleanup, err := wireCmd(grpcHelper)
	if err != nil {
		return err
	}
	defer cleanup()

	// start and wait for stop signal
	return cmd.Execute()
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
