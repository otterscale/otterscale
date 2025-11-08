<script lang="ts">
	import { PieChart, Text } from 'layerchart';
	import { PrometheusDriver } from 'prometheus-query';
	import { onDestroy, onMount } from 'svelte';

	import type { Scope } from '$lib/api/scope/v1/scope_pb';
	import ComponentLoading from '$lib/components/custom/chart/component-loading.svelte';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Card from '$lib/components/ui/card';
	import * as Chart from '$lib/components/ui/chart/index.js';
	import { m } from '$lib/paraglide/messages';

	// Props
	let {
		client,
		scope,
		isReloading = $bindable()
	}: { client: PrometheusDriver; scope: Scope; isReloading: boolean } = $props();

	// Constants
	const CHART_TITLE = m.osd_type();
	const CHART_DESCRIPTION = m.osd_type_distribution();
	const UNKNOWN_DEVICE_CLASS = 'unknown';

	// Predefined device class configurations
	const DEVICE_CLASS_CONFIGS: Record<string, { label: string; color: string }> = {
		hdd: { label: 'HDD', color: 'var(--chart-1)' },
		ssd: { label: 'SSD', color: 'var(--chart-2)' },
		nvme: { label: 'NVMe', color: 'var(--chart-3)' }
	} as const;

	// Chart colors for unknown device classes
	const FALLBACK_CHART_COLORS = [
		'var(--chart-4)',
		'var(--chart-5)',
		'var(--chart-1)',
		'var(--chart-2)',
		'var(--chart-3)'
	] as const;

	// Utility functions
	function getDeviceClassConfig(deviceClass: string, index: number) {
		if (deviceClass in DEVICE_CLASS_CONFIGS) {
			return DEVICE_CLASS_CONFIGS[deviceClass];
		}

		return {
			label: deviceClass.toUpperCase(),
			color: FALLBACK_CHART_COLORS[index % FALLBACK_CHART_COLORS.length]
		};
	}

	function generateChartConfig(deviceClasses: string[]): Chart.ChartConfig {
		const config: Record<string, { label: string; color: string }> = {};

		deviceClasses.forEach((deviceClass, index) => {
			config[deviceClass] = getDeviceClassConfig(deviceClass, index);
		});

		return config as Chart.ChartConfig;
	}

	// Queries
	const queries = $derived({
		osdTypeCount: `count by (device_class) (ceph_osd_metadata{juju_model_uuid=~"${scope.uuid}"})`
	});

	// Auto Update
	type ChartDataItem = {
		deviceClass: string;
		count: number;
		color: string;
		fill: string;
		label: string;
	};

	let response = $state({
		chartData: [] as ChartDataItem[],
		total: 0,
		chartConfig: {} as Chart.ChartConfig
	});
	let isLoading = $state(true);
	const reloadManager = new ReloadManager(fetch);
	async function fetch() {
		try {
			const prometheusResponse = await client.instantQuery(queries.osdTypeCount);

			if (!prometheusResponse?.result?.length) {
				response = { chartData: [], total: 0, chartConfig: {} as Chart.ChartConfig };
				return;
			}

			const chartData = prometheusResponse.result.map((series, index) => {
				const deviceClass = series.metric?.labels.device_class || UNKNOWN_DEVICE_CLASS;
				const count = Number(series.value?.value || 0);
				const config = getDeviceClassConfig(deviceClass, index);

				return {
					deviceClass,
					count,
					color: config.color,
					fill: config.color,
					label: config.label
				};
			});

			const total = chartData.reduce((sum, item) => sum + item.count, 0);
			const deviceClasses = chartData.map((item) => item.deviceClass);
			const chartConfig = generateChartConfig(deviceClasses);

			response = { chartData, total, chartConfig };
		} catch (error) {
			console.error('Failed to fetch OSD type data:', error);
			response = {
				chartData: [],
				total: 0,
				chartConfig: {} as Chart.ChartConfig
			};
		}
	}

	$effect(() => {
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
	onDestroy(() => {
		reloadManager.stop();
	});
</script>

{#if isLoading}
	<ComponentLoading />
{:else}
	<Card.Root class="h-full gap-2">
		<Card.Header class="items-center">
			<Card.Title>{CHART_TITLE}</Card.Title>
			<Card.Description>{CHART_DESCRIPTION}</Card.Description>
		</Card.Header>
		<Card.Content class="flex-1">
			{#if response.chartData.length > 0}
				<Chart.Container config={response.chartConfig} class="mx-auto aspect-square max-h-[200px]">
					<PieChart
						data={response.chartData}
						key="deviceClass"
						value="count"
						c="color"
						innerRadius={60}
						padding={28}
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
				<div class="flex h-full items-center justify-center text-muted-foreground">
					<p>No OSD data available</p>
				</div>
			{/if}
		</Card.Content>
	</Card.Root>
{/if}
