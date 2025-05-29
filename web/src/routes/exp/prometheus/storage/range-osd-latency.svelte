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
	import { formatCapacity, formatNetworkIO } from '$lib/formatter';

	let renderContext: 'svg' | 'canvas' = 'svg';
	let debug = false;

	let {
		client,
		scope: scope,
		timeRange
	}: { client: PrometheusDriver; scope: Scope; timeRange: TimeRange } = $props();

	const step = 1 * 60;

	const applyQuery = $derived(
		`
		avg(ceph_osd_apply_latency_ms{juju_model_uuid=~"${scope.uuid}"})
		`
	);
	const commitQuery = $derived(
		`
		avg(ceph_osd_commit_latency_ms{juju_model_uuid=~"${scope.uuid}"})
		`
	);

	let serieses: Map<string, SampleValue[] | undefined> = $state(new Map());

	let mounted = $state(false);
	onMount(async () => {
		try {
			const applyResponse = await fetchRange(client, timeRange, step, applyQuery);
			serieses.set('apply', applyResponse);

			const commitResponse = await fetchRange(client, timeRange, step, commitQuery);
			serieses.set('commit', commitResponse);

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
				<h1 class="text-3xl">Latency</h1>
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
							{ key: 'apply', color: 'hsl(var(--color-primary))' },
							{ key: 'commit', color: 'hsl(var(--color-secondary))' }
						]}
						legend={{
							classes: { root: '-mb-[50px] w-full overflow-auto' }
						}}
						props={{
							yAxis: {
								format: (v: number) => {
									return `${v} ms`;
								}
							},
							tooltip: {
								root: { class: 'bg-white/60 p-3 rounded shadow-lg' },
								header: { class: 'font-light' },
								item: {
									format: (v: number) => {
										return `${v} ms`;
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
