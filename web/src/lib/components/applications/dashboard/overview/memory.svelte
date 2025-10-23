<script lang="ts">
	import { scaleUtc } from 'd3-scale';
	import { curveNatural } from 'd3-shape';
	import { Area, AreaChart, LinearGradient } from 'layerchart';
	import { PrometheusDriver, SampleValue } from 'prometheus-query';
	import { onDestroy, onMount } from 'svelte';

	import type { Scope } from '$lib/api/scope/v1/scope_pb';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Card from '$lib/components/ui/card';
	import * as Chart from '$lib/components/ui/chart/index.js';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { formatCapacity } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';

	let {
		prometheusDriver,
		scope,
		isReloading = $bindable(),
	}: { prometheusDriver: PrometheusDriver; scope: Scope; isReloading: boolean } = $props();

	let memoryUsages: SampleValue[] = $state([]);
	const memoryUsagesConfiguration = {
		usage: { label: 'Usage', color: 'var(--chart-1)' },
	} satisfies Chart.ChartConfig;
	let allocatableNodesMemory = $state(0);
	let memoryRequests = $state(0);
	let memoryLimits = $state(0);
	const { value: allocatableNodesMemoryValue, unit: allocatableNodesMemoryUnit } = $derived(
		formatCapacity(allocatableNodesMemory),
	);
	const { value: memoryRequestsValue, unit: memoryRequestsUnit } = $derived(formatCapacity(memoryRequests));
	const { value: memoryLimitsValue, unit: memoryLimitsUnit } = $derived(formatCapacity(memoryLimits));

	function fetch() {
		prometheusDriver
			.rangeQuery(
				`
				sum(
				container_memory_rss{container!="",job="kubelet",juju_model_uuid="${scope.uuid}",metrics_path="/metrics/cadvisor"}
				)
						`,
				Date.now() - 60 * 60 * 1000,
				Date.now(),
				2 * 60,
			)
			.then((response) => {
				memoryUsages = response.result[0].values;
			});
		prometheusDriver
			.instantQuery(
				`
				sum(
					kube_node_status_allocatable{job="kube-state-metrics",juju_model_uuid="${scope.uuid}",resource="memory"}
				)
				`,
			)
			.then((response) => {
				allocatableNodesMemory = response.result[0].value.value;
			});
		prometheusDriver
			.instantQuery(
				`
				sum(
					namespace_memory:kube_pod_container_resource_requests:sum{juju_model_uuid="${scope.uuid}"}
				)
				`,
			)
			.then((response) => {
				memoryRequests = response.result[0].value.value;
			});
		prometheusDriver
			.instantQuery(
				`
				sum(
					namespace_memory:kube_pod_container_resource_limits:sum{juju_model_uuid="${scope.uuid}"}
				)
				`,
			)
			.then((response) => {
				memoryLimits = response.result[0].value.value;
			});
	}

	const reloadManager = new ReloadManager(fetch);

	let isLoading = $state(true);
	onMount(async () => {
		await fetch();
		isLoading = false;
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

{#if isLoading}
	Loading
{:else}
	<Card.Root class="h-full gap-2">
		<Card.Header>
			<Card.Title>{m.memory_usage()}</Card.Title>
			<Card.Action class="text-muted-foreground flex flex-col gap-0.5 text-sm">
				<div class="flex justify-between gap-2">
					<p>{m.requests()}</p>
					<Tooltip.Provider>
						<Tooltip.Root>
							<Tooltip.Trigger>
								<p class="font-mono">{Math.round((memoryRequests * 100) / allocatableNodesMemory)}%</p>
							</Tooltip.Trigger>
							<Tooltip.Content>
								{memoryRequestsValue}
								{memoryRequestsUnit} / {allocatableNodesMemoryValue}
								{allocatableNodesMemoryUnit}
							</Tooltip.Content>
						</Tooltip.Root>
					</Tooltip.Provider>
				</div>
				<div class="flex justify-between gap-2">
					<p>{m.limits()}</p>
					<Tooltip.Provider>
						<Tooltip.Root>
							<Tooltip.Trigger>
								<p class="font-mono">{Math.round((memoryLimits * 100) / allocatableNodesMemory)}%</p>
							</Tooltip.Trigger>
							<Tooltip.Content>
								{memoryLimitsValue}
								{memoryLimitsUnit} / {allocatableNodesMemoryValue}
								{allocatableNodesMemoryUnit}
							</Tooltip.Content>
						</Tooltip.Root>
					</Tooltip.Provider>
				</div>
			</Card.Action>
		</Card.Header>
		<Card.Content>
			<Chart.Container config={memoryUsagesConfiguration}>
				<AreaChart
					data={memoryUsages}
					x="time"
					xScale={scaleUtc()}
					yPadding={[0, 25]}
					series={[
						{
							key: 'value',
							label: memoryUsagesConfiguration.usage.label,
							color: memoryUsagesConfiguration.usage.color,
						},
					]}
					props={{
						area: {
							curve: curveNatural,
							'fill-opacity': 0.4,
							line: { class: 'stroke-1' },
							motion: 'tween',
						},
						xAxis: {
							format: (v: Date) =>
								`${v.getHours().toString().padStart(2, '0')}:${v.getMinutes().toString().padStart(2, '0')}`,
						},
						yAxis: { format: () => '' },
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
									minute: 'numeric',
								});
							}}
						>
							{#snippet formatter({ item, name, value })}
								{@const { value: capacity, unit } = formatCapacity(Number(value))}
								<div
									style="--color-bg: {item.color}"
									class="aspect-square h-full w-fit shrink-0 border-(--color-border) bg-(--color-bg)"
								></div>
								<div class="flex flex-1 shrink-0 items-center justify-between text-xs leading-none">
									<div class="grid gap-1.5">
										<span class="text-muted-foreground">{name}</span>
									</div>
									<p class="font-mono">{capacity} {unit}</p>
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
