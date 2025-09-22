<script lang="ts">
	import { scaleBand } from 'd3-scale';
	import { BarChart, Highlight, type ChartContextValue } from 'layerchart';
	import { InstantVector, PrometheusDriver } from 'prometheus-query';
	import { onMount } from 'svelte';
	import { cubicInOut } from 'svelte/easing';

	import type { Scope } from '$lib/api/scope/v1/scope_pb';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Card from '$lib/components/ui/card';
	import * as Chart from '$lib/components/ui/chart/index.js';
	import { m } from '$lib/paraglide/messages';

	let {
		prometheusDriver,
		scope,
		isReloading = $bindable(),
	}: { prometheusDriver: PrometheusDriver; scope: Scope; isReloading: boolean } = $props();

	let allocatableGPUs: Record<string, number> = $state({});
	let occupiedGPUs: Record<string, number> = $state({});

	const usages = $derived(
		Object.keys(allocatableGPUs).map((node) => ({
			node,
			free: allocatableGPUs[node] ? (allocatableGPUs[node] - occupiedGPUs[node]) / allocatableGPUs[node] : 0,
			occupied: allocatableGPUs[node] ? occupiedGPUs[node] / allocatableGPUs[node] : 0,
		})),
	);

	const configuration = {
		occupied: { label: 'Occupied', color: 'var(--chart-1)' },
		free: { label: 'Free', color: 'var(--chart-2)' },
	} satisfies Chart.ChartConfig;

	let context = $state<ChartContextValue>();

	async function fetch() {
		await prometheusDriver
			.instantQuery(
				`
				sum by (node) (kube_node_status_allocatable{juju_model_uuid="${scope.uuid}",resource="split_nvidia_com_gpu"})
				`,
			)
			.then((response) => {
				const instanceVectors: InstantVector[] = response.result;
				allocatableGPUs = Object.fromEntries(
					instanceVectors.map((instanceVector) => [
						(instanceVector.metric.labels as { node?: string }).node,
						Number(instanceVector.value.value),
					]),
				);
			});
		await prometheusDriver
			.instantQuery(
				`
				sum by (node) (kube_pod_container_resource_limits{juju_model_uuid="${scope.uuid}",resource="split_nvidia_com_gpu"})
				`,
			)
			.then((response) => {
				const instanceVectors: InstantVector[] = response.result;
				occupiedGPUs = Object.fromEntries(
					instanceVectors.map((instanceVector) => [
						(instanceVector.metric.labels as { node?: string }).node,
						Number(instanceVector.value.value),
					]),
				);
			});
	}

	const reloadManager = new ReloadManager(fetch);

	let isLoading = $state(true);
	onMount(async () => {
		await fetch();
		isLoading = false;
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
	<Card.Root>
		<Card.Header>
			<Card.Title>{m.vgpu()}</Card.Title>
			<Card.Description>
				{m.vgpu_usage_description()}
			</Card.Description>
		</Card.Header>
		<Card.Content>
			<Chart.Container config={configuration} class="w-full">
				<BarChart
					bind:context
					data={usages}
					orientation="horizontal"
					yScale={scaleBand().padding(0.25)}
					y="node"
					axis="y"
					padding={{ right: 16 }}
					rule={false}
					series={[
						{
							key: 'free',
							label: configuration.free.label,
							color: configuration.free.color,
						},
						{
							key: 'occupied',
							label: configuration.occupied.label,
							color: configuration.occupied.color,
						},
					]}
					seriesLayout="stack"
					legend
					props={{
						bars: {
							stroke: 'none',
							initialY: context?.height,
							initialHeight: 0,
							motion: {
								y: { type: 'tween', duration: 500, easing: cubicInOut },
								height: { type: 'tween', duration: 500, easing: cubicInOut },
							},
						},
						highlight: { area: false },
						yAxis: {
							tickLabelProps: {
								textAnchor: 'start',
								dx: 6,
								class: 'stroke-none fill-background!',
							},
							tickLength: 0,
						},
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
									<p class="font-mono">{Number(value) * 100} %</p>
								</div>
							{/snippet}
						</Chart.Tooltip>
					{/snippet}
				</BarChart>
			</Chart.Container>
		</Card.Content>
	</Card.Root>
{/if}
