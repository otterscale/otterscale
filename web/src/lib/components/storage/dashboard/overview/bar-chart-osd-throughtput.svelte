<script lang="ts">
	import type { Scope } from '$lib/api/scope/v1/scope_pb';
	import ComponentLoading from '$lib/components/custom/chart/component-loading.svelte';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Card from '$lib/components/ui/card';
	import * as Chart from '$lib/components/ui/chart/index.js';
	import { formatIO } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';
	import { getLocale } from '$lib/paraglide/runtime';
	import { BarChart, Highlight, type ChartContextValue } from 'layerchart';
	import { PrometheusDriver, type SampleValue } from 'prometheus-query';
	import { onMount } from 'svelte';
	import { cubicInOut } from 'svelte/easing';

	// Props
	let {
		client,
		scope,
		isReloading = $bindable()
	}: { client: PrometheusDriver; scope: Scope; isReloading: boolean } = $props();

	// Constants
	const CHART_TITLE = m.osd_throughPut();
	const CHART_DESCRIPTION = `${m.read()}/${m.write()}`;

	// Time range calculation
	const STEP_SECONDS = 60 * 60; // 1 hour
	const TIME_RANGE_HOURS = 24; // 24 hours of data
	const MILLISECONDS_PER_HOUR = 60 * 60 * 1000;
	const endTime = Date.now();
	const startTime = endTime - TIME_RANGE_HOURS * MILLISECONDS_PER_HOUR;

	// Chart configuration
	const chartConfig = {
		Read: {
			label: m.read(),
			color: 'var(--chart-1)'
		},
		Write: {
			label: m.write(),
			color: 'var(--chart-2)'
		}
	} satisfies Chart.ChartConfig;

	// Type
	type ChartKey = keyof typeof chartConfig;
	type TrafficData = {
		date: Date;
		Read: number;
		Write: number;
	};
	type MetricsResponse = {
		traffics: TrafficData[] | [];
		latestReadValue: number | undefined;
		latestWriteValue: number | undefined;
		latestReadUnit: string | undefined;
		latestWriteUnit: string | undefined;
	};

	// State
	let activeChart = $state<ChartKey>('Read');
	let context = $state<ChartContextValue>();

	// Derived state
	const queries = $derived({
		Read: `sum(irate(ceph_osd_op_r_out_bytes{juju_model_uuid=~"${scope.uuid}"}[1h]))`,
		Write: `sum(irate(ceph_osd_op_w_in_bytes{juju_model_uuid=~"${scope.uuid}"}[1h]))`
	});

	const activeSeries = $derived([
		{
			key: activeChart,
			label: chartConfig[activeChart].label,
			color: chartConfig[activeChart].color
		}
	]);

	// Helper functions
	function extractMetricValues(response: any): SampleValue[] {
		return response.result[0]?.values ?? [];
	}

	function extractLatestValue(response: any): number {
		return response.result[0]?.value?.value ?? 0;
	}

	function combineTrafficData(reads: SampleValue[], writes: SampleValue[]): TrafficData[] {
		return reads.map((sample: SampleValue, index: number) => ({
			date: sample.time,
			Read: sample.value,
			Write: writes[index]?.value ?? 0
		}));
	}

	// Auto Update
	let response = $state({} as MetricsResponse);
	let isLoading = $state(true);
	const reloadManager = new ReloadManager(fetch);

	// Data fetching function
	async function fetch(): Promise<void> {
		try {
			const [readResponse, writeResponse, latestReadResponse, latestWriteResponse] =
				await Promise.all([
					client.rangeQuery(queries.Read, startTime, endTime, STEP_SECONDS),
					client.rangeQuery(queries.Write, startTime, endTime, STEP_SECONDS),
					client.instantQuery(queries.Read),
					client.instantQuery(queries.Write)
				]);

			const reads = extractMetricValues(readResponse);
			const writes = extractMetricValues(writeResponse);
			const latestReadValue = extractLatestValue(latestReadResponse);
			const latestWriteValue = extractLatestValue(latestWriteResponse);

			const { value: readValue, unit: readUnit } = formatIO(latestReadValue);
			const { value: writeValue, unit: writeUnit } = formatIO(latestWriteValue);

			const traffics = combineTrafficData(reads, writes);

			response = {
				traffics,
				latestReadValue: readValue,
				latestWriteValue: writeValue,
				latestReadUnit: readUnit,
				latestWriteUnit: writeUnit
			};
		} catch (error) {
			console.error('Failed to fetch throughput metrics:', error);
			response = {
				traffics: [],
				latestReadValue: undefined,
				latestWriteValue: undefined,
				latestReadUnit: undefined,
				latestWriteUnit: undefined
			};
		}
	}

	$effect(() => {
		isReloading;
		if (isReloading) {
			reloadManager.restart();
		} else {
			reloadManager.stop();
		}
	});
	onMount(() => {
		fetch();
		isLoading = false;
	});
</script>

{#if isLoading}
	<ComponentLoading />
{:else}
	<Card.Root class="h-full gap-2">
		<Card.Header class="flex flex-col items-stretch space-y-0 border-b p-0 sm:flex-row">
			<div class="flex flex-1 flex-col justify-center gap-1 px-6 py-5 sm:py-6">
				<Card.Title>{CHART_TITLE}</Card.Title>
				<Card.Description>{CHART_DESCRIPTION}</Card.Description>
			</div>
			<div class="flex">
				{#each ['Read', 'Write'] as key (key)}
					{@const chart = key as ChartKey}
					{@const isActive = activeChart === chart}
					{@const latestValue =
						key === 'Read' ? response.latestReadValue : response.latestWriteValue}
					{@const latestUnit = key === 'Read' ? response.latestReadUnit : response.latestWriteUnit}
					<button
						data-active={isActive}
						class="data-[active=true]:bg-muted/50 relative z-30 flex flex-1 flex-col justify-center gap-1 border-t px-6 py-4 text-left even:border-l sm:border-l sm:border-t-0 sm:px-8 sm:py-6"
						onclick={() => (activeChart = chart)}
					>
						<span class="text-muted-foreground text-xs">
							{chartConfig[chart].label}
						</span>
						<span class="flex items-end gap-1 text-lg font-bold leading-none sm:text-3xl">
							{latestValue}
							<span class="text-muted-foreground text-xs">{latestUnit}</span>
						</span>
					</button>
				{/each}
			</div>
		</Card.Header>

		<Card.Content class="px-2 sm:p-6">
			<Chart.Container config={chartConfig} class="aspect-auto h-[150px] w-full">
				<BarChart
					bind:context
					data={response.traffics}
					x="date"
					axis="x"
					series={activeSeries}
					props={{
						bars: {
							stroke: 'none',
							rounded: 'none',
							initialY: context?.height,
							initialHeight: 0,
							motion: {
								y: { type: 'tween', duration: 500, easing: cubicInOut },
								height: { type: 'tween', duration: 500, easing: cubicInOut }
							}
						},
						highlight: { area: { fill: 'none' } },
						xAxis: {
							format: (d: Date) => {
								return d.toLocaleDateString(getLocale(), {
									month: 'numeric',
									day: 'numeric'
								});
							},
							ticks: 24
						}
					}}
				>
					{#snippet belowMarks()}
						<Highlight area={{ class: 'fill-muted' }} />
					{/snippet}
					{#snippet tooltip()}
						<Chart.Tooltip
							nameKey="views"
							labelFormatter={(time: Date) => {
								return time.toLocaleDateString('en-US', {
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
									class="border-(--color-border) bg-(--color-bg) aspect-square h-full w-fit shrink-0"
								></div>
								<div class="flex flex-1 shrink-0 items-center justify-between text-xs leading-none">
									<div class="grid gap-1.5">
										<span class="text-muted-foreground">{name}</span>
									</div>
									<p class="font-mono">{io} {unit}</p>
								</div>
							{/snippet}
						</Chart.Tooltip>
					{/snippet}
				</BarChart>
			</Chart.Container>
		</Card.Content>
	</Card.Root>
{/if}
