<script lang="ts">
	import { type Machine } from '$lib/api/machine/v1/machine_pb';
	import { Badge } from '$lib/components/ui/badge';
	import * as Card from '$lib/components/ui/card';
	import { Progress } from '$lib/components/ui/progress/index.js';
	import { formatCapacity } from '$lib/formatter';

	let { machines }: { machines: Machine[] } = $props();

	// Calculate statistics
	const totalMachines = $derived(machines.length);
	const uniqueStatuses = $derived([...new Set(machines.map((m) => m.status))]);
	const totalDisks = $derived(
		machines.reduce((acc, machine) => acc + machine.blockDevices.length, 0)
	);
	const storageFormatted = $derived(
		formatCapacity(machines.reduce((acc, machine) => acc + machine.storageMb, 0) * 1024 * 1024)
	);
	const machinesOn = $derived(machines.filter((m) => m.powerState === 'on').length);
	const machinesDeployed = $derived(machines.filter((m) => m.status === 'Deployed').length);
	const powerOnPercentage = $derived(Math.round((machinesOn / totalMachines) * 100));
	const deploymentPercentage = $derived(Math.round((machinesDeployed / totalMachines) * 100));
</script>

<span class="grid grid-cols-4 gap-4">
	<Card.Root class="h-full">
		<Card.Header class="h-10">
			<Card.Title>MACHINE</Card.Title>
		</Card.Header>
		<Card.Content class="h-30">
			<p class="text-6xl">{totalMachines}</p>
			<div class="flex flex-wrap gap-2 pt-2">
				{#each uniqueStatuses as status}
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
				<span>{storageFormatted.value}</span>
				<span class="text-3xl font-extralight">
					{storageFormatted.unit}
				</span>
				<p class="text-muted-foreground pt-2 text-xs">
					over {totalDisks} disks
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
				{powerOnPercentage}%
			</p>
			<p class="text-muted-foreground text-xs">
				{machinesOn} On over {totalMachines} units
			</p>
		</Card.Content>
		<Card.Footer>
			<Progress value={powerOnPercentage} max={100} />
		</Card.Footer>
	</Card.Root>
	<Card.Root>
		<Card.Header class="h-10">
			<Card.Title>DEPLOYMENT</Card.Title>
		</Card.Header>
		<Card.Content class="h-30">
			<p class="text-3xl">
				{deploymentPercentage}%
			</p>

			<p class="text-muted-foreground text-xs">
				{machinesDeployed} deployed over {totalMachines}
				units
			</p>
		</Card.Content>
		<Card.Footer>
			<Progress value={deploymentPercentage} max={100} />
		</Card.Footer>
	</Card.Root>
</span>
