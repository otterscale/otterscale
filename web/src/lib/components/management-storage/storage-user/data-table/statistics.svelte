<script lang="ts" generics="TData">
	import * as Chart from '$lib/components/custom/chart/templates';
	import { StatisticManager } from '$lib/components/custom/data-table/utils.svelte';
	import { type Table } from '@tanstack/table-core';

	let { table }: { table: Table<TData> } = $props();

	const statisticManager = new StatisticManager(table);
</script>

<div class="grid w-full gap-4 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4">
	<span class="col-span-1">
		<Chart.Text>
			{#snippet title()}
				User
			{/snippet}
			{#snippet content()}
				<div class="flex justify-between">
					<div class="text-7xl">
						{statisticManager.count('userId')}
					</div>
				</div>
			{/snippet}
		</Chart.Text>
	</span>
</div>
