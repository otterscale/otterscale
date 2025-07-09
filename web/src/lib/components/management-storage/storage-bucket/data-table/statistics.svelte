<script lang="ts" module>
	import { Statistics as Layout } from '$lib/components/custom/chart/layouts/index';
	import * as Chart from '$lib/components/custom/chart/templates';
	import { formatCapacity } from '$lib/formatter';
	import { type Table } from '@tanstack/table-core';
</script>

<script lang="ts" generics="TData">
	let { table }: { table: Table<TData> } = $props();

	const filteredData = $derived(table.getFilteredRowModel().rows.map((row) => row.original));
</script>

<Layout>
	<Chart.Text>
		{#snippet title()}
			Bucket
		{/snippet}
		{#snippet content()}
			{@const nameList = filteredData.map((datum) => datum['name' as keyof TData])}
			{nameList.length}
		{/snippet}
	</Chart.Text>
	<Chart.Text>
		{#snippet title()}
			Usage
		{/snippet}
		{#snippet content()}
			{@const usedBytesList = filteredData.map((datum) =>
				Number(datum['usedBytes' as keyof TData])
			)}
			{@const usedBytesTotal = usedBytesList.reduce((a, current) => a + current, 0)}
			{@const { value, unit } = formatCapacity(Number(usedBytesTotal) / (1024 * 1024))}
			{value}
			{unit}
		{/snippet}
	</Chart.Text>
</Layout>
