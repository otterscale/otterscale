<script lang="ts">
	import Icon from '@iconify/svelte';
	import { scaleUtc } from 'd3-scale';
	import { curveLinear } from 'd3-shape';
	import { LineChart } from 'layerchart';
	import { PrometheusDriver, SampleValue } from 'prometheus-query';
	import { onMount } from 'svelte';

	import { ReloadManager } from '$lib/components/custom/reloader';
	import { buttonVariants } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import * as Chart from '$lib/components/ui/chart';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { m } from '$lib/paraglide/messages';
	import { cn } from '$lib/utils';

	let {
		prometheusDriver,
		scope,
		isReloading = $bindable()
	}: { prometheusDriver: PrometheusDriver; scope: string; isReloading: boolean } = $props();

	let latestLatency: number | undefined = $state(undefined);
	let latencies = $state([] as SampleValue[]);
	const trend = $derived(
		latencies.length > 1 && latencies[latencies.length - 2].value !== 0
			? (latencies[latencies.length - 1].value - latencies[latencies.length - 2].value) /
					latencies[latencies.length - 2].value
			: 0
	);

	const configuration = {
		latency: { label: 'Latency', color: 'var(--chart-1)' }
	} satisfies Chart.ChartConfig;

	async function fetchLatestLatency() {
		try {
			const response = await prometheusDriver.instantQuery(
				`histogram_quantile(0.95, sum by(le) (vllm:e2e_request_latency_seconds_bucket{juju_model="${scope}"}))`
			);
			latestLatency = response.result[0]?.value?.value;
		} catch (error) {
			console.error(`Fail to fetch latest latency in scope ${scope}:`, error);
		}
	}

	async function fetchLatencies() {
		try {
			const response = await prometheusDriver.rangeQuery(
				`histogram_quantile(0.95, sum by(le) (rate(vllm:e2e_request_latency_seconds_bucket{juju_model="${scope}"}[5m])))`,
				Date.now() - 60 * 60 * 1000,
				Date.now(),
				10 * 60
			);
			const sampleValues: SampleValue[] = response.result[0]?.values ?? [];
			latencies =
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
			console.error(`Fail to fetch latencies in scope ${scope}:`, error);
		}
	}

	async function fetch() {
		try {
			await Promise.all([fetchLatestLatency(), fetchLatencies()]);
		} catch (error) {
			console.error(`Fail to fetch data in scope ${scope}:`, error);
		}
	}

	const reloadManager = new ReloadManager(fetch);

	let isLoaded = $state(false);
	onMount(async () => {
		try {
			await fetch();
			isLoaded = true;
		} catch (error) {
			console.error(`Fail to fetch data in scope ${scope}:`, error);
		}
	});

	$effect(() => {
		if (isReloading) {
			reloadManager.restart();
		} else {
			reloadManager.stop();
		}
	});
</script>

<Card.Root class="h-full gap-2">
	<Card.Header>
		<Card.Title class="flex flex-wrap items-center justify-between gap-6">
			<div class="flex items-center gap-2 truncate text-sm font-medium tracking-tight">
				<Icon icon="ph:clock" class="size-4.5" />
				{m.latency()}
			</div>
			<Tooltip.Provider>
				<Tooltip.Root>
					<Tooltip.Trigger class={buttonVariants({ variant: 'ghost', size: 'icon' })}>
						<Icon icon="ph:info" class="size-5 text-muted-foreground" />
					</Tooltip.Trigger>
					<Tooltip.Content>
						<p>{m.llm_dashboard_latency_tooltip()}</p>
					</Tooltip.Content>
				</Tooltip.Root>
			</Tooltip.Provider>
		</Card.Title>
	</Card.Header>
	{#if !isLoaded}
		<Card.Content>
			<div class="flex h-9 w-full items-center justify-center">
				<Icon icon="svg-spinners:6-dots-rotate" class="size-10" />
			</div>
		</Card.Content>
	{:else if latestLatency == undefined}
		<Card.Content>
			<div class="flex h-full w-full flex-col items-center justify-center">
				<Icon icon="ph:chart-bar-fill" class="size-6 animate-pulse text-muted-foreground" />
				<p class="p-0 text-xs text-muted-foreground">{m.no_data_display()}</p>
			</div>
		</Card.Content>
	{:else}
		<Card.Content class="flex flex-wrap items-center justify-between gap-6">
			<div class="flex flex-col gap-0.5">
				<div class="text-3xl font-bold">{latestLatency.toFixed(2)}</div>
				<p class="text-sm text-muted-foreground">{m.second()}</p>
			</div>
			<Chart.Container config={configuration} class="h-full w-20">
				<LineChart
					data={latencies}
					x="time"
					xScale={scaleUtc()}
					axis={false}
					series={[
						{
							key: 'value',
							label: configuration.latency.label,
							color: configuration.latency.color
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
									<p class="font-mono">{Number(value).toFixed(2)} {m.second()}</p>
								</div>
							{/snippet}
						</Chart.Tooltip>
					{/snippet}
				</LineChart>
			</Chart.Container>
		</Card.Content>
		<Card.Footer
			class={cn(
				'flex flex-wrap items-center justify-end text-sm leading-none font-medium',
				trend >= 0 ? 'text-emerald-500 dark:text-emerald-400' : 'text-rose-500 dark:text-rose-400'
			)}
		>
			{Math.abs(trend).toFixed(2)} %
			{#if trend >= 0}
				<Icon icon="ph:caret-up" />
			{:else}
				<Icon icon="ph:caret-down" />
			{/if}
		</Card.Footer>
	{/if}
</Card.Root>
