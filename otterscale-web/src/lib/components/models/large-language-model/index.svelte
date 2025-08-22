<script lang="ts" module>
	import * as Loading from '$lib/components/custom/loading';
	import { ReloadManager, Reloader } from '$lib/components/custom/reloader';
	import { onDestroy, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import { DataTable } from './data-table/index';

	import { listLargeLanguageModels, type LargeLangeageModel } from './protobuf.svelte';
</script>

<script lang="ts">
	const largeLanguageModels = writable<LargeLangeageModel[]>([]);
	let isMounted = $state(false);
	const reloadManager = new ReloadManager(() => {
		largeLanguageModels.set(listLargeLanguageModels());
	});

	onMount(() => {
		largeLanguageModels.set(listLargeLanguageModels());
		isMounted = true;
		reloadManager.start();
	});
	onDestroy(() => {
		reloadManager.stop();
	});
</script>

<main class="space-y-4 py-4">
	{#if isMounted}
		<Reloader {reloadManager} />
		<DataTable {largeLanguageModels} />
	{:else}
		<Loading.DataTable />
	{/if}
</main>
