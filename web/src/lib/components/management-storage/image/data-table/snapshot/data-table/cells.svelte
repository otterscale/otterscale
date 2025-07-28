<script lang="ts" module>
	import type { Image_Snapshot } from '$gen/api/storage/v1/storage_pb';
	import TableRowPicker from '$lib/components/custom/data-table/data-table-row-pickers/cell.svelte';
	import { Badge } from '$lib/components/ui/badge';
	import Icon from '@iconify/svelte';
	import type { Row } from '@tanstack/table-core';

	export const cells = {
		_row_picker: _row_picker,
		name: name,
		protect: protect
	};
</script>

{#snippet _row_picker(row: Row<Image_Snapshot>)}
	<TableRowPicker {row} />
{/snippet}

{#snippet name(row: Row<Image_Snapshot>)}
	{row.original.name}
{/snippet}

{#snippet protect(row: Row<Image_Snapshot>)}
	<div class="flex justify-end">
		{#if row.original.protected}
			<Icon icon="ph:circle" class="text-primary" />
		{:else}
			<Icon icon="ph:x" class="text-destructive" />
		{/if}
	</div>
{/snippet}
