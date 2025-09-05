<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onDestroy, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import { DataTable } from './data-table/index';
	import { Statistics } from './statistics';

	import { ApplicationService, type Application } from '$lib/api/application/v1/application_pb';
	import * as Loading from '$lib/components/custom/loading';
	import { ReloadManager } from '$lib/components/custom/reloader';
</script>

<script lang="ts">
	let { scopeUuid, facilityName }: { scopeUuid: string; facilityName: string } = $props();

	const transport: Transport = getContext('transport');
	let isMounted = $state(false);

	const applications = writable<Application[]>([]);

	const applicationClient = createClient(ApplicationService, transport);
	const reloadManager = new ReloadManager(() => {
		applicationClient
			.listApplications({
				scopeUuid: scopeUuid,
				facilityName: facilityName,
			})
			.then((response) => {
				applications.set(response.applications);
			});
	});

	onMount(() => {
		applicationClient
			.listApplications({
				scopeUuid: scopeUuid,
				facilityName: facilityName,
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

<main class="space-y-4 py-4">
	{#if isMounted}
		<Statistics {scopeUuid} {facilityName} />
		<DataTable {applications} {reloadManager} />
	{:else}
		<Loading.DataTable />
	{/if}
</main>
