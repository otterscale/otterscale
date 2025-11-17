import type { ColumnDef } from '@tanstack/table-core';

import type { VirtualMachine_Snapshot } from '$lib/api/instance/v1/instance_pb';
import { renderSnippet } from '$lib/components/ui/data-table/index.js';
import { m } from '$lib/paraglide/messages';

import { cells } from './cells.svelte';
import { headers } from './headers.svelte';
import type { ReloadManager } from '$lib/components/custom/reloader';

const messages = {
	name: m.name(),
	namespace: m.namespace(),
	sourceName: m.source_name(),
	phase: m.phase(),
	ready: m.ready(),
	createTime: m.create_time()
};

function getColumns(
	scope: string,
	reloadManager: ReloadManager
): ColumnDef<VirtualMachine_Snapshot>[] {
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
			accessorKey: 'namespace',
			header: ({ column }) => {
				return renderSnippet(headers.namespace, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.namespace, row);
			}
		},
		{
			accessorKey: 'sourceName',
			header: ({ column }) => {
				return renderSnippet(headers.sourceName, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.sourceName, row);
			}
		},
		{
			accessorKey: 'phase',
			header: ({ column }) => {
				return renderSnippet(headers.phase, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.phase, row);
			}
		},
		{
			accessorKey: 'ready',
			header: ({ column }) => {
				return renderSnippet(headers.ready, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.ready, row);
			}
		},
		{
			accessorKey: 'createTime',
			header: ({ column }) => {
				return renderSnippet(headers.createTime, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.createTime, row);
			}
		},

		{
			accessorKey: 'actions',
			header: ({ column }) => {
				return renderSnippet(headers.actions, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.actions, { row, scope, reloadManager });
			},
			enableHiding: false
		}
	];
}

export { getColumns, messages };
