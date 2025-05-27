<script lang="ts" module>
	import { getLocalTimeZone, today } from '@internationalized/date';

	const DEFAULT_DURATION = 7;
	const DEFAULT_END_POINT = today(getLocalTimeZone());
	const DEFAULT_START_POINT = DEFAULT_END_POINT.subtract({ days: DEFAULT_DURATION });
</script>

<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';

	import ScopePicker from './utils/scope-picker.svelte';
	import InstancePicker from './utils/instance-picker.svelte';
	import { DateTimestampPicker } from '$lib/components/custom/date-timestamp-range-picker';
	import { type TimeRange } from '$lib/components/custom/date-timestamp-range-picker';

	import QuickUptime from './uptime/quick-metric.svelte';
	import QuickCPU from './cpu/quick-metric.svelte';
	import BasicCPU from './cpu/basic-metric.svelte';
	import BasicRAM from './ram/basic-metric.svelte';
	import QuickRAM from './ram/quick-metric.svelte';
	import QuickSWAP from './swap/quick-metric.svelte';
	import BasicSWAP from './swap/basic-metric.svelte';
	import QuickRootFS from './root-fs/quick-metric.svelte';

	import NetworkTrafficBasic from './network/traffic-basic.svelte';
	import type { Scope } from '$gen/api/nexus/v1/nexus_pb';

	let {
		client,
		scopes,
		instances
	}: { client: PrometheusDriver; scopes: Scope[]; instances: string[] } = $props();

	let selectedTimeRange = $state({
		start: DEFAULT_START_POINT.toDate(getLocalTimeZone()),
		end: DEFAULT_END_POINT.toDate(getLocalTimeZone())
	} as TimeRange);
	let selectedScope = $state(scopes[1]);
	let selectedInstance = $state(instances[0]);
</script>

<main class="no-user-select grid gap-4 p-4">
	<div class="mr-auto flex items-center gap-2">
		<ScopePicker bind:selectedScope {scopes} />
		<InstancePicker bind:selectedInstance {instances} />
		<DateTimestampPicker bind:value={selectedTimeRange} />
	</div>
	{#key selectedScope}
		{#key selectedInstance}
			<p class="text-xl font-bold">Quick Metric</p>
			<div class="grid w-full gap-4 sm:grid-cols-1 md:grid-cols-3 lg:grid-cols-5">
				<span class="col-span-1">
					<QuickUptime {client} scope={selectedScope} instance={selectedInstance} />
				</span>
				<span class="col-span-1">
					<QuickCPU {client} scope={selectedScope} instance={selectedInstance} />
				</span>
				<span class="col-span-1">
					<QuickRAM {client} scope={selectedScope} instance={selectedInstance} />
				</span>
				<span class="col-span-1">
					<QuickSWAP {client} scope={selectedScope} instance={selectedInstance} />
				</span>
				<span class="col-span-1">
					<QuickRootFS {client} scope={selectedScope} instance={selectedInstance} />
				</span>
			</div>
			<p class="text-xl font-bold">Basic Metric</p>
			<div class="grid w-full gap-4 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-3">
				<span class="col-span-1">
					<BasicCPU
						{client}
						scope={selectedScope}
						instance={selectedInstance}
						timeRange={selectedTimeRange}
					/>
				</span>
				<span class="col-span-1">
					<BasicRAM {client} scope={selectedScope} instance={selectedInstance} timeRange={selectedTimeRange}/>
				</span>
				<span class="col-span-1">
					<BasicSWAP {client} scope={selectedScope} instance={selectedInstance} timeRange={selectedTimeRange}/>
				</span>
				<span class="col-span-1">
					<NetworkTrafficBasic {client} scope={selectedScope} instance={selectedInstance} timeRange={selectedTimeRange}/>
				</span>
			</div>
		{/key}
	{/key}
</main>

<style>
	.no-user-select {
		user-select: none;
		-webkit-user-select: none;
		-moz-user-select: none;
		-ms-user-select: none;
	}
</style>
