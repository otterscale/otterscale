<script lang="ts">
	import { ArcChart, Text } from 'layerchart';
	import { PrometheusDriver } from 'prometheus-query';
	import { onDestroy, onMount } from 'svelte';

	import type { Scope } from '$lib/api/scope/v1/scope_pb';
	import ComponentLoading from '$lib/components/custom/chart/component-loading.svelte';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Card from '$lib/components/ui/card';
	import * as Chart from '$lib/components/ui/chart/index.js';
	import { formatCapacity } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';

	// Props
	let {
		client,
		scope,
		isReloading = $bindable()
	}: { client: PrometheusDriver; scope: Scope; isReloading: boolean } = $props();

	// Constants
	const CHART_TITLE = m.capacity();
	const CHART_DESCRIPTION = m.remaining_capacity();
	const chartConfig = { data: { color: 'var(--chart-2)' } } satisfies Chart.ChartConfig;

	// Queries
	const queries = $derived({
		used: `ceph_cluster_total_used_bytes{juju_model_uuid=~"${scope.uuid}"}`,
		total: `ceph_cluster_total_bytes{juju_model_uuid=~"${scope.uuid}"}`
	});

	// Auto Update
	let response = $state(
		{} as {
			usedValue: number | undefined;
			usedUnit: string | undefined;
			totalValue: number | undefined;
			totalUnit: string | undefined;
			usage: { value: number }[];
		}
	);
	let isLoading = $state(true);
	const reloadManager = new ReloadManager(fetch);

	// Data fetching function
	async function fetch() {
		const [usedResponse, totalResponse] = await Promise.all([
			client.instantQuery(queries.used),
			client.instantQuery(queries.total)
		]);

		const usedValue = usedResponse.result[0]?.value?.value;
		const totalValue = totalResponse.result[0]?.value?.value;

		const usedCapacity = usedValue ? formatCapacity(usedValue) : null;
		const totalCapacity = totalValue ? formatCapacity(totalValue) : null;
		const usageValue = usedValue / totalValue;
		const usagePercentage = usageValue != null ? usageValue * 100 : null;

		response = {
			usedValue: usedCapacity ? Math.round(usedCapacity.value) : undefined,
			usedUnit: usedCapacity ? usedCapacity.unit : undefined,
			totalValue: totalCapacity ? Math.round(totalCapacity.value) : undefined,
			totalUnit: totalCapacity ? totalCapacity.unit : undefined,
			usage: usagePercentage !== null ? [{ value: usagePercentage }] : [{ value: NaN }]
		};
	}

	// Effects
	$effect(() => {
		if (isReloading) {
			reloadManager.restart();
		} else {
			reloadManager.stop();
		}
	});

	onMount(() => {
		fetch();
		isLoading = false;
	});
	onDestroy(() => {
		reloadManager.stop();
	});
</script>

{#if isLoading}
	<ComponentLoading />
{:else}
	<Card.Root class="h-full gap-2">
		<Card.Header class="items-center">
			<Card.Title>{CHART_TITLE}</Card.Title>
			<Card.Description>{CHART_DESCRIPTION}</Card.Description>
		</Card.Header>
		<Card.Content class="flex-1">
			<Chart.Container config={chartConfig} class="mx-auto aspect-square max-h-[200px]">
				<ArcChart
					data={response.usage}
					outerRadius={88}
					innerRadius={66}
					trackOuterRadius={83}
					trackInnerRadius={72}
					padding={40}
					range={[90, -270]}
					maxValue={100}
					series={[
						{
							key: 'data',
							color: chartConfig.data.color
						}
					]}
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
							value={response.usedValue}
							textAnchor="middle"
							verticalAnchor="middle"
							class="fill-foreground text-4xl! font-bold"
							dy={3}
						/>
						<Text
							value={response.usedUnit}
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
{/if}
