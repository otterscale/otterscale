<script lang="ts">
	import Icon from '@iconify/svelte';
	import { scaleUtc } from 'd3-scale';
	import { curveLinear } from 'd3-shape';
	import { LineChart } from 'layerchart';
	import type { PrometheusDriver, SampleValue } from 'prometheus-query';
	import { onDestroy, onMount } from 'svelte';

	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Card from '$lib/components/ui/card';
	import * as Chart from '$lib/components/ui/chart';
	import * as Select from '$lib/components/ui/select';
	import { m } from '$lib/paraglide/messages';

	let {
		prometheusDriver,
		scope,
		isReloading = $bindable()
	}: { prometheusDriver: PrometheusDriver; scope: string; isReloading: boolean } = $props();

	let gpuCount: SampleValue | undefined = $state(undefined);
	let avgGpuTemperature: SampleValue | undefined = $state(undefined);
	let gpuTempHistory = $state([] as SampleValue[]);

	const timeRanges = {
		'10min': { label: 'Last 10 minutes', duration: 10 * 60 * 1000, step: 10 },
		'1hour': { label: 'Last 1 hour', duration: 60 * 60 * 1000, step: 60 },
		'1day': { label: 'Last 1 day', duration: 24 * 60 * 60 * 1000, step: 900 },
		'1week': { label: 'Last 1 week', duration: 7 * 24 * 60 * 60 * 1000, step: 10800 }
	};
	let selectedTimeRange = $state('1hour');

	const chartConfig = {
		temp: { label: 'Temperature', color: 'var(--chart-1)' }
	} satisfies Chart.ChartConfig;

	async function fetchGpuInfo() {
		const countResponse = await prometheusDriver.instantQuery(
			`count(DCGM_FI_DEV_GPU_UTIL{juju_model="${scope}"})`
		);
		gpuCount = countResponse.result[0]?.value ?? undefined;

		const tempResponse = await prometheusDriver.instantQuery(
			`avg(DCGM_FI_DEV_GPU_TEMP{juju_model="${scope}"})`
		);
		avgGpuTemperature = tempResponse.result[0]?.value ?? undefined;
	}

	async function fetchGpuTempHistory() {
		try {
			const range = timeRanges[selectedTimeRange as keyof typeof timeRanges];
			const response = await prometheusDriver.rangeQuery(
				`avg(DCGM_FI_DEV_GPU_TEMP{juju_model="${scope}"})`,
				Date.now() - range.duration,
				Date.now(),
				range.step
			);
			const sampleValues: SampleValue[] = response.result[0]?.values ?? [];
			gpuTempHistory =
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
			console.error(`Fail to fetch GPU temp history in scope ${scope}:`, error);
		}
	}

	async function fetch() {
		try {
			await Promise.all([fetchGpuInfo(), fetchGpuTempHistory()]);
		} catch (error) {
			console.error('Failed to fetch GPU info:', error);
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
		<Card.Title>NVIDIA GPUs</Card.Title>
		<Card.Description class="z-10 flex items-center justify-between">
			<span>Cluster NVIDIA GPU Cards</span>
			<Select.Root
				type="single"
				value={selectedTimeRange}
				onValueChange={(v) => {
					if (v) {
						selectedTimeRange = v;
						fetchGpuTempHistory();
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
	{:else if !gpuCount && !avgGpuTemperature}
		<div class="flex h-full w-full flex-col items-center justify-center">
			<Icon icon="ph:chart-bar-fill" class="size-6 animate-pulse text-muted-foreground" />
			<p class="p-0 text-xs text-muted-foreground">{m.no_data_display()}</p>
		</div>
	{:else}
		<Card.Content class="flex items-center justify-between gap-4">
			<div class="flex gap-8">
				<div class="flex flex-col">
					<span class="text-3xl font-bold">{gpuCount?.value ?? 'N/A'}</span>
					<span class="text-1xl font-medium tracking-wider text-muted-foreground uppercase"
						>Total NVIDIA GPUs</span
					>
				</div>
				<div class="flex flex-col">
					<span class="text-3xl font-bold"
						>{avgGpuTemperature ? Math.round(Number(avgGpuTemperature.value)) : 'N/A'}°C</span
					>
					<span class="text-1xl font-medium tracking-wider text-muted-foreground uppercase"
						>Average Temperature</span
					>
				</div>
			</div>
			<div class="flex flex-col items-end gap-1">
				<Chart.Container config={chartConfig} class="h-10 w-40">
					<LineChart
						data={gpuTempHistory}
						x="time"
						xScale={scaleUtc()}
						axis={false}
						series={[
							{
								key: 'value',
								label: chartConfig.temp.label,
								color: chartConfig.temp.color
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
										<p class="font-mono">{Number(value).toFixed(2)} °C</p>
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
