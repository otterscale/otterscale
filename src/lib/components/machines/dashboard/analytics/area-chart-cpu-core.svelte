<script lang="ts">
	import Icon from '@iconify/svelte';
	import { scaleUtc } from 'd3-scale';
	import { curveMonotoneX } from 'd3-shape';
	import { Area, AreaChart, ChartClipPath } from 'layerchart';
	import { PrometheusDriver } from 'prometheus-query';
	import { cubicInOut } from 'svelte/easing';

	import {
		formatTimeRange,
		formatXAxisDate,
		getXAxisTicks
	} from '$lib/components/custom/chart/units/formatter';
	import * as Statistics from '$lib/components/custom/data-table/statistics/index';
	import {
		type DataPoint,
		fetchFlattenedRange,
		generateChartConfig,
		getSeries
	} from '$lib/components/custom/prometheus';
	import ChartContainer from '$lib/components/ui/chart/chart-container.svelte';
	import * as Chart from '$lib/components/ui/chart/index.js';
	import Input from '$lib/components/ui/input/input.svelte';
	import { m } from '$lib/paraglide/messages';
	import { cn } from '$lib/utils';

	let { client, fqdn }: { client: PrometheusDriver; fqdn: string } = $props();

	// Configuration
	const STEP_SECONDS = 60;
	const TIME_RANGE_HOURS = 1;
	const TOP_HIGHLIGHT_COUNT = 3;

	// Time range
	const endTime = new Date();
	const startTime = new Date(endTime.getTime() - TIME_RANGE_HOURS * 60 * 60 * 1000);

	// Prometheus query
	const query = $derived(
		`sum by (cpu) (rate(node_cpu_seconds_total{instance=~"${fqdn}", mode!="idle"}[5m])) * 100`
	);

	// Filter state
	let topk = $state(10);

	/**
	 * Calculate total usage for each CPU core across the entire time range
	 */
	function calculateTopKSeries(data: DataPoint[], k: number): string[] {
		const config = generateChartConfig(data);
		const seriesKeys = Object.keys(config);

		return seriesKeys
			.map((key) => ({
				key,
				total: data.reduce((sum, datum) => sum + (Number(datum[key]) || 0), 0)
			}))
			.sort((a, b) => b.total - a.total)
			.slice(0, k)
			.map((item) => item.key);
	}

	/**
	 * Filter data to only include top K CPU cores
	 */
	function filterDataByTopK(data: DataPoint[], topKSeries: string[]): DataPoint[] {
		return data.map((datum) => {
			const filtered: DataPoint = { date: datum.date };
			topKSeries.forEach((key) => {
				filtered[key] = datum[key];
			});
			return filtered;
		});
	}

	/**
	 * Check if a CPU core is in the top 3
	 */
	function isTopSeries(name: string, topSeries: string[]): boolean {
		return topSeries.includes(name);
	}
</script>

<Statistics.Root type="count" class="overflow-visible">
	<Statistics.Header>
		<div class="flex">
			<Statistics.Title class="h-8 **:data-[slot=data-table-statistics-title-icon]:size-6">
				<div class="flex flex-col justify-between">
					{m.cpu()}
					<p class="text-sm text-muted-foreground">{m.core()}/{m.processor()}</p>
				</div>
			</Statistics.Title>
			<div class="relative ml-auto">
				<span class="absolute top-1/2 left-3 -translate-y-1/2 items-center">
					<Icon icon="ph:funnel-duotone" />
				</span>
				<Input type="number" bind:value={topk} min={0} step={5} class="h-8 w-22 pl-9 text-lg" />
			</div>
		</div>
	</Statistics.Header>

	<Statistics.Content class="min-h-16">
		{#await fetchFlattenedRange(client, query, startTime, endTime, STEP_SECONDS)}
			<!-- Loading state -->
			<div class="flex h-[250px] w-full items-center justify-center">
				<Icon icon="svg-spinners:blocks-wave" class="m-8 size-32 text-muted-foreground/50" />
			</div>
		{:then rawData}
			{#if rawData.length === 0}
				<!-- Empty state -->
				<div class="flex h-[250px] w-full flex-col items-center justify-center">
					<Icon icon="ph:chart-line-fill" class="size-60 animate-pulse text-muted-foreground" />
					<p class="text-base text-muted-foreground">{m.no_data_display()}</p>
				</div>
			{:else}
				<!-- Calculate top K CPU cores based on total usage -->
				{@const topKSeries = calculateTopKSeries(rawData, topk)}
				{@const filteredData = filterDataByTopK(rawData, topKSeries)}
				{@const chartConfig = generateChartConfig(filteredData)}
				{@const top3Series = topKSeries.slice(0, TOP_HIGHLIGHT_COUNT)}

				<ChartContainer config={chartConfig} class="aspect-auto h-[250px] w-full">
					<AreaChart
						data={filteredData}
						x="date"
						xScale={scaleUtc()}
						series={getSeries(chartConfig)}
						seriesLayout="stack"
						props={{
							area: {
								curve: curveMonotoneX,
								'fill-opacity': 0.4,
								line: { class: 'stroke-1' },
								motion: 'tween'
							},
							xAxis: {
								ticks: getXAxisTicks(formatTimeRange(TIME_RANGE_HOURS)),
								format: (date: Date) => formatXAxisDate(date, formatTimeRange(TIME_RANGE_HOURS))
							},
							yAxis: {
								format: () => ''
							}
						}}
					>
						{#snippet marks({ series: chartSeries, getAreaProps })}
							<!-- Define gradients for each series -->
							<defs>
								{#each chartSeries as series (series.key)}
									{@const gradientId = series.key.replace(/\s+/g, '')}
									<linearGradient id="fill{gradientId}" x1="0" y1="0" x2="0" y2="1">
										<stop offset="5%" stop-color={series.color} stop-opacity={1.0} />
										<stop offset="95%" stop-color={series.color} stop-opacity={0.4} />
									</linearGradient>
								{/each}
							</defs>

							<!-- Animated chart areas -->
							<ChartClipPath
								initialWidth={0}
								motion={{
									width: {
										type: 'tween',
										duration: 1000,
										easing: cubicInOut
									}
								}}
							>
								{#each chartSeries as series, index (series.key)}
									{@const gradientId = series.key.replace(/\s+/g, '')}
									<Area {...getAreaProps(series, index)} fill="url(#fill{gradientId})" />
								{/each}
							</ChartClipPath>
						{/snippet}

						{#snippet tooltip()}
							<Chart.Tooltip
								labelFormatter={(v: Date) => {
									return v.toLocaleDateString('en-US', {
										year: 'numeric',
										month: 'long',
										day: 'numeric',
										hour: 'numeric',
										minute: 'numeric'
									});
								}}
							>
								{#snippet formatter({ name, value })}
									{@const isTop = isTopSeries(name, top3Series)}
									<div
										class="flex w-full shrink-0 items-center justify-between gap-4 leading-none"
										style="--color-bg: var(--color-{name})"
									>
										{#if value !== undefined && value !== null}
											<span class="flex w-full items-center gap-1">
												<Icon icon="ph:square-fill" class="text-(--color-bg)" />
												<p class={isTop ? 'fond-bold text-destructive' : 'text-foreground'}>
													{name}
												</p>
											</span>
											<p
												class={cn(
													'font-mono font-medium whitespace-nowrap tabular-nums',
													isTop ? 'fond-bold text-destructive' : 'text-foreground'
												)}
											>
												{Number(value).toFixed(2)}%
											</p>
										{/if}
									</div>
								{/snippet}
							</Chart.Tooltip>
						{/snippet}
					</AreaChart>
				</ChartContainer>
			{/if}
		{:catch}
			<!-- Error state -->
			<div class="flex h-[250px] w-full flex-col items-center justify-center">
				<Icon icon="ph:chart-bar-fill" class="size-24 animate-pulse text-muted-foreground" />
				<p class="text-base text-muted-foreground">{m.no_data_display()}</p>
			</div>
		{/await}
	</Statistics.Content>
</Statistics.Root>
