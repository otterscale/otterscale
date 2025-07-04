<script lang="ts">
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import { createClient, type Transport } from '@connectrpc/connect';
	import {
		ApplicationService,
		type Application_Release,
		type Application_Chart
	} from '$gen/api/application/v1/application_pb';
	import { Store } from '$lib/components/otterscale/index';
	import { PageLoading } from '$lib/components/otterscale/ui/index';

	const transport: Transport = getContext('transport');
	const applicationClient = createClient(ApplicationService, transport);

	const chartsStore = writable<Application_Chart[]>([]);
	const chartsLoading = writable(true);
	async function fetchCharts() {
		try {
			const response = await applicationClient.listCharts({});
			chartsStore.set(response.charts);
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			chartsLoading.set(false);
		}
	}

	const releasesStore = writable<Application_Release[]>([]);
	const releasesLoading = writable(true);
	async function fetchReleases() {
		try {
			const response = await applicationClient.listReleases({});
			releasesStore.set(response.releases);
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			releasesLoading.set(false);
		}
	}

	let mounted = false;
	onMount(async () => {
		try {
			await fetchCharts();
			await fetchReleases();
		} catch (error) {
			console.error('Error during initial data load:', error);
		}

		mounted = true;
	});
</script>

{#if mounted}
	<Store charts={$chartsStore} releases={$releasesStore} />
{:else}
	<PageLoading />
{/if}
