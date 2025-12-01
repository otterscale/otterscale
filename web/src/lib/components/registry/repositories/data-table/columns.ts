import type { ColumnDef } from '@tanstack/table-core';

import type { Repository } from '$lib/api/registry/v1/registry_pb';
import type { ReloadManager } from '$lib/components/custom/reloader';
import { renderSnippet } from '$lib/components/ui/data-table/index.js';
import { m } from '$lib/paraglide/messages';

import { cells } from './cells.svelte';
import { headers } from './headers.svelte';

const messages = {
	name: m.name(),
	manifest_count: m.manifest_count(),
	size_bytes: m.size(),
	latest_tag: m.latest_tag()
};

function getColumns(scope: string, reloadManager: ReloadManager): ColumnDef<Repository>[] {
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
			accessorKey: 'name',
			header: ({ column }) => {
				return renderSnippet(headers.name, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.name, row);
			}
		},
		{
			accessorKey: 'manifest_count',
			header: ({ column }) => {
				return renderSnippet(headers.manifest_count, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.manifest_count, { row, scope, reloadManager });
			}
		},
		{
			accessorKey: 'size_bytes',
			header: ({ column }) => {
				return renderSnippet(headers.size_bytes, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.size_bytes, row);
			}
		},
		{
			accessorKey: '	',
			header: ({ column }) => {
				return renderSnippet(headers.latest_tag, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.latest_tag, row);
			}
		}
		// {
		// 	accessorKey: 'actions',
		// 	header: ({ column }) => {
		// 		return renderSnippet(headers.actions, column);
		// 	},
		// 	cell: ({ row }) => {
		// 		return renderSnippet(cells.actions, { data: { row, scope, reloadManager } });
		// 	},
		// 	enableHiding: false
		// }
	];
}

export { getColumns, messages };
