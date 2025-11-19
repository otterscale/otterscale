import type { ColumnDef } from '@tanstack/table-core';

import type { Network_IPRange } from '$lib/api/network/v1/network_pb';
import type { ReloadManager } from '$lib/components/custom/reloader';
import { renderSnippet } from '$lib/components/ui/data-table/index.js';
import { m } from '$lib/paraglide/messages';

import { cells } from './cells.svelte';
import { headers } from './headers.svelte';

const messages = {
	type: m.type(),
	startIp: m.start_ip(),
	endIp: m.end_ip(),
	comment: m.comment()
};

function getColumns(reloadManager: ReloadManager): ColumnDef<Network_IPRange>[] {
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
			accessorKey: 'startIp',
			header: ({ column }) => {
				return renderSnippet(headers.startIp, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.startIp, row);
			}
		},
		{
			accessorKey: 'endIp',
			header: ({ column }) => {
				return renderSnippet(headers.endIp, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.endIp, row);
			}
		},
		{
			accessorKey: 'type',
			header: ({ column }) => {
				return renderSnippet(headers.type, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.type, row);
			},
			filterFn: 'arrIncludesSome'
		},
		{
			accessorKey: 'comment',
			header: ({ column }) => {
				return renderSnippet(headers.comment, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.comment, row);
			}
		},
		{
			accessorKey: 'actions',
			header: ({ column }) => {
				return renderSnippet(headers.actions, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.actions, { row, reloadManager });
			},
			enableHiding: false
		}
	];
}

export { getColumns, messages };
