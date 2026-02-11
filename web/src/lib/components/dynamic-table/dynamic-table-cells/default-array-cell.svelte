<script lang="ts">
	import type { JsonValue } from '@bufbuild/protobuf';
	import { type Column, type Row } from '@tanstack/table-core';

	import { Badge } from '$lib/components/ui/badge';

	let {
		row,
		column
	}: {
		row: Row<Record<string, JsonValue>>;
		column: Column<Record<string, JsonValue>>;
	} = $props();

	const data = $derived(row.original[column.id] as JsonValue[]);
	const hasObjectItem = $derived(data && data.some((value) => value && typeof value === 'object'));
</script>

{#if data && data.length > 0}
	{#if hasObjectItem}
		{data.length}
	{:else}
		<div class="flex items-center gap-1">
			{#each data as datum, index (index)}
				<Badge variant="outline">{datum}</Badge>
			{/each}
		</div>
	{/if}
{/if}
