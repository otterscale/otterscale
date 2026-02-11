<script lang="ts">
	import type { JsonValue } from '@bufbuild/protobuf';
	import Circle from '@lucide/svelte/icons/circle';
	import CircleX from '@lucide/svelte/icons/circle-x';
	import { type Column, type Row } from '@tanstack/table-core';

	let {
		row,
		column
	}: {
		row: Row<Record<string, JsonValue>>;
		column: Column<Record<string, JsonValue>>;
	} = $props();

	const data = $derived(row.original[column.id] as boolean);
</script>

{#if data === true}
	<Circle class="inline-block size-4 text-primary" />
{:else if data === false}
	<CircleX class="inline-block size-6 text-destructive" />
{/if}
