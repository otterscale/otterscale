<script lang="ts">
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import { createClient, type Transport } from '@connectrpc/connect';

	import {
		Nexus,
		type Application_Release,
		type Application_Chart,
		type Facility,
		type Facility_Charm,
		type Facility_Charm_Artifact,
		type Scope
	} from '$gen/api/nexus/v1/nexus_pb';

	import { Store } from '$lib/components/otterscale/index';
	import { PageLoading } from '$lib/components/otterscale/ui/index';

	const transport: Transport = getContext('transportNEW');
	const client = createClient(Nexus, transport);

	const chartsStore = writable<Application_Chart[]>([]);
	const chartsLoading = writable(true);
	async function fetchCharts() {
		try {
			const response = await client.listCharts({});
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
			const response = await client.listReleases({});
			releasesStore.set(response.releases);
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			releasesLoading.set(false);
		}
	}

	const charmsStore = writable<Facility_Charm[]>([]);
	const charmsLoading = writable(true);
	async function fetchCharms() {
		try {
			const response = await client.listCharms({});
			charmsStore.set(response.charms);
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			charmsLoading.set(false);
		}
	}

	const scopesStore = writable<Scope[]>([]);
	const scopesLoading = writable(true);
	async function fetchScopes() {
		try {
			const response = await client.listScopes({});
			scopesStore.set(response.scopes);
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			scopesLoading.set(false);
		}
	}

	const facilitiesStore = writable<Facility[]>([]);
	const facilitiesLoading = writable(true);
	async function fetchFacilities(scopeUuid: string) {
		try {
			const response = await client.listFacilities({
				scopeUuid: scopeUuid
			});
			facilitiesStore.set(response.facilities);
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			facilitiesLoading.set(false);
		}
	}

	let mounted = false;
	onMount(async () => {
		try {
			await fetchCharts();
			await fetchReleases();
			await fetchCharms();
			await fetchScopes();
			for (const scope of $scopesStore) {
				const facilitiesResponse = await client.listFacilities({
					scopeUuid: scope.uuid
				});
				facilitiesStore.update((facilities) => [...facilities, ...facilitiesResponse.facilities]);
			}
		} catch (error) {
			console.error('Error during initial data load:', error);
		}

		mounted = true;
	});
</script>

{#if mounted}
	<Store
		charts={$chartsStore}
		releases={$releasesStore}
		charms={$charmsStore}
		facilities={$facilitiesStore}
	/>
{:else}
	<PageLoading />
{/if}
