<script lang="ts">
	import * as Chart from '$lib/components/ui/chart/index.js';
	import { scaleUtc } from 'd3-scale';
	import { BarChart, type ChartContextValue, Highlight } from 'layerchart';
	import { cubicInOut } from 'svelte/easing';

	let {
		activeChart = $bindable(),
		data,
		chartConfig,
		timeRange = '90d',
	}: {
		activeChart: string;
		data: Array<any>;
		chartConfig: any;
		timeRange?: string;
	} = $props();
	let context = $state<ChartContextValue>();
	const activeSeries = $derived([
		{
			key: activeChart,
			label: chartConfig[activeChart].label,
			color: chartConfig[activeChart].color,
		},
	]);
</script>

<Chart.Container config={chartConfig} class="aspect-auto h-[250px] w-full">
	<BarChart
		bind:context
		{data}
		x="date"
		axis="x"
		series={activeSeries}
		props={{
			bars: {
				stroke: 'none',
				rounded: 'none',
				// use the height of the chart to animate the bars
				initialY: context?.height,
				initialHeight: 0,
				motion: {
					y: { type: 'tween', duration: 500, easing: cubicInOut },
					height: { type: 'tween', duration: 500, easing: cubicInOut },
				},
			},
			highlight: { area: { fill: 'none' } },
			xAxis: {
				format: (d: Date) => {
					return d.toLocaleDateString('en-US', {
						month: 'short',
						day: '2-digit',
					});
				},
				ticks: (scale) => scaleUtc(scale.domain(), scale.range()).ticks(),
			},
		}}
	>
		{#snippet belowMarks()}
			<Highlight area={{ class: 'fill-muted' }} />
		{/snippet}
		{#snippet tooltip()}
			<Chart.Tooltip
				nameKey="views"
				labelFormatter={(v: Date) => {
					return v.toLocaleDateString('en-US', {
						month: 'short',
						day: 'numeric',
						year: 'numeric',
					});
				}}
			/>
		{/snippet}
	</BarChart>
</Chart.Container>
