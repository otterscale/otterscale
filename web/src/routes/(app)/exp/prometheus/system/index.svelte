<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';
	import { getLocalTimeZone, today, now } from '@internationalized/date';

	import ScopePicker from '../utils/scope-picker.svelte';
	import InstancePicker from '../utils/instance-picker.svelte';
	import {
		DateTimestampPicker,
		type TimeRange
	} from '$lib/components/custom/date-timestamp-range-picker';

	import { default as QuickUptime } from './uptime/quick-metric.svelte';
	import { default as QuickCPU } from './cpu/quick-metric.svelte';
	import { default as BasicCPU } from './cpu/basic-metric.svelte';
	import { default as BasicRAM } from './ram/basic-metric.svelte';
	import { default as QuickRAM } from './ram/quick-metric.svelte';
	import { default as QuickSWAP } from './swap/quick-metric.svelte';
	import { default as BasicSWAP } from './swap/basic-metric.svelte';
	import { default as QuickRootFS } from './root-fs/quick-metric.svelte';
	import { default as NetworkTrafficBasic } from './network/traffic-basic.svelte';

	import type { Scope } from '$gen/api/nexus/v1/nexus_pb';

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
		<InstancePicker bind:selectedInstance {instances} />
	</div>
	<div class="mr-auto flex flex-wrap items-center gap-2">
		<DateTimestampPicker bind:value={selectedTimeRange} />
	</div>
	{#key selectedScope}
		{#key selectedInstance}
			<p class="text-xl font-bold">Quick Metric</p>
			<div class="grid w-full gap-4 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4">
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
			{#key selectedTimeRange}
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
			{/key}
		{/key}
	{/key}
</div>

<style>
	.no-user-select {
		user-select: none;
		-webkit-user-select: none;
		-moz-user-select: none;
		-ms-user-select: none;
	}
</style>
