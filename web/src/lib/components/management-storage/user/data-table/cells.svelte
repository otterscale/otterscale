<script lang="ts" module>
	import type { User } from '$gen/api/storage/v1/storage_pb';
	import TableRowPicker from '$lib/components/custom/data-table/data-table-row-pickers/cell.svelte';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import Icon from '@iconify/svelte';
	import type { Row } from '@tanstack/table-core';

	export const cells = {
		_row_picker: _row_picker,
		id: id,
		name: name,
		suspended: suspended
	};
</script>

{#snippet _row_picker(row: Row<User>)}
	<TableRowPicker {row} />
{/snippet}

{#snippet id(row: Row<User>)}
	{row.original.id}
{/snippet}

{#snippet name(row: Row<User>)}
	<Badge variant="outline">{row.original.name}</Badge>
{/snippet}

{#snippet suspended(row: Row<User>)}
	{#if row.original.suspended}
		<Icon icon="ph:x" class="text-destructive" />
	{:else}
		<Icon icon="ph:check" class="text-success" />
	{/if}
{/snippet}
