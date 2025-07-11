<script lang="ts" module>
	import type { User_Key } from '$gen/api/storage/v1/storage_pb';
	import TableRowPicker from '$lib/components/custom/data-table/data-table-row-pickers/cell.svelte';
	import type { Row } from '@tanstack/table-core';

	export const cells = {
		_row_picker: _row_picker,
		access: access
	};
</script>

{#snippet _row_picker(row: Row<User_Key>)}
	<TableRowPicker {row} />
{/snippet}

{#snippet access(row: Row<User_Key>)}
	{row.original.accessKey}
{/snippet}
