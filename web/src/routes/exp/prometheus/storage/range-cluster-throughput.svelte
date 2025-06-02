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
	import * as Template from '../utils/templates';

	let renderContext: 'svg' | 'canvas' = 'svg';
	let debug = false;

	let {
		client,
		scope: scope,
		timeRange
	}: { client: PrometheusDriver; scope: Scope; timeRange: TimeRange } = $props();

	const step = 1 * 60;

	const writeQuery = $derived(
		`
		sum(irate(ceph_osd_op_w_in_bytes{juju_model_uuid=~"${scope.uuid}"}[5m]))
		`
	);
	const readQuery = $derived(
		`
		sum(irate(ceph_osd_op_r_out_bytes{juju_model_uuid=~"${scope.uuid}"}[5m]))
		`
	);

	let serieses: Map<string, SampleValue[] | undefined> = $state(new Map());

	let mounted = $state(false);
	onMount(async () => {
		try {
			const writeResponse = await fetchRange(client, timeRange, step, writeQuery);
			serieses.set('write', writeResponse);

			const readResponse = await fetchRange(client, timeRange, step, readQuery);
			serieses.set('read', readResponse);

			console.log(serieses);

			mounted = true;
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

{#if mounted}
	{@const data = integrateSerieses(serieses)}
	<Template.Area title="Throughput">
		{#snippet content()}
			{#if data.length === 0}
				<NoData type="area" />
			{:else}
				<div class="h-[200px] w-full resize overflow-visible">
					<AreaChart
						{data}
						x="time"
						series={[
							{ key: 'write', color: 'hsl(var(--color-primary))' },
							{ key: 'read', color: 'hsl(var(--color-secondary))' }
						]}
						legend={{
							classes: { root: '-mb-[50px] w-full overflow-auto' }
						}}
						props={{
							yAxis: {
								format: (v: number) => {
									const throughput = formatNetworkIO(v);
									return `${throughput.value} ${throughput.unit}`;
								}
							},
							tooltip: {
								root: { class: 'bg-white/60 p-3 rounded shadow-lg' },
								header: { class: 'font-light' },
								item: {
									format: (v: number) => {
										const throughput = formatNetworkIO(v);
										return `${throughput.value} ${throughput.unit}`;
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
