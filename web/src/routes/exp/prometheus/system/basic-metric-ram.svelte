<script lang="ts">
	import { PrometheusDriver, SampleValue } from 'prometheus-query';
	import { AreaChart } from 'layerchart';
	import * as Card from '$lib/components/ui/card';
	import Icon from '@iconify/svelte';
	import { onMount } from 'svelte';
	import { integrateSerieses } from '../utils';
	import { formatCapacity } from '$lib/formatter';
	import ComponentLoading from '$lib/components/otterscale/ui/component-loading.svelte';
	import * as HoverCard from '$lib/components/ui/hover-card/index.js';
	import { Button } from '$lib/components/ui/button';
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import NoData from '../utils/empty.svelte';
	import type { TimeRange } from '$lib/components/custom/date-timestamp-range-picker';

	let renderContext: 'svg' | 'canvas' = 'svg';
	let debug = false;

	let {
		client,
		scope: scope,
		instance: instance,
		timeRange
	}: { client: PrometheusDriver; scope: Scope; instance: string; timeRange: TimeRange } = $props();

	const step = 1 * 60;

	const totalQuery = $derived(
		`
		node_memory_MemTotal_bytes{instance="${instance}",juju_model_uuid=~"${scope.uuid}"}
		`
	);
	const usedQuery = $derived(
		`
			node_memory_MemTotal_bytes{instance="${instance}",juju_model_uuid=~"${scope.uuid}"}
		-
			node_memory_MemFree_bytes{instance="${instance}",juju_model_uuid=~"${scope.uuid}"}
		-
		(
				node_memory_Cached_bytes{instance="${instance}",juju_model_uuid=~"${scope.uuid}"}
			+
				node_memory_Buffers_bytes{instance="${instance}",juju_model_uuid=~"${scope.uuid}"}
			+
			node_memory_SReclaimable_bytes{instance="${instance}",juju_model_uuid=~"${scope.uuid}"}
		)
		`
	);
	const cacheAndBufferQuery = $derived(
		`
			node_memory_Cached_bytes{instance="${instance}",juju_model_uuid=~"${scope.uuid}"}
		+
			node_memory_Buffers_bytes{instance="${instance}",juju_model_uuid=~"${scope.uuid}"}
		+
		node_memory_SReclaimable_bytes{instance="${instance}",juju_model_uuid=~"${scope.uuid}"}
		`
	);
	const freeQuery = $derived(
		`
		node_memory_MemFree_bytes{instance="${instance}",juju_model_uuid=~"${scope.uuid}"}
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

	<Card.Root class="col-span-1 h-full w-full border-none shadow-none">
		<Card.Header class="h-[100px]">
			<Card.Title class="flex">
				<h1 class="text-3xl">RAM</h1>
				<HoverCard.Root>
					<HoverCard.Trigger>
						<Button variant="ghost" size="icon" class="hover:bg-muted">
							<Icon icon="ph:info" />
						</Button>
					</HoverCard.Trigger>
					<HoverCard.Content class="w-fit max-w-[38w] text-xs text-muted-foreground">
						Basic Memory Information
					</HoverCard.Content>
				</HoverCard.Root>
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
		</Card.Content>
		<Card.Footer class="h-[150px]"></Card.Footer>
	</Card.Root>
{:else}
	<ComponentLoading />
{/if}
