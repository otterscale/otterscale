<script lang="ts">
	import { type Machine } from '$lib/api/machine/v1/machine_pb';
	import Content from '$lib/components/custom/chart/content/text/text-large.svelte';
	import ContentSubtitle from '$lib/components/custom/chart/content/text/text-with-subtitle.svelte';
	import Layout from '$lib/components/custom/chart/layout/small-flexible-height.svelte';
	import Title from '$lib/components/custom/chart/title.svelte';
	import { Progress } from '$lib/components/ui/progress/index.js';
	import { formatCapacity, formatProgressColor } from '$lib/formatter';

	let { machines }: { machines: Machine[] } = $props();

	// Calculate statistics
	const totalMachines = $derived(machines.length);
	const totalDisks = $derived(machines.reduce((acc, machine) => acc + machine.blockDevices.length, 0));
	const storageFormatted = $derived(
		formatCapacity(machines.reduce((acc, machine) => acc + machine.storageMb, 0) * 1024 * 1024),
	);
	const machinesOn = $derived(machines.filter((m) => m.powerState === 'on').length);
	const machinesDeployed = $derived(machines.filter((m) => m.status === 'Deployed').length);
	const powerOnPercentage = $derived(Math.round((machinesOn / totalMachines) * 100));
	const deploymentPercentage = $derived(Math.round((machinesDeployed / totalMachines) * 100));
</script>

<div class="grid w-full gap-3 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4 xl:grid-cols-5">
	<Layout>
		{#snippet title()}
			<Title title="NODE" />
		{/snippet}

		{#snippet content()}
			<Content value={totalMachines} />
		{/snippet}
	</Layout>

	<Layout>
		{#snippet title()}
			<Title title="STORAGE" />
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
				over {totalDisks} disks
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

	<Layout>
		{#snippet title()}
			<Title title="DEPLOYMENT" />
		{/snippet}

		{#snippet content()}
			<ContentSubtitle
				value={deploymentPercentage}
				unit="%"
				subtitle={`${machinesDeployed} deployed over ${totalMachines} units`}
			/>
		{/snippet}

		{#snippet footer()}
			<Progress value={deploymentPercentage} max={100} class={formatProgressColor(deploymentPercentage)} />
		{/snippet}
	</Layout>
</div>
