<script lang="ts" generics="TData">
	import type { TestResult } from '$gen/api/bist/v1/bist_pb';
	import { Statistics as Layout } from '$lib/components/custom/chart/layouts/index';
	import * as Chart from '$lib/components/custom/chart/templates';
	import { StatisticManager } from '$lib/components/custom/data-table/utils.svelte';
	import { type Table } from '@tanstack/table-core';

	let { table }: { table: Table<TestResult> } = $props();

	const statisticManager = new StatisticManager(table);
</script>

<Layout>
	<Chart.Text>
		{#snippet title()}
			OST
		{/snippet}
		{#snippet description()}
			Internal/External Storage Tester
		{/snippet}
		{#snippet content()}
			<div class="flex justify-between">
				<div class="text-7xl">
					{statisticManager.count('name')}
				</div>
			</div>
		{/snippet}
	</Chart.Text>
</Layout>

