<script lang="ts">
	import { getContext, onDestroy, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { page } from '$app/state';
	import { FacilityService, type Facility } from '$lib/api/facility/v1/facility_pb';
	import ContainerImage from '$lib/assets/container.jpg';
	import DiskImage from '$lib/assets/disk.jpg';
	import { m } from '$lib/paraglide/messages';
	import { dynamicPaths } from '$lib/path';
	import { activeScope, breadcrumb, currentCeph, currentKubernetes } from '$lib/stores';

	// Types
	interface StateConfig {
		state: string;
		details: string;
		color: string;
		textClass: string;
		icon: string;
	}

	interface ServiceState {
		state: StateConfig;
		controlPlane?: ServiceComponent;
		worker?: ServiceComponent;
		mon?: ServiceComponent;
		osd?: ServiceComponent;
	}

	interface ServiceComponent {
		channel: string;
		version: string;
		allUnits: number;
		activeUnits: number;
	}

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
	let kubernetes: ServiceState = {
		state: {} as StateConfig,
		controlPlane: createEmptyComponent(),
		worker: createEmptyComponent()
	};

	let ceph: ServiceState = {
		state: {} as StateConfig,
		mon: createEmptyComponent(),
		osd: createEmptyComponent()
	};

	// Helper functions
	function createEmptyComponent(): ServiceComponent {
		return {
			channel: '',
			version: '',
			allUnits: 0,
			activeUnits: 0
		};
	}

	function toStateConfig(state: string, details: string): StateConfig {
		const configs = {
			active: {
				state: 'active',
				details,
				color: 'green',
				textClass: 'text-green-600 dark:text-green-400',
				icon: 'ph:check-bold'
			},
			waiting: {
				state: 'waiting',
				details,
				color: 'yellow',
				textClass: 'text-yellow-600 dark:text-yellow-400',
				icon: 'ph:spinner-gap'
			},
			blocked: {
				state: 'blocked',
				details,
				color: 'red',
				textClass: 'text-red-600 dark:text-red-400',
				icon: 'ph:exclamation-mark'
			},
			maintenance: {
				state: 'maintenance',
				details,
				color: 'blue',
				textClass: 'text-blue-600 dark:text-blue-400',
				icon: 'ph:wrench'
			}
		};

		return (
			configs[state.toLowerCase() as keyof typeof configs] || {
				state,
				details,
				color: 'gray',
				textClass: 'text-gray-600 dark:text-gray-400',
				icon: 'ph:question-bold'
			}
		);
	}

	function updateServiceComponent(facility: Facility): ServiceComponent {
		return {
			channel: facility.channel || '',
			version: facility.version || '',
			allUnits: facility.units.length,
			activeUnits: facility.units.filter((u) => u.workloadStatus?.state === 'active').length
		};
	}

	function updateServiceState(
		facility: Facility,
		currentState: StateConfig,
		onlyIfNotActive = false
	): StateConfig {
		const status = facility.status;
		if (!status) return currentState;

		if (onlyIfNotActive && status.state === 'active') return currentState;

		return toStateConfig(status.state, status.details);
	}

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

	function processKubernetesFacilities(facilities: Facility[]) {
		// Control plane
		const controlPlaneFacility = facilities.find(
			(f) => f.charmName.includes('kubernetes-control-plane') && f.units.length > 0
		);
		if (controlPlaneFacility) {
			kubernetes.controlPlane = updateServiceComponent(controlPlaneFacility);
			kubernetes.state = updateServiceState(controlPlaneFacility, kubernetes.state);
		}

		// Worker
		const workerFacility = facilities.find(
			(f) => f.charmName.includes('kubernetes-worker') && f.units.length > 0
		);
		if (workerFacility) {
			kubernetes.worker = updateServiceComponent(workerFacility);
			kubernetes.state = updateServiceState(workerFacility, kubernetes.state, true);
		}
	}

	function processCephFacilities(facilities: Facility[]) {
		// Monitor
		const monFacility = facilities.find(
			(f) => f.charmName.includes('ceph-mon') && f.units.length > 0
		);
		if (monFacility) {
			ceph.mon = updateServiceComponent(monFacility);
			ceph.state = updateServiceState(monFacility, ceph.state);
		}

		// OSD
		const osdFacility = facilities.find(
			(f) => f.charmName.includes('ceph-osd') && f.units.length > 0
		);
		if (osdFacility) {
			ceph.osd = updateServiceComponent(osdFacility);
			ceph.state = updateServiceState(osdFacility, ceph.state, true);
		}
	}

	onMount(async () => {
		const unsubscribe = activeScope.subscribe(async (scope) => {
			if (scope) {
				await fetchFacilities(scope.uuid);

				const facilities = $facilitiesStore;
				processKubernetesFacilities(facilities);
				processCephFacilities(facilities);
			}
		});

		onDestroy(() => unsubscribe());
	});
</script>

<div class="mx-auto max-w-7xl min-w-7xl">
	<div class="grid w-full grid-cols-2 gap-4 sm:gap-6 lg:grid-cols-4">
		{#if $currentKubernetes}
			<!-- Kubernetes Main Card -->
			<a
				href={dynamicPaths.setupScopeKubernetes(page.params.scope).url}
				class="group relative col-span-2 row-span-2 overflow-clip rounded-lg shadow-sm sm:max-lg:col-span-1"
			>
				<img
					src={ContainerImage}
					alt="container"
					class="absolute h-full w-full object-cover object-center"
				/>
				<div
					class="bg-primary/20 text-primary-foreground hover:bg-primary/30 relative flex h-full w-full flex-col items-start justify-between p-4 transition-colors md:p-6 lg:p-10"
				>
					<div class="flex items-center gap-4">
						<Icon icon="logos:kubernetes" class="size-14" />
						<div class="flex flex-col">
							<span class="text-2xl font-semibold">Kubernetes</span>
							<span class="text-md">
								{kubernetes.controlPlane?.channel}
								{kubernetes.controlPlane?.version}
							</span>
						</div>
					</div>
					<div class="flex items-center text-xs font-medium md:text-base lg:text-lg">
						{m.details()}
						<Icon
							icon="ph:arrow-right-bold"
							class="ml-2 size-6 transition-transform group-hover:translate-x-0.5"
						/>
					</div>
				</div>
			</a>

			<!-- Control Planes Card -->
			<div
				class="bg-muted flex flex-col justify-between rounded-lg p-4 shadow-sm sm:justify-end md:p-6 lg:p-10"
			>
				<div class="mb-8 text-3xl sm:mb-2 lg:text-5xl">
					{kubernetes.controlPlane?.activeUnits}
					{#if kubernetes.controlPlane?.activeUnits !== kubernetes.controlPlane?.allUnits}
						-> {kubernetes.controlPlane?.allUnits}
					{/if}
				</div>
				<div class="text-xs md:text-base lg:text-lg">{m.control_planes()}</div>
			</div>

			<!-- Kubernetes Status Card -->
			<div
				class="bg-accent relative row-span-2 flex flex-col justify-between overflow-hidden rounded-lg p-4 shadow-sm md:p-6 lg:p-10"
			>
				<div
					class="text-primary/5 absolute text-8xl tracking-tight text-nowrap uppercase group-hover:hidden"
				>
					<Icon icon={kubernetes.state.icon} class="size-84" />
				</div>
				<div class="mb-8 flex flex-col space-y-2 text-3xl sm:mb-2 lg:text-5xl">
					<span
						class="flex space-x-2 truncate overflow-visible capitalize {kubernetes.state.textClass}"
					>
						<span>{kubernetes.state.state}</span>
					</span>
					<div
						class="text-muted-foreground text-xs tracking-tight capitalize md:text-base lg:text-lg"
					>
						{kubernetes.state.details}
					</div>
				</div>
			</div>

			<!-- Workers Card -->
			<div
				class="bg-accent flex flex-col justify-between rounded-lg p-4 shadow-sm sm:justify-end md:p-6 lg:p-10"
			>
				<div class="mb-8 text-3xl sm:mb-2 lg:text-5xl">
					{kubernetes.worker?.activeUnits}
					{#if kubernetes.worker?.activeUnits !== kubernetes.worker?.allUnits}
						-> {kubernetes.worker?.allUnits}
					{/if}
				</div>
				<div class="text-xs md:text-base lg:text-lg">{m.workers()}</div>
			</div>
		{/if}

		{#if $currentCeph}
			<!-- Monitors Card -->
			<div
				class="bg-accent flex flex-col justify-between rounded-lg p-4 shadow-sm sm:justify-end md:p-6 lg:p-10"
			>
				<div class="mb-8 text-3xl sm:mb-2 lg:text-5xl">
					{ceph.mon?.activeUnits}
					{#if ceph.mon?.activeUnits !== ceph.mon?.allUnits}
						-> {ceph.mon?.allUnits}
					{/if}
				</div>
				<div class="text-xs md:text-base lg:text-lg">{m.monitors()}</div>
			</div>

			<!-- Ceph Status Card -->
			<div
				class="bg-accent relative row-span-2 flex flex-col justify-between overflow-hidden rounded-lg p-4 shadow-sm md:p-6 lg:p-10"
			>
				<div
					class="text-primary/5 absolute text-8xl tracking-tight text-nowrap uppercase group-hover:hidden"
				>
					<Icon icon={ceph.state.icon} class="size-84" />
				</div>
				<div class="mb-8 flex flex-col space-y-2 text-3xl sm:mb-2 lg:text-5xl">
					<span class="flex space-x-2 truncate overflow-visible capitalize {ceph.state.textClass}">
						<span>{ceph.state.state}</span>
					</span>
					<div
						class="text-muted-foreground text-xs tracking-tight capitalize md:text-base lg:text-lg"
					>
						{ceph.state.details}
					</div>
				</div>
			</div>

			<!-- Ceph Main Card -->
			<a
				href={dynamicPaths.setupScopeCeph(page.params.scope).url}
				class="group relative col-span-2 row-span-2 overflow-clip rounded-lg shadow-sm sm:max-lg:col-span-1"
			>
				<img src={DiskImage} alt="disk" class="absolute h-full w-full object-cover object-center" />
				<div
					class="bg-primary/20 text-primary-foreground hover:bg-primary/30 relative flex h-full w-full flex-col items-start justify-between gap-4 p-4 transition-colors md:flex-row md:items-end md:p-6 lg:p-10"
				>
					<div class="flex items-center gap-4">
						<Icon icon="simple-icons:ceph" class="size-14 text-[#f0424d]" />
						<div class="flex flex-col">
							<span class="text-2xl font-semibold">Ceph</span>
							<span class="text-md">
								{ceph.mon?.channel}
								{ceph.mon?.version}
							</span>
						</div>
					</div>
					<div class="flex shrink-0 items-center text-xs font-medium md:text-base lg:text-lg">
						{m.details()}
						<Icon
							icon="ph:arrow-right-bold"
							class="ml-2 size-6 transition-transform group-hover:translate-x-0.5"
						/>
					</div>
				</div>
			</a>

			<!-- OSDs Card -->
			<div
				class="bg-accent flex flex-col justify-between rounded-lg p-4 shadow-sm sm:justify-end md:p-6 lg:p-10"
			>
				<div class="mb-8 text-3xl sm:mb-2 lg:text-5xl">
					{ceph.osd?.activeUnits}
					{#if ceph.osd?.activeUnits !== ceph.osd?.allUnits}
						-> {ceph.osd?.allUnits}
					{/if}
				</div>
				<div class="text-xs md:text-base lg:text-lg">{m.osds()}</div>
			</div>
		{/if}
	</div>
</div>
