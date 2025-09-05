<script lang="ts">
	import { ScatterChart } from 'layerchart';

	import * as Chart from '$lib/components/ui/chart/index.js';

	let {
		data,
		xKey,
		yKey,
	}: {
		data: Record<string, any>;
		xKey: string;
		yKey: string;
	} = $props();

	const chartConfigNew = {} satisfies Chart.ChartConfig;
</script>

<Chart.Container config={chartConfigNew} class="aspect-auto h-[250px] w-full">
	<ScatterChart
		legend
		x={xKey}
		y={yKey}
		series={Object.values(data)}
		props={{
			xAxis: {
				format: () => '',
			},
			yAxis: {
				format: () => '',
			},
		}}
	></ScatterChart>
</Chart.Container>
