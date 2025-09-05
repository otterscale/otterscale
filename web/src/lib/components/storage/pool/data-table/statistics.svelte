<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import { type Table } from '@tanstack/table-core';

	import Content from '$lib/components/custom/chart/content/text/text-large.svelte';
	import Description from '$lib/components/custom/chart/description.svelte';
	import Layout from '$lib/components/custom/chart/layout/small-flexible-height.svelte';
	import Title from '$lib/components/custom/chart/title.svelte';
	import { formatCapacity } from '$lib/formatter';
</script>

<script lang="ts" generics="TData">
	let { table }: { table: Table<TData> } = $props();

	const filteredData = $derived(table.getFilteredRowModel().rows.map((row) => row.original));

	const quotaTotal = $derived(
		filteredData.reduce((sum, datum) => sum + Number(datum['quotaBytes' as keyof TData]), 0),
	);
	const usedTotal = $derived(filteredData.reduce((sum, datum) => sum + Number(datum['usedBytes' as keyof TData]), 0));
</script>

<div class="grid w-full gap-3 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4 xl:grid-cols-5">
	<Layout>
		{#snippet title()}
			<Title title="Pool" />
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

		{#snippet description()}
			{@const { value: totalValue, unit: totalUnit } = formatCapacity(quotaTotal)}
			{@const { value: usedValue, unit: usedUnit } = formatCapacity(usedTotal)}
			<Description description={`${usedValue} ${usedUnit} / ${totalValue} ${totalUnit}`} />
		{/snippet}

		{#snippet content()}
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
