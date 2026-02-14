<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import type { Writable } from 'svelte/store';
	import { toast } from 'svelte-sonner';

	import { type Chart, RegistryService } from '$lib/api/registry/v1/registry_pb';
	import { Button } from '$lib/components/ui/button';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let {
		scope,
		charts
	}: {
		scope: string;
		charts: Writable<Chart[]>;
	} = $props();

	const transport: Transport = getContext('transport');
	const registryClient = createClient(RegistryService, transport);

	let isSyncing = $state(false);

	function handleSync() {
		const promise = async () => {
			isSyncing = true;
			try {
				const getRegistryURLResponse = await registryClient.getRegistryURL({
					scope: scope
				});

				await registryClient.syncArtifactHub({
					registryUrl: getRegistryURLResponse.registryUrl
				});

				const response = await registryClient.listCharts({
					scope: scope
				});

				charts.set(response.charts.sort((p, n) => p.name.localeCompare(n.name)));
			} finally {
				isSyncing = false;
			}
		};

		toast.promise(promise(), {
			loading: 'Synchronizing with Artifact Hub...',
			success: 'Synchronization successful!',
			error: (error) => {
				console.error('Failed to synchronize with Artifact Hub:', error);
				const message = error instanceof ConnectError ? error.message : 'Unknown error';
				return `Synchronization failed: ${message}`;
			}
		});
	}
</script>

<Button variant="outline" class="flex items-center" disabled={isSyncing} onclick={handleSync}>
	<Icon icon="ph:arrows-clockwise" class={isSyncing ? 'animate-spin' : ''} />
	{m.sync()}
</Button>
