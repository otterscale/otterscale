<script lang="ts" generics="TData">
	import * as Chart from '$lib/components/custom/chart/templates';
	import { formatCapacity } from '$lib/formatter';
	import { type Table } from '@tanstack/table-core';

	let { table }: { table: Table<TData> } = $props();

	const filteredData = $derived(table.getFilteredRowModel().rows.map((row) => row.original));
</script>

<div class="grid w-full gap-4 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4">
	<span class="col-span-1">
		<Chart.Text>
			{#snippet title()}
				OSD
			{/snippet}
			{#snippet content()}
				{@const idList = filteredData.map((datum) => datum['id' as keyof TData])}
				<div class="flex justify-between">
					<div class="text-7xl">
						{idList.length}
					</div>
				</div>
			{/snippet}
		</Chart.Text>
	</span>
	<span class="col-span-1">
		<Chart.Text>
			{#snippet title()}
				Size
			{/snippet}
			{#snippet content()}
				{@const sizeList = filteredData.map((datum) => datum['size' as keyof TData] as number)}
				{@const { value, unit } = formatCapacity(sizeList.reduce((a, value) => a + value, 0))}
				<span class="text-7xl">
					{value}
				</span>
				<span class="text-6xl">
					{unit}
				</span>
			{/snippet}
		</Chart.Text>
	</span>
</div>
