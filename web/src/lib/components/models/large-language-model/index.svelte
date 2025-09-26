<script lang="ts" module>
	import { onDestroy, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import Alert from './alert.svelte';
	import { DataTable } from './data-table/index';
	import { listLargeLanguageModels, type LargeLangeageModel } from './protobuf.svelte';

	import * as Loading from '$lib/components/custom/loading';
	import { ReloadManager } from '$lib/components/custom/reloader';
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
	<Alert />
	{#if isMounted}
		<DataTable {largeLanguageModels} {reloadManager} />
	{:else}
		<Loading.DataTable />
	{/if}
</main>
