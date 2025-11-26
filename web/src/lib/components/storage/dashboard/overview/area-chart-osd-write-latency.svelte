<script lang="ts">
	import Icon from '@iconify/svelte';
	import { scaleUtc } from 'd3-scale';
	import { curveLinear } from 'd3-shape';
	import { LineChart } from 'layerchart';
	import { PrometheusDriver } from 'prometheus-query';
	import { onDestroy, onMount } from 'svelte';
	import { SvelteDate } from 'svelte/reactivity';

	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Card from '$lib/components/ui/card';
	import * as Chart from '$lib/components/ui/chart/index.js';
	import * as Select from '$lib/components/ui/select/index.js';
	import { m } from '$lib/paraglide/messages';

	// Props
	let {
		client,
		scope,
		isReloading = $bindable()
	}: { client: PrometheusDriver; scope: string; isReloading: boolean } = $props();

	// Types
	type TimeInterval = 'day' | 'week' | 'month';
	type TimeRangeConfig = {
		count: number;
		label: string;
		stepSize: string;
	};

	// Constants
	const CHART_TITLE = m.osd();
	const CHART_DESCRIPTION = m.write_latencies();
	const CHART_CONFIG = {
		latency: {
			label: 'Write Latency (ms)',
			color: 'var(--chart-1)'
		}
	} satisfies Chart.ChartConfig;
	const TIME_INTERVALS: Record<TimeInterval, TimeRangeConfig> = {
		day: { count: 7, label: m.last_7_days(), stepSize: '1d' },
		week: { count: 5, label: m.last_5_weeks(), stepSize: '1w' },
		month: { count: 6, label: m.last_6_months(), stepSize: '1M' }
	};

	// Query
	const PROMETHEUS_QUERY = (scope: string) =>
		`quantile(0.95, (rate(ceph_osd_op_w_latency_sum{juju_model="${scope}"}[5m]) / ` +
		`on(ceph_daemon) rate(ceph_osd_op_w_latency_count{juju_model="${scope}"}[5m]) * 1000))`;

	// State
	let selectedInterval = $state<TimeInterval>('day');

	// Computed
	const timeRange = $derived(TIME_INTERVALS[selectedInterval]);

	// Helper functions
	function calculateTimeRange(
		interval: TimeInterval,
		index: number,
		today: Date
	): { start: Date; end: Date } {
		const start = new SvelteDate(today);
		const end = new SvelteDate(today);

		switch (interval) {
			case 'day': {
				start.setUTCDate(today.getUTCDate() - index);
				start.setUTCHours(0, 0, 0, 0);

				if (index === 0) {
					end.setTime(today.getTime());
				} else {
					end.setTime(start.getTime());
					end.setUTCHours(23, 59, 59, 999);
				}
				break;
			}

			case 'week': {
				const weeksBack = index * 7;
				start.setUTCDate(today.getUTCDate() - weeksBack);

				// Align to start of week (Sunday)
				const dayOfWeek = start.getUTCDay();
				start.setUTCDate(start.getUTCDate() - dayOfWeek);
				start.setUTCHours(0, 0, 0, 0);

				if (index === 0) {
					end.setTime(today.getTime());
				} else {
					end.setTime(start.getTime());
					end.setUTCDate(end.getUTCDate() + 6);
					end.setUTCHours(23, 59, 59, 999);
				}
				break;
			}

			case 'month': {
				const targetMonth = today.getUTCMonth() - index;
				start.setUTCMonth(targetMonth, 1);
				start.setUTCHours(0, 0, 0, 0);

				if (index === 0) {
					end.setTime(today.getTime());
				} else {
					end.setUTCMonth(targetMonth + 1, 0);
					end.setUTCHours(23, 59, 59, 999);
				}
				break;
			}
		}

		return { start, end };
	}

	function getXAxisFormat(interval: TimeInterval) {
		const formatters: Record<TimeInterval, (v: Date) => string> = {
			day: (v: Date) => v.toLocaleDateString('en-US', { day: 'numeric' }),
			week: (v: Date) => {
				const month = v.toLocaleString('en-US', { month: 'short' });
				const weekNum = Math.ceil(v.getUTCDate() / 7);
				return `${month}-W${weekNum}`;
			},
			month: (v: Date) => {
				const month = v.toLocaleString('en-US', { month: 'short' });
				const weekNum = Math.ceil(v.getUTCDate() / 7);
				return `${month}-W${weekNum}`;
			}
		};

		return formatters[interval];
	}

	async function fetchLatencyForPeriod(
		start: Date,
		end: Date
	): Promise<{ date: Date; latency: number }> {
		try {
			const query = PROMETHEUS_QUERY(scope);
			const response = await client.rangeQuery(
				query,
				start.getTime(),
				end.getTime(),
				timeRange.stepSize
			);

			const values = response.result?.[0]?.values;
			if (values?.length > 0) {
				const avgLatency =
					values.reduce(
						(sum: number, v: { value: number | string }) => sum + Number(v.value || 0),
						0
					) / values.length;

				return { date: start, latency: avgLatency };
			}
		} catch (error) {
			console.error('Error fetching latency data:', error);
		}

		return { date: start, latency: 0 };
	}

	// Data fetching
	let response = $state([] as { date: Date; latency: number }[]);
	let isLoading = $state(true);
	const reloadManager = new ReloadManager(fetch);

	async function fetch(): Promise<void> {
		const today = new SvelteDate();
		const promises = Array.from({ length: timeRange.count }, (_, i) => {
			const index = timeRange.count - 1 - i;
			const { start, end } = calculateTimeRange(selectedInterval, index, today);
			return fetchLatencyForPeriod(start, end);
		});

		response = await Promise.all(promises);
	}

	$effect(() => {
		if (isReloading) {
			reloadManager.restart();
		} else {
			reloadManager.stop();
		}
	});

	$effect(() => {
		void selectedInterval;
		fetch();
	});

	onMount(() => {
		fetch();
		isLoading = false;
	});
	onDestroy(() => {
		reloadManager.stop();
	});
</script>

<Card.Root class="h-[162px]">
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
	{#if isLoading}
		<div class="flex h-9 w-full items-center justify-center">
			<Icon icon="svg-spinners:6-dots-rotate" class="size-10" />
		</div>
	{:else if response.length === 0}
		<div class="flex h-full w-full flex-col items-center justify-center">
			<Icon icon="ph:chart-bar-fill" class="size-6 animate-pulse text-muted-foreground" />
			<p class="p-0 text-xs text-muted-foreground">{m.no_data_display()}</p>
		</div>
	{:else}
		<Card.Content>
			<Chart.Container config={CHART_CONFIG} class="h-16 w-full p-1">
				<LineChart
					data={response}
					x="date"
					xScale={scaleUtc()}
					axis="x"
					series={[
						{
							key: 'latency',
							label: 'Write Latency',
							color: CHART_CONFIG.latency.color
						}
					]}
					props={{
						spline: { curve: curveLinear, motion: 'tween', strokeWidth: 2 },
						xAxis: {
							format: getXAxisFormat(selectedInterval),
							ticks: response.length
						},
						yAxis: { format: () => '' },
						highlight: { points: { r: 4 } }
					}}
				>
					{#snippet tooltip()}
						<Chart.Tooltip hideLabel />
					{/snippet}
				</LineChart>
			</Chart.Container>
		</Card.Content>
	{/if}
</Card.Root>
