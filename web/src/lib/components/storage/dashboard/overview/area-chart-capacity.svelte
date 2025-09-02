<script lang="ts">
	import type { Scope } from '$lib/api/scope/v1/scope_pb';
	import ComponentLoading from '$lib/components/custom/chart/component-loading.svelte';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Card from '$lib/components/ui/card';
	import * as Chart from '$lib/components/ui/chart/index.js';
	import * as Select from '$lib/components/ui/select/index.js';
	import { formatCapacity } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';
	import { scaleUtc } from 'd3-scale';
	import { curveStep } from 'd3-shape';
	import { AreaChart } from 'layerchart';
	import { PrometheusDriver } from 'prometheus-query';
	import { onMount } from 'svelte';

	// Props
	let {
		client,
		scope,
		isReloading = $bindable()
	}: { client: PrometheusDriver; scope: Scope; isReloading: boolean } = $props();

	// Types
	type TimeInterval = 'day' | 'week' | 'month';

	// Constants
	const CHART_TITLE = m.capacity();
	const CHART_DESCRIPTION = m.capacity_usage_changes();
	const CHART_CONFIG = {
		used: { label: 'Used', color: 'var(--chart-1)' },
		total: { label: 'Total', color: 'var(--chart-3)' }
	} satisfies Chart.ChartConfig;

	const TIME_INTERVALS: Record<TimeInterval, { count: number; label: string }> = {
		day: { count: 7, label: m.last_7_days() },
		week: { count: 5, label: m.last_5_weeks() },
		month: { count: 6, label: m.last_6_months() }
	};

	// State
	let selectedInterval = $state<TimeInterval>('day');

	// Computed
	const timeRange = $derived(TIME_INTERVALS[selectedInterval]);

	// Helper functions
	function calculateIntervals(
		selectedInterval: TimeInterval,
		index: number,
		today: Date
	): { start: Date; end: Date } {
		let intervalStart: Date;
		let intervalEnd: Date;

		switch (selectedInterval) {
			case 'day': {
				intervalStart = new Date(today);
				intervalStart.setUTCDate(today.getUTCDate() - index);
				intervalStart.setUTCHours(0, 0, 0, 0);

				intervalEnd = new Date(intervalStart);
				if (index === 0) {
					intervalEnd = new Date(today);
				} else {
					intervalEnd.setUTCHours(23, 59, 59, 999);
				}
				break;
			}
			case 'week': {
				intervalStart = new Date(today);
				intervalStart.setUTCDate(today.getUTCDate() - index * 7);
				const dayOfWeek = intervalStart.getUTCDay();
				intervalStart.setUTCDate(intervalStart.getUTCDate() - dayOfWeek);
				intervalStart.setUTCHours(0, 0, 0, 0);

				intervalEnd = new Date(intervalStart);
				if (index === 0) {
					intervalEnd = new Date(today);
				} else {
					intervalEnd.setUTCDate(intervalEnd.getUTCDate() + 6);
					intervalEnd.setUTCHours(23, 59, 59, 999);
				}
				break;
			}
			case 'month':
			default: {
				const year = today.getUTCFullYear();
				const month = today.getUTCMonth() - index;
				intervalStart = new Date(Date.UTC(year, month, 1, 0, 0, 0, 0));

				if (index === 0) {
					intervalEnd = new Date(today);
				} else {
					intervalEnd = new Date(Date.UTC(year, month + 1, 0, 23, 59, 59, 999));
				}
				break;
			}
		}

		return { start: intervalStart, end: intervalEnd };
	}

	function getXAxisFormat(interval: TimeInterval) {
		const formatters: Record<TimeInterval, (v: Date) => string> = {
			day: (v: Date) => v.toLocaleDateString('en-US', { day: 'numeric' }),
			week: (v: Date) => {
				const month = v.toLocaleString('en-US', { month: 'short' });
				const weekNum = Math.ceil(v.getUTCDate() / 7);
				return `${month}-W${weekNum}`;
			},
			month: (v: Date) => v.toLocaleDateString('en-US', { month: 'short' })
		};

		return formatters[interval];
	}

	function getYAxisDomain(
		data: { date: Date; used: number; total: number; available: number }[]
	): [number, number] {
		const maxTotal = Math.max(...data.map((d) => d.total || 0));
		return [0, maxTotal];
	}

	async function fetchMetricForInterval(
		intervalStart: Date,
		intervalEnd: Date
	): Promise<{ date: Date; used: number; total: number; available: number }> {
		const endTimestamp = Math.floor(intervalEnd.getTime() / 1000);

		const queries = {
			used: `ceph_cluster_total_used_bytes{juju_model_uuid=~"${scope.uuid}"} @ ${endTimestamp}`,
			total: `ceph_cluster_total_bytes{juju_model_uuid=~"${scope.uuid}"} @ ${endTimestamp}`
		};

		try {
			const [usedResponse, totalResponse] = await Promise.all([
				client.instantQuery(queries.used),
				client.instantQuery(queries.total)
			]);

			const usedValue = Number(usedResponse.result[0]?.value?.value || 0);
			const totalValue = Number(totalResponse.result[0]?.value?.value || 0);

			return {
				date: intervalStart,
				used: usedValue,
				total: totalValue,
				available: totalValue - usedValue
			};
		} catch (error) {
			return {
				date: intervalStart,
				used: 0,
				total: 0,
				available: 0
			};
		}
	}
	// Auto Update
	let response = $state([] as { date: Date; used: number; total: number; available: number }[]);
	let isLoading = $state(true);
	const reloadManager = new ReloadManager(fetch);
	async function fetch(): Promise<void> {
		const today = new Date();
		const promises = [];

		for (let i = timeRange.count - 1; i >= 0; i--) {
			const { start, end } = calculateIntervals(selectedInterval, i, today);
			promises.push(fetchMetricForInterval(start, end));
		}

		response = await Promise.all(promises);
	}

	$effect(() => {
		isReloading;
		if (isReloading) {
			reloadManager.restart();
		} else {
			reloadManager.stop();
		}
	});

	$effect(() => {
		selectedInterval;
		fetch();
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
		<Card.Header class="flex items-center">
			<div class="grid flex-1 gap-1 text-center sm:text-left">
				<Card.Title>{CHART_TITLE}</Card.Title>
				<Card.Description>{CHART_DESCRIPTION}</Card.Description>
			</div>

			<Select.Root type="single" bind:value={selectedInterval}>
				<Select.Trigger class="w-fit rounded-lg sm:ml-auto" aria-label="Select time range">
					{timeRange.label}
				</Select.Trigger>
				<Select.Content class="rounded-xl">
					<Select.Item value="day" class="rounded-lg">{TIME_INTERVALS.day.label}</Select.Item>
					<Select.Item value="week" class="rounded-lg">{TIME_INTERVALS.week.label}</Select.Item>
					<Select.Item value="month" class="rounded-lg">{TIME_INTERVALS.month.label}</Select.Item>
				</Select.Content>
			</Select.Root>
		</Card.Header>

		<Card.Content>
			<Chart.Container config={CHART_CONFIG} class="h-[200px] w-full">
				<AreaChart
					data={response}
					x="date"
					xScale={scaleUtc()}
					yDomain={getYAxisDomain(response)}
					series={[
						{
							key: 'used',
							label: 'Used',
							color: CHART_CONFIG.used.color
						}
					]}
					seriesLayout="stack"
					props={{
						area: {
							curve: curveStep,
							'fill-opacity': 0.4,
							line: { class: 'stroke-1' },
							motion: 'tween'
						},
						xAxis: {
							format: getXAxisFormat(selectedInterval),
							ticks: response.length
						},
						yAxis: { format: () => '' }
					}}
				>
					{#snippet tooltip()}
						<Chart.Tooltip
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
								{@const { value: io, unit } = formatCapacity(Number(value))}
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
				</AreaChart>
			</Chart.Container>
		</Card.Content>
	</Card.Root>
{/if}
