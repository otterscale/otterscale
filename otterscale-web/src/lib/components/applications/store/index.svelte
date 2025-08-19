<script lang="ts" module>
	import {
		ApplicationService,
		type Application_Chart,

		type Application_Release

	} from '$lib/api/application/v1/application_pb';
	import * as Loading from '$lib/components/custom/loading';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import { CommerceStore } from './commerce-store/index';
</script>

<script lang="ts">
	const transport: Transport = getContext('transport');

	let charts = $state(writable<Application_Chart[]>([]));
	let releases = $state(writable<Application_Release[]>([]));
	let isMounted = $state(false);

	const applicationClient = createClient(ApplicationService, transport);

	onMount(async () => {
		await applicationClient
			.listCharts({})
			.then((response) => {
				charts.set(response.charts.sort((p, n) => p.name.localeCompare(n.name)));
				isMounted = true;
			})
			.catch((error) => {
				console.error('Error during initial data load:', error);
			});
		await applicationClient
			.listReleases({})
			.then((response) => {
				releases.set(response.releases);
				isMounted = true;
			})
	});
</script>

{#if isMounted}
	<CommerceStore bind:charts />
{:else}
	<Loading.ApplicationStore />
{/if}
