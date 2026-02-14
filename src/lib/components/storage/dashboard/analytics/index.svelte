<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';

	import type { Scope } from '$lib/api/scope/v1/scope_pb';

	import { default as AreaCapacity } from './area-chart-capacity.svelte';
	import { default as AreaIOPS } from './area-chart-iops.svelte';
	import { default as AreaOSDLatency } from './area-chart-osd-latency.svelte';
	import { default as AreaPGStates } from './area-chart-placement-group-states.svelte';
	import { default as AreaRecovery } from './area-chart-recovery.svelte';
	import { default as AreaStuckPGs } from './area-chart-stuck-placement-groups.svelte';
	import { default as AreaThroughput } from './area-chart-throughput.svelte';
	import { default as TextClusterHealth } from './text-chart-cluster-health.svelte';
	import { default as TextCPU } from './text-chart-cpu.svelte';
	import { default as TextDisk } from './text-chart-disk.svelte';
	import { default as TextIOPSRead } from './text-chart-iops-read.svelte';
	import { default as TextIOPSWrite } from './text-chart-iops-write.svelte';
	import { default as TextNetwork } from './text-chart-network.svelte';
	import { default as TextOSD } from './text-chart-osd.svelte';
	import { default as TextPool } from './text-chart-pool.svelte';
	import { default as TextRAM } from './text-chart-ram.svelte';
	import { default as TextThroughputRead } from './text-chart-throughput-read.svelte';
	import { default as TextThroughputWrite } from './text-chart-throughput-write.svelte';
	import { default as UsageCapacity } from './usage-rate-chart-capacity.svelte';
	import { default as UsageOSDInOut } from './usage-rate-chart-osd-in-out.svelte';
	import { default as UsageOSDQuorum } from './usage-rate-chart-osd-quorum.svelte';
	import { default as UsageOSDUpDown } from './usage-rate-chart-osd-up-down.svelte';

	let { client, scope }: { client: PrometheusDriver; scope: Scope } = $props();
</script>

<div class="flex flex-col gap-4">
	<div class="grid w-full gap-4 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-3">
		<span class="col-span-1">
			<TextClusterHealth {client} {scope} />
		</span>
		<span class="col-span-1">
			<TextOSD {client} {scope} />
		</span>
		<span class="col-span-1">
			<TextPool {client} {scope} />
		</span>
	</div>

	<div class="grid w-full gap-4 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-3">
		<span class="col-span-1">
			<UsageOSDInOut {client} {scope} />
		</span>
		<span class="col-span-1">
			<UsageOSDUpDown {client} {scope} />
		</span>
		<span class="col-span-1">
			<UsageOSDQuorum {client} {scope} />
		</span>
	</div>

	<div class="grid w-full gap-4 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4">
		<span class="col-span-2">
			<UsageCapacity {client} {scope} />
		</span>
		<span class="col-span-2">
			<AreaCapacity {client} {scope} />
		</span>
	</div>

	<div class="grid w-full gap-4 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4">
		<span class="col-span-1">
			<TextThroughputRead {client} {scope} />
		</span>
		<span class="col-span-1">
			<TextThroughputWrite {client} {scope} />
		</span>
		<span class="col-span-2">
			<AreaThroughput {client} {scope} />
		</span>
	</div>

	<div class="grid w-full gap-4 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4">
		<span class="col-span-1">
			<TextIOPSRead {client} {scope} />
		</span>
		<span class="col-span-1">
			<TextIOPSWrite {client} {scope} />
		</span>
		<span class="col-span-2">
			<AreaIOPS {client} {scope} />
		</span>
	</div>

	<div class="grid w-full gap-4 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4">
		<span class="col-span-1">
			<TextCPU {client} {scope} />
		</span>
		<span class="col-span-1">
			<TextRAM {client} {scope} />
		</span>
		<span class="col-span-1">
			<TextDisk {client} {scope} />
		</span>
		<span class="col-span-1">
			<TextNetwork {client} {scope} />
		</span>
	</div>

	<div class="grid w-full gap-4 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4">
		<span class="col-span-2">
			<AreaOSDLatency {client} {scope} />
		</span>
		<span class="col-span-2">
			<AreaRecovery {client} {scope} />
		</span>
	</div>

	<div class="grid w-full gap-4 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4">
		<span class="col-span-2">
			<AreaPGStates {client} {scope} />
		</span>
		<span class="col-span-2">
			<AreaStuckPGs {client} {scope} />
		</span>
	</div>
</div>
