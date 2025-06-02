<script lang="ts">
	import ComponentLoading from '$lib/components/otterscale/ui/component-loading.svelte';
	import { PrometheusDriver, SampleValue } from 'prometheus-query';
	import { AreaChart } from 'layerchart';
	import { onMount } from 'svelte';
	import { integrateSerieses } from '../utils';
	import NoData from '../utils/empty.svelte';
	import type { Scope } from '$gen/api/scope/v1/scope_pb';
	import type { TimeRange } from '$lib/components/custom/date-timestamp-range-picker';
	import { fetchRange } from '../utils';
	import * as Template from '../utils/templates';

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

			mounted = true;
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

{#if mounted}
	{@const data = integrateSerieses(serieses)}
	<Template.Area title="OSD Latency">
		{#snippet content()}
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
		{/snippet}
	</Template.Area>
{:else}
	<ComponentLoading />
{/if}
