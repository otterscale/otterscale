<script lang="ts">
	import Icon from '@iconify/svelte';
	import { type Table } from '@tanstack/table-core';

	import { type VirtualMachine } from '$lib/api/instance/v1/instance_pb';
	import * as Card from '$lib/components/ui/card';
	import { m } from '$lib/paraglide/messages';

	let {
		table,
		scope: _
	}: {
		table: Table<VirtualMachine>;
		scope: string;
	} = $props();

	const filteredVirtualMachines = $derived(
		table.getFilteredRowModel().rows.map((row) => row.original)
	);

	// Calculate statistics
	// const totalMachines = $derived(virtualMachines.length);
	// const totalDisks = $derived(virtualMachines.reduce((acc, vm) => acc + vm.disks.length, 0));
	// const totalDataVolumes = $derived(

	// );
	// const storageFormatted = $derived(
	// 	formatCapacity(
	// 		virtualMachines.reduce((acc, vm) => {
	// 			return (
	// 				acc +
	// 				vm.disks.reduce((diskAcc, disk) => {
	// 					if (disk.sourceData.case === 'dataVolume') {
	// 						return diskAcc + Number(disk.sourceData.value.sizeBytes);
	// 					}
	// 					return diskAcc;
	// 				}, 0)
	// 			);
	// 		}, 0),
	// 	),
	// );
	// const machinesOn = $derived(virtualMachines.filter((vm) => vm.status === 'Running').length);
	// const powerOnPercentage = $derived(
	// 	totalMachines === 0 ? 0 : Math.round((machinesOn / totalMachines) * 100)
	// );
</script>

<div class="grid w-full gap-3 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4 xl:grid-cols-5">
	<!-- 	
	<Layout>
		{#snippet title()}
			<Title title="VIRTUAL MACHINE" />
		{/snippet}

		{#snippet content()}
			<Content value={totalMachines} />
		{/snippet}
	</Layout>

	<Layout>
		{#snippet title()}
			<Title title="DISK" />
		{/snippet}

		{#snippet content()}
			<Content value={totalDisks} />
		{/snippet}
	</Layout> -->

	<!-- <Layout>
		{#snippet title()}
			<Title title="DATA VOLUME" />
		{/snippet}

		{#snippet content()}
			<Content>
				<span>{storageFormatted.value}</span>
				<span class="text-3xl font-extralight">
					{storageFormatted.unit}
				</span>
			</Content>
		{/snippet}

		{#snippet footer()}
			<p class="text-muted-foreground text-xs">
				over {totalDataVolumes} Data Volumes
			</p>
		{/snippet}
	</Layout> -->
	<!-- 
	<Layout>
		{#snippet title()}
			<Title title="POWER ON" />
		{/snippet}

		{#snippet content()}
			<ContentSubtitle
				value={powerOnPercentage}
				unit="%"
				subtitle={`${machinesOn} On over ${totalMachines} units`}
			/>
		{/snippet}

		{#snippet footer()}
			<Progress
				value={powerOnPercentage}
				max={100}
				class={formatProgressColor(powerOnPercentage)}
			/>
		{/snippet}
	</Layout> -->

	{#snippet VirtualMachines()}
		{@const title = m.virtual_machine()}
		{@const titleIcon = 'ph:chart-bar-bold'}
		{@const backgroundIcon = 'ph:desktop-tower'}
		{@const virtualMachines = filteredVirtualMachines.length}
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
				<p class="text-7xl font-semibold">{virtualMachines}</p>
			</Card.Content>
			<div
				class="absolute top-0 -right-16 text-8xl tracking-tight text-nowrap text-primary/5 uppercase group-hover:hidden"
			>
				<Icon icon={backgroundIcon} class="size-72" />
			</div>
		</Card.Root>
	{/snippet}
	{@render VirtualMachines()}

	{#snippet power()}
		{@const title = m.power()}
		{@const titleIcon = 'ph:chart-bar-bold'}
		{@const backgroundIcon = 'ph:cube'}
		{@const on = filteredVirtualMachines.filter((vm) => vm.status === 'Running').length}
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
				<p class="text-7xl font-semibold">{on}</p>
			</Card.Content>
			<div
				class="absolute top-0 -right-16 text-8xl tracking-tight text-nowrap text-primary/5 uppercase group-hover:hidden"
			>
				<Icon icon={backgroundIcon} class="size-72" />
			</div>
		</Card.Root>
	{/snippet}
	{@render power()}
</div>
