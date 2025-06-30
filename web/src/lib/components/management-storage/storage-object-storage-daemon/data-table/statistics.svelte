<script lang="ts" generics="TData">
	import { Statistics as Layout } from '$lib/components/custom/chart/layouts/index';
	import * as Chart from '$lib/components/custom/chart/templates';
	import { formatCapacity } from '$lib/formatter';
	import { type Table } from '@tanstack/table-core';

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
			Size
		{/snippet}
		{#snippet content()}
			{@const sizeList = filteredData.map((datum) => datum['size' as keyof TData] as number)}
			{@const { value, unit } = formatCapacity(sizeList.reduce((a, value) => a + value, 0))}
			{Number(value).toFixed(1)}
			<span class="font-light">
				{unit}
			</span>
		{/snippet}
	</Chart.Text>
</Layout>
