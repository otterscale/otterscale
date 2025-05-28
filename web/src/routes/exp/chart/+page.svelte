<script lang="ts">
	import { AreaChart, BarChart, LineChart, PieChart } from 'layerchart';
	import { group, flatGroup } from 'd3-array';
	import { curveCatmullRom } from 'd3-shape';
	import { createDateSeries, longData } from './gen';

	let renderContext: 'svg' | 'canvas' = 'svg';
	let debug = false;

	const dataByYear = group(longData, (d) => d.year);
	const data = dataByYear.get(2019) ?? [];

	const dateSeriesData = createDateSeries({
		count: 10,
		min: 20,
		max: 100,
		value: 'integer',
		keys: ['value', 'baseline']
	});
</script>

<div class="h-[300px] rounded border p-4">
	<AreaChart
		data={dateSeriesData}
		x="date"
		y="value"
		props={{ area: { curve: curveCatmullRom } }}
		{renderContext}
		{debug}
	/>
</div>

<div class="h-[300px] rounded border p-4">
	<BarChart
		data={dateSeriesData}
		x="date"
		series={[
			{
				key: 'value',
				color: 'hsl(var(--color-primary))'
			}
		]}
		{renderContext}
		{debug}
	/>
</div>

<div class="h-[300px] rounded border p-4">
	<LineChart
		data={dateSeriesData}
		x="date"
		y="value"
		props={{ spline: { curve: curveCatmullRom } }}
		{renderContext}
		{debug}
	/>
</div>

<div class="h-[300px] resize overflow-auto rounded border">
	<PieChart
		{data}
		key="fruit"
		value="value"
		range={[-90, 90]}
		outerRadius={300 / 2}
		innerRadius={-20}
		cornerRadius={10}
		padAngle={0.02}
		props={{ group: { y: 80 } }}
		{renderContext}
		{debug}
	/>
</div>
