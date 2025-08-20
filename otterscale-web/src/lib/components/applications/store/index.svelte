<script lang="ts">
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
	import { currentKubernetes } from '$lib/stores';

	const transport: Transport = getContext('transport');

	const charts = writable<Application_Chart[]>([]);
	const releases = writable<Application_Release[]>([]);
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
				scopeUuid: $currentKubernetes?.scopeUuid,
				facilityName: $currentKubernetes?.name
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
	<CommerceStore {charts} {releases} />
{:else}
	<Loading.ApplicationStore />
{/if}
