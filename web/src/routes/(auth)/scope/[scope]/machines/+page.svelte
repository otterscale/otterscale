<script lang="ts">
	import { ReloadManager, Reloader } from '$lib/components/custom/reloader';
	import { page } from '$app/state';
	import { env } from '$env/dynamic/public';
	import { EnvironmentService } from '$lib/api/environment/v1/environment_pb';
	import { MachineService, type Machine } from '$lib/api/machine/v1/machine_pb';
	import { Dashboard } from '$lib/components/machines/dashboard';
	import { buttonVariants } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import * as Chart from '$lib/components/ui/chart';
	import * as Tabs from '$lib/components/ui/tabs';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { formatCapacity } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';
	import { getLocale } from '$lib/paraglide/runtime';
	import { dynamicPaths } from '$lib/path';
	import { breadcrumb } from '$lib/stores';
	import { cn } from '$lib/utils';
	import { timestampDate } from '@bufbuild/protobuf/wkt';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { scaleUtc } from 'd3-scale';
	import { curveLinear, curveNatural } from 'd3-shape';
	import { Area, AreaChart, LinearGradient, LineChart, PieChart, Text } from 'layerchart';
	import { PrometheusDriver, SampleValue } from 'prometheus-query';
	import { getContext, onDestroy, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	// Set breadcrumb navigation
	breadcrumb.set({ parents: [], current: dynamicPaths.machines(page.params.scope) });

	let totalCPUCores = $state(0);
	let totalMemoryBytes = $state(0);
	let totalStorageBytes = $state(0);
	let totalDisks = $state(0);
	let totalNodes = $state(0);

	let cpuUsages = $state([] as SampleValue[]);
	const cpuUsagesConfiguration = {
		usage: { label: 'value', color: 'var(--chart-1)' }
	} satisfies Chart.ChartConfig;
	const cpuUsagesTrend = $derived(
		cpuUsages.length > 0
			? (cpuUsages[cpuUsages.length - 1].value - cpuUsages[cpuUsages.length - 2].value) /
					cpuUsages[cpuUsages.length - 2].value
			: 0
	);
	let memoryUsages = $state([] as SampleValue[]);
	const memoryUsagesConfiguration = {
		usage: { label: 'value', color: 'var(--chart-1)' }
	} satisfies Chart.ChartConfig;
	const memoryUsagesTrend = $derived(
		memoryUsages.length > 0
			? (memoryUsages[memoryUsages.length - 1].value -
					memoryUsages[memoryUsages.length - 2].value) /
					memoryUsages[memoryUsages.length - 2].value
			: 0
	);
	let storageUsages = $state([] as SampleValue[]);
	const storageUsagesConfiguration = {
		usage: { label: 'value', color: 'var(--chart-1)' }
	} satisfies Chart.ChartConfig;
	const storageUsagesTrend = $derived(
		storageUsages.length > 0
			? (storageUsages[storageUsages.length - 1].value -
					storageUsages[storageUsages.length - 2].value) /
					storageUsages[storageUsages.length - 2].value
			: 0
	);

	const nodeProportionsConfiguration = {
		nodes: { label: 'Nodes' },
		physical: { label: 'Physical', color: 'var(--chart-1)' },
		virtual: { label: 'Virtual', color: 'var(--chart-2)' }
	} satisfies Chart.ChartConfig;
	let nodeProportions: {
		node: string;
		nodes: number;
		color: string;
	}[] = $state([]);

	const systemLoadConfiguration = {
		one: { label: '1 min', color: 'var(--chart-1)' },
		five: { label: '5 min', color: 'var(--chart-2)' }
	} satisfies Chart.ChartConfig;
	let ones = $state([] as SampleValue[]);
	let fives = $state({} as SampleValue[]);
	const systemLoads = $derived(
		ones.map((sample, index) => ({
			time: sample.time,
			one: sample.value,
			five: fives[index]?.value ?? 0
		}))
	);

	const nodesConfiguration = {
		node: { label: 'Node', color: 'var(--chart-1)' }
	} satisfies Chart.ChartConfig;
	let nodes: {
		date: Date;
		node: number;
	}[] = $state([]);

	const transport: Transport = getContext('transport');
	const machineClient = createClient(MachineService, transport);
	const environmentService = createClient(EnvironmentService, transport);
	const machinesStore = writable<Machine[]>([]);
	let prometheusDriver: PrometheusDriver | null = $state(null);

	function fetch() {
		environmentService.getPrometheus({}).then((response) => {
			prometheusDriver = new PrometheusDriver({
				endpoint: `${env.PUBLIC_API_URL}/prometheus`,
				baseURL: response.baseUrl
			});
			if (prometheusDriver) {
				prometheusDriver
					.rangeQuery(
						`1 - (sum(irate(node_cpu_seconds_total{mode="idle"}[2m])) / sum(irate(node_cpu_seconds_total[2m])))`,
						Date.now() - 10 * 60 * 1000,
						Date.now(),
						2 * 60
					)
					.then((response) => {
						cpuUsages = response.result[0]?.values;
					});
				prometheusDriver
					.rangeQuery(
						`sum(node_memory_MemTotal_bytes - node_memory_MemFree_bytes - (node_memory_Cached_bytes + node_memory_Buffers_bytes + node_memory_SReclaimable_bytes)) / sum(node_memory_MemTotal_bytes)`,
						Date.now() - 10 * 60 * 1000,
						Date.now(),
						2 * 60
					)
					.then((response) => {
						memoryUsages = response.result[0]?.values;
					});
				prometheusDriver
					.rangeQuery(
						`1 - sum(node_filesystem_avail_bytes{mountpoint="/"}) / sum(node_filesystem_size_bytes{mountpoint="/"})`,
						Date.now() - 10 * 60 * 1000,
						Date.now(),
						2 * 60
					)
					.then((response) => {
						storageUsages = response.result[0]?.values;
					});
				prometheusDriver
					.rangeQuery(`sum(node_load1)`, Date.now() - 24 * 60 * 60 * 1000, Date.now(), 2 * 60)
					.then((response) => {
						ones = response.result[0]?.values;
					});
				prometheusDriver
					.rangeQuery(`sum(node_load5)`, Date.now() - 24 * 60 * 60 * 1000, Date.now(), 2 * 60)
					.then((response) => {
						fives = response.result[0]?.values;
					});
			}
		});

		machineClient.listMachines({}).then((response) => {
			machinesStore.set(response.machines);
		});
	}

	const reloadManager = new ReloadManager(fetch);

	let mounted = $state(false);
	onMount(async () => {
		try {
			await environmentService.getPrometheus({}).then((response) => {
				prometheusDriver = new PrometheusDriver({
					endpoint: `${env.PUBLIC_API_URL}/prometheus`,
					baseURL: response.baseUrl
				});
				if (prometheusDriver) {
					prometheusDriver
						.rangeQuery(
							`1 - (sum(irate(node_cpu_seconds_total{mode="idle"}[2m])) / sum(irate(node_cpu_seconds_total[2m])))`,
							Date.now() - 10 * 60 * 1000,
							Date.now(),
							2 * 60
						)
						.then((response) => {
							cpuUsages = response.result[0]?.values;
						});
					prometheusDriver
						.rangeQuery(
							`sum(node_memory_MemTotal_bytes - node_memory_MemFree_bytes - (node_memory_Cached_bytes + node_memory_Buffers_bytes + node_memory_SReclaimable_bytes)) / sum(node_memory_MemTotal_bytes)`,
							Date.now() - 10 * 60 * 1000,
							Date.now(),
							2 * 60
						)
						.then((response) => {
							memoryUsages = response.result[0]?.values;
						});
					prometheusDriver
						.rangeQuery(
							`1 - sum(node_filesystem_avail_bytes{mountpoint="/"}) / sum(node_filesystem_size_bytes{mountpoint="/"})`,
							Date.now() - 10 * 60 * 1000,
							Date.now(),
							2 * 60
						)
						.then((response) => {
							storageUsages = response.result[0]?.values;
						});
					prometheusDriver
						.rangeQuery(`sum(node_load1)`, Date.now() - 24 * 60 * 60 * 1000, Date.now(), 2 * 60)
						.then((response) => {
							ones = response.result[0]?.values;
						});
					prometheusDriver
						.rangeQuery(`sum(node_load5)`, Date.now() - 24 * 60 * 60 * 1000, Date.now(), 2 * 60)
						.then((response) => {
							fives = response.result[0]?.values;
						});
				}
			});

			await machineClient.listMachines({}).then((response) => {
				machinesStore.set(response.machines);
			});

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

			nodeProportions = [
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
			nodes = months.map((month) => ({
				date: new Date(month + '-01'),
				node: monthlyCounts[month] || 0
			}));
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
		mounted = true;

		reloadManager.start();
	});
	onDestroy(() => {
		reloadManager.stop();
	});
</script>

<div class="mx-auto grid w-full gap-6">
	<div class="grid gap-1">
		<h1 class="text-2xl font-bold tracking-tight md:text-3xl">{m.dashboard()}</h1>
		<p class="text-muted-foreground">
			{m.machine_dashboard_description()}
		</p>
	</div>

	<Tabs.Root value="overview">
		<div class="flex justify-between gap-2">
			<Tabs.List>
				<Tabs.Trigger value="overview">{m.overview()}</Tabs.Trigger>
				<Tabs.Trigger value="analytics" disabled>{m.analytics()}</Tabs.Trigger>
			</Tabs.List>
			<Reloader {reloadManager} />
		</div>
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
									<p>{m.machine_dashboard_total_cpu_tooltip()}</p>
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
					<Chart.Container config={cpuUsagesConfiguration} class="h-full w-20">
						<LineChart
							data={cpuUsages}
							x="time"
							xScale={scaleUtc()}
							axis={false}
							series={[
								{
									key: 'value',
									label: 'usage',
									color: cpuUsagesConfiguration.usage.color
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
								<Chart.Tooltip hideLabel>
									{#snippet formatter({ item, name, value })}
										<div
											style="--color-bg: {item.color}"
											class="aspect-square h-full w-fit shrink-0 border-(--color-border) bg-(--color-bg)"
										></div>
										<div
											class="flex flex-1 shrink-0 items-center justify-between text-xs leading-none"
										>
											<div class="grid gap-1.5">
												<span class="text-muted-foreground">{name}</span>
											</div>
											<p class="font-mono">{(Number(value) * 100).toFixed(2)} %</p>
										</div>
									{/snippet}
								</Chart.Tooltip>
							{/snippet}
						</LineChart>
					</Chart.Container>
				</Card.Content>
				<Card.Footer
					class={cn(
						'flex flex-wrap items-center justify-end text-sm leading-none font-medium',
						cpuUsagesTrend >= 0
							? 'text-emerald-500 dark:text-emerald-400'
							: 'text-rose-500 dark:text-rose-400'
					)}
				>
					{Math.abs(cpuUsagesTrend).toFixed(2)} %
					{#if cpuUsagesTrend >= 0}
						<Icon icon="ph:caret-up" />
					{:else}
						<Icon icon="ph:caret-down" />
					{/if}
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
									<p>{m.machine_dashboard_total_memory_tooltip()}</p>
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
						<Chart.Container config={memoryUsagesConfiguration} class="h-full w-20">
							<LineChart
								data={memoryUsages}
								x="time"
								xScale={scaleUtc()}
								axis={false}
								series={[
									{
										key: 'value',
										label: 'usage',
										color: memoryUsagesConfiguration.usage.color
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
									<Chart.Tooltip hideLabel>
										{#snippet formatter({ item, name, value })}
											<div
												style="--color-bg: {item.color}"
												class="aspect-square h-full w-fit shrink-0 border-(--color-border) bg-(--color-bg)"
											></div>
											<div
												class="flex flex-1 shrink-0 items-center justify-between text-xs leading-none"
											>
												<div class="grid gap-1.5">
													<span class="text-muted-foreground">{name}</span>
												</div>
												<p class="font-mono">{(Number(value) * 100).toFixed(2)} %</p>
											</div>
										{/snippet}
									</Chart.Tooltip>
								{/snippet}
							</LineChart>
						</Chart.Container>
					</div>
					<p class="text-muted-foreground text-sm lowercase">{m.total_memory()}</p>
				</Card.Content>
				<Card.Footer
					class={cn(
						'flex flex-wrap items-center justify-end text-sm leading-none font-medium',
						memoryUsagesTrend >= 0
							? 'text-emerald-500 dark:text-emerald-400'
							: 'text-red-500 dark:text-red-400'
					)}
				>
					{Math.abs(memoryUsagesTrend).toFixed(2)} %
					{#if memoryUsagesTrend >= 0}
						<Icon icon="ph:caret-up" />
					{:else}
						<Icon icon="ph:caret-down" />
					{/if}
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
									<p>{m.machine_dashboard_total_storage_tooltip()}</p>
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
						<Chart.Container config={storageUsagesConfiguration} class="h-full w-20">
							<LineChart
								data={storageUsages}
								x="time"
								xScale={scaleUtc()}
								axis={false}
								series={[
									{
										key: 'value',
										label: 'usage',
										color: storageUsagesConfiguration.usage.color
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
									<Chart.Tooltip hideLabel>
										{#snippet formatter({ item, name, value })}
											<div
												style="--color-bg: {item.color}"
												class="aspect-square h-full w-fit shrink-0 border-(--color-border) bg-(--color-bg)"
											></div>
											<div
												class="flex flex-1 shrink-0 items-center justify-between text-xs leading-none"
											>
												<div class="grid gap-1.5">
													<span class="text-muted-foreground">{name}</span>
												</div>
												<p class="font-mono">{(Number(value) * 100).toFixed(2)} %</p>
											</div>
										{/snippet}
									</Chart.Tooltip>
								{/snippet}
							</LineChart>
						</Chart.Container>
					</div>
					<p class="text-muted-foreground text-sm uppercase">{m.over_n_disks({ totalDisks })}</p>
				</Card.Content>
				<Card.Footer
					class={cn(
						'flex flex-wrap items-center justify-end text-sm leading-none font-medium',
						storageUsagesTrend >= 0
							? 'text-emerald-500 dark:text-emerald-400'
							: 'text-red-500 dark:text-red-400'
					)}
				>
					{Math.abs(storageUsagesTrend).toFixed(2)} %
					{#if storageUsagesTrend >= 0}
						<Icon icon="ph:caret-up" />
					{:else}
						<Icon icon="ph:caret-down" />
					{/if}
				</Card.Footer>
			</Card.Root>

			<Card.Root class="col-span-3 gap-2">
				<Card.Header>
					<Card.Title class="flex flex-wrap items-center justify-between gap-6">
						<div
							class="flex flex-col items-start gap-0.5 truncate text-sm font-medium tracking-tight"
						>
							<p class="text-muted-foreground text-xs uppercase">{m.over_the_past_6_months()}</p>
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
									<p>{m.machine_dashboard_nodes_tooltip()}</p>
								</Tooltip.Content>
							</Tooltip.Root>
						</Tooltip.Provider>
					</Card.Title>
				</Card.Header>
				<Card.Content>
					<Chart.Container config={nodesConfiguration} class="h-[80px] w-full px-2 pt-2">
						<LineChart
							points={{ r: 4 }}
							data={nodes}
							x="date"
							xScale={scaleUtc()}
							axis="x"
							series={[
								{
									key: 'node',
									label: 'Node',
									color: nodesConfiguration.node.color
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
					<Card.Title>{m.system_load()}</Card.Title>
					<Card.Description>
						{m.machine_dashboard_system_loads_tooltip()}
					</Card.Description>
				</Card.Header>
				<Card.Content>
					<Chart.Container config={systemLoadConfiguration} class="h-[200px] w-full">
						<AreaChart
							data={systemLoads}
							x="time"
							xScale={scaleUtc()}
							yPadding={[0, 25]}
							series={[
								{
									key: 'one',
									label: systemLoadConfiguration.one.label,
									color: systemLoadConfiguration.one.color
								},
								{
									key: 'five',
									label: systemLoadConfiguration.five.label,
									color: systemLoadConfiguration.five.color
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
											year: 'numeric',
											month: 'short',
											day: 'numeric',
											hour: 'numeric',
											minute: 'numeric'
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
						<div class="flex items-center gap-1 truncate text-sm font-medium tracking-tight">
							<Icon icon="ph:cube" class="size-4.5" />
							{m.node_distribution()}
						</div>
					</Card.Title>
					<Card.Description class="text-xs">
						{m.node_distribution_description()}
					</Card.Description>
				</Card.Header>
				<Card.Content class="flex-1">
					<Chart.Container
						config={nodeProportionsConfiguration}
						class="mx-auto aspect-square max-h-[250px]"
					>
						<PieChart
							data={nodeProportions}
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
		</Tabs.Content>
	</Tabs.Root>
</div>
