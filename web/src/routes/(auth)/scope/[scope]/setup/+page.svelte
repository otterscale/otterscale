<script lang="ts">
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext } from 'svelte';
	import { writable } from 'svelte/store';

	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { type Facility,FacilityService } from '$lib/api/facility/v1/facility_pb';
	import { SetupScope } from '$lib/components/setup';
	import { m } from '$lib/paraglide/messages';
	import { activeScope, breadcrumbs } from '$lib/stores';

	// Set breadcrumbs navigation
	breadcrumbs.set([
		{
			title: m.setup(),
			url: resolve('/(auth)/scope/[scope]/setup', { scope: page.params.scope! })
		}
	]);

	// API setup
	const transport: Transport = getContext('transport');
	const facilityClient = createClient(FacilityService, transport);
	const facilitiesStore = writable<Facility[]>([]);

	// State & Timer
	let autoRefresh = $state(false);
	let refreshInterval: NodeJS.Timeout | null = null;

	async function fetchFacilities(scope: string) {
		try {
			const response = await facilityClient.listFacilities({
				scope: scope
			});
			facilitiesStore.set(response.facilities);
		} catch (error) {
			console.error('Error fetching facilities:', error);
		}
	}

	$effect(() => {
		const scope = $activeScope;

		// Clear existing interval
		if (refreshInterval) {
			clearInterval(refreshInterval);
			refreshInterval = null;
		}

		if (scope) {
			fetchFacilities(scope.name);

			// Setup auto-refresh if enabled
			if (autoRefresh) {
				refreshInterval = setInterval(() => {
					fetchFacilities(scope.name);
				}, 3000);
			}
		}

		// Cleanup on effect destruction
		return () => {
			if (refreshInterval) {
				clearInterval(refreshInterval);
				refreshInterval = null;
			}
		};
	});
</script>

<div class="mx-auto max-w-7xl min-w-7xl">
	<SetupScope facilities={$facilitiesStore} bind:autoRefresh />
</div>
