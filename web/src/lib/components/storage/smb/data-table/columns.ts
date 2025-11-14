import type { ColumnDef } from '@tanstack/table-core';

import type { SMBShare } from '$lib/api/storage/v1/storage_pb';
import { getSortingFunction } from '$lib/components/custom/data-table/core';
import type { ReloadManager } from '$lib/components/custom/reloader';
import { renderSnippet } from '$lib/components/ui/data-table/index.js';
import { m } from '$lib/paraglide/messages';

import { cells } from './cells.svelte';
import { headers } from './headers.svelte';

const messages = {
	name: m.name(),
	namespace: m.namespace(),
	status: m.status(),
	size: m.size(),
	browsable: m.browsable(),
	readOnly: m.read_only(),
	guestOk: m.guest_accessible(),
	mapToGuest: m.map_to_guest(),
	securityMode: m.security_mode(),
	auth: m.auth(),
	validUsers: m.valid_users(),
	actions: m.actions()
};

function getColumns(
	scope: string,
	namespace: string,
	reloadManager: ReloadManager
): ColumnDef<SMBShare>[] {
	return [
		{
			id: 'select',
			header: ({ table }) => renderSnippet(headers.row_picker, table),
			cell: ({ row }) => renderSnippet(cells.row_picker, row),
			enableSorting: false,
			enableHiding: false
		},
		{
			accessorKey: 'name',
			header: ({ column }) => renderSnippet(headers.name, column),
			cell: ({ row }) => renderSnippet(cells.name, row)
		},
		{
			accessorKey: 'namespace',
			header: ({ column }) => renderSnippet(headers.namespace, column),
			cell: ({ row }) => renderSnippet(cells.namespace, row)
		},
		{
			accessorKey: 'status',
			header: ({ column }) => renderSnippet(headers.status, column),
			cell: ({ row }) => renderSnippet(cells.status, row)
		},
		{
			accessorKey: 'size',
			header: ({ column }) => renderSnippet(headers.size, column),
			cell: ({ row }) => renderSnippet(cells.size, row),
			sortingFn: (previousRow, nextRow) =>
				getSortingFunction(
					Number(previousRow.original.sizeBytes),
					Number(nextRow.original.sizeBytes),
					(p, n) => p < n,
					(p, n) => p === n
				)
		},
		{
			accessorKey: 'browsable',
			header: ({ column }) => renderSnippet(headers.browsable, column),
			cell: ({ row }) => renderSnippet(cells.browsable, row),
			enableSorting: false
		},
		{
			accessorKey: 'read_only',
			header: ({ column }) => renderSnippet(headers.read_only, column),
			cell: ({ row }) => renderSnippet(cells.read_only, row),
			enableSorting: false
		},
		{
			accessorKey: 'guest_ok',
			header: ({ column }) => renderSnippet(headers.guest_ok, column),
			cell: ({ row }) => renderSnippet(cells.guest_ok, row),
			enableSorting: false
		},
		{
			accessorKey: 'map_to_guest',
			header: ({ column }) => renderSnippet(headers.map_to_guest, column),
			cell: ({ row }) => renderSnippet(cells.map_to_guest, row),
			sortingFn: (previousRow, nextRow) =>
				getSortingFunction(
					Number(previousRow.original.mapToGuest),
					Number(nextRow.original.mapToGuest),
					(p, n) => p < n,
					(p, n) => p === n
				)
		},
		{
			accessorKey: 'security_mode',
			header: ({ column }) => renderSnippet(headers.security_mode, column),
			cell: ({ row }) => renderSnippet(cells.security_mode, row),
			sortingFn: (previousRow, nextRow) =>
				getSortingFunction(
					Number(previousRow.original.securityMode),
					Number(nextRow.original.securityMode),
					(p, n) => p < n,
					(p, n) => p === n
				)
		},
		{
			accessorKey: 'valid_users',
			header: ({ column }) => renderSnippet(headers.valid_users, column),
			cell: ({ row }) => renderSnippet(cells.valid_users, row),
			sortingFn: (previousRow, nextRow) =>
				getSortingFunction(
					previousRow.original.validUsers.length,
					nextRow.original.validUsers.length,
					(p, n) => p < n,
					(p, n) => p === n
				)
		},
		{
			accessorKey: 'actions',
			header: ({ column }) => renderSnippet(headers.actions, column),
			cell: ({ row }) => renderSnippet(cells.actions, { row, scope, namespace, reloadManager }),
			enableHiding: false
		}
	];
}

export { getColumns, messages };
