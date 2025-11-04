<script lang="ts">
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { scaleUtc } from 'd3-scale';
	import { curveLinear } from 'd3-shape';
	import { LineChart } from 'layerchart';
	import { PrometheusDriver, SampleValue } from 'prometheus-query';
	import { getContext, onDestroy, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import { page } from '$app/state';
	import { MachineService, type Machine } from '$lib/api/machine/v1/machine_pb';
	import type { Scope } from '$lib/api/scope/v1/scope_pb';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import { buttonVariants } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import * as Chart from '$lib/components/ui/chart';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { formatCapacity } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';
	import { cn } from '$lib/utils';

	let {
		prometheusDriver,
		scope,
		isReloading = $bindable(),
	}: { prometheusDriver: PrometheusDriver; scope: Scope; isReloading: boolean } = $props();

	const transport: Transport = getContext('transport');
	const machineClient = createClient(MachineService, transport);

	const machines = writable<Machine[]>([]);
	const scopeMachines = $derived(
		$machines.filter((m) => m.workloadAnnotations['juju-machine-id']?.startsWith(page.params.scope!)),
	);
	const totalMemoryBytes = $derived(scopeMachines.reduce((sum, m) => sum + Number(m.memoryMb ?? 0), 0) * 1024 * 1024);

	let memoryUsages = $state([] as SampleValue[]);
	const memoryUsagesTrend = $derived(
		memoryUsages.length > 1 && memoryUsages[memoryUsages.length - 2].value !== 0
			? (memoryUsages[memoryUsages.length - 1].value - memoryUsages[memoryUsages.length - 2].value) /
					memoryUsages[memoryUsages.length - 2].value
			: 0,
	);

	const memoryUsagesConfiguration = {
		usage: { label: 'value', color: 'var(--chart-1)' },
	} satisfies Chart.ChartConfig;

	async function fetch() {
		prometheusDriver
			.rangeQuery(
				`sum(node_memory_MemTotal_bytes{juju_model_uuid="${scope.uuid}"} - node_memory_MemFree_bytes{juju_model_uuid="${scope.uuid}"} - (node_memory_Cached_bytes{juju_model_uuid="${scope.uuid}"} + node_memory_Buffers_bytes{juju_model_uuid="${scope.uuid}"} + node_memory_SReclaimable_bytes{juju_model_uuid="${scope.uuid}"})) / sum(node_memory_MemTotal_bytes{juju_model_uuid="${scope.uuid}"})`,
				Date.now() - 10 * 60 * 1000,
				Date.now(),
				2 * 60,
			)
			.then((response) => {
				memoryUsages = response.result[0]?.values ?? [];
			});

		machineClient.listMachines({}).then((response) => {
			machines.set(response.machines);
		});
	}

	const reloadManager = new ReloadManager(fetch);

	let isLoading = $state(true);
	onMount(async () => {
		await fetch();
		isLoading = false;
	});
	onDestroy(() => {
		reloadManager.stop();
	});

	$effect(() => {
		if (isReloading) {
			reloadManager.restart();
		} else {
			reloadManager.stop();
		}
	});
</script>

{#if isLoading}
	Loading
{:else}
	<Card.Root class="h-full gap-2">
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
								color: memoryUsagesConfiguration.usage.color,
							},
						]}
						props={{
							spline: { curve: curveLinear, motion: 'tween', strokeWidth: 2 },
							xAxis: {
								format: (v: Date) => v.toLocaleDateString('en-US', { month: 'short' }),
							},
							highlight: { points: { r: 4 } },
						}}
					>
						{#snippet tooltip()}
							<Chart.Tooltip hideLabel>
								{#snippet formatter({ item, name, value })}
									<div
										style="--color-bg: {item.color}"
										class="aspect-square h-full w-fit shrink-0 border-(--color-border) bg-(--color-bg)"
									></div>
									<div class="flex flex-1 shrink-0 items-center justify-between text-xs leading-none">
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
				memoryUsagesTrend >= 0 ? 'text-emerald-500 dark:text-emerald-400' : 'text-red-500 dark:text-red-400',
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
{/if}
