<script lang="ts">
	import Icon from '@iconify/svelte';
	import { scaleUtc } from 'd3-scale';
	import { curveNatural } from 'd3-shape';
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

	const configuration = {
		usage: { label: 'Usage', color: 'var(--chart-2)' }
	} satisfies Chart.ChartConfig;

	let cpuUsages: SampleValue[] = $state([]);
	async function fetchCPUUsage() {
		const response = await prometheusDriver.rangeQuery(
			`avg(rate(kubevirt_vmi_cpu_usage_seconds_total{juju_model="${scope}"}[5m]))`,
			new SvelteDate().setMinutes(0, 0, 0) - 60 * 60 * 1000,
			new SvelteDate().setMinutes(0, 0, 0),
			2 * 60
		);
		cpuUsages = response.result[0]?.values ?? [];
	}

	let cpuWait: SampleValue | undefined = $state(undefined);
	async function fetchCPUWait() {
		const response = await prometheusDriver.instantQuery(
			`avg(rate(kubevirt_vmi_vcpu_wait_seconds_total{juju_model="${scope}"}[5m]))`
		);
		cpuWait = response.result[0]?.value ?? undefined;
	}

	let cpuDelay: SampleValue | undefined = $state(undefined);
	async function fetchCPUDelay() {
		const response = await prometheusDriver.instantQuery(
			`avg(rate(kubevirt_vmi_vcpu_delay_seconds_total{juju_model="${scope}"}[5m]))`
		);
		cpuDelay = response.result[0]?.value ?? undefined;
	}

	let isLoaded = $state(false);
	async function fetch() {
		try {
			await Promise.all([fetchCPUUsage(), fetchCPUWait(), fetchCPUDelay()]);
		} catch (error) {
			console.error('Failed to fetch cpu data:', error);
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

	onMount(async () => {
		await fetch();
		isLoaded = true;
	});
	onDestroy(() => {
		reloadManager.stop();
	});
</script>

<Card.Root class="h-full gap-2">
	<Card.Header>
		<Card.Title>{m.cpu_usage()}</Card.Title>
		<Card.Action class="flex flex-col gap-0.5 text-sm text-muted-foreground">
			{#if cpuWait}
				<div class="flex justify-between gap-2">
					<p>{m.wait()}</p>
					<p class="font-mono">{(Number(cpuWait.value) * 100).toFixed(2)}%</p>
				</div>
			{/if}
			{#if cpuDelay}
				<div class="flex justify-between gap-2">
					<p>{m.delay()}</p>
					<p class="font-mono">{(Number(cpuDelay.value) * 100).toFixed(2)}%</p>
				</div>
			{/if}
		</Card.Action>
	</Card.Header>
	<Card.Content class="h-full">
		{#if !isLoaded}
			<div class="flex h-full w-full items-center justify-center border">
				<Icon icon="svg-spinners:6-dots-rotate" class="size-24" />
			</div>
		{:else if !cpuUsages?.length}
			<div class="flex h-full w-full flex-col items-center justify-center">
				<Icon icon="ph:chart-line-fill" class="size-60 animate-pulse text-muted-foreground" />
				<p class="text-base text-muted-foreground">{m.no_data_display()}</p>
			</div>
		{:else}
			<Chart.Container config={configuration}>
				<AreaChart
					data={cpuUsages}
					x="time"
					xScale={scaleUtc()}
					yPadding={[0, 25]}
					series={[
						{
							key: 'value',
							label: configuration.usage.label,
							color: configuration.usage.color
						}
					]}
					props={{
						area: {
							curve: curveNatural,
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
								<div class="flex flex-1 shrink-0 items-center justify-between text-xs leading-none">
									<div class="grid gap-1.5">
										<span class="text-muted-foreground">{name}</span>
									</div>
									<p class="font-mono">{(Number(value) * 100).toFixed(2)}%</p>
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
		{/if}
	</Card.Content>
</Card.Root>
