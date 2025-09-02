<script lang="ts">
	import type { Scope } from '$lib/api/scope/v1/scope_pb';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Card from '$lib/components/ui/card';
	import * as Chart from '$lib/components/ui/chart/index.js';
	import { formatIO } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';
	import { scaleUtc } from 'd3-scale';
	import { BarChart, Highlight, type ChartContextValue } from 'layerchart';
	import { PrometheusDriver, SampleValue } from 'prometheus-query';
	import { onMount } from 'svelte';
	import { cubicInOut } from 'svelte/easing';

	let {
		prometheusDriver,
		scope,
		isReloading = $bindable()
	}: {
		prometheusDriver: PrometheusDriver;
		scope: Scope;
		isReloading: boolean;
	} = $props();

	let receives = $state([] as SampleValue[]);
	let transmits = $state([] as SampleValue[]);
	let latestReceive = $state({} as number);
	let latestTransmit = $state({} as number);
	let activeTraffic = $state<keyof typeof trafficsConfigurations>('receive');

	const traffics = $derived(
		receives.map((sample, index) => ({
			time: sample.time,
			receive: sample.value,
			transmit: transmits[index]?.value ?? 0
		}))
	);
	const latestTraffics = $derived({
		receive: latestReceive,
		transmit: latestTransmit
	});
	const trafficsConfigurations = {
		views: { label: 'Traffic', color: '' },
		receive: { label: 'Receive', color: 'var(--chart-1)' },
		transmit: { label: 'Transmit', color: 'var(--chart-2)' }
	} satisfies Chart.ChartConfig;

	let trafficsContext = $state<ChartContextValue>();

	const activeTrafficConfiguration = $derived([
		{
			key: activeTraffic,
			label: trafficsConfigurations[activeTraffic].label,
			color: trafficsConfigurations[activeTraffic].color
		}
	]);

	function fetch() {
		prometheusDriver
			.rangeQuery(
				`sum(irate(node_network_receive_bytes_total{juju_model_uuid="${scope.uuid}"}[4m]))`,
				new Date().setMinutes(0, 0, 0) - 24 * 60 * 60 * 1000,
				new Date().setMinutes(0, 0, 0),
				2 * 60
			)
			.then((response) => {
				receives = response.result[0].values;
			});
		prometheusDriver
			.rangeQuery(
				`sum(irate(node_network_transmit_bytes_total{juju_model_uuid="${scope.uuid}"}[4m]))`,
				new Date().setMinutes(0, 0, 0) - 24 * 60 * 60 * 1000,
				new Date().setMinutes(0, 0, 0),
				2 * 60
			)
			.then((response) => {
				transmits = response.result[0].values;
			});
		prometheusDriver
			.instantQuery(
				`sum(irate(node_network_receive_bytes_total{juju_model_uuid="${scope.uuid}"}[4m]))`
			)
			.then((response) => {
				latestReceive = response.result[0].value.value;
			});
		prometheusDriver
			.instantQuery(
				`sum(irate(node_network_transmit_bytes_total{juju_model_uuid="${scope.uuid}"}[4m]))`
			)
			.then((response) => {
				latestTransmit = response.result[0].value.value;
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
	<Card.Root class="h-full gap-2">
		<Card.Header class="flex flex-col items-stretch space-y-0 border-b p-0 sm:flex-row">
			<div class="flex flex-1 flex-col justify-center gap-1 px-6 py-5 sm:py-6">
				<Card.Title>{m.network_traffic()}</Card.Title>
			</div>
			<div class="flex">
				{#each ['receive', 'transmit'] as key (key)}
					{@const chart = key as keyof typeof trafficsConfigurations}
					{@const { value, unit } = formatIO(latestTraffics[key as keyof typeof latestTraffics])}
					<button
						data-active={activeTraffic === chart}
						class="data-[active=true]:bg-muted/50 relative z-30 flex flex-1 flex-col justify-center gap-1 border-t px-6 py-4 text-left even:border-l sm:border-t-0 sm:border-l sm:px-8 sm:py-6"
						onclick={() => (activeTraffic = chart)}
					>
						<span class="text-muted-foreground text-xs">
							{trafficsConfigurations[chart].label}
						</span>
						<span class="flex items-end gap-1 text-lg leading-none font-bold sm:text-3xl">
							{value.toLocaleString()}
							<span class="text-xs">{unit}</span>
						</span>
					</button>
				{/each}
			</div>
		</Card.Header>
		<Card.Content class="px-6 pt-6">
			<Chart.Container config={trafficsConfigurations} class="aspect-auto h-[120px] w-full">
				<BarChart
					bind:context={trafficsContext}
					data={traffics}
					x="time"
					axis="x"
					series={activeTrafficConfiguration}
					props={{
						bars: {
							stroke: 'none',
							rounded: 'none',
							// use the height of the chart to animate the bars
							initialY: trafficsContext?.height,
							initialHeight: 0,
							motion: {
								y: { type: 'tween', duration: 500, easing: cubicInOut },
								height: { type: 'tween', duration: 500, easing: cubicInOut }
							}
						},
						highlight: { area: { fill: 'none' } },
						xAxis: {
							format: (v: Date) =>
								`${v.getHours().toString().padStart(2, '0')}:${v.getMinutes().toString().padStart(2, '0')}`,
							ticks: (scale) => scaleUtc(scale.domain(), scale.range()).ticks()
						}
					}}
				>
					{#snippet belowMarks()}
						<Highlight area={{ class: 'fill-muted' }} />
					{/snippet}
					{#snippet tooltip()}
						<Chart.Tooltip
							nameKey="views"
							labelFormatter={(time: Date) => {
								return time.toLocaleDateString('en-US', {
									year: 'numeric',
									month: 'short',
									day: 'numeric',
									hour: 'numeric',
									minute: 'numeric'
								});
							}}
						>
							{#snippet formatter({ item, name, value })}
								{@const { value: io, unit } = formatIO(Number(value))}
								<div
									style="--color-bg: {item.color}"
									class="aspect-square h-full w-fit shrink-0 border-(--color-border) bg-(--color-bg)"
								></div>
								<div class="flex flex-1 shrink-0 items-center justify-between text-xs leading-none">
									<div class="grid gap-1.5">
										<span class="text-muted-foreground">{name}</span>
									</div>
									<p class="font-mono">{io} {unit}</p>
								</div>
							{/snippet}
						</Chart.Tooltip>
					{/snippet}
				</BarChart>
			</Chart.Container>
		</Card.Content>
	</Card.Root>
{/if}
