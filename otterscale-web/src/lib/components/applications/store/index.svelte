<script lang="ts" module>
	import {
		ApplicationService,
		type Application_Chart
	} from '$lib/api/application/v1/application_pb';
	import * as Loading from '$lib/components/custom/loading';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import { CommerceStore } from './commerce-store/index';
</script>

<script lang="ts">
	const transport: Transport = getContext('transport');
	const applicationClient = createClient(ApplicationService, transport);

	const charts = writable<Application_Chart[]>([]);

	let isMounted = $state(false);
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
	});
</script>

{#if isMounted}
	<CommerceStore {charts} />
{:else}
	<Loading.ApplicationStore />
{/if}
