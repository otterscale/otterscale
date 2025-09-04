<script lang="ts">
	import type { TestResult } from '$lib/api/bist/v1/bist_pb';
	import Content from '$lib/components/custom/chart/content/scatter/scatter.svelte';
	import Description from '$lib/components/custom/chart/description.svelte';
	import Layout from '$lib/components/custom/chart/layout/small.svelte';
	import Title from '$lib/components/custom/chart/title.svelte';
	import { BistDashboardManager } from '$lib/components/settings/built-in-self-test/utils/bistManager';
	import { type Table } from '@tanstack/table-core';
	import Pickers from './pickers.svelte';

	let { table }: { table: Table<TestResult> } = $props();
	let mode = $state('get');
	const dashboardManager = new BistDashboardManager(table);
	const { get: getTmp, put: putTmp, delete: deleteTmp } = dashboardManager.getWarpOutputs();
</script>

<Pickers bind:selectedMode={mode} />

<div class="my-6 grid w-full gap-4 sm:grid-cols-1 md:grid-cols-1 lg:grid-cols-3">
	<Layout>
		{#snippet title()}
			<Title title="Throughput" />
		{/snippet}

		{#snippet description()}
			<Description description="Fastest" />
		{/snippet}

		{#snippet content()}
			<Content
				xKey="bytesFastest"
				yKey="completedAt"
				data={mode === 'get' ? getTmp : mode === 'put' ? putTmp : deleteTmp}
			/>
		{/snippet}
	</Layout>
	<Layout>
		{#snippet title()}
			<Title title="Throughput" />
		{/snippet}

		{#snippet description()}
			<Description description="Slowest" />
		{/snippet}

		{#snippet content()}
			<Content
				xKey="bytesSlowest"
				yKey="completedAt"
				data={mode === 'get' ? getTmp : mode === 'put' ? putTmp : deleteTmp}
			/>
		{/snippet}
	</Layout>
	<Layout>
		{#snippet title()}
			<Title title="Throughput" />
		{/snippet}

		{#snippet description()}
			<Description description="Median" />
		{/snippet}

		{#snippet content()}
			<Content
				xKey="latency"
				yKey="bytesMedian"
				data={mode === 'get' ? getTmp : mode === 'put' ? putTmp : deleteTmp}
			/>
		{/snippet}
	</Layout>
</div>
