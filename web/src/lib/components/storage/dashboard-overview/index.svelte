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

	let { client, scope }: { client: PrometheusDriver; scope: Scope } = $props();
</script>

<TextClusterHealth {client} {scope} />
<TextTimeTillFull {client} {scope} />
<UsageCapacity {client} {scope} />
<AreaCapacity {client} {scope} />
<PieOSDType {client} {scope} />
<TextQuorum {client} {scope} />
<TextOSDs {client} {scope} />
<BarOSDThroughtput {client} {scope} />
<BarOSDIOPS {client} {scope} />
<AreaOSDReadLatency {client} {scope} />
<AreaOSDWriteLatency {client} {scope} />
