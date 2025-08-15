<script lang="ts">
	import type { Scope } from '$lib/api/scope/v1/scope_pb';
	import { type Machine } from '$lib/api/machine/v1/machine_pb';
	import { default as Pickers } from '$lib/components/machines/units/machine-picker.svelte';
	import { PrometheusDriver } from 'prometheus-query';
	import { default as TextClusterHealth } from './text-chart-cluster-health.svelte';
	import { default as TextOSD } from './text-chart-osd.svelte';
	import { default as TextPool } from './text-chart-pool.svelte';
	import { default as UsageOSDInOut } from './usage-rate-chart-osd-in-out.svelte';
	import { default as UsageOSDUpDown } from './usage-rate-chart-osd-up-down.svelte';
	import { default as UsageCapacity } from './usage-rate-chart-capacity.svelte';
	import { default as AreaCapacity } from './area-chart-capacity.svelte';

	//
	// import { default as CPUAverage } from './area-chart-cpu-average.svelte';
	// import { default as CPUCoreProcessor } from './area-chart-cpu-core.svelte';
	// import { default as DiskIOTime } from './area-chart-disk-io-time.svelte';
	// import { default as BasicDisk } from './area-chart-disk-rw.svelte';
	// import { default as NetworkReceived } from './area-chart-network-received.svelte';
	// import { default as NetworkTransmitted } from './area-chart-network-transmitted.svelte';
	// import { default as BasicRAM } from './area-chart-ram.svelte';
	// import { default as UsageRateCPU } from './usage-rate-chart-cpu.svelte';
	// import { default as UsageRateRAM } from './usage-rate-chart-ram.svelte';
	// import { default as UsageRateRootFS } from './usage-rate-chart-root-fs.svelte';
	// import { default as UsageRateSWAP } from './usage-rate-chart-swap.svelte';

	let { client, scope }: { client: PrometheusDriver; scope: Scope } = $props();
</script>

<div class="flex flex-col gap-4">
	<div class="grid w-full gap-4 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4">
		<span class="col-span-1">
			<TextClusterHealth {client} {scope} />
		</span>
		<span class="col-span-1">
			<TextOSD {client} {scope} />
		</span>
		<span class="col-span-1">
			<TextPool {client} {scope} />
		</span>
		<!-- <span class="col-span-1">
			<TextPool {client} {scope} />
		</span> -->
	</div>

	<div class="grid w-full gap-4 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4">
		<span class="col-span-1">
			<UsageOSDInOut {client} {scope} />
		</span>
		<span class="col-span-1">
			<UsageOSDUpDown {client} {scope} />
		</span>
		<span class="col-span-1">
			<UsageOSDUpDown {client} {scope} />
		</span>
		<span class="col-span-1">
			<!-- <UsageOSDUpDown {client} {scope} /> -->
		</span>
	</div>

	<div class="grid w-full gap-4 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4">
		<span class="col-span-1">
			<UsageCapacity {client} {scope} />
		</span>
		<span class="col-span-2">
			<AreaCapacity {client} {scope} />
		</span>
	</div>

	<!-- <div class="grid w-full gap-4 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4">
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
	</div> -->

	<!-- <div class="grid w-full gap-4 sm:grid-cols-1 md:grid-cols-2">
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
	</div> -->
</div>
