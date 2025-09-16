<script lang="ts">
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { scaleUtc } from 'd3-scale';
	import { curveLinear } from 'd3-shape';
	import { LineChart } from 'layerchart';
	import { PrometheusDriver, SampleValue } from 'prometheus-query';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import { page } from '$app/state';
	import { MachineService, type Machine } from '$lib/api/machine/v1/machine_pb';
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
		isReloading = $bindable(),
	}: { prometheusDriver: PrometheusDriver; scope: Scope; isReloading: boolean } = $props();

	const transport: Transport = getContext('transport');
	const machineClient = createClient(MachineService, transport);

	const machines = writable<Machine[]>([]);
	const scopeMachines = $derived(
		$machines.filter((m) => m.workloadAnnotations['juju-machine-id']?.startsWith(page.params.scope!)),
	);
	const totalGPUs = $derived(scopeMachines.reduce((sum, m) => sum + Number(m.gpuDevices.length ?? 0), 0));
	let gpus = $state([] as SampleValue[]);
	const gpusTrend = $derived(
		gpus.length > 0 ? (gpus[gpus.length - 1].value - gpus[gpus.length - 2].value) / gpus[gpus.length - 2].value : 0,
	);

	const gpuUsagesConfiguration = {
		usage: { label: 'value', color: 'var(--chart-1)' },
	} satisfies Chart.ChartConfig;

	async function fetch() {
		prometheusDriver
			.rangeQuery(
				`sum(irate(DCGM_FI_DEV_GPU_UTIL{scope_uuid="${scope.uuid}"}[2m]))`,
				Date.now() - 10 * 60 * 1000,
				Date.now(),
				2 * 60,
			)
			.then((response) => {
				gpus = response.result ? response.result[0].values : [0];
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
							<Icon icon="ph:info" class="text-muted-foreground size-5" />
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
				<p class="text-muted-foreground text-sm">{m.pieces()}</p>
			</div>
			<Chart.Container config={gpuUsagesConfiguration} class="h-full w-20">
				<LineChart
					data={gpus}
					x="time"
					xScale={scaleUtc()}
					axis={false}
					series={[
						{
							key: 'value',
							label: 'usage',
							color: gpuUsagesConfiguration.usage.color,
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
		</Card.Content>
		<Card.Footer
			class={cn(
				'flex flex-wrap items-center justify-end text-sm leading-none font-medium',
				gpusTrend >= 0 ? 'text-emerald-500 dark:text-emerald-400' : 'text-rose-500 dark:text-rose-400',
			)}
		>
			{Math.abs(gpusTrend).toFixed(2)} %
			{#if gpusTrend >= 0}
				<Icon icon="ph:caret-up" />
			{:else}
				<Icon icon="ph:caret-down" />
			{/if}
		</Card.Footer>
	</Card.Root>
{/if}
