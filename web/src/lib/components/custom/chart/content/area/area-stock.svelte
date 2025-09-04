<script lang="ts">
	import { scaleUtc } from 'd3-scale';
	import { curveNatural } from 'd3-shape';
	import { Area, AreaChart, ChartClipPath } from 'layerchart';
	import { cubicInOut } from 'svelte/easing';

	import {
		formatXAxisDate,
		formatTooltipDate,
		getXAxisTicks,
		type TimeRange,
	} from '$lib/components/custom/chart/units/formatter';
	import {
		generateChartConfig,
		getSeries,
		type DataPoint,
		type ChartConfig,
	} from '$lib/components/custom/prometheus';
	import ChartContainer from '$lib/components/ui/chart/chart-container.svelte';
	import * as Chart from '$lib/components/ui/chart/index.js';

	// Constants
	const DEFAULT_CHART_HEIGHT = 'h-[250px]';
	const GRADIENT_OPACITY_START = 1.0;
	const GRADIENT_OPACITY_END = 0.1;
	const AREA_FILL_OPACITY = 0.4;
	const ANIMATION_DURATION = 1000;

	/**
	 * Component props interface
	 */
	interface Props {
		/** Chart data points with date and metric values */
		data: Array<DataPoint>;
		/** Optional chart configuration object */
		chartConfig?: ChartConfig;
		/** Time range for x-axis formatting (1h, 1d, 7d, 30d) */
		timeRange?: TimeRange;
		/** Additional CSS classes for the chart container */
		class?: string;
		/** Optional value formatter function that returns formatted value and unit */
		valueFormatter?: (value: number) => { value: number; unit: string };
	}

	let { data, chartConfig, timeRange, valueFormatter, class: className = '' }: Props = $props();

	// Derived reactive values
	const computedChartConfig = $derived(() => chartConfig ?? generateChartConfig(data));
	const series = $derived(() => getSeries(computedChartConfig()));
</script>

<ChartContainer config={computedChartConfig()} class="aspect-auto {DEFAULT_CHART_HEIGHT} w-full {className}">
	<AreaChart
		legend
		{data}
		x="date"
		xScale={scaleUtc()}
		series={series()}
		seriesLayout="stack"
		props={{
			area: {
				curve: curveNatural,
				'fill-opacity': AREA_FILL_OPACITY,
				line: { class: 'stroke-1' },
				motion: 'tween',
			},
			xAxis: {
				ticks: getXAxisTicks(timeRange),
				format: (date: Date) => formatXAxisDate(date, timeRange),
			},
			yAxis: {
				format: () => '',
			},
		}}
	>
		{#snippet marks({ series: chartSeries, getAreaProps })}
			<defs>
				{#each chartSeries as s (s.key)}
					{@const key = s.key.replace(/\s+/g, '')}
					<linearGradient id="fill{key}" x1="0" y1="0" x2="0" y2="1">
						<stop offset="5%" stop-color={s.color} stop-opacity={GRADIENT_OPACITY_START} />
						<stop offset="95%" stop-color={s.color} stop-opacity={GRADIENT_OPACITY_END} />
					</linearGradient>
				{/each}
			</defs>

			<ChartClipPath
				initialWidth={0}
				motion={{
					width: {
						type: 'tween',
						duration: ANIMATION_DURATION,
						easing: cubicInOut,
					},
				}}
			>
				{#each chartSeries as s, i (s.key)}
					{@const key = s.key.replace(/\s+/g, '')}
					<Area {...getAreaProps(s, i)} fill="url(#fill{key})" />
				{/each}
			</ChartClipPath>
		{/snippet}

		{#snippet tooltip()}
			{#if valueFormatter}
				<Chart.Tooltip labelFormatter={(date: Date) => formatTooltipDate(date, timeRange)}>
					{#snippet formatter({ name, value })}
						<div
							style="--color-bg: var(--color-{name})"
							class="h-full w-1 shrink-0 rounded-[2px] border-(--color-border) bg-(--color-bg)"
						></div>
						<div class="flex flex-1 shrink-0 items-center justify-between leading-none">
							<div class="grid gap-1.5">
								<span class="text-muted-foreground"> {name} </span>
							</div>
							{#if value !== undefined && value !== null}
								{@const formatted = valueFormatter(Number(value))}
								<span class="text-foreground font-mono font-medium tabular-nums">
									{formatted.value.toLocaleString()}
									{formatted.unit}
								</span>
							{/if}
						</div>
					{/snippet}
				</Chart.Tooltip>
			{:else}
				<Chart.Tooltip labelFormatter={(date: Date) => formatTooltipDate(date, timeRange)} indicator="line" />
			{/if}
		{/snippet}
	</AreaChart>
</ChartContainer>
