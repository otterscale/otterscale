<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { type Writable } from 'svelte/store';
	import { toast } from 'svelte-sonner';

	import type { ExtensionsBundleType } from './types';

	import type { Extension } from '$lib/api/orchestrator/v1/orchestrator_pb';
	import { OrchestratorService } from '$lib/api/orchestrator/v1/orchestrator_pb';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import Button from '$lib/components/ui/button/button.svelte';
	import { Progress } from '$lib/components/ui/progress/index.js';
	import { m } from '$lib/paraglide/messages';
	import { activeScope } from '$lib/stores';
	import { cn } from '$lib/utils';

	const extensionsBundleConfigurations: Record<
		ExtensionsBundleType,
		{
			name: string;
			description: string;
			icon: string;
		}
	> = {
		general: {
			name: m.general(),
			description: m.extensions_general_bundle_description(),
			icon: 'ph:cube',
		},
		model: {
			name: m.model(),
			description: m.extensions_model_bundle_description(),
			icon: 'ph:robot',
		},
		instance: {
			name: m.virtual_machine(),
			description: m.extensions_virtual_machine_bundle_description(),
			icon: 'ph:desktop-tower',
		},
		storage: {
			name: m.storage(),
			description: m.extensions_storage_bundle_description(),
			icon: 'ph:hard-drives',
		},
	};
</script>

<script lang="ts">
	let {
		scope,
		facility,
		extensionsBundle,
		extensions,
		updator,
	}: {
		scope: string;
		facility: string;
		extensionsBundle: ExtensionsBundleType;
		extensions: Writable<Extension[]>;
		updator: () => void;
	} = $props();

	const transport: Transport = getContext('transport');
	const orchestratorClient = createClient(OrchestratorService, transport);

	const installed = $derived($extensions.filter((extension) => extension.current).length);
	const required = $derived($extensions.length);
</script>

<div class="flex w-full flex-col gap-4">
	<Progress
		value={required ? (installed * 100) / required : 0}
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
				<Icon icon={extensionsBundleConfigurations[extensionsBundle].icon} class="size-8" />
			</div>
			<div>
				<div class="flex items-center gap-1">
					<h3 class="text-lg font-bold">{extensionsBundleConfigurations[extensionsBundle].name}</h3>
					<Badge>{$activeScope.name}</Badge>
				</div>

				<p class="text-muted-foreground mt-1 text-sm">
					{extensionsBundleConfigurations[extensionsBundle].description}
				</p>
			</div>
		</div>
		<div class="ml-auto flex flex-col items-center justify-between gap-4">
			<p class="text-muted-foreground whitespace-nowrap">{installed} / {required}</p>
			{#if installed < required}
				<div class="ml-auto">
					<Button
						class="w-full"
						size="sm"
						onclick={() => {
							toast.promise(
								() =>
									orchestratorClient.installExtensions({
										scope: scope,
										facility: facility,
										charts: $extensions
											.filter((extension) => extension.latest && !extension.current)
											.map((extension) => extension.latest!),
									}),
								{
									loading: `Installing ${extensionsBundle} extensions.`,
									success: () => {
										updator();
										return `Successfully installed ${extensionsBundle} extensions`;
									},
									error: (error) => {
										let message = `Failed to install ${extensionsBundle} extensions`;
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
									orchestratorClient.upgradeExtensions({
										scope: scope,
										facility: facility,
										charts: $extensions
											.filter((extension) => extension.latest && extension.current)
											.map((extension) => extension.latest!),
									}),
								{
									loading: `Upgrading ${extensionsBundle} extensions`,
									success: () => {
										updator();
										return `Successfully upgraded ${extensionsBundle} extensions`;
									},
									error: (error) => {
										let message = `Failed to upgrade ${extensionsBundle} extensions`;
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
						{m.extensions_upgrade()}
					</Button>
				</div>
			{/if}
		</div>
	</div>
</div>
