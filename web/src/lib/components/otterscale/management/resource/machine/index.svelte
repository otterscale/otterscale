<script lang="ts">
	import Icon from '@iconify/svelte';
	import { writable } from 'svelte/store';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import * as Table from '$lib/components/ui/table/index.js';
	import { Progress } from '$lib/components/ui/progress/index.js';
	import { formatCapacity } from '$lib/formatter';
	import { cn } from '$lib/utils';

	import PowerOnMachine from './power-on.svelte';
	import PowerOffMachine from './power-off.svelte';
	import CreateMachine from './create.svelte';
	import DeleteMachine from './delete.svelte';
	import RemoveTags from './remove-tags.svelte';
	import AddTags from './add-tags.svelte';

	import { Nexus, type Machine, type Scope, type Tag } from '$gen/api/nexus/v1/nexus_pb';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onMount } from 'svelte';
	let {
		machines
	}: {
		machines: Machine[];
	} = $props();

	const machineSubvalueContentClass = cn('text-xs font-extralight');

	const transport: Transport = getContext('transportNEW');
	const client = createClient(Nexus, transport);

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
			<Card.Header class="h-10">
				<Card.Title>MACHINE</Card.Title>
			</Card.Header>
			<Card.Content class="h-30">
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
			<Card.Header class="h-10">
				<Card.Title>STORAGE</Card.Title>
			</Card.Header>
			<Card.Content class="h-30">
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
			<Card.Header class="h-10">
				<Card.Title>POWER ON</Card.Title>
			</Card.Header>
			<Card.Content class="h-30">
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
			<Card.Header class="h-10">
				<Card.Title>DEPLOYMENT</Card.Title>
			</Card.Header>
			<Card.Content class="h-30">
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
