<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';
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

	async function synchronize() {
		try {
			const getRegistryURLResponse = await registryClient.getRegistryURL({
				scope: scope
			});
			registryClient.syncArtifactHub({
				registryUrl: getRegistryURLResponse.registryUrl
			});
		} catch (error) {
			console.error('Failed to synchronize with Artifact Hub:', error);
		}
	}

	onMount(async () => {});
</script>

<Button
	class="flex h-8 items-center gap-2"
	onclick={() => {
		toast.promise(
			async () => {
				await synchronize();
				const response = await registryClient.listCharts({
					scope: scope
				});
				charts.set(response.charts.sort((p, n) => p.name.localeCompare(n.name)));
			},
			{
				loading: 'Synchronizing with Artifact Hub...',
				success: 'Synchronization successful!',
				error: (error) => {
					console.error('Failed to synchronize with Artifact Hub:', error);
					return `Synchronization failed: ${(error as ConnectError).message}`;
				}
			}
		);
		synchronize();
	}}
>
	<Icon icon="ph:box-arrow-down" class="size-3" />
	{m.sync()}
</Button>
