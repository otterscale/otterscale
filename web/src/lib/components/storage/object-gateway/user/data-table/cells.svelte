<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import type { Row } from '@tanstack/table-core';

	import type { User } from '$lib/api/storage/v1/storage_pb';
	import { Cells } from '$lib/components/custom/data-table/core';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Table from '$lib/components/custom/table/index.js';
	import { Key } from '$lib/components/storage/object-gateway/user/key';
	import Badge from '$lib/components/ui/badge/badge.svelte';

	import Actions from './cell-actions.svelte';

	export const cells = {
		row_picker,
		id,
		name,
		suspended,
		keys,
		actions
	};
</script>

{#snippet row_picker(row: Row<User>)}
	<Table.Cell alignClass="items-center">
		<Cells.RowPicker {row} />
	</Table.Cell>
{/snippet}

{#snippet id(row: Row<User>)}
	<Table.Cell alignClass="items-start">
		{row.original.id}
	</Table.Cell>
{/snippet}

{#snippet name(row: Row<User>)}
	<Table.Cell alignClass="items-start">
		<Badge variant="outline">{row.original.name}</Badge>
	</Table.Cell>
{/snippet}

{#snippet suspended(row: Row<User>)}
	<Table.Cell alignClass="items-end">
		{#if row.original.suspended}
			<Icon icon="ph:circle" class="text-primary" />
		{:else}
			<Icon icon="ph:x" class="text-destructive" />
		{/if}
	</Table.Cell>
{/snippet}

{#snippet keys(data: { row: Row<User>; scope: string; reloadManager: ReloadManager })}
	<Table.Cell alignClass="items-end">
		<Key user={data.row.original} scope={data.scope} reloadManager={data.reloadManager} />
	</Table.Cell>
{/snippet}

{#snippet actions(data: { row: Row<User>; scope: string; reloadManager: ReloadManager })}
	<Table.Cell alignClass="items-start">
		<Actions user={data.row.original} scope={data.scope} reloadManager={data.reloadManager} />
	</Table.Cell>
{/snippet}
