<script lang="ts">
	import { PrometheusDriver, SampleValue } from 'prometheus-query';
	import { AreaChart } from 'layerchart';
	import { onMount } from 'svelte';
	import { integrateSerieses } from '../utils';
	import { formatCapacity } from '$lib/formatter';
	import ComponentLoading from '$lib/components/otterscale/ui/component-loading.svelte';
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import * as Empty from '../utils/empty';
	import type { TimeRange } from '$lib/components/custom/date-timestamp-range-picker';
	import * as Template from '../utils/templates';

	let renderContext: 'svg' | 'canvas' = 'svg';
	let debug = false;

	let {
		client,
		scope: scope,
		timeRange
	}: { client: PrometheusDriver; scope: Scope; timeRange: TimeRange } = $props();

	const step = 1 * 60;

	const totalQuery = $derived(
		`
		node_memory_MemTotal_bytes{juju_model_uuid=~"${scope.uuid}"}
		`
	);
	const usedQuery = $derived(
		`
			node_memory_MemTotal_bytes{juju_model_uuid=~"${scope.uuid}"}
		-
			node_memory_MemFree_bytes{juju_model_uuid=~"${scope.uuid}"}
		-
		(
				node_memory_Cached_bytes{juju_model_uuid=~"${scope.uuid}"}
			+
				node_memory_Buffers_bytes{juju_model_uuid=~"${scope.uuid}"}
			+
			node_memory_SReclaimable_bytes{juju_model_uuid=~"${scope.uuid}"}
		)
		`
	);
	const cacheAndBufferQuery = $derived(
		`
			node_memory_Cached_bytes{juju_model_uuid=~"${scope.uuid}"}
		+
			node_memory_Buffers_bytes{juju_model_uuid=~"${scope.uuid}"}
		+
		node_memory_SReclaimable_bytes{juju_model_uuid=~"${scope.uuid}"}
		`
	);
	const freeQuery = $derived(
		`
		node_memory_MemFree_bytes{juju_model_uuid=~"${scope.uuid}"}
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
			const totalResponse = await fetch(totalQuery);
			serieses.set('total', totalResponse);

			const usedResponse = await fetch(usedQuery);
			serieses.set('used', usedResponse);

			const cacheAndBufferResponse = await fetch(cacheAndBufferQuery);
			serieses.set('cacheAndBuffer', cacheAndBufferResponse);

			const freeResponse = await fetch(freeQuery);
			serieses.set('free', freeResponse);

			mounted = true;
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

{#if mounted}
	{@const data = integrateSerieses(serieses)}
	<Template.Area title="RAM">
		{#snippet hint()}
			<p>Basic Memory Information</p>
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
							{ key: 'total', color: 'hsl(var(--color-primary))' },
							{ key: 'used', color: 'hsl(var(--color-secondary))' },
							{ key: 'cacheAndBuffer', label: 'cache + buffer', color: 'hsl(var(--color-info))' },
							{ key: 'free', color: 'hsl(var(--color-success))' }
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
							yAxis: {
								format: (v: number) => `${(v / 1024 / 1024 / 1024).toFixed(0)} GiB`
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
