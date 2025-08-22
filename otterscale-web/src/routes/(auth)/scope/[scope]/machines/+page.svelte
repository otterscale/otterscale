<script lang="ts">
	import { page } from '$app/state';
	import { env } from '$env/dynamic/public';
	import { EnvironmentService } from '$lib/api/environment/v1/environment_pb';
	import { MachineService, type Machine } from '$lib/api/machine/v1/machine_pb';
	import { Dashboard } from '$lib/components/machines/dashboard';
	import { dynamicPaths } from '$lib/path';
	import { breadcrumb } from '$lib/stores';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { PrometheusDriver } from 'prometheus-query';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import * as Tabs from '$lib/components/ui/tabs';
	import * as Card from '$lib/components/ui/card';
	import { LineChart } from 'layerchart';
	import TrendingUpIcon from '@lucide/svelte/icons/trending-up';
	import { scaleUtc } from 'd3-scale';
	import { curveNatural } from 'd3-shape';
	import * as Chart from '$lib/components/ui/chart';
	import { PieChart, Text } from 'layerchart';
	import { Area, AreaChart, LinearGradient } from 'layerchart';
	import { formatCapacity } from '$lib/formatter';
	import { timestampDate } from '@bufbuild/protobuf/wkt';
	import { getLocale } from '$lib/paraglide/runtime';
	import { m } from '$lib/paraglide/messages';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { buttonVariants } from '$lib/components/ui/button';
	import * as HoverCard from '$lib/components/ui/hover-card';
	import { scaleBand } from 'd3-scale';
	import { curveLinear } from 'd3-shape';
	import { Separator } from '$lib/components/ui/separator';

	let totalCPUCores = $state(0);
	let totalMemoryBytes = $state(0);
	let totalStorageBytes = $state(0);
	let totalDisks = $state(0);
	let totalNodes = $state(0);

	const lineChartConfig = {
		node: { label: 'Node', color: 'var(--chart-1)' }
	} satisfies Chart.ChartConfig;

	let lineChartData: {
		date: Date;
		node: number;
	}[] = $state([]);

	const pieChartConfig = {
		nodes: { label: 'Nodes' },
		physical: { label: 'Physical', color: 'var(--chart-1)' },
		virtual: { label: 'Virtual', color: 'var(--chart-2)' }
	} satisfies Chart.ChartConfig;

	let pieChartData: {
		node: string;
		nodes: number;
		color: string;
	}[] = $state([]);

	//

	const chartData3 = [
		{ date: new Date('2024-01-01'), desktop: 186, mobile: 80 },
		{ date: new Date('2024-02-01'), desktop: 305, mobile: 200 },
		{ date: new Date('2024-03-01'), desktop: 237, mobile: 120 },
		{ date: new Date('2024-04-01'), desktop: 73, mobile: 190 },
		{ date: new Date('2024-05-01'), desktop: 209, mobile: 130 },
		{ date: new Date('2024-06-01'), desktop: 214, mobile: 140 }
	];

	const chartConfig3 = {
		desktop: { label: 'Desktop', color: 'var(--chart-1)' },
		mobile: { label: 'Mobile', color: 'var(--chart-2)' }
	} satisfies Chart.ChartConfig;

	const exampleChartData = [
		{ date: new Date('2024-01-01'), desktop: 186 },
		{ date: new Date('2024-02-01'), desktop: 305 },
		{ date: new Date('2024-03-01'), desktop: 237 },
		{ date: new Date('2024-04-01'), desktop: 73 },
		{ date: new Date('2024-05-01'), desktop: 209 },
		{ date: new Date('2024-06-01'), desktop: 214 }
	];

	const exampleChartConfig = {
		desktop: { label: 'Desktop', color: 'var(--chart-1)' }
	} satisfies Chart.ChartConfig;

	// Set breadcrumb navigation
	breadcrumb.set({ parents: [], current: dynamicPaths.machines(page.params.scope) });

	const transport: Transport = getContext('transport');
	const machineClient = createClient(MachineService, transport);
	const environmentService = createClient(EnvironmentService, transport);
	const machinesStore = writable<Machine[]>([]);
	let prometheusDriver: PrometheusDriver | null = null;

	async function initializePrometheusDriver(): Promise<PrometheusDriver> {
		try {
			const response = await environmentService.getPrometheus({});
			return new PrometheusDriver({
				endpoint: `${env.PUBLIC_API_URL}/prometheus`,
				baseURL: response.baseUrl
			});
		} catch (error) {
			console.error('Error initializing Prometheus driver:', error);
			throw error;
		}
	}

	async function fetchMachines(): Promise<void> {
		try {
			const response = await machineClient.listMachines({});
			machinesStore.set(response.machines);
		} catch (error) {
			console.error('Error fetching machines:', error);
		}
	}

	let mounted = $state(false);
	onMount(async () => {
		try {
			prometheusDriver = await initializePrometheusDriver();
			await fetchMachines();
			const scopeMachines = $machinesStore.filter((m) =>
				m.workloadAnnotations['juju-machine-id']?.startsWith(page.params.scope!)
			);
			totalNodes = scopeMachines.length;
			totalCPUCores = scopeMachines.reduce((sum, m) => sum + Number(m.cpuCount ?? 0), 0);
			totalMemoryBytes =
				scopeMachines.reduce((sum, m) => sum + Number(m.memoryMb ?? 0), 0) * 1024 * 1024;

			const blockDevices = scopeMachines.flatMap((m) => m.blockDevices).filter((d) => !d.bootDisk);
			totalDisks = blockDevices.length;
			totalStorageBytes =
				blockDevices.reduce((sum, m) => sum + Number(m.storageMb ?? 0), 0) * 1024 * 1024;

			const virtualNodes = scopeMachines.filter((m) => m.tags.includes('virtual')).length;
			const physicalNodes = scopeMachines.length - virtualNodes;

			pieChartData = [
				{ node: 'physical', nodes: physicalNodes, color: 'var(--color-physical)' },
				{ node: 'virtual', nodes: virtualNodes, color: 'var(--color-virtual)' }
			];

			const monthlyCounts: Record<string, number> = {};
			scopeMachines.forEach((m) => {
				const dateStr = m.lastCommissioned
					? timestampDate(m.lastCommissioned).toISOString().slice(0, 7) // yyyy-mm
					: null;
				if (dateStr) {
					monthlyCounts[dateStr] = (monthlyCounts[dateStr] || 0) + 1;
				}
			});

			const now = new Date();
			const months: string[] = [];
			for (let i = 5; i >= 0; i--) {
				const d = new Date(now.getFullYear(), now.getMonth() - i + 1, 1);
				months.push(d.toISOString().slice(0, 7));
			}
			lineChartData = months.map((month) => ({
				date: new Date(month + '-01'),
				node: monthlyCounts[month] || 0
			}));
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
		mounted = true;
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
			class="grid auto-rows-auto grid-cols-3 gap-5 pt-4 md:grid-cols-6 lg:grid-cols-9"
		>
			<Card.Root class="col-span-2 gap-2">
				<Card.Header>
					<Card.Title class="flex flex-wrap items-center justify-between gap-6">
						<div class="flex items-center gap-2 truncate text-sm font-medium tracking-tight">
							<Icon icon="ph:cpu" class="size-4.5" />
							{m.cpu()}
						</div>
						<Tooltip.Provider>
							<Tooltip.Root>
								<Tooltip.Trigger class={buttonVariants({ variant: 'ghost', size: 'icon' })}>
									<Icon icon="ph:info" class="text-muted-foreground size-5" />
								</Tooltip.Trigger>
								<Tooltip.Content>
									<p>Total number of CPU cores across all machines.</p>
								</Tooltip.Content>
							</Tooltip.Root>
						</Tooltip.Provider>
					</Card.Title>
				</Card.Header>
				<Card.Content class="flex flex-wrap items-center justify-between gap-6">
					<div class="flex flex-col gap-0.5">
						<div class="text-3xl font-bold">{totalCPUCores}</div>
						<p class="text-muted-foreground text-sm">{m.cores()}</p>
					</div>
					{@render cpuLineChart()}
				</Card.Content>
				<Card.Footer
					class="flex flex-wrap items-center justify-end text-sm leading-none font-medium text-emerald-500 dark:text-emerald-400"
				>
					15.54% <Icon icon="ph:caret-up" />
				</Card.Footer>
			</Card.Root>

			<Card.Root class="col-span-2 gap-2">
				<Card.Header>
					<Card.Title class="flex flex-wrap items-center justify-between gap-6">
						<div class="flex items-center gap-2 truncate text-sm font-medium tracking-tight">
							<Icon icon="ph:memory" class="size-4.5" />
							{m.memory()}
						</div>
						<Tooltip.Provider>
							<Tooltip.Root>
								<Tooltip.Trigger class={buttonVariants({ variant: 'ghost', size: 'icon' })}>
									<Icon icon="ph:info" class="text-muted-foreground size-5" />
								</Tooltip.Trigger>
								<Tooltip.Content>
									<p>Total amount of memory across all machines.</p>
								</Tooltip.Content>
							</Tooltip.Root>
						</Tooltip.Provider>
					</Card.Title>
				</Card.Header>
				<Card.Content class="flex flex-col gap-0.5">
					<div class="flex flex-wrap items-center justify-between gap-6">
						<div class="text-3xl font-bold">
							{formatCapacity(totalMemoryBytes).value}
							{formatCapacity(totalMemoryBytes).unit}
						</div>
						!!CHART!!
					</div>
					<p class="text-muted-foreground text-sm">123</p>
				</Card.Content>
				<Card.Footer
					class="flex flex-wrap items-center justify-end text-sm leading-none font-medium text-emerald-500 dark:text-emerald-400"
				>
					15.54% <Icon icon="ph:caret-up" />
				</Card.Footer>
			</Card.Root>

			<Card.Root class="col-span-2 gap-2">
				<Card.Header>
					<Card.Title class="flex flex-wrap items-center justify-between gap-6">
						<div class="flex items-center gap-2 truncate text-sm font-medium tracking-tight">
							<Icon icon="ph:hard-drives" class="size-4.5" />
							{m.storage()}
						</div>
						<Tooltip.Provider>
							<Tooltip.Root>
								<Tooltip.Trigger class={buttonVariants({ variant: 'ghost', size: 'icon' })}>
									<Icon icon="ph:info" class="text-muted-foreground size-5" />
								</Tooltip.Trigger>
								<Tooltip.Content>
									<p>
										Total storage capacity across all machines, including both used and available
										disk space, but excluding boot disks.
									</p>
								</Tooltip.Content>
							</Tooltip.Root>
						</Tooltip.Provider>
					</Card.Title>
				</Card.Header>
				<Card.Content class="flex flex-col gap-0.5">
					<div class="flex flex-wrap items-center justify-between gap-6">
						<div class="text-3xl font-bold">
							{formatCapacity(totalStorageBytes).value}
							{formatCapacity(totalStorageBytes).unit}
						</div>
						!!CHART!!
					</div>
					<p class="text-muted-foreground text-sm">OVER {totalDisks} DISKS</p>
				</Card.Content>
				<Card.Footer
					class="flex flex-wrap items-center justify-end text-sm leading-none font-medium text-emerald-500 dark:text-emerald-400"
				>
					15.54% <Icon icon="ph:caret-up" />
				</Card.Footer>
			</Card.Root>

			<Card.Root class="col-span-3 gap-2">
				<Card.Header>
					<Card.Title class="flex flex-wrap items-center justify-between gap-6">
						<div
							class="flex flex-col items-start gap-0.5 truncate text-sm font-medium tracking-tight"
						>
							<p class="text-muted-foreground text-xs">OVER THE PAST 6 MONTHS</p>
							<div class="flex items-center gap-1 text-lg font-medium">
								<Icon icon="ph:trend-up" class="size-4.5" /> +{totalNodes}
							</div>
						</div>
						<Tooltip.Provider>
							<Tooltip.Root>
								<Tooltip.Trigger class={buttonVariants({ variant: 'ghost', size: 'icon' })}>
									<Icon icon="ph:info" class="text-muted-foreground size-5" />
								</Tooltip.Trigger>
								<Tooltip.Content>
									<p>Shows the monthly trend of node count over the past 6 months.</p>
								</Tooltip.Content>
							</Tooltip.Root>
						</Tooltip.Provider>
					</Card.Title>
				</Card.Header>
				<Card.Content>
					<Chart.Container config={lineChartConfig} class="h-[80px] w-full px-2 pt-2">
						<LineChart
							points={{ r: 4 }}
							data={lineChartData}
							x="date"
							xScale={scaleUtc()}
							axis="x"
							series={[
								{
									key: 'node',
									label: 'Node',
									color: lineChartConfig.node.color
								}
							]}
							props={{
								spline: { curve: curveNatural, motion: 'tween', strokeWidth: 2 },
								highlight: {
									points: {
										motion: 'none',
										r: 6
									}
								},
								xAxis: {
									format: (v: Date) => v.toLocaleDateString(getLocale(), { month: 'short' })
								}
							}}
						>
							{#snippet tooltip()}
								<Chart.Tooltip hideLabel />
							{/snippet}
						</LineChart>
					</Chart.Container>
				</Card.Content>
			</Card.Root>

			<Card.Root class="col-span-6">
				<Card.Header>
					<Card.Title>Load Average LINUX</Card.Title>
					<Card.Description>xxx</Card.Description>
				</Card.Header>
				<Card.Content>
					{@render exampleLineChart()}
				</Card.Content>
			</Card.Root>

			<Card.Root class="col-span-3 flex flex-col">
				<Card.Header class="gap-0.5">
					<Card.Title>
						<div class="flex items-center gap-1 truncate text-sm font-medium tracking-tight">
							<Icon icon="ph:percent" class="size-4.5" />
							Node Distribution
						</div>
					</Card.Title>
					<Card.Description class="text-xs">
						Shows the proportion of physical and virtual machines in your environment.
					</Card.Description>
				</Card.Header>
				<Card.Content class="flex-1">
					<Chart.Container config={pieChartConfig} class="mx-auto aspect-square max-h-[250px]">
						<PieChart
							data={pieChartData}
							key="node"
							value="nodes"
							c="color"
							innerRadius={60}
							padding={28}
							props={{ pie: { motion: 'tween' } }}
						>
							{#snippet aboveMarks()}
								<Text
									value={String(totalNodes)}
									textAnchor="middle"
									verticalAnchor="middle"
									class="fill-foreground text-3xl! font-bold"
									dy={3}
								/>
								<Text
									value="Nodes"
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
		</Tabs.Content>
		<Tabs.Content value="analytics">
			{#if mounted && prometheusDriver && $machinesStore.length > 0}
				{@const filteredMachines = $machinesStore.filter((machine) =>
					machine.workloadAnnotations?.['juju-machine-id']?.includes('-machine-')
				)}
				{@const allMachine = {
					fqdn: filteredMachines.map((machine) => machine.fqdn).join('|'),
					id: 'All Machine'
				} as Machine}
				{@const machines = [allMachine, ...filteredMachines]}
				<Dashboard client={prometheusDriver} {machines} />
			{:else}
				<div class="flex items-center justify-center p-8">
					<Icon icon="mdi:loading" class="animate-spin text-2xl" />
					<span class="ml-2">Loading machines...</span>
				</div>
			{/if}
		</Tabs.Content>
	</Tabs.Root>
</div>

{#snippet exampleLineChart()}
	<Chart.Container config={chartConfig3} class="h-[200px] w-full">
		<AreaChart
			data={chartData3}
			x="date"
			xScale={scaleUtc()}
			yPadding={[0, 25]}
			series={[
				{
					key: 'mobile',
					label: 'Mobile',
					color: 'var(--color-mobile)'
				},
				{
					key: 'desktop',
					label: 'Desktop',
					color: 'var(--color-desktop)'
				}
			]}
			seriesLayout="stack"
			props={{
				area: {
					curve: curveNatural,
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
				<Chart.Tooltip
					indicator="dot"
					labelFormatter={(v: Date) => {
						return v.toLocaleDateString('en-US', {
							month: 'long'
						});
					}}
				/>
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
{/snippet}

{#snippet cpuLineChart()}
	<Chart.Container config={exampleChartConfig} class="h-full w-20">
		<LineChart
			data={exampleChartData}
			x="date"
			xScale={scaleUtc()}
			axis={false}
			series={[
				{
					key: 'desktop',
					label: 'Desktop',
					color: exampleChartConfig.desktop.color
				}
			]}
			props={{
				spline: { curve: curveLinear, motion: 'tween', strokeWidth: 2 },
				xAxis: {
					format: (v: Date) => v.toLocaleDateString('en-US', { month: 'short' })
				},
				highlight: { points: { r: 4 } }
			}}
		>
			{#snippet tooltip()}
				<Chart.Tooltip hideLabel />
			{/snippet}
		</LineChart>
	</Chart.Container>
{/snippet}

{#snippet memoryLineChart()}
	<Chart.Container config={exampleChartConfig} class="h-full w-20">
		<LineChart
			data={exampleChartData}
			x="date"
			xScale={scaleUtc()}
			axis={false}
			series={[
				{
					key: 'desktop',
					label: 'Desktop',
					color: exampleChartConfig.desktop.color
				}
			]}
			props={{
				spline: { curve: curveLinear, motion: 'tween', strokeWidth: 2 },
				xAxis: {
					format: (v: Date) => v.toLocaleDateString('en-US', { month: 'short' })
				},
				highlight: { points: { r: 4 } }
			}}
		>
			{#snippet tooltip()}
				<Chart.Tooltip hideLabel />
			{/snippet}
		</LineChart>
	</Chart.Container>
{/snippet}

{#snippet storageLineChart()}
	<Chart.Container config={exampleChartConfig} class="h-full w-20">
		<LineChart
			data={exampleChartData}
			x="date"
			xScale={scaleUtc()}
			axis={false}
			series={[
				{
					key: 'desktop',
					label: 'Desktop',
					color: exampleChartConfig.desktop.color
				}
			]}
			props={{
				spline: { curve: curveLinear, motion: 'tween', strokeWidth: 2 },
				xAxis: {
					format: (v: Date) => v.toLocaleDateString('en-US', { month: 'short' })
				},
				highlight: { points: { r: 4 } }
			}}
		>
			{#snippet tooltip()}
				<Chart.Tooltip hideLabel />
			{/snippet}
		</LineChart>
	</Chart.Container>
{/snippet}
