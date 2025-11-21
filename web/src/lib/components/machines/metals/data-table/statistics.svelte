<script lang="ts">
	import Icon from '@iconify/svelte';
	import { type Table } from '@tanstack/table-core';

	import { type Machine } from '$lib/api/machine/v1/machine_pb';
	import { formatProgressColor } from '$lib/components/custom/progress/utils.svelte';
	import * as Card from '$lib/components/ui/card';
	import Progress from '$lib/components/ui/progress/progress.svelte';
	import { formatBigNumber, formatCapacity, formatPercentage } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';
	import { cn } from '$lib/utils';

	let {
		table
	}: {
		table: Table<Machine>;
	} = $props();

	const filteredMachines = $derived(table.getFilteredRowModel().rows.map((row) => row.original));
</script>

<div class="grid w-full gap-3 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4 xl:grid-cols-5">
	{#snippet Nodes()}
		{@const title = m.node()}
		{@const titleIcon = 'ph:chart-bar-bold'}
		{@const backgroundIcon = 'ph:devices'}
		{@const nodes = filteredMachines.length}
		<Card.Root class="relative overflow-hidden">
			<Card.Header class="gap-3">
				<Card.Title class="flex items-center gap-2 font-medium">
					<div
						class="flex size-8 shrink-0 items-center justify-center rounded-md bg-primary/10 text-primary"
					>
						<Icon icon={titleIcon} class="size-5" />
					</div>
					<p class="font-bold">{title}</p>
				</Card.Title>
			</Card.Header>
			<Card.Content class="lg:max-[1100px]:flex-col lg:max-[1100px]:items-start">
				<p class="text-8xl font-semibold">{nodes}</p>
			</Card.Content>
			<div
				class="absolute -top-2 -right-16 text-8xl tracking-tight text-nowrap text-primary/5 uppercase group-hover:hidden"
			>
				<Icon icon={backgroundIcon} class="size-68" />
			</div>
		</Card.Root>
	{/snippet}
	{@render Nodes()}

	{#snippet PowerOnMachines()}
		{@const title = m.power()}
		{@const titleIcon = 'ph:chart-pie-bold'}
		{@const backgroundIcon = 'ph:power'}
		{@const powerOnMachines = filteredMachines.filter((m) => m.powerState === 'on').length}
		{@const totalMachines = filteredMachines.length}
		{@const percentage = formatPercentage(powerOnMachines, totalMachines)}
		<Card.Root class="relative overflow-hidden">
			<Card.Header class="gap-3">
				<Card.Title>
					<div class="flex items-center gap-2 font-medium">
						<div
							class="flex size-8 shrink-0 items-center justify-center rounded-md bg-primary/10 text-primary"
						>
							<Icon icon={titleIcon} class="size-5" />
						</div>
						<p class="font-bold">{title}</p>
					</div>
				</Card.Title>
			</Card.Header>
			<Card.Content class="lg:max-[1100px]:flex-col lg:max-[1100px]:items-start">
				<div class="space-y-1">
					<p class="text-6xl font-semibold">{percentage ? `${percentage} %` : 'NaN'}</p>
					<p class="text-3xl text-muted-foreground">
						{formatBigNumber(powerOnMachines)}/{formatBigNumber(totalMachines)}
					</p>
				</div>
			</Card.Content>
			<Progress
				value={Number(percentage ?? 0)}
				max={100}
				class={cn(
					formatProgressColor(powerOnMachines, totalMachines, 'LTB'),
					'absolute top-0 left-0 h-2 rounded-none'
				)}
			/>
			<div
				class="absolute top-0 -right-16 text-8xl tracking-tight text-nowrap text-primary/5 uppercase group-hover:hidden"
			>
				<Icon icon={backgroundIcon} class="size-72" />
			</div>
		</Card.Root>
	{/snippet}
	{@render PowerOnMachines()}

	{#snippet Deployments()}
		{@const title = m.deployments()}
		{@const titleIcon = 'ph:chart-bar-bold'}
		{@const backgroundIcon = 'ph:check'}
		{@const deployedMachines = filteredMachines.filter(
			(machine) => machine.status === 'Deployed'
		).length}
		{@const totalMachines = filteredMachines.length}
		{@const percentage = formatPercentage(deployedMachines, totalMachines)}
		<Card.Root class="relative overflow-hidden">
			<Card.Header class="gap-3">
				<Card.Title>
					<div class="flex items-center gap-2 font-medium">
						<div
							class="flex size-8 shrink-0 items-center justify-center rounded-md bg-primary/10 text-primary"
						>
							<Icon icon={titleIcon} class="size-5" />
						</div>
						<p class="font-bold">{title}</p>
					</div>
				</Card.Title>
			</Card.Header>
			<Card.Content class="lg:max-[1100px]:flex-col lg:max-[1100px]:items-start">
				<div class="space-y-1">
					<p class="text-6xl font-semibold">{percentage ? `${percentage} %` : 'NaN'}</p>
					<p class="text-3xl text-muted-foreground">
						{formatBigNumber(deployedMachines)}/{formatBigNumber(totalMachines)}
					</p>
				</div>
			</Card.Content>
			<Progress
				value={Number(percentage ?? 0)}
				max={100}
				class={cn(
					formatProgressColor(deployedMachines, totalMachines, 'LTB'),
					'absolute top-0 left-0 h-2 rounded-none'
				)}
			/>
			<div
				class="absolute top-0 -right-16 text-8xl tracking-tight text-nowrap text-primary/5 uppercase group-hover:hidden"
			>
				<Icon icon={backgroundIcon} class="size-72" />
			</div>
		</Card.Root>
	{/snippet}
	{@render Deployments()}

	{#snippet Disks()}
		{@const title = m.disks()}
		{@const titleIcon = 'ph:chart-bar-bold'}
		{@const backgroundIcon = 'ph:disc'}
		{@const disks = filteredMachines.reduce((a, machine) => a + machine.blockDevices.length, 0)}
		{@const { value: storageValue, unit: storageUnit } = formatCapacity(
			filteredMachines.reduce((a, machine) => a + machine.storageMb * 1024 ** 2, 0)
		)}
		<Card.Root class="relative overflow-hidden">
			<Card.Header class="gap-3">
				<Card.Title class="flex items-center gap-2 font-medium">
					<div
						class="flex size-8 shrink-0 items-center justify-center rounded-md bg-primary/10 text-primary"
					>
						<Icon icon={titleIcon} class="size-5" />
					</div>
					<p class="font-bold">{title}</p>
				</Card.Title>
			</Card.Header>
			<Card.Content class="lg:max-[1100px]:flex-col lg:max-[1100px]:items-start">
				<p class="text-6xl font-semibold">{disks}</p>
				<span class="flex items-end gap-1">
					<p class="text-4xl font-semibold">{storageValue}</p>
					<p class="text-3xl text-muted-foreground">{storageUnit}</p>
				</span>
			</Card.Content>
			<div
				class="absolute top-2 -right-16 text-8xl tracking-tight text-nowrap text-primary/5 uppercase group-hover:hidden"
			>
				<Icon icon={backgroundIcon} class="size-72" />
			</div>
		</Card.Root>
	{/snippet}
	{@render Disks()}
</div>
