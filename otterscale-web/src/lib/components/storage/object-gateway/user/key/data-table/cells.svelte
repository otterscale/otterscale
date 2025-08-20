<script lang="ts" module>
	import type { User_Key } from '$lib/api/storage/v1/storage_pb';
	import { Cell as RowPicker } from '$lib/components/custom/data-table/data-table-row-pickers';
	import type { Row } from '@tanstack/table-core';
	import Actions from './cell-actions.svelte';

	export const cells = {
		row_picker,
		accessKey,
		actions
	};
</script>

{#snippet row_picker(row: Row<User_Key>)}
	<RowPicker {row} />
{/snippet}

{#snippet accessKey(row: Row<User_Key>)}
	{row.original.accessKey}
{/snippet}

{#snippet actions(row: Row<User_Key>)}
	<Actions key={row.original} />
{/snippet}
