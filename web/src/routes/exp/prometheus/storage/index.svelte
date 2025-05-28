<script lang="ts">
	import { PrometheusDriver } from 'prometheus-query';
	import { getLocalTimeZone, today, now } from '@internationalized/date';

	import ScopePicker from '../utils/scope-picker.svelte';
	import InstancePicker from '../utils/instance-picker.svelte';
	import {
		DateTimestampPicker,
		type TimeRange
	} from '$lib/components/custom/date-timestamp-range-picker';

	import { default as PoolMeta } from './pool/meta/index.svelte';
	import { default as PoolRawCapacity } from './pool/raw-capacity/index.svelte';
	import { default as PoolCompression } from './pool/compression/index.svelte';
	import { default as OSDHosts } from './osd/hosts/index.svelte';
	import { default as OSDPhysicalIOPS } from './osd/io/index.svelte';

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
		<p class="text-xl font-bold">Number</p>
		<div class="grid w-full gap-4 sm:grid-cols-1 md:grid-cols-3 lg:grid-cols-5">
			<span class="col-span-1">
				<PoolMeta {client} scope={selectedScope} />
			</span>
			<span class="col-span-1">
				<OSDHosts {client} scope={selectedScope} />
			</span>
		</div>
		<p class="text-xl font-bold">Usage</p>
		<div class="grid w-full gap-4 sm:grid-cols-1 md:grid-cols-3 lg:grid-cols-5">
			<span class="col-span-1">
				<PoolRawCapacity {client} scope={selectedScope} />
			</span>

			<span class="col-span-1">
				<OSDPhysicalIOPS {client} scope={selectedScope} />
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
</div>

<style>
	.no-user-select {
		user-select: none;
		-webkit-user-select: none;
		-moz-user-select: none;
		-ms-user-select: none;
	}
</style>
