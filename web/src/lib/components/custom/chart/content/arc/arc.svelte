<script lang="ts">
	import { ArcChart, Text } from 'layerchart';
	import * as Chart from '$lib/components/ui/chart/index.js';

	let {
		data = [{ value: NaN }],
		subtitle,
	}: {
		data?: Array<any>;
		subtitle?: string;
	} = $props();

	const chartConfig = {
		data: { color: 'var(--chart-4)' },
	} satisfies Chart.ChartConfig;
</script>

<Chart.Container config={chartConfig} class="mx-auto aspect-square h-[250px] w-full">
	<ArcChart
		{data}
		outerRadius={-5}
		innerRadius={-12}
		padding={30}
		range={[-120, 120]}
		maxValue={100}
		cornerRadius={20}
		series={[
			{
				key: 'data',
				color: chartConfig.data.color,
			},
		]}
		props={{
			arc: { track: { fill: 'var(--muted)' }, motion: 'tween' },
			tooltip: { context: { hideDelay: 350 } },
		}}
		tooltip={false}
	>
		{#snippet belowMarks()}
			<circle cx="0" cy="0" r="60" class="fill-background" />
		{/snippet}

		{#snippet aboveMarks()}
			<Text
				value={!isNaN(data[0].value) ? `${Math.round(data[0].value)}%` : 'NaN'}
				textAnchor="middle"
				verticalAnchor="middle"
				class="fill-foreground text-4xl! font-bold"
				dy={3}
			/>
			{#if subtitle}
				<Text
					value={subtitle}
					textAnchor="middle"
					verticalAnchor="middle"
					class="fill-muted-foreground!"
					dy={22}
				/>
			{/if}
		{/snippet}
	</ArcChart>
</Chart.Container>
