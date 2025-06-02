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

	let renderContext: 'svg' | 'canvas' = 'svg';
	let debug = false;

	let {
		client,
		scope: scope,
		instance: instance,
		timeRange
	}: { client: PrometheusDriver; scope: Scope; instance: string; timeRange: TimeRange } = $props();

	const step = 1 * 60;

	const systemQuery = $derived(
		`
        sum by (instance) (
            irate(
            node_cpu_seconds_total{instance="${instance}",juju_model_uuid=~"${scope.uuid}",mode="system"}[4m]
            )
        )
        / on (instance) group_left ()
        sum by (instance) (
            (irate(node_cpu_seconds_total{instance="${instance}",juju_model_uuid=~"${scope.uuid}"}[4m]))
        )
		`
	);
	const userQuery = $derived(
		`
        sum by (instance) (
            irate(
            node_cpu_seconds_total{instance="${instance}",juju_model_uuid=~"${scope.uuid}",mode="user"}[4m]
            )
        )
        / on (instance) group_left ()
        sum by (instance) (
            (irate(node_cpu_seconds_total{instance="${instance}",juju_model_uuid=~"${scope.uuid}"}[4m]))
        )
		`
	);
	const iowaitQuery = $derived(
		`
        sum by (instance) (
            irate(
            node_cpu_seconds_total{instance="${instance}",juju_model_uuid=~"${scope.uuid}",mode="iowait"}[4m]
            )
        )
        / on (instance) group_left ()
        sum by (instance) (
            (irate(node_cpu_seconds_total{instance="${instance}",juju_model_uuid=~"${scope.uuid}"}[4m]))
        )
		`
	);
	const irqsQuery = $derived(
		`
        sum by (instance) (
            irate(
            node_cpu_seconds_total{instance="${instance}",juju_model_uuid=~"${scope.uuid}",mode=~".*irq"}[4m]
            )
        )
        / on (instance) group_left ()
        sum by (instance) (
            (irate(node_cpu_seconds_total{instance="${instance}",juju_model_uuid=~"${scope.uuid}"}[4m]))
        )
		`
	);
	const idleQuery = $derived(
		`
        sum by (instance) (
            irate(
            node_cpu_seconds_total{instance="${instance}",juju_model_uuid=~"${scope.uuid}",mode="idle"}[4m]
            )
        )
        / on (instance) group_left ()
        sum by (instance) (
            (irate(node_cpu_seconds_total{instance="${instance}",juju_model_uuid=~"${scope.uuid}"}[4m]))
        )
		`
	);
	const otherQuery = $derived(
		`
		sum by (instance) (
			irate(
			node_cpu_seconds_total{instance="${instance}",juju_model_uuid=~"${scope.uuid}",mode!="idle",mode!="iowait",mode!="irq",mode!="softirq",mode!="system",mode!="user"}[4m]
			)
		)
		/ on (instance) group_left ()
		sum by (instance) (
			(irate(node_cpu_seconds_total{instance="${instance}",juju_model_uuid=~"${scope.uuid}"}[4m]))
		)
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
			const systemResponse = await fetch(systemQuery);
			serieses.set('system', systemResponse);

			const userResponse = await fetch(userQuery);
			serieses.set('user', userResponse);

			const iowaitResponse = await fetch(iowaitQuery);
			serieses.set('iowait', iowaitResponse);

			const irqseResponse = await fetch(irqsQuery);
			serieses.set('irqs', irqseResponse);

			const idleResponse = await fetch(idleQuery);
			serieses.set('idle', idleResponse);

			const otherResponse = await fetch(otherQuery);
			serieses.set('other', otherResponse);

			mounted = true;
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

{#if mounted}
	{@const data = integrateSerieses(serieses)}
	<Template.Area title="CPU">
		{#snippet hint()}
			<p>Basic CPU Information</p>
		{/snippet}
		{#snippet content()}
			{#if data.length === 0}
				<NoData type="area" />
			{:else}
				<div class="h-[200px] w-full resize overflow-visible">
					<AreaChart
						{data}
						x="time"
						yDomain={[0, 1]}
						series={[
							{ key: 'system', color: 'hsl(var(--color-primary))' },
							{ key: 'user', color: 'hsl(var(--color-secondary))' },
							{ key: 'iowait', color: 'hsl(var(--color-info))' },
							{ key: 'irqs', color: 'hsl(var(--color-success))' },
							{ key: 'idle', color: 'hsl(var(--color-warning))' },
							{ key: 'other', color: 'hsl(var(--color-danger))' }
						]}
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
