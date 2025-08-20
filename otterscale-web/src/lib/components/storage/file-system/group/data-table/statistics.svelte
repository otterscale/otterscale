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
			<Title title="Pool" />
		{/snippet}

		{#snippet content()}
			{@const nameList = filteredData.map((datum) => datum['name' as keyof TData])}
			<Content value={nameList.length} />
		{/snippet}
	</Layout>
</div>
