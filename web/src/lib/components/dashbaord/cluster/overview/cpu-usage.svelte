<script lang="ts">
	import Icon from '@iconify/svelte';
	import { scaleUtc } from 'd3-scale';
	import { curveLinear } from 'd3-shape';
	import { ArcChart, Text } from 'layerchart';
	import { LineChart } from 'layerchart';
	import type { PrometheusDriver, SampleValue } from 'prometheus-query';
	import { onDestroy, onMount } from 'svelte';

	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Card from '$lib/components/ui/card';
	import * as Chart from '$lib/components/ui/chart';
	import { Progress } from '$lib/components/ui/progress';
	import * as Select from '$lib/components/ui/select';
	import { formatCapacity } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';

	let {
		prometheusDriver,
		scope,
		isReloading = $bindable()
	}: { prometheusDriver: PrometheusDriver; scope: string; isReloading: boolean } = $props();

	let avgCpuUsage: SampleValue | undefined = $state(undefined);
	let cpuHistory = $state([] as SampleValue[]);

	const timeRanges = {
		'10min': { label: 'Last 10 minutes', duration: 10 * 60 * 1000, step: 10 },
		'1hour': { label: 'Last 1 hour', duration: 60 * 60 * 1000, step: 60 },
		'1day': { label: 'Last 1 day', duration: 24 * 60 * 60 * 1000, step: 900 },
		'1week': { label: 'Last 1 week', duration: 7 * 24 * 60 * 60 * 1000, step: 10800 }
	};
	let selectedTimeRange = $state('1hour');

	const chartConfig = {
		cpu: { label: 'CPU Usage', color: 'var(--chart-1)' }
	} satisfies Chart.ChartConfig;

	async function fetchCpuUsage() {
		const response = await prometheusDriver.instantQuery(
			`sum(rate(node_cpu_seconds_total{mode!="idle", juju_model="${scope}"}[5m])) / sum(rate(node_cpu_seconds_total{juju_model="${scope}"}[5m])) * 100`
		);
		avgCpuUsage = response.result[0]?.value ?? undefined;
		console.log(avgCpuUsage);
	}

	async function fetchCpuHistory() {
		try {
			const range = timeRanges[selectedTimeRange as keyof typeof timeRanges];
			const response = await prometheusDriver.rangeQuery(
				`sum(rate(node_cpu_seconds_total{mode!="idle", juju_model="${scope}"}[5m])) / sum(rate(node_cpu_seconds_total{juju_model="${scope}"}[5m])) * 100`,
				Date.now() - range.duration,
				Date.now(),
				range.step
			);
			const sampleValues: SampleValue[] = response.result[0]?.values ?? [];
			cpuHistory =
				sampleValues.length > 0
					? sampleValues.map(
							(sampleValue) =>
								({
									time: sampleValue.time,
									value: sampleValue && !isNaN(Number(sampleValue.value)) ? sampleValue.value : 0
								}) as SampleValue
						)
					: [];
		} catch (error) {
			console.error(`Fail to fetch CPU history in scope ${scope}:`, error);
		}
	}

	async function fetch() {
		try {
			await Promise.all([fetchCpuUsage(), fetchCpuHistory()]);
		} catch (error) {
			console.error('Failed to fetch CPU usage:', error);
		}
	}

	const reloadManager = new ReloadManager(fetch);

	$effect(() => {
		if (isReloading) {
			reloadManager.restart();
		} else {
			reloadManager.stop();
		}
	});

	let isLoaded = $state(false);
	onMount(async () => {
		await fetch();
		isLoaded = true;
	});
	onDestroy(() => {
		reloadManager.stop();
	});
</script>

<Card.Root class="relative h-full min-h-[140px] gap-2 overflow-hidden">
	<Icon
		icon="ph:cpu"
		class="absolute -right-10 bottom-0 -z-0 size-36 text-8xl tracking-tight text-nowrap text-primary/5 uppercase group-hover:hidden"
	/>
	<Card.Header>
		<Card.Title>CPU Usage (%)</Card.Title>
		<Card.Description class="z-10 flex items-center justify-between">
			<span>Cluster Average CPU Usage (%)</span>
			<Select.Root
				type="single"
				value={selectedTimeRange}
				onValueChange={(v) => {
					if (v) {
						selectedTimeRange = v;
						fetchCpuHistory();
					}
				}}
			>
				<Select.Trigger class="h-6 w-[130px] px-2 text-[12px]">
					{timeRanges[selectedTimeRange as keyof typeof timeRanges].label}
				</Select.Trigger>
				<Select.Content class="z-50">
					{#each Object.entries(timeRanges) as [key, { label }]}
						<Select.Item value={key} {label} class="text-xs">
							{label}
						</Select.Item>
					{/each}
				</Select.Content>
			</Select.Root>
		</Card.Description>
	</Card.Header>
	{#if !isLoaded}
		<div class="flex h-9 w-full items-center justify-center">
			<Icon icon="svg-spinners:6-dots-rotate" class="size-10" />
		</div>
	{:else if !avgCpuUsage}
		<div class="flex h-full w-full flex-col items-center justify-center">
			<Icon icon="ph:chart-bar-fill" class="size-6 animate-pulse text-muted-foreground" />
			<p class="p-0 text-xs text-muted-foreground">{m.no_data_display()}</p>
		</div>
	{:else}
		<Card.Content>
			<Chart.Container class="h-[200px] w-full px-2 pt-2" config={chartConfig}>
				<ArcChart
					data={[{ value: avgCpuUsage.value }]}
					outerRadius={88}
					innerRadius={66}
					trackOuterRadius={83}
					trackInnerRadius={72}
					padding={40}
					range={[90, -270]}
					maxValue={100}
					series={[
						{
							key: 'value',
							color: chartConfig.cpu.color
						}
					]}
					props={{
						arc: { track: { fill: 'var(--muted)' }, motion: 'tween' },
						tooltip: { context: { hideDelay: 350 } }
					}}
					tooltip={false}
				>
					{#snippet belowMarks()}
						<circle cx="0" cy="0" r="80" class="fill-background" />
					{/snippet}
					{#snippet aboveMarks()}
						<Text
							value={avgCpuUsage ? `${Math.round(Number(avgCpuUsage.value))}%` : undefined}
							textAnchor="middle"
							verticalAnchor="middle"
							class="fill-foreground text-4xl! font-bold"
						/>
					{/snippet}
				</ArcChart>
			</Chart.Container>
		</Card.Content>
		<!-- <Card.Content class="flex items-center justify-between gap-4">
			<div class="flex-1">
				<div class="text-3xl font-bold">
					{avgCpuUsage ? Math.round(Number(avgCpuUsage.value)) + '%' : 'N/A'}
				</div>
				{#if avgCpuUsage}
					<Progress value={Number(avgCpuUsage.value)} class="mt-2 h-2" />
				{/if}
			</div>
			<div class="flex flex-col items-end gap-1">
				<Chart.Container config={chartConfig} class="h-10 w-40">
					<LineChart
						data={cpuHistory}
						x="time"
						xScale={scaleUtc()}
						axis={false}
						series={[
							{
								key: 'value',
								label: chartConfig.cpu.label,
								color: chartConfig.cpu.color
							}
						]}
						props={{
							spline: { curve: curveLinear, motion: 'tween', strokeWidth: 2 },
							xAxis: {
								format: (v: Date) => v.toLocaleDateString('en-US', { month: 'short' })
							},
							highlight: { points: { r: 4 } }
						}}
					>
						{#snippet tooltip()}
							<Chart.Tooltip hideLabel>
								{#snippet formatter({ item, name, value })}
									<div
										style="--color-bg: {item.color}"
										class="aspect-square h-full w-fit shrink-0 border-(--color-border) bg-(--color-bg)"
									></div>
									<div
										class="flex flex-1 shrink-0 items-center justify-between gap-2 text-xs leading-none"
									>
										<div class="grid gap-1.5">
											<span class="text-muted-foreground">{name}</span>
										</div>
										<p class="font-mono">{Number(value).toFixed(2)} %</p>
									</div>
								{/snippet}
							</Chart.Tooltip>
						{/snippet}
					</LineChart>
				</Chart.Container>
			</div>
		</Card.Content> -->
	{/if}
</Card.Root>
