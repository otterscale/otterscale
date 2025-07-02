<script lang="ts" generics="TData">
	import { Statistics as Layout } from '$lib/components/custom/chart/layouts/index';
	import * as Chart from '$lib/components/custom/chart/templates';
	import { type Table } from '@tanstack/table-core';

	let { table }: { table: Table<TData> } = $props();

	const filteredData = $derived(table.getFilteredRowModel().rows.map((row) => row.original));
</script>

<Layout>
	<Chart.Text>
		{#snippet title()}
			Monitor
		{/snippet}
		{#snippet content()}
			{@const idList = filteredData.map((datum) => datum['name' as keyof TData])}
			{idList.length}
		{/snippet}
	</Chart.Text>
</Layout>
