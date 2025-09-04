<script lang="ts" module>
	import { BISTService, type TestResult } from '$lib/api/bist/v1/bist_pb';
	import * as Loading from '$lib/components/custom/loading';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onDestroy, onMount, setContext, type Snippet } from 'svelte';
	import { writable } from 'svelte/store';
	import { DataTable } from './data-table';
</script>

<script lang="ts">
	let { trigger }: { trigger: Snippet } = $props();

	const transport: Transport = getContext('transport');

	const testResults = writable<TestResult[]>([]);
	let isMounted = $state(false);

	const bistClient = createClient(BISTService, transport);
	const reloadManager = new ReloadManager(() => {
		bistClient.listTestResults({}).then((response) => {
			testResults.set(response.testResults);
		});
	});
	setContext('reloadManager', reloadManager);

	onMount(() => {
		bistClient
			.listTestResults({})
			.then((response) => {
				testResults.set(response.testResults);
				isMounted = true;
			})
			.catch((error) => {
				console.error('Error during initial data load:', error);
			});

		reloadManager.start();
	});
	onDestroy(() => {
		reloadManager.stop();
	});
</script>

<main class="space-y-4 py-4">
	{#if isMounted}
		{@render trigger()}
		<DataTable {testResults} {reloadManager} />
	{:else}
		<Loading.DataTable />
	{/if}
</main>
