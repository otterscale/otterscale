<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';

	import { default as CPUAverage } from './area-chart-cpu-average.svelte';
	import { default as CPUCoreProcessor } from './area-chart-cpu-core.svelte';
	import { default as DiskIOTime } from './area-chart-disk-io-time.svelte';
	import { default as BasicDisk } from './area-chart-disk-rw.svelte';
	import { default as NetworkReceived } from './area-chart-network-received.svelte';
	import { default as NetworkTransmitted } from './area-chart-network-transmitted.svelte';
	import { default as BasicRAM } from './area-chart-ram.svelte';
	import { default as UsageRateUptime } from './text-chart-uptime.svelte';
	import { default as UsageRateCPU } from './usage-rate-chart-cpu.svelte';
	import { default as UsageRateRAM } from './usage-rate-chart-ram.svelte';
	import { default as UsageRateRootFS } from './usage-rate-chart-root-fs.svelte';
	import { default as UsageRateSWAP } from './usage-rate-chart-swap.svelte';

	import { type Machine } from '$lib/api/machine/v1/machine_pb';
	import { default as Pickers } from '$lib/components/machines/units/machine-picker.svelte';

	let { client, machines }: { client: PrometheusDriver; machines: Machine[] } = $props();
	let selectedMachine = $state(machines[0]);
</script>

<div class="flex flex-col gap-4">
	<div class="mr-auto flex flex-wrap items-center gap-2">
		<Pickers bind:selectedMachine {machines} />
	</div>
	{#key selectedMachine}
		<div class="grid w-full gap-3 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-5">
			<span class="col-span-1">
				<UsageRateUptime {client} machine={selectedMachine} />
			</span>
			<span class="col-span-1">
				<UsageRateCPU {client} machine={selectedMachine} />
			</span>
			<span class="col-span-1">
				<UsageRateRAM {client} machine={selectedMachine} />
			</span>
			<span class="col-span-1">
				<UsageRateSWAP {client} machine={selectedMachine} />
			</span>
			<span class="col-span-1">
				<UsageRateRootFS {client} machine={selectedMachine} />
			</span>
		</div>

		<div class="grid w-full gap-4 sm:grid-cols-1 md:grid-cols-1 lg:grid-cols-2">
			<span class="col-span-1">
				<CPUCoreProcessor {client} machine={selectedMachine} />
			</span>
			<span class="col-span-1">
				<CPUAverage {client} machine={selectedMachine} />
			</span>
			<span class="col-span-1">
				<BasicRAM {client} machine={selectedMachine} />
			</span>
			<span class="col-span-1">
				<BasicDisk {client} machine={selectedMachine} />
			</span>
			<span class="col-span-1">
				<DiskIOTime {client} machine={selectedMachine} />
			</span>
			<span class="col-span-1">
				<NetworkReceived {client} machine={selectedMachine} />
			</span>
			<span class="col-span-1">
				<NetworkTransmitted {client} machine={selectedMachine} />
			</span>
		</div>
	{/key}
</div>
