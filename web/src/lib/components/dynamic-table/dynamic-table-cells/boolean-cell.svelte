<script lang="ts">
	import type { JsonValue } from '@bufbuild/protobuf';
	import { CircleIcon, XIcon } from '@lucide/svelte';
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
	<CircleIcon class="inline-block size-4 text-primary" />
{:else if data === false}
	<XIcon class="inline-block size-4 text-destructive" />
{/if}
