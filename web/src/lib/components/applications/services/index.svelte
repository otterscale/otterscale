<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onDestroy, onMount, setContext } from 'svelte';
	import { writable } from 'svelte/store';

	import { DataTable } from './data-table/index';
	import type { Service } from './types';

	import { ApplicationService } from '$lib/api/application/v1/application_pb';
	import * as Loading from '$lib/components/custom/loading';
	import { ReloadManager } from '$lib/components/custom/reloader';
</script>

<script lang="ts">
	let { scope, facility }: { scope: string; facility: string } = $props();

	const transport: Transport = getContext('transport');
	let isMounted = $state(false);

	const services = writable<Service[]>([]);

	const applicationClient = createClient(ApplicationService, transport);
	const reloadManager = new ReloadManager(() => {
		applicationClient
			.listApplications({
				scope: scope,
				facility: facility,
			})
			.then((response) => {
				services.set(
					response.applications.flatMap((application) =>
						application.services.map((service) => ({
							...service,
							publicAddress: response.publicAddress,
						})),
					),
				);
			})
			.catch((error) => {
				console.error('Error during data loading:', error);
			});
	});
	setContext('reloadManager', reloadManager);

	onMount(() => {
		applicationClient
			.listApplications({
				scope: scope,
				facility: facility,
			})
			.then((response) => {
				services.set(
					response.applications.flatMap((application) =>
						application.services.map((service) => ({
							...service,
							publicAddress: response.publicAddress,
						})),
					),
				);
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
		<DataTable {services} {reloadManager} />
	{:else}
		<Loading.DataTable />
	{/if}
</main>
