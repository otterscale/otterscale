<script lang="ts">
	import type { Scope } from '$lib/api/scope/v1/scope_pb';
	import { PrometheusDriver } from 'prometheus-query';
	import { default as AreaCapacity } from './area-chart-capacity.svelte';
	import { default as AreaOSDReadLatency } from './area-chart-osd-read-latency.svelte';
	import { default as AreaOSDWriteLatency } from './area-chart-osd-write-latency.svelte';
	import { default as BarOSDIOPS } from './bar-chart-osd-iops.svelte';
	import { default as BarOSDThroughtput } from './bar-chart-osd-throughtput.svelte';
	import { default as PieOSDType } from './pie-chart-osd-type.svelte';
	import { default as TextClusterHealth } from './text-chart-cluster-health.svelte';
	import { default as TextOSDs } from './text-chart-osds.svelte';
	import { default as TextQuorum } from './text-chart-quorum.svelte';
	import { default as TextTimeTillFull } from './text-chart-time-till-full.svelte';
	import { default as UsageCapacity } from './usage-chart-capacity.svelte';

	let {
		client,
		scope,
		isReloading = $bindable(),
	}: { client: PrometheusDriver; scope: Scope; isReloading: boolean } = $props();
</script>

<div class="grid auto-rows-auto grid-cols-2 gap-5 pt-4 md:grid-cols-4 lg:grid-cols-10">
	<div class="col-span-2">
		<TextClusterHealth {client} {scope} />
	</div>
	<div class="col-span-2">
		<TextTimeTillFull {client} {scope} />
	</div>
	<div class="col-span-2 row-span-2">
		<UsageCapacity {client} {scope} />
	</div>
	<div class="col-span-2 row-span-2">
		<AreaCapacity {client} {scope} />
	</div>
	<div class="col-span-2 row-span-2">
		<PieOSDType {client} {scope} />
	</div>
	<div class="col-span-2">
		<TextQuorum {client} {scope} />
	</div>
	<div class="col-span-2">
		<TextOSDs {client} {scope} />
	</div>
	<div class="col-span-4 row-span-2">
		<BarOSDThroughtput {client} {scope} bind:isReloading />
	</div>
	<div class="col-span-4 row-span-2">
		<BarOSDIOPS {client} {scope} />
	</div>
	<div class="col-span-2">
		<AreaOSDReadLatency {client} {scope} />
	</div>
	<div class="col-span-2">
		<AreaOSDWriteLatency {client} {scope} />
	</div>
</div>
