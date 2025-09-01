<script lang="ts">
	import { page } from '$app/state';
	import { env } from '$env/dynamic/public';
	import { EnvironmentService } from '$lib/api/environment/v1/environment_pb';
	import { MachineService, type Machine } from '$lib/api/machine/v1/machine_pb';
	import { Reloader, ReloadManager } from '$lib/components/custom/reloader';
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
	import CPU from './cpu.svelte';
	import Memory from './memory.svelte';
	import Storage from './storage.svelte';
	import Nodes from './nodes.svelte';

	let { prometheusDriver, isReloading = $bindable() }: { prometheusDriver: PrometheusDriver; isReloading: boolean } =
		$props();

	let totalCPUCores = $state(0);
	let totalMemoryBytes = $state(0);
	let totalStorageBytes = $state(0);
	let totalDisks = $state(0);
	let totalNodes = $state(0);

	let cpuUsages = $state([] as SampleValue[]);
	const cpuUsagesConfiguration = {
		usage: { label: 'value', color: 'var(--chart-1)' },
	} satisfies Chart.ChartConfig;
	const cpuUsagesTrend = $derived(
		cpuUsages.length > 0
			? (cpuUsages[cpuUsages.length - 1].value - cpuUsages[cpuUsages.length - 2].value) /
					cpuUsages[cpuUsages.length - 2].value
			: 0,
	);
	let memoryUsages = $state([] as SampleValue[]);
	const memoryUsagesConfiguration = {
		usage: { label: 'value', color: 'var(--chart-1)' },
	} satisfies Chart.ChartConfig;
	const memoryUsagesTrend = $derived(
		memoryUsages.length > 0
			? (memoryUsages[memoryUsages.length - 1].value - memoryUsages[memoryUsages.length - 2].value) /
					memoryUsages[memoryUsages.length - 2].value
			: 0,
	);
	let storageUsages = $state([] as SampleValue[]);
	const storageUsagesConfiguration = {
		usage: { label: 'value', color: 'var(--chart-1)' },
	} satisfies Chart.ChartConfig;
	const storageUsagesTrend = $derived(
		storageUsages.length > 0
			? (storageUsages[storageUsages.length - 1].value - storageUsages[storageUsages.length - 2].value) /
					storageUsages[storageUsages.length - 2].value
			: 0,
	);

	const nodeProportionsConfiguration = {
		nodes: { label: 'Nodes' },
		physical: { label: 'Physical', color: 'var(--chart-1)' },
		virtual: { label: 'Virtual', color: 'var(--chart-2)' },
	} satisfies Chart.ChartConfig;
	let nodeProportions: {
		node: string;
		nodes: number;
		color: string;
	}[] = $state([]);

	const systemLoadConfiguration = {
		one: { label: '1 min', color: 'var(--chart-1)' },
		five: { label: '5 min', color: 'var(--chart-2)' },
	} satisfies Chart.ChartConfig;
	let ones = $state([] as SampleValue[]);
	let fives = $state({} as SampleValue[]);
	const systemLoads = $derived(
		ones.map((sample, index) => ({
			time: sample.time,
			one: sample.value,
			five: fives[index]?.value ?? 0,
		})),
	);

	const nodesConfiguration = {
		node: { label: 'Node', color: 'var(--chart-1)' },
	} satisfies Chart.ChartConfig;
	let nodes: {
		date: Date;
		node: number;
	}[] = $state([]);

	const transport: Transport = getContext('transport');
	const machineClient = createClient(MachineService, transport);
	const environmentService = createClient(EnvironmentService, transport);
	const machinesStore = writable<Machine[]>([]);

	async function fetch() {
		prometheusDriver
			.rangeQuery(
				`1 - (sum(irate(node_cpu_seconds_total{mode="idle"}[2m])) / sum(irate(node_cpu_seconds_total[2m])))`,
				Date.now() - 10 * 60 * 1000,
				Date.now(),
				2 * 60,
			)
			.then((response) => {
				cpuUsages = response.result[0]?.values;
			});
		prometheusDriver
			.rangeQuery(
				`sum(node_memory_MemTotal_bytes - node_memory_MemFree_bytes - (node_memory_Cached_bytes + node_memory_Buffers_bytes + node_memory_SReclaimable_bytes)) / sum(node_memory_MemTotal_bytes)`,
				Date.now() - 10 * 60 * 1000,
				Date.now(),
				2 * 60,
			)
			.then((response) => {
				memoryUsages = response.result[0]?.values;
			});
		prometheusDriver
			.rangeQuery(
				`1 - sum(node_filesystem_avail_bytes{mountpoint="/"}) / sum(node_filesystem_size_bytes{mountpoint="/"})`,
				Date.now() - 10 * 60 * 1000,
				Date.now(),
				2 * 60,
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

		machineClient.listMachines({}).then((response) => {
			machinesStore.set(response.machines);
		});
	}

	const reloadManager = new ReloadManager(fetch);

	let mounted = $state(false);
	onMount(async () => {
		try {
			await fetch();

			const scopeMachines = $machinesStore.filter((m) =>
				m.workloadAnnotations['juju-machine-id']?.startsWith(page.params.scope!),
			);
			totalNodes = scopeMachines.length;
			// totalCPUCores = scopeMachines.reduce((sum, m) => sum + Number(m.cpuCount ?? 0), 0);
			totalMemoryBytes = scopeMachines.reduce((sum, m) => sum + Number(m.memoryMb ?? 0), 0) * 1024 * 1024;

			const blockDevices = scopeMachines.flatMap((m) => m.blockDevices).filter((d) => !d.bootDisk);
			totalDisks = blockDevices.length;
			totalStorageBytes = blockDevices.reduce((sum, m) => sum + Number(m.storageMb ?? 0), 0) * 1024 * 1024;

			const virtualNodes = scopeMachines.filter((m) => m.tags.includes('virtual')).length;
			const physicalNodes = scopeMachines.length - virtualNodes;

			nodeProportions = [
				{ node: 'physical', nodes: physicalNodes, color: 'var(--color-physical)' },
				{ node: 'virtual', nodes: virtualNodes, color: 'var(--color-virtual)' },
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
				node: monthlyCounts[month] || 0,
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

<div class="grid auto-rows-auto grid-cols-3 gap-5 pt-4 md:grid-cols-6 lg:grid-cols-9">
	<div class="col-span-2">
		<CPU {prometheusDriver} {isReloading} />
	</div>

	<div class="col-span-2">
		<Memory {prometheusDriver} {isReloading} />
	</div>

	<div class="col-span-2">
		<Storage {prometheusDriver} {isReloading} />
	</div>

	<div class="col-span-3">
		<Nodes {isReloading} />
	</div>

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
							color: systemLoadConfiguration.one.color,
						},
						{
							key: 'five',
							label: systemLoadConfiguration.five.label,
							color: systemLoadConfiguration.five.color,
						},
					]}
					seriesLayout="stack"
					props={{
						area: {
							curve: curveNatural,
							'fill-opacity': 0.4,
							line: { class: 'stroke-1' },
							motion: 'tween',
						},
						xAxis: {
							format: (v: Date) => v.toLocaleDateString('en-US', { month: 'short' }),
						},
						yAxis: { format: () => '' },
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
									minute: 'numeric',
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
			<Chart.Container config={nodeProportionsConfiguration} class="mx-auto aspect-square max-h-[250px]">
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
</div>
