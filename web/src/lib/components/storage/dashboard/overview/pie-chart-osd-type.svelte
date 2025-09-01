<script lang="ts">
	import type { Scope } from '$lib/api/scope/v1/scope_pb';
	import ComponentLoading from '$lib/components/custom/chart/component-loading.svelte';
	import * as Card from '$lib/components/ui/card';
	import * as Chart from '$lib/components/ui/chart/index.js';
	import { PieChart, Text } from 'layerchart';
	import { PrometheusDriver } from 'prometheus-query';
	import { m } from '$lib/paraglide/messages';

	// Props
	let { client, scope }: { client: PrometheusDriver; scope: Scope } = $props();

	// Constants
	const CHART_TITLE = m.osd_type();
	const CHART_DESCRIPTION = m.osd_type_distribution();
	const UNKNOWN_DEVICE_CLASS = 'unknown';

	// Predefined device class configurations
	const DEVICE_CLASS_CONFIGS: Record<string, { label: string; color: string }> = {
		hdd: { label: 'HDD', color: 'var(--chart-1)' },
		ssd: { label: 'SSD', color: 'var(--chart-2)' },
		nvme: { label: 'NVMe', color: 'var(--chart-3)' },
	} as const;

	// Chart colors for unknown device classes
	const FALLBACK_CHART_COLORS = [
		'var(--chart-4)',
		'var(--chart-5)',
		'var(--chart-1)',
		'var(--chart-2)',
		'var(--chart-3)',
	] as const;

	// Utility functions
	function getDeviceClassConfig(deviceClass: string, index: number) {
		if (deviceClass in DEVICE_CLASS_CONFIGS) {
			return DEVICE_CLASS_CONFIGS[deviceClass];
		}

		return {
			label: deviceClass.toUpperCase(),
			color: FALLBACK_CHART_COLORS[index % FALLBACK_CHART_COLORS.length],
		};
	}

	function generateChartConfig(deviceClasses: string[]): Chart.ChartConfig {
		const config: Record<string, { label: string; color: string }> = {};

		deviceClasses.forEach((deviceClass, index) => {
			config[deviceClass] = getDeviceClassConfig(deviceClass, index);
		});

		return config as Chart.ChartConfig;
	}

	function transformResponseData(response: any, chartConfig: Chart.ChartConfig) {
		const chartData = response.result.map((series: any, index: number) => {
			const deviceClass = series.metric?.labels.device_class || UNKNOWN_DEVICE_CLASS;
			const count = Number(series.value?.value || 0);
			const config = getDeviceClassConfig(deviceClass, index);

			return {
				deviceClass,
				count,
				color: config.color,
				fill: config.color,
				label: config.label,
			};
		});

		const total = chartData.reduce((sum, item) => sum + item.count, 0);

		return { chartData, total };
	}

	// Queries
	const queries = $derived({
		osdTypeCount: `count by (device_class) (ceph_osd_metadata{juju_model_uuid=~"${scope.uuid}"})`,
	});

	// Data fetching function
	async function fetchMetrics() {
		try {
			const response = await client.instantQuery(queries.osdTypeCount);

			if (!response?.result?.length) {
				return {
					chartData: [],
					total: 0,
					chartConfig: {} as Chart.ChartConfig,
				};
			}

			// Extract device classes and generate chart configuration
			const deviceClasses = response.result.map(
				(series: any) => series.metric?.labels.device_class || UNKNOWN_DEVICE_CLASS,
			);
			const chartConfig = generateChartConfig(deviceClasses);

			// Transform response data
			const { chartData, total } = transformResponseData(response, chartConfig);

			return { chartData, total, chartConfig };
		} catch (error) {
			console.error('Failed to fetch OSD type data:', error);
			return {
				chartData: [],
				total: 0,
				chartConfig: {} as Chart.ChartConfig,
			};
		}
	}
	// Component constants
	const CHART_INNER_RADIUS = 60;
	const CHART_PADDING = 28;
	const CHART_CLASSES = 'mx-auto aspect-square max-h-[200px]';
	const CARD_CLASSES = 'gap-2';
</script>

{#await fetchMetrics()}
	<ComponentLoading />
{:then response}
	<Card.Root class={CARD_CLASSES}>
		<Card.Header class="items-center">
			<Card.Title>{CHART_TITLE}</Card.Title>
			<Card.Description>{CHART_DESCRIPTION}</Card.Description>
		</Card.Header>
		<Card.Content class="flex-1">
			{#if response.chartData.length > 0}
				<Chart.Container config={response.chartConfig} class={CHART_CLASSES}>
					<PieChart
						data={response.chartData}
						key="deviceClass"
						value="count"
						c="color"
						innerRadius={CHART_INNER_RADIUS}
						padding={CHART_PADDING}
						props={{ pie: { motion: 'tween' } }}
					>
						{#snippet aboveMarks()}
							<Text
								value={String(response.total)}
								textAnchor="middle"
								verticalAnchor="middle"
								class="fill-foreground text-3xl! font-bold"
								dy={3}
							/>
							<Text
								value="Total OSDs"
								textAnchor="middle"
								verticalAnchor="middle"
								class="fill-muted-foreground! text-muted-foreground"
								dy={22}
							/>
						{/snippet}
						{#snippet tooltip()}
							<Chart.Tooltip hideLabel />
						{/snippet}
					</PieChart>
				</Chart.Container>
			{:else}
				<div class="text-muted-foreground flex h-full items-center justify-center">
					<p>No OSD data available</p>
				</div>
			{/if}
		</Card.Content>
	</Card.Root>
{:catch error}
	<Card.Root class={CARD_CLASSES}>
		<Card.Header class="items-center">
			<Card.Title>{CHART_TITLE}</Card.Title>
			<Card.Description>{CHART_DESCRIPTION}</Card.Description>
		</Card.Header>
		<Card.Content class="flex-1">
			<div class="text-destructive flex h-full items-center justify-center">
				<p>Failed to load chart data</p>
			</div>
		</Card.Content>
	</Card.Root>
{/await}
