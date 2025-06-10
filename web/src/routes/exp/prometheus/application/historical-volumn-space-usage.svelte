<script lang="ts">
	import ComponentLoading from '$lib/components/otterscale/ui/component-loading.svelte';
	import { PrometheusDriver, SampleValue } from 'prometheus-query';
	import { AreaChart } from 'layerchart';
	import { onMount } from 'svelte';
	import { integrateSerieses } from '../utils';
	import NoData from '../utils/empty.svelte';
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import type { TimeRange } from '$lib/components/custom/date-timestamp-range-picker';
	import * as Template from '../utils/templates';
	import { fetchRange } from '../utils';

	let renderContext: 'svg' | 'canvas' = 'svg';
	let debug = false;

	let {
		client,
		scope: scope,
		timeRange
	}: { client: PrometheusDriver; scope: Scope; timeRange: TimeRange } = $props();

	const step = 1 * 60;

	const usageQuery = $derived(
		`
		(
			sum without (instance, node) (
			topk(
				1,
				(
				kubelet_volume_stats_capacity_bytes{job="kubelet",juju_model_uuid=~"${scope.uuid}",metrics_path="/metrics"}
				)
			)
			)
		-
			sum without (instance, node) (
			topk(
				1,
				(
				kubelet_volume_stats_available_bytes{job="kubelet",juju_model_uuid=~"${scope.uuid}",metrics_path="/metrics"}
				)
			)
			)
		)
		`
	);

	let serieses: Map<string, SampleValue[] | undefined> = $state(new Map());

	let mounted = $state(false);
	onMount(async () => {
		try {
			const usageResponse = await fetchRange(client, timeRange, step, usageQuery);
			serieses.set('usage', usageResponse);

			mounted = true;
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

{#if mounted}
	{@const data = integrateSerieses(serieses)}
	<Template.Area title="Volumn Space">
		{#snippet content()}
			{#if data.length === 0}
				<NoData type="area" />
			{:else}
				<div class="h-[200px] w-full resize overflow-visible">
					<AreaChart
						{data}
						x="time"
						yDomain={[0, 1]}
						series={[{ key: 'usage', color: 'hsl(var(--color-primary))' }]}
						legend={{
							classes: { root: '-mb-[50px] w-full overflow-auto' }
						}}
						props={{
							yAxis: { format: 'percent' },
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
