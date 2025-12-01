<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import { ApplicationService, type Release } from '$lib/api/application/v1/application_pb';
	import { type Chart, RegistryService } from '$lib/api/registry/v1/registry_pb';
	import * as Loading from '$lib/components/custom/loading';

	import { CommerceStore } from './commerce-store/index';
</script>

<script lang="ts">
	let { scope }: { scope: string } = $props();

	const transport: Transport = getContext('transport');
	const applicationClient = createClient(ApplicationService, transport);
	const registryClient = createClient(RegistryService, transport);

	const charts = writable<Chart[]>([]);
	const releases = writable<Release[]>([]);

	async function fetchCharts() {
		try {
			const response = await registryClient.listCharts({});
			charts.set(response.charts.sort((p, n) => p.name.localeCompare(n.name)));
		} catch (error) {
			console.error('Failed to fetch charts:', error);
		}
	}

	async function fetchReleases() {
		try {
			const response = await applicationClient.listReleases({
				scope: scope
			});
			releases.set(response.releases);
		} catch (error) {
			console.error('Failed to fetch releases:', error);
		}
	}

	async function fetch() {
		try {
			await Promise.all([fetchCharts(), fetchReleases()]);
		} catch (error) {
			console.error('Error during data loading:', error);
		}
	}

	let isLoaded = $state(false);
	onMount(async () => {
		await fetch();
		isLoaded = true;
	});
</script>

{#if isLoaded}
	<CommerceStore {scope} {charts} {releases} />
{:else}
	<Loading.ApplicationStore />
{/if}
