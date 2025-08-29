<script lang="ts">
	import type { Scope } from '$lib/api/scope/v1/scope_pb';
	import * as Card from '$lib/components/ui/card';
	import * as Chart from '$lib/components/ui/chart/index.js';
	import { getLocale } from '$lib/paraglide/runtime';
	import { scaleUtc } from 'd3-scale';
	import { curveLinear } from 'd3-shape';
	import { BarChart, type ChartContextValue, Highlight, LineChart } from 'layerchart';
	import { PrometheusDriver } from 'prometheus-query';
	import { cubicInOut } from 'svelte/easing';
	import { default as AreaCapacity } from './area-chart-capacity.svelte';
	import { default as PieOSDType } from './pie-chart-osd-type.svelte';
	import { default as TextClusterHealth } from './text-chart-cluster-health.svelte';
	import { default as TextOSDs } from './text-chart-osds.svelte';
	import { default as TextQuorum } from './text-chart-quorum.svelte';
	import { default as TextTimeTillFull } from './text-chart-time-till-full.svelte';
	import { default as UsageCapacity } from './usage-chart-capacity.svelte';
	import { default as AreaOSDReadLatency } from './area-chart-osd-read-latency.svelte';
	import { default as AreaOSDWriteLatency } from './area-chart-osd-write-latency.svelte';
	import { default as BarOSDThroughtput } from './bar-chart-osd-throughtput.svelte';
	import { default as BarOSDIOPS } from './bar-chart-osd-iops.svelte';

	import {
		chartConfig3,
		chartConfig4,
		chartConfig5,
		chartData2,
		chartData3,
		chartData4,
		chartData5
	} from './test-data';

	let { client, scope }: { client: PrometheusDriver; scope: Scope } = $props();

	// OSD Type Chart Configuration
	const osdTypeChartConfig = {
		hdd: { label: 'HDD', color: 'var(--chart-1)' },
		ssd: { label: 'SSD', color: 'var(--chart-2)' },
		nvme: { label: 'NVMe', color: 'var(--chart-3)' }
	} satisfies Chart.ChartConfig;

	// Fetch OSD Type Data
	async function fetchOSDTypeData() {
		try {
			const query = `count by (device_class) (ceph_osd_metadata{juju_model_uuid=~"${scope.uuid}"})`;
			const response = await client.instantQuery(query);

			const chartData = response.result.map((series, index) => {
				const deviceClass = series.metric?.device_class || 'unknown';
				const count = Number(series.value?.value || 0);

				// Map device class to color
				let color = 'var(--chart-4)'; // default color
				if (deviceClass in osdTypeChartConfig) {
					color = osdTypeChartConfig[deviceClass as keyof typeof osdTypeChartConfig].color;
				}

				return {
					deviceClass,
					count,
					color,
					fill: color
				};
			});

			const total = chartData.reduce((sum, item) => sum + item.count, 0);

			return {
				chartData,
				total
			};
		} catch (error) {
			console.error('Failed to fetch OSD type data:', error);
			// Return empty data on error
			return {
				chartData: [],
				total: 0
			};
		}
	}

	let context3 = $state<ChartContextValue>();

	let activeChart3 = $state<keyof typeof chartConfig3>('desktop');
	const total3 = $derived({
		desktop: chartData3.reduce((acc, curr) => acc + curr.desktop, 0),
		mobile: chartData3.reduce((acc, curr) => acc + curr.mobile, 0)
	});
	const activeSeries3 = $derived([
		{
			key: activeChart3,
			label: chartConfig3[activeChart3].label,
			color: chartConfig3[activeChart3].color
		}
	]);

	let context4 = $state<ChartContextValue>();
	let activeChart4 = $state<keyof typeof chartConfig4>('desktop');
	const total4 = $derived({
		desktop: chartData4.reduce((acc, curr) => acc + curr.desktop, 0),
		mobile: chartData4.reduce((acc, curr) => acc + curr.mobile, 0)
	});
	const activeSeries4 = $derived([
		{
			key: activeChart4,
			label: chartConfig4[activeChart4].label,
			color: chartConfig4[activeChart4].color
		}
	]);
</script>

<TextClusterHealth {client} {scope} />
<TextTimeTillFull {client} {scope} />
<UsageCapacity {client} {scope} />
<AreaCapacity {client} {scope} />
<PieOSDType {client} {scope} />
<TextQuorum {client} {scope} />
<TextOSDs {client} {scope} />
<BarOSDThroughtput {client} {scope} />
<BarOSDIOPS {client} {scope} />
<AreaOSDReadLatency {client} {scope} />
<AreaOSDWriteLatency {client} {scope} />
