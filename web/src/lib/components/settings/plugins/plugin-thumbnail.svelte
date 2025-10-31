<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { writable, type Writable } from 'svelte/store';
	import { toast } from 'svelte-sonner';

	import type { PluginsBundleType } from './types';

	import type { Plugin } from '$lib/api/orchestrator/v1/orchestrator_pb';
	import { OrchestratorService } from '$lib/api/orchestrator/v1/orchestrator_pb';
	import Button from '$lib/components/ui/button/button.svelte';
	import { Progress } from '$lib/components/ui/progress/index.js';
	import { m } from '$lib/paraglide/messages';
	import { cn } from '$lib/utils';

	const pluginsBundleConfigurations: Record<
		PluginsBundleType,
		{
			name: string;
			description: string;
			icon: string;
		}
	> = {
		general: {
			name: 'Dashboards',
			description:
				'Create interactive dashboards with customizable widgets, real-time charts, and drill-down insights from metrics and logs, plus role-based access controls for secure collaboration.',
			icon: 'ph:gauge',
		},
		model: {
			name: 'Models',
			description: 'Enable vLLM plugins (tokenizers, backends, batching, logging).',
			icon: 'ph:robot',
		},
		instance: {
			name: 'Virtual Machines',
			description:
				'Provision and manage virtual machines with scalable resource allocation, snapshots, and secure networking.',
			icon: 'ph:desktop-tower',
		},
		storage: {
			name: 'Storages',
			description:
				'Provide scalable, redundant storage for stateful workloads — block, file, and object stores with dynamic provisioning, snapshotting, and backup integrations.',
			icon: 'ph:hard-drives',
		},
	};
</script>

<script lang="ts">
	let {
		scope,
		facility,
		pluginsBundle,
		plugins,
	}: { scope: string; facility: string; pluginsBundle: PluginsBundleType; plugins: Plugin[] } = $props();

	const transport: Transport = getContext('transport');
	const orchestratorClient = createClient(OrchestratorService, transport);

	const modelPlugins: Writable<Plugin[]> = writable([]);
	const generalPlugins: Writable<Plugin[]> = writable([]);

	orchestratorClient
		.listModelPlugins({ scope: scope, facility: facility })
		.then((respoonse) => {
			modelPlugins.set(respoonse.plugins);
		})
		.catch((error) => {
			console.error('Failed to fetch model plugins:', error);
		});
	orchestratorClient
		.listGeneralPlugins({ scope: scope, facility: facility })
		.then((respoonse) => {
			generalPlugins.set(respoonse.plugins);
		})
		.catch((error) => {
			console.error('Failed to fetch general plugins:', error);
		});

	const installed = $derived(plugins.filter((plugin) => plugin.current).length);
	const required = $derived(plugins.length);
</script>

<div class="flex w-full flex-col gap-4">
	<Progress
		value={(installed * 100) / required}
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
				<Icon icon={pluginsBundleConfigurations[pluginsBundle].icon} class="size-8" />
			</div>
			<div>
				<h3 class="text-lg font-bold">{pluginsBundleConfigurations[pluginsBundle].name}</h3>
				<p class="text-muted-foreground mt-1 text-sm">
					{pluginsBundleConfigurations[pluginsBundle].description}
				</p>
			</div>
		</div>
		<div class="ml-auto flex flex-col justify-between gap-4">
			<p class="text-muted-foreground whitespace-nowrap">{installed} over {required}</p>
			{#if installed < required}
				<div class="ml-auto">
					<Button
						class="w-full"
						size="sm"
						onclick={() => {
							toast.promise(
								() =>
									orchestratorClient.installPlugins({
										scope: scope,
										facility: facility,
										charts: $modelPlugins
											.filter((modelPlugin) => modelPlugin.latest && !modelPlugin.current)
											.map((modelPlugin) => modelPlugin.latest!),
									}),
								{
									loading: `Installing plugins of model`,
									success: () => {
										orchestratorClient
											.listModelPlugins({ scope: scope, facility: facility })
											.then((response) => {
												modelPlugins.set(response.plugins);
											});
										return `Successfully installed plugins`;
									},
									error: (error) => {
										let message = `Failed to install plugins`;
										toast.error(message, {
											description: (error as ConnectError).message.toString(),
											duration: Number.POSITIVE_INFINITY,
										});
										return message;
									},
								},
							);
						}}
					>
						{m.install()}
					</Button>
				</div>
			{:else}
				<div class="ml-auto">
					<Button
						class="w-full"
						size="sm"
						onclick={() => {
							toast.promise(
								() =>
									orchestratorClient.upgradePlugins({
										scope: scope,
										facility: facility,
										charts: $modelPlugins
											.filter((modelPlugin) => modelPlugin.latest && modelPlugin.current)
											.map((modelPlugin) => modelPlugin.latest!),
									}),
								{
									loading: `Upgrading plugins of model`,
									success: () => {
										orchestratorClient
											.listModelPlugins({ scope: scope, facility: facility })
											.then((response) => {
												modelPlugins.set(response.plugins);
											});
										return `Successfully upgraded plugins`;
									},
									error: (error) => {
										let message = `Failed to upgrade plugins`;
										toast.error(message, {
											description: (error as ConnectError).message.toString(),
											duration: Number.POSITIVE_INFINITY,
										});
										return message;
									},
								},
							);
						}}
					>
						{m.plugins_upgrade()}
					</Button>
				</div>
			{/if}
		</div>
	</div>
</div>
