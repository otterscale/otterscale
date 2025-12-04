<script lang="ts" module>
	import type { Row } from '@tanstack/table-core';

	import type { User, User_Key } from '$lib/api/storage/v1/storage_pb';
	import { Cells } from '$lib/components/custom/data-table/core';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Table from '$lib/components/custom/table/index.js';

	import Actions from './cell-actions.svelte';

	export const cells = {
		row_picker,
		accessKey,
		actions
	};
</script>

{#snippet row_picker(row: Row<User_Key>)}
	<Table.Cell alignClass="items-center">
		<Cells.RowPicker {row} />
	</Table.Cell>
{/snippet}

{#snippet accessKey(row: Row<User_Key>)}
	<Table.Cell alignClass="items-start">
		{row.original.accessKey}
	</Table.Cell>
{/snippet}

{#snippet actions(data: {
	row: Row<User_Key>;
	user: User;
	scope: string;
	reloadManager: ReloadManager;
})}
	<Table.Cell alignClass="items-start">
		<Actions
			key={data.row.original}
			user={data.user}
			scope={data.scope}
			reloadManager={data.reloadManager}
		/>
	</Table.Cell>
{/snippet}
