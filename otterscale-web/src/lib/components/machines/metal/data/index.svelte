<script lang="ts" module>
	import { type Machine } from '$lib/api/machine/v1/machine_pb';
	import { Layout } from '$lib/components/custom/instance';
	import { Badge } from '$lib/components/ui/badge';
	import Icon from '@iconify/svelte';
	import { type Writable } from 'svelte/store';
	import Alert from './alert.svelte';
	import StatisticCPU from './statistic-cpu.svelte';
	import StatisticMemory from './statistic-memory.svelte';
	import StatisticPower from './statistic-power.svelte';
	import StatisticStorage from './statistic-storage.svelte';
	import TableBlockDevice from './table-block-device.svelte';
	import TableChassis from './table-chassis.svelte';
	import TableMainboard from './table-mainboard.svelte';
	import TableNetwork from './table-network.svelte';
	import TableSystem from './table-system.svelte';
</script>

<script lang="ts">
	let {
		machine
	}: {
		machine: Writable<Machine>;
	} = $props();
</script>

<main class="space-y-4 py-4">
	{#if $machine.statusMessage !== 'Deployed'}
		<Alert {machine} />
	{/if}

	<div class="space-y-4 py-4">
		<div class="flex items-end gap-2 text-5xl">
			<p class="text-muted-foreground">{$machine.id}</p>
			{$machine.fqdn}
		</div>
		<div class="flex flex-wrap gap-1 overflow-visible">
			{#each $machine.tags as tag}
				<Badge variant="outline">
					<Icon icon="ph:tag" />
					{tag}
				</Badge>
			{/each}
		</div>
	</div>

	<Layout.Statistics>
		<StatisticPower {machine} />
		<StatisticCPU {machine} />
		<StatisticMemory {machine} />
		<StatisticStorage {machine} />
	</Layout.Statistics>

	<Layout.Tables>
		<Layout.Table.Root open={true}>
			<Layout.Table.Trigger>
				<Icon icon="ph:desktop" />
				System
			</Layout.Table.Trigger>
			<Layout.Table.Content>
				<TableSystem {machine} />
			</Layout.Table.Content>
		</Layout.Table.Root>

		<Layout.Table.Root open={true}>
			<Layout.Table.Trigger>
				<Icon icon="ph:circuitry" />
				Mainboard
			</Layout.Table.Trigger>
			<Layout.Table.Content>
				<TableMainboard {machine} />
			</Layout.Table.Content>
		</Layout.Table.Root>

		<Layout.Table.Root open={true}>
			<Layout.Table.Trigger>
				<Icon icon="ph:computer-tower" />
				Chassis
			</Layout.Table.Trigger>
			<Layout.Table.Content>
				<TableChassis {machine} />
			</Layout.Table.Content>
		</Layout.Table.Root>

		<Layout.Table.Root>
			<Layout.Table.Trigger>
				<Icon icon="ph:hard-drives" />
				Block Device
			</Layout.Table.Trigger>
			<Layout.Table.Content>
				<TableBlockDevice {machine} />
			</Layout.Table.Content>
		</Layout.Table.Root>

		<Layout.Table.Root>
			<Layout.Table.Trigger>
				<Icon icon="ph:network" />
				Network
			</Layout.Table.Trigger>
			<Layout.Table.Content>
				<TableNetwork {machine} />
			</Layout.Table.Content>
		</Layout.Table.Root>
	</Layout.Tables>
</main>
