<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onDestroy, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import { RegistryService, type Repository } from '$lib/api/registry/v1/registry_pb';
	import * as Loading from '$lib/components/custom/loading';
	import { ReloadManager } from '$lib/components/custom/reloader';

	import { DataTable } from './data-table/index';
</script>

<script lang="ts">
	let { scope }: { scope: string } = $props();

	const transport: Transport = getContext('transport');
	const registryClient = createClient(RegistryService, transport);

	const repositories = writable<Repository[]>([]);
	async function fetch() {
		try {
			const response = await registryClient.listRepositories({
				scope
			});
			repositories.set(response.repositories);
		} catch (error) {
			console.error('Failed to fetch repositories:', error);
		}
	}

	const reloadManager = new ReloadManager(fetch);

	let isLoaded = $state(false);
	onMount(async () => {
		await fetch();
		isLoaded = true;
	});
	onDestroy(() => {
		reloadManager.stop();
	});
</script>

<main class="space-y-4 py-4">
	{#if isLoaded}
		<DataTable {repositories} {scope} {reloadManager} />
	{:else}
		<Loading.DataTable />
	{/if}
</main>
