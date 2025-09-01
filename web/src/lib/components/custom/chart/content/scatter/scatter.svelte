<script lang="ts">
	import * as Chart from '$lib/components/ui/chart/index.js';
	import { scaleUtc } from 'd3-scale';
	import { ScatterChart, type ChartContextValue } from 'layerchart';

	let {
		data,
		xKey,
		yKey,
		timeRange = '90d',
	}: {
		data: Record<string, any>;
		xKey: string;
		yKey: string;
		timeRange?: string;
	} = $props();
	let context = $state<ChartContextValue>();

	const chartConfigNew = {
		// '1-5m-8KB-1GB-1': { label: '1-5m-8KB-1GB-1', color: 'var(--chart-3)' },
		// '32-1m-4KB-1GB-1': { label: '32-1m-4KB-1GB-1', color: 'var(--chart-5)' }
		// 'groupNames-3': { label: 'groupNames-3', color: 'var(--chart-5)' },
		// 'groupNames-4': { label: 'groupNames-4', color: 'var(--chart-6)' },
		// 'groupNames-5': { label: 'groupNames-5', color: 'var(--chart-7)' }
	} satisfies Chart.ChartConfig;
</script>

<!-- legend={{
			classes: { root: '-mb-[50px] w-full overflow-auto' }
		}} -->
<!-- <Chart.Container config={computedChartConfig()} class="aspect-auto h-[250px] w-full"> -->
<Chart.Container config={chartConfigNew} class="aspect-auto h-[250px] w-full">
	<ScatterChart
		legend
		x={xKey}
		y={yKey}
		series={Object.values(data)}
		props={{
			xAxis: {
				// ticks: (scale) => scaleUtc(scale.domain(), scale.range()).ticks(4),
				// format: (d: Date) => {
				// 	return d.toLocaleTimeString('en-US', {
				// 		hour: 'numeric',
				// 		minute: '2-digit'
				// 	});
				// }
				format: () => '',
			},
			yAxis: {
				format: () => '',
			},
		}}
	>
		<!-- {#snippet tooltip()}
			<Chart.Tooltip
				nameKey="views"
				labelFormatter={(v: Date) => {
					return v.toLocaleDateString('en-US', {
						month: 'short',
						day: 'numeric',
						year: 'numeric'
					});
				}}
			/>
			<Chart.Tooltip />
		{/snippet} -->
	</ScatterChart>
</Chart.Container>
