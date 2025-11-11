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

	let latestAvailableModels = $state(0);
	let availableModels = $state([] as SampleValue[]);
	const trend = $derived(
		availableModels.length > 1 && availableModels[availableModels.length - 2].value !== 0
			? (availableModels[availableModels.length - 1].value -
					availableModels[availableModels.length - 2].value) /
					availableModels[availableModels.length - 2].value
			: 0
	);

	const configuration = {
		usage: { label: 'Models', color: 'var(--chart-1)' }
	} satisfies Chart.ChartConfig;

	async function fetch() {
		prometheusDriver
			.instantQuery(`count by(endpoint) (vllm:gpu_cache_usage_perc{juju_model="${scope}"})`)
			.then((response) => {
				latestAvailableModels = response.result[0]?.value?.value;
			});
		prometheusDriver
			.rangeQuery(
				`count by(endpoint) (vllm:gpu_cache_usage_perc{juju_model="${scope}"})`,
				Date.now() - 10 * 60 * 1000,
				Date.now(),
				2 * 60
			)
			.then((response) => {
				availableModels = response.result[0]?.values ?? [];
			});
	}

	const reloadManager = new ReloadManager(fetch);

	let isLoading = $state(true);
	onMount(async () => {
		try {
			await fetch();
			isLoading = false;
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

{#if isLoading}
	Loading
{:else}
	<Card.Root class="h-full gap-2">
		<Card.Header>
			<Card.Title class="flex flex-wrap items-center justify-between gap-6">
				<div class="flex items-center gap-2 truncate text-sm font-medium tracking-tight">
					<Icon icon="ph:robot" class="size-4.5" />
					{m.models()}
				</div>
				<Tooltip.Provider>
					<Tooltip.Root>
						<Tooltip.Trigger class={buttonVariants({ variant: 'ghost', size: 'icon' })}>
							<Icon icon="ph:info" class="size-5 text-muted-foreground" />
						</Tooltip.Trigger>
						<Tooltip.Content>
							<p>{m.llm_dashboard_models_tooltip()}</p>
						</Tooltip.Content>
					</Tooltip.Root>
				</Tooltip.Provider>
			</Card.Title>
		</Card.Header>
		<Card.Content class="flex flex-wrap items-center justify-between gap-6">
			<div class="flex flex-col gap-0.5">
				<div class="text-3xl font-bold">{latestAvailableModels}</div>
				<p class="text-sm text-muted-foreground">{m.models()}</p>
			</div>
			<Chart.Container config={configuration} class="h-full w-20">
				<LineChart
					data={availableModels}
					x="time"
					xScale={scaleUtc()}
					axis={false}
					series={[
						{
							key: 'value',
							label: configuration.usage.label,
							color: configuration.usage.color
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
									<p class="font-mono">{value}</p>
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
	</Card.Root>
{/if}
