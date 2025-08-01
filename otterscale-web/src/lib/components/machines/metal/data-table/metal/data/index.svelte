<script lang="ts">
	import { type Machine } from '$lib/api/machine/v1/machine_pb';
	import { Button } from '$lib/components/ui/button';
	import Icon from '@iconify/svelte';
	import { type Writable } from 'svelte/store';
	import * as Layout from '../layout';
	import Alert from './alert.svelte';
	import StatisticCPU from './statistic-cpu.svelte';
	import StatisticMachine from './statistic-machine.svelte';
	import StatisticMemory from './statistic-memory.svelte';
	import StatisticPower from './statistic-power.svelte';
	import StatisticStorage from './statistic-storage.svelte';
	import TableBlockDevice from './table-block-device.svelte';
	import TableChassis from './table-chassis.svelte';
	import TableMainboard from './table-mainboard.svelte';
	import TableNetwork from './table-network.svelte';
	import TableSystem from './table-system.svelte';

	let {
		machine
	}: {
		machine: Writable<Machine>;
	} = $props();
</script>

<main class="space-y-4">
	{#if $machine.statusMessage !== 'Deployed'}
		<Alert {machine} />
	{/if}

	<Layout.Statistics>
		<StatisticMachine {machine} />
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
