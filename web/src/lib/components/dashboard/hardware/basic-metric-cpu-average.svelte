<script lang="ts">
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import type { TimeRange } from '$lib/components/custom/date-timestamp-range-picker';
	import ComponentLoading from '$lib/components/otterscale/ui/component-loading.svelte';
	import { formatTime } from '$lib/formatter';
	import { AreaChart } from 'layerchart';
	import { PrometheusDriver, SampleValue } from 'prometheus-query';
	import { onMount } from 'svelte';
	import { integrateSerieses } from '../utils';
	import * as Empty from '../utils/empty';
	import * as Template from '../utils/templates';

	let renderContext: 'svg' | 'canvas' = 'svg';
	let debug = false;

	let {
		client,
		scope: scope,
		timeRange
	}: { client: PrometheusDriver; scope: Scope; timeRange: TimeRange } = $props();

	const step = 1 * 60;
	const loadAverage1m = $derived(
		`
		node_load1{instance="juju-1eb21e-0-lxd-1"}
		/ on(instance) group_left() count by (instance) (node_cpu_seconds_total{instance="juju-1eb21e-0-lxd-1", mode="idle"})
		`
	);
	const loadAverage5m = $derived(
		`
		node_load5{instance="juju-1eb21e-0-lxd-1"}
		/ on(instance) group_left() count by (instance) (node_cpu_seconds_total{instance="juju-1eb21e-0-lxd-1", mode="idle"})
		`
	);
	const loadAverage15m = $derived(
		`
		node_load15{instance="juju-1eb21e-0-lxd-1"}
		/ on(instance) group_left() count by (instance) (node_cpu_seconds_total{instance="juju-1eb21e-0-lxd-1", mode="idle"})
		`
	);

	async function fetch(query: string) {
		try {
			let sampleSpace = [] as SampleValue[];

			const response = await client.rangeQuery(
				query,
				timeRange.start.getTime(),
				timeRange.end.getTime(),
				step
			);
			response.result.forEach((series) => {
				series.values.forEach((sampleValue: SampleValue) => {
					sampleSpace.push(sampleValue);
				});
			});

			sampleSpace.sort((p, n) => p.time.getTime() - n.time.getTime());

			return sampleSpace;
		} catch (error) {
			console.error('Error fetching:', error);
		}
	}

	let serieses: Map<string, SampleValue[] | undefined> = $state(new Map());

	let mounted = $state(false);

	onMount(async () => {
		try {
			const loadAverage1mResponse = await fetch(loadAverage1m);
			serieses.set('1m avg', loadAverage1mResponse);

			const loadAverage5mResponse = await fetch(loadAverage5m);
			serieses.set('5m avg', loadAverage5mResponse);

			const loadAverage15mResponse = await fetch(loadAverage15m);
			serieses.set('15m avg', loadAverage15mResponse);

			mounted = true;
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

{#if mounted}
	{@const data = integrateSerieses(serieses)}
	<Template.Area title="CPU">
		{#snippet hint()}
			<p>Basic CPU Information</p>
		{/snippet}
		{#snippet description()}
			<p class="text-xl">Load Average</p>
		{/snippet}
		{#snippet content()}
			{#if data.length === 0}
				<Empty.Area />
			{:else}
				<div class="h-[200px] w-full resize overflow-visible">
					<AreaChart
						{data}
						x="time"
						yDomain={[0, 1]}
						series={[
							{ key: '1m avg', color: 'hsl(var(--color-primary))' },
							{ key: '5m avg', color: 'hsl(var(--color-secondary))' },
							{ key: '15m avg', color: 'hsl(var(--color-info))' }
						]}
						legend={{
							classes: { root: '-mb-[50px] w-full overflow-auto' }
						}}
						props={{
							xAxis: { format: formatTime },
							yAxis: { format: 'percentRound' },
							tooltip: {
								root: { class: 'bg-white/60 p-3 rounded shadow-lg' },
								header: { class: 'font-light' },
								item: { format: 'percent' }
							}
						}}
						{renderContext}
						{debug}
					/>
				</div>
			{/if}
		{/snippet}
	</Template.Area>
{:else}
	<ComponentLoading />
{/if}
