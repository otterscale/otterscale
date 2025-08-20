<!-- <script lang="ts">
	import { Statistics as Layout } from '$lib/components/custom/chart/layouts/index';
	import * as Chart from '$lib/components/custom/chart/templates';
	import Icon from '@iconify/svelte';
	import { type Table } from '@tanstack/table-core';
	import { formatCapacity } from '$lib/formatter';
	import Icon from '@iconify/svelte';
</script>

<script lang="ts" generics="TData">
	let { table }: { table: Table<TData> } = $props();

	const filteredData = $derived(table.getFilteredRowModel().rows.map((row) => row.original));

	const quotaTotal = $derived(
		filteredData.reduce((sum, datum) => sum + Number(datum['quotaBytes' as keyof TData]), 0)
	);
	const usedTotal = $derived(
		filteredData.reduce((sum, datum) => sum + Number(datum['usedBytes' as keyof TData]), 0)
	);
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
