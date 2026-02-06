<script lang="ts">
	import Icon from '@iconify/svelte';
	import { ArcChart } from 'layerchart';
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

	let gpuUtilization: SampleValue | undefined = $state(undefined);
	async function fetchMemoryUsage() {
		const usageResponse = await prometheusDriver.instantQuery(
			`
			avg(sum(Device_utilization_desc_of_container{juju_model="${scope}"}) by (deviceuuid, vdeviceid, podname, podnamespace))
			`
		);
		gpuUtilization = usageResponse.result[0]?.value ?? undefined;
	}

	async function fetch() {
		try {
			await Promise.all([fetchMemoryUsage()]);
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
	<Icon
		icon="ph:graphics-card"
		class="absolute -right-10 bottom-0 -z-0 size-36 text-8xl tracking-tight text-nowrap text-primary/5 uppercase group-hover:hidden"
	/>
	<Card.Header>
		<Card.Title>GPU Utilization</Card.Title>
		<Card.Description class="z-10 flex flex-col items-end">
			<p>utilization: {Math.round(Number(gpuUtilization?.value ?? 0) * 100)} %</p>
		</Card.Description>
	</Card.Header>
	{#if !isLoaded}
		<div class="flex h-9 w-full items-center justify-center">
			<Icon icon="svg-spinners:6-dots-rotate" class="size-10" />
		</div>
	{:else if !gpuUtilization}
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
					maxValue={1}
					series={[
						{ key: 'usage', data: [{ key: 'usage', ...gpuUtilization }], color: 'var(--chart-1)' }
					]}
					props={{
						arc: { track: { fill: 'var(--muted)' }, motion: 'tween' },
						tooltip: { context: { hideDelay: 350 } }
					}}
				>
					{#snippet tooltip()}
						<Chart.Tooltip hideLabel nameKey="key" />
					{/snippet}
				</ArcChart>
			</Chart.Container>
		</Card.Content>
	{/if}
</Card.Root>
