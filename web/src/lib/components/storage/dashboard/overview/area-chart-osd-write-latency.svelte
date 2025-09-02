<script lang="ts">
	import type { Scope } from '$lib/api/scope/v1/scope_pb';
	import ComponentLoading from '$lib/components/custom/chart/component-loading.svelte';
	import ErrorLayout from '$lib/components/custom/chart/layout/standard-error.svelte';
	import * as Card from '$lib/components/ui/card';
	import * as Chart from '$lib/components/ui/chart/index.js';
	import * as Select from '$lib/components/ui/select/index.js';
	import { formatCapacity } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';
	import { getLocale } from '$lib/paraglide/runtime';
	import { scaleUtc } from 'd3-scale';
	import { curveBasis, curveLinear, curveStep } from 'd3-shape';
	import { LineChart } from 'layerchart';
	import { PrometheusDriver } from 'prometheus-query';

	// Props
	let { client, scope }: { client: PrometheusDriver; scope: Scope } = $props();

	// Types
	type TimeInterval = 'day' | 'week' | 'month';

	interface MetricData {
		date: Date;
		latency: number;
	}

	interface TimeRangeConfig {
		count: number;
		label: string;
		stepSize: string;
	}

	// Constants
	const CHART_TITLE = m.osd();
	const CHART_DESCRIPTION = m.write_latencies();
	const CHART_CONFIG = {
		latency: {
			label: 'Write Latency (ms)',
			color: 'var(--chart-1)',
		},
	} satisfies Chart.ChartConfig;

	const TIME_INTERVALS: Record<TimeInterval, TimeRangeConfig> = {
		day: { count: 7, label: m.last_7_days(), stepSize: '1d' },
		week: { count: 5, label: m.last_5_weeks(), stepSize: '1w' },
		month: { count: 6, label: m.last_6_months(), stepSize: '1M' },
	};

	const PROMETHEUS_QUERY = (uuid: string) =>
		`quantile(0.95, (rate(ceph_osd_op_w_latency_sum{juju_model_uuid=~"${uuid}"}[5m]) / ` +
		`on(ceph_daemon) rate(ceph_osd_op_w_latency_count{juju_model_uuid=~"${uuid}"}[5m]) * 1000))`;

	// State
	let selectedInterval = $state<TimeInterval>('day');

	// Computed
	const timeRange = $derived(TIME_INTERVALS[selectedInterval]);

	// Helper functions
	function calculateTimeRange(interval: TimeInterval, index: number): { start: Date; end: Date } {
		const today = new Date();
		const start = new Date(today);
		const end = new Date(today);

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

	async function fetchLatencyForPeriod(start: Date, end: Date): Promise<{ date: Date; latency: number }> {
		try {
			const query = PROMETHEUS_QUERY(scope.uuid);
			const response = await client.rangeQuery(query, start, end, timeRange.stepSize);

			if (response.result?.[0]?.values?.length > 0) {
				// Get the average of all values in the time series
				const values = response.result[0].values;
				const avgLatency =
					values.reduce((sum: number, v: { value: number | string }) => sum + Number(v.value || 0), 0) /
					values.length;

				return { date: start, latency: avgLatency };
			}
		} catch (error) {
			console.error('Error fetching latency data:', error);
		}

		return { date: start, latency: 0 };
	}

	async function fetchMetrics(): Promise<{ date: Date; latency: number }[]> {
		const promises = Array.from({ length: timeRange.count }, (_, i) => {
			const index = timeRange.count - 1 - i;
			const { start, end } = calculateTimeRange(selectedInterval, index);
			return fetchLatencyForPeriod(start, end);
		});

		return Promise.all(promises);
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
			},
		};

		return formatters[interval];
	}
</script>

{#key selectedInterval}
	{#await fetchMetrics()}
		<ComponentLoading />
	{:then response}
		<Card.Root class="gap-2">
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
				<Chart.Container config={CHART_CONFIG} class="h-[64px] w-full">
					<LineChart
						data={response}
						x="date"
						xScale={scaleUtc()}
						axis="x"
						series={[
							{
								key: 'latency',
								label: 'Write Latency',
								color: CHART_CONFIG.latency.color,
							},
						]}
						props={{
							spline: { curve: curveLinear, motion: 'tween', strokeWidth: 2 },
							xAxis: {
								format: getXAxisFormat(selectedInterval),
								ticks: response.length,
							},
							yAxis: { format: () => '' },
							highlight: { points: { r: 4 } },
						}}
					>
						{#snippet tooltip()}
							<Chart.Tooltip hideLabel />
						{/snippet}
					</LineChart>
				</Chart.Container>
			</Card.Content>
		</Card.Root>
	{:catch error}
		<Card.Root class="gap-2">
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
				<Chart.Container config={CHART_CONFIG} class="h-[64px] w-full" />
			</Card.Content>
		</Card.Root>
	{/await}
{/key}
