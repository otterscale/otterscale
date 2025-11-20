<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onDestroy, onMount, type Snippet } from 'svelte';
	import { writable } from 'svelte/store';

	import {
		ConfigurationService,
		type TestResult
	} from '$lib/api/configuration/v1/configuration_pb';
	import * as Loading from '$lib/components/custom/loading';
	import { ReloadManager } from '$lib/components/custom/reloader';

	import { DataTable } from './data-table';
	import Pickers from './pickers.svelte';
</script>

<script lang="ts">
	let { scope, selectedTab, trigger }: { scope: string; selectedTab: string; trigger: Snippet } =
		$props();

	let mode = $state('get');

	const transport: Transport = getContext('transport');
	const client = createClient(ConfigurationService, transport);
	const testResults = writable<TestResult[]>([]);

	async function fetch() {
		try {
			const response = await client.listTestResults({});
			testResults.set(response.testResults.filter((result) => result.kind.case === 'warp'));
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	}
	const reloadManager = new ReloadManager(fetch);

	let isMounted = $state(false);
	onMount(async () => {
		await fetch();
		isMounted = true;
	});
	onDestroy(() => {
		reloadManager.stop();
	});
</script>

<main>
	{#if isMounted && selectedTab === 'object-storage-test'}
		<div class="flex w-full items-center justify-between gap-2">
			{@render trigger()}
			<Pickers bind:selectedMode={mode} />
		</div>
		<DataTable {scope} {mode} {testResults} {reloadManager} />
	{:else}
		<Loading.DataTable />
	{/if}
</main>
