<script lang="ts">
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import type { TimeRange } from '$lib/components/custom/date-timestamp-range-picker';
	import ComponentLoading from '$lib/components/otterscale/ui/component-loading.svelte';
	import { formatLatency, formatTime } from '$lib/formatter';
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

	const step = 1 * 120;
	const ioTimeQuery = $derived(
		`
		sum by (instance) (rate(node_disk_io_time_seconds_total{instance="juju-1eb21e-0-lxd-1", device=~"(/dev/)?(mmcblk.p.+|nvme.+|rbd.+|sd.+|vd.+|xvd.+|dm-.+|dasd.+)"}[5m]))
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
			const ioTimeResponse = await fetch(ioTimeQuery);
			serieses.set('io time', ioTimeResponse);

			mounted = true;
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

{#if mounted}
	{@const data = integrateSerieses(serieses)}
	<Template.Area title="Disk I/O">
		{#snippet hint()}
			<p>Historical view of disk I/O over the past hour</p>
		{/snippet}
		{#snippet content()}
			{#if data.length === 0}
				<Empty.Area />
			{:else}
				<div class="h-[200px] w-full resize overflow-visible">
					<AreaChart
						{data}
						x="time"
						series={[{ key: 'io time', color: 'hsl(var(--color-primary))' }]}
						legend={{
							classes: { root: '-mb-[50px] w-full overflow-auto' }
						}}
						props={{
							tooltip: {
								root: { class: 'bg-white/60 p-3 rounded shadow-lg' },
								header: { class: 'font-light' },
								item: {
									format: (v: number) => {
										const time = formatLatency(v);
										return `${time.value} ${time.unit}`;
									}
								}
							},
							xAxis: { format: formatTime },
							yAxis: {
								format: (v: number) => {
									const time = formatLatency(v);
									return `${Number(time.value).toFixed(0)} ${time.unit}`;
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
