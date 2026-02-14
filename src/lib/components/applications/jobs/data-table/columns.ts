import { timestampDate } from '@bufbuild/protobuf/wkt';
import type { ColumnDef } from '@tanstack/table-core';

import type { Job } from '$lib/api/application/v1/application_pb';
import { getSortingFunction } from '$lib/components/custom/data-table/core';
import { renderSnippet } from '$lib/components/ui/data-table/index.js';
import { m } from '$lib/paraglide/messages';

import { cells } from './cells.svelte';
import { headers } from './headers.svelte';

const messages = {
	name: m.name(),
	namespace: m.namespace(),
	ready: m.ready(),
	succeeded: m.succeeded(),
	failed: m.failed(),
	terminating: m.terminating(),
	status: m.status(),
	conditions: m.conditions(),
	startedAt: m.started_at(),
	completedAt: m.completed_at()
};

function getColumns(): ColumnDef<Job>[] {
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
			accessorKey: 'status',
			header: ({ column }) => {
				return renderSnippet(headers.status, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.status, row);
			}
		},
		{
			accessorKey: 'conditions',
			header: ({ column }) => {
				return renderSnippet(headers.conditions, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.conditions, row);
			}
		},
		{
			accessorKey: 'active',
			header: ({ column }) => {
				return renderSnippet(headers.active, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.active, row);
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
			accessorKey: 'terminating',
			header: ({ column }) => {
				return renderSnippet(headers.terminating, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.terminating, row);
			}
		},
		{
			accessorKey: 'succeeded',
			header: ({ column }) => {
				return renderSnippet(headers.succeeded, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.succeeded, row);
			}
		},
		{
			accessorKey: 'failed',
			header: ({ column }) => {
				return renderSnippet(headers.failed, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.failed, row);
			}
		},
		{
			accessorKey: 'startedAt',
			header: ({ column }) => {
				return renderSnippet(headers.startedAt, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.startedAt, row);
			},
			sortingFn: (previousRow, nextRow) =>
				getSortingFunction(
					previousRow.original.startedAt,
					nextRow.original.startedAt,
					(p, n) => timestampDate(p) < timestampDate(n),
					(p, n) => timestampDate(p) === timestampDate(n)
				)
		},
		{
			accessorKey: 'completedAt',
			header: ({ column }) => {
				return renderSnippet(headers.completedAt, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.completedAt, row);
			},
			sortingFn: (previousRow, nextRow) =>
				getSortingFunction(
					previousRow.original.completedAt,
					nextRow.original.completedAt,
					(p, n) => timestampDate(p) < timestampDate(n),
					(p, n) => timestampDate(p) === timestampDate(n)
				)
		}
	];
}

export { getColumns, messages };
