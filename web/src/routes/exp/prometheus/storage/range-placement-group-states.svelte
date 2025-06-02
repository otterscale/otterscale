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

	const totalQuery = $derived(
		`
		sum(ceph_pg_total{juju_model_uuid=~"${scope.uuid}"})
		`
	);
	const activeQuery = $derived(
		`
		sum(ceph_pg_active{juju_model_uuid=~"${scope.uuid}"})
		`
	);
	const inactiveQuery = $derived(
		`
		sum(
		ceph_pg_total{juju_model_uuid=~"${scope.uuid}"} - ceph_pg_active{juju_model_uuid=~"${scope.uuid}"}
		)

		`
	);
	const undersizedQuery = $derived(
		`
		sum(ceph_pg_undersized{juju_model_uuid=~"${scope.uuid}"})
		`
	);
	const degradedQuery = $derived(
		`
		sum(ceph_pg_degraded{juju_model_uuid=~"${scope.uuid}"})
		`
	);
	const inconsistentQuery = $derived(
		`
		sum(ceph_pg_inconsistent{juju_model_uuid=~"${scope.uuid}"})
		`
	);
	const downQuery = $derived(
		`
		sum(ceph_pg_down{juju_model_uuid=~"${scope.uuid}"})
		`
	);

	let serieses: Map<string, SampleValue[] | undefined> = $state(new Map());

	let mounted = $state(false);
	onMount(async () => {
		try {
			const totalResponse = await fetchRange(client, timeRange, step, totalQuery);
			serieses.set('total', totalResponse);

			const activeResponse = await fetchRange(client, timeRange, step, activeQuery);
			serieses.set('active', activeResponse);

			const inactiveResponse = await fetchRange(client, timeRange, step, inactiveQuery);
			serieses.set('inactive', inactiveResponse);

			const undersizedResponse = await fetchRange(client, timeRange, step, undersizedQuery);
			serieses.set('undersized', undersizedResponse);

			const degradedResponse = await fetchRange(client, timeRange, step, degradedQuery);
			serieses.set('degraded', degradedResponse);

			const inconsistentResponse = await fetchRange(client, timeRange, step, inconsistentQuery);
			serieses.set('inconsistent', inconsistentResponse);

			const downResponse = await fetchRange(client, timeRange, step, downQuery);
			serieses.set('down', downResponse);

			console.log(serieses);

			mounted = true;
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

{#if mounted}
	{@const data = integrateSerieses(serieses)}
	<Template.Area title="Placement Group States">
		{#snippet content()}
			{#if data.length === 0}
				<NoData type="area" />
			{:else}
				<div class="h-[200px] w-full resize overflow-visible">
					<AreaChart
						{data}
						x="time"
						series={[
							{ key: 'total', color: 'hsl(var(--color-primary))' },
							{ key: 'active', color: 'hsl(var(--color-secondary))' },
							{ key: 'inactive', color: 'hsl(var--color-neutral))' },
							{ key: 'undersized', color: 'hsl(var(--color-info))' },
							{ key: 'degraded', color: 'hsl(var(--color-success))' },
							{ key: 'inconsistent', color: 'hsl(var(--color-warning))' },
							{ key: 'down', color: 'hsl(var(--color-danger))' }
						]}
						legend={{
							classes: { root: '-mb-[50px] w-full overflow-auto' }
						}}
						props={{
							tooltip: {
								root: { class: 'bg-white/60 p-3 rounded shadow-lg' },
								header: { class: 'font-light' }
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
