<script lang="ts">
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import type { TimeRange } from '$lib/components/custom/date-timestamp-range-picker';
	import ComponentLoading from '$lib/components/otterscale/ui/component-loading.svelte';
	import { formatCapacity, formatTime } from '$lib/formatter';
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
	const readBytesQuery = $derived(
		`
		sum by (instance) (rate(node_disk_read_bytes_total{instance="juju-1eb21e-0-lxd-1", device=~"(/dev/)?(mmcblk.p.+|nvme.+|rbd.+|sd.+|vd.+|xvd.+|dm-.+|dasd.+)"}[5m]))
		`
	);
	const writeBytesQuery = $derived(
		`
		sum by (instance) (rate(node_disk_written_bytes_total{instance="juju-1eb21e-0-lxd-1", device=~"(/dev/)?(mmcblk.p.+|nvme.+|rbd.+|sd.+|vd.+|xvd.+|dm-.+|dasd.+)"}[5m]))
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
			const readBytesResponse = await fetch(readBytesQuery);
			serieses.set('read', readBytesResponse);

			const writeBytesResponse = await fetch(writeBytesQuery);
			serieses.set('write', writeBytesResponse);

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
						series={[
							{ key: 'read', color: 'hsl(var(--color-primary))' },
							{ key: 'write', color: 'hsl(var(--color-secondary))' }
						]}
						legend={{
							classes: { root: '-mb-[50px] w-full overflow-auto' }
						}}
						props={{
							tooltip: {
								root: { class: 'bg-white/60 p-3 rounded shadow-lg' },
								header: { class: 'font-light' },
								item: {
									format: (v: number) => {
										const capacity = formatCapacity(v / 1024 / 1024);
										return `${capacity.value} ${capacity.unit}`;
									}
								}
							},
							xAxis: { format: formatTime },
							yAxis: {
								format: (v: number) => {
									const capacity = formatCapacity(v / 1024 / 1024);
									return `${Number(capacity.value).toFixed(0)} ${capacity.unit}`;
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
