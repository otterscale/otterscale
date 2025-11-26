import type { ColumnDef } from '@tanstack/table-core';

import type { Model } from '$lib/api/model/v1/model_pb';
import type { ReloadManager } from '$lib/components/custom/reloader';
import { renderSnippet } from '$lib/components/ui/data-table/index.js';
import { m } from '$lib/paraglide/messages';

import { cells } from './cells.svelte';
import { headers } from './headers.svelte';

const messages = {
	name: m.name(),
	namespace: m.namespace(),
	status: m.status(),
	description: m.description(),
	first_deployed_at: m.time(),
	last_deployed_at: m.time(),
	chart_version: m.version(),
	app_version: m.version(),
	prefill: m.prefill(),
	decode: m.decode(),
	gpu_relation: m.gpu_relation()
};

function getColumns(scope: string, reloadManager: ReloadManager): ColumnDef<Model>[] {
	return [
		{
			id: 'expand',
			header: ({ table }) => {
				return renderSnippet(headers.row_expander, table);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.row_expander, row);
			},
			enableSorting: false,
			enableHiding: false
		},
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
			accessorKey: 'id',
			header: ({ column }) => {
				return renderSnippet(headers.id, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.id, row);
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
			accessorKey: 'status',
			header: ({ column }) => {
				return renderSnippet(headers.status, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.status, row);
			}
		},
		{
			accessorKey: 'description',
			header: ({ column }) => {
				return renderSnippet(headers.description, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.description, row);
			}
		},

		{
			accessorKey: 'chart_version',
			header: ({ column }) => {
				return renderSnippet(headers.chart_version, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.chart_version, row);
			}
		},
		{
			accessorKey: 'app_version',
			header: ({ column }) => {
				return renderSnippet(headers.app_version, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.app_version, row);
			}
		},
		{
			accessorKey: 'first_deployed_at',
			header: ({ column }) => {
				return renderSnippet(headers.first_deployed_at, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.first_deployed_at, row);
			}
		},
		{
			accessorKey: 'last_deployed_at',
			header: ({ column }) => {
				return renderSnippet(headers.last_deployed_at, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.last_deployed_at, row);
			}
		},
		{
			accessorKey: 'prefill',
			header: ({ column }) => {
				return renderSnippet(headers.prefill, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.prefill, row);
			},
			enableHiding: false
		},
		{
			accessorKey: 'decode',
			header: ({ column }) => {
				return renderSnippet(headers.decode, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.decode, row);
			},
			enableHiding: false
		},
		{
			accessorKey: 'gpu_relation',
			header: ({ column }) => {
				return renderSnippet(headers.gpu_relation, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.gpu_relation, { row, scope, reloadManager });
			},
			enableHiding: false
		},
		{
			accessorKey: 'action',
			header: ({ column }) => {
				return renderSnippet(headers.action, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.action, { row, scope, reloadManager });
			},
			enableHiding: false
		}
	];
}
export { getColumns, messages };
