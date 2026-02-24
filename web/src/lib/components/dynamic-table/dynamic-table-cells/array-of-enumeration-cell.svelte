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
</script>

{#if data && data.length > 0}
	<div class="flex items-center gap-1">
		{#each data as datum, index (index)}
			<Badge variant="outline">{datum}</Badge>
		{/each}
	</div>
{:else}
	<Badge variant="secondary">No Data</Badge>
{/if}
