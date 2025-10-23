<script lang="ts">
	import { type Table } from '@tanstack/table-core';

	import type { TestResult } from '$lib/api/configuration/v1/configuration_pb';
	import Content from '$lib/components/custom/chart/content/scatter/scatter.svelte';
	import Description from '$lib/components/custom/chart/description.svelte';
	import Layout from '$lib/components/custom/chart/layout/small.svelte';
	import Title from '$lib/components/custom/chart/title.svelte';
	import { BistDashboardManager } from '$lib/components/settings/built-in-self-test/utils/bistManager';

	let { mode, table }: { mode: string; table: Table<TestResult> } = $props();
	const dashboardManager = new BistDashboardManager(table);
	const { read: readTmp, write: writeTmp, trim: trimTmp } = dashboardManager.getFioOutputs();
</script>

<div class="my-6 grid w-full gap-4 sm:grid-cols-1 md:grid-cols-1 lg:grid-cols-3">
	<Layout>
		{#snippet title()}
			<Title title="Bandwidth" />
		{/snippet}

		{#snippet description()}
			<Description description="Bandwidth Bytes" />
		{/snippet}

		{#snippet content()}
			<Content
				xKey="bandwidthBytes"
				yKey="completedAt"
				data={mode === 'read' ? readTmp : mode === 'write' ? writeTmp : trimTmp}
			/>
		{/snippet}
	</Layout>
	<Layout>
		{#snippet title()}
			<Title title="IOPS" />
		{/snippet}

		{#snippet description()}
			<Description description="IO Per Second" />
		{/snippet}

		{#snippet content()}
			<Content
				xKey="ioPerSecond"
				yKey="completedAt"
				data={mode === 'read' ? readTmp : mode === 'write' ? writeTmp : trimTmp}
			/>
		{/snippet}
	</Layout>
	<Layout>
		{#snippet title()}
			<Title title="Latency" />
		{/snippet}

		{#snippet description()}
			<Description description="Mean Latency" />
		{/snippet}

		{#snippet content()}
			<Content
				xKey="latency"
				yKey="completedAt"
				data={mode === 'read' ? readTmp : mode === 'write' ? writeTmp : trimTmp}
			/>
		{/snippet}
	</Layout>
</div>
