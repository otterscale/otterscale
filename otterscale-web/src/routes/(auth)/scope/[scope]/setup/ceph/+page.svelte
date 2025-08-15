<script lang="ts">
	import { getContext, onDestroy, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { page } from '$app/state';
	import { FacilityService, type Facility } from '$lib/api/facility/v1/facility_pb';
	import { SetupScopeGrid } from '$lib/components/setup';
	import { m } from '$lib/paraglide/messages';
	import { dynamicPaths } from '$lib/path';
	import { activeScope, breadcrumb } from '$lib/stores';

	// Configuration for Ceph services
	const CEPH_SERVICES = {
		monitors: {
			name: 'ceph-mon',
			icon: 'ph:binoculars',
			title: m.monitors(),
			gridClass: 'col-span-3 row-span-2'
		},
		osds: {
			name: 'ceph-osd',
			icon: 'ph:hard-drives',
			title: m.osds(),
			gridClass: 'col-span-3 row-span-2'
		},
		fileSystem: {
			name: 'ceph-fs',
			icon: 'ph:tree-view',
			title: m.file_system(),
			gridClass: 'col-span-2'
		},
		objectGateway: {
			name: 'ceph-radosgw',
			icon: 'ph:traffic-sign',
			title: m.object_gateway(),
			gridClass: 'col-span-2'
		},
		networkFileSystem: {
			name: 'ceph-nfs',
			icon: 'ph:network',
			title: m.network_file_system(),
			gridClass: 'col-span-2'
		}
	} as const;

	// Set breadcrumb navigation
	breadcrumb.set({
		parents: [dynamicPaths.setupScope(page.params.scope)],
		current: dynamicPaths.setupScopeCeph(page.params.scope)
	});

	// API setup
	const transport: Transport = getContext('transport');
	const facilityClient = createClient(FacilityService, transport);
	const facilitiesStore = writable<Facility[]>([]);

	// State
	let autoRefresh = $state(true);

	async function fetchFacilities(uuid: string) {
		try {
			const response = await facilityClient.listFacilities({ scopeUuid: uuid });
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
	<SetupScopeGrid facilities={$facilitiesStore} services={CEPH_SERVICES} bind:autoRefresh />
</div>
