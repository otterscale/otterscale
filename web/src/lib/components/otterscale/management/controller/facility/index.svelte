<script lang="ts">
	import AddFacilityUnits from './add-units.svelte';

	import UpdateFacility from './update.svelte';
	import DeleteFacility from './delete.svelte';
	import ExposeFacility from './expose.svelte';
	import AddCephUnits from './add-ceph-units.svelte';
	import AddKubernetesUnits from './add-kubernetes-units.svelte';

	import * as HoverCard from '$lib/components/ui/hover-card/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import * as Card from '$lib/components/ui/card';
	import { Progress } from '$lib/components/ui/progress/index.js';
	import * as Tabs from '$lib/components/ui/tabs';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { ManagementFacilityActions } from '$lib/components/otterscale/index';

	import {
		Nexus,
		type Machine_Placement,
		type AddFacilityUnitsRequest,
		type CreateFacilityRequest,
		type DeleteFacilityRequest,
		type Facility,
		type Facility_Status,
		type Facility_Unit,
		type Machine,
		type Machine_Constraint,
		type UpdateFacilityRequest
	} from '$gen/api/nexus/v1/nexus_pb';

	import Icon from '@iconify/svelte';

	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import * as AlertDialog from '$lib/components/ui/alert-dialog';
	import * as Popover from '$lib/components/ui/popover';
	import * as Select from '$lib/components/ui/select/index.js';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import * as Table from '$lib/components/ui/table/index.js';

	import { toast } from 'svelte-sonner';

	let {
		scopeUuid
	}: {
		scopeUuid: string;
	} = $props();

	const transport: Transport = getContext('transportNEW');
	const client = createClient(Nexus, transport);

	const facilitiesStore = writable<Facility[]>([]);
	const facilitiesLoading = writable(true);
	async function fetchFacilities() {
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

	const machinesStore = writable<Machine[]>([]);
	const machinesLoading = writable(true);
	async function fetchMachines() {
		try {
			const response = await client.listMachines({});
			machinesStore.set(response.machines);
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			machinesLoading.set(false);
		}
	}

	const FACILITY_CATEGORYS = {
		KUBERNETES: 'Kubernetes',
		CEPH: 'Ceph',
		OTHER: 'Customs'
	} as const;

	function getCategory(facility: Facility) {
		const charmName = facility.charmName;
		switch (true) {
			case charmName.includes('kubeapi-load-balancer'):
			case charmName.includes('kubernetes-worker'):
			case charmName.includes('etcd'):
			case charmName.includes('keepalived'):
			case charmName.includes('kubernetes-control-plane'):
				return FACILITY_CATEGORYS.KUBERNETES;
			case charmName.includes('ceph-mon'):
				return FACILITY_CATEGORYS.CEPH;
			default:
				return FACILITY_CATEGORYS.OTHER;
		}
	}

	function getLogoByCategory(category: string) {
		switch (category) {
			case FACILITY_CATEGORYS.KUBERNETES:
				return 'logos:kubernetes';
			case FACILITY_CATEGORYS.CEPH:
				return 'simple-icons:ceph';
			default:
				return '';
		}
	}

	const collapsibleOpen = $state(
		Object.values(FACILITY_CATEGORYS).reduce(
			(m, category) => {
				m[category] = true;
				return m;
			},
			{} as Record<string, boolean>
		)
	);

	let mounted = false;
	onMount(async () => {
		try {
			await fetchFacilities();
			await fetchMachines();
		} catch (error) {
			console.error('Error during initial data load:', error);
		}

		mounted = true;
	});
</script>

{#snippet collapsibilityHandler(facilityCategory: string)}
	<Button
		class="size-3"
		variant="ghost"
		onclick={() => (collapsibleOpen[facilityCategory] = !collapsibleOpen[facilityCategory])}
	>
		<Icon icon={collapsibleOpen[facilityCategory] ? 'ph:caret-up' : 'ph:caret-down'} />
	</Button>
{/snippet}

{#snippet AddBundleUnits(facilityCategory: string)}
	{#if facilityCategory == FACILITY_CATEGORYS.KUBERNETES}
		<AddKubernetesUnits {scopeUuid} ceph={{} as Facility} machines={$machinesStore} />
	{:else if facilityCategory == FACILITY_CATEGORYS.CEPH}
		<AddCephUnits {scopeUuid} ceph={{} as Facility} machines={$machinesStore} />
	{/if}
{/snippet}

<main>
	{@render StatisticFacilities($facilitiesStore)}

	<Table.Root>
		<Table.Header>
			<Table.Row class="*:text-sm *:font-light">
				<Table.Head>
					NAME
					<p class="text-xs">VERSION</p>
				</Table.Head>
				<Table.Head>UNITS</Table.Head>
				<Table.Head>STATUS</Table.Head>
				<Table.Head>REVISION</Table.Head>
				<Table.Head>CHARM NAME</Table.Head>
				<Table.Head></Table.Head>
				<Table.Head></Table.Head>
			</Table.Row>
		</Table.Header>
		<Table.Body>
			{#each Object.values(FACILITY_CATEGORYS) as facilityCategory}
				{#if $facilitiesStore.some((a) => getCategory(a) === facilityCategory)}
					<Table.Row class="border-none">
						<Table.Cell colspan={8} class="rounded-lg bg-secondary">
							<span class="flex items-center justify-between">
								<span class="flex items-center gap-2">
									<Icon icon={getLogoByCategory(facilityCategory)} class="size-7" />
									<h1 class="text-start text-lg">{facilityCategory}</h1>
								</span>
								<span class="flex items-center gap-2">
									{@render collapsibilityHandler(facilityCategory)}
								</span>
							</span>
						</Table.Cell>
					</Table.Row>

					{#if collapsibleOpen[facilityCategory]}
						<Table.Row class="border-none hover:bg-transparent">
							<Table.Cell colspan={8}>
								<span class="flex justify-end">
									{@render AddBundleUnits(facilityCategory)}
								</span>
							</Table.Cell>
						</Table.Row>
						{@const facilitiesByCategory = $facilitiesStore.filter(
							(a) => getCategory(a) === facilityCategory
						)}
						{#each facilitiesByCategory as facilityByCategory}
							<Table.Row class="border-none">
								<Table.Cell>
									<span class="flex items-center gap-1">
										<div>
											<span class="flex items-center gap-1">
												{facilityByCategory.name}
												{#if facilityByCategory.charmName.includes('kubernetes-control-plane')}
													<a
														href={`/management/scope/${scopeUuid}/facility/${facilityByCategory.name}`}
														target="_blank"
													>
														<Icon icon="ph:arrow-square-out" />
													</a>
												{/if}
											</span>
											{#if facilityByCategory.version}
												<p class="text-xs font-light text-muted-foreground">
													{facilityByCategory.version}
												</p>
											{/if}
										</div>
									</span>
								</Table.Cell>
								<Table.Cell>
									<span class="flex items-center justify-between">
										<div>
											{#each facilityByCategory.units.sort( (previous, present) => previous.name.localeCompare(present.name) ) as unit}
												<span class="flex items-center gap-1">
													<Badge variant="outline">
														{unit.name}
													</Badge>
													{@render ReadUnit(unit)}
												</span>
											{/each}
										</div>
										<DropdownMenu.Root>
											<DropdownMenu.Trigger>
												<Button variant="ghost">
													<Icon icon="ph:dots-three-vertical" />
												</Button>
											</DropdownMenu.Trigger>
											<DropdownMenu.Content>
												<DropdownMenu.Item onSelect={(e) => e.preventDefault()}>
													<AddFacilityUnits
														{scopeUuid}
														{facilityByCategory}
														machines={$machinesStore}
													/>
												</DropdownMenu.Item>
											</DropdownMenu.Content>
										</DropdownMenu.Root>
									</span>
								</Table.Cell>
								<Table.Cell>
									{#if facilityByCategory.status}
										<span class="flex items-center gap-1">
											<Badge variant="outline">
												{facilityByCategory.status.state}
											</Badge>
											{@render ReadStatus(facilityByCategory.status)}
										</span>
									{/if}
								</Table.Cell>
								<Table.Cell>
									{facilityByCategory.revision}
								</Table.Cell>
								<Table.Cell>
									{facilityByCategory.charmName}
								</Table.Cell>
								<Table.Cell>
									<div class="flex justify-end">
										<ManagementFacilityActions {scopeUuid} facilityName={facilityByCategory.name} />
									</div>
								</Table.Cell>

								<Table.Cell>
									<div class="flex justify-end">
										<DropdownMenu.Root>
											<DropdownMenu.Trigger>
												<Icon icon="ph:dots-three-vertical" />
											</DropdownMenu.Trigger>
											<DropdownMenu.Content>
												<DropdownMenu.Item onSelect={(e) => e.preventDefault()}>
													<UpdateFacility {scopeUuid} {facilityByCategory} />
												</DropdownMenu.Item>
												<DropdownMenu.Item onSelect={(e) => e.preventDefault()}>
													<DeleteFacility {scopeUuid} {facilityByCategory} />
												</DropdownMenu.Item>
												<DropdownMenu.Item onSelect={(e) => e.preventDefault()}>
													<ExposeFacility {scopeUuid} {facilityByCategory} />
												</DropdownMenu.Item>
											</DropdownMenu.Content>
										</DropdownMenu.Root>
									</div>
								</Table.Cell>
							</Table.Row>
						{/each}
					{/if}
				{/if}
			{/each}
		</Table.Body>
	</Table.Root>
</main>

{#snippet StatisticFacilities(facilities: Facility[])}
	{@const numberOfApplications = facilities.length}
	{@const numberOfHealthApplications = facilities.filter((a) =>
		a.status ? a.status.state === 'active' : false
	).length}
	{@const health = (numberOfHealthApplications * 100) / numberOfApplications || 0}
	<div class="grid grid-cols-4 gap-3 *:border-none *:shadow-none">
		<Card.Root>
			<Card.Header>
				<Card.Title>Facility</Card.Title>
			</Card.Header>
			<Card.Content class="text-7xl">
				{numberOfApplications}
			</Card.Content>
		</Card.Root>
		<Card.Root>
			<Card.Header>
				<Card.Title>Health</Card.Title>
			</Card.Header>
			<Card.Content>
				<p class="text-3xl">
					{Math.round(health)}%
				</p>
				<p class="text-xs text-muted-foreground">
					{numberOfHealthApplications} active over {numberOfApplications} facilities
				</p>
			</Card.Content>
			<Card.Footer>
				<Progress
					value={health}
					max={100}
					class={`${
						health > 62
							? 'bg-green-50 *:bg-green-700'
							: health > 38
								? 'bg-yellow-50 *:bg-yellow-500'
								: 'bg-red-50 *:bg-red-700'
					}`}
				/>
			</Card.Footer>
		</Card.Root>
	</div>
{/snippet}

{#snippet ReadUnit(unit: Facility_Unit)}
	<HoverCard.Root openDelay={13}>
		<HoverCard.Trigger>
			<Icon icon="ph:info" class="size-4 text-blue-800" />
		</HoverCard.Trigger>
		<HoverCard.Content class="min-w-max">
			<Table.Root>
				<Table.Body>
					<Table.Row>
						<Table.Cell class="text-xs font-light">Name</Table.Cell>
						<Table.Cell class="text-right">{unit.name}</Table.Cell>
					</Table.Row>
					<Table.Row>
						<Table.Cell class="text-xs font-light">Version</Table.Cell>
						<Table.Cell class="text-right">
							{#if unit.version}
								<Badge variant="outline">
									{unit.version}
								</Badge>
							{/if}
						</Table.Cell>
					</Table.Row>
					<Table.Row>
						<Table.Cell class="text-xs font-light">Leader</Table.Cell>
						<Table.Cell class="text-right">
							<div class="flex justify-end">
								<Icon
									icon={unit.leader ? 'ph:circle' : 'ph:x'}
									class={unit.leader ? 'text-blue-500' : 'text-red-500'}
								/>
							</div>
						</Table.Cell>
					</Table.Row>
					<Table.Row>
						<Table.Cell class="text-xs font-light">IP Address</Table.Cell>
						<Table.Cell class="text-right">
							{unit.ipAddress}
						</Table.Cell>
					</Table.Row>
					<Table.Row>
						<Table.Cell class="text-xs font-light">Ports</Table.Cell>
						<Table.Cell class="text-right">
							{#each unit.ports as port}
								<Badge variant="outline">
									{port}
								</Badge>
							{/each}
						</Table.Cell>
					</Table.Row>
					<Table.Row>
						<Table.Cell class="text-xs font-light">Agent</Table.Cell>
						<Table.Cell class="text-right">
							{#if unit.agentStatus}
								<Badge variant="outline">
									{unit.agentStatus.state}
								</Badge>
							{/if}
						</Table.Cell>
					</Table.Row>
					<Table.Row>
						<Table.Cell class="text-xs font-light">Workload</Table.Cell>
						<Table.Cell class="text-right">
							{#if unit.workloadStatus}
								<Badge variant="outline">
									{unit.workloadStatus.state}
								</Badge>
							{/if}
						</Table.Cell>
					</Table.Row>
					<Table.Row>
						<Table.Cell class="text-xs font-light">Machine</Table.Cell>
						<Table.Cell class="flex justify-end">
							{#if unit.workloadStatus?.state == 'active' && unit.machineId}
								<Badge variant="outline">
									{unit.machineId}
								</Badge>
								<!-- <span class="flex items-center gap-1">
									<a href={`/management/machine/${unit.machineSystemId}`}>
										<Icon icon="ph:arrow-square-out" />
									</a>
								</span> -->
							{/if}
						</Table.Cell>
					</Table.Row>
					{#if unit.subordinates.length > 0}
						<Table.Row>
							<Table.Cell class="text-xs font-light">Depending</Table.Cell>
							<Table.Cell class="flex justify-end">
								{#each unit.subordinates as depending}
									<span class="flex items-center gap-1">
										<Badge variant="outline">
											{depending.name}
										</Badge>
										{@render ReadUnit(depending)}
									</span>
								{/each}
							</Table.Cell>
						</Table.Row>
					{/if}
				</Table.Body>
			</Table.Root>
		</HoverCard.Content>
	</HoverCard.Root>
{/snippet}

{#snippet ReadStatus(status: Facility_Status)}
	<HoverCard.Root openDelay={13}>
		<HoverCard.Trigger>
			<Icon icon="ph:info" class="size-4 text-blue-800" />
		</HoverCard.Trigger>
		<HoverCard.Content class="min-w-max">
			<Table.Root>
				<Table.Body>
					<Table.Row>
						<Table.Cell class="text-xs font-light">State</Table.Cell>
						<Table.Cell class="text-right">
							{#if status.state}
								<Badge variant="outline">
									{status.state}
								</Badge>
							{/if}
						</Table.Cell>
					</Table.Row>
					<Table.Row>
						<Table.Cell class="text-xs font-light">Details</Table.Cell>
						<Table.Cell class="text-right">{status.details}</Table.Cell>
					</Table.Row>
					<Table.Row>
						<Table.Cell class="text-xs font-light">Create Time</Table.Cell>
						<Table.Cell class="text-right">
							{#if status.createdAt}
								{status.createdAt.seconds}
							{/if}
						</Table.Cell>
					</Table.Row>
				</Table.Body>
			</Table.Root>
		</HoverCard.Content>
	</HoverCard.Root>
{/snippet}
