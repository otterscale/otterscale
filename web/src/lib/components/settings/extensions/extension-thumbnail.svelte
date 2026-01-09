<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { type Writable } from 'svelte/store';
	import { toast } from 'svelte-sonner';

	import type { Extension } from '$lib/api/orchestrator/v1/orchestrator_pb';
	import { OrchestratorService } from '$lib/api/orchestrator/v1/orchestrator_pb';
	import * as AlertDialog from '$lib/components/ui/alert-dialog/index.js';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import { buttonVariants } from '$lib/components/ui/button/index.js';
	import { Progress } from '$lib/components/ui/progress/index.js';
	import { m } from '$lib/paraglide/messages';
	import { cn } from '$lib/utils';

	import type { ExtensionsBundleType } from './types';
	import InstallRookConfig from './install-rook-config.svelte';

	const extensionsBundleConfigurations: Record<
		ExtensionsBundleType,
		{
			name: string;
			description: string;
			icon: string;
		}
	> = {
		metrics: {
			name: m.metrics(),
			description: m.extensions_metrics_bundle_description(),
			icon: 'ph:speedometer'
		},
		serviceMesh: {
			name: m.service_mesh(),
			description: m.extensions_service_mesh_bundle_description(),
			icon: 'ph:network'
		},
		model: {
			name: m.model(),
			description: m.extensions_model_bundle_description(),
			icon: 'ph:robot'
		},
		registry: {
			name: m.container_registry(),
			description: m.extensions_registry_bundle_description(),
			icon: 'ph:shipping-container'
		},
		instance: {
			name: m.virtual_machine(),
			description: m.extensions_virtual_machine_bundle_description(),
			icon: 'ph:desktop-tower'
		},
		storage: {
			name: m.storage(),
			description: m.extensions_storage_bundle_description(),
			icon: 'ph:hard-drives'
		}
	};
</script>

<script lang="ts">
	let {
		scope,
		extensionsBundle,
		extensions,
		updator
	}: {
		scope: string;
		extensionsBundle: ExtensionsBundleType;
		extensions: Writable<Extension[]>;
		updator: () => void;
	} = $props();

	const transport: Transport = getContext('transport');
	const orchestratorClient = createClient(OrchestratorService, transport);

	const installed = $derived($extensions.filter((extension) => extension.current).length);
	const required = $derived($extensions.length);
	let openConfirmDialog = $state(false);
	let openConfigDialog = $state(false);
	let isLoading = $state(false);
	let extensionArguments = $state<Record<string, Record<string, string>>>({});
	const isInstallMode = $derived(installed < required);

	function handleInstallClick() {
		if (isInstallMode && extensionsBundle === 'storage') {
			openConfigDialog = true;
		} else {
			openConfirmDialog = true;
		}
	}

	function onConfirm() {
		isLoading = true;
		toast.promise(
			() =>
				orchestratorClient.installOrUpgradeExtensions({
					scope: scope,
					manifests: $extensions
						.filter(
							(extension) =>
								extension.latest && (isInstallMode ? !extension.current : extension.current)
						)
						.map((extension) => extension.latest!),
					arguments: extensionArguments
				}),
			{
				loading: isInstallMode
					? `Installing ${extensionsBundle} extensions.`
					: `Upgrading ${extensionsBundle} extensions`,
				success: () => {
					isLoading = false;
					updator();
					return isInstallMode
						? `Successfully installed ${extensionsBundle} extensions`
						: `Successfully upgraded ${extensionsBundle} extensions`;
				},
				error: (error) => {
					isLoading = false;
					let message = isInstallMode
						? `Failed to install ${extensionsBundle} extensions`
						: `Failed to upgrade ${extensionsBundle} extensions`;
					toast.error(message, {
						description: (error as ConnectError).message.toString(),
						duration: Number.POSITIVE_INFINITY
					});
					return message;
				}
			}
		);
		openConfirmDialog = false;
	}

	function handleConfigConfirm(config: Record<string, string>) {
		extensionArguments = {
			'rook-ceph-cluster': { values: config }
		};
		openConfirmDialog = true;
	}
</script>

<div class="flex w-full flex-col gap-4">
	<Progress
		value={required ? (installed * 100) / required : 0}
		class={cn(
			installed == required
				? 'bg-green-700/20 **:data-[slot="progress-indicator"]:bg-green-700'
				: installed > 0
					? 'bg-yellow-500/20 **:data-[slot="progress-indicator"]:bg-yellow-500'
					: 'bg-red-700/20 **:data-[slot="progress-indicator"]:bg-red-700'
		)}
	/>
	<div class="flex items-center gap-2">
		<div class="flex h-full items-center justify-between gap-4">
			<div class="rounded-lg bg-primary/10 p-2">
				<Icon icon={extensionsBundleConfigurations[extensionsBundle].icon} class="size-8" />
			</div>
			<div>
				<div class="flex items-center gap-1">
					<h3 class="text-lg font-bold">{extensionsBundleConfigurations[extensionsBundle].name}</h3>
					<Badge>{scope}</Badge>
				</div>

				<p class="mt-1 text-sm text-muted-foreground">
					{extensionsBundleConfigurations[extensionsBundle].description}
				</p>
			</div>
		</div>
		<div class="ml-auto flex flex-col items-center justify-between gap-4">
			<p class="whitespace-nowrap text-muted-foreground">{installed} / {required}</p>
			<div class="ml-auto flex gap-2">
				<button
					class={cn(buttonVariants({ variant: 'default', size: 'sm' }), 'w-full')}
					disabled={isLoading}
					onclick={(e) => {
						e.stopPropagation();
						handleInstallClick();
					}}
				>
					{isInstallMode ? m.install() : m.extensions_upgrade()}
				</button>

				<AlertDialog.Root bind:open={openConfirmDialog}>
					<AlertDialog.Trigger class="hidden" />
					<AlertDialog.Content>
						<AlertDialog.Header>
							<AlertDialog.Title>
								{isInstallMode ? m.install() : m.extensions_upgrade()}
							</AlertDialog.Title>
							<AlertDialog.Description>
								{isInstallMode
									? m.install_extensions_description({ extensionsBundle: extensionsBundle })
									: m.upgrade_extensions_description({ extensionsBundle: extensionsBundle })}
							</AlertDialog.Description>
						</AlertDialog.Header>
						<AlertDialog.Footer>
							<AlertDialog.Cancel>{m.cancel()}</AlertDialog.Cancel>
							<AlertDialog.Action onclick={onConfirm}>
								{m.confirm()}
							</AlertDialog.Action>
						</AlertDialog.Footer>
					</AlertDialog.Content>
				</AlertDialog.Root>
			</div>
		</div>
	</div>
</div>

{#if extensionsBundle === 'storage'}
	<InstallRookConfig
		bind:open={openConfigDialog}
		{scope}
		{extensionsBundle}
		{extensions}
		{updator}
		onConfirm={handleConfigConfirm}
	/>
{/if}
