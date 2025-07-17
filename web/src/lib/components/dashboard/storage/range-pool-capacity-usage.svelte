<script lang="ts">
	import ComponentLoading from '$lib/components/otterscale/ui/component-loading.svelte';
	import { PrometheusDriver, SampleValue } from 'prometheus-query';
	import { AreaChart } from 'layerchart';
	import { onMount } from 'svelte';
	import { integrateSerieses } from '../utils';
	import * as Empty from '../utils/empty';
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import type { TimeRange } from '$lib/components/custom/date-timestamp-range-picker';
	import { fetchBatchRange } from '../utils';
	import { formatCapacity } from '$lib/formatter';
	import * as Template from '../utils/templates';

	let renderContext: 'svg' | 'canvas' = 'svg';
	let debug = false;

	let {
		client,
		scope: scope,
		timeRange
	}: { client: PrometheusDriver; scope: Scope; timeRange: TimeRange } = $props();

	const step = 1 * 60;

	const query = $derived(
		`
		ceph_pool_bytes_used{juju_model_uuid=~"${scope.uuid}"}
		* on (pool_id) group_right ()
		ceph_pool_metadata{juju_model_uuid=~"${scope.uuid}"}
		`
	);

	let serieses: Map<string, SampleValue[] | undefined> = $state(new Map());

	let mounted = $state(false);
	onMount(async () => {
		try {
			const batchResponses = await fetchBatchRange(client, timeRange, step, query);

			if (batchResponses) {
				for (const [metric, smapleSpace] of batchResponses) {
					const labels = metric.labels as { name: string };
					serieses.set(labels.name, smapleSpace);
				}
			}

			mounted = true;
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

{#if mounted}
	{@const data = integrateSerieses(serieses)}
	<Template.Area title="Capacity">
		{#snippet hint()}
			<p>
				Historical view of capacity usage, to help identify growth and trends in pool consumption
			</p>
		{/snippet}
		{#snippet description()}
			<p class="text-xl">Pool</p>
		{/snippet}
		{#snippet content()}
			{#if data.length === 0}
				<Empty.Area />
			{:else}
				<div class="h-[200px] w-full resize overflow-visible">
					<AreaChart
						{data}
						x="time"
						series={Array.from(serieses.keys()).map((key) => {
							return {
								key: key,
								color: `hsl(${Math.random() * 360}, 60%, 40%)`
							};
						})}
						legend={{
							classes: { root: '-mb-[50px] w-full overflow-auto' }
						}}
						props={{
							yAxis: {
								format: (v: number) => {
									const capacity = formatCapacity(v / 1024 / 1024);
									return `${capacity.value} ${capacity.unit}`;
								}
							},
							tooltip: {
								root: { class: 'bg-white/60 p-3 rounded shadow-lg' },
								header: { class: 'font-light' },
								item: {
									format: (v: number) => {
										const capacity = formatCapacity(v / 1024 / 1024);
										return `${capacity.value} ${capacity.unit}`;
									}
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
