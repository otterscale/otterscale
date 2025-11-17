import type { ColumnDef } from '@tanstack/table-core';

import type { User_Key } from '$lib/api/storage/v1/storage_pb';
import type { User } from '$lib/api/storage/v1/storage_pb';
import type { ReloadManager } from '$lib/components/custom/reloader';
import { renderSnippet } from '$lib/components/ui/data-table/index.js';
import { m } from '$lib/paraglide/messages';

import { cells } from './cells.svelte';
import { headers } from './headers.svelte';

const messages = {
	accessKey: m.access_key()
};

function getColumns(
	scope: string,
	user: User,
	reloadManager: ReloadManager
): ColumnDef<User_Key>[] {
	return [
		{
			id: 'select',
			header: ({ table }) => {
				return renderSnippet(headers.row_picker, table);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.row_picker, row);
			},
			enableSorting: false,
			enableHiding: false
		},
		{
			accessorKey: 'accessKey',
			header: ({ column }) => {
				return renderSnippet(headers.accessKey, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.accessKey, row);
			}
		},
		{
			accessorKey: 'actions',
			header: ({ column }) => {
				return renderSnippet(headers.actions, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.actions, { row, user, scope, reloadManager });
			},
			enableHiding: false
		}
	];
}

export { getColumns, messages };
