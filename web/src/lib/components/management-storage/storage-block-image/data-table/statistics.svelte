<script lang="ts" generics="TData">
	import { Statistics as Layout } from '$lib/components/custom/chart/layouts/index';
	import * as Chart from '$lib/components/custom/chart/templates';
	import Icon from '@iconify/svelte';
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
			Usage
		{/snippet}
		{#snippet content()}
			{@const quotaBytesList = filteredData.map((datum) =>
				Number(datum['quotaBytes' as keyof TData])
			)}
			{@const quotaBytesTotal = quotaBytesList.reduce((a, current) => a + current, 0)}
			{@const usedBytesList = filteredData.map((datum) =>
				Number(datum['usedBytes' as keyof TData])
			)}
			{@const usedBytesTotal = usedBytesList.reduce((a, current) => a + current, 0)}
			{#if quotaBytesTotal}
				{((usedBytesTotal / quotaBytesTotal) * 100).toFixed(2)}%
			{:else}
				<Icon icon="ph:infinity" />
			{/if}
		{/snippet}
	</Chart.Text>
</Layout>
