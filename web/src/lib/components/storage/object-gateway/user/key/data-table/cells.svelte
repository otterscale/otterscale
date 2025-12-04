<script lang="ts" module>
	import type { Row } from '@tanstack/table-core';

	import type { User, User_Key } from '$lib/api/storage/v1/storage_pb';
	import { Cells } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { ReloadManager } from '$lib/components/custom/reloader';

	import Actions from './cell-actions.svelte';

	export const cells = {
		row_picker,
		accessKey,
		actions
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

{#snippet actions(data: {
	row: Row<User_Key>;
	user: User;
	scope: string;
	reloadManager: ReloadManager;
})}
	<Layout.Cell class="items-start">
		<Actions
			key={data.row.original}
			user={data.user}
			scope={data.scope}
			reloadManager={data.reloadManager}
		/>
	</Layout.Cell>
{/snippet}
