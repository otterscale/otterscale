<!-- <script lang="ts">
	import { Statistics as Layout } from '$lib/components/custom/chart/layouts/index';
	import * as Chart from '$lib/components/custom/chart/templates';
	import { formatCapacity } from '$lib/formatter';
	import { type Table } from '@tanstack/table-core';
</script>

<script lang="ts" generics="TData">
	let { table }: { table: Table<TData> } = $props();

	const filteredData = $derived(table.getFilteredRowModel().rows.map((row) => row.original));
</script>

<div class="grid grid-cols-5 gap-3">
	<Layout>
		{#snippet title()}
			<Title title="Usage" />
		{/snippet}

		{#snippet content()}
			{@const nameList = filteredData.map((datum) => datum['name' as keyof TData])}
			<Content value={nameList.length} />
		{/snippet}
	</Layout>
	<Layout>
		{#snippet title()}
			<Title title="Bucket" />
		{/snippet}

		{#snippet content()}
			{@const bucketBytesList = filteredData.map((datum) =>
				Number(datum['usedBytes' as keyof TData])
			)}
			{@const bucketBytesTotal = bucketBytesList.reduce((a, current) => a + current, 0)}
			{@const { value, unit } = formatCapacity(Number(bucketBytesTotal) / (1024 * 1024))}
			<Content {value} {unit} />
		{/snippet}
	</Layout>
</div>
