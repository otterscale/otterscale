<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onDestroy, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import { ApplicationService, type Job } from '$lib/api/application/v1/application_pb';
	import * as Loading from '$lib/components/custom/loading';
	import { ReloadManager } from '$lib/components/custom/reloader';

	import { DataTable } from './data-table/index';
</script>

<script lang="ts">
	let { scope }: { scope: string } = $props();

	const transport: Transport = getContext('transport');
	const applicationClient = createClient(ApplicationService, transport);

	let selectedNamespace = $state('default');

	const jobs = writable<Job[]>([]);
	async function fetch() {
		try {
			const response = await applicationClient.listJobs({
				scope: scope,
				namespace: selectedNamespace
			});
			jobs.set(response.jobs);
		} catch (error) {
			console.error('Failed to fetch jobs:', error);
		}
	}
	const reloadManager = new ReloadManager(fetch);

	let isMounted = $state(false);
	onMount(async () => {
		await fetch();
		isMounted = true;
		reloadManager.start();
	});
	onDestroy(() => {
		reloadManager.stop();
	});
	$effect(() => {
		if (selectedNamespace) {
			reloadManager.force();
		}
	});
</script>

<main class="space-y-4 py-4">
	{#if isMounted}
		<DataTable {jobs} {scope} bind:selectedNamespace {reloadManager} />
	{:else}
		<Loading.DataTable />
	{/if}
</main>
