<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import {
		type Application_Chart,
		type Application_Release,
		ApplicationService
	} from '$lib/api/application/v1/application_pb';
	import * as Loading from '$lib/components/custom/loading';

	import { CommerceStore } from './commerce-store/index';
</script>

<script lang="ts">
	let { scope }: { scope: string } = $props();

	const transport: Transport = getContext('transport');

	let charts = $state(writable<Application_Chart[]>([]));
	let releases = $state(writable<Application_Release[]>([]));
	let isChartsLoading = $state(true);
	let isReleasesLoading = $state(true);
	let isMounted = $derived(!isChartsLoading && !isReleasesLoading);

	const applicationClient = createClient(ApplicationService, transport);

	onMount(async () => {
		await applicationClient
			.listCharts({})
			.then((response) => {
				charts.set(response.charts.sort((p, n) => p.name.localeCompare(n.name)));
				isChartsLoading = false;
			})
			.catch((error) => {
				console.error('Error during initial data load:', error);
			});
		await applicationClient
			.listReleases({
				scope: scope
			})
			.then((response) => {
				releases.set(response.releases);
				isReleasesLoading = false;
			})
			.catch((error) => {
				console.error('Error during initial data load:', error);
			});
	});
</script>

{#if isMounted}
	<CommerceStore {scope} bind:charts bind:releases />
{:else}
	<Loading.ApplicationStore />
{/if}
