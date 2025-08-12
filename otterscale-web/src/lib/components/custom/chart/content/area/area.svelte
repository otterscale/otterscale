<script lang="ts">
	import { scaleUtc } from 'd3-scale';
	import { Area, AreaChart, ChartClipPath } from 'layerchart';
	import { curveNatural } from 'd3-shape';
	import ChartContainer from '$lib/components/ui/chart/chart-container.svelte';
	import * as Chart from '$lib/components/ui/chart/index.js';
	import { cubicInOut } from 'svelte/easing';
	import {
		generateChartConfig,
		getSeries,
		type DataPoint,
		type ChartConfig
	} from '$lib/components/custom/prometheus';
	import {
		formatXAxisDate,
		formatTooltipDate,
		getXAxisTicks,
		type TimeRange
	} from '$lib/components/custom/chart/units/formatter';

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
	}

	let { data, chartConfig, timeRange, class: className = '' }: Props = $props();

	// Derived reactive values
	const computedChartConfig = $derived(() => chartConfig ?? generateChartConfig(data));
	const series = $derived(() => getSeries(computedChartConfig()));
</script>

<ChartContainer
	config={computedChartConfig()}
	class="aspect-auto {DEFAULT_CHART_HEIGHT} w-full {className}"
>
	<AreaChart
		legend
		{data}
		x="date"
		xScale={scaleUtc()}
		series={series()}
		props={{
			area: {
				curve: curveNatural,
				'fill-opacity': AREA_FILL_OPACITY,
				line: { class: 'stroke-1' },
				motion: 'tween'
			},
			xAxis: {
				ticks: getXAxisTicks(timeRange),
				format: (date: Date) => formatXAxisDate(date, timeRange)
			},
			yAxis: {
				format: () => ''
			}
		}}
	>
		{#snippet marks({ series: chartSeries, getAreaProps })}
			<defs>
				{#each chartSeries as s, i (s.key)}
					<linearGradient id="fill{s.key}" x1="0" y1="0" x2="0" y2="1">
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
						easing: cubicInOut
					}
				}}
			>
				{#each chartSeries as s, i (s.key)}
					<Area {...getAreaProps(s, i)} fill="url(#fill{s.key})" />
				{/each}
			</ChartClipPath>
		{/snippet}

		{#snippet tooltip()}
			<Chart.Tooltip
				labelFormatter={(date: Date) => formatTooltipDate(date, timeRange)}
				indicator="line"
			/>
		{/snippet}
	</AreaChart>
</ChartContainer>
