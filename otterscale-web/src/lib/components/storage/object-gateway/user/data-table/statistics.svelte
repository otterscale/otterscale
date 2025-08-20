<!-- <script lang="ts">
	import { Statistics as Layout } from '$lib/components/custom/chart/layouts/index';
	import * as Chart from '$lib/components/custom/chart/templates';
	import { type Table } from '@tanstack/table-core';
	import { formatCapacity } from '$lib/formatter';
	import Icon from '@iconify/svelte';
</script>

<script lang="ts" generics="TData">
	let { table }: { table: Table<TData> } = $props();

	const filteredData = $derived(table.getFilteredRowModel().rows.map((row) => row.original));
</script>

<div class="grid grid-cols-5 gap-3">
	<Layout>
		{#snippet title()}
			<Title title="User" />
		{/snippet}

		{#snippet content()}
			{@const nameList = filteredData.map((datum) => datum['userId' as keyof TData])}
			<Content value={nameList.length} />
		{/snippet}
	</Layout>
</div>
