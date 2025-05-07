<script lang="ts">
	import * as m from '$lib/paraglide/messages.js';
	import { page } from '$app/state';
	import Icon from '@iconify/svelte';
	import { Badge } from '$lib/components/ui/badge';
	import * as Card from '$lib/components/ui/card/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import * as Table from '$lib/components/ui/table/index.js';
	import { Progress } from '$lib/components/ui/progress/index.js';
	import { formatCapacity } from '$lib/formatter';
	import { cn } from '$lib/utils';
	import PowerOffMachine from './power-off.svelte';
	import CreateMachine from './create.svelte';
	import DeleteMachine from './delete.svelte';
	import { Nexus, type Machine, type Tag } from '$gen/api/nexus/v1/nexus_pb';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	let {
		machines
	}: {
		machines: Machine[];
	} = $props();

	const machineSubvalueContentClass = cn('text-xs font-extralight');

	const transport: Transport = getContext('transportNEW');
	const client = createClient(Nexus, transport);

	const tagStore = writable<Tag[]>();
	async function fetchTags() {
		try {
			const response = await client.listTags({});
			tagStore.set(response.tags);
		} catch (error) {
			console.error('Error fetching tags:', error);
		}
	}

	async function refreshMachines() {
		while (page.url.searchParams.get('intervals')) {
			await new Promise((resolve) =>
				setTimeout(resolve, 1000 * Number(page.url.searchParams.get('intervals')))
			);
			console.log(`Refresh machines`);

			try {
				const response = await client.listMachines({});
				machines = response.machines;
			} catch (error) {
				console.error('Error refreshing:', error);
			}
		}
	}

	onMount(async () => {
		try {
			await fetchTags();
			refreshMachines();
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

<div>
	{@render StatisticMachines()}
	<div class="p-4">
		<Table.Root>
			<Table.Header class="bg-muted/50">
				<Table.Row class="*:text-xs *:font-light [&>th]:py-2 [&>th]:align-top">
					<Table.Head>
						<div>FQDN</div>
						<div>IP</div>
					</Table.Head>
					<Table.Head>POWER</Table.Head>
					<Table.Head>STATUS</Table.Head>
					<Table.Head class="text-right ">
						<div>CORES</div>
						<div>ARCH</div>
					</Table.Head>
					<Table.Head class="text-end  ">RAM</Table.Head>
					<Table.Head>DISKS</Table.Head>
					<Table.Head class="text-end">STORAGE</Table.Head>
				</Table.Row>
			</Table.Header>
			<Table.Body>
				{#each machines.sort( (previous, present) => previous.fqdn.localeCompare(present.fqdn) ) as machine, index}
					<Table.Row class="*:truncate *:whitespace-nowrap [&>td]:align-top">
						<Table.Cell>
							<div class="flex justify-between">
								<span>
									<a href={`/management/machine/${machine.id}?intervals=5`}>
										<div class="flex items-center gap-1">
											<p>{machine.fqdn}</p>
											<Icon icon="ph:arrow-square-out" />
										</div>
									</a>
									<div class={machineSubvalueContentClass}>
										{machine.ipAddresses.join(', ')}
									</div>
								</span>
							</div>
						</Table.Cell>
						<Table.Cell>
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
						</Table.Cell>
						<Table.Cell>
							<Badge variant="outline">
								{machine.status}
							</Badge>
							<p class="text-xs font-light">
								{#if machine.statusMessage != 'Deployed'}
									<span class="flex items-center gap-1">
										{#if machine.status.toLowerCase() === 'commissioning' || machine.status.toLowerCase() === 'deploying' || machine.status.toLowerCase() === 'disk_erasing' || machine.status.toLowerCase() === 'entering_rescue_mode' || machine.status.toLowerCase() === 'exiting_rescue_mode' || machine.status.toLowerCase() === 'releasing' || machine.status.toLowerCase() === 'testing'}
											<Icon icon="ph:spinner" class="animate-spin" />
										{/if}
										{machine.statusMessage}
									</span>
								{:else}
									{`${machine.osystem} ${machine.hweKernel} ${machine.distroSeries}`}
								{/if}
							</p>
						</Table.Cell>
						<Table.Cell>
							<div class="text-right">
								<div>{machine.cpuCount}</div>
								<div class={machineSubvalueContentClass}>
									{machine.architecture}
								</div>
							</div>
						</Table.Cell>
						<Table.Cell>
							<div class="flex items-end justify-end space-x-1">
								<div>{formatCapacity(machine.memoryMb).value}</div>
								<div class="text-xs font-extralight">
									{formatCapacity(machine.memoryMb).unit}
								</div>
							</div>
						</Table.Cell>
						<Table.Cell class="text-center ">{machine.blockDevices.length}</Table.Cell>
						<Table.Cell>
							<div class="flex items-end justify-end space-x-1">
								<div>{formatCapacity(machine.storageMb).value}</div>
								<div class="text-xs font-extralight">
									{formatCapacity(machine.storageMb).unit}
								</div>
							</div>
						</Table.Cell>
						<Table.Cell>
							<DropdownMenu.Root>
								<DropdownMenu.Trigger>
									<Icon icon="ph:dots-three" class="size-6" />
								</DropdownMenu.Trigger>
								<DropdownMenu.Content>
									<DropdownMenu.Item onSelect={(e) => e.preventDefault()}>
										<CreateMachine
											bind:machines
											{machine}
											disabled={!!machine.workloadAnnotations['juju-model-uuid']}
										/>
									</DropdownMenu.Item>
									<DropdownMenu.Item onSelect={(e) => e.preventDefault()}>
										<DeleteMachine
											bind:machines
											{machine}
											disabled={!machine.workloadAnnotations['juju-model-uuid']}
										/>
									</DropdownMenu.Item>
									<DropdownMenu.Separator />
									<DropdownMenu.Item onSelect={(e) => e.preventDefault()}>
										<PowerOffMachine
											bind:machine={machines[index]}
											disabled={machine.powerState.toLowerCase() !== 'on'}
										/>
									</DropdownMenu.Item>
									<DropdownMenu.Separator />
									<DropdownMenu.Sub>
										<DropdownMenu.SubTrigger><Icon icon="ph:tag" /> Tag</DropdownMenu.SubTrigger>
										<DropdownMenu.SubContent>
											{#each $tagStore as tag}
												<DropdownMenu.CheckboxItem
													checked={machine.tags.includes(tag.name)}
													class="capitalize"
												>
													{tag.name}
												</DropdownMenu.CheckboxItem>
											{/each}
										</DropdownMenu.SubContent>
									</DropdownMenu.Sub>
								</DropdownMenu.Content>
							</DropdownMenu.Root>
						</Table.Cell>
					</Table.Row>
				{/each}
			</Table.Body>
		</Table.Root>
	</div>
</div>

{#snippet StatisticMachines()}
	<span class="grid grid-cols-4 gap-4">
		<Card.Root class="h-full">
			<Card.Header class="h-10">
				<Card.Title>MACHINE</Card.Title>
			</Card.Header>
			<Card.Content class="h-30">
				<p class="text-6xl">{machines.length}</p>
				<div class="flex flex-wrap gap-2 pt-2">
					{#each [...new Set(machines.map((m) => m.status))] as status}
						<Badge variant="outline">
							{status}: {machines.filter((m) => m.status === status).length}
						</Badge>
					{/each}
				</div>
			</Card.Content>
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
					<p class="pt-2 text-xs text-muted-foreground">
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
