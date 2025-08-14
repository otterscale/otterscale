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

	// Configuration for Kubernetes services
	const KUBERNETES_SERVICES = {
		kubernetesControlPlane: {
			name: 'kubernetes-control-plane',
			icon: 'ph:compass',
			title: m.control_planes(),
			gridClass: 'col-span-3 row-span-2'
		},
		kubernetesWorker: {
			name: 'kubernetes-worker',
			icon: 'ph:cube',
			title: m.workers(),
			gridClass: 'col-span-3 row-span-2'
		},
		kubeapiLoadBalancer: {
			name: 'kubeapi-load-balancer',
			icon: 'ph:scales',
			title: m.load_balancers(),
			gridClass: 'col-span-2'
		},
		etcd: {
			name: 'etcd',
			icon: 'ph:brackets-curly',
			title: 'etcd',
			gridClass: 'col-span-2'
		},
		easyrsa: {
			name: 'easyrsa',
			icon: 'ph:certificate',
			title: 'easyrsa',
			gridClass: 'col-span-2'
		}
	} as const;

	// Set breadcrumb navigation
	breadcrumb.set({
		parents: [dynamicPaths.setupScope(page.params.scope)],
		current: dynamicPaths.setupScopeKubernetes(page.params.scope)
	});

	// API setup
	const transport: Transport = getContext('transport');
	const facilityClient = createClient(FacilityService, transport);
	const facilitiesStore = writable<Facility[]>([]);

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
	<SetupScopeGrid facilities={$facilitiesStore} services={KUBERNETES_SERVICES} />
</div>
