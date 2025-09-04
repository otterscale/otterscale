<script lang="ts">
	import { page } from '$app/state';
	import { MachineService, type Machine } from '$lib/api/machine/v1/machine_pb';
	import type { Scope } from '$lib/api/scope/v1/scope_pb';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Card from '$lib/components/ui/card';
	import * as Chart from '$lib/components/ui/chart';
	import { m } from '$lib/paraglide/messages';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { PieChart, Text } from 'layerchart';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	let { scope, isReloading = $bindable() }: { scope: Scope; isReloading: boolean } = $props();

	const transport: Transport = getContext('transport');
	const machineClient = createClient(MachineService, transport);

	const machines = writable<Machine[]>([]);
	const scopeMachines = $derived(
		$machines.filter((m) => m.workloadAnnotations['juju-machine-id']?.startsWith(page.params.scope!)),
	);
	const totalNodes = $derived(scopeMachines.length);

	const virtualNodes = $derived(scopeMachines.filter((m) => m.tags.includes('virtual')).length);
	const physicalNodes = $derived(scopeMachines.length - virtualNodes);

	const nodeProportions = $derived([
		{ node: 'physical', nodes: physicalNodes, color: 'var(--color-physical)' },
		{ node: 'virtual', nodes: virtualNodes, color: 'var(--color-virtual)' },
	]);

	const nodeProportionsConfiguration = {
		nodes: { label: 'Nodes' },
		physical: { label: 'Physical', color: 'var(--chart-1)' },
		virtual: { label: 'Virtual', color: 'var(--chart-2)' },
	} satisfies Chart.ChartConfig;

	async function fetch() {
		machineClient.listMachines({ scopeUuid: scope.uuid }).then((response) => {
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
	<Card.Root class="flex h-full flex-col">
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
{/if}
