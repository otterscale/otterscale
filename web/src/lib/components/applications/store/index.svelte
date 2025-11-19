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

	const charts = writable<Application_Chart[]>([]);
	const releases = writable<Application_Release[]>([]);

	let isChartsLoaded = $state(false);
	let isReleasesLoaded = $state(false);
	const isMounted = $derived(isChartsLoaded && isReleasesLoaded);

	const applicationClient = createClient(ApplicationService, transport);

	onMount(async () => {
		await applicationClient
			.listCharts({})
			.then((response) => {
				charts.set(response.charts.sort((p, n) => p.name.localeCompare(n.name)));
				isChartsLoaded = true;
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
				isReleasesLoaded = true;
			})
			.catch((error) => {
				console.error('Error during initial data load:', error);
			});
	});
</script>

{#if isMounted}
	<CommerceStore {scope} {charts} {releases} />
{:else}
	<Loading.ApplicationStore />
{/if}
