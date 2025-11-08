<script lang="ts">
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext } from 'svelte';
	import { writable } from 'svelte/store';

	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { FacilityService, type Facility } from '$lib/api/facility/v1/facility_pb';
	import { SetupScopeGrid } from '$lib/components/setup';
	import { m } from '$lib/paraglide/messages';
	import { activeScope, breadcrumbs } from '$lib/stores';

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

	// Set breadcrumbs navigation
	breadcrumbs.set([
		{
			title: m.setup(),
			url: resolve('/(auth)/scope/[scope]/setup', { scope: page.params.scope! })
		},
		{
			title: m.kubernetes(),
			url: resolve('/(auth)/scope/[scope]/setup/kubernetes', { scope: page.params.scope! })
		}
	]);

	// API setup
	const transport: Transport = getContext('transport');
	const facilityClient = createClient(FacilityService, transport);
	const facilities = writable<Facility[]>([]);

	// State & Timer
	let autoRefresh = $state(false);
	let refreshInterval: NodeJS.Timeout | null = null;

	async function fetchFacilities(scope: string) {
		try {
			const response = await facilityClient.listFacilities({ scope: scope });
			facilities.set(response.facilities);
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
	<SetupScopeGrid {facilities} services={KUBERNETES_SERVICES} bind:autoRefresh />
</div>
