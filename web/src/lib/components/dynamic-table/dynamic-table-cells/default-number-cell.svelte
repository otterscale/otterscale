<script lang="ts">
	import type { JsonValue } from '@bufbuild/protobuf';
	import { type Column, type Row } from '@tanstack/table-core';

	let {
		row,
		column
	}: {
		row: Row<Record<string, JsonValue>>;
		column: Column<Record<string, JsonValue>>;
	} = $props();

	const data = $derived(row.original[column.id] as number);
</script>

{data}
