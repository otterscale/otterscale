<script lang="ts">
	import { type Table } from '@tanstack/table-core';

	import type { Model } from '$lib/api/model/v1/model_pb';
	import * as Statistics from '$lib/components/custom/data-table/statistics/index';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import { m } from '$lib/paraglide/messages';

	let { serviceUri, table }: { serviceUri: string; table: Table<Model> } = $props();

	const filteredModels = $derived(table.getFilteredRowModel().rows.map((row) => row.original));
</script>

<div class="grid w-full gap-3 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4 xl:grid-cols-5">
	{#if serviceUri}
		<Statistics.Root type="data">
			<Statistics.Header>
				<Statistics.Title>{m.service()}</Statistics.Title>
				<p class="text-center text-lg font-semibold">{serviceUri}</p>
			</Statistics.Header>
			<Statistics.Footer class="space-x-1">
				<Badge>chat</Badge>
				<Badge>completions</Badge>
				<Badge>embedding</Badge>
			</Statistics.Footer>
			<Statistics.Background icon="ph:squares-four" />
		</Statistics.Root>
	{/if}
	<Statistics.Root type="count">
		{@const models = filteredModels.length}
		<Statistics.Header>
			<Statistics.Title>{m.models()}</Statistics.Title>
		</Statistics.Header>
		<Statistics.Content>
			<p class="text-6xl font-semibold">{models}</p>
		</Statistics.Content>
		<Statistics.Background icon="ph:robot" />
	</Statistics.Root>
</div>
