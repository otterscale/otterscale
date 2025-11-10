<script lang="ts">
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

	let receives = $state([] as SampleValue[]);
	let transmits = $state([] as SampleValue[]);
	const traffics = $derived(
		receives.map((sample, index) => ({
			time: sample.time,
			receive: sample.value,
			transmit: transmits[index]?.value ?? 0
		}))
	);
	const trafficsConfigurations = {
		receive: { label: 'Receive', color: 'var(--chart-1)' },
		transmit: { label: 'Transmit', color: 'var(--chart-2)' }
	} satisfies Chart.ChartConfig;

	function fetch() {
		prometheusDriver
			.rangeQuery(
				`
						sum(
						irate(
							container_network_receive_bytes_total{job="kubelet",juju_model="${scope}",metrics_path="/metrics/cadvisor",namespace=~".+"}[4m]
						)
						)
						`,
				new SvelteDate().setMinutes(0, 0, 0) - 1 * 60 * 60 * 1000,
				new SvelteDate().setMinutes(0, 0, 0),
				2 * 60
			)
			.then((response) => {
				receives = response.result[0]?.values ?? [];
			});
		prometheusDriver
			.rangeQuery(
				`
						sum(
						irate(
							container_network_transmit_bytes_total{job="kubelet",juju_model="${scope}",metrics_path="/metrics/cadvisor",namespace=~".+"}[4m]
						)
						)
						`,
				new SvelteDate().setMinutes(0, 0, 0) - 1 * 60 * 60 * 1000,
				new SvelteDate().setMinutes(0, 0, 0),
				2 * 60
			)
			.then((response) => {
				transmits = response.result[0]?.values ?? [];
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
			<Card.Title>{m.network_bandwidth()}</Card.Title>
			<Card.Description>{m.receive_and_transmit()}</Card.Description>
		</Card.Header>
		<Card.Content>
			<Chart.Container config={trafficsConfigurations}>
				<AreaChart
					data={traffics}
					x="time"
					xScale={scaleUtc()}
					yPadding={[0, 25]}
					series={[
						{
							key: 'receive',
							label: trafficsConfigurations.receive.label,
							color: trafficsConfigurations.receive.color
						},
						{
							key: 'transmit',
							label: trafficsConfigurations.transmit.label,
							color: trafficsConfigurations.transmit.color
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
		</Card.Content>
	</Card.Root>
{/if}
