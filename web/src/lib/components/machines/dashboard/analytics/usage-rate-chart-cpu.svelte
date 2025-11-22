<script lang="ts">
	import Icon from '@iconify/svelte';
	import { ArcChart, Text } from 'layerchart';
	import { PrometheusDriver } from 'prometheus-query';

	import * as Statistics from '$lib/components/custom/data-table/statistics/index';
	import * as Chart from '$lib/components/ui/chart/index.js';
	import { formatPercentage } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';

	let { client, fqdn }: { client: PrometheusDriver; fqdn: string } = $props();

	const queries = $derived({
		count: `count(count by (cpu, instance) (node_cpu_seconds_total{instance=~"${fqdn}"}))`,
		using: `sum(irate(node_cpu_seconds_total{instance=~"${fqdn}",mode!="idle"}[6m]))`,
		total: `sum(irate(node_cpu_seconds_total{instance=~"${fqdn}"}[6m]))`
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
			<Statistics.Title>{m.cpu()}</Statistics.Title>
			{#await client.instantQuery(queries.count) then response}
				{#if response.result[0]?.value?.value}
					{@const value = response.result[0].value.value}
					<div class="flex items-center gap-1 text-xl">
						<p class="font-bold">{value} {m.cpu()}</p>
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
							<Text
								value={`${percentage} %`}
								textAnchor="middle"
								verticalAnchor="middle"
								class="fill-foreground text-4xl! font-bold"
							/>
						{/snippet}
					</ArcChart>
				</Chart.Container>
				<Statistics.Progress numerator={using} denominator={total} target="STB" />
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
	<Statistics.Background icon="ph:cpu" class="-top-20 right-20 size-10" />
</Statistics.Root>
