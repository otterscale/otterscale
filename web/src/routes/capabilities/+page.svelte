<script lang="ts">
	import { onMount } from 'svelte';
	import { browser } from '$app/environment';

	interface Capability {
		category: string;
		name: string;
		description: string;
		features: string[];
		available: boolean;
	}

	interface CapabilitiesResponse {
		platform_name: string;
		platform_description: string;
		capabilities: Capability[];
		use_cases: string[];
		documentation_url: string;
		version: string;
	}

	let capabilities: CapabilitiesResponse | null = null;
	let language = 'en';
	let loading = true;
	let error: string | null = null;

	// Detect browser language
	if (browser) {
		const browserLang = navigator.language.toLowerCase();
		if (browserLang.startsWith('zh')) {
			language = 'zh';
		}
	}

	async function loadCapabilities() {
		try {
			loading = true;
			error = null;
			
			// Mock capabilities data - in real implementation this would call the API
			capabilities = getCapabilities(language);
			
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to load capabilities';
		} finally {
			loading = false;
		}
	}

	function getCapabilities(lang: string): CapabilitiesResponse {
		const isZh = lang.startsWith('zh');
		
		return {
			platform_name: "OtterScale",
			platform_description: isZh 
				? "統一基礎設施，賦能創新 - 超融合基礎設施平台" 
				: "Unifying Infrastructure, Empowering Innovation - Hyper-Converged Infrastructure Platform",
			capabilities: [
				{
					category: isZh ? "虛擬化管理" : "Virtualization Management",
					name: isZh ? "虛擬機生命週期管理" : "VM Lifecycle Management",
					description: isZh ? "創建、啟動、停止、暫停、遷移虛擬機" : "Create, start, stop, pause, migrate virtual machines",
					features: isZh 
						? ["KVM/QEMU集成", "GPU直通", "熱遷移", "快照管理"]
						: ["KVM/QEMU Integration", "GPU Passthrough", "Live Migration", "Snapshot Management"],
					available: true
				},
				{
					category: isZh ? "容器編排" : "Container Orchestration",
					name: isZh ? "Kubernetes原生支援" : "Kubernetes Native Support",
					description: isZh ? "部署和管理容器化應用程序" : "Deploy and manage containerized applications",
					features: isZh 
						? ["Juju Charm部署", "工作負載管理", "服務網格", "自動擴展"]
						: ["Juju Charm Deployment", "Workload Management", "Service Mesh", "Auto Scaling"],
					available: true
				},
				{
					category: isZh ? "存儲服務" : "Storage Services",
					name: isZh ? "分佈式存儲" : "Distributed Storage",
					description: isZh ? "基於Ceph的可擴展塊、對象和文件存儲" : "Ceph-based scalable block, object, and file storage",
					features: isZh 
						? ["S3兼容對象存儲", "高性能塊存儲", "POSIX文件系統", "備份與恢復"]
						: ["S3-Compatible Object Storage", "High-Performance Block Storage", "POSIX File Systems", "Backup & Recovery"],
					available: true
				},
				{
					category: isZh ? "網絡" : "Networking",
					name: isZh ? "軟件定義網絡" : "Software-Defined Networking",
					description: isZh ? "虛擬網絡、子網和路由" : "Virtual networks, subnets, and routing",
					features: isZh 
						? ["負載均衡", "防火牆管理", "VPN集成", "網絡隔離"]
						: ["Load Balancing", "Firewall Management", "VPN Integration", "Network Isolation"],
					available: true
				},
				{
					category: isZh ? "基礎設施管理" : "Infrastructure Management",
					name: isZh ? "裸機配置" : "Bare Metal Provisioning",
					description: isZh ? "MAAS集成進行物理服務器管理" : "MAAS integration for physical server management",
					features: isZh 
						? ["資源分配", "高可用性", "自動故障轉移", "水平擴展"]
						: ["Resource Allocation", "High Availability", "Automatic Failover", "Horizontal Scaling"],
					available: true
				},
				{
					category: isZh ? "監控與可觀測性" : "Monitoring & Observability",
					name: isZh ? "全面監控" : "Comprehensive Monitoring",
					description: isZh ? "基於Prometheus的監控和Grafana可視化" : "Prometheus-based monitoring with Grafana visualization",
					features: isZh 
						? ["指標收集", "告警系統", "日誌聚合", "分佈式追蹤"]
						: ["Metrics Collection", "Alerting System", "Log Aggregation", "Distributed Tracing"],
					available: true
				},
				{
					category: isZh ? "安全與訪問控制" : "Security & Access Control",
					name: isZh ? "企業級安全" : "Enterprise Security",
					description: isZh ? "基於角色的訪問控制和企業認證" : "Role-based access control and enterprise authentication",
					features: isZh 
						? ["RBAC", "LDAP/AD集成", "單點登錄", "數據加密", "審計日誌"]
						: ["RBAC", "LDAP/AD Integration", "Single Sign-On", "Data Encryption", "Audit Logging"],
					available: true
				},
				{
					category: isZh ? "應用市場" : "Application Marketplace",
					name: isZh ? "精選應用程序" : "Curated Applications",
					description: isZh ? "預配置的應用程序，可立即部署" : "Pre-configured applications ready for deployment",
					features: isZh 
						? ["Charm商店", "自定義應用程序", "應用程序生命週期", "一鍵部署"]
						: ["Charm Store", "Custom Applications", "Application Lifecycle", "One-Click Deployment"],
					available: true
				},
				{
					category: isZh ? "API與集成" : "API & Integration",
					name: isZh ? "全面API支援" : "Comprehensive API Support",
					description: isZh ? "REST API和gRPC服務覆蓋所有平台操作" : "REST API and gRPC services for all platform operations",
					features: isZh 
						? ["RESTful API", "gRPC服務", "CLI工具", "Webhook支持", "Terraform提供商"]
						: ["RESTful API", "gRPC Services", "CLI Tools", "Webhook Support", "Terraform Provider"],
					available: true
				}
			],
			use_cases: isZh ? [
				"企業數據中心 - 多租戶基礎設施和資源優化",
				"開發與測試 - CI/CD集成和環境配置",
				"邊緣計算 - 分佈式部署和本地處理",
				"雲遷移 - 混合雲和工作負載遷移"
			] : [
				"Enterprise Data Centers - Multi-tenant infrastructure and resource optimization",
				"Development & Testing - CI/CD integration and environment provisioning",
				"Edge Computing - Distributed deployment and local processing",
				"Cloud Migration - Hybrid cloud and workload migration"
			],
			documentation_url: "https://otterscale.github.io",
			version: "v0.6.0"
		};
	}

	onMount(() => {
		loadCapabilities();
	});

	function switchLanguage(newLang: string) {
		language = newLang;
		loadCapabilities();
	}
</script>

<svelte:head>
	<title>{capabilities?.platform_name || 'OtterScale'} - {language.startsWith('zh') ? '功能能力' : 'Capabilities'}</title>
	<meta name="description" content={capabilities?.platform_description || 'OtterScale platform capabilities'} />
</svelte:head>

<div class="container mx-auto px-4 py-8 max-w-6xl">
	<!-- Language Toggle -->
	<div class="flex justify-end mb-6">
		<div class="flex space-x-2">
			<button 
				class="px-3 py-1 rounded {language === 'en' ? 'bg-blue-500 text-white' : 'bg-gray-200 text-gray-700'}"
				on:click={() => switchLanguage('en')}
			>
				English
			</button>
			<button 
				class="px-3 py-1 rounded {language === 'zh' ? 'bg-blue-500 text-white' : 'bg-gray-200 text-gray-700'}"
				on:click={() => switchLanguage('zh')}
			>
				中文
			</button>
		</div>
	</div>

	{#if loading}
		<div class="flex justify-center items-center h-64">
			<div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500"></div>
		</div>
	{:else if error}
		<div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
			<strong>{language.startsWith('zh') ? '錯誤' : 'Error'}:</strong> {error}
		</div>
	{:else if capabilities}
		<!-- Header -->
		<div class="text-center mb-12">
			<h1 class="text-4xl font-bold text-gray-900 mb-4">
				🦦 {capabilities.platform_name}
			</h1>
			<p class="text-xl text-gray-600 max-w-4xl mx-auto">
				{capabilities.platform_description}
			</p>
		</div>

		<!-- Core Capabilities -->
		<section class="mb-12">
			<h2 class="text-3xl font-bold text-gray-900 mb-8 flex items-center">
				📋 {language.startsWith('zh') ? '核心功能' : 'Core Capabilities'}
			</h2>
			<div class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
				{#each capabilities.capabilities as capability}
					<div class="bg-white rounded-lg shadow-md p-6 border border-gray-200 hover:shadow-lg transition-shadow">
						<div class="mb-4">
							<h3 class="text-lg font-semibold text-gray-900 mb-2">
								{capability.category}
							</h3>
							<h4 class="text-md font-medium text-blue-600 mb-3">
								{capability.name}
							</h4>
							<p class="text-gray-600 text-sm">
								{capability.description}
							</p>
						</div>
						
						{#if capability.features.length > 0}
							<div class="mb-4">
								<p class="text-sm font-medium text-gray-700 mb-2">
									{language.startsWith('zh') ? '功能特性' : 'Features'}:
								</p>
								<div class="flex flex-wrap gap-1">
									{#each capability.features as feature}
										<span class="bg-blue-100 text-blue-800 text-xs px-2 py-1 rounded">
											{feature}
										</span>
									{/each}
								</div>
							</div>
						{/if}

						<div class="flex items-center justify-between">
							<span class="text-sm text-gray-500">
								{language.startsWith('zh') ? '狀態' : 'Status'}:
							</span>
							<span class="flex items-center text-sm {capability.available ? 'text-green-600' : 'text-yellow-600'}">
								{capability.available ? '✅' : '⏳'}
								{capability.available 
									? (language.startsWith('zh') ? '可用' : 'Available')
									: (language.startsWith('zh') ? '規劃中' : 'Planned')
								}
							</span>
						</div>
					</div>
				{/each}
			</div>
		</section>

		<!-- Use Cases -->
		<section class="mb-12">
			<h2 class="text-3xl font-bold text-gray-900 mb-8 flex items-center">
				🎯 {language.startsWith('zh') ? '使用場景' : 'Use Cases'}
			</h2>
			<div class="grid gap-4 md:grid-cols-2">
				{#each capabilities.use_cases as useCase}
					<div class="bg-gradient-to-r from-blue-50 to-indigo-50 rounded-lg p-4 border border-blue-200">
						<p class="text-gray-800">• {useCase}</p>
					</div>
				{/each}
			</div>
		</section>

		<!-- Footer -->
		<section class="bg-gray-50 rounded-lg p-6 text-center">
			<div class="flex flex-col md:flex-row justify-center items-center space-y-4 md:space-y-0 md:space-x-8">
				<div class="flex items-center space-x-2">
					<span class="text-lg">📚</span>
					<span class="text-gray-600">{language.startsWith('zh') ? '文檔' : 'Documentation'}:</span>
					<a href={capabilities.documentation_url} class="text-blue-600 hover:text-blue-800 underline" target="_blank" rel="noopener noreferrer">
						{capabilities.documentation_url}
					</a>
				</div>
				<div class="flex items-center space-x-2">
					<span class="text-lg">🏷️</span>
					<span class="text-gray-600">{language.startsWith('zh') ? '版本' : 'Version'}:</span>
					<span class="font-mono text-blue-600">{capabilities.version}</span>
				</div>
			</div>
		</section>
	{/if}
</div>

<style>
	/* Additional responsive styles if needed */
	@media (max-width: 640px) {
		.container {
			padding-left: 1rem;
			padding-right: 1rem;
		}
	}
</style>