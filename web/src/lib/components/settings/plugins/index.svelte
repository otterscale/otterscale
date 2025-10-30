<script lang="ts" module>
	import { timestampDate } from '@bufbuild/protobuf/wkt';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';

	import { OrchestratorService, type Plugin } from '$lib/api/orchestrator/v1/orchestrator_pb';
	import * as Accordion from '$lib/components/ui/accordion/index.js';
	import Button from '$lib/components/ui/button/button.svelte';
	import * as Card from '$lib/components/ui/card';
	import { Progress } from '$lib/components/ui/progress/index.js';
	import { cn } from '$lib/utils';

	export type PluginConfiguration = {
		name: string;
		description: string;
		icon: string;
	};
	export type PlatformConfigurations = Record<
		'General' | 'Model' | 'Instance' | 'Storage',
		{
			name: string;
			plugins: PluginConfiguration[];
			description: string;
			icon: string;
		}
	>;

	export const platformConfigurations: PlatformConfigurations = {
		General: {
			name: 'Dashboards',
			plugins: [
				{
					name: 'kube-prometheus-stack',
					description:
						'Prometheus is an open-source monitoring system that scrapes metrics, stores time-series data, and supports queries and alerts.',
					icon: 'ph:gauge',
				},
			],
			description:
				'Create interactive dashboards with customizable widgets, real-time charts, and drill-down insights from metrics and logs, plus role-based access controls for secure collaboration.',
			icon: 'ph:gauge',
		},
		Model: {
			name: 'Models',
			plugins: [
				{
					name: 'gpu-operator',
					description:
						'Installs and manages NVIDIA GPU drivers, device plugins, and related components to enable GPU-accelerated workloads on the cluster.',
					icon: 'ph:graphics-card',
				},
				{
					name: 'llm-d-infra',
					description:
						'llm-d-infra are the infrastructure components surrounding the llm-d system - a Kubernetes-native high-performance distributed LLM inference framework',
					icon: 'ph:robot',
				},
			],
			description: 'Enable vLLM plugins (tokenizers, backends, batching, logging).',
			icon: 'ph:robot',
		},

		Instance: {
			name: 'Virtual Machines',
			plugins: [
				{
					name: 'kubevirt-infra',
					description:
						'Provides KubeVirt components and tooling to run and manage virtual machines on Kubernetes — VM lifecycle, networking and storage integration, snapshots, and secure isolation.',
					icon: 'ph:desktop-tower',
				},
			],
			description:
				'Provision and manage virtual machines with scalable resource allocation, snapshots, and secure networking.',
			icon: 'ph:desktop-tower',
		},
		Storage: {
			name: 'Storages',
			plugins: [
				{
					name: 'samba-operator',
					description:
						'Installs and manages Samba/CIFS file share services to provide SMB-compatible file storage for workloads and expose persistent file shares to clients.',
					icon: 'ph:hard-drive',
				},
				{
					name: 'nfs-operator',
					description:
						'Deploys and manages NFS servers and exports, enabling scalable network file storage with POSIX semantics for stateful applications.',
					icon: 'ph:hard-drive',
				},
			],
			description:
				'Provide scalable, redundant storage for stateful workloads — block, file, and object stores with dynamic provisioning, snapshotting, and backup integrations.',
			icon: 'ph:hard-drives',
		},
	};
</script>

<script lang="ts">
	let { scope, facility }: { scope: string; facility: string } = $props();

	const transport: Transport = getContext('transport');
	const orchestratorService = createClient(OrchestratorService, transport);

	let chartNameToInstalledPlugins: Record<string, Plugin[] | undefined> = $state({});
	const installedCharts = $derived(Object.keys(chartNameToInstalledPlugins));

	orchestratorService
		.listPlugins({ scope: scope, facility: facility })
		.then((respoonse) => {
			chartNameToInstalledPlugins = Object.groupBy(respoonse.plugins, (plugin) => plugin?.chart?.name ?? '');
		})
		.catch((error) => {
			console.error('Failed to fetch plugins:', error);
		});
</script>

<Accordion.Root
	type="multiple"
	class="group bg-card text-card-foreground w-full overflow-hidden rounded-lg border transition-all duration-300"
>
	{#each Object.entries(platformConfigurations) as [_, pluginConfiguration], index}
		<Accordion.Item value={String(index)} class="p-6">
			<Accordion.Trigger>
				{@render Thumbnial(
					pluginConfiguration.icon,
					pluginConfiguration.name,
					pluginConfiguration.description,
					pluginConfiguration.plugins.filter((plugin) => installedCharts.includes(plugin.name)).length,
					pluginConfiguration.plugins.length,
				)}
			</Accordion.Trigger>
			<Accordion.Content class="flex flex-col gap-4 text-balance">
				{#each pluginConfiguration.plugins as requirement, index}
					<div class="relative mx-auto hidden w-full space-y-12 pt-10 md:block">
						{@render Node(
							index % 2 ? 'right' : 'left',
							requirement.icon,
							requirement.name,
							requirement.description,
							installedCharts.includes(requirement.name),
						)}
					</div>
				{/each}
			</Accordion.Content>
		</Accordion.Item>
	{/each}
</Accordion.Root>

{#snippet Thumbnial(icon: string, title: string, description: string, installed: number, required: number)}
	{@const percentage = (installed * 100) / required}
	<div class="flex w-full flex-col gap-4">
		<Progress
			value={percentage}
			class={cn(
				installed == required
					? 'bg-green-700/20 **:data-[slot="progress-indicator"]:bg-green-700'
					: installed > 0
						? 'bg-yellow-500/20 **:data-[slot="progress-indicator"]:bg-yellow-500'
						: 'bg-red-700/20 **:data-[slot="progress-indicator"]:bg-red-700',
			)}
		/>
		<div class="flex items-center gap-2">
			<div class="flex h-full items-center justify-between gap-4">
				<div class="bg-primary/10 rounded-lg p-2">
					<Icon {icon} class="size-8" />
				</div>
				<div>
					<h3 class="text-lg font-bold">{title}</h3>
					<p class="text-muted-foreground mt-1 text-sm">
						{description}
					</p>
				</div>
			</div>
			<div class="ml-auto flex flex-col justify-between gap-4">
				<p class="text-muted-foreground whitespace-nowrap">{installed} over {required}</p>
				{#if installed < required}
					<Button size="sm">Install</Button>
				{:else}
					<Button size="sm">Upgrade</Button>
				{/if}
			</div>
		</div>
	</div>
{/snippet}

{#snippet Node(
	alignment: 'left' | 'right' = 'right',
	icon: string,
	name: string,
	description: string,
	installed: boolean,
	action?: string,
)}
	<div
		class={alignment == 'right'
			? 'relative flex flex-row-reverse items-center gap-8'
			: 'relative flex items-center gap-8'}
	>
		<div
			class="bg-primary text-primary-foreground absolute left-1/2 z-10 flex h-12 w-12 -translate-x-1/2 transform items-center justify-center rounded-full font-bold"
		>
			<Icon {icon} class="size-6" />
		</div>
		<div class={alignment == 'right' ? 'w-1/2 pr-16' : 'w-1/2 pl-16'}>
			<Card.Root class={cn(installed ? 'bg-green-50' : 'bg-red-50', 'p-0')}>
				<Card.Content class="space-y-2 p-5">
					<div class="flex flex-row-reverse items-center justify-end gap-2">
						<Icon
							icon={installed ? 'ph:check-circle' : 'ph:minus-circle'}
							class={cn(installed ? 'text-green-500' : 'text-red-500', 'size-6')}
						/>
						<h3 class="text-base font-bold">{name}</h3>
					</div>
					<p class="text-muted-foreground text-sm font-light">
						{description}
					</p>
					{#if chartNameToInstalledPlugins[name]}
						{@const installedPlugins = chartNameToInstalledPlugins[name]}
						<div>
							{#each installedPlugins as installedPlugin}
								<div class="flex items-center gap-2 space-y-2">
									<span class="bg-secondary rounded-full border p-2">
										<Icon icon="ph:cube" class="text-secondary-foreground size-5" />
									</span>
									<div>
										<span class="flex items-center space-x-1 font-semibold">
											<p class="text-muted-foreground">{installedPlugin.namespace}</p>
											<p>{installedPlugin.name}</p>
										</span>
										{#if installedPlugin.lastDeployedAt}
											<p class="text-muted-foreground text-xs">
												{timestampDate(installedPlugin.lastDeployedAt).toLocaleString()}
											</p>
										{/if}
									</div>
								</div>
							{/each}
						</div>
					{/if}
					{#if installed}
						<div class="ml-auto">
							<Button class="w-full" size="sm" href={action}>Upgrade</Button>
						</div>
					{:else}
						<div class="ml-auto">
							<Button class="w-full" size="sm" href={action}>Install</Button>
						</div>
					{/if}
				</Card.Content>
			</Card.Root>
		</div>
		<div class="w-1/2"></div>
	</div>
{/snippet}
