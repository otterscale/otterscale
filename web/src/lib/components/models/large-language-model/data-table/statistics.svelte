<script lang="ts" module>
	import { type Table } from '@tanstack/table-core';

	import Content from '$lib/components/custom/chart/content/text/text-large.svelte';
	import Layout from '$lib/components/custom/chart/layout/small-flexible-height.svelte';
	import Title from '$lib/components/custom/chart/title.svelte';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts" generics="TData">
	let { table }: { table: Table<TData> } = $props();

	const filteredData = $derived(table.getFilteredRowModel().rows.map((row) => row.original));
</script>

<div class="grid w-full gap-3 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4 xl:grid-cols-5">
	<Layout>
		{#snippet title()}
			<Title title={m.models()} />
		{/snippet}

		{#snippet content()}
			{@const nameList = filteredData.map((datum) => datum['name' as keyof TData])}
			<Content value={nameList.length} />
		{/snippet}
	</Layout>
	<Layout>
		{#snippet title()}
			<Title title={m.requests()} />
		{/snippet}

		{#snippet content()}
			{@const requestsList = filteredData.map((datum) => datum['requests' as keyof TData])}
			{@const totalRequests = requestsList.reduce((a, requests) => a + Number(requests), 0)}
			<Content value={totalRequests} />
		{/snippet}
	</Layout>
</div>
