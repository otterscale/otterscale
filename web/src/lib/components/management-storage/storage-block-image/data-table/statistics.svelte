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
			Image
		{/snippet}
		{#snippet content()}
			{@const nameList = filteredData.map((datum) => datum['name' as keyof TData])}
			{nameList.length}
		{/snippet}
	</Chart.Text>
	<Chart.Text>
		{#snippet title()}
			Objects
		{/snippet}
		{#snippet description()}
			{@const objectList = filteredData.map((datum) => datum['objects' as keyof TData] as number)}
			{objectList.reduce((a, value) => a + value, 0)} units
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
