<script lang="ts">
	import { type Table } from '@tanstack/table-core';

	import { type Machine } from '$lib/api/machine/v1/machine_pb';
	import * as Statistics from '$lib/components/custom/data-table/statistics/index';
	import { formatBigNumber, formatCapacity, formatPercentage } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';

	let {
		table
	}: {
		table: Table<Machine>;
	} = $props();

	const filteredMachines = $derived(table.getFilteredRowModel().rows.map((row) => row.original));
</script>

<div class="grid w-full gap-3 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4 xl:grid-cols-5">
	<Statistics.Root type="count">
		<Statistics.Header>
			<Statistics.Title>{m.node()}</Statistics.Title>
		</Statistics.Header>
		<Statistics.Content>
			<p class="text-8xl font-semibold">{filteredMachines.length}</p>
		</Statistics.Content>
		<Statistics.Background icon="ph:devices" />
	</Statistics.Root>

	<Statistics.Root type="ratio">
		{@const powerOnMachines = filteredMachines.filter((m) => m.powerState === 'on').length}
		{@const totalMachines = filteredMachines.length}
		{@const percentage = formatPercentage(powerOnMachines, totalMachines, 0)}
		<Statistics.Header>
			<Statistics.Title>{m.power()}</Statistics.Title>
		</Statistics.Header>
		<Statistics.Content>
			<p class="text-6xl font-semibold">{percentage ? `${percentage} %` : 'NaN'}</p>
			<p class="text-3xl text-muted-foreground">
				{formatBigNumber(powerOnMachines)}/{formatBigNumber(totalMachines)}
			</p>
		</Statistics.Content>
		<Statistics.Background icon="ph:power" />
		<Statistics.Progress numerator={powerOnMachines} denominator={totalMachines} target="LTB" />
	</Statistics.Root>

	<Statistics.Root type="ratio">
		{@const deployedMachines = filteredMachines.filter(
			(machine) => machine.status === 'Deployed'
		).length}
		{@const totalMachines = filteredMachines.length}
		{@const percentage = formatPercentage(deployedMachines, totalMachines, 0)}
		<Statistics.Header>
			<Statistics.Title>{m.deployments()}</Statistics.Title>
		</Statistics.Header>
		<Statistics.Content>
			<p class="text-6xl font-semibold">{percentage ? `${percentage} %` : 'NaN'}</p>
			<p class="text-3xl text-muted-foreground">
				{formatBigNumber(deployedMachines)}/{formatBigNumber(totalMachines)}
			</p>
		</Statistics.Content>
		<Statistics.Background icon="ph:check" />
		<Statistics.Progress numerator={deployedMachines} denominator={totalMachines} target="LTB" />
	</Statistics.Root>

	<Statistics.Root type="count">
		{@const disks = filteredMachines.reduce((a, machine) => a + machine.blockDevices.length, 0)}
		{@const { value: storageValue, unit: storageUnit } = formatCapacity(
			filteredMachines.reduce((a, machine) => a + machine.storageMb * 1024 ** 2, 0)
		)}
		<Statistics.Header>
			<Statistics.Title>{m.disk()}</Statistics.Title>
		</Statistics.Header>
		<Statistics.Content>
			<p class="text-6xl font-semibold">{disks}</p>
			<span class="flex items-end gap-1">
				<p class="text-4xl font-semibold">{storageValue}</p>
				<p class="text-3xl text-muted-foreground">{storageUnit}</p>
			</span>
		</Statistics.Content>
		<Statistics.Background icon="ph:disc" />
	</Statistics.Root>
</div>
