<script lang="ts">
	import Icon from '@iconify/svelte';
	import { scaleUtc } from 'd3-scale';
	import { curveMonotoneX } from 'd3-shape';
	import { Area, AreaChart, LinearGradient } from 'layerchart';
	import { PrometheusDriver, SampleValue } from 'prometheus-query';
	import { onDestroy, onMount } from 'svelte';
	import { SvelteDate } from 'svelte/reactivity';

	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Card from '$lib/components/ui/card';
	import * as Chart from '$lib/components/ui/chart/index.js';
	import { m } from '$lib/paraglide/messages';

	let {
		prometheusDriver,
		scope,
		isReloading = $bindable()
	}: { prometheusDriver: PrometheusDriver; scope: string; isReloading: boolean } = $props();

	let gpuUtilizations = $state([] as SampleValue[]);

	const gpuUtilizationsConfigurations = {
		usage: { label: 'Usage', color: 'var(--chart-1)' }
	} satisfies Chart.ChartConfig;

	async function fetchGPUUtilizations() {
		const response = await prometheusDriver.rangeQuery(
			`
			sum(avg(DCGM_FI_PROF_GR_ENGINE_ACTIVE{juju_model="${scope}"}) by (Hostname, gpu)) / count(count(DCGM_FI_PROF_GR_ENGINE_ACTIVE{juju_model="${scope}"}) by (Hostname, device))
			`,
			new SvelteDate().setMinutes(0, 0, 0) - 1 * 60 * 60 * 1000,
			new SvelteDate().setMinutes(0, 0, 0),
			2 * 60
		);
		gpuUtilizations = response.result[0]?.values ?? [];
		console.log(gpuUtilizations);
	}

	async function fetch() {
		try {
			await Promise.all([fetchGPUUtilizations()]);
		} catch (error) {
			console.error('Error fetching GPU memory usage data:', error);
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
	<Card.Root>
		<Card.Header class="h-[42px]">
			<Card.Title>GPU Utilization</Card.Title>
			<!-- <Card.Description>{m.receive_and_transmit()}</Card.Description> -->
		</Card.Header>
		<Card.Content>
			<div class="flex h-[230px] w-full items-center justify-center">
				<Icon icon="svg-spinners:6-dots-rotate" class="size-12" />
			</div>
		</Card.Content>
	</Card.Root>
{:else}
	<Card.Root class="h-full gap-2">
		<Card.Header>
			<Card.Title>GPU Utilization</Card.Title>
			<!-- <Card.Description>{m.receive_and_transmit()}</Card.Description> -->
		</Card.Header>
		<Card.Content>
			<Chart.Container class="h-[230px] w-full px-2 pt-2" config={gpuUtilizationsConfigurations}>
				<AreaChart
					data={gpuUtilizations}
					x="time"
					xScale={scaleUtc()}
					yPadding={[0, 25]}
					series={[
						{
							key: 'usage',
							label: gpuUtilizationsConfigurations.usage.label,
							color: gpuUtilizationsConfigurations.usage.color
						}
					]}
					props={{
						area: {
							curve: curveMonotoneX,
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
						<Chart.Tooltip
							indicator="dot"
							labelFormatter={(v: Date) => {
								return v.toLocaleDateString('en-US', {
									year: 'numeric',
									month: 'short',
									day: 'numeric',
									hour: 'numeric',
									minute: 'numeric'
								});
							}}
						>
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
	</Card.Root>
{/if}
