<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { PrometheusDriver } from 'prometheus-query';
	import { getContext } from 'svelte';
	import { type Writable } from 'svelte/store';
	import { toast } from 'svelte-sonner';

	import { env } from '$env/dynamic/public';
	import { EnvironmentService } from '$lib/api/environment/v1/environment_pb';
	import type { Extension } from '$lib/api/orchestrator/v1/orchestrator_pb';
	import { OrchestratorService } from '$lib/api/orchestrator/v1/orchestrator_pb';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import { ScrollArea } from '$lib/components/ui/scroll-area/index.js';
	import { Separator } from '$lib/components/ui/separator/index.js';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import { m } from '$lib/paraglide/messages';

	import type { ExtensionsBundleType } from './types';

	// API response types
	interface Node {
		id: string;
		name: string;
		status?: string;
	}

	interface Device {
		id: string;
		name: string;
		path: string;
		size?: string;
		type?: string;
	}

	interface NodeDevices {
		node: Node;
		devices: Device[];
	}

	interface NodeConfig {
		nodeId: string;
		nodeName: string;
		enabled: boolean;
		selectedDevices: string[];
	}

	const isDev = env.PUBLIC_DEV_MODE === 'true';
	// Mock data for development environment
	const MOCK_NODES_WITH_DEVICES: NodeDevices[] = [
		{
			node: { id: 'node-1', name: 'worker-node-01', status: 'Ready' },
			devices: [
				{ id: 'dev-1-1', name: 'Samsung 870 EVO', path: '/dev/sda', size: '500GB', type: 'SSD' },
				{ id: 'dev-1-2', name: 'WD Blue', path: '/dev/sdb', size: '1TB', type: 'HDD' },
				{ id: 'dev-1-3', name: 'Seagate BarraCuda', path: '/dev/sdc', size: '1TB', type: 'HDD' }
			]
		},
		{
			node: { id: 'node-2', name: 'worker-node-02', status: 'Ready' },
			devices: [
				{ id: 'dev-2-1', name: 'Crucial MX500', path: '/dev/sda', size: '500GB', type: 'SSD' },
				{ id: 'dev-2-2', name: 'Toshiba P300', path: '/dev/sdb', size: '2TB', type: 'HDD' }
			]
		},
		{
			node: { id: 'node-3', name: 'worker-node-03', status: 'Ready' },
			devices: [
				{ id: 'dev-3-1', name: 'Kingston A2000', path: '/dev/sda', size: '1TB', type: 'SSD' },
				{ id: 'dev-3-2', name: 'PASCARI AI100E', path: '/dev/nvme0n1', size: '2TB', type: 'NVMe' }
			]
		}
	];
</script>

<script lang="ts">
	let {
		scope,
		extensionsBundle,
		extensions,
		updator,
		open = $bindable(false),
		onConfirm
	}: {
		scope: string;
		extensionsBundle: ExtensionsBundleType;
		extensions: Writable<Extension[]>;
		updator: () => void;
		open: boolean;
		onConfirm?: (config: any) => void;
	} = $props();

	const transport: Transport = getContext('transport');
	const orchestratorClient = createClient(OrchestratorService, transport);
	const environmentService = createClient(EnvironmentService, transport);

	let isLoading = $state(false);
	let loadingData = $state(false);
	let loadError = $state<string | null>(null);
	let nodesWithDevices = $state<NodeDevices[]>([]);
	let nodeConfigs = $state<NodeConfig[]>([]);
	let prometheusDriver = $state<PrometheusDriver | null>(null);

	function getDeviceType(rotational: string): string {
		if (rotational === '0') {
			return 'SSD';
		}
		return 'HDD';
	}

	// Helper function to add timeout to promises
	function withTimeout<T>(promise: Promise<T>, timeoutMs: number, errorMsg: string): Promise<T> {
		return Promise.race([
			promise,
			new Promise<T>((_, reject) => setTimeout(() => reject(new Error(errorMsg)), timeoutMs))
		]);
	}

	// Load all nodes with their devices from Prometheus
	async function loadNodesAndDevices() {
		loadingData = true;
		loadError = null;

		try {
			if (!prometheusDriver) {
				const response = await withTimeout(
					environmentService.getPrometheus({}),
					10000,
					'Timeout: Failed to connect to Prometheus'
				);
				prometheusDriver = new PrometheusDriver({
					endpoint: `${env.PUBLIC_API_URL}/prometheus`,
					baseURL: response.baseUrl
				});
			}

			// Query node_disk_info metric filtered by scope with timeout
			const diskInfoResponse = await withTimeout(
				prometheusDriver.instantQuery(
					`node_disk_info{scope="${scope}",device=~"(^sd[a-z]$)|(^nvme[0-9]+n[0-9]+$)"}`
				),
				10000,
				'Timeout: Prometheus query took too long'
			);

			// Parse the response and group by node
			const nodeMap = new Map<string, NodeDevices>();

			if (diskInfoResponse.result && Array.isArray(diskInfoResponse.result)) {
				for (const item of diskInfoResponse.result) {
					const metric = item.metric;
					const hostname = metric.agent_hostname || metric.instance;
					const device = metric.device;
					const model = metric.model || 'Unknown';
					const rotational = metric.rotational || '1';

					// Determine device type
					let deviceType = getDeviceType(rotational);
					if (device.startsWith('nvme')) {
						deviceType = 'NVMe';
					}

					// Get or create node entry
					if (!nodeMap.has(hostname)) {
						nodeMap.set(hostname, {
							node: {
								id: hostname,
								name: hostname,
								status: 'Ready'
							},
							devices: []
						});
					}

					const nodeEntry = nodeMap.get(hostname)!;

					// Add device to node
					nodeEntry.devices.push({
						id: `${hostname}-${device}`,
						name: model,
						path: `/dev/${device}`,
						type: deviceType,
						size: undefined // Size info not available in node_disk_info
					});
				}
			}

			nodesWithDevices = Array.from(nodeMap.values());

			// Use mock data in development environment if no real data
			if (nodesWithDevices.length === 0 && isDev) {
				nodesWithDevices = MOCK_NODES_WITH_DEVICES;
			}

			// Sort nodes by name
			nodesWithDevices.sort((a, b) => a.node.name.localeCompare(b.node.name));

			// Sort devices within each node
			nodesWithDevices.forEach((nodeEntry) => {
				nodeEntry.devices.sort((a, b) => a.path.localeCompare(b.path));
			});

			// Initialize node configs
			nodeConfigs = nodesWithDevices.map((nd) => ({
				nodeId: nd.node.id,
				nodeName: nd.node.name,
				enabled: true,
				selectedDevices: []
			}));
		} catch (error) {
			const errorMessage =
				(error as ConnectError).message?.toString() || (error as Error).message || 'Unknown error';
			loadError = errorMessage;

			// Use mock data in development environment on error
			if (isDev) {
				nodesWithDevices = MOCK_NODES_WITH_DEVICES;

				nodeConfigs = nodesWithDevices.map((nd) => ({
					nodeId: nd.node.id,
					nodeName: nd.node.name,
					enabled: true,
					selectedDevices: []
				}));

				// Clear error in dev mode
				loadError = null;

				toast.error(m.rook_storage_failed_load_error(), {
					description: m.rook_storage_using_mock_data({ error: errorMessage })
				});
			} else {
				toast.error(m.rook_storage_failed_load_error(), {
					description: errorMessage
				});
			}
		} finally {
			await Promise.resolve();
			loadingData = false;
		}
	}

	// Load data when dialog opens
	$effect(() => {
		if (open && nodesWithDevices.length === 0) {
			loadNodesAndDevices();
		}
	});

	// Toggle node enabled status
	function toggleNode(nodeId: string) {
		const config = nodeConfigs.find((c) => c.nodeId === nodeId);
		if (config) {
			config.enabled = !config.enabled;
			if (!config.enabled) {
				config.selectedDevices = [];
			}
		}
	}

	// Toggle device selection
	function toggleDevice(nodeId: string, deviceId: string) {
		const config = nodeConfigs.find((c) => c.nodeId === nodeId);
		if (config && config.enabled) {
			const index = config.selectedDevices.indexOf(deviceId);
			if (index > -1) {
				config.selectedDevices = config.selectedDevices.filter((id) => id !== deviceId);
			} else {
				config.selectedDevices = [...config.selectedDevices, deviceId];
			}
		}
	}

	// Select/deselect all devices for a node
	function toggleAllDevices(nodeId: string) {
		const config = nodeConfigs.find((c) => c.nodeId === nodeId);
		const nodeData = nodesWithDevices.find((nd) => nd.node.id === nodeId);

		if (config && nodeData && config.enabled) {
			const allDeviceIds = nodeData.devices.map((d) => d.id);
			if (config.selectedDevices.length === allDeviceIds.length) {
				config.selectedDevices = [];
			} else {
				config.selectedDevices = allDeviceIds;
			}
		}
	}

	// Validate configuration
	function validateConfig(): boolean {
		const enabledNodes = nodeConfigs.filter((c) => c.enabled);
		if (enabledNodes.length === 0) {
			toast.error(m.rook_storage_enable_one_node());
			return false;
		}

		const nodesWithDevices = enabledNodes.filter((c) => c.selectedDevices.length > 0);
		if (nodesWithDevices.length === 0) {
			toast.error(m.rook_storage_select_one_device());
			return false;
		}

		return true;
	}

	function handleConfirm() {
		if (!validateConfig()) {
			return;
		}

		// Convert configuration to Helm ValuesMap format
		// Format: key-value pairs that will be injected into the chart
		const valuesMap: Record<string, string> = {
			'cephClusterSpec.storage.useAllNodes': 'false'
		};

		// Build the nodes configuration
		const enabledNodeConfigs = nodeConfigs.filter((c) => c.enabled && c.selectedDevices.length > 0);

		enabledNodeConfigs.forEach((c, nodeIndex) => {
			const nodeData = nodesWithDevices.find((nd) => nd.node.id === c.nodeId);
			const selectedDeviceNames =
				nodeData?.devices.filter((d) => c.selectedDevices.includes(d.id)).map((d) => d.path) || [];

			// Set node name: cephClusterSpec.storage.nodes[0].name=worker-node-01
			valuesMap[`cephClusterSpec.storage.nodes[${nodeIndex}].name`] = c.nodeName;

			// Set devices: cephClusterSpec.storage.nodes[0].devices[0].name=/dev/sda
			selectedDeviceNames.forEach((devicePath, deviceIndex) => {
				valuesMap[`cephClusterSpec.storage.nodes[${nodeIndex}].devices[${deviceIndex}].name`] =
					devicePath;
			});
		});

		if (isDev) {
			console.log('Rook Storage Helm ValuesMap:', valuesMap);
		}

		// Close dialog and return configuration
		open = false;

		// Call the onConfirm callback with the ValuesMap
		if (onConfirm) {
			onConfirm(valuesMap);
		}
	}

	function handleCancel() {
		open = false;
	}

	// Calculate selected statistics
	const selectedStats = $derived(() => {
		const enabledCount = nodeConfigs.filter((c) => c.enabled).length;
		const totalDevices = nodeConfigs
			.filter((c) => c.enabled)
			.reduce((sum, c) => sum + c.selectedDevices.length, 0);
		return { nodes: enabledCount, devices: totalDevices };
	});
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="max-h-[85vh] sm:max-w-[700px]">
		<Dialog.Header>
			<Dialog.Title>{m.rook_storage_config_title()}</Dialog.Title>
			<Dialog.Description>
				{m.rook_storage_config_description()}
			</Dialog.Description>
		</Dialog.Header>

		{#if loadingData}
			<div class="flex items-center justify-center py-8">
				<Icon icon="ph:spinner" class="size-8 animate-spin text-muted-foreground" />
				<span class="ml-2 text-muted-foreground">{m.rook_storage_loading_nodes()}</span>
			</div>
		{:else if loadError && nodesWithDevices.length === 0}
			<!-- Error state -->
			<div class="flex flex-col items-center justify-center px-4 py-12">
				<Icon icon="ph:warning-circle" class="mb-4 size-16 text-destructive" />
				<h3 class="mb-2 text-lg font-semibold">{m.rook_storage_failed_load_nodes()}</h3>
				<p class="mb-4 max-w-md text-center text-sm text-muted-foreground">
					{m.rook_storage_failed_load_description()}
				</p>
				<p
					class="mb-4 max-w-md rounded bg-muted p-2 text-center font-mono text-xs text-muted-foreground"
				>
					{loadError}
				</p>
				<Button variant="outline" onclick={loadNodesAndDevices}>
					<Icon icon="ph:arrow-clockwise" class="mr-2 size-4" />
					{m.rook_storage_retry()}
				</Button>
			</div>
		{:else if nodesWithDevices.length === 0}
			<!-- Empty state -->
			<div class="flex flex-col items-center justify-center px-4 py-12">
				<Icon icon="ph:database" class="mb-4 size-16 text-muted-foreground" />
				<h3 class="mb-2 text-lg font-semibold">{m.rook_storage_no_nodes_title()}</h3>
				<p class="max-w-md text-center text-sm text-muted-foreground">
					{m.rook_storage_no_nodes_description()}
				</p>
			</div>
		{:else}
			<ScrollArea class="max-h-[50vh] pr-4">
				<div class="space-y-4 py-4">
					{#each nodesWithDevices as { node, devices }, index}
						{@const config = nodeConfigs.find((c) => c.nodeId === node.id)}
						{#if config}
							<div class="rounded-lg border p-4 {config.enabled ? 'bg-background' : 'bg-muted/50'}">
								<!-- Node header -->
								<div class="mb-3 flex items-center justify-between">
									<div class="flex items-center gap-3">
										<Checkbox
											checked={config.enabled}
											onCheckedChange={() => toggleNode(node.id)}
										/>
										<div>
											<div class="flex items-center gap-2">
												<h4 class="font-semibold {!config.enabled ? 'text-muted-foreground' : ''}">
													{node.name}
												</h4>
												{#if node.status}
													<Badge variant="outline" class="text-xs">
														{node.status}
													</Badge>
												{/if}
											</div>
											{#if config.enabled}
												<p class="mt-1 text-xs text-muted-foreground">
													Selected {config.selectedDevices.length} / {devices.length} devices
												</p>
											{:else}
												<p class="mt-1 text-xs text-muted-foreground">Node disabled</p>
											{/if}
										</div>
									</div>
									{#if config.enabled && devices.length > 0}
										<Button variant="ghost" size="sm" onclick={() => toggleAllDevices(node.id)}>
											{config.selectedDevices.length === devices.length
												? m.rook_storage_deselect_all()
												: m.rook_storage_select_all()}
										</Button>
									{/if}
								</div>

								<!-- Device list -->
								{#if config.enabled}
									<Separator class="mb-3" />
									{#if devices.length > 0}
										<div class="grid grid-cols-3 gap-3 px-2">
											{#each devices as device}
												<label
													class="flex cursor-pointer flex-col items-center gap-2 rounded-lg border-2 p-3 transition-all {config.selectedDevices.includes(
														device.id
													)
														? 'border-primary bg-primary/5'
														: 'border-border hover:border-primary/50 hover:bg-muted/50'}"
												>
													<Checkbox
														checked={config.selectedDevices.includes(device.id)}
														onCheckedChange={() => toggleDevice(node.id, device.id)}
														class="self-end"
													/>
													<Icon
														icon="ph:hard-drive"
														class="size-12 {config.selectedDevices.includes(device.id)
															? 'text-primary'
															: 'text-muted-foreground'}"
													/>
													<div class="w-full text-center">
														<div class="truncate text-sm font-semibold">{device.name}</div>
														<div class="truncate text-xs text-muted-foreground">{device.path}</div>
														<div class="mt-1.5 flex items-center justify-center gap-1.5">
															{#if device.size}
																<Badge variant="secondary" class="text-xs">
																	{device.size}
																</Badge>
															{/if}
															{#if device.type}
																<Badge variant="outline" class="text-xs">
																	{device.type}
																</Badge>
															{/if}
														</div>
													</div>
												</label>
											{/each}
										</div>
									{:else}
										<p class="pl-8 text-sm text-muted-foreground">
											No devices available on this node
										</p>
									{/if}
								{/if}
							</div>
						{/if}
					{/each}
				</div>
			</ScrollArea>

			<!-- Summary -->
			<div class="mt-4 rounded-lg bg-muted/50 p-3">
				<div class="flex items-center justify-between text-sm">
					<span class="text-muted-foreground">{m.rook_storage_selected_summary()}</span>
					<div class="flex items-center gap-4">
						<span>
							<strong>{selectedStats().nodes}</strong>
							{m.rook_storage_nodes_count()}
						</span>
						<span>
							<strong>{selectedStats().devices}</strong>
							{m.rook_storage_devices_count()}
						</span>
					</div>
				</div>
			</div>
		{/if}

		<Dialog.Footer>
			<Button variant="outline" onclick={handleCancel} disabled={isLoading}>
				{m.cancel()}
			</Button>
			<Button
				onclick={handleConfirm}
				disabled={isLoading || loadingData || selectedStats().devices === 0}
			>
				{m.confirm()}
			</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
