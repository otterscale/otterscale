<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onDestroy, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import { ApplicationService } from '$lib/api/application/v1/application_pb';
	import * as Loading from '$lib/components/custom/loading';
	import { ReloadManager } from '$lib/components/custom/reloader';

	import { DataTable } from './data-table/index';
	import type { Service } from './types';
</script>

<script lang="ts">
	let { scope }: { scope: string } = $props();

	const transport: Transport = getContext('transport');
	let isMounted = $state(false);

	const services = writable<Service[]>([]);

	const applicationClient = createClient(ApplicationService, transport);
	const reloadManager = new ReloadManager(() => {
		applicationClient
			.listApplications({
				scope: scope
			})
			.then((response) => {
				services.set(
					response.applications.flatMap((application) =>
						application.services.map((service) => ({
							...service,
							hostname: response.hostname
						}))
					)
				);
			})
			.catch((error) => {
				console.error('Error during data loading:', error);
			});
	});

	onMount(() => {
		applicationClient
			.listApplications({
				scope: scope
			})
			.then((response) => {
				services.set(
					response.applications.flatMap((application) =>
						application.services.map((service) => ({
							...service,
							hostname: response.hostname
						}))
					)
				);
				isMounted = true;
			})
			.catch((error) => {
				console.error('Error during initial data load:', error);
			});

		reloadManager.start(); //TODO: Remove
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
