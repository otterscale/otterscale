import type { ColumnDef } from '@tanstack/table-core';

import { cells } from './cells.svelte';
import { headers } from './headers.svelte';

import type { Network_IPRange } from '$lib/api/network/v1/network_pb';
import { renderSnippet } from '$lib/components/ui/data-table/index.js';
import { m } from '$lib/paraglide/messages';

const messages = {
	type: m.type(),
	startIp: m.start_ip(),
	endIp: m.end_ip(),
	comment: m.comment()
};

const columns: ColumnDef<Network_IPRange>[] = [
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
	}
];

export { columns, messages };
