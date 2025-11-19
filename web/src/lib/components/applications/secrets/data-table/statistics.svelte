<script lang="ts">
	import type { Table } from '@tanstack/table-core';

	import type { Secret } from '$lib/api/application/v1/application_pb';
	import Content from '$lib/components/custom/chart/content/text/text-large.svelte';
	import Layout from '$lib/components/custom/chart/layout/small-flexible-height.svelte';
	import Title from '$lib/components/custom/chart/title.svelte';

	let { table }: { table: Table<Secret> } = $props();
</script>

<div class="grid w-full gap-3 sm:grid-cols-1 md:grid-cols-3 lg:grid-cols-5">
	<Layout>
		{#snippet title()}
			<Title title="Secrets" />
		{/snippet}

		{#snippet content()}
			<Content value={table.getCoreRowModel().rows.length} />
		{/snippet}
	</Layout>

	<Layout>
		{#snippet title()}
			<Title title="Types" />
		{/snippet}

		{#snippet content()}
			<Content
				value={new Set([...table.getCoreRowModel().rows.map((row) => row.getValue('type'))]).size}
			/>
		{/snippet}
	</Layout>
</div>
