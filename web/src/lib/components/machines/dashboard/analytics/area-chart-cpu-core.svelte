<script lang="ts">
	import Icon from '@iconify/svelte';
	import { scaleUtc } from 'd3-scale';
	import { curveNatural } from 'd3-shape';
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

	const STEP_SECONDS = 60;
	const TIME_RANGE_HOURS = 1;

	const endTime = new Date();
	const startTime = new Date(endTime.getTime() - TIME_RANGE_HOURS * 60 * 60 * 1000);

	let topk = $state(10);

	const query = $derived(
		`topk(${topk}, sum by (cpu) (rate(node_cpu_seconds_total{instance=~"${fqdn}", mode!="idle"}[5m])) * 100)`
	);
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
			<div class="flex h-full w-full items-center justify-center">
				<Icon icon="svg-spinners:blocks-wave" class="m-8 size-32 text-muted-foreground/50" />
			</div>
		{:then response}
			{#if response.length === 0}
				<div class="flex h-full w-full flex-col items-center justify-center">
					<Icon icon="ph:chart-line-fill" class="size-60 animate-pulse text-muted-foreground" />
					<p class="text-base text-muted-foreground">{m.no_data_display()}</p>
				</div>
			{:else}
				{@const data = response}
				{@const chartConfig = generateChartConfig(data)}
				<ChartContainer config={chartConfig} class="aspect-auto h-[300px] w-full">
					<AreaChart
						{data}
						x="date"
						xScale={scaleUtc()}
						series={getSeries(chartConfig)}
						seriesLayout="stack"
						props={{
							area: {
								curve: curveNatural,
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
							<defs>
								{#each chartSeries as series (series.key)}
									{@const key = series.key.replace(/\s+/g, '')}
									<linearGradient id="fill{key}" x1="0" y1="0" x2="0" y2="1">
										<stop offset="5%" stop-color={series.color} stop-opacity={1.0} />
										<stop offset="95%" stop-color={series.color} stop-opacity={0.4} />
									</linearGradient>
								{/each}
							</defs>

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
									{@const key = series.key.replace(/\s+/g, '')}
									<Area {...getAreaProps(series, index)} fill="url(#fill{key})" />
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
								{#snippet formatter({ name, value, payload })}
									{@const top3 = payload
										.map((datum) => datum.value)
										.sort((p, n) => n - p)
										.slice(0, 3)}
									<div
										class="flex w-full shrink-0 items-center justify-between gap-4 leading-none"
										style="--color-bg: var(--color-{name})"
									>
										{#if value !== undefined && value !== null}
											<span class="flex w-full items-center gap-1">
												<Icon icon="ph:square-fill" class="text-(--color-bg)" />
												<p
													class={top3.includes(value)
														? 'fond-bold text-destructive'
														: 'text-foreground'}
												>
													{name}
												</p>
											</span>
											<p
												class={cn(
													'font-mono font-medium whitespace-nowrap tabular-nums',
													top3.includes(value) ? 'fond-bold text-destructive' : 'text-foreground'
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
			<div class="flex h-full w-full flex-col items-center justify-center">
				<Icon icon="ph:chart-bar-fill" class="size-24 animate-pulse text-muted-foreground" />
				<p class="text-base text-muted-foreground">{m.no_data_display()}</p>
			</div>
		{/await}
	</Statistics.Content>
</Statistics.Root>
