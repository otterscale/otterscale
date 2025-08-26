<script lang="ts">
	import { page } from '$app/state';
	import { env } from '$env/dynamic/public';
	import { EnvironmentService } from '$lib/api/environment/v1/environment_pb';
	import Loading from '$lib/components/custom/loading/report.svelte';
	import { Dashboard } from '$lib/components/storage/dashboard';
	import * as Card from '$lib/components/ui/card';
	import * as Chart from '$lib/components/ui/chart/index.js';
	import * as Tabs from '$lib/components/ui/tabs';
	import { getLocale } from '$lib/paraglide/runtime';
	import { dynamicPaths } from '$lib/path';
	import { activeScope, breadcrumb } from '$lib/stores';
	import { createClient, type Transport } from '@connectrpc/connect';
	import ActivityIcon from '@lucide/svelte/icons/activity';
	import { scaleUtc } from 'd3-scale';
	import { curveLinear, curveStep } from 'd3-shape';
	import { ArcChart, AreaChart, BarChart, type ChartContextValue, Highlight, LineChart, PieChart, Text } from 'layerchart';
	import { PrometheusDriver } from 'prometheus-query';
	import { getContext, onMount } from 'svelte';
	import { cubicInOut } from 'svelte/easing';

	// Set breadcrumb navigation
	breadcrumb.set({ parents: [], current: dynamicPaths.storage(page.params.scope) });

	//
	const chartData1 = [{ browser: 'safari', visitors: 1260, color: 'var(--color-safari)' }];
	const chartConfig1 = {
		visitors: { label: 'Visitors' },
		safari: { label: 'Safari', color: 'var(--chart-2)' }
	} satisfies Chart.ChartConfig;

	//
	const chartData2 = [
		{ date: new Date('2024-01-01'), desktop: 10 },
		{ date: new Date('2024-02-01'), desktop: 12 },
		{ date: new Date('2024-03-01'), desktop: 24 },
		{ date: new Date('2024-04-01'), desktop: 28 },
		{ date: new Date('2024-05-01'), desktop: 36 },
		{ date: new Date('2024-06-01'), desktop: 48 }
	];
	const chartConfig2 = {
		desktop: { label: 'Desktop', color: 'var(--chart-1)', icon: ActivityIcon }
	} satisfies Chart.ChartConfig;

	//
	const chartData3 = [
		{ date: new Date('2024-04-01'), desktop: 222, mobile: 150 },
		{ date: new Date('2024-04-02'), desktop: 97, mobile: 180 },
		{ date: new Date('2024-04-03'), desktop: 167, mobile: 120 },
		{ date: new Date('2024-04-04'), desktop: 242, mobile: 260 },
		{ date: new Date('2024-04-05'), desktop: 373, mobile: 290 },
		{ date: new Date('2024-04-06'), desktop: 301, mobile: 340 },
		{ date: new Date('2024-04-07'), desktop: 245, mobile: 180 },
		{ date: new Date('2024-04-08'), desktop: 409, mobile: 320 },
		{ date: new Date('2024-04-09'), desktop: 59, mobile: 110 },
		{ date: new Date('2024-04-10'), desktop: 261, mobile: 190 },
		{ date: new Date('2024-04-11'), desktop: 327, mobile: 350 },
		{ date: new Date('2024-04-12'), desktop: 292, mobile: 210 },
		{ date: new Date('2024-04-13'), desktop: 342, mobile: 380 },
		{ date: new Date('2024-04-14'), desktop: 137, mobile: 220 },
		{ date: new Date('2024-04-15'), desktop: 120, mobile: 170 },
		{ date: new Date('2024-04-16'), desktop: 138, mobile: 190 },
		{ date: new Date('2024-04-17'), desktop: 446, mobile: 360 },
		{ date: new Date('2024-04-18'), desktop: 364, mobile: 410 },
		{ date: new Date('2024-04-19'), desktop: 243, mobile: 180 },
		{ date: new Date('2024-04-20'), desktop: 89, mobile: 150 },
		{ date: new Date('2024-04-21'), desktop: 137, mobile: 200 },
		{ date: new Date('2024-04-22'), desktop: 224, mobile: 170 },
		{ date: new Date('2024-04-23'), desktop: 138, mobile: 230 },
		{ date: new Date('2024-04-24'), desktop: 387, mobile: 290 },
		{ date: new Date('2024-04-25'), desktop: 215, mobile: 250 },
		{ date: new Date('2024-04-26'), desktop: 75, mobile: 130 },
		{ date: new Date('2024-04-27'), desktop: 383, mobile: 420 },
		{ date: new Date('2024-04-28'), desktop: 122, mobile: 180 },
		{ date: new Date('2024-04-29'), desktop: 315, mobile: 240 },
		{ date: new Date('2024-04-30'), desktop: 454, mobile: 380 },
		{ date: new Date('2024-05-01'), desktop: 165, mobile: 220 },
		{ date: new Date('2024-05-02'), desktop: 293, mobile: 310 },
		{ date: new Date('2024-05-03'), desktop: 247, mobile: 190 },
		{ date: new Date('2024-05-04'), desktop: 385, mobile: 420 },
		{ date: new Date('2024-05-05'), desktop: 481, mobile: 390 },
		{ date: new Date('2024-05-06'), desktop: 498, mobile: 520 },
		{ date: new Date('2024-05-07'), desktop: 388, mobile: 300 },
		{ date: new Date('2024-05-08'), desktop: 149, mobile: 210 },
		{ date: new Date('2024-05-09'), desktop: 227, mobile: 180 },
		{ date: new Date('2024-05-10'), desktop: 293, mobile: 330 },
		{ date: new Date('2024-05-11'), desktop: 335, mobile: 270 },
		{ date: new Date('2024-05-12'), desktop: 197, mobile: 240 },
		{ date: new Date('2024-05-13'), desktop: 197, mobile: 160 },
		{ date: new Date('2024-05-14'), desktop: 448, mobile: 490 },
		{ date: new Date('2024-05-15'), desktop: 473, mobile: 380 },
		{ date: new Date('2024-05-16'), desktop: 338, mobile: 400 },
		{ date: new Date('2024-05-17'), desktop: 499, mobile: 420 },
		{ date: new Date('2024-05-18'), desktop: 315, mobile: 350 },
		{ date: new Date('2024-05-19'), desktop: 235, mobile: 180 },
		{ date: new Date('2024-05-20'), desktop: 177, mobile: 230 },
		{ date: new Date('2024-05-21'), desktop: 82, mobile: 140 },
		{ date: new Date('2024-05-22'), desktop: 81, mobile: 120 },
		{ date: new Date('2024-05-23'), desktop: 252, mobile: 290 },
		{ date: new Date('2024-05-24'), desktop: 294, mobile: 220 },
		{ date: new Date('2024-05-25'), desktop: 201, mobile: 250 },
		{ date: new Date('2024-05-26'), desktop: 213, mobile: 170 },
		{ date: new Date('2024-05-27'), desktop: 420, mobile: 460 },
		{ date: new Date('2024-05-28'), desktop: 233, mobile: 190 },
		{ date: new Date('2024-05-29'), desktop: 78, mobile: 130 },
		{ date: new Date('2024-05-30'), desktop: 340, mobile: 280 },
		{ date: new Date('2024-05-31'), desktop: 178, mobile: 230 },
		{ date: new Date('2024-06-01'), desktop: 178, mobile: 200 },
		{ date: new Date('2024-06-02'), desktop: 470, mobile: 410 },
		{ date: new Date('2024-06-03'), desktop: 103, mobile: 160 },
		{ date: new Date('2024-06-04'), desktop: 439, mobile: 380 },
		{ date: new Date('2024-06-05'), desktop: 88, mobile: 140 },
		{ date: new Date('2024-06-06'), desktop: 294, mobile: 250 },
		{ date: new Date('2024-06-07'), desktop: 323, mobile: 370 },
		{ date: new Date('2024-06-08'), desktop: 385, mobile: 320 },
		{ date: new Date('2024-06-09'), desktop: 438, mobile: 480 },
		{ date: new Date('2024-06-10'), desktop: 155, mobile: 200 },
		{ date: new Date('2024-06-11'), desktop: 92, mobile: 150 },
		{ date: new Date('2024-06-12'), desktop: 492, mobile: 420 },
		{ date: new Date('2024-06-13'), desktop: 81, mobile: 130 },
		{ date: new Date('2024-06-14'), desktop: 426, mobile: 380 },
		{ date: new Date('2024-06-15'), desktop: 307, mobile: 350 },
		{ date: new Date('2024-06-16'), desktop: 371, mobile: 310 },
		{ date: new Date('2024-06-17'), desktop: 475, mobile: 520 },
		{ date: new Date('2024-06-18'), desktop: 107, mobile: 170 },
		{ date: new Date('2024-06-19'), desktop: 341, mobile: 290 },
		{ date: new Date('2024-06-20'), desktop: 408, mobile: 450 },
		{ date: new Date('2024-06-21'), desktop: 169, mobile: 210 },
		{ date: new Date('2024-06-22'), desktop: 317, mobile: 270 },
		{ date: new Date('2024-06-23'), desktop: 480, mobile: 530 },
		{ date: new Date('2024-06-24'), desktop: 132, mobile: 180 },
		{ date: new Date('2024-06-25'), desktop: 141, mobile: 190 },
		{ date: new Date('2024-06-26'), desktop: 434, mobile: 380 },
		{ date: new Date('2024-06-27'), desktop: 448, mobile: 490 },
		{ date: new Date('2024-06-28'), desktop: 149, mobile: 200 },
		{ date: new Date('2024-06-29'), desktop: 103, mobile: 160 },
		{ date: new Date('2024-06-30'), desktop: 446, mobile: 400 }
	];
	const chartConfig3 = {
		views: { label: 'Page Views', color: '' },
		desktop: { label: 'Desktop', color: 'var(--chart-1)' },
		mobile: { label: 'Mobile', color: 'var(--chart-3)' }
	} satisfies Chart.ChartConfig;
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

	//
	const chartData4 = [
		{ date: new Date('2024-04-01'), desktop: 222, mobile: 150 },
		{ date: new Date('2024-04-02'), desktop: 97, mobile: 180 },
		{ date: new Date('2024-04-03'), desktop: 167, mobile: 120 },
		{ date: new Date('2024-04-04'), desktop: 242, mobile: 260 },
		{ date: new Date('2024-04-05'), desktop: 373, mobile: 290 },
		{ date: new Date('2024-04-06'), desktop: 301, mobile: 340 },
		{ date: new Date('2024-04-07'), desktop: 245, mobile: 180 },
		{ date: new Date('2024-04-08'), desktop: 409, mobile: 320 },
		{ date: new Date('2024-04-09'), desktop: 59, mobile: 110 },
		{ date: new Date('2024-04-10'), desktop: 261, mobile: 190 },
		{ date: new Date('2024-04-11'), desktop: 327, mobile: 350 },
		{ date: new Date('2024-04-12'), desktop: 292, mobile: 210 },
		{ date: new Date('2024-04-13'), desktop: 342, mobile: 380 },
		{ date: new Date('2024-04-14'), desktop: 137, mobile: 220 },
		{ date: new Date('2024-04-15'), desktop: 120, mobile: 170 },
		{ date: new Date('2024-04-16'), desktop: 138, mobile: 190 },
		{ date: new Date('2024-04-17'), desktop: 446, mobile: 360 },
		{ date: new Date('2024-04-18'), desktop: 364, mobile: 410 },
		{ date: new Date('2024-04-19'), desktop: 243, mobile: 180 },
		{ date: new Date('2024-04-20'), desktop: 89, mobile: 150 },
		{ date: new Date('2024-04-21'), desktop: 137, mobile: 200 },
		{ date: new Date('2024-04-22'), desktop: 224, mobile: 170 },
		{ date: new Date('2024-04-23'), desktop: 138, mobile: 230 },
		{ date: new Date('2024-04-24'), desktop: 387, mobile: 290 },
		{ date: new Date('2024-04-25'), desktop: 215, mobile: 250 },
		{ date: new Date('2024-04-26'), desktop: 75, mobile: 130 },
		{ date: new Date('2024-04-27'), desktop: 383, mobile: 420 },
		{ date: new Date('2024-04-28'), desktop: 122, mobile: 180 },
		{ date: new Date('2024-04-29'), desktop: 315, mobile: 240 },
		{ date: new Date('2024-04-30'), desktop: 454, mobile: 380 },
		{ date: new Date('2024-05-01'), desktop: 165, mobile: 220 },
		{ date: new Date('2024-05-02'), desktop: 293, mobile: 310 },
		{ date: new Date('2024-05-03'), desktop: 247, mobile: 190 },
		{ date: new Date('2024-05-04'), desktop: 385, mobile: 420 },
		{ date: new Date('2024-05-05'), desktop: 481, mobile: 390 },
		{ date: new Date('2024-05-06'), desktop: 498, mobile: 520 },
		{ date: new Date('2024-05-07'), desktop: 388, mobile: 300 },
		{ date: new Date('2024-05-08'), desktop: 149, mobile: 210 },
		{ date: new Date('2024-05-09'), desktop: 227, mobile: 180 },
		{ date: new Date('2024-05-10'), desktop: 293, mobile: 330 },
		{ date: new Date('2024-05-11'), desktop: 335, mobile: 270 },
		{ date: new Date('2024-05-12'), desktop: 197, mobile: 240 },
		{ date: new Date('2024-05-13'), desktop: 197, mobile: 160 },
		{ date: new Date('2024-05-14'), desktop: 448, mobile: 490 },
		{ date: new Date('2024-05-15'), desktop: 473, mobile: 380 },
		{ date: new Date('2024-05-16'), desktop: 338, mobile: 400 },
		{ date: new Date('2024-05-17'), desktop: 499, mobile: 420 },
		{ date: new Date('2024-05-18'), desktop: 315, mobile: 350 },
		{ date: new Date('2024-05-19'), desktop: 235, mobile: 180 },
		{ date: new Date('2024-05-20'), desktop: 177, mobile: 230 },
		{ date: new Date('2024-05-21'), desktop: 82, mobile: 140 },
		{ date: new Date('2024-05-22'), desktop: 81, mobile: 120 },
		{ date: new Date('2024-05-23'), desktop: 252, mobile: 290 },
		{ date: new Date('2024-05-24'), desktop: 294, mobile: 220 },
		{ date: new Date('2024-05-25'), desktop: 201, mobile: 250 },
		{ date: new Date('2024-05-26'), desktop: 213, mobile: 170 },
		{ date: new Date('2024-05-27'), desktop: 420, mobile: 460 },
		{ date: new Date('2024-05-28'), desktop: 233, mobile: 190 },
		{ date: new Date('2024-05-29'), desktop: 78, mobile: 130 },
		{ date: new Date('2024-05-30'), desktop: 340, mobile: 280 },
		{ date: new Date('2024-05-31'), desktop: 178, mobile: 230 },
		{ date: new Date('2024-06-01'), desktop: 178, mobile: 200 },
		{ date: new Date('2024-06-02'), desktop: 470, mobile: 410 },
		{ date: new Date('2024-06-03'), desktop: 103, mobile: 160 },
		{ date: new Date('2024-06-04'), desktop: 439, mobile: 380 },
		{ date: new Date('2024-06-05'), desktop: 88, mobile: 140 },
		{ date: new Date('2024-06-06'), desktop: 294, mobile: 250 },
		{ date: new Date('2024-06-07'), desktop: 323, mobile: 370 },
		{ date: new Date('2024-06-08'), desktop: 385, mobile: 320 },
		{ date: new Date('2024-06-09'), desktop: 438, mobile: 480 },
		{ date: new Date('2024-06-10'), desktop: 155, mobile: 200 },
		{ date: new Date('2024-06-11'), desktop: 92, mobile: 150 },
		{ date: new Date('2024-06-12'), desktop: 492, mobile: 420 },
		{ date: new Date('2024-06-13'), desktop: 81, mobile: 130 },
		{ date: new Date('2024-06-14'), desktop: 426, mobile: 380 },
		{ date: new Date('2024-06-15'), desktop: 307, mobile: 350 },
		{ date: new Date('2024-06-16'), desktop: 371, mobile: 310 },
		{ date: new Date('2024-06-17'), desktop: 475, mobile: 520 },
		{ date: new Date('2024-06-18'), desktop: 107, mobile: 170 },
		{ date: new Date('2024-06-19'), desktop: 341, mobile: 290 },
		{ date: new Date('2024-06-20'), desktop: 408, mobile: 450 },
		{ date: new Date('2024-06-21'), desktop: 169, mobile: 210 },
		{ date: new Date('2024-06-22'), desktop: 317, mobile: 270 },
		{ date: new Date('2024-06-23'), desktop: 480, mobile: 530 },
		{ date: new Date('2024-06-24'), desktop: 132, mobile: 180 },
		{ date: new Date('2024-06-25'), desktop: 141, mobile: 190 },
		{ date: new Date('2024-06-26'), desktop: 434, mobile: 380 },
		{ date: new Date('2024-06-27'), desktop: 448, mobile: 490 },
		{ date: new Date('2024-06-28'), desktop: 149, mobile: 200 },
		{ date: new Date('2024-06-29'), desktop: 103, mobile: 160 },
		{ date: new Date('2024-06-30'), desktop: 446, mobile: 400 }
	];
	const chartConfig4 = {
		views: { label: 'Page Views', color: '' },
		desktop: { label: 'Desktop', color: 'var(--chart-2)' },
		mobile: { label: 'Mobile', color: 'var(--chart-4)' }
	} satisfies Chart.ChartConfig;
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

	//
	const chartData5 = [
		{ date: new Date('2024-01-01'), desktop: 186 },
		{ date: new Date('2024-02-01'), desktop: 305 },
		{ date: new Date('2024-03-01'), desktop: 237 },
		{ date: new Date('2024-04-01'), desktop: 73 },
		{ date: new Date('2024-05-01'), desktop: 209 },
		{ date: new Date('2024-06-01'), desktop: 214 }
	];
	const chartConfig5 = {
		desktop: { label: 'Desktop', color: 'var(--chart-5)' }
	} satisfies Chart.ChartConfig;

	//
	const chartData6 = [
		{ browser: 'PCIe', visitors: 275, color: 'var(--color-chrome)' },
		{ browser: 'SSD', visitors: 200, color: 'var(--color-safari)' }
	];
	const chartConfig6 = {
		visitors: { label: 'Visitors' },
		chrome: { label: 'Chrome', color: 'var(--chart-1)' },
		safari: { label: 'Safari', color: 'var(--chart-2)' },
		firefox: { label: 'Firefox', color: 'var(--chart-3)' },
		edge: { label: 'Edge', color: 'var(--chart-4)' },
		other: { label: 'Other', color: 'var(--chart-5)' }
	} satisfies Chart.ChartConfig;
	const totalVisitors6 = chartData6.reduce((acc, curr) => acc + curr.visitors, 0);

	//
	const transport: Transport = getContext('transport');
	const environmentService = createClient(EnvironmentService, transport);

	let prometheusDriver = $state<PrometheusDriver | null>(null);
	let mounted = $state(false);

	async function initializePrometheusDriver(): Promise<PrometheusDriver> {
		const response = await environmentService.getPrometheus({});
		return new PrometheusDriver({
			endpoint: `${env.PUBLIC_API_URL}/prometheus`,
			baseURL: response.baseUrl
		});
	}

	onMount(async () => {
		try {
			prometheusDriver = await initializePrometheusDriver();
		} catch (error) {
			console.error('Failed to initialize Prometheus driver:', error);
		} finally {
			mounted = true;
		}
	});
</script>

<div class="mx-auto grid w-full gap-6">
	<div class="grid gap-1">
		<h1 class="text-2xl font-bold tracking-tight md:text-3xl">Dashboard</h1>
		<p class="text-muted-foreground">description</p>
	</div>

	<Tabs.Root value="overview">
		<Tabs.List>
			<Tabs.Trigger value="overview">Overview</Tabs.Trigger>
			<Tabs.Trigger value="analytics">Analytics</Tabs.Trigger>
		</Tabs.List>
		<Tabs.Content
			value="overview"
			class="grid auto-rows-auto grid-cols-2 gap-5 pt-4 md:grid-cols-4 lg:grid-cols-10"
		>
			<Card.Root class="col-span-2 gap-2">
				<Card.Header class="items-center">
					<Card.Title>Health</Card.Title>
					<Card.Description>xxx</Card.Description>
				</Card.Header>
				<Card.Content class="flex-1">OK</Card.Content>
			</Card.Root>

			<Card.Root class="col-span-2 gap-2">
				<Card.Header class="items-center">
					<Card.Title>Time till full (Grafana 有)</Card.Title>
					<Card.Description>
						<!-- hover info -->
						<!-- Time till pool is full assuming the average fill rate of the last 6 hours -->
					</Card.Description>
				</Card.Header>
				<Card.Content class="flex-1">∞ year</Card.Content>
			</Card.Root>

			<Card.Root class="col-span-2 row-span-2 gap-2">
				<Card.Header class="items-center">
					<Card.Title>剩餘使用容量</Card.Title>
					<Card.Description>xxx</Card.Description>
				</Card.Header>
				<Card.Content class="flex-1">
					<Chart.Container config={chartConfig1} class="mx-auto aspect-square max-h-[200px]">
						<ArcChart
							label="browser"
							value="visitors"
							outerRadius={88}
							innerRadius={66}
							trackOuterRadius={83}
							trackInnerRadius={72}
							padding={40}
							range={[90, -270]}
							maxValue={chartData1[0].visitors * 4}
							series={chartData1.map((d) => ({
								key: d.browser,
								color: d.color,
								data: [d]
							}))}
							props={{
								arc: { track: { fill: 'var(--muted)' }, motion: 'tween' },
								tooltip: { context: { hideDelay: 350 } }
							}}
							tooltip={false}
						>
							{#snippet belowMarks()}
								<circle cx="0" cy="0" r="80" class="fill-background" />
							{/snippet}
							{#snippet aboveMarks()}
								<Text
									value={String(chartData1[0].visitors)}
									textAnchor="middle"
									verticalAnchor="middle"
									class="fill-foreground text-4xl! font-bold"
									dy={3}
								/>
								<Text
									value="Visitors"
									textAnchor="middle"
									verticalAnchor="middle"
									class="fill-muted-foreground!"
									dy={22}
								/>
							{/snippet}
						</ArcChart>
					</Chart.Container>
				</Card.Content>
			</Card.Root>

			<Card.Root class="col-span-2 row-span-2 gap-2">
				<Card.Header>
					<Card.Title>capacity 成長 by 天</Card.Title>
					<Card.Description>xxx</Card.Description>
				</Card.Header>
				<Card.Content>
					<Chart.Container config={chartConfig2} class="h-[200px] w-full">
						<AreaChart
							data={chartData2}
							x="date"
							xScale={scaleUtc()}
							series={[
								{
									key: 'desktop',
									label: 'Desktop',
									color: chartConfig2.desktop.color
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
									format: (v: Date) => v.toLocaleDateString('en-US', { month: 'short' })
								},
								yAxis: { format: () => '' }
							}}
						>
							{#snippet tooltip()}
								<Chart.Tooltip hideLabel />
							{/snippet}
						</AreaChart>
					</Chart.Container>
				</Card.Content>
			</Card.Root>

			<Card.Root class="col-span-2 row-span-2 gap-2">
				<Card.Header class="items-center">
					<Card.Title>OSD Type</Card.Title>
					<Card.Description>xxx</Card.Description>
				</Card.Header>
				<Card.Content class="flex-1">
					<Chart.Container config={chartConfig6} class="mx-auto aspect-square max-h-[200px]">
						<PieChart
							data={chartData6}
							key="browser"
							value="visitors"
							c="color"
							innerRadius={60}
							padding={28}
							props={{ pie: { motion: 'tween' } }}
						>
							{#snippet aboveMarks()}
								<Text
									value={String(totalVisitors6)}
									textAnchor="middle"
									verticalAnchor="middle"
									class="fill-foreground text-3xl! font-bold"
									dy={3}
								/>
								<Text
									value="Visitors"
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
				</Card.Content>
			</Card.Root>

			<Card.Root class="col-span-2 gap-2">
				<Card.Header class="items-center">
					<Card.Title>Monitors In Quorum</Card.Title>
					<Card.Description>xxx</Card.Description>
				</Card.Header>
				<Card.Content class="flex-1">1 / 3</Card.Content>
			</Card.Root>

			<Card.Root class="col-span-2 gap-2">
				<Card.Header class="items-center">
					<Card.Title>OSDs IN & UP</Card.Title>
					<Card.Description>xxx</Card.Description>
				</Card.Header>
				<Card.Content class="flex-1">2 / 2</Card.Content>
			</Card.Root>

			<Card.Root class="col-span-4 row-span-2 gap-2">
				<Card.Header class="flex flex-col items-stretch space-y-0 border-b p-0 sm:flex-row">
					<div class="flex flex-1 flex-col justify-center gap-1 px-6 py-5 sm:py-6">
						<Card.Title>Current Throughput</Card.Title>
						<Card.Description>xxx</Card.Description>
					</div>
					<div class="flex">
						{#each ['desktop', 'mobile'] as key (key)}
							{@const chart = key as keyof typeof chartConfig4}
							<button
								data-active={activeChart4 === chart}
								class="data-[active=true]:bg-muted/50 relative z-30 flex flex-1 flex-col justify-center gap-1 border-t px-6 py-4 text-left even:border-l sm:border-l sm:border-t-0 sm:px-8 sm:py-6"
								onclick={() => (activeChart4 = chart)}
							>
								<span class="text-muted-foreground text-xs">
									{chartConfig4[chart].label}
								</span>
								<span class="flex items-end gap-1 text-lg font-bold leading-none sm:text-3xl">
									{total4[key as keyof typeof total4].toLocaleString()}
									<span class="text-muted-foreground text-xs">B/s</span>
								</span>
							</button>
						{/each}
					</div>
				</Card.Header>
				<Card.Content class="px-2 sm:p-6">
					<Chart.Container config={chartConfig4} class="aspect-auto h-[150px] w-full">
						<BarChart
							bind:context={context4}
							data={chartData4}
							x="date"
							axis="x"
							series={activeSeries4}
							props={{
								bars: {
									stroke: 'none',
									rounded: 'none',
									// use the height of the chart to animate the bars
									initialY: context4?.height,
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
											month: 'short',
											day: '2-digit'
										});
									},
									ticks: (scale) => scaleUtc(scale.domain(), scale.range()).ticks()
								}
							}}
						>
							{#snippet belowMarks()}
								<Highlight area={{ class: 'fill-muted' }} />
							{/snippet}
							{#snippet tooltip()}
								<Chart.Tooltip
									nameKey="views"
									labelFormatter={(v: Date) => {
										return v.toLocaleDateString(getLocale(), {
											month: 'short',
											day: 'numeric',
											year: 'numeric'
										});
									}}
								/>
							{/snippet}
						</BarChart>
					</Chart.Container>
				</Card.Content>
			</Card.Root>

			<Card.Root class="col-span-4 row-span-2 gap-2">
				<Card.Header class="flex flex-col items-stretch space-y-0 border-b p-0 sm:flex-row">
					<div class="flex flex-1 flex-col justify-center gap-1 px-6 py-5 sm:py-6">
						<Card.Title>Current IOPS</Card.Title>
						<Card.Description>xxx</Card.Description>
					</div>
					<div class="flex">
						{#each ['desktop', 'mobile'] as key (key)}
							{@const chart = key as keyof typeof chartConfig3}
							<button
								data-active={activeChart3 === chart}
								class="data-[active=true]:bg-muted/50 relative z-30 flex flex-1 flex-col justify-center gap-1 border-t px-6 py-4 text-left even:border-l sm:border-l sm:border-t-0 sm:px-8 sm:py-6"
								onclick={() => (activeChart3 = chart)}
							>
								<span class="text-muted-foreground text-xs">
									{chartConfig3[chart].label}
								</span>
								<span class="flex items-end gap-1 text-lg font-bold leading-none sm:text-3xl">
									{total3[key as keyof typeof total3].toLocaleString()}
									<span class="text-muted-foreground text-xs">ops/s</span>
								</span>
							</button>
						{/each}
					</div>
				</Card.Header>
				<Card.Content class="px-2 sm:p-6">
					<Chart.Container config={chartConfig3} class="aspect-auto h-[150px] w-full">
						<BarChart
							bind:context={context3}
							data={chartData3}
							x="date"
							axis="x"
							series={activeSeries3}
							props={{
								bars: {
									stroke: 'none',
									rounded: 'none',
									// use the height of the chart to animate the bars
									initialY: context3?.height,
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
											month: 'short',
											day: '2-digit'
										});
									},
									ticks: (scale) => scaleUtc(scale.domain(), scale.range()).ticks()
								}
							}}
						>
							{#snippet belowMarks()}
								<Highlight area={{ class: 'fill-muted' }} />
							{/snippet}
							{#snippet tooltip()}
								<Chart.Tooltip
									nameKey="views"
									labelFormatter={(v: Date) => {
										return v.toLocaleDateString(getLocale(), {
											month: 'short',
											day: 'numeric',
											year: 'numeric'
										});
									}}
								/>
							{/snippet}
						</BarChart>
					</Chart.Container>
				</Card.Content>
			</Card.Root>

			<Card.Root class="col-span-2 gap-2">
				<Card.Header>
					<Card.Title>OSD Read Latencies</Card.Title>
					<Card.Description>xxx</Card.Description>
				</Card.Header>
				<Card.Content>
					<Chart.Container config={chartConfig5} class="h-[64px] w-full">
						<LineChart
							data={chartData5}
							x="date"
							xScale={scaleUtc()}
							axis="x"
							series={[
								{
									key: 'desktop',
									label: 'Desktop',
									color: chartConfig5.desktop.color
								}
							]}
							props={{
								spline: { curve: curveLinear, motion: 'tween', strokeWidth: 2 },
								xAxis: {
									format: (v: Date) => v.toLocaleDateString(getLocale(), { month: 'short' })
								},
								highlight: { points: { r: 4 } }
							}}
						>
							{#snippet tooltip()}
								<Chart.Tooltip hideLabel />
							{/snippet}
						</LineChart>
					</Chart.Container>
				</Card.Content>
			</Card.Root>

			<Card.Root class="col-span-2 gap-2">
				<Card.Header>
					<Card.Title>OSD Write Latencies</Card.Title>
					<Card.Description>xxx</Card.Description>
				</Card.Header>
				<Card.Content>
					<Chart.Container config={chartConfig5} class="h-[64px] w-full">
						<LineChart
							data={chartData5}
							x="date"
							xScale={scaleUtc()}
							axis="x"
							series={[
								{
									key: 'desktop',
									label: 'Desktop',
									color: chartConfig5.desktop.color
								}
							]}
							props={{
								spline: { curve: curveLinear, motion: 'tween', strokeWidth: 2 },
								xAxis: {
									format: (v: Date) => v.toLocaleDateString(getLocale(), { month: 'short' })
								},
								highlight: { points: { r: 4 } }
							}}
						>
							{#snippet tooltip()}
								<Chart.Tooltip hideLabel />
							{/snippet}
						</LineChart>
					</Chart.Container>
				</Card.Content>
			</Card.Root>
		</Tabs.Content>
		<Tabs.Content value="analytics">
			{#if mounted && prometheusDriver && $activeScope}
				<Dashboard client={prometheusDriver} scope={$activeScope} />
			{:else}
				{console.log(mounted)}
				{console.log(prometheusDriver)}
				{console.log(activeScope)}
				<Loading />
			{/if}
		</Tabs.Content>
	</Tabs.Root>
</div>
