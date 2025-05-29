<script lang="ts">
	import ComponentLoading from '$lib/components/otterscale/ui/component-loading.svelte';
	import { PrometheusDriver, SampleValue } from 'prometheus-query';
	import { AreaChart } from 'layerchart';
	import * as Card from '$lib/components/ui/card';
	import Icon from '@iconify/svelte';
	import { onMount } from 'svelte';
	import { integrateSerieses } from '../utils';
	import * as HoverCard from '$lib/components/ui/hover-card/index.js';
	import { Button } from '$lib/components/ui/button';
	import NoData from '../utils/empty.svelte';
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import type { TimeRange } from '$lib/components/custom/date-timestamp-range-picker';
	import { fetchRange } from '../utils';
	import { formatCapacity } from '$lib/formatter';

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
		ceph_cluster_total_bytes{juju_model_uuid=~"${scope.uuid}"}
		`
	);
	const usedQuery = $derived(
		`
        ceph_cluster_total_used_bytes{juju_model_uuid=~"${scope.uuid}"}
		`
	);

	let serieses: Map<string, SampleValue[] | undefined> = $state(new Map());

	let mounted = $state(false);
	onMount(async () => {
		try {
			const totalResponse = await fetchRange(client, timeRange, step, totalQuery);
			serieses.set('total', totalResponse);

			const usedResponse = await fetchRange(client, timeRange, step, usedQuery);
			serieses.set('used', usedResponse);

			console.log(serieses);

			mounted = true;
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

{#if mounted}
	{@const data = integrateSerieses(serieses)}

	<Card.Root class="col-span-1 h-full w-full border-none shadow-none">
		<Card.Header class="h-[100px]">
			<Card.Title class="flex">
				<h1 class="text-3xl">Capacity</h1>
			</Card.Title>
			<Card.Description></Card.Description>
		</Card.Header>
		<Card.Content class="h-[200px]">
			{#if data.length === 0}
				<NoData type="area" />
			{:else}
				<div class="h-[200px] w-full resize overflow-visible">
					<AreaChart
						{data}
						x="time"
						series={[
							{ key: 'total', color: 'hsl(var(--color-primary))' },
							{ key: 'used', color: 'hsl(var(--color-secondary))' }
						]}
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
		</Card.Content>
		<Card.Footer class="h-[150px]"></Card.Footer>
	</Card.Root>
{:else}
	<ComponentLoading />
{/if}
