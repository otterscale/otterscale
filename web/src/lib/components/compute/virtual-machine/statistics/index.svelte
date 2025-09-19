<script lang="ts">
	import { type VirtualMachine, VirtualMachine_status } from '$lib/api/kubevirt/v1/kubevirt_pb';
	import Content from '$lib/components/custom/chart/content/text/text-large.svelte';
	import ContentSubtitle from '$lib/components/custom/chart/content/text/text-with-subtitle.svelte';
	import Layout from '$lib/components/custom/chart/layout/small-flexible-height.svelte';
	import Title from '$lib/components/custom/chart/title.svelte';
	import { Progress } from '$lib/components/ui/progress/index.js';
	import { formatCapacity, formatProgressColor } from '$lib/formatter';

	let { virtualMachines }: { virtualMachines: VirtualMachine[] } = $props();

	// Calculate statistics
	const totalMachines = $derived(virtualMachines.length);
	const totalDisks = $derived(virtualMachines.reduce((acc, vm) => acc + vm.disks.length, 0));
	const totalDataVolumes = $derived(
		virtualMachines.reduce(
			(acc, vm) => acc + vm.disks.filter((disk) => disk.sourceData.case === 'dataVolume').length,
			0,
		),
	);
	const storageFormatted = $derived(
		formatCapacity(
			virtualMachines.reduce((acc, vm) => {
				return (
					acc +
					vm.disks.reduce((diskAcc, disk) => {
						if (disk.sourceData.case === 'dataVolume') {
							return diskAcc + Number(disk.sourceData.value.sizeBytes);
						}
						return diskAcc;
					}, 0)
				);
			}, 0),
		),
	);
	const machinesOn = $derived(
		virtualMachines.filter((vm) => vm.statusPhase === VirtualMachine_status.RUNNING).length,
	);
	const powerOnPercentage = $derived(totalMachines === 0 ? 0 : Math.round((machinesOn / totalMachines) * 100));
</script>

<div class="grid w-full gap-3 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4 xl:grid-cols-5">
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
	</Layout>

	<Layout>
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
	</Layout>

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
			<Progress value={powerOnPercentage} max={100} class={formatProgressColor(powerOnPercentage)} />
		{/snippet}
	</Layout>
</div>
