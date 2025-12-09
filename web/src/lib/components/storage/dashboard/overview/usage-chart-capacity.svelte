<script lang="ts">
	import Icon from '@iconify/svelte';
	import { ArcChart, Text } from 'layerchart';
	import { PrometheusDriver } from 'prometheus-query';
	import { onDestroy, onMount } from 'svelte';

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
	}: { client: PrometheusDriver; scope: string; isReloading: boolean } = $props();

	// Constants
	const CHART_TITLE = m.capacity();
	const CHART_DESCRIPTION = m.remaining_capacity();
	const chartConfig = { data: { color: 'var(--chart-2)' } } satisfies Chart.ChartConfig;

	// Queries
	const queries = $derived({
		used: `ceph_cluster_total_used_bytes{juju_model="${scope}"}`,
		total: `ceph_cluster_total_bytes{juju_model="${scope}"}`
	});

	// Auto Update
	let response = $state(
		{} as {
			usedValue: number | undefined;
			usedUnit: string | undefined;
			availableValue: number | undefined;
			availableUnit: string | undefined;
			totalValue: number | undefined;
			totalUnit: string | undefined;
			availablePercentage: { value: number }[];
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
		const availableBytes = totalValue - usedValue;
		const availableCapacity = availableBytes ? formatCapacity(availableBytes) : null;
		const totalCapacity = totalValue ? formatCapacity(totalValue) : null;
		const availableRatio = availableBytes / totalValue;
		const availablePercentage = availableRatio != null ? availableRatio * 100 : null;

		response = {
			usedValue: usedCapacity ? Math.round(usedCapacity.value) : undefined,
			usedUnit: usedCapacity ? usedCapacity.unit : undefined,
			availableValue: availableCapacity ? Math.round(availableCapacity.value) : undefined,
			availableUnit: availableCapacity ? availableCapacity.unit : undefined,
			totalValue: totalCapacity ? Math.round(totalCapacity.value) : undefined,
			totalUnit: totalCapacity ? totalCapacity.unit : undefined,
			availablePercentage: availablePercentage !== null ? [{ value: availablePercentage }] : [{ value: NaN }]
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

	onMount(async () => {
		await fetch();
		isLoading = false;
	});
	onDestroy(() => {
		reloadManager.stop();
	});
</script>

<Card.Root class="h-full gap-2">
	<Card.Header class="h-[42px]">
		<Card.Title>{CHART_TITLE}</Card.Title>
		<Card.Description>{CHART_DESCRIPTION}</Card.Description>
	</Card.Header>
	{#if isLoading}
		<Card.Content>
			<div class="flex h-[200px] w-full items-center justify-center">
				<Icon icon="svg-spinners:6-dots-rotate" class="size-12" />
			</div>
		</Card.Content>
	{:else if response.availableValue == undefined && response.availableUnit == undefined}
		<Card.Content>
			<div class="flex h-[200px] w-full flex-col items-center justify-center">
				<Icon icon="ph:chart-line-fill" class="size-50 animate-pulse text-muted-foreground" />
				<p class="text-base text-muted-foreground">{m.no_data_display()}</p>
			</div>
		</Card.Content>
	{:else}
		<Card.Content>
			<Chart.Container class="h-[200px] w-full px-2 pt-2" config={chartConfig}>
				<ArcChart
					data={response.availablePercentage}
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
							value="{response.availableValue} {response.availableUnit}"
							textAnchor="middle"
							verticalAnchor="middle"
							class="fill-foreground text-4xl! font-bold"
							dy={3}
						/>
						<Text
							value="Used {response.usedValue} {response.usedUnit}"
							textAnchor="middle"
							verticalAnchor="middle"
							class="fill-muted-foreground!"
							dy={25}
						/>
					{/snippet}
				</ArcChart>
			</Chart.Container>
		</Card.Content>
	{/if}
</Card.Root>
