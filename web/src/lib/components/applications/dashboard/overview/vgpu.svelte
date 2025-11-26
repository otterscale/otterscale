<script lang="ts">
	import Icon from '@iconify/svelte';
	import { scaleBand } from 'd3-scale';
	import { BarChart, type ChartContextValue, Highlight } from 'layerchart';
	import { InstantVector, PrometheusDriver } from 'prometheus-query';
	import { onDestroy, onMount } from 'svelte';
	import { cubicInOut } from 'svelte/easing';

	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Card from '$lib/components/ui/card';
	import * as Chart from '$lib/components/ui/chart/index.js';
	import { m } from '$lib/paraglide/messages';

	let {
		prometheusDriver,
		scope,
		isReloading = $bindable()
	}: { prometheusDriver: PrometheusDriver; scope: string; isReloading: boolean } = $props();

	let memoryUsage: Record<string, number>[] = $state([]);

	const configuration = {
		usage: { label: 'Usage', color: 'var(--chart-1)' }
	} satisfies Chart.ChartConfig;

	let context = $state<ChartContextValue>();

	async function fetch() {
		try {
			const response = await prometheusDriver.instantQuery(
				`
				topk(10, avg by (nodeid) (nodeGPUMemoryPercentage{juju_model="${scope}"}))
				`
			);
			const instanceVectors: InstantVector[] = response.result;
			memoryUsage = instanceVectors
				.sort((p, n) => n.value.value - p.value.value)
				.map((instanceVector) =>
					Object.fromEntries([
						['node', (instanceVector.metric.labels as { nodeid?: string }).nodeid],
						['usage', instanceVector.value.value]
					])
				);
		} catch (error) {
			console.error('Failed to fetch VGPU memory usage:', error);
		}
	}

	const reloadManager = new ReloadManager(fetch);

	let isLoaded = $state(false);
	onMount(async () => {
		await fetch();
		isLoaded = true;
	});
	onDestroy(() => {
		reloadManager.stop();
	});

	$effect(() => {
		if (isReloading) {
			reloadManager.restart();
		} else {
			reloadManager.stop();
		}
	});
</script>

{#if !isLoaded}
	<Card.Root class="h-full gap-2">
		<Card.Header class="h-[42px]">
			<Card.Title>{m.vgpu()}</Card.Title>
			<Card.Description>
				{m.vgpu_usage_description()}
			</Card.Description>
		</Card.Header>
		<Card.Content>
			<div class="flex h-[230px] w-full items-center justify-center">
				<Icon icon="svg-spinners:6-dots-rotate" class="size-12" />
			</div>
		</Card.Content>
	</Card.Root>
{:else}
	<Card.Root>
		<Card.Header>
			<Card.Title>{m.vgpu()}</Card.Title>
			<Card.Description>
				{m.vgpu_usage_description()}
			</Card.Description>
		</Card.Header>
		<Card.Content>
			<Chart.Container class="h-[230px] w-full px-2 pt-2" config={configuration}>
				<BarChart
					bind:context
					data={memoryUsage}
					orientation="horizontal"
					yScale={scaleBand().padding(0.25)}
					y="node"
					axis="y"
					padding={{ right: 16 }}
					rule={false}
					series={[
						{
							key: 'usage',
							label: configuration.usage.label,
							color: configuration.usage.color
						}
					]}
					props={{
						bars: {
							stroke: 'none',
							initialY: context?.height,
							initialHeight: 0,
							motion: {
								y: { type: 'tween', duration: 500, easing: cubicInOut },
								height: { type: 'tween', duration: 500, easing: cubicInOut }
							}
						},
						highlight: { area: false },
						yAxis: {
							tickLabelProps: {
								textAnchor: 'start',
								dx: 8,
								class: 'stroke-none fill-background'
							},
							tickLength: 0
						}
					}}
				>
					{#snippet belowMarks()}
						<Highlight area={{ class: 'fill-muted' }} />
					{/snippet}

					{#snippet tooltip()}
						<Chart.Tooltip indicator="dot">
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
									<p class="font-mono">{(Number(value) * 100).toFixed(2)} %</p>
								</div>
							{/snippet}
						</Chart.Tooltip>
					{/snippet}
				</BarChart>
			</Chart.Container>
		</Card.Content>
	</Card.Root>
{/if}
