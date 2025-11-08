<script lang="ts">
	import { scaleBand } from 'd3-scale';
	import { BarChart } from 'layerchart';
	import { cubicInOut } from 'svelte/easing';

	import * as Chart from '$lib/components/ui/chart/index.js';
	import { formatIO, formatLatencyNano, formatSecond } from '$lib/formatter';

	let {
		data,
		xKey,
		yKey,
		labelKey,
		addFormatter
	}: {
		data: Record<string, any>[] | Record<string, any>;
		xKey: string;
		yKey: string;
		labelKey?: string;
		addFormatter?: 'io' | 'second' | 'nanoSecond';
	} = $props();

	// Convert object to array if needed
	const dataArray = $derived(Array.isArray(data) ? data : Object.values(data));

	// data processing: get top 5 based on xKey
	const getTop5Data = (data: any[], key: string) => {
		// const colors = ['var(--chart-1)', 'var(--chart-1)', 'var(--chart-3)', 'var(--chart-4)', 'var(--chart-5)'];
		return data
			.filter((item) => item && item[key] != null && !isNaN(Number(item[key])))
			.sort((a, b) => Number(b[key]) - Number(a[key]))
			.slice(0, 5)
			.map((item, index) => ({
				...item,
				id: `item-${index}`,
				// color: colors[index],
				color: 'oklch(0.81 0.10 252)'
			}));
	};

	const top5Data = $derived(getTop5Data(dataArray, xKey));

	// Create a fixed array of 5 slots for Y axis
	// const fixedYData = $derived.by(() => {
	// 	const result = Array(5).fill(null);
	// 	top5Data.forEach((item, index) => {
	// 		if (index < 5) {
	// 			result[index] = item;
	// 		}
	// 	});
	// 	return result.map(
	// 		(item, index) => item || { [yKey]: '', [xKey]: 0, id: `empty-${index}`, color: 'oklch(0.81 0.10 252)' },
	// 	);
	// });

	const chartConfig = {
		views: { label: 'Page Views', color: '' },
		value: { label: 'Value', color: '' }
	} satisfies Chart.ChartConfig;
</script>

<Chart.Container config={chartConfig} class="aspect-auto h-[185px] w-full">
	<BarChart
		data={top5Data}
		orientation="horizontal"
		yScale={scaleBand().padding(0.25)}
		y={yKey}
		x={xKey}
		cRange={top5Data.map((c) => c.color)}
		c="color"
		padding={{ right: 16 }}
		grid={false}
		rule={false}
		axis="y"
		props={{
			bars: {
				stroke: 'none',
				radius: 5,
				rounded: 'all',
				initialWidth: 0,
				initialX: 0,
				motion: {
					x: { type: 'tween', duration: 500, easing: cubicInOut },
					width: { type: 'tween', duration: 500, easing: cubicInOut }
				}
			},
			highlight: { area: { fill: 'none' } },
			yAxis: {
				tickLabelProps: {
					textAnchor: 'start',
					dx: 6,
					class: 'stroke-none fill-background!'
				},
				tickLength: 0
			}
		}}
	>
		{#snippet tooltip()}
			<Chart.Tooltip nameKey="views">
				{#snippet formatter({ item, name, value })}
					{@const formattedData = (() => {
						if (addFormatter === 'io') {
							const { value: formattedValue, unit } = formatIO(Number(value));
							return { display: `${formattedValue} ${unit}` };
						} else if (addFormatter === 'second') {
							const { value: formattedValue, unit } = formatSecond(Number(value));
							return { display: `${formattedValue} ${unit}` };
						} else if (addFormatter === 'nanoSecond') {
							const { value: formattedValue, unit } = formatLatencyNano(Number(value));
							return { display: `${formattedValue} ${unit}` };
						} else {
							return { display: value };
						}
					})()}
					<div
						style="--color-bg: {item.color}"
						class="aspect-square h-full w-fit shrink-0 border-(--color-border) bg-(--color-bg)"
					></div>
					<div class="flex flex-1 shrink-0 items-center justify-between text-xs leading-none">
						<div class="grid gap-1.5">
							<span class="text-muted-foreground">{labelKey ? labelKey : name}</span>
						</div>
						<p class="font-mono">{formattedData.display}</p>
					</div>
				{/snippet}
			</Chart.Tooltip>
		{/snippet}
	</BarChart>
</Chart.Container>
