<script lang="ts" module>
	import { ApplicationService, type Application } from '$lib/api/application/v1/application_pb';
	import * as Loading from '$lib/components/custom/loading';
	import * as Reloader from '$lib/components/custom/reloader';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onDestroy, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import { DataTable } from './data-table/index';
	import { Statistics } from './statistics';
</script>

<script lang="ts">
	let { scopeUuid, facilityName }: { scopeUuid: string; facilityName: string } = $props();

	const transport: Transport = getContext('transport');
	const applicationClient = createClient(ApplicationService, transport);

	const applications = writable<Application[]>([]);
	const reloadManager = new Reloader.ReloadManager(() => {
		applicationClient
			.listApplications({
				scopeUuid: scopeUuid,
				facilityName: facilityName
			})
			.then((response) => {
				applications.set(response.applications);
			});
	});

	let isMounted = $state(false);
	onMount(() => {
		applicationClient
			.listApplications({
				scopeUuid: scopeUuid,
				facilityName: facilityName
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
		<Loading.DataTable />
	{:else}
		<Reloader.Root {reloadManager} />
		<Statistics {scopeUuid} {facilityName} />
		<DataTable {applications} />
	{/if}
</main>
