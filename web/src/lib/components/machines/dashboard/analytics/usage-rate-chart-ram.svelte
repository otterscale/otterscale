<script lang="ts">
	import Icon from '@iconify/svelte';
	import { ArcChart, Text } from 'layerchart';
	import { PrometheusDriver } from 'prometheus-query';

	import * as Statistics from '$lib/components/custom/data-table/statistics/index';
	import * as Chart from '$lib/components/ui/chart/index.js';
	import { formatCapacity, formatPercentage } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';

	let { client, fqdn }: { client: PrometheusDriver; fqdn: string } = $props();

	const queries = $derived({
		overall: `sum(node_memory_MemTotal_bytes{instance=~"${fqdn}"})`,
		using: `sum(node_memory_MemTotal_bytes{instance=~"${fqdn}"}) - sum(node_memory_MemAvailable_bytes{instance=~"${fqdn}"})`,
		total: `sum(node_memory_MemTotal_bytes{instance=~"${fqdn}"})`
	});

	async function fetchUsage() {
		const using = await client.instantQuery(queries.using);
		const total = await client.instantQuery(queries.total);

		return { using, total };
	}
</script>

<Statistics.Root type="ratio">
	<Statistics.Header>
		<div class="flex justify-between gap-4">
			<Statistics.Title>{m.ram()}</Statistics.Title>
			{#await client.instantQuery(queries.overall) then response}
				{#if response.result[0]?.value?.value}
					{@const { value, unit } = formatCapacity(response.result[0].value.value)}
					<div class="flex items-center gap-1 text-xl">
						<p class="font-bold">{value} {unit}</p>
					</div>
				{/if}
			{/await}
		</div>
	</Statistics.Header>
	<Statistics.Content class="min-h-20">
		{#await fetchUsage()}
			<div class="flex h-full w-full items-center justify-center">
				<Icon icon="svg-spinners:6-dots-rotate" class="m-8 size-16" />
			</div>
		{:then response}
			{#if response.using.result[0]?.value && response.total.result[0]?.value}
				{@const chartConfig = {
					data: { color: 'var(--chart-3)' }
				} satisfies Chart.ChartConfig}
				{@const using = response.using.result[0].value.value}
				{@const total = response.total.result[0].value.value}
				{@const value = using / total}
				{@const data = [{ value: value }]}
				<Chart.Container
					config={chartConfig}
					class="mx-auto my-auto aspect-square h-[250px] w-full"
				>
					<ArcChart
						{data}
						innerRadius={-15}
						cornerRadius={15}
						range={[-120, 120]}
						maxValue={1}
						series={[
							{
								key: 'data',
								color: chartConfig.data.color
							}
						]}
						props={{
							arc: { track: { fill: 'var(--muted)' }, motion: 'tween' }
						}}
						tooltip={false}
					>
						{#snippet aboveMarks()}
							{@const percentage = formatPercentage(using, total, 1)}
							{@const { value: usingValue, unit: usingUnit } = formatCapacity(using)}
							{@const { value: totalValue, unit: totalUnit } = formatCapacity(total)}
							<Text
								value={`${percentage} %`}
								textAnchor="middle"
								verticalAnchor="middle"
								class="fill-foreground text-3xl! font-bold"
								dy={-15}
							/>
							<Text
								value={`${usingValue} ${usingUnit}/${totalValue} ${totalUnit}`}
								textAnchor="middle"
								verticalAnchor="middle"
								class="font-base text-md! text-muted-foreground"
								dy={15}
							/>
						{/snippet}
					</ArcChart>
				</Chart.Container>
			{:else}
				<div class="flex h-full w-full flex-col items-center justify-center">
					<Icon icon="ph:chart-bar-fill" class="size-24 animate-pulse text-muted-foreground" />
					<p class="text-base text-muted-foreground">{m.no_data_display()}</p>
				</div>
			{/if}
		{:catch}
			<div class="flex h-full w-full flex-col items-center justify-center">
				<Icon icon="ph:chart-bar-fill" class="size-24 animate-pulse text-muted-foreground" />
				<p class="text-base text-muted-foreground">{m.no_data_display()}</p>
			</div>
		{/await}
	</Statistics.Content>
</Statistics.Root>
