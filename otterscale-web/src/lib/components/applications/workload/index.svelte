<script lang="ts" module>
	import { ApplicationService, type Application } from '$lib/api/application/v1/application_pb';
	import { DataTable as DataTableLoading } from '$lib/components/custom/loading';
	import * as Reloader from '$lib/components/custom/reloader';
	import { currentKubernetes } from '$lib/stores';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onDestroy, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import { DataTable } from './data-table/index';
</script>

<script lang="ts">
	const transport: Transport = getContext('transport');
	const applicationClient = createClient(ApplicationService, transport);

	const applications = writable<Application[]>([]);
	const reloadManager = new Reloader.ReloadManager(() => {
		applicationClient
			.listApplications({
				scopeUuid: $currentKubernetes?.scopeUuid,
				facilityName: $currentKubernetes?.name
			})
			.then((response) => {
				applications.set(response.applications);
			});
	});

	let isMounted = $state(false);
	onMount(() => {
		console.log($currentKubernetes?.scopeUuid, $currentKubernetes?.name);
		applicationClient
			.listApplications({
				scopeUuid: $currentKubernetes?.scopeUuid,
				facilityName: $currentKubernetes?.name
			})
			.then((response) => {
				applications.set(response.applications);
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

<main class="space-y-4">
	{#if !isMounted}
		<DataTableLoading />
	{:else}
		<div class="flex justify-end">
			<Reloader.Root {reloadManager} />
		</div>
		<DataTable {applications} />
	{/if}
</main>
