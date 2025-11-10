<script lang="ts">
	import Icon from '@iconify/svelte';

	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { PremiumTier_Level } from '$lib/api/environment/v1/environment_pb';
	import { type Facility } from '$lib/api/facility/v1/facility_pb';
	import ContainerImage from '$lib/assets/container.jpg';
	import DiskImage from '$lib/assets/disk.jpg';
	import { Button } from '$lib/components/ui/button';
	import { Label } from '$lib/components/ui/label';
	import { Switch } from '$lib/components/ui/switch';
	import { m } from '$lib/paraglide/messages';
	import { currentCeph, currentKubernetes, premiumTier } from '$lib/stores';

	// Props
	let {
		scope,
		facilities,
		autoRefresh = $bindable(true)
	}: {
		scope: string;
		facilities: Facility[];
		autoRefresh: boolean;
	} = $props();

	// Types
	interface StateConfig {
		state: string;
		details: string;
		color: string;
		textClass: string;
		icon: string;
	}

	interface ServiceComponent {
		channel: string;
		version: string;
		allUnits: number;
		activeUnits: number;
	}

	interface ServiceState {
		state: StateConfig;
		controlPlane?: ServiceComponent;
		worker?: ServiceComponent;
		mon?: ServiceComponent;
		osd?: ServiceComponent;
	}

	// Constants
	const STATE_CONFIGS = {
		active: {
			color: 'green',
			textClass: 'text-green-600 dark:text-green-400',
			icon: 'ph:check-bold'
		},
		waiting: {
			color: 'yellow',
			textClass: 'text-yellow-600 dark:text-yellow-400',
			icon: 'ph:spinner-gap'
		},
		blocked: {
			color: 'red',
			textClass: 'text-red-600 dark:text-red-400',
			icon: 'ph:exclamation-mark'
		},
		maintenance: {
			color: 'blue',
			textClass: 'text-blue-600 dark:text-blue-400',
			icon: 'ph:wrench'
		}
	} as const;

	const CHARM_NAMES = {
		kubernetesControlPlane: 'kubernetes-control-plane',
		kubernetesWorker: 'kubernetes-worker',
		cephMon: 'ceph-mon',
		cephOsd: 'ceph-osd'
	} as const;

	// Utility functions
	function createEmptyComponent(): ServiceComponent {
		return { channel: '', version: '', allUnits: 0, activeUnits: 0 };
	}

	function toStateConfig(state: string, details: string): StateConfig {
		const config = STATE_CONFIGS[state.toLowerCase() as keyof typeof STATE_CONFIGS];
		return config
			? { state, details, ...config }
			: {
					state,
					details,
					color: 'gray',
					textClass: 'text-gray-600 dark:text-gray-400',
					icon: 'ph:question-bold'
				};
	}

	function createServiceComponent(facility: Facility): ServiceComponent {
		const activeUnits = facility.units.filter((u) => u.workloadStatus?.state === 'active').length;
		return {
			channel: facility.channel || '',
			version: facility.version || '',
			allUnits: facility.units.length,
			activeUnits
		};
	}

	function updateServiceState(
		facility: Facility,
		currentState: StateConfig,
		onlyIfNotActive = false
	): StateConfig {
		const status = facility.status;
		if (!status || (onlyIfNotActive && status.state === 'active')) {
			return currentState;
		}
		return toStateConfig(status.state, status.details);
	}

	function findFacilityByCharm(charmName: string): Facility | undefined {
		return facilities.find((f) => f.charmName.includes(charmName) && f.units.length > 0);
	}

	function processServiceFacilities(charmMappings: Record<string, string>): ServiceState {
		const serviceState: ServiceState = {
			state: toStateConfig('fetching', 'loading...')
		};

		Object.entries(charmMappings).forEach(([key, charmName]) => {
			const facility = findFacilityByCharm(charmName);
			if (facility) {
				(serviceState as Record<string, any>)[key] = createServiceComponent(facility);
				const isSecondary = key !== Object.keys(charmMappings)[0];
				serviceState.state = updateServiceState(facility, serviceState.state, isSecondary);
			} else {
				(serviceState as Record<string, any>)[key] = createEmptyComponent();
			}
		});

		return serviceState;
	}

	// Computed values
	let kubernetes: ServiceState = $state({} as ServiceState);
	let ceph: ServiceState = $state({} as ServiceState);

	$effect(() => {
		kubernetes = processServiceFacilities({
			controlPlane: CHARM_NAMES.kubernetesControlPlane,
			worker: CHARM_NAMES.kubernetesWorker
		});

		ceph = processServiceFacilities({
			mon: CHARM_NAMES.cephMon,
			osd: CHARM_NAMES.cephOsd
		});
	});
</script>

<!-- Header Controls -->
<div class="mx-auto max-w-7xl min-w-7xl">
	<div class="grid w-full grid-cols-2 gap-4 sm:gap-6 lg:grid-cols-4">
		<div class="col-span-2 flex justify-end space-x-4 rounded-lg sm:space-x-6 lg:col-span-4">
			<Button variant="ghost" disabled={$premiumTier.level === PremiumTier_Level.BASIC}>
				<Icon icon="ph:plus" class="size-4" />
				{m.add_node()}
			</Button>
			<div class="flex items-center space-x-2">
				<Switch id="auto-update" bind:checked={autoRefresh} />
				<Label for="auto-update">{m.auto_update()}</Label>
			</div>
		</div>

		<!-- Kubernetes Section -->
		{#if $currentKubernetes}
			{@render kubernetesCards()}
		{/if}

		<!-- Ceph Section -->
		{#if $currentCeph}
			{@render cephCards()}
		{/if}
	</div>
</div>

<!-- Snippet definitions -->
{#snippet unitCount(component: ServiceComponent | undefined)}
	<div class="mb-8 flex space-x-1 text-3xl sm:mb-2 lg:text-5xl">
		<span>{component?.activeUnits}</span>

		{#if component && component.allUnits > component.activeUnits}
			<Icon icon="ph:arrow-up-bold" class="size-6 animate-bounce" />
		{/if}
	</div>
{/snippet}
<!-- TODO: copy here -->
{#snippet statusCard(serviceState: ServiceState)}
	<div
		class="relative row-span-2 flex flex-col justify-between overflow-hidden rounded-lg bg-muted p-4 shadow-sm md:p-6 lg:p-10"
	>
		<div
			class="absolute text-8xl tracking-tight text-nowrap text-primary/5 uppercase group-hover:hidden"
		>
			<Icon icon={serviceState.state?.icon} class="size-84" />
		</div>
		<div class="z-10 mb-8 flex flex-col space-y-2 text-3xl sm:mb-2 lg:text-5xl">
			<span
				class="flex space-x-2 truncate overflow-visible capitalize {serviceState.state?.textClass}"
			>
				<span>{serviceState.state?.state}</span>
			</span>
			<div class="text-xs tracking-tight text-muted-foreground capitalize md:text-base lg:text-lg">
				{serviceState.state?.details}
			</div>
		</div>
	</div>
{/snippet}

{#snippet kubernetesCards()}
	<!-- Kubernetes Main Card -->
	<a
		href={resolve('/(auth)/scope/[scope]/setup/kubernetes', { scope: scope })}
		class="group relative col-span-2 row-span-2 overflow-clip rounded-lg shadow-sm sm:max-lg:col-span-1"
	>
		<img
			src={ContainerImage}
			alt="container"
			class="absolute h-full w-full object-cover object-center"
		/>
		<div
			class="relative flex h-full w-full flex-col items-start justify-between bg-primary/20 p-4 text-primary-foreground transition-colors hover:bg-primary/30 md:p-6 lg:p-10 dark:text-primary"
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
		class="flex flex-col justify-between rounded-lg bg-muted p-4 shadow-sm sm:justify-end md:p-6 lg:p-10"
	>
		{@render unitCount(kubernetes.controlPlane)}
		<div class="flex items-center space-x-2 text-xs md:text-base lg:text-lg">
			<Icon icon="ph:compass" class="size-6" />
			<span>{m.control_planes()}</span>
		</div>
	</div>

	<!-- Kubernetes Status Card -->
	{@render statusCard(kubernetes)}

	<!-- Workers Card -->
	<div
		class="flex flex-col justify-between rounded-lg bg-muted p-4 shadow-sm sm:justify-end md:p-6 lg:p-10"
	>
		{@render unitCount(kubernetes.worker)}
		<div class="flex items-center space-x-2 text-xs md:text-base lg:text-lg">
			<Icon icon="ph:cube" class="size-6" />
			<span>{m.workers()}</span>
		</div>
	</div>
{/snippet}

{#snippet cephCards()}
	<!-- Monitors Card -->
	<div
		class="flex flex-col justify-between rounded-lg bg-muted p-4 shadow-sm sm:justify-end md:p-6 lg:p-10"
	>
		{@render unitCount(ceph.mon)}
		<div class="flex items-center space-x-2 text-xs md:text-base lg:text-lg">
			<Icon icon="ph:binoculars" class="size-6" />
			<span>{m.monitors()}</span>
		</div>
	</div>

	<!-- Ceph Status Card -->
	{@render statusCard(ceph)}

	<!-- Ceph Main Card -->
	<a
		href={resolve('/(auth)/scope/[scope]/setup/ceph', { scope: scope })}
		class="group relative col-span-2 row-span-2 overflow-clip rounded-lg shadow-sm sm:max-lg:col-span-1"
	>
		<img src={DiskImage} alt="disk" class="absolute h-full w-full object-cover object-center" />
		<div
			class="relative flex h-full w-full flex-col items-start justify-between gap-4 bg-primary/20 p-4 text-primary-foreground transition-colors hover:bg-primary/30 md:flex-row md:items-end md:p-6 lg:p-10 dark:text-primary"
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
		class="flex flex-col justify-between rounded-lg bg-muted p-4 shadow-sm sm:justify-end md:p-6 lg:p-10"
	>
		{@render unitCount(ceph.osd)}
		<div class="flex items-center space-x-2 text-xs md:text-base lg:text-lg">
			<Icon icon="ph:hard-drives" class="size-6" />
			<span>{m.osds()}</span>
		</div>
	</div>
{/snippet}
