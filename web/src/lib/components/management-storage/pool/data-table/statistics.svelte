<script lang="ts" module>
	import { Statistics as Layout } from '$lib/components/custom/chart/layouts/index';
	import * as Chart from '$lib/components/custom/chart/templates';
	import Icon from '@iconify/svelte';
	import { type Table } from '@tanstack/table-core';
</script>

<script lang="ts" generics="TData">
	import { formatCapacity } from '$lib/formatter';

	let { table }: { table: Table<TData> } = $props();

	const filteredData = $derived(table.getFilteredRowModel().rows.map((row) => row.original));
</script>

<Layout>
	<Chart.Text>
		{#snippet title()}
			Pool
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
		{#snippet description()}
			{@const quotaBytesList = filteredData.map((datum) =>
				Number(datum['quotaBytes' as keyof TData])
			)}
			{@const quotaBytesTotal = quotaBytesList.reduce((a, current) => a + current, 0)}
			{@const { value: totalValue, unit: totalUnit } = formatCapacity(quotaBytesTotal)}
			{@const usedBytesList = filteredData.map((datum) =>
				Number(datum['usedBytes' as keyof TData])
			)}
			{@const usedBytesTotal = usedBytesList.reduce((a, current) => a + current, 0)}
			{@const { value: usedValue, unit: usedUnit } = formatCapacity(usedBytesTotal)}
			{usedValue}
			{usedUnit}
			/
			{totalValue}
			{totalUnit}
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
			{#if quotaBytesTotal && usedBytesTotal / quotaBytesTotal < 1}
				{((usedBytesTotal / quotaBytesTotal) * 100).toFixed(2)}%
			{:else}
				<Icon icon="ph:infinity" />
			{/if}
		{/snippet}
	</Chart.Text>
</Layout>
