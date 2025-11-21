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
	import { formatIO } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';

	let {
		prometheusDriver,
		scope,
		isReloading = $bindable()
	}: { prometheusDriver: PrometheusDriver; scope: string; isReloading: boolean } = $props();

	const configuration = {
		read: { label: 'Read', color: 'var(--chart-1)' },
		write: { label: 'Write', color: 'var(--chart-2)' }
	} satisfies Chart.ChartConfig;

	let reads: SampleValue[] = $state([] as SampleValue[]);
	async function fetchReads() {
		const response = await prometheusDriver.rangeQuery(
			`avg(rate(kubevirt_vmi_storage_read_traffic_bytes_total{juju_model="${scope}"}[5m]))`,
			new SvelteDate().setMinutes(0, 0, 0) - 60 * 60 * 1000,
			new SvelteDate().setMinutes(0, 0, 0),
			2 * 60
		);
		reads = response.result[0]?.values ?? [];
	}

	let writes: SampleValue[] = $state([] as SampleValue[]);
	async function fetchWrites() {
		const response = await prometheusDriver.rangeQuery(
			`avg(rate(kubevirt_vmi_storage_write_traffic_bytes_total{juju_model="${scope}"}[5m]))`,
			new SvelteDate().setMinutes(0, 0, 0) - 60 * 60 * 1000,
			new SvelteDate().setMinutes(0, 0, 0),
			2 * 60
		);
		writes = response.result[0]?.values ?? [];
	}

	let isLoaded = $state(false);
	async function fetchData() {
		try {
			await Promise.all([fetchReads(), fetchWrites()]);
			isLoaded = true;
		} catch (error) {
			console.error('Failed to fetch network data:', error);
		}
	}
	const reloadManager = new ReloadManager(fetchData);

	const throughputs = $derived(
		reads.map((sample, index) => ({
			time: sample.time,
			read: sample.value,
			write: writes[index]?.value ?? 0
		}))
	);

	$effect(() => {
		if (isReloading) {
			reloadManager.restart();
		} else {
			reloadManager.stop();
		}
	});

	onMount(async () => {
		await fetchData();
	});
	onDestroy(() => {
		reloadManager.stop();
	});
</script>

<Card.Root class="h-full gap-2">
	<Card.Header>
		<Card.Title>{m.storage_throughPut()}</Card.Title>
		<Card.Description>{m.read_and_write()}</Card.Description>
	</Card.Header>
	<Card.Content class="h-full">
		{#if !isLoaded}
			<div class="flex h-full w-full items-center justify-center border">
				<Icon icon="svg-spinners:6-dots-rotate" class="size-24" />
			</div>
		{:else if !reads?.length || !writes?.length}
			<div class="flex h-full w-full flex-col items-center justify-center">
				<Icon icon="ph:chart-line-fill" class="size-60 animate-pulse text-muted-foreground" />
				<p class="text-base text-muted-foreground">{m.no_data_display()}</p>
			</div>
		{:else}
			<Chart.Container config={configuration}>
				<AreaChart
					data={throughputs}
					x="time"
					xScale={scaleUtc()}
					yPadding={[0, 25]}
					series={[
						{
							key: 'read',
							label: configuration.read.label,
							color: configuration.read.color
						},
						{
							key: 'write',
							label: configuration.write.label,
							color: configuration.write.color
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
								{@const { value: io, unit } = formatIO(Number(value))}
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
									<p class="font-mono">{io} {unit}</p>
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
