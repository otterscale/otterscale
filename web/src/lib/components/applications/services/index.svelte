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
	const applicationClient = createClient(ApplicationService, transport);
	const services = writable<Service[]>([]);

	async function fetch() {
		try {
			const response = await applicationClient.listApplications({
				scope: scope
			});
			services.set(
				response.applications.flatMap((application) =>
					application.services.map((service) => ({
						...service,
						hostname: response.hostname
					}))
				)
			);
		} catch (error) {
			console.error('Failed to fetch services:', error);
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

<main class="space-y-4 py-4">
	{#if isMounted}
		<DataTable {services} {reloadManager} />
	{:else}
		<Loading.DataTable />
	{/if}
</main>
