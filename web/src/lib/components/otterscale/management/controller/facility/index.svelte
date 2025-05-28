<script lang="ts">
	import AddFacilityUnits from './add-units.svelte';
	import UpdateFacility from './update.svelte';
	import DeleteFacility from './delete.svelte';
	import ExposeFacility from './expose.svelte';
	import AddCephUnits from './add-ceph-units.svelte';
	import AddKubernetesUnits from './add-kubernetes-units.svelte';
	import { page } from '$app/state';
	import * as HoverCard from '$lib/components/ui/hover-card/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import * as Card from '$lib/components/ui/card';
	import { Progress } from '$lib/components/ui/progress/index.js';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import { createClient, type Transport } from '@connectrpc/connect';
	import {
		ManagementScopeComboBox,
		ManagementScopeCreate,
		ManagementFacilityActions
	} from '$lib/components/otterscale/index';

	import {
		Nexus,
		type Facility,
		type Facility_Info,
		type Facility_Status,
		type Facility_Unit,
		type Scope
	} from '$gen/api/nexus/v1/nexus_pb';

	import Icon from '@iconify/svelte';

	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as Table from '$lib/components/ui/table/index.js';

	let scopeUuid = $state(page.url.searchParams.get('scope') || '');

	const transport: Transport = getContext('transport');
	const client = createClient(Nexus, transport);

	const kubernetesesStore = writable<Facility_Info[]>([]);
	async function fetchKuberneteses() {
		try {
			const response = await client.listKuberneteses({
				scopeUuid: scopeUuid
			});
			kubernetesesStore.set(response.kuberneteses);
		} catch (error) {
			console.error('Error fetching:', error);
		}
	}

	const scopesStore = writable<Scope[]>([]);
	const scopesLoading = writable(true);
	async function fetchScopes() {
		try {
			const response = await client.listScopes({});
			scopesStore.set(response.scopes);

			let defaultScope = response.scopes.find((s) => s.name === 'default');
			if (defaultScope) {
				scopeUuid = defaultScope.uuid;
			}
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			scopesLoading.set(false);
		}
	}

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

	async function refreshFacilties() {
		while (page.url.searchParams.get('intervals')) {
			await new Promise((resolve) =>
				setTimeout(resolve, 1000 * Number(page.url.searchParams.get('intervals')))
			);
			console.log(`Refresh facilities`);

			try {
				const response = await client.listFacilities({
					scopeUuid: scopeUuid
				});
				facilitiesStore.set(response.facilities);
			} catch (error) {
				console.error('Error fetching:', error);
			}
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
			case charmName.includes('kubernetes'):
			case charmName.includes('kubeapi-load-balancer'):
			case charmName.includes('etcd'):
			case charmName.includes('keepalived'):
				return FACILITY_CATEGORYS.KUBERNETES;
			case charmName.includes('ceph'):
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

	let mounted = $state(false);
	onMount(async () => {
		try {
			await fetchScopes();
			await fetchKuberneteses();
			await fetchFacilities();
			refreshFacilties();
		} catch (error) {
			console.error('Error during initial data load:', error);
		}

		mounted = true;
	});

	async function handleChange() {
		await fetchFacilities();
		await fetchKuberneteses();
	}
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

{#snippet AddBundleUnits(facilityCategory: string, facilities: Facility[])}
	{#if facilityCategory == FACILITY_CATEGORYS.KUBERNETES}
		{@const masterKubernetes = facilities.find((f) =>
			f.charmName.includes('kubernetes-control-plane')
		)}
		{#if masterKubernetes}
			<AddKubernetesUnits {scopeUuid} kubernetes={masterKubernetes} />
		{/if}
	{:else if facilityCategory == FACILITY_CATEGORYS.CEPH}
		{@const masterCeph = facilities.find((f) => f.charmName.includes('ceph-mon'))}
		{#if masterCeph}
			<AddCephUnits {scopeUuid} ceph={masterCeph} />
		{/if}
	{/if}
{/snippet}

<main>
	{@render StatisticFacilities($facilitiesStore)}

	<div class="flex justify-end space-x-2 py-4">
		{#if mounted}
			{@const selected = $kubernetesesStore.find((k) => k.scopeUuid === scopeUuid)}
			{#if selected}
				<Button href="/management/facility/{selected.facilityName}?scope={selected.scopeUuid}">
					<Icon icon="ph:rocket-launch" />
					Applications
				</Button>
			{/if}
		{/if}
		<ManagementScopeCreate label={true} scopes={$scopesStore} />
		<ManagementScopeComboBox
			scopes={$scopesStore}
			bind:selected={scopeUuid}
			onSelect={handleChange}
		/>
	</div>

	<div class="p-4">
		<Table.Root>
			<Table.Header class="bg-muted/50">
				<Table.Row class="*:text-xs *:font-light [&>th]:py-2 [&>th]:align-top">
					<Table.Head>
						NAME
						<p>VERSION</p>
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
							<Table.Cell colspan={8} class="border-b hover:bg-transparent">
								<span class="flex items-center justify-between">
									<span class="flex items-center gap-1">
										<Icon icon={getLogoByCategory(facilityCategory)} class="h-[23px] w-fit" />
										<h1 class="text-start text-[23px]">{facilityCategory}</h1>
									</span>
									<span class="flex items-center gap-2">
										{@render collapsibilityHandler(facilityCategory)}
									</span>
								</span>
							</Table.Cell>
						</Table.Row>

						{#if collapsibleOpen[facilityCategory]}
							{@const facilitiesByCategory = $facilitiesStore
								.filter((a) => getCategory(a) === facilityCategory)
								.sort((previous, present) => previous.name.localeCompare(present.name))}
							<Table.Row class="border-none hover:bg-transparent">
								<Table.Cell colspan={8}>
									<span class="flex justify-end">
										{@render AddBundleUnits(facilityCategory, facilitiesByCategory)}
									</span>
								</Table.Cell>
							</Table.Row>
							{#each facilitiesByCategory as facilityByCategory, index}
								<Table.Row class="border-none">
									<Table.Cell>
										<span class="flex items-center gap-1">
											<div>
												<span class="flex items-center gap-1">
													{facilityByCategory.name}
													{#if facilityByCategory.charmName.includes('kubernetes-control-plane')}
														<a
															href={`/management/facility/${facilityByCategory.name}?scope=${scopeUuid}`}
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
															bind:facilityByCategory={facilitiesByCategory[index]}
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
											<ManagementFacilityActions
												{scopeUuid}
												facilityName={facilityByCategory.name}
											/>
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
	</div>
</main>

{#snippet StatisticFacilities(facilities: Facility[])}
	{@const numberOfApplications = facilities.length}
	{@const numberOfHealthApplications = facilities.filter((a) =>
		a.status ? a.status.state === 'active' : false
	).length}
	{@const health = (numberOfHealthApplications * 100) / numberOfApplications || 0}
	<div class="grid grid-cols-4 gap-4">
		<Card.Root>
			<Card.Header>
				<Card.Title>FACILITY</Card.Title>
			</Card.Header>
			<Card.Content class="text-7xl">
				{numberOfApplications}
			</Card.Content>
		</Card.Root>
		<Card.Root>
			<Card.Header>
				<Card.Title>HEALTH</Card.Title>
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
				<Table.Body class="*:text-xs [&>tr>th]:text-right [&>tr>th]:font-light">
					<Table.Row>
						<Table.Head>Name</Table.Head>
						<Table.Cell>{unit.name}</Table.Cell>
					</Table.Row>
					<Table.Row>
						<Table.Head>Version</Table.Head>
						<Table.Cell>
							{#if unit.version}
								<Badge variant="outline">
									{unit.version}
								</Badge>
							{/if}
						</Table.Cell>
					</Table.Row>
					<Table.Row>
						<Table.Head>Leader</Table.Head>
						<Table.Cell>
							<Icon
								icon={unit.leader ? 'ph:circle' : 'ph:x'}
								class={unit.leader ? 'text-blue-500' : 'text-red-500'}
							/>
						</Table.Cell>
					</Table.Row>
					<Table.Row>
						<Table.Head>IP Address</Table.Head>
						<Table.Cell>
							{unit.ipAddress}
						</Table.Cell>
					</Table.Row>
					<Table.Row>
						<Table.Head>Ports</Table.Head>
						<Table.Cell>
							{#each unit.ports as port}
								<Badge variant="outline">
									{port}
								</Badge>
							{/each}
						</Table.Cell>
					</Table.Row>
					<Table.Row>
						<Table.Head>Agent</Table.Head>
						<Table.Cell>
							{#if unit.agentStatus}
								<Badge variant="outline">
									{unit.agentStatus.state}
								</Badge>
							{/if}
						</Table.Cell>
					</Table.Row>
					<Table.Row>
						<Table.Head>Workload</Table.Head>
						<Table.Cell>
							{#if unit.workloadStatus}
								<Badge variant="outline">
									{unit.workloadStatus.state}
								</Badge>
							{/if}
						</Table.Cell>
					</Table.Row>
					<Table.Row>
						<Table.Head>Machine</Table.Head>
						<Table.Cell>
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
							<Table.Head>Depending</Table.Head>
							<Table.Cell>
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
				<Table.Body class="*:text-xs [&>tr>th]:text-right [&>tr>th]:font-light">
					<Table.Row>
						<Table.Head>State</Table.Head>
						<Table.Cell>
							{#if status.state}
								<Badge variant="outline">
									{status.state}
								</Badge>
							{/if}
						</Table.Cell>
					</Table.Row>
					<Table.Row>
						<Table.Head>Details</Table.Head>
						<Table.Cell>{status.details}</Table.Cell>
					</Table.Row>
					<Table.Row>
						<Table.Head>Create Time</Table.Head>
						<Table.Cell>
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
