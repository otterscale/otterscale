<script lang="ts">
	import Icon from '@iconify/svelte';
	import { ArcChart, Text } from 'layerchart';
	import type { PrometheusDriver, SampleValue } from 'prometheus-query';
	import { onDestroy, onMount } from 'svelte';

	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Chart from '$lib/components/ui/chart/index.js';
	import { formatCapacity } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';

	let {
		prometheusDriver,
		scope,
		isReloading = $bindable()
	}: { prometheusDriver: PrometheusDriver; scope: string; isReloading: boolean } = $props();

	let memoryUsage: SampleValue | undefined = $state(undefined);
	async function fetchMemoryUsage() {
		const usageResponse = await prometheusDriver.instantQuery(
			`
			100 * (sum(node_memory_MemTotal_bytes{juju_model="${scope}", container!=""} - node_memory_MemAvailable_bytes{juju_model="${scope}", container!=""}))
			/
			sum(node_memory_MemTotal_bytes{juju_model="${scope}", container!=""})
			`
		);
		memoryUsage = usageResponse.result[0]?.value ?? undefined;
	}

	let memoryRequest: SampleValue | undefined = $state(undefined);
	async function fetchMemoryRequest() {
		const response = await prometheusDriver.instantQuery(
			`
			100 * sum(kube_pod_container_resource_requests{resource="memory", unit="byte", juju_model="${scope}", container!=""})
			/
			sum(kube_node_status_allocatable{cluster!="", juju_model="${scope}", resource="memory"})
			`
		);
		memoryRequest = response.result[0]?.value ?? undefined;
	}

	let allocatableMemory: SampleValue | undefined = $state(undefined);
	async function fetchAllocatableMemory() {
		const response = await prometheusDriver.instantQuery(
			`
			sum(kube_node_status_allocatable{cluster!="", juju_model="${scope}", resource="memory"})
			`
		);
		allocatableMemory = response.result[0]?.value?.value ?? undefined;
	}

	async function fetch() {
		try {
			await Promise.all([fetchMemoryUsage(), fetchMemoryRequest(), fetchAllocatableMemory()]);
		} catch (error) {
			console.error('Failed to fetch CPU usage:', error);
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

	const chartConfig = {
		usage: { label: 'Usage' },
		request: { label: 'Request' }
	} satisfies Chart.ChartConfig;
</script>

<Card.Root class="relative h-full min-h-[140px] gap-2 overflow-hidden">
	<Card.Header class="h-20">
		<Card.Title>{m.memory()}</Card.Title>
		<Card.Description class="flex">
			{m.cluster_dashboard_memory_description()}
		</Card.Description>
	</Card.Header>
	{#if !isLoaded}
		<div class="flex h-9 w-full items-center justify-center">
			<Icon icon="svg-spinners:6-dots-rotate" class="size-10" />
		</div>
	{:else if !memoryUsage}
		<div class="flex h-full w-full flex-col items-center justify-center">
			<Icon icon="ph:chart-bar-fill" class="size-6 animate-pulse text-muted-foreground" />
			<p class="p-0 text-xs text-muted-foreground">{m.no_data_display()}</p>
		</div>
	{:else}
		<Card.Content class="flex-1">
			<Chart.Container config={chartConfig} class="mx-auto aspect-square max-h-[250px]">
				<ArcChart
					value="value"
					outerRadius={-23}
					innerRadius={-13}
					padding={23}
					range={[180, -180]}
					maxValue={100}
					series={[
						{
							key: 'request',
							data: [{ key: 'request', ...memoryRequest }],
							color: 'var(--chart-2)'
						},
						{ key: 'usage', data: [{ key: 'usage', ...memoryUsage }], color: 'var(--chart-1)' }
					]}
					props={{
						arc: { track: { fill: 'var(--muted)' }, motion: 'tween' },
						tooltip: { context: { hideDelay: 350 } }
					}}
					tooltip={false}
				>
					{#snippet aboveMarks()}
						{@const { value, unit } = formatCapacity(Number(allocatableMemory))}
						<Text
							{value}
							textAnchor="middle"
							verticalAnchor="middle"
							class="fill-foreground text-3xl! font-bold"
						/>
						<Text
							value={unit}
							textAnchor="middle"
							verticalAnchor="middle"
							class="fill-foreground text-xl! font-bold"
							dy={30}
						/>
					{/snippet}
				</ArcChart>
			</Chart.Container>
			<Card.Footer class="mt-auto w-full">
				<div class="mx-auto grid w-fit grid-cols-2 py-2">
					<p class="col-start-1 row-start-1">
						<span class="mr-2 inline-block aspect-square size-3 bg-chart-1 align-middle"></span>
						usage
					</p>
					<p class="col-start-2 row-start-1 ml-auto">
						{Math.round(Number(memoryUsage?.value ?? 0))} %
					</p>
					<p class="col-start-1 row-start-2">
						<span class="mr-2 inline-block aspect-square size-3 bg-chart-2 align-middle"></span>
						request
					</p>
					<p class="col-start-2 row-start-2 ml-auto">
						{Math.round(Number(memoryRequest?.value ?? 0))} %
					</p>
				</div>
			</Card.Footer>
		</Card.Content>
	{/if}
</Card.Root>
