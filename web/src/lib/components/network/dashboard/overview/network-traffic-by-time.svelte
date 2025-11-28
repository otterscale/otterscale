<script lang="ts">
	import Icon from '@iconify/svelte';
	import { scaleBand } from 'd3-scale';
	import { BarChart, type ChartContextValue, Highlight } from 'layerchart';
	import { PrometheusDriver, SampleValue } from 'prometheus-query';
	import { onDestroy, onMount } from 'svelte';
	import { cubicInOut } from 'svelte/easing';
	import { SvelteDate } from 'svelte/reactivity';

	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Card from '$lib/components/ui/card';
	import * as Chart from '$lib/components/ui/chart/index.js';
	import { formatCapacity } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';

	let {
		prometheusDriver,
		isReloading = $bindable()
	}: {
		prometheusDriver: PrometheusDriver;
		isReloading: boolean;
	} = $props();

	let receivesByTime = $state([] as SampleValue[]);
	let transmitsByTime = $state([] as SampleValue[]);
	const trafficsByTime = $derived(
		receivesByTime?.map((sample, index) => ({
			time: sample.time,
			receive: sample.value,
			transmit: transmitsByTime?.[index]?.value ?? 0
		})) ?? []
	);
	let trafficsByTimeContext = $state<ChartContextValue>();
	const trafficsByTimeConfiguration = {
		receive: { label: 'Receive', color: 'var(--chart-1)' },
		transmit: { label: 'Transmit', color: 'var(--chart-2)' }
	} satisfies Chart.ChartConfig;

	let isLoaded = $state(false);

	async function fetchReceive() {
		const response = await prometheusDriver.rangeQuery(
			`sum(increase(node_network_receive_bytes_total[1h]))`,
			new SvelteDate().setHours(0, 0, 0, 0) - 24 * 60 * 60 * 1000,
			new SvelteDate().setHours(0, 0, 0, 0) + 24 * 60 * 60 * 1000,
			1 * 60 * 60
		);
		receivesByTime = response.result[0]?.values ?? [];
	}

	async function fetchTransmit() {
		const response = await prometheusDriver.rangeQuery(
			`sum(increase(node_network_transmit_bytes_total[1h]))`,
			new SvelteDate().setHours(0, 0, 0, 0) - 24 * 60 * 60 * 1000,
			new SvelteDate().setHours(0, 0, 0, 0) + 24 * 60 * 60 * 1000,
			1 * 60 * 60
		);
		transmitsByTime = response.result[0]?.values ?? [];
	}

	async function fetch() {
		try {
			await Promise.all([fetchReceive(), fetchTransmit()]);
			isLoaded = true;
		} catch (error) {
			console.error('Failed to fetch network traffic data:', error);
		}
	}

	const reloadManager = new ReloadManager(fetch);

	onMount(async () => {
		await fetch();
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

<Card.Root class="h-full gap-2">
	<Card.Header>
		<Card.Title>{m.total_upload_and_download()}</Card.Title>
		<Card.Description>
			<p class="lowercase">{m.over_each_time({ unit: m.hour() })}</p>
		</Card.Description>
	</Card.Header>
	{#if !isLoaded}
		<div class="flex h-[200px] w-full items-center justify-center">
			<Icon icon="svg-spinners:6-dots-rotate" class="m-4 size-12" />
		</div>
	{:else}
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
	{/if}
</Card.Root>
