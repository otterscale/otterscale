<script lang="ts">
	import type { Machine } from '$gen/api/machine/v1/machine_pb';
	import type { TimeRange } from '$lib/components/custom/date-timestamp-range-picker';
	import ComponentLoading from '$lib/components/otterscale/ui/component-loading.svelte';
	import { formatNetworkIO, formatTime } from '$lib/formatter';
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
	machine,
	timeRange
}: { client: PrometheusDriver; machine: Machine; timeRange: TimeRange } = $props();

	const step = 1 * 60;
	const query = $derived(
		`
		rate(node_network_transmit_bytes_total{instance=~"${machine.fqdn}", device!="lo"}[5m])
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
	const colors = [
		'hsl(var(--color-primary))',
		'hsl(var(--color-secondary))',
		'hsl(var(--color-info))',
		'hsl(var(--color-success))',
		'hsl(var(--color-danger))'
	];

	let mounted = $state(false);
	onMount(async () => {
		try {
			const response = await fetch(query);
			response?.forEach((value: SampleValue[], key: Record<string, string>) => {
				keys.push(key['device']);
				serieses.set(key['device'], value);
			});

			mounted = true;
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

{#if mounted}
	{@const data = integrateSerieses(serieses)}
	<Template.Area title="Network Transmitted">
		{#snippet hint()}
			<p>Historical view of network transmitted over the past hour</p>
		{/snippet}
		{#snippet content()}
			{#if data.length === 0}
				<Empty.Area />
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
										const io = formatNetworkIO(v);
										return `${io.value} ${io.unit}`;
									}
								}
							},
							xAxis: { format: formatTime },
							yAxis: {
								format: (v: number) => {
									const io = formatNetworkIO(v);
									return `${io.value.toFixed(0)} ${io.unit}`;
								}
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
