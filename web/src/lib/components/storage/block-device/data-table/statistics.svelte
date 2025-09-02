<script lang="ts" module>
	import Content from '$lib/components/custom/chart/content/text/text-large.svelte';
	import Layout from '$lib/components/custom/chart/layout/small-flexible-height.svelte';
	import Title from '$lib/components/custom/chart/title.svelte';
	import Description from '$lib/components/custom/chart/description.svelte';
	import { type Table } from '@tanstack/table-core';
	import { formatCapacity } from '$lib/formatter';
	import Icon from '@iconify/svelte';
</script>

<script lang="ts" generics="TData">
	let { table }: { table: Table<TData> } = $props();

	const filteredData = $derived(table.getFilteredRowModel().rows.map((row) => row.original));
</script>

<div class="grid w-full gap-3 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4 xl:grid-cols-5">
	<Layout>
		{#snippet title()}
			<Title title="Image" />
		{/snippet}

		{#snippet content()}
			{@const nameList = filteredData.map((datum) => datum['name' as keyof TData])}
			<Content value={nameList.length} />
		{/snippet}
	</Layout>
	<Layout>
		{#snippet title()}
			<Title title="Usage" />
		{/snippet}

		{#snippet content()}
			{@const quotaList = filteredData.map((datum) => Number(datum['quotaBytes' as keyof TData]))}
			{@const quotaTotal = quotaList.reduce((a, current) => a + current, 0)}
			{@const usedList = filteredData.map((datum) => Number(datum['usedBytes' as keyof TData]))}
			{@const usedTotal = usedList.reduce((a, current) => a + current, 0)}
			{#if quotaTotal && usedTotal / quotaTotal < 1}
				<Content value={Math.round((usedTotal / quotaTotal) * 100)} unit="%" />
			{:else}
				<Content>
					<Icon icon="ph:infinity" />
				</Content>
			{/if}
		{/snippet}
	</Layout>
</div>
