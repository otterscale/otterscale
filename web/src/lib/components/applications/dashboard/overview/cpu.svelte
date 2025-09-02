<script lang="ts">
	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Card from '$lib/components/ui/card';
	import * as Chart from '$lib/components/ui/chart/index.js';
	import { m } from '$lib/paraglide/messages';
	import { currentKubernetes } from '$lib/stores';
	import { scaleUtc } from 'd3-scale';
	import { curveNatural } from 'd3-shape';
	import { Area, AreaChart, LinearGradient } from 'layerchart';
	import { PrometheusDriver, SampleValue } from 'prometheus-query';
	import { onMount } from 'svelte';

	let { prometheusDriver, isReloading = $bindable() }: { prometheusDriver: PrometheusDriver; isReloading: boolean } =
		$props();

	let cpuUsages: SampleValue[] = $state([]);
	const cpuUsagesConfiguration = {
		usage: { label: 'Usage', color: 'var(--chart-2)' },
	} satisfies Chart.ChartConfig;
	let cpuRequests = $state(0);
	let cpuLimits = $state(0);

	function fetch() {
		prometheusDriver
			.rangeQuery(
				`
						sum(
						node_namespace_pod_container:container_cpu_usage_seconds_total:sum_irate{juju_model_uuid="${$currentKubernetes?.scopeUuid}"}
						)
						`,
				Date.now() - 60 * 60 * 1000,
				Date.now(),
				2 * 60,
			)
			.then((response) => {
				cpuUsages = response.result[0].values;
			});
		prometheusDriver
			.instantQuery(
				`
						sum(
							namespace_cpu:kube_pod_container_resource_requests:sum{juju_model_uuid="${$currentKubernetes?.scopeUuid}"}
						)
						/
						sum(
							kube_node_status_allocatable{job="kube-state-metrics",juju_model_uuid="${$currentKubernetes?.scopeUuid}",resource="cpu"}
						)
						`,
			)
			.then((response) => {
				cpuRequests = response.result[0].value.value;
			});
		prometheusDriver
			.instantQuery(
				`
						sum(
							namespace_cpu:kube_pod_container_resource_limits:sum{juju_model_uuid="${$currentKubernetes?.scopeUuid}"}
						)
						/
						sum(
							kube_node_status_allocatable{job="kube-state-metrics",juju_model_uuid="${$currentKubernetes?.scopeUuid}",resource="cpu"}
						)
						`,
			)
			.then((response) => {
				cpuLimits = response.result[0].value.value;
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
			<Card.Title>{m.cpu_usage()}</Card.Title>
			<Card.Action class="text-muted-foreground flex flex-col gap-0.5 text-sm">
				<div class="flex justify-between gap-2">
					<p>{m.requests()}</p>
					<p class="font-mono">{Math.round(cpuRequests * 100)}%</p>
				</div>
				<div class="flex justify-between gap-2">
					<p>{m.limits()}</p>
					<p class="font-mono">{Math.round(cpuLimits * 100)}%</p>
				</div>
			</Card.Action>
		</Card.Header>
		<Card.Content>
			<Chart.Container config={cpuUsagesConfiguration}>
				<AreaChart
					data={cpuUsages}
					x="time"
					xScale={scaleUtc()}
					yPadding={[0, 25]}
					series={[
						{
							key: 'value',
							label: cpuUsagesConfiguration.usage.label,
							color: cpuUsagesConfiguration.usage.color,
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
