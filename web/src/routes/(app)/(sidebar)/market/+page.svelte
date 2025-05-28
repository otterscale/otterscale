<script lang="ts">
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { ScopeService, type Scope } from '$gen/api/scope/v1/scope_pb';
	import {
		FacilityService,
		type Facility,
		type Facility_Charm
	} from '$gen/api/facility/v1/facility_pb';
	import {
		ApplicationService,
		type Application_Release,
		type Application_Chart
	} from '$gen/api/application/v1/application_pb';
	import { Store } from '$lib/components/otterscale/index';
	import { PageLoading } from '$lib/components/otterscale/ui/index';

	const transport: Transport = getContext('transport');
	const scopeClient = createClient(ScopeService, transport);
	const facilityClient = createClient(FacilityService, transport);
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

	const charmsStore = writable<Facility_Charm[]>([]);
	const charmsLoading = writable(true);
	async function fetchCharms() {
		try {
			const response = await facilityClient.listCharms({});
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
			const response = await scopeClient.listScopes({});
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
			const response = await facilityClient.listFacilities({
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
				const facilitiesResponse = await facilityClient.listFacilities({
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
