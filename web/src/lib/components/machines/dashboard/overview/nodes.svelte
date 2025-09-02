<script lang="ts">
	import { page } from '$app/state';
	import { MachineService, type Machine } from '$lib/api/machine/v1/machine_pb';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import { buttonVariants } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import * as Chart from '$lib/components/ui/chart';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { m } from '$lib/paraglide/messages';
	import { getLocale } from '$lib/paraglide/runtime';
	import { timestampDate } from '@bufbuild/protobuf/wkt';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { scaleUtc } from 'd3-scale';
	import { curveNatural } from 'd3-shape';
	import { LineChart } from 'layerchart';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	const now = new Date();
	const months: string[] = [];
	for (let i = 5; i >= 0; i--) {
		const d = new Date(now.getFullYear(), now.getMonth() - i + 1, 1);
		months.push(d.toISOString().slice(0, 7));
	}

	let { isReloading = $bindable() }: { isReloading: boolean } = $props();

	const transport: Transport = getContext('transport');
	const machineClient = createClient(MachineService, transport);

	const machines = writable<Machine[]>([]);
	const scopeMachines = $derived(
		$machines.filter((m) =>
			m.workloadAnnotations['juju-machine-id']?.startsWith(page.params.scope!)
		)
	);
	const totalNodes = $derived(scopeMachines.length);
	const monthlyCounts = $derived(
		scopeMachines.reduce(
			(acc, m) => {
				const dateStr = m.lastCommissioned
					? timestampDate(m.lastCommissioned).toISOString().slice(0, 7)
					: null;
				if (dateStr) {
					acc[dateStr] = (acc[dateStr] || 0) + 1;
				}
				return acc;
			},
			{} as Record<string, number>
		)
	);
	const nodes = $derived(
		months.map((month) => ({
			date: new Date(month + '-01'),
			node: monthlyCounts[month] || 0
		}))
	);

	const nodesConfiguration = {
		node: { label: 'Node', color: 'var(--chart-1)' }
	} satisfies Chart.ChartConfig;

	async function fetch() {
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
				<div class="flex flex-col items-start gap-0.5 truncate text-sm font-medium tracking-tight">
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
{/if}
