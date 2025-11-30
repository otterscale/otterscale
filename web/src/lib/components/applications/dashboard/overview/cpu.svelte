<script lang="ts">
	import Icon from '@iconify/svelte';
	import { scaleUtc } from 'd3-scale';
	import { curveMonotoneX } from 'd3-shape';
	import { Area, AreaChart, LinearGradient } from 'layerchart';
	import { PrometheusDriver, SampleValue } from 'prometheus-query';
	import { onDestroy, onMount } from 'svelte';

	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Card from '$lib/components/ui/card';
	import * as Chart from '$lib/components/ui/chart/index.js';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { m } from '$lib/paraglide/messages';

	let {
		prometheusDriver,
		scope,
		isReloading = $bindable()
	}: { prometheusDriver: PrometheusDriver; scope: string; isReloading: boolean } = $props();

	let cpuUsages: SampleValue[] = $state([]);
	const cpuUsagesConfiguration = {
		usage: { label: 'Usage', color: 'var(--chart-2)' }
	} satisfies Chart.ChartConfig;
	let allocatableNodesCPU = $state(0);
	let cpuRequests = $state(0);
	let cpuLimits = $state(0);

	async function fetchCPUUsages() {
		const response = await prometheusDriver.rangeQuery(
			`
			sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_irate{juju_model="${scope}"})
			`,
			Date.now() - 60 * 60 * 1000,
			Date.now(),
			2 * 60
		);
		cpuUsages = response.result[0]?.values ?? [];
	}

	async function fetchAllocatableNodesCPU() {
		const response = await prometheusDriver.instantQuery(
			`
			sum(kube_node_status_allocatable{job="kube-state-metrics",juju_model="${scope}",resource="cpu"})
			`
		);
		allocatableNodesCPU = response.result[0]?.value?.value;
	}

	async function fetchCPURequests() {
		const response = await prometheusDriver.instantQuery(
			`
			sum(namespace_cpu:kube_pod_container_resource_requests:sum{juju_model="${scope}"})
			`
		);
		cpuRequests = response.result[0]?.value?.value;
	}

	async function fetchCPULimits() {
		const response = await prometheusDriver.instantQuery(
			`
			sum(namespace_cpu:kube_pod_container_resource_limits:sum{juju_model="${scope}"})
			`
		);
		cpuLimits = response.result[0]?.value?.value;
	}

	async function fetch() {
		try {
			await Promise.all([
				fetchCPUUsages(),
				fetchAllocatableNodesCPU(),
				fetchCPURequests(),
				fetchCPULimits()
			]);
		} catch (error) {
			console.error('Failed to fetch CPU metrics:', error);
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
			<Card.Title>{m.cpu_usage()}</Card.Title>
		</Card.Header>
		<Card.Content>
			<div class="flex h-[200px] w-full items-center justify-center">
				<Icon icon="svg-spinners:6-dots-rotate" class="size-12" />
			</div>
		</Card.Content>
	</Card.Root>
{:else}
	<Card.Root class="h-full gap-2">
		<Card.Header>
			<Card.Title>{m.cpu_usage()}</Card.Title>
			<Card.Action class="flex flex-col gap-0.5 text-sm text-muted-foreground">
				<div class="flex justify-between gap-2">
					<p>{m.requests()}</p>
					<Tooltip.Provider>
						<Tooltip.Root>
							<Tooltip.Trigger>
								<p class="font-mono">{Math.round((cpuRequests * 100) / allocatableNodesCPU)}%</p>
							</Tooltip.Trigger>
							<Tooltip.Content>
								{cpuRequests.toFixed(2)} / {allocatableNodesCPU}
							</Tooltip.Content>
						</Tooltip.Root>
					</Tooltip.Provider>
				</div>
				<div class="flex justify-between gap-2">
					<p>{m.limits()}</p>
					<Tooltip.Provider>
						<Tooltip.Root>
							<Tooltip.Trigger>
								<p class="font-mono">{Math.round((cpuLimits * 100) / allocatableNodesCPU)}%</p>
							</Tooltip.Trigger>
							<Tooltip.Content>
								{cpuLimits.toFixed(2)} / {allocatableNodesCPU}
							</Tooltip.Content>
						</Tooltip.Root>
					</Tooltip.Provider>
				</div>
			</Card.Action>
		</Card.Header>
		<Card.Content>
			<Chart.Container class="h-[200px] w-full px-2 pt-2" config={cpuUsagesConfiguration}>
				<AreaChart
					data={cpuUsages}
					x="time"
					xScale={scaleUtc()}
					yPadding={[0, 25]}
					series={[
						{
							key: 'value',
							label: cpuUsagesConfiguration.usage.label,
							color: cpuUsagesConfiguration.usage.color
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
		</Card.Content>
	</Card.Root>
{/if}
