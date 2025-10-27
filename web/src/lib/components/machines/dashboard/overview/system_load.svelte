<script lang="ts">
	import { scaleUtc } from 'd3-scale';
	import { curveNatural } from 'd3-shape';
	import { Area, AreaChart, LinearGradient } from 'layerchart';
	import { PrometheusDriver, SampleValue } from 'prometheus-query';
	import { onDestroy, onMount } from 'svelte';

	import type { Scope } from '$lib/api/scope/v1/scope_pb';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Card from '$lib/components/ui/card';
	import * as Chart from '$lib/components/ui/chart';
	import { m } from '$lib/paraglide/messages';

	let {
		prometheusDriver,
		scope,
		isReloading = $bindable(),
	}: { prometheusDriver: PrometheusDriver; scope: Scope; isReloading: boolean } = $props();

	const systemLoadConfiguration = {
		one: { label: '1 min', color: 'var(--chart-1)' },
		five: { label: '5 min', color: 'var(--chart-2)' },
	} satisfies Chart.ChartConfig;

	let ones = $state([] as SampleValue[]);
	let fives = $state([] as SampleValue[]);
	const systemLoads = $derived(
		ones.map((sample, index) => ({
			time: sample.time,
			one: sample.value,
			five: fives[index]?.value ?? 0,
		})),
	);

	async function fetch() {
		prometheusDriver
			.rangeQuery(
				`sum(node_load1{juju_model_uuid="${scope.uuid}"})`,
				Date.now() - 24 * 60 * 60 * 1000,
				Date.now(),
				2 * 60,
			)
			.then((response) => {
				ones = response.result[0]?.values;
			});
		prometheusDriver
			.rangeQuery(
				`sum(node_load5{juju_model_uuid="${scope.uuid}"})`,
				Date.now() - 24 * 60 * 60 * 1000,
				Date.now(),
				2 * 60,
			)
			.then((response) => {
				fives = response.result[0]?.values;
			});
	}

	const reloadManager = new ReloadManager(fetch);

	let isLoading = $state(true);
	onMount(async () => {
		await fetch();
		isLoading = false;
	});
	onDestroy(() => {
		reloadManager.stop();
	});

	$effect(() => {
		if (isReloading) {
			reloadManager.restart();
		} else {
			reloadManager.stop();
		}
	});
</script>

{#if isLoading}
	Loading
{:else}
	<Card.Root class="h-full">
		<Card.Header>
			<Card.Title>{m.system_load()}</Card.Title>
			<Card.Description>
				{m.machine_dashboard_system_loads_tooltip()}
			</Card.Description>
		</Card.Header>
		<Card.Content>
			<Chart.Container config={systemLoadConfiguration} class="h-[200px] w-full">
				<AreaChart
					data={systemLoads}
					x="time"
					xScale={scaleUtc()}
					yPadding={[0, 25]}
					series={[
						{
							key: 'one',
							label: systemLoadConfiguration.one.label,
							color: systemLoadConfiguration.one.color,
						},
						{
							key: 'five',
							label: systemLoadConfiguration.five.label,
							color: systemLoadConfiguration.five.color,
						},
					]}
					props={{
						area: {
							curve: curveNatural,
							'fill-opacity': 0.4,
							line: { class: 'stroke-1' },
							motion: 'tween',
						},
						xAxis: {
							format: (v: Date) => v.toLocaleDateString('en-US', { month: 'short' }),
						},
						yAxis: { format: () => '' },
					}}
				>
					{#snippet tooltip()}
						<Chart.Tooltip
							indicator="dot"
							labelFormatter={(v: Date) => {
								return v.toLocaleDateString('en-US', {
									year: 'numeric',
									month: 'short',
									day: 'numeric',
									hour: 'numeric',
									minute: 'numeric',
								});
							}}
						/>
					{/snippet}
					{#snippet marks({ series, getAreaProps })}
						{#each series as s, i (s.key)}
							<LinearGradient
								stops={[s.color ?? '', 'color-mix(in lch, ' + s.color + ' 10%, transparent)']}
								vertical
							>
								{#snippet children({ gradient })}
									<Area {...getAreaProps(s, i)} fill={gradient} />
								{/snippet}
							</LinearGradient>
						{/each}
					{/snippet}
				</AreaChart>
			</Chart.Container>
		</Card.Content>
	</Card.Root>
{/if}
