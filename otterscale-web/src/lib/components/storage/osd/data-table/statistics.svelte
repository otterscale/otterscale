<script lang="ts" module>
	import Content from '$lib/components/custom/chart/content/text/text-large.svelte';
	import Layout from '$lib/components/custom/chart/layout/small-flexible-height.svelte';
	import Title from '$lib/components/custom/chart/title.svelte';
	import { type Table } from '@tanstack/table-core';
</script>

<script lang="ts" generics="TData">
	let { table }: { table: Table<TData> } = $props();

	const filteredData = $derived(table.getFilteredRowModel().rows.map((row) => row.original));
</script>

<div class="grid grid-cols-5 gap-3">
	<Layout>
		{#snippet title()}
			<Title title="OSD" />
		{/snippet}

		{#snippet content()}
			{@const idList = filteredData.map((datum) => datum['id' as keyof TData])}
			<Content value={idList.length} />
		{/snippet}
	</Layout>
	<Layout>
		{#snippet title()}
			<Title title="Usage" />
		{/snippet}

		{#snippet content()}
			{@const usedList = filteredData.map((datum) => Number(datum['usedBytes' as keyof TData]))}
			{@const usedTotal = usedList.reduce((a, value) => a + value, 0)}
			{@const sizeList = filteredData.map((datum) => Number(datum['sizeBytes' as keyof TData]))}
			{@const sizeTotal = sizeList.reduce((a, value) => a + value, 0)}
			<Content value={Number((usedTotal * 100) / sizeTotal).toFixed(1)} unit={'%'} />
		{/snippet}
	</Layout>
</div>
