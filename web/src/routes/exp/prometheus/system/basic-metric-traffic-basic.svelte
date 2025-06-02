<script lang="ts">
	import { PrometheusDriver, SampleValue } from 'prometheus-query';
	import { AreaChart, Svg, Axis, Points, Highlight } from 'layerchart';
	import * as Card from '$lib/components/ui/card';
	import Icon from '@iconify/svelte';
	import { onMount } from 'svelte';
	import ComponentLoading from '$lib/components/otterscale/ui/component-loading.svelte';
	import * as HoverCard from '$lib/components/ui/hover-card/index.js';
	import { Button } from '$lib/components/ui/button';
	import { integrateSerieses } from '../utils';
	import { formatNetworkIO } from '$lib/formatter';
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import NoData from '../utils/empty.svelte';
	import type { TimeRange } from '$lib/components/custom/date-timestamp-range-picker';
	import * as Template from '../utils/templates';

	let renderContext: 'svg' | 'canvas' = 'svg';
	let debug = false;

	let {
		client,
		scope: scope,
		instance: instance,
		timeRange
	}: { client: PrometheusDriver; scope: Scope; instance: string; timeRange: TimeRange } = $props();

	const step = 1 * 60;

	const receiveQuery = $derived(
		`
		irate(
			node_network_receive_bytes_total{device="lo",instance="${instance}",juju_model_uuid=~"${scope.uuid}"}[4m]
		)
		*
		8
		`
	);

	const transmitQuery = $derived(
		`
		irate(
			node_network_transmit_bytes_total{device="lo",instance="${instance}",juju_model_uuid=~"${scope.uuid}"}[4m]
		)
		*
		8
		`
	);

	let sampleSpaces = $state(new Map<Record<string, string>, SampleValue[]>());

	async function fetch(query: string) {
		try {
			sampleSpaces.clear();

			const response = await client.rangeQuery(
				query,
				timeRange.start.getTime(),
				timeRange.end.getTime(),
				step
			);
			response.result.forEach((series) => {
				const label = series.metric.labels;

				if (!sampleSpaces.has(label)) {
					sampleSpaces.set(label, []);
				}

				series.values.forEach((sampleValue: SampleValue) => {
					sampleSpaces.get(label)?.push(sampleValue);
				});
			});

			return sampleSpaces;
		} catch (error) {
			console.error('Error fetching:', error);
		}
	}

	let serieses: Map<string, SampleValue[] | undefined> = $state(new Map());
	let keys: string[] = $state([]);
	const colors = ['hsl(var(--color-primary))', 'hsl(var(--color-secondary))'];

	let mounted = $state(false);
	onMount(async () => {
		try {
			const receiveResponse = await fetch(receiveQuery);
			receiveResponse?.forEach((value: SampleValue[], key: Record<string, string>) => {
				const k = `receive ${key.device}`;
				keys.push(k);
				serieses.set(k, value);
			});
			const transmitResponse = await fetch(transmitQuery);
			transmitResponse?.forEach((value: SampleValue[], key: Record<string, string>) => {
				const k = `transmit ${key.device}`;
				keys.push(k);
				serieses.set(k, value);
			});

			mounted = true;
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

{#if mounted}
	{@const data = integrateSerieses(serieses)}
	<Template.Area title="Network Traffic">
		{#snippet hint()}
			<p>Basic Network Information per Interface</p>
		{/snippet}
		{#snippet content()}
			{#if data.length === 0}
				<NoData type="area" />
			{:else}
				<div class="h-[200px] w-full resize overflow-visible">
					<AreaChart
						{data}
						x="time"
						series={keys.map((k, i) => ({
							key: k,
							color: colors[i]
						}))}
						legend={{
							classes: { root: '-mb-[50px] w-full overflow-auto' }
						}}
						props={{
							tooltip: {
								root: { class: 'bg-white/60 p-3 rounded shadow-lg max-h-[50vh] overflow-auto' },
								header: { class: 'font-light' },
								item: {
									format: (v: number) => {
										const capacity = formatNetworkIO(v);
										return `${capacity.value} ${capacity.unit}`;
									}
								}
							},
							yAxis: {
								format: (v: number) => `${(v / 1024 / 1024).toFixed(0)} Mb/s`
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
