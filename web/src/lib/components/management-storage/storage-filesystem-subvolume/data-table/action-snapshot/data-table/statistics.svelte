<script lang="ts" generics="TData">
	import * as Chart from '$lib/components/custom/chart/templates';
	import { type Table } from '@tanstack/table-core';

	let { table }: { table: Table<TData> } = $props();

	const filteredData = $derived(table.getFilteredRowModel().rows.map((row) => row.original));
</script>

<div class="grid w-full gap-4 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4">
	<span class="col-span-1">
		<Chart.Text>
			{#snippet title()}
				Snapshot
			{/snippet}
			{#snippet content()}
				{@const nameList = filteredData.map((datum) => datum['name' as keyof TData])}
				<div class="flex justify-between">
					<div class="text-7xl">
						{nameList.length}
					</div>
				</div>
			{/snippet}
		</Chart.Text>
	</span>
</div>
