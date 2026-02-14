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
	import { formatIO } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';

	let {
		prometheusDriver,
		scope,
		isReloading = $bindable()
	}: { prometheusDriver: PrometheusDriver; scope: string; isReloading: boolean } = $props();

	const configuration = {
		receive: { label: 'Receive', color: 'var(--chart-1)' },
		transmit: { label: 'Transmit', color: 'var(--chart-2)' }
	} satisfies Chart.ChartConfig;

	let receives: SampleValue[] = $state([]);
	async function fetchReceives() {
		const response = await prometheusDriver.rangeQuery(
			`avg(rate(kubevirt_vmi_network_receive_bytes_total{juju_model="${scope}"}[5m]))`,
			new SvelteDate().setMinutes(0, 0, 0) - 60 * 60 * 1000,
			new SvelteDate().setMinutes(0, 0, 0),
			2 * 60
		);
		receives = response.result[0]?.values ?? [];
	}

	let transmits: SampleValue[] = $state([]);
	async function fetchTransmits() {
		const response = await prometheusDriver.rangeQuery(
			`avg(rate(kubevirt_vmi_network_transmit_bytes_total{juju_model="${scope}"}[5m]))`,
			new SvelteDate().setMinutes(0, 0, 0) - 60 * 60 * 1000,
			new SvelteDate().setMinutes(0, 0, 0),
			2 * 60
		);
		transmits = response.result[0]?.values ?? [];
	}

	async function fetch() {
		try {
			await Promise.all([fetchReceives(), fetchTransmits()]);
		} catch (error) {
			console.error('Failed to fetch network data:', error);
		}
	}

	const reloadManager = new ReloadManager(fetch);

	const traffics = $derived(
		receives && transmits
			? receives.map((sample, index) => ({
					time: sample.time,
					receive: sample.value,
					transmit: transmits[index]?.value ?? null
				}))
			: []
	);

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
</script>

{#if !isLoaded}
	<Card.Root>
		<Card.Header class="h-[42px]">
			<Card.Title>{m.network_bandwidth()}</Card.Title>
			<Card.Description>{m.receive_and_transmit()}</Card.Description>
		</Card.Header>
		<Card.Content>
			<div class="flex h-[230px] w-full items-center justify-center">
				<Icon icon="svg-spinners:6-dots-rotate" class="size-12" />
			</div>
		</Card.Content>
	</Card.Root>
{:else if !receives?.length || !transmits?.length}
	<Card.Root>
		<Card.Header class="h-[42px]">
			<Card.Title>{m.network_bandwidth()}</Card.Title>
			<Card.Description>{m.receive_and_transmit()}</Card.Description>
		</Card.Header>
		<Card.Content>
			<div class="flex h-[230px] w-full flex-col items-center justify-center">
				<Icon icon="ph:chart-line-fill" class="size-50 animate-pulse text-muted-foreground" />
				<p class="text-base text-muted-foreground">{m.no_data_display()}</p>
			</div>
		</Card.Content>
	</Card.Root>
{:else}
	<Card.Root>
		<Card.Header>
			<Card.Title>{m.network_bandwidth()}</Card.Title>
			<Card.Description>{m.receive_and_transmit()}</Card.Description>
		</Card.Header>
		<Card.Content class="h-full">
			<Chart.Container class="h-[230px] w-full px-2 pt-2" config={configuration}>
				<AreaChart
					data={traffics}
					x="time"
					xScale={scaleUtc()}
					yPadding={[0, 25]}
					series={[
						{
							key: 'receive',
							label: configuration.receive.label,
							color: configuration.receive.color
						},
						{
							key: 'transmit',
							label: configuration.transmit.label,
							color: configuration.transmit.color
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
		</Card.Content>
	</Card.Root>
{/if}
