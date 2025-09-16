<script lang="ts">
	import { scaleUtc } from 'd3-scale';
	import { curveNatural } from 'd3-shape';
	import { Area, AreaChart, LinearGradient } from 'layerchart';
	import { PrometheusDriver, SampleValue } from 'prometheus-query';
	import { onMount } from 'svelte';

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

	const requestsConfiguration = {
		ninety_five: { label: '95', color: 'var(--chart-1)' },
		ninety_nine: { label: '99', color: 'var(--chart-2)' },
	} satisfies Chart.ChartConfig;

	let ninety_fives = $state([] as SampleValue[]);
	let ninety_nines = $state([] as SampleValue[]);
	const requests = $derived(
		ninety_fives.map((sample, index) => ({
			time: sample.time,
			running: sample.value,
			waiting: ninety_nines[index]?.value ?? 0,
		})),
	);

	async function fetch() {
		prometheusDriver
			.rangeQuery(
				`histogram_quantile(0.95, sum by(le) (rate(vllm:time_to_first_token_seconds_bucket{scope_uuid="${scope.uuid}"}[2m])))`,
				Date.now() - 24 * 60 * 60 * 1000,
				Date.now(),
				2 * 60,
			)
			.then((response) => {
				ninety_fives = response.result[0]?.values;
			});
		prometheusDriver
			.rangeQuery(
				`histogram_quantile(0.99, sum by(le) (rate(vllm:time_to_first_token_seconds_bucket{scope_uuid="${scope.uuid}"}[2m])))`,
				Date.now() - 24 * 60 * 60 * 1000,
				Date.now(),
				2 * 60,
			)
			.then((response) => {
				ninety_nines = response.result[0]?.values;
			});
	}

	const reloadManager = new ReloadManager(fetch);

	let isLoading = $state(true);
	onMount(async () => {
		await fetch();
		isLoading = false;
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
			<Card.Title>{m.time_to_first_token()}</Card.Title>
			<Card.Description>
				{m.llm_dashboard_time_to_first_token_tooltip()}
			</Card.Description>
		</Card.Header>
		<Card.Content>
			<Chart.Container config={requestsConfiguration} class="h-[200px] w-full">
				<AreaChart
					data={requests}
					x="time"
					xScale={scaleUtc()}
					yPadding={[0, 25]}
					series={[
						{
							key: 'ninety_five',
							label: requestsConfiguration.ninety_five.label,
							color: requestsConfiguration.ninety_five.color,
						},
						{
							key: 'ninety_nine',
							label: requestsConfiguration.ninety_nine.label,
							color: requestsConfiguration.ninety_nine.color,
						},
					]}
					seriesLayout="stack"
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
