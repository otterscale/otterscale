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
	import * as Chart from '$lib/components/ui/chart/index.js';
	import { PieChart, Text } from 'layerchart';
	import { Area, AreaChart, LinearGradient } from 'layerchart';
	import { formatCapacity } from '$lib/formatter';
	import { timestampDate } from '@bufbuild/protobuf/wkt';
	import { getLocale } from '$lib/paraglide/runtime';
	import { m } from '$lib/paraglide/messages';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { buttonVariants } from '$lib/components/ui/button';

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

	let mounted = false;

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
			<Card.Content class="flex flex-col gap-0.5">
				<div class="flex flex-wrap items-center justify-between gap-6">
					<div class="text-3xl font-bold">{totalCPUCores}</div>
					!!CHART!!
				</div>
				<p class="text-muted-foreground text-sm">{m.cores()}</p>
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
									Total storage capacity across all machines, including both used and available disk
									space, but excluding boot disks.
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
						<div class="flex items-center gap-2 text-lg font-medium">
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
			</Card.Content>
		</Card.Root>

		<Card.Root class="col-span-3 flex flex-col">
			<Card.Header class="gap-0.5">
				<Card.Title>
					<div class="flex items-center gap-2 truncate text-sm font-medium tracking-tight">
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
	<Tabs.Content value="analytics">AAA</Tabs.Content>
</Tabs.Root>

<!-- {#if mounted && prometheusDriver && $machinesStore.length > 0}
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
{/if} -->

<div class="space-y-4 p-4">
	<div class="mb-2 flex flex-col items-start justify-between space-y-2 md:flex-row md:items-center">
		<h1 class="text-2xl font-bold tracking-tight">Dashboard</h1>
		<div class="flex items-center space-x-2">
			<button
				class="focus-visible:ring-ring [&amp;_svg]:pointer-events-none [&amp;_svg]:size-4 [&amp;_svg]:shrink-0 bg-primary text-primary-foreground hover:bg-primary/90 inline-flex h-9 items-center justify-center gap-2 rounded-md px-4 py-2 text-sm font-medium whitespace-nowrap shadow-sm transition-colors focus-visible:ring-1 focus-visible:outline-hidden disabled:pointer-events-none disabled:opacity-50"
				><svg
					xmlns="http://www.w3.org/2000/svg"
					width="24"
					height="24"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="2"
					stroke-linecap="round"
					stroke-linejoin="round"
					class="tabler-icon tabler-icon-download"
					><path d="M4 17v2a2 2 0 0 0 2 2h12a2 2 0 0 0 2 -2v-2"></path><path d="M7 11l5 5l5 -5"
					></path><path d="M12 4l0 12"></path></svg
				>Download</button
			><button
				class="focus-visible:ring-ring [&amp;_svg]:pointer-events-none [&amp;_svg]:size-4 [&amp;_svg]:shrink-0 border-input bg-background hover:bg-accent hover:text-accent-foreground text-muted-foreground inline-flex h-9 items-center justify-start gap-2 rounded-md border px-4 py-2 text-left text-sm font-normal whitespace-nowrap shadow-xs transition-colors focus-visible:ring-1 focus-visible:outline-hidden disabled:pointer-events-none disabled:opacity-50"
				type="button"
				aria-haspopup="dialog"
				aria-expanded="false"
				aria-controls="radix-:r9i:"
				data-state="closed"
				><svg
					xmlns="http://www.w3.org/2000/svg"
					width="24"
					height="24"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="2"
					stroke-linecap="round"
					stroke-linejoin="round"
					class="lucide lucide-calendar"
					><path d="M8 2v4"></path><path d="M16 2v4"></path><rect
						width="18"
						height="18"
						x="3"
						y="4"
						rx="2"
					></rect><path d="M3 10h18"></path></svg
				><span>Pick a date</span></button
			>
		</div>
	</div>
	<div dir="ltr" data-orientation="vertical" class="space-y-4">
		<div class="w-full overflow-x-auto pb-2">
			<div
				role="tablist"
				aria-orientation="vertical"
				class="bg-muted text-muted-foreground inline-flex h-9 items-center justify-center rounded-lg p-1"
				tabindex="0"
				data-orientation="vertical"
				style="outline: none;"
			>
				<button
					type="button"
					role="tab"
					aria-selected="true"
					aria-controls="radix-:r9j:-content-overview"
					data-state="active"
					id="radix-:r9j:-trigger-overview"
					class="ring-offset-background focus-visible:ring-ring data-[state=active]:bg-background data-[state=active]:text-foreground flex items-center justify-center gap-2 rounded-md px-3 py-1 text-sm font-medium whitespace-nowrap transition-all focus-visible:ring-2 focus-visible:ring-offset-2 focus-visible:outline-hidden disabled:pointer-events-none disabled:opacity-50 data-[state=active]:shadow-sm"
					tabindex="-1"
					data-orientation="vertical"
					data-radix-collection-item=""
					><svg
						xmlns="http://www.w3.org/2000/svg"
						width="14"
						height="14"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="2"
						stroke-linecap="round"
						stroke-linejoin="round"
						class="tabler-icon tabler-icon-settings-2"
						><path
							d="M19.875 6.27a2.225 2.225 0 0 1 1.125 1.948v7.284c0 .809 -.443 1.555 -1.158 1.948l-6.75 4.27a2.269 2.269 0 0 1 -2.184 0l-6.75 -4.27a2.225 2.225 0 0 1 -1.158 -1.948v-7.285c0 -.809 .443 -1.554 1.158 -1.947l6.75 -3.98a2.33 2.33 0 0 1 2.25 0l6.75 3.98h-.033z"
						></path><path d="M12 12m-3 0a3 3 0 1 0 6 0a3 3 0 1 0 -6 0"></path></svg
					>Overview</button
				><button
					type="button"
					role="tab"
					aria-selected="false"
					aria-controls="radix-:r9j:-content-analytics"
					data-state="inactive"
					id="radix-:r9j:-trigger-analytics"
					class="ring-offset-background focus-visible:ring-ring data-[state=active]:bg-background data-[state=active]:text-foreground flex items-center justify-center gap-2 rounded-md px-3 py-1 text-sm font-medium whitespace-nowrap transition-all focus-visible:ring-2 focus-visible:ring-offset-2 focus-visible:outline-hidden disabled:pointer-events-none disabled:opacity-50 data-[state=active]:shadow-sm"
					tabindex="-1"
					data-orientation="vertical"
					data-radix-collection-item=""
					><svg
						xmlns="http://www.w3.org/2000/svg"
						width="16"
						height="16"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="2"
						stroke-linecap="round"
						stroke-linejoin="round"
						class="tabler-icon tabler-icon-analyze"
						><path d="M20 11a8.1 8.1 0 0 0 -6.986 -6.918a8.095 8.095 0 0 0 -8.019 3.918"
						></path><path d="M4 13a8.1 8.1 0 0 0 15 3"></path><path
							d="M19 16m-1 0a1 1 0 1 0 2 0a1 1 0 1 0 -2 0"
						></path><path d="M5 8m-1 0a1 1 0 1 0 2 0a1 1 0 1 0 -2 0"></path><path
							d="M12 12m-3 0a3 3 0 1 0 6 0a3 3 0 1 0 -6 0"
						></path></svg
					>Analytics</button
				><button
					type="button"
					role="tab"
					aria-selected="false"
					aria-controls="radix-:r9j:-content-reports"
					data-state="inactive"
					data-disabled=""
					id="radix-:r9j:-trigger-reports"
					class="ring-offset-background focus-visible:ring-ring data-[state=active]:bg-background data-[state=active]:text-foreground flex items-center justify-center gap-2 rounded-md px-3 py-1 text-sm font-medium whitespace-nowrap transition-all focus-visible:ring-2 focus-visible:ring-offset-2 focus-visible:outline-hidden disabled:pointer-events-none disabled:opacity-50 data-[state=active]:shadow-sm"
					tabindex="-1"
					data-orientation="vertical"
					data-radix-collection-item=""
					><svg
						xmlns="http://www.w3.org/2000/svg"
						width="16"
						height="16"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="2"
						stroke-linecap="round"
						stroke-linejoin="round"
						class="tabler-icon tabler-icon-file-report"
						><path d="M17 17m-4 0a4 4 0 1 0 8 0a4 4 0 1 0 -8 0"></path><path d="M17 13v4h4"
						></path><path d="M12 3v4a1 1 0 0 0 1 1h4"></path><path
							d="M11.5 21h-6.5a2 2 0 0 1 -2 -2v-14a2 2 0 0 1 2 -2h7l5 5v2m0 3v4"
						></path></svg
					>Reports</button
				><button
					type="button"
					role="tab"
					aria-selected="false"
					aria-controls="radix-:r9j:-content-notifications"
					data-state="inactive"
					data-disabled=""
					id="radix-:r9j:-trigger-notifications"
					class="ring-offset-background focus-visible:ring-ring data-[state=active]:bg-background data-[state=active]:text-foreground flex items-center justify-center gap-2 rounded-md px-3 py-1 text-sm font-medium whitespace-nowrap transition-all focus-visible:ring-2 focus-visible:ring-offset-2 focus-visible:outline-hidden disabled:pointer-events-none disabled:opacity-50 data-[state=active]:shadow-sm"
					tabindex="-1"
					data-orientation="vertical"
					data-radix-collection-item=""
					><svg
						xmlns="http://www.w3.org/2000/svg"
						width="16"
						height="16"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="2"
						stroke-linecap="round"
						stroke-linejoin="round"
						class="tabler-icon tabler-icon-notification"
						><path d="M10 6h-3a2 2 0 0 0 -2 2v9a2 2 0 0 0 2 2h9a2 2 0 0 0 2 -2v-3"></path><path
							d="M17 7m-3 0a3 3 0 1 0 6 0a3 3 0 1 0 -6 0"
						></path></svg
					>Notifications</button
				>
			</div>
		</div>
		<div
			data-state="active"
			data-orientation="vertical"
			role="tabpanel"
			aria-labelledby="radix-:r9j:-trigger-overview"
			id="radix-:r9j:-content-overview"
			tabindex="0"
			class="ring-offset-background focus-visible:ring-ring mt-2 space-y-4 focus-visible:ring-2 focus-visible:ring-offset-2 focus-visible:outline-hidden"
			style="animation-duration: 0s;"
		>
			<div class="grid auto-rows-auto grid-cols-3 gap-5 md:grid-cols-6 lg:grid-cols-9">
				<div
					class="bg-card text-card-foreground col-span-3 h-full rounded-xl border shadow-sm lg:col-span-2 xl:col-span-2"
				>
					<div class="flex flex-row items-center justify-between gap-5 space-y-0 p-6 pt-4 pb-2">
						<div class="flex items-center gap-2 truncate text-sm font-medium tracking-tight">
							<svg
								xmlns="http://www.w3.org/2000/svg"
								width="16"
								height="16"
								viewBox="0 0 24 24"
								fill="none"
								stroke="currentColor"
								stroke-width="2"
								stroke-linecap="round"
								stroke-linejoin="round"
								class="tabler-icon tabler-icon-subscript"
								><path d="M5 7l8 10m-8 0l8 -10"></path><path
									d="M21 20h-4l3.5 -4a1.73 1.73 0 0 0 -3.5 -2"
								></path></svg
							>New Subscriptions
						</div>
						<button data-state="closed"
							><svg
								xmlns="http://www.w3.org/2000/svg"
								width="24"
								height="24"
								viewBox="0 0 24 24"
								fill="none"
								stroke="currentColor"
								stroke-width="2"
								stroke-linecap="round"
								stroke-linejoin="round"
								class="tabler-icon tabler-icon-info-circle text-muted-foreground scale-90 stroke-[1.25]"
								><path d="M3 12a9 9 0 1 0 18 0a9 9 0 0 0 -18 0"></path><path d="M12 9h.01"
								></path><path d="M11 12h1v4h1"></path></svg
							><span class="sr-only">More Info</span></button
						>
					</div>
					<div class="flex h-[calc(100%_-_48px)] flex-col justify-between p-6 py-4">
						<div class="flex flex-col">
							<div class="flex flex-wrap items-center justify-between gap-6">
								<div class="text-3xl font-bold">4,682</div>
								<div
									data-slot="chart"
									data-chart="chart-r9p"
									class="[&amp;_.recharts-cartesian-axis-tick_text]:fill-muted-foreground [&amp;_.recharts-cartesian-grid_line[stroke='#ccc']]:stroke-border/50 [&amp;_.recharts-curve.recharts-tooltip-cursor]:stroke-border [&amp;_.recharts-polar-grid_[stroke='#ccc']]:stroke-border [&amp;_.recharts-radial-bar-background-sector]:fill-muted [&amp;_.recharts-rectangle.recharts-tooltip-cursor]:fill-muted [&amp;_.recharts-reference-line_[stroke='#ccc']]:stroke-border [&amp;_.recharts-dot[stroke='#fff']]:stroke-transparent [&amp;_.recharts-layer]:outline-hidden [&amp;_.recharts-sector]:outline-hidden [&amp;_.recharts-sector[stroke='#fff']]:stroke-transparent [&amp;_.recharts-surface]:outline-hidden flex aspect-video w-[70px] justify-center text-xs"
								>
									<style>
										[data-chart='chart-r9p'] {
											--color-month: var(--chart-1);
										}

										.dark [data-chart='chart-r9p'] {
											--color-month: var(--chart-1);
										}
									</style>
									<div
										class="recharts-responsive-container"
										style="width: 100%; height: 100%; min-width: 0px;"
									>
										<div
											class="recharts-wrapper"
											style="position: relative; cursor: default; width: 100%; height: 100%; max-height: 39px; max-width: 70px;"
										>
											<svg
												tabindex="0"
												role="application"
												class="recharts-surface"
												width="70"
												height="39"
												viewBox="0 0 70 39"
												style="width: 100%; height: 100%;"
												><title></title><desc></desc><defs
													><clipPath id="recharts152-clip"
														><rect x="5" y="5" height="29" width="60"></rect></clipPath
													></defs
												><g class="recharts-layer recharts-line"
													><path
														stroke="var(--color-month)"
														stroke-width="1.5"
														width="60"
														height="29"
														fill="none"
														class="recharts-curve recharts-line-curve"
														stroke-dasharray="100.97642517089844px 0px"
														d="M5,15.875L15,6.359L25,12.522L35,27.384L45,15.059L55,33.094L65,14.606"
													></path><g class="recharts-layer"></g></g
												></svg
											>
										</div>
									</div>
								</div>
							</div>
							<p class="text-muted-foreground text-xs">Since Last week</p>
						</div>
						<div class="flex flex-wrap items-center justify-between gap-5">
							<div class="text-sm font-semibold">Details</div>
							<div class="flex items-center gap-1 text-emerald-500 dark:text-emerald-400">
								<p class="text-[13px] leading-none font-medium">15.54%</p>
								<svg
									xmlns="http://www.w3.org/2000/svg"
									width="18"
									height="18"
									viewBox="0 0 24 24"
									fill="currentColor"
									stroke="none"
									class="tabler-icon tabler-icon-caret-up-filled"
									><path
										d="M11.293 7.293a1 1 0 0 1 1.32 -.083l.094 .083l6 6l.083 .094l.054 .077l.054 .096l.017 .036l.027 .067l.032 .108l.01 .053l.01 .06l.004 .057l.002 .059l-.002 .059l-.005 .058l-.009 .06l-.01 .052l-.032 .108l-.027 .067l-.07 .132l-.065 .09l-.073 .081l-.094 .083l-.077 .054l-.096 .054l-.036 .017l-.067 .027l-.108 .032l-.053 .01l-.06 .01l-.057 .004l-.059 .002h-12c-.852 0 -1.297 -.986 -.783 -1.623l.076 -.084l6 -6z"
									></path></svg
								>
							</div>
						</div>
					</div>
				</div>
				<div
					class="bg-card text-card-foreground col-span-3 h-full rounded-xl border shadow-sm lg:col-span-2 xl:col-span-2"
				>
					<div class="flex flex-row items-center justify-between gap-5 space-y-0 p-6 pt-4 pb-2">
						<div class="flex items-center gap-2 truncate text-sm font-medium tracking-tight">
							<svg
								xmlns="http://www.w3.org/2000/svg"
								width="16"
								height="16"
								viewBox="0 0 24 24"
								fill="none"
								stroke="currentColor"
								stroke-width="2"
								stroke-linecap="round"
								stroke-linejoin="round"
								class="tabler-icon tabler-icon-menu-order"
								><path d="M4 10h16"></path><path d="M4 14h16"></path><path d="M9 18l3 3l3 -3"
								></path><path d="M9 6l3 -3l3 3"></path></svg
							>New Orders
						</div>
						<button data-state="closed"
							><svg
								xmlns="http://www.w3.org/2000/svg"
								width="24"
								height="24"
								viewBox="0 0 24 24"
								fill="none"
								stroke="currentColor"
								stroke-width="2"
								stroke-linecap="round"
								stroke-linejoin="round"
								class="tabler-icon tabler-icon-info-circle text-muted-foreground scale-90 stroke-[1.25]"
								><path d="M3 12a9 9 0 1 0 18 0a9 9 0 0 0 -18 0"></path><path d="M12 9h.01"
								></path><path d="M11 12h1v4h1"></path></svg
							><span class="sr-only">More Info</span></button
						>
					</div>
					<div class="flex h-[calc(100%_-_48px)] flex-col justify-between p-6 py-4">
						<div class="flex flex-col">
							<div class="flex flex-wrap items-center justify-between gap-6">
								<div class="text-3xl font-bold">1,226</div>
								<div
									data-slot="chart"
									data-chart="chart-r9r"
									class="[&amp;_.recharts-cartesian-axis-tick_text]:fill-muted-foreground [&amp;_.recharts-cartesian-grid_line[stroke='#ccc']]:stroke-border/50 [&amp;_.recharts-curve.recharts-tooltip-cursor]:stroke-border [&amp;_.recharts-polar-grid_[stroke='#ccc']]:stroke-border [&amp;_.recharts-radial-bar-background-sector]:fill-muted [&amp;_.recharts-rectangle.recharts-tooltip-cursor]:fill-muted [&amp;_.recharts-reference-line_[stroke='#ccc']]:stroke-border [&amp;_.recharts-dot[stroke='#fff']]:stroke-transparent [&amp;_.recharts-layer]:outline-hidden [&amp;_.recharts-sector]:outline-hidden [&amp;_.recharts-sector[stroke='#fff']]:stroke-transparent [&amp;_.recharts-surface]:outline-hidden flex aspect-video w-[70px] justify-center text-xs"
								>
									<style>
										[data-chart='chart-r9r'] {
											--color-month: var(--chart-2);
										}

										.dark [data-chart='chart-r9r'] {
											--color-month: var(--chart-2);
										}
									</style>
									<div
										class="recharts-responsive-container"
										style="width: 100%; height: 100%; min-width: 0px;"
									>
										<div
											class="recharts-wrapper"
											style="position: relative; cursor: default; width: 100%; height: 100%; max-height: 39px; max-width: 70px;"
										>
											<svg
												tabindex="0"
												role="application"
												class="recharts-surface"
												width="70"
												height="39"
												viewBox="0 0 70 39"
												style="width: 100%; height: 100%;"
												><title></title><desc></desc><defs
													><clipPath id="recharts155-clip"
														><rect x="5" y="5" height="29" width="60"></rect></clipPath
													></defs
												><g class="recharts-layer recharts-line"
													><path
														stroke="var(--color-month)"
														stroke-width="1.5"
														width="60"
														height="29"
														fill="none"
														class="recharts-curve recharts-line-curve"
														stroke-dasharray="80.24917602539062px 0px"
														d="M5,17.144L15,6.359L25,12.522L35,27.384L45,15.059L55,14.606L65,14.606"
													></path><g class="recharts-layer"></g></g
												></svg
											>
										</div>
									</div>
								</div>
							</div>
							<p class="text-muted-foreground text-xs">Since Last week</p>
						</div>
						<div class="flex flex-wrap items-center justify-between gap-5">
							<div class="text-sm font-semibold">Details</div>
							<div class="flex items-center gap-1 text-red-500 dark:text-red-400">
								<p class="text-[13px] leading-none font-medium">40.2%</p>
								<svg
									xmlns="http://www.w3.org/2000/svg"
									width="24"
									height="24"
									viewBox="0 0 24 24"
									fill="currentColor"
									stroke="none"
									class="tabler-icon tabler-icon-caret-down-filled"
									><path
										d="M18 9c.852 0 1.297 .986 .783 1.623l-.076 .084l-6 6a1 1 0 0 1 -1.32 .083l-.094 -.083l-6 -6l-.083 -.094l-.054 -.077l-.054 -.096l-.017 -.036l-.027 -.067l-.032 -.108l-.01 -.053l-.01 -.06l-.004 -.057v-.118l.005 -.058l.009 -.06l.01 -.052l.032 -.108l.027 -.067l.07 -.132l.065 -.09l.073 -.081l.094 -.083l.077 -.054l.096 -.054l.036 -.017l.067 -.027l.108 -.032l.053 -.01l.06 -.01l.057 -.004l12.059 -.002z"
									></path></svg
								>
							</div>
						</div>
					</div>
				</div>
				<div
					class="bg-card text-card-foreground col-span-3 h-full rounded-xl border shadow-sm lg:col-span-2 xl:col-span-2"
				>
					<div class="flex flex-row items-center justify-between gap-5 space-y-0 p-6 pt-4 pb-2">
						<div class="flex items-center gap-2 truncate text-sm font-medium tracking-tight">
							<svg
								xmlns="http://www.w3.org/2000/svg"
								width="16"
								height="16"
								viewBox="0 0 24 24"
								fill="none"
								stroke="currentColor"
								stroke-width="2"
								stroke-linecap="round"
								stroke-linejoin="round"
								class="tabler-icon tabler-icon-gift"
								><path
									d="M3 8m0 1a1 1 0 0 1 1 -1h16a1 1 0 0 1 1 1v2a1 1 0 0 1 -1 1h-16a1 1 0 0 1 -1 -1z"
								></path><path d="M12 8l0 13"></path><path
									d="M19 12v7a2 2 0 0 1 -2 2h-10a2 2 0 0 1 -2 -2v-7"
								></path><path
									d="M7.5 8a2.5 2.5 0 0 1 0 -5a4.8 8 0 0 1 4.5 5a4.8 8 0 0 1 4.5 -5a2.5 2.5 0 0 1 0 5"
								></path></svg
							>Avg Order Revenue
						</div>
						<button data-state="closed"
							><svg
								xmlns="http://www.w3.org/2000/svg"
								width="24"
								height="24"
								viewBox="0 0 24 24"
								fill="none"
								stroke="currentColor"
								stroke-width="2"
								stroke-linecap="round"
								stroke-linejoin="round"
								class="tabler-icon tabler-icon-info-circle text-muted-foreground scale-90 stroke-[1.25]"
								><path d="M3 12a9 9 0 1 0 18 0a9 9 0 0 0 -18 0"></path><path d="M12 9h.01"
								></path><path d="M11 12h1v4h1"></path></svg
							><span class="sr-only">More Info</span></button
						>
					</div>
					<div class="flex h-[calc(100%_-_48px)] flex-col justify-between p-6 py-4">
						<div class="flex flex-col">
							<div class="flex flex-wrap items-center justify-between gap-6">
								<div class="text-3xl font-bold">1,080</div>
								<div
									data-slot="chart"
									data-chart="chart-r9t"
									class="[&amp;_.recharts-cartesian-axis-tick_text]:fill-muted-foreground [&amp;_.recharts-cartesian-grid_line[stroke='#ccc']]:stroke-border/50 [&amp;_.recharts-curve.recharts-tooltip-cursor]:stroke-border [&amp;_.recharts-polar-grid_[stroke='#ccc']]:stroke-border [&amp;_.recharts-radial-bar-background-sector]:fill-muted [&amp;_.recharts-rectangle.recharts-tooltip-cursor]:fill-muted [&amp;_.recharts-reference-line_[stroke='#ccc']]:stroke-border [&amp;_.recharts-dot[stroke='#fff']]:stroke-transparent [&amp;_.recharts-layer]:outline-hidden [&amp;_.recharts-sector]:outline-hidden [&amp;_.recharts-sector[stroke='#fff']]:stroke-transparent [&amp;_.recharts-surface]:outline-hidden flex aspect-video w-[70px] justify-center text-xs"
								>
									<style>
										[data-chart='chart-r9t'] {
											--color-month: #6366f1;
										}

										.dark [data-chart='chart-r9t'] {
											--color-month: #6366f1;
										}
									</style>
									<div
										class="recharts-responsive-container"
										style="width: 100%; height: 100%; min-width: 0px;"
									>
										<div
											class="recharts-wrapper"
											style="position: relative; cursor: default; width: 100%; height: 100%; max-height: 39px; max-width: 70px;"
										>
											<svg
												tabindex="0"
												role="application"
												class="recharts-surface"
												width="70"
												height="39"
												viewBox="0 0 70 39"
												style="width: 100%; height: 100%;"
												><title></title><desc></desc><defs
													><clipPath id="recharts158-clip"
														><rect x="5" y="5" height="29" width="60"></rect></clipPath
													></defs
												><g class="recharts-layer recharts-line"
													><path
														stroke="var(--color-month)"
														stroke-width="1.5"
														width="60"
														height="29"
														fill="none"
														class="recharts-curve recharts-line-curve"
														stroke-dasharray="88.83061981201172px 0px"
														d="M5,29.167L15,21.917L25,10.8L35,25.01L45,13.797L55,19.5L65,5"
													></path><g class="recharts-layer"></g></g
												></svg
											>
										</div>
									</div>
								</div>
							</div>
							<p class="text-muted-foreground text-xs">Since Last week</p>
						</div>
						<div class="flex flex-wrap items-center justify-between gap-5">
							<div class="text-sm font-semibold">Details</div>
							<div class="flex items-center gap-1 text-emerald-500 dark:text-emerald-400">
								<p class="text-[13px] leading-none font-medium">10.8%</p>
								<svg
									xmlns="http://www.w3.org/2000/svg"
									width="18"
									height="18"
									viewBox="0 0 24 24"
									fill="currentColor"
									stroke="none"
									class="tabler-icon tabler-icon-caret-up-filled"
									><path
										d="M11.293 7.293a1 1 0 0 1 1.32 -.083l.094 .083l6 6l.083 .094l.054 .077l.054 .096l.017 .036l.027 .067l.032 .108l.01 .053l.01 .06l.004 .057l.002 .059l-.002 .059l-.005 .058l-.009 .06l-.01 .052l-.032 .108l-.027 .067l-.07 .132l-.065 .09l-.073 .081l-.094 .083l-.077 .054l-.096 .054l-.036 .017l-.067 .027l-.108 .032l-.053 .01l-.06 .01l-.057 .004l-.059 .002h-12c-.852 0 -1.297 -.986 -.783 -1.623l.076 -.084l6 -6z"
									></path></svg
								>
							</div>
						</div>
					</div>
				</div>
				<div class="col-span-3">
					<div class="bg-card text-card-foreground h-full rounded-xl border shadow-sm">
						<div class="flex flex-row items-center justify-between space-y-0 p-6 pb-2">
							<div class="text-sm font-normal tracking-tight">Total Revenue</div>
						</div>
						<div class="h-[calc(100%_-_52px)] p-6 pt-0 pb-0">
							<div class="text-2xl font-bold">$15,231.89</div>
							<p class="text-muted-foreground text-xs">+20.1% from last month</p>
							<div
								data-slot="chart"
								data-chart="chart-r9u"
								class="[&amp;_.recharts-cartesian-axis-tick_text]:fill-muted-foreground [&amp;_.recharts-cartesian-grid_line[stroke='#ccc']]:stroke-border/50 [&amp;_.recharts-curve.recharts-tooltip-cursor]:stroke-border [&amp;_.recharts-polar-grid_[stroke='#ccc']]:stroke-border [&amp;_.recharts-radial-bar-background-sector]:fill-muted [&amp;_.recharts-rectangle.recharts-tooltip-cursor]:fill-muted [&amp;_.recharts-reference-line_[stroke='#ccc']]:stroke-border [&amp;_.recharts-dot[stroke='#fff']]:stroke-transparent [&amp;_.recharts-layer]:outline-hidden [&amp;_.recharts-sector]:outline-hidden [&amp;_.recharts-sector[stroke='#fff']]:stroke-transparent [&amp;_.recharts-surface]:outline-hidden flex aspect-video h-[80px] w-full justify-center text-xs"
							>
								<style>
									[data-chart='chart-r9u'] {
										--color-revenue: var(--primary);
										--color-subscription: var(--primary);
									}

									.dark [data-chart='chart-r9u'] {
										--color-revenue: var(--primary);
										--color-subscription: var(--primary);
									}
								</style>
								<div
									class="recharts-responsive-container"
									style="width: 100%; height: 100%; min-width: 0px;"
								>
									<div
										class="recharts-wrapper"
										style="position: relative; cursor: default; width: 100%; height: 100%; max-height: 80px; max-width: 389px;"
									>
										<svg
											class="recharts-surface"
											width="389"
											height="80"
											viewBox="0 0 389 80"
											style="width: 100%; height: 100%;"
											><title></title><desc></desc><defs
												><clipPath id="recharts161-clip"
													><rect x="10" y="5" height="75" width="369"></rect></clipPath
												></defs
											><g class="recharts-layer recharts-line"
												><path
													stroke-width="2"
													stroke="var(--color-revenue)"
													width="369"
													height="75"
													fill="none"
													class="recharts-curve recharts-line-curve"
													d="M10,52.143C27.571,46.779,45.143,41.415,62.714,41.415C80.286,41.415,97.857,52.679,115.429,54.821C133,56.964,150.571,56.964,168.143,58.036C185.714,59.107,203.286,61.25,220.857,61.25C238.429,61.25,256,56.18,273.571,54.286C291.143,52.391,308.714,52.818,326.286,49.882C343.857,46.946,361.429,28.016,379,9.085"
												></path><g class="recharts-layer"></g><g
													class="recharts-layer recharts-line-dots"
													><circle
														r="3"
														stroke-width="2"
														stroke="var(--color-revenue)"
														width="369"
														height="75"
														fill="#fff"
														cx="10"
														cy="52.14285714285714"
														class="recharts-dot recharts-line-dot"
													></circle><circle
														r="3"
														stroke-width="2"
														stroke="var(--color-revenue)"
														width="369"
														height="75"
														fill="#fff"
														cx="62.714285714285715"
														cy="41.41517857142857"
														class="recharts-dot recharts-line-dot"
													></circle><circle
														r="3"
														stroke-width="2"
														stroke="var(--color-revenue)"
														width="369"
														height="75"
														fill="#fff"
														cx="115.42857142857143"
														cy="54.821428571428584"
														class="recharts-dot recharts-line-dot"
													></circle><circle
														r="3"
														stroke-width="2"
														stroke="var(--color-revenue)"
														width="369"
														height="75"
														fill="#fff"
														cx="168.14285714285714"
														cy="58.035714285714285"
														class="recharts-dot recharts-line-dot"
													></circle><circle
														r="3"
														stroke-width="2"
														stroke="var(--color-revenue)"
														width="369"
														height="75"
														fill="#fff"
														cx="220.85714285714286"
														cy="61.25"
														class="recharts-dot recharts-line-dot"
													></circle><circle
														r="3"
														stroke-width="2"
														stroke="var(--color-revenue)"
														width="369"
														height="75"
														fill="#fff"
														cx="273.57142857142856"
														cy="54.285714285714285"
														class="recharts-dot recharts-line-dot"
													></circle><circle
														r="3"
														stroke-width="2"
														stroke="var(--color-revenue)"
														width="369"
														height="75"
														fill="#fff"
														cx="326.2857142857143"
														cy="49.88214285714285"
														class="recharts-dot recharts-line-dot"
													></circle><circle
														r="3"
														stroke-width="2"
														stroke="var(--color-revenue)"
														width="369"
														height="75"
														fill="#fff"
														cx="379"
														cy="9.08482142857143"
														class="recharts-dot recharts-line-dot"
													></circle></g
												></g
											></svg
										>
									</div>
								</div>
							</div>
						</div>
					</div>
				</div>
				<div class="col-span-3 md:col-span-6">
					<div class="bg-card text-card-foreground h-full rounded-xl border shadow-sm">
						<div class="flex flex-col space-y-1.5 p-6">
							<div class="leading-none font-semibold tracking-tight">Sale Activity - Monthly</div>
							<div class="text-muted-foreground text-sm">
								Showing total sales for the last 6 months
							</div>
						</div>
						<div class="h-[calc(100%_-_90px)] p-6 pt-0">
							<div
								class="recharts-responsive-container"
								style="width: 100%; height: 100%; min-width: 0px;"
							>
								<div
									data-slot="chart"
									data-chart="chart-raf"
									class="[&amp;_.recharts-cartesian-axis-tick_text]:fill-muted-foreground [&amp;_.recharts-cartesian-grid_line[stroke='#ccc']]:stroke-border/50 [&amp;_.recharts-curve.recharts-tooltip-cursor]:stroke-border [&amp;_.recharts-polar-grid_[stroke='#ccc']]:stroke-border [&amp;_.recharts-radial-bar-background-sector]:fill-muted [&amp;_.recharts-rectangle.recharts-tooltip-cursor]:fill-muted [&amp;_.recharts-reference-line_[stroke='#ccc']]:stroke-border [&amp;_.recharts-dot[stroke='#fff']]:stroke-transparent [&amp;_.recharts-layer]:outline-hidden [&amp;_.recharts-sector]:outline-hidden [&amp;_.recharts-sector[stroke='#fff']]:stroke-transparent [&amp;_.recharts-surface]:outline-hidden flex aspect-video justify-center text-xs"
									style="height: 100%; width: 100%; max-height: 239px; max-width: 849px;"
								>
									<style>
										[data-chart='chart-raf'] {
											--color-desktop: var(--chart-1);
											--color-mobile: var(--chart-2);
										}

										.dark [data-chart='chart-raf'] {
											--color-desktop: var(--chart-1);
											--color-mobile: var(--chart-2);
										}
									</style>
									<div
										class="recharts-responsive-container"
										style="width: 100%; height: 100%; min-width: 0px;"
									>
										<div
											class="recharts-wrapper"
											style="position: relative; cursor: default; width: 100%; height: 100%; max-height: 239px; max-width: 849px;"
										>
											<svg
												tabindex="0"
												role="application"
												class="recharts-surface"
												width="849"
												height="239"
												viewBox="0 0 849 239"
												style="width: 100%; height: 100%;"
												><title></title><desc></desc><defs
													><clipPath id="recharts169-clip"
														><rect x="12" y="0" height="209" width="825"></rect></clipPath
													></defs
												><g class="recharts-cartesian-grid"
													><g class="recharts-cartesian-grid-horizontal"
														><line
															stroke="#ccc"
															fill="none"
															x="12"
															y="0"
															width="825"
															height="209"
															x1="12"
															y1="209"
															x2="837"
															y2="209"
														></line><line
															stroke="#ccc"
															fill="none"
															x="12"
															y="0"
															width="825"
															height="209"
															x1="12"
															y1="156.75"
															x2="837"
															y2="156.75"
														></line><line
															stroke="#ccc"
															fill="none"
															x="12"
															y="0"
															width="825"
															height="209"
															x1="12"
															y1="104.5"
															x2="837"
															y2="104.5"
														></line><line
															stroke="#ccc"
															fill="none"
															x="12"
															y="0"
															width="825"
															height="209"
															x1="12"
															y1="52.25"
															x2="837"
															y2="52.25"
														></line><line
															stroke="#ccc"
															fill="none"
															x="12"
															y="0"
															width="825"
															height="209"
															x1="12"
															y1="0"
															x2="837"
															y2="0"
														></line></g
													></g
												><g class="recharts-layer recharts-cartesian-axis recharts-xAxis xAxis"
													><g class="recharts-cartesian-axis-ticks"
														><g class="recharts-layer recharts-cartesian-axis-tick"
															><text
																orientation="bottom"
																width="825"
																height="30"
																stroke="none"
																x="12"
																y="223"
																class="recharts-text recharts-cartesian-axis-tick-value"
																text-anchor="middle"
																fill="#666"><tspan x="12" dy="0.71em">Jan</tspan></text
															></g
														><g class="recharts-layer recharts-cartesian-axis-tick"
															><text
																orientation="bottom"
																width="825"
																height="30"
																stroke="none"
																x="177"
																y="223"
																class="recharts-text recharts-cartesian-axis-tick-value"
																text-anchor="middle"
																fill="#666"><tspan x="177" dy="0.71em">Feb</tspan></text
															></g
														><g class="recharts-layer recharts-cartesian-axis-tick"
															><text
																orientation="bottom"
																width="825"
																height="30"
																stroke="none"
																x="342"
																y="223"
																class="recharts-text recharts-cartesian-axis-tick-value"
																text-anchor="middle"
																fill="#666"><tspan x="342" dy="0.71em">Mar</tspan></text
															></g
														><g class="recharts-layer recharts-cartesian-axis-tick"
															><text
																orientation="bottom"
																width="825"
																height="30"
																stroke="none"
																x="507"
																y="223"
																class="recharts-text recharts-cartesian-axis-tick-value"
																text-anchor="middle"
																fill="#666"><tspan x="507" dy="0.71em">Apr</tspan></text
															></g
														><g class="recharts-layer recharts-cartesian-axis-tick"
															><text
																orientation="bottom"
																width="825"
																height="30"
																stroke="none"
																x="672"
																y="223"
																class="recharts-text recharts-cartesian-axis-tick-value"
																text-anchor="middle"
																fill="#666"><tspan x="672" dy="0.71em">May</tspan></text
															></g
														><g class="recharts-layer recharts-cartesian-axis-tick"
															><text
																orientation="bottom"
																width="825"
																height="30"
																stroke="none"
																x="837"
																y="223"
																class="recharts-text recharts-cartesian-axis-tick-value"
																text-anchor="middle"
																fill="#666"><tspan x="837" dy="0.71em">Jun</tspan></text
															></g
														></g
													></g
												><defs
													><linearGradient id="fillDesktop" x1="0" y1="0" x2="0" y2="1"
														><stop offset="5%" stop-color="var(--color-desktop)" stop-opacity="0.8"
														></stop><stop
															offset="95%"
															stop-color="var(--color-desktop)"
															stop-opacity="0.1"
														></stop></linearGradient
													><linearGradient id="fillMobile" x1="0" y1="0" x2="0" y2="1"
														><stop offset="5%" stop-color="var(--color-mobile)" stop-opacity="0.8"
														></stop><stop
															offset="95%"
															stop-color="var(--color-mobile)"
															stop-opacity="0.1"
														></stop></linearGradient
													></defs
												><g class="recharts-layer recharts-area"
													><g class="recharts-layer"
														><path
															fill="url(#fillMobile)"
															fill-opacity="0.4"
															width="825"
															height="209"
															stroke="none"
															class="recharts-curve recharts-area-area"
															d="M12,181.133C67,159.4,122,137.667,177,139.333C232,141,287,166.067,342,167.2C397,168.333,452,145.533,507,142.817C562,140.1,617,157.467,672,163.717C727,169.967,782,165.1,837,160.233L837,209C782,209,727,209,672,209C617,209,562,209,507,209C452,209,397,209,342,209C287,209,232,209,177,209C122,209,67,209,12,209Z"
														></path><path
															fill="none"
															fill-opacity="0.4"
															stroke="var(--color-mobile)"
															width="825"
															height="209"
															class="recharts-curve recharts-area-curve"
															d="M12,181.133C67,159.4,122,137.667,177,139.333C232,141,287,166.067,342,167.2C397,168.333,452,145.533,507,142.817C562,140.1,617,157.467,672,163.717C727,169.967,782,165.1,837,160.233"
														></path></g
													></g
												><g class="recharts-layer recharts-area"
													><g class="recharts-layer"
														><path
															fill="url(#fillDesktop)"
															fill-opacity="0.4"
															width="825"
															height="209"
															stroke="none"
															class="recharts-curve recharts-area-area"
															d="M12,116.343C67,76.514,122,36.686,177,33.092C232,29.498,287,62.139,342,84.645C397,107.151,452,119.522,507,117.388C562,115.254,617,98.616,672,90.915C727,83.214,782,84.452,837,85.69L837,160.233C782,165.1,727,169.967,672,163.717C617,157.467,562,140.1,507,142.817C452,145.533,397,168.333,342,167.2C287,166.067,232,141,177,139.333C122,137.667,67,159.4,12,181.133Z"
														></path><path
															fill="none"
															fill-opacity="0.4"
															stroke="var(--color-desktop)"
															width="825"
															height="209"
															class="recharts-curve recharts-area-curve"
															d="M12,116.343C67,76.514,122,36.686,177,33.092C232,29.498,287,62.139,342,84.645C397,107.151,452,119.522,507,117.388C562,115.254,617,98.616,672,90.915C727,83.214,782,84.452,837,85.69"
														></path></g
													></g
												></svg
											>
											<div
												tabindex="-1"
												class="recharts-tooltip-wrapper recharts-tooltip-wrapper-right recharts-tooltip-wrapper-bottom"
												style="visibility: hidden; pointer-events: none; position: absolute; top: 0px; left: 0px; transform: translate(186.4px, 76px);"
											></div>
										</div>
									</div>
								</div>
							</div>
						</div>
					</div>
				</div>
				<div class="col-span-3 md:col-span-6 lg:col-span-3">
					<div class="bg-card text-card-foreground h-full rounded-xl border shadow-sm">
						<div class="flex flex-row items-center justify-between space-y-0 p-6 pb-2">
							<div class="text-sm font-normal tracking-tight">Subscriptions</div>
						</div>
						<div class="h-full p-6 pt-0">
							<div class="text-2xl font-bold">+2350</div>
							<p class="text-muted-foreground text-xs">+180.1% from last month</p>
							<div
								data-slot="chart"
								data-chart="chart-r9v"
								class="[&amp;_.recharts-cartesian-axis-tick_text]:fill-muted-foreground [&amp;_.recharts-cartesian-grid_line[stroke='#ccc']]:stroke-border/50 [&amp;_.recharts-curve.recharts-tooltip-cursor]:stroke-border [&amp;_.recharts-polar-grid_[stroke='#ccc']]:stroke-border [&amp;_.recharts-radial-bar-background-sector]:fill-muted [&amp;_.recharts-rectangle.recharts-tooltip-cursor]:fill-muted [&amp;_.recharts-reference-line_[stroke='#ccc']]:stroke-border [&amp;_.recharts-dot[stroke='#fff']]:stroke-transparent [&amp;_.recharts-layer]:outline-hidden [&amp;_.recharts-sector]:outline-hidden [&amp;_.recharts-sector[stroke='#fff']]:stroke-transparent [&amp;_.recharts-surface]:outline-hidden mt-6 flex aspect-video h-[calc(100%_-_120px)] max-h-[205px] w-full justify-center text-xs"
							>
								<style>
									[data-chart='chart-r9v'] {
										--color-revenue: var(--chart-2);
										--color-subscription: var(--chart-1);
									}

									.dark [data-chart='chart-r9v'] {
										--color-revenue: var(--chart-2);
										--color-subscription: var(--chart-1);
									}
								</style>
								<div
									class="recharts-responsive-container"
									style="width: 100%; height: 100%; min-width: 0px;"
								>
									<div
										class="recharts-wrapper"
										style="position: relative; cursor: default; width: 100%; height: 100%; max-height: 205px; max-width: 389px;"
									>
										<svg
											class="recharts-surface"
											width="389"
											height="205"
											viewBox="0 0 389 205"
											style="width: 100%; height: 100%;"
											><title></title><desc></desc><defs
												><clipPath id="recharts164-clip"
													><rect x="5" y="5" height="195" width="379"></rect></clipPath
												></defs
											><g class="recharts-layer recharts-bar"
												><g class="recharts-layer recharts-bar-rectangles"
													><g class="recharts-layer"
														><g class="recharts-layer recharts-bar-rectangle"
															><path
																x="9.7375"
																y="122"
																width="16"
																height="78"
																radius="4"
																fill="var(--color-subscription)"
																class="recharts-rectangle"
																d="M 9.7375,126
            A 4,4,0,0,1,13.7375,122
            L 21.7375,122
            A 4,4,0,0,1,25.7375,126
            L 25.7375,196
            A 4,4,0,0,1,21.7375,200
            L 13.7375,200
            A 4,4,0,0,1,9.7375,196 Z"
															></path></g
														><g class="recharts-layer recharts-bar-rectangle"
															><path
																x="57.1125"
																y="102.5"
																width="16"
																height="97.5"
																radius="4"
																fill="var(--color-subscription)"
																class="recharts-rectangle"
																d="M 57.1125,106.5
            A 4,4,0,0,1,61.1125,102.5
            L 69.1125,102.5
            A 4,4,0,0,1,73.1125,106.5
            L 73.1125,196
            A 4,4,0,0,1,69.1125,200
            L 61.1125,200
            A 4,4,0,0,1,57.1125,196 Z"
															></path></g
														><g class="recharts-layer recharts-bar-rectangle"
															><path
																x="104.4875"
																y="135"
																width="16"
																height="65"
																radius="4"
																fill="var(--color-subscription)"
																class="recharts-rectangle"
																d="M 104.4875,139
            A 4,4,0,0,1,108.4875,135
            L 116.4875,135
            A 4,4,0,0,1,120.4875,139
            L 120.4875,196
            A 4,4,0,0,1,116.4875,200
            L 108.4875,200
            A 4,4,0,0,1,104.4875,196 Z"
															></path></g
														><g class="recharts-layer recharts-bar-rectangle"
															><path
																x="151.8625"
																y="109.64999999999999"
																width="16"
																height="90.35000000000001"
																radius="4"
																fill="var(--color-subscription)"
																class="recharts-rectangle"
																d="M 151.8625,113.64999999999999
            A 4,4,0,0,1,155.8625,109.64999999999999
            L 163.8625,109.64999999999999
            A 4,4,0,0,1,167.8625,113.64999999999999
            L 167.8625,196
            A 4,4,0,0,1,163.8625,200
            L 155.8625,200
            A 4,4,0,0,1,151.8625,196 Z"
															></path></g
														><g class="recharts-layer recharts-bar-rectangle"
															><path
																x="199.2375"
																y="138.575"
																width="16"
																height="61.42500000000001"
																radius="4"
																fill="var(--color-subscription)"
																class="recharts-rectangle"
																d="M 199.2375,142.575
            A 4,4,0,0,1,203.2375,138.575
            L 211.2375,138.575
            A 4,4,0,0,1,215.2375,142.575
            L 215.2375,196
            A 4,4,0,0,1,211.2375,200
            L 203.2375,200
            A 4,4,0,0,1,199.2375,196 Z"
															></path></g
														><g class="recharts-layer recharts-bar-rectangle"
															><path
																x="246.6125"
																y="122.325"
																width="16"
																height="77.675"
																radius="4"
																fill="var(--color-subscription)"
																class="recharts-rectangle"
																d="M 246.6125,126.325
            A 4,4,0,0,1,250.6125,122.325
            L 258.6125,122.325
            A 4,4,0,0,1,262.6125,126.325
            L 262.6125,196
            A 4,4,0,0,1,258.6125,200
            L 250.6125,200
            A 4,4,0,0,1,246.6125,196 Z"
															></path></g
														><g class="recharts-layer recharts-bar-rectangle"
															><path
																x="293.9875"
																y="109.64999999999999"
																width="16"
																height="90.35000000000001"
																radius="4"
																fill="var(--color-subscription)"
																class="recharts-rectangle"
																d="M 293.9875,113.64999999999999
            A 4,4,0,0,1,297.9875,109.64999999999999
            L 305.9875,109.64999999999999
            A 4,4,0,0,1,309.9875,113.64999999999999
            L 309.9875,196
            A 4,4,0,0,1,305.9875,200
            L 297.9875,200
            A 4,4,0,0,1,293.9875,196 Z"
															></path></g
														><g class="recharts-layer recharts-bar-rectangle"
															><path
																x="341.3625"
																y="138.575"
																width="16"
																height="61.42500000000001"
																radius="4"
																fill="var(--color-subscription)"
																class="recharts-rectangle"
																d="M 341.3625,142.575
            A 4,4,0,0,1,345.3625,138.575
            L 353.3625,138.575
            A 4,4,0,0,1,357.3625,142.575
            L 357.3625,196
            A 4,4,0,0,1,353.3625,200
            L 345.3625,200
            A 4,4,0,0,1,341.3625,196 Z"
															></path></g
														></g
													></g
												><g class="recharts-layer"></g></g
											><g class="recharts-layer recharts-bar"
												><g class="recharts-layer recharts-bar-rectangles"
													><g class="recharts-layer"
														><g class="recharts-layer recharts-bar-rectangle"
															><path
																x="29.7375"
																y="166.20000000000002"
																width="16"
																height="33.79999999999998"
																radius="4"
																fill="var(--color-revenue)"
																class="recharts-rectangle"
																d="M 29.7375,170.20000000000002
            A 4,4,0,0,1,33.7375,166.20000000000002
            L 41.7375,166.20000000000002
            A 4,4,0,0,1,45.7375,170.20000000000002
            L 45.7375,196
            A 4,4,0,0,1,41.7375,200
            L 33.7375,200
            A 4,4,0,0,1,29.7375,196 Z"
															></path></g
														><g class="recharts-layer recharts-bar-rectangle"
															><path
																x="77.1125"
																y="153.2"
																width="16"
																height="46.80000000000001"
																radius="4"
																fill="var(--color-revenue)"
																class="recharts-rectangle"
																d="M 77.1125,157.2
            A 4,4,0,0,1,81.1125,153.2
            L 89.1125,153.2
            A 4,4,0,0,1,93.1125,157.2
            L 93.1125,196
            A 4,4,0,0,1,89.1125,200
            L 81.1125,200
            A 4,4,0,0,1,77.1125,196 Z"
															></path></g
														><g class="recharts-layer recharts-bar-rectangle"
															><path
																x="124.4875"
																y="57.00000000000001"
																width="16"
																height="143"
																radius="4"
																fill="var(--color-revenue)"
																class="recharts-rectangle"
																d="M 124.4875,61.00000000000001
            A 4,4,0,0,1,128.4875,57.00000000000001
            L 136.4875,57.00000000000001
            A 4,4,0,0,1,140.4875,61.00000000000001
            L 140.4875,196
            A 4,4,0,0,1,136.4875,200
            L 128.4875,200
            A 4,4,0,0,1,124.4875,196 Z"
															></path></g
														><g class="recharts-layer recharts-bar-rectangle"
															><path
																x="171.8625"
																y="96"
																width="16"
																height="104"
																radius="4"
																fill="var(--color-revenue)"
																class="recharts-rectangle"
																d="M 171.8625,100
            A 4,4,0,0,1,175.8625,96
            L 183.8625,96
            A 4,4,0,0,1,187.8625,100
            L 187.8625,196
            A 4,4,0,0,1,183.8625,200
            L 175.8625,200
            A 4,4,0,0,1,171.8625,196 Z"
															></path></g
														><g class="recharts-layer recharts-bar-rectangle"
															><path
																x="219.2375"
																y="167.50000000000003"
																width="16"
																height="32.49999999999997"
																radius="4"
																fill="var(--color-revenue)"
																class="recharts-rectangle"
																d="M 219.2375,171.50000000000003
            A 4,4,0,0,1,223.2375,167.50000000000003
            L 231.2375,167.50000000000003
            A 4,4,0,0,1,235.2375,171.50000000000003
            L 235.2375,196
            A 4,4,0,0,1,231.2375,200
            L 223.2375,200
            A 4,4,0,0,1,219.2375,196 Z"
															></path></g
														><g class="recharts-layer recharts-bar-rectangle"
															><path
																x="266.6125"
																y="17.999999999999996"
																width="16"
																height="182"
																radius="4"
																fill="var(--color-revenue)"
																class="recharts-rectangle"
																d="M 266.6125,21.999999999999996
            A 4,4,0,0,1,270.6125,17.999999999999996
            L 278.6125,17.999999999999996
            A 4,4,0,0,1,282.6125,21.999999999999996
            L 282.6125,196
            A 4,4,0,0,1,278.6125,200
            L 270.6125,200
            A 4,4,0,0,1,266.6125,196 Z"
															></path></g
														><g class="recharts-layer recharts-bar-rectangle"
															><path
																x="313.9875"
																y="163.6"
																width="16"
																height="36.400000000000006"
																radius="4"
																fill="var(--color-revenue)"
																class="recharts-rectangle"
																d="M 313.9875,167.6
            A 4,4,0,0,1,317.9875,163.6
            L 325.9875,163.6
            A 4,4,0,0,1,329.9875,167.6
            L 329.9875,196
            A 4,4,0,0,1,325.9875,200
            L 317.9875,200
            A 4,4,0,0,1,313.9875,196 Z"
															></path></g
														><g class="recharts-layer recharts-bar-rectangle"
															><path
																x="361.3625"
																y="114.20000000000002"
																width="16"
																height="85.79999999999998"
																radius="4"
																fill="var(--color-revenue)"
																class="recharts-rectangle"
																d="M 361.3625,118.20000000000002
            A 4,4,0,0,1,365.3625,114.20000000000002
            L 373.3625,114.20000000000002
            A 4,4,0,0,1,377.3625,118.20000000000002
            L 377.3625,196
            A 4,4,0,0,1,373.3625,200
            L 365.3625,200
            A 4,4,0,0,1,361.3625,196 Z"
															></path></g
														></g
													></g
												><g class="recharts-layer"></g></g
											></svg
										>
									</div>
								</div>
							</div>
						</div>
					</div>
				</div>
				<div class="col-span-3 md:col-span-6 lg:col-span-5 xl:col-span-6">
					<div class="bg-card text-card-foreground h-full rounded-xl border shadow-sm">
						<div class="flex flex-col space-y-1.5 p-6">
							<div class="text-xl font-semibold tracking-tight">Payments</div>
							<div class="text-muted-foreground text-sm">Manage your payments.</div>
						</div>
						<div class="h-[calc(100%_-_102px)] p-6 pt-0">
							<div class="mb-4 flex items-center gap-4">
								<input
									class="border-input file:text-foreground placeholder:text-muted-foreground focus-visible:ring-ring flex h-9 w-full max-w-sm rounded-md border bg-transparent px-3 py-1 text-base shadow-xs transition-colors file:border-0 file:bg-transparent file:text-sm file:font-medium focus-visible:ring-1 focus-visible:outline-hidden disabled:cursor-not-allowed disabled:opacity-50 md:text-sm"
									placeholder="Filter emails..."
									value=""
								/><button
									class="focus-visible:ring-ring [&amp;_svg]:pointer-events-none [&amp;_svg]:size-4 [&amp;_svg]:shrink-0 border-input bg-background hover:bg-accent hover:text-accent-foreground ml-auto inline-flex h-9 items-center justify-center gap-2 rounded-md border px-4 py-2 text-sm font-medium whitespace-nowrap shadow-xs transition-colors focus-visible:ring-1 focus-visible:outline-hidden disabled:pointer-events-none disabled:opacity-50"
									type="button"
									id="radix-:ra0:"
									aria-haspopup="menu"
									aria-expanded="false"
									data-state="closed"
									>Columns <svg
										xmlns="http://www.w3.org/2000/svg"
										width="24"
										height="24"
										viewBox="0 0 24 24"
										fill="none"
										stroke="currentColor"
										stroke-width="2"
										stroke-linecap="round"
										stroke-linejoin="round"
										class="lucide lucide-chevron-down"><path d="m6 9 6 6 6-6"></path></svg
									></button
								>
							</div>
							<div class="h-[calc(100%_-_52px)] rounded-md border">
								<div class="relative w-full overflow-auto">
									<table class="w-full caption-bottom text-sm">
										<thead class="[&amp;_tr]:border-b"
											><tr
												class="group/row hover:bg-muted/50 data-[state=selected]:bg-muted border-b transition-colors"
												><th
													class="text-muted-foreground [&amp;:has([role=checkbox])]:pr-0 [&amp;&gt;[role=checkbox]]:translate-y-[2px] [&amp;:has([role=checkbox])]:pl-3 h-10 px-2 text-left align-middle font-medium"
													><button
														type="button"
														role="checkbox"
														aria-checked="false"
														data-state="unchecked"
														value="on"
														class="peer border-primary focus-visible:ring-ring data-[state=checked]:bg-primary data-[state=checked]:text-primary-foreground h-4 w-4 shrink-0 rounded-sm border shadow-sm focus-visible:ring-1 focus-visible:outline-hidden disabled:cursor-not-allowed disabled:opacity-50"
														aria-label="Select all"
													></button></th
												><th
													class="text-muted-foreground [&amp;:has([role=checkbox])]:pr-0 [&amp;&gt;[role=checkbox]]:translate-y-[2px] [&amp;:has([role=checkbox])]:pl-3 h-10 px-2 text-left align-middle font-medium"
													>Status</th
												><th
													class="text-muted-foreground [&amp;:has([role=checkbox])]:pr-0 [&amp;&gt;[role=checkbox]]:translate-y-[2px] [&amp;:has([role=checkbox])]:pl-3 h-10 px-2 text-left align-middle font-medium"
													><button
														class="focus-visible:ring-ring [&amp;_svg]:pointer-events-none [&amp;_svg]:size-4 [&amp;_svg]:shrink-0 hover:bg-accent hover:text-accent-foreground inline-flex h-9 items-center justify-center gap-2 rounded-md px-4 py-2 text-sm font-medium whitespace-nowrap transition-colors focus-visible:ring-1 focus-visible:outline-hidden disabled:pointer-events-none disabled:opacity-50"
														>Email<svg
															xmlns="http://www.w3.org/2000/svg"
															width="24"
															height="24"
															viewBox="0 0 24 24"
															fill="none"
															stroke="currentColor"
															stroke-width="2"
															stroke-linecap="round"
															stroke-linejoin="round"
															class="lucide lucide-arrow-up-down"
															><path d="m21 16-4 4-4-4"></path><path d="M17 20V4"></path><path
																d="m3 8 4-4 4 4"
															></path><path d="M7 4v16"></path></svg
														></button
													></th
												><th
													class="text-muted-foreground [&amp;:has([role=checkbox])]:pr-0 [&amp;&gt;[role=checkbox]]:translate-y-[2px] [&amp;:has([role=checkbox])]:pl-3 h-10 px-2 text-left align-middle font-medium"
													><div class="text-right">Amount</div></th
												><th
													class="text-muted-foreground [&amp;:has([role=checkbox])]:pr-0 [&amp;&gt;[role=checkbox]]:translate-y-[2px] [&amp;:has([role=checkbox])]:pl-3 h-10 px-2 text-left align-middle font-medium"
												></th></tr
											></thead
										><tbody class="[&amp;_tr:last-child]:border-0"
											><tr
												class="group/row hover:bg-muted/50 data-[state=selected]:bg-muted border-b transition-colors"
												data-state="false"
												><td
													class="[&amp;:has([role=checkbox])]:pr-0 [&amp;&gt;[role=checkbox]]:translate-y-[2px] [&amp;:has([role=checkbox])]:pl-3 p-2 align-middle"
													><button
														type="button"
														role="checkbox"
														aria-checked="false"
														data-state="unchecked"
														value="on"
														class="peer border-primary focus-visible:ring-ring data-[state=checked]:bg-primary data-[state=checked]:text-primary-foreground h-4 w-4 shrink-0 rounded-sm border shadow-sm focus-visible:ring-1 focus-visible:outline-hidden disabled:cursor-not-allowed disabled:opacity-50"
														aria-label="Select row"
													></button></td
												><td
													class="[&amp;:has([role=checkbox])]:pr-0 [&amp;&gt;[role=checkbox]]:translate-y-[2px] [&amp;:has([role=checkbox])]:pl-3 p-2 align-middle"
													><div class="capitalize">success</div></td
												><td
													class="[&amp;:has([role=checkbox])]:pr-0 [&amp;&gt;[role=checkbox]]:translate-y-[2px] [&amp;:has([role=checkbox])]:pl-3 p-2 align-middle"
													><div class="lowercase">ken99@yahoo.com</div></td
												><td
													class="[&amp;:has([role=checkbox])]:pr-0 [&amp;&gt;[role=checkbox]]:translate-y-[2px] [&amp;:has([role=checkbox])]:pl-3 p-2 align-middle"
													><div class="text-right font-medium">$316.00</div></td
												><td
													class="[&amp;:has([role=checkbox])]:pr-0 [&amp;&gt;[role=checkbox]]:translate-y-[2px] [&amp;:has([role=checkbox])]:pl-3 p-2 align-middle"
													><button
														class="focus-visible:ring-ring [&amp;_svg]:pointer-events-none [&amp;_svg]:size-4 [&amp;_svg]:shrink-0 hover:bg-accent hover:text-accent-foreground inline-flex h-8 w-8 items-center justify-center gap-2 rounded-md p-0 text-sm font-medium whitespace-nowrap transition-colors focus-visible:ring-1 focus-visible:outline-hidden disabled:pointer-events-none disabled:opacity-50"
														type="button"
														id="radix-:ra2:"
														aria-haspopup="menu"
														aria-expanded="false"
														data-state="closed"
														><span class="sr-only">Open menu</span><svg
															xmlns="http://www.w3.org/2000/svg"
															width="24"
															height="24"
															viewBox="0 0 24 24"
															fill="none"
															stroke="currentColor"
															stroke-width="2"
															stroke-linecap="round"
															stroke-linejoin="round"
															class="lucide lucide-ellipsis"
															><circle cx="12" cy="12" r="1"></circle><circle cx="19" cy="12" r="1"
															></circle><circle cx="5" cy="12" r="1"></circle></svg
														></button
													></td
												></tr
											><tr
												class="group/row hover:bg-muted/50 data-[state=selected]:bg-muted border-b transition-colors"
												data-state="false"
												><td
													class="[&amp;:has([role=checkbox])]:pr-0 [&amp;&gt;[role=checkbox]]:translate-y-[2px] [&amp;:has([role=checkbox])]:pl-3 p-2 align-middle"
													><button
														type="button"
														role="checkbox"
														aria-checked="false"
														data-state="unchecked"
														value="on"
														class="peer border-primary focus-visible:ring-ring data-[state=checked]:bg-primary data-[state=checked]:text-primary-foreground h-4 w-4 shrink-0 rounded-sm border shadow-sm focus-visible:ring-1 focus-visible:outline-hidden disabled:cursor-not-allowed disabled:opacity-50"
														aria-label="Select row"
													></button></td
												><td
													class="[&amp;:has([role=checkbox])]:pr-0 [&amp;&gt;[role=checkbox]]:translate-y-[2px] [&amp;:has([role=checkbox])]:pl-3 p-2 align-middle"
													><div class="capitalize">success</div></td
												><td
													class="[&amp;:has([role=checkbox])]:pr-0 [&amp;&gt;[role=checkbox]]:translate-y-[2px] [&amp;:has([role=checkbox])]:pl-3 p-2 align-middle"
													><div class="lowercase">Abe45@gmail.com</div></td
												><td
													class="[&amp;:has([role=checkbox])]:pr-0 [&amp;&gt;[role=checkbox]]:translate-y-[2px] [&amp;:has([role=checkbox])]:pl-3 p-2 align-middle"
													><div class="text-right font-medium">$242.00</div></td
												><td
													class="[&amp;:has([role=checkbox])]:pr-0 [&amp;&gt;[role=checkbox]]:translate-y-[2px] [&amp;:has([role=checkbox])]:pl-3 p-2 align-middle"
													><button
														class="focus-visible:ring-ring [&amp;_svg]:pointer-events-none [&amp;_svg]:size-4 [&amp;_svg]:shrink-0 hover:bg-accent hover:text-accent-foreground inline-flex h-8 w-8 items-center justify-center gap-2 rounded-md p-0 text-sm font-medium whitespace-nowrap transition-colors focus-visible:ring-1 focus-visible:outline-hidden disabled:pointer-events-none disabled:opacity-50"
														type="button"
														id="radix-:ra4:"
														aria-haspopup="menu"
														aria-expanded="false"
														data-state="closed"
														><span class="sr-only">Open menu</span><svg
															xmlns="http://www.w3.org/2000/svg"
															width="24"
															height="24"
															viewBox="0 0 24 24"
															fill="none"
															stroke="currentColor"
															stroke-width="2"
															stroke-linecap="round"
															stroke-linejoin="round"
															class="lucide lucide-ellipsis"
															><circle cx="12" cy="12" r="1"></circle><circle cx="19" cy="12" r="1"
															></circle><circle cx="5" cy="12" r="1"></circle></svg
														></button
													></td
												></tr
											><tr
												class="group/row hover:bg-muted/50 data-[state=selected]:bg-muted border-b transition-colors"
												data-state="false"
												><td
													class="[&amp;:has([role=checkbox])]:pr-0 [&amp;&gt;[role=checkbox]]:translate-y-[2px] [&amp;:has([role=checkbox])]:pl-3 p-2 align-middle"
													><button
														type="button"
														role="checkbox"
														aria-checked="false"
														data-state="unchecked"
														value="on"
														class="peer border-primary focus-visible:ring-ring data-[state=checked]:bg-primary data-[state=checked]:text-primary-foreground h-4 w-4 shrink-0 rounded-sm border shadow-sm focus-visible:ring-1 focus-visible:outline-hidden disabled:cursor-not-allowed disabled:opacity-50"
														aria-label="Select row"
													></button></td
												><td
													class="[&amp;:has([role=checkbox])]:pr-0 [&amp;&gt;[role=checkbox]]:translate-y-[2px] [&amp;:has([role=checkbox])]:pl-3 p-2 align-middle"
													><div class="capitalize">processing</div></td
												><td
													class="[&amp;:has([role=checkbox])]:pr-0 [&amp;&gt;[role=checkbox]]:translate-y-[2px] [&amp;:has([role=checkbox])]:pl-3 p-2 align-middle"
													><div class="lowercase">Monserrat44@gmail.com</div></td
												><td
													class="[&amp;:has([role=checkbox])]:pr-0 [&amp;&gt;[role=checkbox]]:translate-y-[2px] [&amp;:has([role=checkbox])]:pl-3 p-2 align-middle"
													><div class="text-right font-medium">$837.00</div></td
												><td
													class="[&amp;:has([role=checkbox])]:pr-0 [&amp;&gt;[role=checkbox]]:translate-y-[2px] [&amp;:has([role=checkbox])]:pl-3 p-2 align-middle"
													><button
														class="focus-visible:ring-ring [&amp;_svg]:pointer-events-none [&amp;_svg]:size-4 [&amp;_svg]:shrink-0 hover:bg-accent hover:text-accent-foreground inline-flex h-8 w-8 items-center justify-center gap-2 rounded-md p-0 text-sm font-medium whitespace-nowrap transition-colors focus-visible:ring-1 focus-visible:outline-hidden disabled:pointer-events-none disabled:opacity-50"
														type="button"
														id="radix-:ra6:"
														aria-haspopup="menu"
														aria-expanded="false"
														data-state="closed"
														><span class="sr-only">Open menu</span><svg
															xmlns="http://www.w3.org/2000/svg"
															width="24"
															height="24"
															viewBox="0 0 24 24"
															fill="none"
															stroke="currentColor"
															stroke-width="2"
															stroke-linecap="round"
															stroke-linejoin="round"
															class="lucide lucide-ellipsis"
															><circle cx="12" cy="12" r="1"></circle><circle cx="19" cy="12" r="1"
															></circle><circle cx="5" cy="12" r="1"></circle></svg
														></button
													></td
												></tr
											><tr
												class="group/row hover:bg-muted/50 data-[state=selected]:bg-muted border-b transition-colors"
												data-state="false"
												><td
													class="[&amp;:has([role=checkbox])]:pr-0 [&amp;&gt;[role=checkbox]]:translate-y-[2px] [&amp;:has([role=checkbox])]:pl-3 p-2 align-middle"
													><button
														type="button"
														role="checkbox"
														aria-checked="false"
														data-state="unchecked"
														value="on"
														class="peer border-primary focus-visible:ring-ring data-[state=checked]:bg-primary data-[state=checked]:text-primary-foreground h-4 w-4 shrink-0 rounded-sm border shadow-sm focus-visible:ring-1 focus-visible:outline-hidden disabled:cursor-not-allowed disabled:opacity-50"
														aria-label="Select row"
													></button></td
												><td
													class="[&amp;:has([role=checkbox])]:pr-0 [&amp;&gt;[role=checkbox]]:translate-y-[2px] [&amp;:has([role=checkbox])]:pl-3 p-2 align-middle"
													><div class="capitalize">failed</div></td
												><td
													class="[&amp;:has([role=checkbox])]:pr-0 [&amp;&gt;[role=checkbox]]:translate-y-[2px] [&amp;:has([role=checkbox])]:pl-3 p-2 align-middle"
													><div class="lowercase">carmella@hotmail.com</div></td
												><td
													class="[&amp;:has([role=checkbox])]:pr-0 [&amp;&gt;[role=checkbox]]:translate-y-[2px] [&amp;:has([role=checkbox])]:pl-3 p-2 align-middle"
													><div class="text-right font-medium">$721.00</div></td
												><td
													class="[&amp;:has([role=checkbox])]:pr-0 [&amp;&gt;[role=checkbox]]:translate-y-[2px] [&amp;:has([role=checkbox])]:pl-3 p-2 align-middle"
													><button
														class="focus-visible:ring-ring [&amp;_svg]:pointer-events-none [&amp;_svg]:size-4 [&amp;_svg]:shrink-0 hover:bg-accent hover:text-accent-foreground inline-flex h-8 w-8 items-center justify-center gap-2 rounded-md p-0 text-sm font-medium whitespace-nowrap transition-colors focus-visible:ring-1 focus-visible:outline-hidden disabled:pointer-events-none disabled:opacity-50"
														type="button"
														id="radix-:ra8:"
														aria-haspopup="menu"
														aria-expanded="false"
														data-state="closed"
														><span class="sr-only">Open menu</span><svg
															xmlns="http://www.w3.org/2000/svg"
															width="24"
															height="24"
															viewBox="0 0 24 24"
															fill="none"
															stroke="currentColor"
															stroke-width="2"
															stroke-linecap="round"
															stroke-linejoin="round"
															class="lucide lucide-ellipsis"
															><circle cx="12" cy="12" r="1"></circle><circle cx="19" cy="12" r="1"
															></circle><circle cx="5" cy="12" r="1"></circle></svg
														></button
													></td
												></tr
											></tbody
										>
									</table>
								</div>
							</div>
						</div>
					</div>
				</div>
				<div class="col-span-3 md:col-span-6 lg:col-span-4 xl:col-span-3">
					<div class="bg-card text-card-foreground h-full rounded-xl border shadow-sm">
						<div class="flex flex-col space-y-1.5 p-6">
							<div class="leading-none font-semibold tracking-tight">Team Members</div>
							<div class="text-muted-foreground truncate text-sm">
								Invite your team members to collaborate.
							</div>
						</div>
						<div class="grid gap-6 p-6 pt-0">
							<div class="flex flex-wrap items-center justify-between gap-2">
								<div class="flex items-center space-x-4">
									<span class="relative flex h-8 w-8 shrink-0 overflow-hidden rounded-full"
										><img
											class="aspect-square h-full w-full"
											alt="Image"
											src="/avatars/avatar-1.png"
										/></span
									>
									<div>
										<p class="text-sm leading-none font-medium">Dale Komen</p>
										<p class="text-muted-foreground text-sm">dale@example.com</p>
									</div>
								</div>
								<button
									class="focus-visible:ring-ring [&amp;_svg]:pointer-events-none [&amp;_svg]:size-4 [&amp;_svg]:shrink-0 border-input bg-background hover:bg-accent hover:text-accent-foreground inline-flex h-7 items-center justify-center gap-2 rounded-md border px-3 py-0 text-xs font-medium whitespace-nowrap shadow-xs transition-colors focus-visible:ring-1 focus-visible:outline-hidden disabled:pointer-events-none disabled:opacity-50"
									type="button"
									aria-haspopup="dialog"
									aria-expanded="false"
									aria-controls="radix-:raa:"
									data-state="closed"
									>Member <svg
										xmlns="http://www.w3.org/2000/svg"
										width="24"
										height="24"
										viewBox="0 0 24 24"
										fill="none"
										stroke="currentColor"
										stroke-width="2"
										stroke-linecap="round"
										stroke-linejoin="round"
										class="lucide lucide-chevron-down text-muted-foreground"
										><path d="m6 9 6 6 6-6"></path></svg
									></button
								>
							</div>
							<div class="flex flex-wrap items-center justify-between gap-2">
								<div class="flex items-center space-x-4">
									<span class="relative flex h-8 w-8 shrink-0 overflow-hidden rounded-full"
										><img
											class="aspect-square h-full w-full"
											alt="Image"
											src="/avatars/avatar-5.png"
										/></span
									>
									<div>
										<p class="text-sm leading-none font-medium">Sofia Davis</p>
										<p class="text-muted-foreground text-sm">m@example.com</p>
									</div>
								</div>
								<button
									class="focus-visible:ring-ring [&amp;_svg]:pointer-events-none [&amp;_svg]:size-4 [&amp;_svg]:shrink-0 border-input bg-background hover:bg-accent hover:text-accent-foreground inline-flex h-7 items-center justify-center gap-2 rounded-md border px-3 py-0 text-xs font-medium whitespace-nowrap shadow-xs transition-colors focus-visible:ring-1 focus-visible:outline-hidden disabled:pointer-events-none disabled:opacity-50"
									type="button"
									aria-haspopup="dialog"
									aria-expanded="false"
									aria-controls="radix-:rab:"
									data-state="closed"
									>Owner <svg
										xmlns="http://www.w3.org/2000/svg"
										width="24"
										height="24"
										viewBox="0 0 24 24"
										fill="none"
										stroke="currentColor"
										stroke-width="2"
										stroke-linecap="round"
										stroke-linejoin="round"
										class="lucide lucide-chevron-down text-muted-foreground"
										><path d="m6 9 6 6 6-6"></path></svg
									></button
								>
							</div>
							<div class="flex flex-wrap items-center justify-between gap-2">
								<div class="flex items-center space-x-4">
									<span class="relative flex h-8 w-8 shrink-0 overflow-hidden rounded-full"
										><img
											class="aspect-square h-full w-full"
											alt="Image"
											src="/avatars/avatar-4.png"
										/></span
									>
									<div>
										<p class="text-sm leading-none font-medium">Jackson Lee</p>
										<p class="text-muted-foreground text-sm">p@example.com</p>
									</div>
								</div>
								<button
									class="focus-visible:ring-ring [&amp;_svg]:pointer-events-none [&amp;_svg]:size-4 [&amp;_svg]:shrink-0 border-input bg-background hover:bg-accent hover:text-accent-foreground inline-flex h-7 items-center justify-center gap-2 rounded-md border px-3 py-0 text-xs font-medium whitespace-nowrap shadow-xs transition-colors focus-visible:ring-1 focus-visible:outline-hidden disabled:pointer-events-none disabled:opacity-50"
									type="button"
									aria-haspopup="dialog"
									aria-expanded="false"
									aria-controls="radix-:rac:"
									data-state="closed"
									>Member <svg
										xmlns="http://www.w3.org/2000/svg"
										width="24"
										height="24"
										viewBox="0 0 24 24"
										fill="none"
										stroke="currentColor"
										stroke-width="2"
										stroke-linecap="round"
										stroke-linejoin="round"
										class="lucide lucide-chevron-down text-muted-foreground"
										><path d="m6 9 6 6 6-6"></path></svg
									></button
								>
							</div>
							<div class="flex flex-wrap items-center justify-between gap-2">
								<div class="flex items-center space-x-4">
									<span class="relative flex h-8 w-8 shrink-0 overflow-hidden rounded-full"
										><img
											class="aspect-square h-full w-full"
											alt="Image"
											src="/avatars/avatar-3.png"
										/></span
									>
									<div>
										<p class="text-sm leading-none font-medium">Isabella Nguyen</p>
										<p class="text-muted-foreground text-sm">i@example.com</p>
									</div>
								</div>
								<button
									class="focus-visible:ring-ring [&amp;_svg]:pointer-events-none [&amp;_svg]:size-4 [&amp;_svg]:shrink-0 border-input bg-background hover:bg-accent hover:text-accent-foreground inline-flex h-7 items-center justify-center gap-2 rounded-md border px-3 py-0 text-xs font-medium whitespace-nowrap shadow-xs transition-colors focus-visible:ring-1 focus-visible:outline-hidden disabled:pointer-events-none disabled:opacity-50"
									type="button"
									aria-haspopup="dialog"
									aria-expanded="false"
									aria-controls="radix-:rad:"
									data-state="closed"
									>Member <svg
										xmlns="http://www.w3.org/2000/svg"
										width="24"
										height="24"
										viewBox="0 0 24 24"
										fill="none"
										stroke="currentColor"
										stroke-width="2"
										stroke-linecap="round"
										stroke-linejoin="round"
										class="lucide lucide-chevron-down text-muted-foreground"
										><path d="m6 9 6 6 6-6"></path></svg
									></button
								>
							</div>
							<div class="flex flex-wrap items-center justify-between gap-2">
								<div class="flex items-center space-x-4">
									<span class="relative flex h-8 w-8 shrink-0 overflow-hidden rounded-full"
										><img
											class="aspect-square h-full w-full"
											alt="Image"
											src="/avatars/avatar-2.png"
										/></span
									>
									<div>
										<p class="text-sm leading-none font-medium">Hugan Romex</p>
										<p class="text-muted-foreground text-sm">kai@example.com</p>
									</div>
								</div>
								<button
									class="focus-visible:ring-ring [&amp;_svg]:pointer-events-none [&amp;_svg]:size-4 [&amp;_svg]:shrink-0 border-input bg-background hover:bg-accent hover:text-accent-foreground inline-flex h-7 items-center justify-center gap-2 rounded-md border px-3 py-0 text-xs font-medium whitespace-nowrap shadow-xs transition-colors focus-visible:ring-1 focus-visible:outline-hidden disabled:pointer-events-none disabled:opacity-50"
									type="button"
									aria-haspopup="dialog"
									aria-expanded="false"
									aria-controls="radix-:rae:"
									data-state="closed"
									>Member <svg
										xmlns="http://www.w3.org/2000/svg"
										width="24"
										height="24"
										viewBox="0 0 24 24"
										fill="none"
										stroke="currentColor"
										stroke-width="2"
										stroke-linecap="round"
										stroke-linejoin="round"
										class="lucide lucide-chevron-down text-muted-foreground"
										><path d="m6 9 6 6 6-6"></path></svg
									></button
								>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
		<div
			data-state="inactive"
			data-orientation="vertical"
			role="tabpanel"
			aria-labelledby="radix-:r9j:-trigger-analytics"
			hidden=""
			id="radix-:r9j:-content-analytics"
			tabindex="0"
			class="ring-offset-background focus-visible:ring-ring mt-2 space-y-4 focus-visible:ring-2 focus-visible:ring-offset-2 focus-visible:outline-hidden"
		></div>
	</div>
</div>
