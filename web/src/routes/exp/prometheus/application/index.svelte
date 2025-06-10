<script lang="ts">
	import { getLocalTimeZone, now, today } from '@internationalized/date';
	import { PrometheusDriver } from 'prometheus-query';
	import {
		DateTimestampPicker,
		type TimeRange
	} from '$lib/components/custom/date-timestamp-range-picker';
	import ScopePicker from '../utils/scope-picker.svelte';
	import { default as APIServerAvailability30Days } from './api-server-availability-30days/index.svelte';
	import { default as ControllerManagerUp } from './controller-manager-up/index.svelte';
	import { default as ETCDUp } from './etcd-up/index.svelte';
	import { default as HistoricalErrorBudget30Days } from './historical-error-budget-30days.svelte';
	import { default as KubeletContainerRunning } from './kubelet-container-running/index.svelte';
	import { default as KubeletPodRunning } from './kubelet-pod-running/index.svelte';
	import { default as KubeletRunning } from './kubelet-running/index.svelte';
	import { default as ProxyUp } from './proxy-up/index.svelte';
	import { default as SchedulerUp } from './scheduler-up/index.svelte';
	import { default as HistoricalVolumnSpaceUsage } from './historical-volumn-space-usage.svelte';
	import { default as VolumnSpaceUsage } from './volumn-space-usage/index.svelte';

	import type { Scope } from '$gen/api/scope/v1/scope_pb';

	let { client, scopes }: { client: PrometheusDriver; scopes: Scope[] } = $props();

	let selectedTimeRange = $state({
		start: today(getLocalTimeZone()).toDate(getLocalTimeZone()),
		end: now(getLocalTimeZone()).toDate()
	} as TimeRange);
	let selectedScope = $state(scopes[0]);
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
			<div class="grid w-full gap-4 sm:grid-cols-1 md:grid-cols-3">
				<span class="col-span-1">
					<KubeletRunning {client} scope={selectedScope} />
				</span>
				<span class="col-span-1">
					<KubeletPodRunning {client} scope={selectedScope} />
				</span>
				<span class="col-span-1">
					<KubeletContainerRunning {client} scope={selectedScope} />
				</span>
			</div>
			<div class="grid w-full gap-4 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4">
				<span class="col-span-1">
					<ControllerManagerUp {client} scope={selectedScope} />
				</span>
				<span class="col-span-1">
					<ETCDUp {client} scope={selectedScope} />
				</span>
				<span class="col-span-1">
					<ProxyUp {client} scope={selectedScope} />
				</span>
				<span class="col-span-1">
					<SchedulerUp {client} scope={selectedScope} />
				</span>
			</div>

			<div class="grid w-full gap-4 sm:grid-cols-1 md:grid-cols-3">
				<span class="col-span-1">
					<APIServerAvailability30Days {client} scope={selectedScope} />
				</span>
				<span class="sm:col-span-1 md:col-span-2">
					<HistoricalErrorBudget30Days
						{client}
						scope={selectedScope}
						timeRange={selectedTimeRange}
					/>
				</span>
			</div>

			<div class="grid w-full gap-4 sm:grid-cols-1 md:grid-cols-3">
				<span class="col-span-1">
					<VolumnSpaceUsage {client} scope={selectedScope} />
				</span>
				<span class="sm:col-span-1 md:col-span-2">
					<HistoricalVolumnSpaceUsage
						{client}
						scope={selectedScope}
						timeRange={selectedTimeRange}
					/>
				</span>
			</div>
		{/key}
	{/key}
</div>
