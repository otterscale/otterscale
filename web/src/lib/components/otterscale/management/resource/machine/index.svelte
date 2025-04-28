<script lang="ts">
	import Icon from '@iconify/svelte';
	import * as Select from '$lib/components/ui/select/index.js';
	import { Progress } from '$lib/components/ui/progress/index.js';
	import * as HoverCard from '$lib/components/ui/hover-card/index.js';
	import * as Sheet from '$lib/components/ui/sheet';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import * as AlertDialog from '$lib/components/ui/alert-dialog/index.js';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Drawer from '$lib/components/ui/drawer/index.js';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import * as Popover from '$lib/components/ui/popover/index.js';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import * as Table from '$lib/components/ui/table/index.js';
	import * as Tabs from '$lib/components/ui/tabs';

	import { writable } from 'svelte/store';

	import { formatCapacity, formatBigNumber as formatNumber } from '$lib/formatter';
	import { cn } from '$lib/utils';

	import { toast } from 'svelte-sonner';

	import { ManagementNetworkSubnetReservedIPRanges } from '$lib/components/otterscale/index';

	import PowerOnMachine from './power-on.svelte';
	import PowerOffMachine from './power-off.svelte';
	import CreateMachine from './create.svelte';
	import DeleteMachine from './delete.svelte';
	import RemoveTags from './remove-tags.svelte';
	import AddTags from './add-tags.svelte';

	import {
		Nexus,
		type CreateMachineRequest,
		type Scope,
		type CreateIPRangeRequest,
		type CreateNetworkRequest,
		type DeleteIPRangeRequest,
		type DeleteNetworkRequest,
		type Machine,
		type Machine_Placement,
		type Network,
		type Network_Fabric,
		type Network_IPAddress,
		type Network_IPRange,
		type Network_Subnet,
		type Network_VLAN,
		type PowerOffMachineRequest,
		type PowerOnMachineRequest,
		type UpdateFabricRequest,
		type UpdateIPRangeRequest,
		type UpdateSubnetRequest,
		type UpdateVLANRequest,
		type Tag,
		type AddMachineTagsRequest,
		type RemoveMachineTagsRequest
	} from '$gen/api/nexus/v1/nexus_pb';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onMount } from 'svelte';

	let x = {} as Machine_Placement;
	x.type;

	let {
		machines
	}: {
		machines: Machine[];
	} = $props();

	const transport: Transport = getContext('transportNEW');
	const client = createClient(Nexus, transport);

	let powerOnMachineRequest = $state({} as PowerOnMachineRequest);
	let powerOffMachineRequest = $state({} as PowerOffMachineRequest);
	function resetPowerMachineRequest() {
		powerOnMachineRequest = {} as PowerOnMachineRequest;
		powerOffMachineRequest = {} as PowerOffMachineRequest;
	}
	let powerMachineConfirm = $state(false);
	let powerMachineConfirms = $state<Record<string, boolean>>({});
	machines.forEach((m) => {
		powerMachineConfirms[m.id] = false;
	});

	let createMachineRequest = $state({
		enableSsh: true,
		skipBmcConfig: true,
		skipNetworking: true,
		skipStorage: true
	} as CreateMachineRequest);
	function resetCreateMachineRequest() {
		createMachineRequest = {
			enableSsh: true,
			skipBmcConfig: true,
			skipNetworking: true,
			skipStorage: true,
			tags: [] as string[]
		} as CreateMachineRequest;
		createMachineScope = {} as Scope;
	}
	let createMachineScope = $state({} as Scope);
	let createMachineConfirms = $state<Record<string, boolean>>({});
	machines.forEach((m) => {
		createMachineConfirms[m.id] = false;
	});

	let addMachineTagsRequest = $state({ tags: [] as string[] } as AddMachineTagsRequest);
	function resetAddMachineTagsRequest() {
		addMachineTagsRequest = { tags: [] as string[] } as AddMachineTagsRequest;
	}
	let addMachineTagsConfirms = $state<Record<string, boolean>>({});

	machines.forEach((m) => {
		addMachineTagsConfirms[m.id] = false;
	});

	let removeMachineTagsRequest = $state({ tags: [] as string[] } as RemoveMachineTagsRequest);
	function resetRemoveMachineTagsRequest() {
		removeMachineTagsRequest = { tags: [] as string[] } as RemoveMachineTagsRequest;
	}
	let removeMachineTagsConfirms = $state<Record<string, boolean>>({});
	machines.forEach((m) => {
		removeMachineTagsConfirms[m.id] = false;
	});

	// let addMachinesRequest = $state({} as AddMachinesRequest);
	// function resetAddMachinesRequest() {
	// 	addMachinesRequest = {} as AddMachinesRequest;
	// 	resetAddMachinesPlacements();
	// 	resetAddMachinesConstraint();
	// }
	// let addMachinesPlacements = $state([] as string[]);
	// function resetAddMachinesPlacements() {
	// 	addMachinesPlacements = [] as string[];
	// }
	// let addMachinesConstraint = $state({} as Machine_Constraint);
	// function resetAddMachinesConstraint() {
	// 	addMachinesConstraint = {} as Machine_Constraint;
	// }
	// let addMachinesConfirm = $state(false);

	const machineSubvalueContentClass = cn('text-xs font-extralight');

	const scopesStore = writable<Scope[]>([]);
	const scopesLoading = writable(true);
	async function fetchScopes() {
		try {
			const response = await client.listScopes({});
			scopesStore.set(response.scopes);
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			scopesLoading.set(false);
		}
	}

	const tagsStore = writable<Tag[]>();
	const tagsLoading = writable(true);
	async function fetchTags() {
		try {
			const response = await client.listTags({});
			tagsStore.set(response.tags);
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			tagsLoading.set(false);
		}
	}

	let mounted = false;
	onMount(async () => {
		try {
			await fetchScopes();
			await fetchTags();
		} catch (error) {
			console.error('Error during initial data load:', error);
		}

		mounted = true;
	});
</script>

<div class="space-y-3">
	{@render StatisticMachines()}
	<Table.Root>
		<Table.Header>
			<Table.Row class="*:text-xs *:font-light">
				<Table.Head class="align-top">
					<div>FQDN</div>
					<div>IP</div>
				</Table.Head>
				<Table.Head class="align-top">POWER</Table.Head>
				<Table.Head class="align-top">TAGS</Table.Head>
				<Table.Head class="align-top">STATUS</Table.Head>
				<Table.Head class="text-right align-top">
					<div>CORES</div>
					<div>ARCH</div>
				</Table.Head>
				<Table.Head class="align-top">RAM</Table.Head>
				<Table.Head class="align-top">DISKS</Table.Head>
				<Table.Head class="text-end align-top">STORAGE</Table.Head>
			</Table.Row>
		</Table.Header>
		<Table.Body>
			{#each machines as machine}
				<Table.Row>
					<Table.Cell class="align-top">
						<div class="flex items-center justify-between">
							<div class="flex justify-between">
								<span>
									<a href={`/management/machine/${machine.id}`}>
										<div class="flex items-center gap-1">
											{machine.fqdn}
											<Icon icon="ph:arrow-square-out" />
										</div>
									</a>
									<div class={machineSubvalueContentClass}>
										{machine.ipAddresses.join(', ')}
									</div>
								</span>
							</div>
							<DropdownMenu.Root>
								<DropdownMenu.Trigger>
									<Button variant="ghost">
										<Icon icon="ph:dots-three-vertical" />
									</Button>
								</DropdownMenu.Trigger>
								<DropdownMenu.Content>
									<DropdownMenu.Item onSelect={(e) => e.preventDefault()}>
										<CreateMachine {machine} />
									</DropdownMenu.Item>
									<DropdownMenu.Item onSelect={(e) => e.preventDefault()}>
										<DeleteMachine {machine} />
									</DropdownMenu.Item>
								</DropdownMenu.Content>
							</DropdownMenu.Root>
						</div>
					</Table.Cell>
					<Table.Cell class="align-top">
						<div class="flex items-center justify-between">
							<div class="flex items-center gap-1">
								<Icon
									icon={machine.powerState === 'on' ? 'ph:power' : 'ph:power'}
									class={machine.powerState === 'on' ? 'text-green-700' : 'text-red-700'}
								/>
								<div class="flex flex-col items-start">
									<div>{machine.powerState}</div>
									<div class={machineSubvalueContentClass}>{machine.powerType}</div>
								</div>
							</div>
							<DropdownMenu.Root>
								<DropdownMenu.Trigger>
									<Button variant="ghost">
										<Icon icon="ph:dots-three-vertical" />
									</Button>
								</DropdownMenu.Trigger>
								<DropdownMenu.Content>
									<DropdownMenu.Item onSelect={(e) => e.preventDefault()}>
										{#if machine.powerState.toLowerCase() === 'on'}
											<PowerOffMachine {machine} />
										{:else}
											<PowerOnMachine {machine} />
										{/if}
									</DropdownMenu.Item>
								</DropdownMenu.Content>
							</DropdownMenu.Root>
						</div>
					</Table.Cell>
					<Table.Cell class="align-top">
						<div class="flex items-center justify-between">
							<div class="flex gap-1">
								{#each machine.tags as tag}
									<Badge variant="outline">
										{tag}
									</Badge>
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
										<AddTags {machine} tags={$tagsStore} />
									</DropdownMenu.Item>
									<DropdownMenu.Item onSelect={(e) => e.preventDefault()}>
										<RemoveTags {machine} />
									</DropdownMenu.Item>
								</DropdownMenu.Content>
							</DropdownMenu.Root>
						</div>
					</Table.Cell>
					<Table.Cell class="align-top">
						<Badge variant="outline">
							{machine.status}
						</Badge>
						<p class="text-xs font-light">
							{`${machine.osystem} ${machine.hweKernel} ${machine.distroSeries}`}
						</p>
					</Table.Cell>
					<Table.Cell class="align-top">
						<div class="text-right">
							<div>{machine.cpuCount}</div>
							<div class={machineSubvalueContentClass}>
								{machine.architecture}
							</div>
						</div>
					</Table.Cell>
					<Table.Cell class="align-top">
						<div class="flex items-end justify-end space-x-1">
							<div>{formatCapacity(machine.memoryMb).value}</div>
							<div class="text-xs font-extralight">
								{formatCapacity(machine.memoryMb).unit}
							</div>
						</div>
					</Table.Cell>
					<Table.Cell class="text-center align-top">{machine.blockDevices.length}</Table.Cell>
					<Table.Cell class="align-top">
						<div class="flex items-end justify-end space-x-1">
							<div>{formatCapacity(machine.storageMb).value}</div>
							<div class="text-xs font-extralight">
								{formatCapacity(machine.storageMb).unit}
							</div>
						</div>
					</Table.Cell>
				</Table.Row>
			{/each}
		</Table.Body>
	</Table.Root>
</div>

{#snippet StatisticMachines()}
	<span class="grid grid-cols-4 gap-3 *:border-none *:shadow-none">
		<Card.Root class="h-full">
			<Card.Header>
				<Card.Title>Machine</Card.Title>
			</Card.Header>
			<Card.Content>
				<p class="text-7xl">{machines.length}</p>
			</Card.Content>
			<Card.Footer>
				<div class="flex flex-wrap gap-1">
					{#each [...new Set(machines.map((m) => m.status))] as status}
						<Badge variant="outline"
							>{status}: {machines.filter((m) => m.status === status).length}</Badge
						>
					{/each}
				</div>
			</Card.Footer>
		</Card.Root>
		<Card.Root>
			<Card.Header>
				<Card.Title>Storage</Card.Title>
			</Card.Header>
			<Card.Content>
				<div class="text-6xl">
					<span
						>{formatCapacity(machines.reduce((acc, machine) => acc + machine.storageMb, 0))
							.value}</span
					>
					<span class="text-3xl font-extralight">
						{formatCapacity(machines.reduce((acc, machine) => acc + machine.storageMb, 0)).unit}
					</span>
					<p class="text-xs text-muted-foreground">
						over {machines.reduce((acc, machine) => acc + machine.blockDevices.length, 0)} disks
					</p>
				</div>
			</Card.Content>
		</Card.Root>
		<Card.Root>
			<Card.Header>
				<Card.Title>Power On</Card.Title>
			</Card.Header>
			<Card.Content>
				<p class="text-3xl">
					{Math.round(
						(machines.filter((m) => m.powerState === 'on').length / machines.length) * 100
					)}%
				</p>
				<p class="text-xs text-muted-foreground">
					{machines.filter((m) => m.powerState === 'on').length} On over {machines.length} units
				</p>
			</Card.Content>
			<Card.Footer>
				<Progress
					value={machines.filter((m) => m.powerState === 'on').length / machines.length}
					max={1}
				/>
			</Card.Footer>
		</Card.Root>
		<Card.Root>
			<Card.Header>
				<Card.Title>Deployment</Card.Title>
			</Card.Header>
			<Card.Content>
				<p class="text-3xl">
					{Math.round(
						(machines.filter((m) => m.status === 'Deployed').length / machines.length) * 100
					)}%
				</p>

				<p class="text-xs text-muted-foreground">
					{machines.filter((m) => m.status === 'Deployed').length} deployed over {machines.length}
					units
				</p>
			</Card.Content>
			<Card.Footer>
				<Progress
					value={machines.filter((m) => m.status === 'Deployed').length / machines.length}
					max={1}
				/>
			</Card.Footer>
		</Card.Root>
	</span>
{/snippet}
