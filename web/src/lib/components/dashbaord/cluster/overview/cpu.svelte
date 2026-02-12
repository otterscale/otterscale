<script lang="ts">
	import Icon from '@iconify/svelte';
	import { ArcChart, Text } from 'layerchart';
	import type { PrometheusDriver, SampleValue } from 'prometheus-query';
	import { onDestroy, onMount } from 'svelte';

	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Chart from '$lib/components/ui/chart/index.js';
	import { m } from '$lib/paraglide/messages';

	let {
		prometheusDriver,
		scope,
		isReloading = $bindable()
	}: { prometheusDriver: PrometheusDriver; scope: string; isReloading: boolean } = $props();

	let cpuUsage: SampleValue | undefined = $state(undefined);
	async function fetchCPUUsage() {
		const response = await prometheusDriver.instantQuery(
			`100 * sum(rate(node_cpu_seconds_total{mode!="idle", juju_model="${scope}", container!=""}[5m])) / sum(rate(node_cpu_seconds_total{juju_model="${scope}", container!=""}[5m]))`
		);
		cpuUsage = response.result[0]?.value ?? undefined;
	}

	let cpuRequest: SampleValue | undefined = $state(undefined);
	async function fetchCPURequest() {
		const response = await prometheusDriver.instantQuery(
			`
			100 * sum(kube_pod_container_resource_requests{resource="cpu", unit="core", juju_model="${scope}", container!=""})
			/
			sum (kube_node_status_allocatable{cluster!="", juju_model="${scope}", resource="cpu"})
			`
		);
		cpuRequest = response.result[0]?.value ?? undefined;
	}

	let allocatableCPU: SampleValue | undefined = $state(undefined);
	async function fetchAllocatableCPU() {
		const response = await prometheusDriver.instantQuery(
			`
			sum(kube_node_status_allocatable{cluster!="", juju_model="${scope}", resource="cpu"})
			`
		);
		allocatableCPU = response.result[0]?.value?.value ?? undefined;
	}

	async function fetch() {
		try {
			await Promise.all([fetchCPUUsage(), fetchCPURequest(), fetchAllocatableCPU()]);
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
		<Card.Title>{m.cpu()}</Card.Title>
		<Card.Description class="flex">
			{m.cluster_dashboard_cpu_description()}
		</Card.Description>
	</Card.Header>
	{#if !isLoaded}
		<div class="flex h-9 w-full items-center justify-center">
			<Icon icon="svg-spinners:6-dots-rotate" class="size-10" />
		</div>
	{:else if !cpuUsage}
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
						{ key: 'request', data: [{ key: 'request', ...cpuRequest }], color: 'var(--chart-2)' },
						{ key: 'usage', data: [{ key: 'usage', ...cpuUsage }], color: 'var(--chart-1)' }
					]}
					props={{
						arc: { track: { fill: 'var(--muted)' }, motion: 'tween' }
					}}
					tooltip={false}
				>
					{#snippet aboveMarks()}
						<Text
							value={String(allocatableCPU)}
							textAnchor="middle"
							verticalAnchor="middle"
							class="fill-foreground text-3xl! font-bold"
						/>
						<Text
							value="core"
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
						{Math.round(Number(cpuUsage?.value ?? 0))} %
					</p>
					<p class="col-start-1 row-start-2">
						<span class="mr-2 inline-block aspect-square size-3 bg-chart-2 align-middle"></span>
						request
					</p>
					<p class="col-start-2 row-start-2 ml-auto">
						{Math.round(Number(cpuRequest?.value ?? 0))} %
					</p>
				</div>
			</Card.Footer>
		</Card.Content>
	{/if}
</Card.Root>
