<script lang="ts" module>
	import { Statistics as Layout } from '$lib/components/custom/chart/layouts/index';
	import * as Chart from '$lib/components/custom/chart/templates/index';
	import { type Table } from '@tanstack/table-core';
</script>

<script lang="ts" generics="TData">
	let { table }: { table: Table<TData> } = $props();

	const filteredData = $derived(table.getFilteredRowModel().rows.map((row) => row.original));
</script>

<Layout class="lg:grid-cols-2">
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
