<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onDestroy, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import { ApplicationService, type Secret } from '$lib/api/application/v1/application_pb';
	import * as Loading from '$lib/components/custom/loading';
	import { ReloadManager } from '$lib/components/custom/reloader';

	import { DataTable } from './data-table/index';
</script>

<script lang="ts">
	let { scope }: { scope: string } = $props();

	const transport: Transport = getContext('transport');
	const applicationClient = createClient(ApplicationService, transport);

	let selectedNamespace = $state('default');

	const secrets = writable<Secret[]>([]);
	async function fetch() {
		applicationClient
			.listSecrets({
				scope: scope,
				namespace: selectedNamespace
			})
			.then((response) => {
				secrets.set(response.secrets);
			})
			.catch((error) => {
				console.error('Failed to fetch secrets:', error);
			});
	}
	const reloadManager = new ReloadManager(fetch, false);

	let isMounted = $state(false);
	$effect(() => {
		if (selectedNamespace) {
			reloadManager.force();
		}
	});
	onMount(async () => {
		await fetch();
		isMounted = true;
	});
	onDestroy(() => {
		reloadManager.stop();
	});
</script>

<main class="space-y-4 py-4">
	{#if isMounted}
		<DataTable {secrets} {scope} bind:selectedNamespace {reloadManager} />
	{:else}
		<Loading.DataTable />
	{/if}
</main>
