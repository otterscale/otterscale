<script lang="ts">
	import { ApplicationService, type Application } from '$lib/api/application/v1/application_pb';
	import * as Loading from '$lib/components/custom/loading';
	import { ReloadManager, Reloader } from '$lib/components/custom/reloader';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onDestroy, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import { DataTable } from './data-table/index';

	import { listLargeLanguageModels, type LargeLangeageModel } from './protobuf.svelte';

	let { scopeUuid, facilityName }: { scopeUuid: string; facilityName: string } = $props();

	const transport: Transport = getContext('transport');
	let isMounted = $state(false);

	const largeLanguageModels = writable<LargeLangeageModel[]>([]);

	// const applicationClient = createClient(ApplicationService, transport);
	const reloadManager = new ReloadManager(() => {
		largeLanguageModels.set(listLargeLanguageModels());
		// applicationClient
		// 	.listApplications({
		// 		scopeUuid: scopeUuid,
		// 		facilityName: facilityName
		// 	})
		// 	.then((response) => {
		// 		applications.set(response.applications);
		// 	});
	});

	onMount(() => {
		largeLanguageModels.set(listLargeLanguageModels());
		isMounted = true;
		// applicationClient
		// 	.listApplications({
		// 		scopeUuid: scopeUuid,
		// 		facilityName: facilityName
		// 	})
		// 	.then((response) => {
		// 		largeLanguageModels.set(response.applications);
		// 		isMounted = true;
		// 	})
		// 	.catch((error) => {
		// 		console.error('Error during initial data load:', error);
		// 	});

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
