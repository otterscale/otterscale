<script lang="ts">
	import { getContext, onDestroy, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { page } from '$app/state';
	import { FacilityService, type Facility } from '$lib/api/facility/v1/facility_pb';
	import { dynamicPaths } from '$lib/path';
	import { activeScope, breadcrumb } from '$lib/stores';
	import { SetupScope } from '$lib/components/setup';

	// Set breadcrumb navigation
	breadcrumb.set({
		parents: [],
		current: dynamicPaths.setupScope(page.params.scope)
	});

	// API setup
	const transport: Transport = getContext('transport');
	const facilityClient = createClient(FacilityService, transport);
	const facilitiesStore = writable<Facility[]>([]);

	// State
	let autoRefresh = $state(true);

	async function fetchFacilities(uuid: string) {
		try {
			const response = await facilityClient.listFacilities({
				scopeUuid: uuid
			});
			facilitiesStore.set(response.facilities);
		} catch (error) {
			console.error('Error fetching facilities:', error);
		}
	}

	onMount(async () => {
		const unsubscribe = activeScope.subscribe(async (scope) => {
			if (scope) {
				await fetchFacilities(scope.uuid);
			}
		});

		onDestroy(() => unsubscribe());
	});
</script>

<div class="mx-auto max-w-7xl min-w-7xl">
	<SetupScope facilities={$facilitiesStore} bind:autoRefresh />
</div>
