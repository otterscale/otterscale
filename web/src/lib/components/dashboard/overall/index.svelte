<script lang="ts">
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import {
		DateTimestampPicker,
		type TimeRange
	} from '$lib/components/custom/date-timestamp-range-picker';
	import { getLocalTimeZone, now, today } from '@internationalized/date';
	import { PrometheusDriver } from 'prometheus-query';
	import * as Pickers from '../utils/pickers';
	import { default as BasicCPU } from './basic-metric-cpu.svelte';
	import { default as BasicRAM } from './basic-metric-ram.svelte';
	import { default as BasicSWAP } from './basic-metric-swap.svelte';
	import { default as BasicNetworkTrafficBasic } from './basic-metric-traffic-basic.svelte';
	import { default as QuickCPU } from './quick-metric-cpu/index.svelte';
	import { default as QuickRAM } from './quick-metric-ram/index.svelte';
	import { default as QuickRootFS } from './quick-metric-root-fs/index.svelte';
	import { default as QuickSWAP } from './quick-metric-swap/index.svelte';
	import { default as QuickUptime } from './quick-metric-uptime/index.svelte';

	let { client, scopes }: { client: PrometheusDriver; scopes: Scope[] } = $props();

	let selectedTimeRange = $state({
		start: today(getLocalTimeZone()).toDate(getLocalTimeZone()),
		end: now(getLocalTimeZone()).toDate()
	} as TimeRange);
	let selectedScope = $state(scopes[0]);
</script>

<div class="flex flex-col gap-4">
	<div class="mr-auto flex flex-wrap items-center gap-2">
		<Pickers.Scope bind:selectedScope {scopes} />
	</div>
	<div class="mr-auto flex flex-wrap items-center gap-2">
		<DateTimestampPicker bind:value={selectedTimeRange} />
	</div>
	{#key selectedScope}
		<div class="grid w-full gap-4 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4">
			<span class="col-span-1">
				<QuickUptime {client} scope={selectedScope} />
			</span>
		</div>
		<div class="grid w-full gap-4 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4">
			<span class="col-span-1">
				<QuickCPU {client} scope={selectedScope} />
			</span>
			<span class="col-span-1">
				<QuickRAM {client} scope={selectedScope} />
			</span>
			<span class="col-span-1">
				<QuickSWAP {client} scope={selectedScope} />
			</span>
			<span class="col-span-1">
				<QuickRootFS {client} scope={selectedScope} />
			</span>
		</div>
		{#key selectedTimeRange}
			<div class="grid w-full gap-4 sm:grid-cols-1 md:grid-cols-2">
				<span class="col-span-1">
					<BasicCPU {client} scope={selectedScope} timeRange={selectedTimeRange} />
				</span>
				<span class="col-span-1">
					<BasicRAM {client} scope={selectedScope} timeRange={selectedTimeRange} />
				</span>
				<span class="col-span-1">
					<BasicSWAP {client} scope={selectedScope} timeRange={selectedTimeRange} />
				</span>
				<span class="col-span-1">
					<BasicNetworkTrafficBasic {client} scope={selectedScope} timeRange={selectedTimeRange} />
				</span>
			</div>
		{/key}
	{/key}
</div>
