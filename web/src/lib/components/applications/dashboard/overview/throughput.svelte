<script lang="ts">
	import { scaleUtc } from 'd3-scale';
	import { curveNatural } from 'd3-shape';
	import { Area, AreaChart, LinearGradient } from 'layerchart';
	import { PrometheusDriver, SampleValue } from 'prometheus-query';
	import { onMount } from 'svelte';
	import { SvelteDate } from 'svelte/reactivity';

	import type { Scope } from '$lib/api/scope/v1/scope_pb';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Card from '$lib/components/ui/card';
	import * as Chart from '$lib/components/ui/chart/index.js';
	import { formatIO } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';

	let {
		prometheusDriver,
		scope,
		isReloading = $bindable(),
	}: { prometheusDriver: PrometheusDriver; scope: Scope; isReloading: boolean } = $props();

	let reads = $state([] as SampleValue[]);
	let writes = $state([] as SampleValue[]);
	const throughputs = $derived(
		reads.map((sample, index) => ({
			time: sample.time,
			read: sample.value,
			write: writes[index]?.value ?? 0,
		})),
	);
	const throughputsConfigurations = {
		read: { label: 'Read', color: 'var(--chart-1)' },
		write: { label: 'Write', color: 'var(--chart-2)' },
	} satisfies Chart.ChartConfig;
	function fetch() {
		prometheusDriver
			.rangeQuery(
				`
						sum(
						rate(
							container_fs_reads_bytes_total{container!="",device=~"(/dev/)?(mmcblk.p.+|nvme.+|rbd.+|sd.+|vd.+|xvd.+|dm-.+|md.+|dasd.+)",job="kubelet",juju_model_uuid="${scope.uuid}",metrics_path="/metrics/cadvisor",namespace!=""}[4m]
						)
						)
						`,
				new SvelteDate().setMinutes(0, 0, 0) - 1 * 60 * 60 * 1000,
				new SvelteDate().setMinutes(0, 0, 0),
				2 * 60,
			)
			.then((response) => {
				reads = response.result[0].values;
			});
		prometheusDriver
			.rangeQuery(
				`
						sum(
						rate(
							container_fs_writes_bytes_total{container!="",job="kubelet",juju_model_uuid="${scope.uuid}",metrics_path="/metrics/cadvisor",namespace!=""}[4m]
						)
						)
						`,
				new SvelteDate().setMinutes(0, 0, 0) - 1 * 60 * 60 * 1000,
				new SvelteDate().setMinutes(0, 0, 0),
				2 * 60,
			)
			.then((response) => {
				writes = response.result[0].values;
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
	<Card.Root class="h-full gap-2">
		<Card.Header>
			<Card.Title>{m.storage_throughPut()}</Card.Title>
			<Card.Description>{m.read_and_write()}</Card.Description>
		</Card.Header>
		<Card.Content>
			<Chart.Container config={throughputsConfigurations}>
				<AreaChart
					data={throughputs}
					x="time"
					xScale={scaleUtc()}
					yPadding={[0, 25]}
					series={[
						{
							key: 'read',
							label: throughputsConfigurations.read.label,
							color: throughputsConfigurations.read.color,
						},
						{
							key: 'write',
							label: throughputsConfigurations.write.label,
							color: throughputsConfigurations.write.color,
						},
					]}
					seriesLayout="stack"
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
