<script lang="ts">
	import type { Scope } from '$lib/api/scope/v1/scope_pb';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Card from '$lib/components/ui/card';
	import * as Chart from '$lib/components/ui/chart/index.js';
	import { formatCapacity } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';
	import { scaleBand } from 'd3-scale';
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

	let receivesByTime = $state([] as SampleValue[]);
	let transmitsByTime = $state([] as SampleValue[]);
	const trafficsByTime = $derived(
		receivesByTime.map((sample, index) => ({
			time: sample.time,
			receive: sample.value,
			transmit: transmitsByTime[index]?.value ?? 0
		}))
	);
	let trafficsByTimeContext = $state<ChartContextValue>();
	const trafficsByTimeConfiguration = {
		receive: { label: 'Receive', color: 'var(--chart-1)' },
		transmit: { label: 'Transmit', color: 'var(--chart-2)' }
	} satisfies Chart.ChartConfig;

	function fetch() {
		prometheusDriver
			.rangeQuery(
				`sum(increase(node_network_receive_bytes_total{juju_model_uuid="${scope.uuid}"}[1h]))`,
				new Date().setHours(0, 0, 0, 0) - 24 * 60 * 60 * 1000,
				new Date().setHours(0, 0, 0, 0) + 24 * 60 * 60 * 1000,
				1 * 60 * 60
			)
			.then((response) => {
				receivesByTime = response.result[0]?.values;
			});
		prometheusDriver
			.rangeQuery(
				`sum(increase(node_network_transmit_bytes_total{juju_model_uuid="${scope.uuid}"}[1h]))`,
				new Date().setHours(0, 0, 0, 0) - 24 * 60 * 60 * 1000,
				new Date().setHours(0, 0, 0, 0) + 24 * 60 * 60 * 1000,
				1 * 60 * 60
			)
			.then((response) => {
				transmitsByTime = response.result[0]?.values;
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
		<Card.Header>
			<Card.Title>{m.total_upload_and_download()}</Card.Title>
			<Card.Description>
				<p class="lowercase">{m.over_each_time({ unit: m.hour() })}</p>
			</Card.Description>
		</Card.Header>
		<Card.Content>
			<Chart.Container config={trafficsByTimeConfiguration} class="h-[200px] w-full">
				<BarChart
					bind:context={trafficsByTimeContext}
					data={trafficsByTime}
					xScale={scaleBand().padding(0.25)}
					x="time"
					axis="x"
					rule={false}
					series={[
						{
							key: 'receive',
							label: trafficsByTimeConfiguration.receive.label,
							color: trafficsByTimeConfiguration.receive.color,
							props: { rounded: 'bottom' }
						},
						{
							key: 'transmit',
							label: trafficsByTimeConfiguration.transmit.label,
							color: trafficsByTimeConfiguration.transmit.color
						}
					]}
					seriesLayout="stack"
					props={{
						bars: {
							stroke: 'none',
							initialY: trafficsByTimeContext?.height,
							initialHeight: 0,
							motion: {
								y: { type: 'tween', duration: 500, easing: cubicInOut },
								height: { type: 'tween', duration: 500, easing: cubicInOut }
							}
						},
						highlight: { area: false },
						xAxis: {
							format: (v: Date) =>
								`${v.getHours().toString().padStart(2, '0')}:${v.getMinutes().toString().padStart(2, '0')}`,
							ticks: 1
						}
					}}
					legend
				>
					{#snippet belowMarks()}
						<Highlight area={{ class: 'fill-muted' }} />
					{/snippet}
					{#snippet tooltip()}
						<Chart.Tooltip
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
								{@const { value: io, unit } = formatCapacity(Number(value))}
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
