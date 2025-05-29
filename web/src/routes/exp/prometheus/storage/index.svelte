<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';
	import { getLocalTimeZone, today, now } from '@internationalized/date';

	import ScopePicker from '../utils/scope-picker.svelte';
	import InstancePicker from '../utils/instance-picker.svelte';

	import {
		DateTimestampPicker,
		type TimeRange
	} from '$lib/components/custom/date-timestamp-range-picker';

	import { default as ClusterHealthStatus } from './cluster-health/index.svelte';

	import { default as ClusterCapacity } from './cluster-capacity/index.svelte';
	import { default as RangeClusterCapacity } from './range-cluster-capacity.svelte';

	import { default as ClusterMonitor } from './cluster-monitor/index.svelte';

	import { default as ClusterWriteThroughput } from './cluster-write-throughput/index.svelte';
	import { default as ClusterReadThroughput } from './cluster-read-throughput/index.svelte';
	import { default as RangeClusterThroughput } from './range-cluster-throughput.svelte';

	import { default as ClusterReadIOPS } from './cluster-read-iops/index.svelte';
	import { default as ClusterWriteIOPS } from './cluster-write-iops/index.svelte';
	import { default as RangeClusterIOPS } from './range-cluster-iops.svelte';

	import { default as PoolMeta } from './pool-meta/index.svelte';

	import { default as PoolRawCapacity } from './pool-raw-capacity/index.svelte';

	import { default as PoolCompression } from './pool-compression/index.svelte';

	import { default as OSDHosts } from './osd-hosts/index.svelte';
	import { default as OSDHostsIn } from './osd-hosts-in/index.svelte';
	import { default as OSDHostsUp } from './osd-hosts-up/index.svelte';

	import { default as OSDCPUBusy } from './osd-cpu/index.svelte';

	import { default as OSDRAMUtilization } from './osd-ram/index.svelte';

	import { default as OSDCapacity } from './osd-capacity/index.svelte';

	import { default as OSDDiskUtilization } from './osd-disk/index.svelte';

	import { default as OSDPhysicalIOPS } from './osd-physical-iops/index.svelte';

	import { default as OSDNetworkLoad } from './osd-network/index.svelte';

	import { default as RangeOSDLatency } from './range-osd-latency.svelte';

	import { default as RangePlacementGroupStates } from './range-placement-group-states.svelte';

	import type { Scope } from '$gen/api/scope/v1/scope_pb';

	let {
		client,
		scopes,
		instances
	}: { client: PrometheusDriver; scopes: Scope[]; instances: string[] } = $props();

	let selectedTimeRange = $state({
		start: today(getLocalTimeZone()).toDate(getLocalTimeZone()),
		end: now(getLocalTimeZone()).toDate()
	} as TimeRange);
	let selectedScope = $state(scopes[0]);
	let selectedInstance = $state(instances[0]);
</script>

<div class="flex flex-col gap-4">
	<div class="mr-auto flex flex-wrap items-center gap-2">
		<ScopePicker bind:selectedScope {scopes} />
	</div>
	<div class="mr-auto flex flex-wrap items-center gap-2">
		<DateTimestampPicker bind:value={selectedTimeRange} />
	</div>
	{#key selectedScope}
		{#key selectedTimeRange}
			<p class="text-xl font-bold">Cluster</p>
			<div class="grid w-full gap-4 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4">
				<span class="col-span-1">
					<ClusterHealthStatus {client} scope={selectedScope} />
				</span>
				<span class="col-span-1">
					<ClusterCapacity {client} scope={selectedScope} />
				</span>
				<span class="col-span-1">
					<ClusterMonitor {client} scope={selectedScope} />
				</span>
			</div>
			<div class="grid w-full gap-4 sm:grid-cols-1 md:grid-cols-2">
				<span class="col-span-1">
					<RangeClusterCapacity {client} scope={selectedScope} timeRange={selectedTimeRange} />
				</span>
			</div>
			<div class="grid w-full gap-4 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4">
				<span class="col-span-1">
					<ClusterWriteThroughput {client} scope={selectedScope} />
				</span>
				<span class="col-span-1">
					<ClusterReadThroughput {client} scope={selectedScope} />
				</span>
				<span class="col-span-1">
					<ClusterWriteIOPS {client} scope={selectedScope} />
				</span>
				<span class="col-span-1">
					<ClusterReadIOPS {client} scope={selectedScope} />
				</span>
			</div>
			<div class="grid w-full gap-4 sm:grid-cols-1 md:grid-cols-2">
				<span class="col-span-1">
					<RangeClusterThroughput {client} scope={selectedScope} timeRange={selectedTimeRange} />
				</span>
				<span class="col-span-1">
					<RangeClusterIOPS {client} scope={selectedScope} timeRange={selectedTimeRange} />
				</span>
			</div>

			<p class="text-xl font-bold">Pool</p>
			<div class="grid w-full gap-4 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4">
				<span class="col-span-1">
					<PoolMeta {client} scope={selectedScope} />
				</span>
				<span class="col-span-1">
					<PoolRawCapacity {client} scope={selectedScope} />
				</span>
			</div>
			<p class="text-xl font-bold">OSD Host</p>
			<div class="grid w-full gap-4 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4">
				<span class="col-span-1">
					<OSDHosts {client} scope={selectedScope} />
				</span>
				<span class="col-span-1">
					<OSDHostsIn {client} scope={selectedScope} />
				</span>
				<span class="col-span-1">
					<OSDHostsUp {client} scope={selectedScope} />
				</span>
			</div>
			<div class="grid w-full gap-4 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4">
				<span class="col-span-1">
					<OSDCPUBusy {client} scope={selectedScope} />
				</span>
				<span class="col-span-1">
					<OSDRAMUtilization {client} scope={selectedScope} />
				</span>
				<span class="col-span-1">
					<OSDDiskUtilization {client} scope={selectedScope} />
				</span>
				<span class="col-span-1">
					<OSDCapacity {client} scope={selectedScope} />
				</span>
			</div>
			<div class="grid w-full gap-4 sm:grid-cols-1 md:grid-cols-2">
				<span class="col-span-1">
					<RangeOSDLatency {client} scope={selectedScope} timeRange={selectedTimeRange} />
				</span>
				<span class="col-span-1">
					<RangePlacementGroupStates {client} scope={selectedScope} timeRange={selectedTimeRange} />
				</span>
			</div>

			<div class="grid w-full gap-4 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4">
				<span class="col-span-1">
					<OSDPhysicalIOPS {client} scope={selectedScope} />
				</span>
				<span class="col-span-1">
					<OSDNetworkLoad {client} scope={selectedScope} />
				</span>
			</div>

			<!-- {#key selectedTimeRange}
				<p class="text-xl font-bold">Basic Metric</p>
				<div class="grid w-full gap-4 sm:grid-cols-1 md:grid-cols-2">
					<span class="col-span-1">
						<BasicCPU
							{client}
							scope={selectedScope}
							instance={selectedInstance}
							timeRange={selectedTimeRange}
						/>
					</span>
					<span class="col-span-1">
						<BasicRAM
							{client}
							scope={selectedScope}
							instance={selectedInstance}
							timeRange={selectedTimeRange}
						/>
					</span>
					<span class="col-span-1">
						<BasicSWAP
							{client}
							scope={selectedScope}
							instance={selectedInstance}
							timeRange={selectedTimeRange}
						/>
					</span>
					<span class="col-span-1">
						<NetworkTrafficBasic
							{client}
							scope={selectedScope}
							instance={selectedInstance}
							timeRange={selectedTimeRange}
						/>
					</span>
				</div>
			{/key} -->
		{/key}
	{/key}
</div>
