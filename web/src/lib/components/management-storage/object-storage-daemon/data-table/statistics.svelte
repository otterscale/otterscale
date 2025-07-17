<script lang="ts" module>
	import { Statistics as Layout } from '$lib/components/custom/chart/layouts/index';
	import * as Chart from '$lib/components/custom/chart/templates';
	import { type Table } from '@tanstack/table-core';
</script>

<script lang="ts" generics="TData">
	let { table }: { table: Table<TData> } = $props();

	const filteredData = $derived(table.getFilteredRowModel().rows.map((row) => row.original));
</script>

<Layout>
	<Chart.Text>
		{#snippet title()}
			OSD
		{/snippet}
		{#snippet content()}
			{@const idList = filteredData.map((datum) => datum['id' as keyof TData])}
			{idList.length}
		{/snippet}
	</Chart.Text>
	<Chart.Text>
		{#snippet title()}
			Usage
		{/snippet}
		{#snippet content()}
			{@const usedList = filteredData.map((datum) => Number(datum['usedBytes' as keyof TData]))}
			{@const usedTotal = usedList.reduce((a, value) => a + value, 0)}
			{@const sizeList = filteredData.map((datum) => Number(datum['sizeBytes' as keyof TData]))}
			{@const sizeTotal = sizeList.reduce((a, value) => a + value, 0)}

			{Number((usedTotal * 100) / sizeTotal).toFixed(1)}%
		{/snippet}
	</Chart.Text>
</Layout>
