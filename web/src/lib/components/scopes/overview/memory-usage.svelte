<script lang="ts">
	import Icon from '@iconify/svelte';
	import type { PrometheusDriver, SampleValue } from 'prometheus-query';
	import { onDestroy, onMount } from 'svelte';

	import { scaleUtc } from 'd3-scale';
	import { curveLinear } from 'd3-shape';
	import { LineChart } from 'layerchart';

	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Card from '$lib/components/ui/card';
	import * as Chart from '$lib/components/ui/chart';
	import * as Select from '$lib/components/ui/select';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { Progress } from '$lib/components/ui/progress';
	import { m } from '$lib/paraglide/messages';

	let {
		prometheusDriver,
		scope,
		isReloading = $bindable()
	}: { prometheusDriver: PrometheusDriver; scope: string; isReloading: boolean } = $props();

	let memoryUsage: SampleValue | undefined = $state(undefined);
	let memoryTotal: SampleValue | undefined = $state(undefined);
	let memoryHistory = $state([] as SampleValue[]);

	const timeRanges = {
		'10min': { label: 'Last 10 minutes', duration: 10 * 60 * 1000, step: 10 },
		'1hour': { label: 'Last 1 hour', duration: 60 * 60 * 1000, step: 60 },
		'1day': { label: 'Last 1 day', duration: 24 * 60 * 60 * 1000, step: 900 },
		'1week': { label: 'Last 1 week', duration: 7 * 24 * 60 * 60 * 1000, step: 10800 }
	};
	let selectedTimeRange = $state('1hour');

	const chartConfig = {
		memory: { label: 'Memory Usage', color: 'var(--chart-1)' }
	} satisfies Chart.ChartConfig;

	async function fetchMemoryUsage() {
		const usageResponse = await prometheusDriver.instantQuery(
			`sum(node_memory_MemTotal_bytes{juju_model="${scope}", container!=""} - node_memory_MemAvailable_bytes{juju_model="${scope}", container!=""})`
		);
		memoryUsage = usageResponse.result[0]?.value ?? undefined;

		const totalResponse = await prometheusDriver.instantQuery(
			`sum(node_memory_MemTotal_bytes{juju_model="${scope}", container!=""})`
		);
		memoryTotal = totalResponse.result[0]?.value ?? undefined;
	}

	async function fetchMemoryHistory() {
		try {
			const range = timeRanges[selectedTimeRange as keyof typeof timeRanges];
			const response = await prometheusDriver.rangeQuery(
				`sum(node_memory_MemTotal_bytes{juju_model="${scope}", container!=""} - node_memory_MemAvailable_bytes{juju_model="${scope}", container!=""}) / sum(node_memory_MemTotal_bytes{juju_model="${scope}", container!=""}) * 100`,
				Date.now() - range.duration,
				Date.now(),
				range.step
			);
			const sampleValues: SampleValue[] = response.result[0]?.values ?? [];
			memoryHistory =
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
			console.error(`Fail to fetch memory history in scope ${scope}:`, error);
		}
	}

	async function fetch() {
		try {
			await Promise.all([fetchMemoryUsage(), fetchMemoryHistory()]);
		} catch (error) {
			console.error('Failed to fetch memory usage:', error);
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

	function formatBytes(bytes: number): string {
		const gb = bytes / (1024 ** 3);
		return `${gb.toFixed(1)} GB`;
	}
</script>

<Card.Root class="relative h-full min-h-[140px] gap-2 overflow-hidden">
	<Icon
		icon="ph:memory"
		class="absolute -right-10 bottom-0 size-36 text-8xl tracking-tight text-nowrap text-primary/5 uppercase group-hover:hidden -z-0"
	/>
	<Card.Header>
		<Card.Title>Memory Usage</Card.Title>
		<Card.Description class="flex items-center justify-between z-10">
			<span>Cluster Memory Usage (GB)</span>
			<Select.Root
				type="single"
				value={selectedTimeRange}
				onValueChange={(v) => {
					if (v) {
						selectedTimeRange = v;
						fetchMemoryHistory();
					}
				}}
			>
				<Select.Trigger class="h-6 w-[130px] px-2 text-[12px]">
					{timeRanges[selectedTimeRange as keyof typeof timeRanges].label}
				</Select.Trigger>
				<Select.Content class="z-50">
					{#each Object.entries(timeRanges) as [key, { label }]}
						<Select.Item value={key} label={label} class="text-xs">
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
	{:else if !memoryUsage && !memoryTotal}
		<div class="flex h-full w-full flex-col items-center justify-center">
			<Icon icon="ph:chart-bar-fill" class="size-6 animate-pulse text-muted-foreground" />
			<p class="p-0 text-xs text-muted-foreground">{m.no_data_display()}</p>
		</div>
	{:else}
		<Card.Content class="flex items-center justify-between gap-4">
			<div class="flex-1">
				{#if memoryUsage && memoryTotal}
					<div class="text-3xl font-bold">{Math.round((Number(memoryUsage.value) / Number(memoryTotal.value)) * 100) + '%'}</div>
					<Tooltip.Root>
						<Tooltip.Trigger class="w-full cursor-default">
							<Progress value={(Number(memoryUsage.value) / Number(memoryTotal.value)) * 100} class="mt-2 h-2" />
						</Tooltip.Trigger>
						<Tooltip.Content>
							<p>Usage: {formatBytes(Number(memoryUsage.value))}</p>
							<p>Total: {formatBytes(Number(memoryTotal.value))}</p>
						</Tooltip.Content>
					</Tooltip.Root>
				{/if}
			</div>
			<div class="flex flex-col items-end gap-1">
				<Chart.Container config={chartConfig} class="h-10 w-40">
					<LineChart
						data={memoryHistory}
						x="time"
						xScale={scaleUtc()}
						axis={false}
						series={[
							{
								key: 'value',
								label: chartConfig.memory.label,
								color: chartConfig.memory.color
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
		</Card.Content>
	{/if}
</Card.Root>