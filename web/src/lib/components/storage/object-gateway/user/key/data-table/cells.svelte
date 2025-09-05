<script lang="ts" module>
	import type { Row } from '@tanstack/table-core';

	import Actions from './cell-actions.svelte';

	import type { User_Key } from '$lib/api/storage/v1/storage_pb';
	import { Cells } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';

	export const cells = {
		row_picker,
		accessKey,
		actions,
	};
</script>

{#snippet row_picker(row: Row<User_Key>)}
	<Layout.Cell class="items-center">
		<Cells.RowPicker {row} />
	</Layout.Cell>
{/snippet}

{#snippet accessKey(row: Row<User_Key>)}
	<Layout.Cell class="items-start">
		{row.original.accessKey}
	</Layout.Cell>
{/snippet}

{#snippet actions(row: Row<User_Key>)}
	<Layout.Cell class="items-start">
		<Actions key={row.original} />
	</Layout.Cell>
{/snippet}
