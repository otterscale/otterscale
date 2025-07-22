<script lang="ts" generics="TData">
	import type { TestResult } from '$gen/api/bist/v1/bist_pb';
	import * as Chart from '$lib/components/custom/chart/templates';
	import { StatisticManager } from '$lib/components/custom/data-table/utils.svelte';
	import { type Table } from '@tanstack/table-core';
	
	import { Statistics as Layout } from '$lib/components/custom/chart/layouts/index';


    // 確保泛型 TData 被正確設定為 TestResult
    let { table }: { table: Table<TestResult> } = $props();

	const statisticManager = new StatisticManager(table);
</script>

<Layout>
	<Chart.Text>
		{#snippet title()}
			I/O Tester
		{/snippet}
		{#snippet description()}
			Ceph Block and NFS Tester
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
