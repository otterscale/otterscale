<script lang="ts">
	import Icon from '@iconify/svelte';
	import { scaleUtc } from 'd3-scale';
	import { curveStep } from 'd3-shape';
	import { Area, AreaChart, LinearGradient } from 'layerchart';
	import { PrometheusDriver, SampleValue } from 'prometheus-query';
	import { onDestroy, onMount } from 'svelte';

	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Card from '$lib/components/ui/card';
	import * as Chart from '$lib/components/ui/chart';
	import { m } from '$lib/paraglide/messages';

	let {
		prometheusDriver,
		scope,
		isReloading = $bindable()
	}: { prometheusDriver: PrometheusDriver; scope: string; isReloading: boolean } = $props();

	let ninety_five = $state([] as SampleValue[]);
	let ninety_nine = $state([] as SampleValue[]);
	const requestLatencies = $derived(
		ninety_five.map((sample, index) => ({
			time: sample.time,
			ninety_five: sample.value && !isNaN(sample.value) ? sample.value : 0,
			ninety_nine:
				ninety_nine[index] && ninety_nine[index].value && !isNaN(ninety_nine[index].value)
					? ninety_nine[index].value
					: 0
		}))
	);

	const configuration = {
		ninety_five: { label: '95th', color: 'var(--chart-1)' },
		ninety_nine: { label: '99th', color: 'var(--chart-2)' }
	} satisfies Chart.ChartConfig;

	async function fetchNinetyFive() {
		try {
			const response = await prometheusDriver.rangeQuery(
				`histogram_quantile(0.95, sum by(le) (rate(vllm:e2e_request_latency_seconds_bucket{juju_model="${scope}"}[5m])))`,
				Date.now() - 24 * 60 * 60 * 1000,
				Date.now(),
				2 * 60
			);
			ninety_five = response.result[0]?.values ?? [];
		} catch (error) {
			console.error(`Fail to fetch latest running requests in scope ${scope}:`, error);
		}
	}

	async function fetchNinetyNine() {
		try {
			const response = await prometheusDriver.rangeQuery(
				`histogram_quantile(0.99, sum by(le) (rate(vllm:e2e_request_latency_seconds_bucket{juju_model="${scope}"}[5m])))`,
				Date.now() - 24 * 60 * 60 * 1000,
				Date.now(),
				2 * 60
			);
			ninety_nine = response.result[0]?.values ?? [];
		} catch (error) {
			console.error(`Fail to fetch latest waiting requests in scope ${scope}:`, error);
		}
	}

	async function fetch() {
		try {
			await Promise.all([fetchNinetyFive(), fetchNinetyNine()]);
		} catch (error) {
			console.error(`Fail to fetch requests data in scope ${scope}:`, error);
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

<Card.Root class="h-full">
	<Card.Header>
		<Card.Title>{m.requests()}</Card.Title>
		<Card.Description>
			{m.llm_dashboard_requests_tooltip()}
		</Card.Description>
	</Card.Header>
	{#if !isLoaded}
		<Card.Content>
			<div class="flex h-[200px] w-full items-center justify-center">
				<Icon icon="svg-spinners:6-dots-rotate" class="size-12" />
			</div>
		</Card.Content>
	{:else if requestLatencies.length === 0}
		<Card.Content>
			<div class="flex h-[200px] w-full flex-col items-center justify-center">
				<Icon icon="ph:chart-line-fill" class="size-50 animate-pulse text-muted-foreground" />
				<p class="text-base text-muted-foreground">{m.no_data_display()}</p>
			</div>
		</Card.Content>
	{:else}
		<Card.Content>
			<Chart.Container config={configuration} class="h-[200px] w-full">
				<AreaChart
					data={requestLatencies}
					x="time"
					xScale={scaleUtc()}
					yPadding={[0, 25]}
					series={[
						{
							key: 'ninety_five',
							label: configuration.ninety_five.label,
							color: configuration.ninety_five.color
						},
						{
							key: 'ninety_nine',
							label: configuration.ninety_nine.label,
							color: configuration.ninety_nine.color
						}
					]}
					props={{
						area: {
							curve: curveStep,
							'fill-opacity': 0.4,
							line: { class: 'stroke-1' },
							motion: 'tween'
						},
						xAxis: {
							format: (v: Date) =>
								`${v.getHours().toString().padStart(2, '0')}:${v.getMinutes().toString().padStart(2, '0')}`
						},
						yAxis: { format: () => '' }
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
									{#if value}
										<p class="font-mono">{Number(value).toFixed(2)} {m.second()}</p>
									{:else}
										NaN
									{/if}
								</div>
							{/snippet}
						</Chart.Tooltip>
					{/snippet}
					{#snippet marks({ series, getAreaProps })}
						{#each series as s, i (s.key)}
							<LinearGradient
								stops={[s.color ?? '', 'color-mix(in lch, ' + s.color + ' 10%, transparent)']}
								vertical
							>
								{#snippet children({ gradient })}
									<Area {...getAreaProps(s, i)} fill={gradient} />
								{/snippet}
							</LinearGradient>
						{/each}
					{/snippet}
				</AreaChart>
			</Chart.Container>
		</Card.Content>
	{/if}
</Card.Root>
