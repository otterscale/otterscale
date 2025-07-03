<script lang="ts" generics="TData">
	import * as Chart from '$lib/components/custom/chart/templates';
	import { type Table } from '@tanstack/table-core';
	import { Statistics as Layout } from '$lib/components/custom/chart/layouts/index';

	let { table }: { table: Table<TData> } = $props();

	const filteredData = $derived(table.getFilteredRowModel().rows.map((row) => row.original));
</script>

<Layout>
	<Chart.Text>
		{#snippet title()}
			Snapshot
		{/snippet}
		{#snippet content()}
			{@const nameList = filteredData.map((datum) => datum['name' as keyof TData])}
			{nameList.length}
		{/snippet}
	</Chart.Text>
</Layout>
