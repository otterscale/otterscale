<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onDestroy, onMount, setContext, type Snippet } from 'svelte';
	import { writable } from 'svelte/store';

	import { DataTable } from './data-table';
	import Pickers from './pickers.svelte';

	import { ConfigurationService, type TestResult } from '$lib/api/configuration/v1/configuration_pb';
	import * as Loading from '$lib/components/custom/loading';
	import { ReloadManager } from '$lib/components/custom/reloader';
</script>

<script lang="ts">
	let { trigger }: { trigger: Snippet } = $props();

	const transport: Transport = getContext('transport');

	const testResults = writable<TestResult[]>([]);
	let isMounted = $state(false);
	let mode = $state('get');

	const client = createClient(ConfigurationService, transport);
	const reloadManager = new ReloadManager(() => {
		client.listTestResults({}).then((response) => {
			testResults.set(response.testResults);
		});
	});
	setContext('reloadManager', reloadManager);

	onMount(() => {
		client
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
		<div class="flex items-center justify-between gap-2">
			{@render trigger()}
			<Pickers bind:selectedMode={mode} />
		</div>
		<DataTable {mode} {testResults} {reloadManager} />
	{:else}
		<Loading.DataTable />
	{/if}
</main>
