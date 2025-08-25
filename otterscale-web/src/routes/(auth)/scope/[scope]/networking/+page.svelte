<script lang="ts">
	import { page } from '$app/state';
	import { dynamicPaths } from '$lib/path';
	import { breadcrumb } from '$lib/stores';
	import { Construction } from '$lib/components/construction';

	import * as Tabs from '$lib/components/ui/tabs';
	import { ArcChart, Text } from 'layerchart';
	import TrendingUpIcon from '@lucide/svelte/icons/trending-up';

	import * as Chart from '$lib/components/ui/chart/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import { scaleUtc } from 'd3-scale';
	import { BarChart, type ChartContextValue, Highlight } from 'layerchart';
	import { cubicInOut } from 'svelte/easing';
	import { scaleBand } from 'd3-scale';

	import { LineChart } from 'layerchart';
	import { curveLinearClosed } from 'd3-shape';

	// Set breadcrumb navigation
	breadcrumb.set({
		parents: [],
		current: dynamicPaths.networking(page.params.scope)
	});

	//
	const chartData1 = [{ browser: 'safari', visitors: 1260, color: 'var(--color-safari)' }];

	const chartConfig1 = {
		visitors: { label: 'Visitors' },
		safari: { label: 'Safari', color: 'var(--chart-2)' }
	} satisfies Chart.ChartConfig;

	//
	const chartData2 = [
		{ date: new Date('2024-04-01'), receive: 222, transmit: 150 },
		{ date: new Date('2024-04-02'), receive: 97, transmit: 180 },
		{ date: new Date('2024-04-03'), receive: 167, transmit: 120 },
		{ date: new Date('2024-04-04'), receive: 242, transmit: 260 },
		{ date: new Date('2024-04-05'), receive: 373, transmit: 290 },
		{ date: new Date('2024-04-06'), receive: 301, transmit: 340 },
		{ date: new Date('2024-04-07'), receive: 245, transmit: 180 },
		{ date: new Date('2024-04-08'), receive: 409, transmit: 320 },
		{ date: new Date('2024-04-09'), receive: 59, transmit: 110 },
		{ date: new Date('2024-04-10'), receive: 261, transmit: 190 },
		{ date: new Date('2024-04-11'), receive: 327, transmit: 350 },
		{ date: new Date('2024-04-12'), receive: 292, transmit: 210 },
		{ date: new Date('2024-04-13'), receive: 342, transmit: 380 },
		{ date: new Date('2024-04-14'), receive: 137, transmit: 220 },
		{ date: new Date('2024-04-15'), receive: 120, transmit: 170 },
		{ date: new Date('2024-04-16'), receive: 138, transmit: 190 },
		{ date: new Date('2024-04-17'), receive: 446, transmit: 360 },
		{ date: new Date('2024-04-18'), receive: 364, transmit: 410 },
		{ date: new Date('2024-04-19'), receive: 243, transmit: 180 },
		{ date: new Date('2024-04-20'), receive: 89, transmit: 150 },
		{ date: new Date('2024-04-21'), receive: 137, transmit: 200 },
		{ date: new Date('2024-04-22'), receive: 224, transmit: 170 },
		{ date: new Date('2024-04-23'), receive: 138, transmit: 230 },
		{ date: new Date('2024-04-24'), receive: 387, transmit: 290 },
		{ date: new Date('2024-04-25'), receive: 215, transmit: 250 },
		{ date: new Date('2024-04-26'), receive: 75, transmit: 130 },
		{ date: new Date('2024-04-27'), receive: 383, transmit: 420 },
		{ date: new Date('2024-04-28'), receive: 122, transmit: 180 },
		{ date: new Date('2024-04-29'), receive: 315, transmit: 240 },
		{ date: new Date('2024-04-30'), receive: 454, transmit: 380 },
		{ date: new Date('2024-05-01'), receive: 165, transmit: 220 },
		{ date: new Date('2024-05-02'), receive: 293, transmit: 310 },
		{ date: new Date('2024-05-03'), receive: 247, transmit: 190 },
		{ date: new Date('2024-05-04'), receive: 385, transmit: 420 },
		{ date: new Date('2024-05-05'), receive: 481, transmit: 390 },
		{ date: new Date('2024-05-06'), receive: 498, transmit: 520 },
		{ date: new Date('2024-05-07'), receive: 388, transmit: 300 },
		{ date: new Date('2024-05-08'), receive: 149, transmit: 210 },
		{ date: new Date('2024-05-09'), receive: 227, transmit: 180 },
		{ date: new Date('2024-05-10'), receive: 293, transmit: 330 },
		{ date: new Date('2024-05-11'), receive: 335, transmit: 270 },
		{ date: new Date('2024-05-12'), receive: 197, transmit: 240 },
		{ date: new Date('2024-05-13'), receive: 197, transmit: 160 },
		{ date: new Date('2024-05-14'), receive: 448, transmit: 490 },
		{ date: new Date('2024-05-15'), receive: 473, transmit: 380 },
		{ date: new Date('2024-05-16'), receive: 338, transmit: 400 },
		{ date: new Date('2024-05-17'), receive: 499, transmit: 420 },
		{ date: new Date('2024-05-18'), receive: 315, transmit: 350 },
		{ date: new Date('2024-05-19'), receive: 235, transmit: 180 },
		{ date: new Date('2024-05-20'), receive: 177, transmit: 230 },
		{ date: new Date('2024-05-21'), receive: 82, transmit: 140 },
		{ date: new Date('2024-05-22'), receive: 81, transmit: 120 },
		{ date: new Date('2024-05-23'), receive: 252, transmit: 290 },
		{ date: new Date('2024-05-24'), receive: 294, transmit: 220 },
		{ date: new Date('2024-05-25'), receive: 201, transmit: 250 },
		{ date: new Date('2024-05-26'), receive: 213, transmit: 170 },
		{ date: new Date('2024-05-27'), receive: 420, transmit: 460 },
		{ date: new Date('2024-05-28'), receive: 233, transmit: 190 },
		{ date: new Date('2024-05-29'), receive: 78, transmit: 130 },
		{ date: new Date('2024-05-30'), receive: 340, transmit: 280 },
		{ date: new Date('2024-05-31'), receive: 178, transmit: 230 },
		{ date: new Date('2024-06-01'), receive: 178, transmit: 200 },
		{ date: new Date('2024-06-02'), receive: 470, transmit: 410 },
		{ date: new Date('2024-06-03'), receive: 103, transmit: 160 },
		{ date: new Date('2024-06-04'), receive: 439, transmit: 380 },
		{ date: new Date('2024-06-05'), receive: 88, transmit: 140 },
		{ date: new Date('2024-06-06'), receive: 294, transmit: 250 },
		{ date: new Date('2024-06-07'), receive: 323, transmit: 370 },
		{ date: new Date('2024-06-08'), receive: 385, transmit: 320 },
		{ date: new Date('2024-06-09'), receive: 438, transmit: 480 },
		{ date: new Date('2024-06-10'), receive: 155, transmit: 200 },
		{ date: new Date('2024-06-11'), receive: 92, transmit: 150 },
		{ date: new Date('2024-06-12'), receive: 492, transmit: 420 },
		{ date: new Date('2024-06-13'), receive: 81, transmit: 130 },
		{ date: new Date('2024-06-14'), receive: 426, transmit: 380 },
		{ date: new Date('2024-06-15'), receive: 307, transmit: 350 },
		{ date: new Date('2024-06-16'), receive: 371, transmit: 310 },
		{ date: new Date('2024-06-17'), receive: 475, transmit: 520 },
		{ date: new Date('2024-06-18'), receive: 107, transmit: 170 },
		{ date: new Date('2024-06-19'), receive: 341, transmit: 290 },
		{ date: new Date('2024-06-20'), receive: 408, transmit: 450 },
		{ date: new Date('2024-06-21'), receive: 169, transmit: 210 },
		{ date: new Date('2024-06-22'), receive: 317, transmit: 270 },
		{ date: new Date('2024-06-23'), receive: 480, transmit: 530 },
		{ date: new Date('2024-06-24'), receive: 132, transmit: 180 },
		{ date: new Date('2024-06-25'), receive: 141, transmit: 190 },
		{ date: new Date('2024-06-26'), receive: 434, transmit: 380 },
		{ date: new Date('2024-06-27'), receive: 448, transmit: 490 },
		{ date: new Date('2024-06-28'), receive: 149, transmit: 200 },
		{ date: new Date('2024-06-29'), receive: 103, transmit: 160 },
		{ date: new Date('2024-06-30'), receive: 446, transmit: 400 }
	];
	const chartConfig2 = {
		views: { label: 'Page Views', color: '' },
		receive: { label: 'Receive', color: 'var(--chart-1)' },
		transmit: { label: 'Transmit', color: 'var(--chart-2)' }
	} satisfies Chart.ChartConfig;

	let context2 = $state<ChartContextValue>();
	let activeChart = $state<keyof typeof chartConfig2>('receive');
	const total = $derived({
		receive: chartData2.reduce((acc, curr) => acc + curr.receive, 0),
		transmit: chartData2.reduce((acc, curr) => acc + curr.transmit, 0)
	});
	const activeSeries = $derived([
		{
			key: activeChart,
			label: chartConfig2[activeChart].label,
			color: chartConfig2[activeChart].color
		}
	]);

	//
	const chartData3 = [
		{ month: 'January', desktop: 186, mobile: 80 },
		{ month: 'February', desktop: 305, mobile: 200 },
		{ month: 'March', desktop: 237, mobile: 120 },
		{ month: 'April', desktop: 73, mobile: 190 },
		{ month: 'May', desktop: 209, mobile: 130 },
		{ month: 'June', desktop: 214, mobile: 140 }
	];
	const chartConfig3 = {
		desktop: { label: 'Desktop', color: 'var(--chart-1)' },
		mobile: { label: 'Mobile', color: 'var(--chart-2)' }
	} satisfies Chart.ChartConfig;

	let context3 = $state<ChartContextValue>();

	const chartData4 = [
		{ month: 'January', desktop: 186 },
		{ month: 'February', desktop: 305 },
		{ month: 'March', desktop: 237 },
		{ month: 'April', desktop: 273 },
		{ month: 'May', desktop: 209 },
		{ month: 'June', desktop: 214 }
	];
	const chartConfig4 = {
		desktop: { label: 'Desktop', color: 'var(--chart-1)' }
	} satisfies Chart.ChartConfig;
</script>

<div class="mx-auto grid w-full gap-6">
	<div class="grid gap-1">
		<h1 class="text-2xl font-bold tracking-tight md:text-3xl">Dashboard</h1>
		<p class="text-muted-foreground">description</p>
	</div>

	<Tabs.Root value="overview">
		<Tabs.List>
			<Tabs.Trigger value="overview">Overview</Tabs.Trigger>
			<Tabs.Trigger value="analytics" disabled>Analytics</Tabs.Trigger>
		</Tabs.List>
		<Tabs.Content
			value="overview"
			class="grid auto-rows-auto grid-cols-3 gap-5 pt-4 md:grid-cols-6 lg:grid-cols-9"
		>
			<Card.Root class="col-span-1 gap-2">
				<Card.Header>
					<Card.Title>Discovery</Card.Title>
					<Card.Description>Subnet 名稱</Card.Description>
				</Card.Header>
				<Card.Content>On / Off (subnet.active_discovery 欄位)</Card.Content>
			</Card.Root>

			<Card.Root class="col-span-1 gap-2">
				<Card.Header>
					<Card.Title>DHCP</Card.Title>
					<Card.Description>Subnet 名稱</Card.Description>
				</Card.Header>
				<Card.Content>On / Off (vlan.dhcp_on 欄位)</Card.Content>
			</Card.Root>

			<Card.Root class="col-span-2 row-span-2 gap-2">
				<Card.Header class="items-center">
					<Card.Title>目前剩餘可用的 IP 數量</Card.Title>
					<Card.Description>% 數</Card.Description>
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

			<Card.Root class="col-span-5 row-span-2 gap-2">
				<Card.Header class="flex flex-col items-stretch space-y-0 border-b p-0 sm:flex-row">
					<div class="flex flex-1 flex-col justify-center gap-1 px-6 py-5 sm:py-6">
						<Card.Title>Network Traffic</Card.Title>
						<Card.Description>用 fqdn 過濾</Card.Description>
					</div>
					<div class="flex">
						{#each ['receive', 'transmit'] as key (key)}
							{@const chart = key as keyof typeof chartConfig2}
							<button
								data-active={activeChart === chart}
								class="data-[active=true]:bg-muted/50 relative z-30 flex flex-1 flex-col justify-center gap-1 border-t px-6 py-4 text-left even:border-l sm:border-t-0 sm:border-l sm:px-8 sm:py-6"
								onclick={() => (activeChart = chart)}
							>
								<span class="text-muted-foreground text-xs">
									{chartConfig2[chart].label}
								</span>
								<span class="flex items-end gap-1 text-lg leading-none font-bold sm:text-3xl">
									{total[key as keyof typeof total].toLocaleString()}
									<span class="text-xs">Mbps</span>
								</span>
							</button>
						{/each}
					</div>
				</Card.Header>
				<Card.Content class="px-6 pt-6">
					<Chart.Container config={chartConfig2} class="aspect-auto h-[120px] w-full">
						<BarChart
							bind:context={context2}
							data={chartData2}
							x="date"
							axis="x"
							series={activeSeries}
							props={{
								bars: {
									stroke: 'none',
									rounded: 'none',
									// use the height of the chart to animate the bars
									initialY: context2?.height,
									initialHeight: 0,
									motion: {
										y: { type: 'tween', duration: 500, easing: cubicInOut },
										height: { type: 'tween', duration: 500, easing: cubicInOut }
									}
								},
								highlight: { area: { fill: 'none' } },
								xAxis: {
									format: (d: Date) => {
										return d.toLocaleDateString('en-US', {
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
										return v.toLocaleDateString('en-US', {
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
					<Card.Title>DNS Servers</Card.Title>
					<Card.Description>Subnet 名稱</Card.Description>
				</Card.Header>
				<Card.Content>["192.168.1.85"] (subnet.dns_servers 欄位)</Card.Content>
			</Card.Root>

			<Card.Root class="col-span-4 gap-2">
				<Card.Header>
					<Card.Title>上傳下載總和 BY 天</Card.Title>
					<Card.Description>用 fqdn 過濾 過去一周</Card.Description>
				</Card.Header>
				<Card.Content>
					<Chart.Container config={chartConfig3}>
						<BarChart
							bind:context={context3}
							data={chartData3}
							xScale={scaleBand().padding(0.25)}
							x="month"
							axis="x"
							rule={false}
							series={[
								{
									key: 'desktop',
									label: 'Desktop',
									color: chartConfig3.desktop.color,
									props: { rounded: 'bottom' }
								},
								{
									key: 'mobile',
									label: 'Mobile',
									color: chartConfig3.mobile.color
								}
							]}
							seriesLayout="stack"
							props={{
								bars: {
									stroke: 'none',
									initialY: context3?.height,
									initialHeight: 0,
									motion: {
										y: { type: 'tween', duration: 500, easing: cubicInOut },
										height: { type: 'tween', duration: 500, easing: cubicInOut }
									}
								},
								highlight: { area: false },
								xAxis: { format: (d) => d.slice(0, 3) }
							}}
							legend
						>
							{#snippet belowMarks()}
								<Highlight area={{ class: 'fill-muted' }} />
							{/snippet}
							{#snippet tooltip()}
								<Chart.Tooltip />
							{/snippet}
						</BarChart>
					</Chart.Container>
				</Card.Content>
			</Card.Root>
		</Tabs.Content>
		<Tabs.Content value="analytics"></Tabs.Content>
	</Tabs.Root>
</div>
