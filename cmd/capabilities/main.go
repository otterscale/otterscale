package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

type Capability struct {
	Category    string   `json:"category"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Features    []string `json:"features"`
	Available   bool     `json:"available"`
}

type CapabilitiesResponse struct {
	PlatformName        string       `json:"platform_name"`
	PlatformDescription string       `json:"platform_description"`
	Capabilities        []Capability `json:"capabilities"`
	UseCases            []string     `json:"use_cases"`
	DocumentationURL    string       `json:"documentation_url"`
	Version             string       `json:"version"`
}

func getCapabilities(language string) *CapabilitiesResponse {
	isZh := strings.HasPrefix(strings.ToLower(language), "zh")
	
	capabilities := []Capability{
		{
			Category:    ternary(isZh, "虛擬化管理", "Virtualization Management"),
			Name:        ternary(isZh, "虛擬機生命週期管理", "VM Lifecycle Management"),
			Description: ternary(isZh, "創建、啟動、停止、暫停、遷移虛擬機", "Create, start, stop, pause, migrate virtual machines"),
			Features:    ternary(isZh, []string{"KVM/QEMU集成", "GPU直通", "熱遷移", "快照管理"}, []string{"KVM/QEMU Integration", "GPU Passthrough", "Live Migration", "Snapshot Management"}),
			Available:   true,
		},
		{
			Category:    ternary(isZh, "容器編排", "Container Orchestration"),
			Name:        ternary(isZh, "Kubernetes原生支援", "Kubernetes Native Support"),
			Description: ternary(isZh, "部署和管理容器化應用程序", "Deploy and manage containerized applications"),
			Features:    ternary(isZh, []string{"Juju Charm部署", "工作負載管理", "服務網格", "自動擴展"}, []string{"Juju Charm Deployment", "Workload Management", "Service Mesh", "Auto Scaling"}),
			Available:   true,
		},
		{
			Category:    ternary(isZh, "存儲服務", "Storage Services"),
			Name:        ternary(isZh, "分佈式存儲", "Distributed Storage"),
			Description: ternary(isZh, "基於Ceph的可擴展塊、對象和文件存儲", "Ceph-based scalable block, object, and file storage"),
			Features:    ternary(isZh, []string{"S3兼容對象存儲", "高性能塊存儲", "POSIX文件系統", "備份與恢復"}, []string{"S3-Compatible Object Storage", "High-Performance Block Storage", "POSIX File Systems", "Backup & Recovery"}),
			Available:   true,
		},
		{
			Category:    ternary(isZh, "網絡", "Networking"),
			Name:        ternary(isZh, "軟件定義網絡", "Software-Defined Networking"),
			Description: ternary(isZh, "虛擬網絡、子網和路由", "Virtual networks, subnets, and routing"),
			Features:    ternary(isZh, []string{"負載均衡", "防火牆管理", "VPN集成", "網絡隔離"}, []string{"Load Balancing", "Firewall Management", "VPN Integration", "Network Isolation"}),
			Available:   true,
		},
		{
			Category:    ternary(isZh, "基礎設施管理", "Infrastructure Management"),
			Name:        ternary(isZh, "裸機配置", "Bare Metal Provisioning"),
			Description: ternary(isZh, "MAAS集成進行物理服務器管理", "MAAS integration for physical server management"),
			Features:    ternary(isZh, []string{"資源分配", "高可用性", "自動故障轉移", "水平擴展"}, []string{"Resource Allocation", "High Availability", "Automatic Failover", "Horizontal Scaling"}),
			Available:   true,
		},
		{
			Category:    ternary(isZh, "監控與可觀測性", "Monitoring & Observability"),
			Name:        ternary(isZh, "全面監控", "Comprehensive Monitoring"),
			Description: ternary(isZh, "基於Prometheus的監控和Grafana可視化", "Prometheus-based monitoring with Grafana visualization"),
			Features:    ternary(isZh, []string{"指標收集", "告警系統", "日誌聚合", "分佈式追蹤"}, []string{"Metrics Collection", "Alerting System", "Log Aggregation", "Distributed Tracing"}),
			Available:   true,
		},
		{
			Category:    ternary(isZh, "安全與訪問控制", "Security & Access Control"),
			Name:        ternary(isZh, "企業級安全", "Enterprise Security"),
			Description: ternary(isZh, "基於角色的訪問控制和企業認證", "Role-based access control and enterprise authentication"),
			Features:    ternary(isZh, []string{"RBAC", "LDAP/AD集成", "單點登錄", "數據加密", "審計日誌"}, []string{"RBAC", "LDAP/AD Integration", "Single Sign-On", "Data Encryption", "Audit Logging"}),
			Available:   true,
		},
		{
			Category:    ternary(isZh, "應用市場", "Application Marketplace"),
			Name:        ternary(isZh, "精選應用程序", "Curated Applications"),
			Description: ternary(isZh, "預配置的應用程序，可立即部署", "Pre-configured applications ready for deployment"),
			Features:    ternary(isZh, []string{"Charm商店", "自定義應用程序", "應用程序生命週期", "一鍵部署"}, []string{"Charm Store", "Custom Applications", "Application Lifecycle", "One-Click Deployment"}),
			Available:   true,
		},
		{
			Category:    ternary(isZh, "API與集成", "API & Integration"),
			Name:        ternary(isZh, "全面API支援", "Comprehensive API Support"),
			Description: ternary(isZh, "REST API和gRPC服務覆蓋所有平台操作", "REST API and gRPC services for all platform operations"),
			Features:    ternary(isZh, []string{"RESTful API", "gRPC服務", "CLI工具", "Webhook支持", "Terraform提供商"}, []string{"RESTful API", "gRPC Services", "CLI Tools", "Webhook Support", "Terraform Provider"}),
			Available:   true,
		},
	}

	useCases := ternary(isZh,
		[]string{
			"企業數據中心 - 多租戶基礎設施和資源優化",
			"開發與測試 - CI/CD集成和環境配置",
			"邊緣計算 - 分佈式部署和本地處理",
			"雲遷移 - 混合雲和工作負載遷移",
		},
		[]string{
			"Enterprise Data Centers - Multi-tenant infrastructure and resource optimization",
			"Development & Testing - CI/CD integration and environment provisioning",
			"Edge Computing - Distributed deployment and local processing",
			"Cloud Migration - Hybrid cloud and workload migration",
		},
	)

	return &CapabilitiesResponse{
		PlatformName:        "OtterScale",
		PlatformDescription: ternary(isZh, "統一基礎設施，賦能創新 - 超融合基礎設施平台", "Unifying Infrastructure, Empowering Innovation - Hyper-Converged Infrastructure Platform"),
		Capabilities:        capabilities,
		UseCases:            useCases,
		DocumentationURL:    "https://otterscale.github.io",
		Version:             "v0.6.0",
	}
}

func printCapabilities(resp *CapabilitiesResponse) error {
	fmt.Printf("🦦 %s\n", resp.PlatformName)
	fmt.Printf("%s\n\n", resp.PlatformDescription)

	isZh := strings.Contains(resp.PlatformDescription, "統一")
	fmt.Printf("📋 %s:\n", ternary(isZh, "核心功能", "Core Capabilities"))
	for _, cap := range resp.Capabilities {
		fmt.Printf("\n🔹 %s: %s\n", cap.Category, cap.Name)
		fmt.Printf("   %s\n", cap.Description)
		if len(cap.Features) > 0 {
			fmt.Printf("   %s: ", ternary(isZh, "功能特性", "Features"))
			fmt.Printf("%s\n", strings.Join(cap.Features, ", "))
		}
		status := ternary(isZh, "✅ 可用", "✅ Available")
		if !cap.Available {
			status = ternary(isZh, "⏳ 規劃中", "⏳ Planned")
		}
		fmt.Printf("   %s: %s\n", ternary(isZh, "狀態", "Status"), status)
	}

	fmt.Printf("\n🎯 %s:\n", ternary(isZh, "使用場景", "Use Cases"))
	for _, useCase := range resp.UseCases {
		fmt.Printf("• %s\n", useCase)
	}

	fmt.Printf("\n📚 %s: %s\n", ternary(isZh, "文檔", "Documentation"), resp.DocumentationURL)
	fmt.Printf("🏷️  %s: %s\n", ternary(isZh, "版本", "Version"), resp.Version)

	return nil
}

func newCapabilitiesCmd() *cobra.Command {
	var (
		language string
		format   string
	)

	cmd := &cobra.Command{
		Use:   "capabilities",
		Short: "Show OtterScale platform capabilities / 顯示OtterScale平台功能",
		Long: "Display comprehensive information about what OtterScale can do, including:\n" +
			"- Virtualization management\n" +
			"- Container orchestration\n" +
			"- Storage services\n" +
			"- Networking capabilities\n" +
			"- Infrastructure management\n" +
			"- Monitoring and observability\n" +
			"- Security and access control\n" +
			"- Application marketplace\n" +
			"- API and integration options\n\n" +
			"顯示OtterScale能夠執行的全面功能信息，包括：\n" +
			"- 虛擬化管理\n" +
			"- 容器編排\n" +
			"- 存儲服務\n" +
			"- 網絡功能\n" +
			"- 基礎設施管理\n" +
			"- 監控與可觀測性\n" +
			"- 安全與訪問控制\n" +
			"- 應用市場\n" +
			"- API與集成選項",
		RunE: func(cmd *cobra.Command, args []string) error {
			capabilities := getCapabilities(language)

			switch strings.ToLower(format) {
			case "json":
				encoder := json.NewEncoder(os.Stdout)
				encoder.SetIndent("", "  ")
				return encoder.Encode(capabilities)
			default:
				return printCapabilities(capabilities)
			}
		},
	}

	cmd.Flags().StringVarP(&language, "language", "l", "en", "Language for output (en, zh, zh-CN, zh-TW)")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format (text, json)")

	return cmd
}

func main() {
	ctx := context.Background()
	cmd := newCapabilitiesCmd()
	
	if err := cmd.ExecuteContext(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// Helper function for ternary operations
func ternary[T any](condition bool, trueVal T, falseVal T) T {
	if condition {
		return trueVal
	}
	return falseVal
}