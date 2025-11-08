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
	import { type Machine, MachineService } from '$lib/api/machine/v1/machine_pb';
	import type { Scope } from '$lib/api/scope/v1/scope_pb';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import { buttonVariants } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import * as Chart from '$lib/components/ui/chart';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { m } from '$lib/paraglide/messages';
	import { cn } from '$lib/utils';

	let {
		prometheusDriver,
		scope,
		isReloading = $bindable()
	}: { prometheusDriver: PrometheusDriver; scope: Scope; isReloading: boolean } = $props();

	const transport: Transport = getContext('transport');
	const machineClient = createClient(MachineService, transport);

	const machines = writable<Machine[]>([]);
	const scopeMachines = $derived(
		$machines.filter((m) =>
			m.workloadAnnotations['juju-machine-id']?.startsWith(page.params.scope!)
		)
	);
	const totalGPUs = $derived(
		scopeMachines.reduce((a, machine) => a + Number(machine.gpuDevices.length ?? 0), 0)
	);
	let allocatedGPUs = $state([] as SampleValue[]);
	const trend = $derived(
		allocatedGPUs.length > 1 && allocatedGPUs[allocatedGPUs.length - 2].value !== 0
			? (allocatedGPUs[allocatedGPUs.length - 1].value -
					allocatedGPUs[allocatedGPUs.length - 2].value) /
					allocatedGPUs[allocatedGPUs.length - 2].value
			: 0
	);

	const configuration = {
		amounts: { label: 'GPUs', color: 'var(--chart-1)' }
	} satisfies Chart.ChartConfig;

	async function fetch() {
		prometheusDriver
			.rangeQuery(
				`
				count(sum by (node, deviceuuid) (vGPUPodsDeviceAllocated{juju_model_uuid="${scope.uuid}"}) > bool 0)
				`,
				Date.now() - 24 * 60 * 60 * 1000,
				Date.now(),
				60 * 60
			)
			.then((response) => {
				allocatedGPUs = response.result && response.result[0] ? response.result[0].values : [];
			});

		machineClient.listMachines({}).then((response) => {
			machines.set(response.machines);
		});
	}

	const reloadManager = new ReloadManager(fetch);

	let isLoading = $state(true);
	onMount(async () => {
		try {
			await fetch();
			isLoading = false;
		} catch (error) {
			console.error(`Fail to fetch data in scope ${scope}:`, error);
		}
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
					<Icon icon="ph:graphics-card" class="size-4.5" />
					{m.gpu()}
				</div>
				<Tooltip.Provider>
					<Tooltip.Root>
						<Tooltip.Trigger class={buttonVariants({ variant: 'ghost', size: 'icon' })}>
							<Icon icon="ph:info" class="size-5 text-muted-foreground" />
						</Tooltip.Trigger>
						<Tooltip.Content>
							<p>{m.machine_dashboard_gpu_tooltip()}</p>
						</Tooltip.Content>
					</Tooltip.Root>
				</Tooltip.Provider>
			</Card.Title>
		</Card.Header>
		<Card.Content class="flex flex-wrap items-center justify-between gap-6">
			<div class="flex flex-col gap-0.5">
				<div class="text-3xl font-bold">{totalGPUs}</div>
				<p class="text-sm text-muted-foreground">{m.pieces()}</p>
			</div>
			<Chart.Container config={configuration} class="h-full w-20">
				<LineChart
					data={allocatedGPUs}
					x="time"
					xScale={scaleUtc()}
					axis={false}
					series={[
						{
							key: 'value',
							label: configuration.amounts.label,
							color: configuration.amounts.color
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
								<div class="flex flex-1 shrink-0 items-center justify-between text-xs leading-none">
									<div class="grid gap-1.5">
										<span class="text-muted-foreground">{name}</span>
									</div>
									<p class="font-mono">{value}</p>
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
				trend >= 0 ? 'text-emerald-500 dark:text-emerald-400' : 'text-rose-500 dark:text-rose-400'
			)}
		>
			{Math.abs(trend).toFixed(2)} %
			{#if trend >= 0}
				<Icon icon="ph:caret-up" />
			{:else}
				<Icon icon="ph:caret-down" />
			{/if}
		</Card.Footer>
	</Card.Root>
{/if}
