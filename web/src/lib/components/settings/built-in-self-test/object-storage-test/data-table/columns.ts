import { timestampDate } from '@bufbuild/protobuf/wkt';
import { type ColumnDef } from '@tanstack/table-core';

import { page } from '$app/state';
import { type TestResult } from '$lib/api/configuration/v1/configuration_pb';
import { getSortingFunction } from '$lib/components/custom/data-table/core';
import { renderSnippet } from '$lib/components/ui/data-table/index.js';
import { m } from '$lib/paraglide/messages';

import { cells } from './cells.svelte';
import { headers } from './headers.svelte';

const messages = {
	name: m.name(),
	status: m.status(),
	target: m.target(),
	operation: m.operation(),
	duration: m.duration(),
	objectSize: m.object_size(),
	objectCount: m.object_count(),
	throughputFastest: m.throughput_fastest(),
	throughputSlowest: m.throughput_slowest(),
	throughputMedian: m.throughput_median(),
	createdBy: m.created_by(),
	startedAt: m.started_at(),
	completedAt: m.completed_at()
};

const columns: ColumnDef<TestResult>[] = [
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
		accessorKey: 'status',
		header: ({ column }) => {
			return renderSnippet(headers.status, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.status, row);
		}
	},
	{
		accessorKey: 'target',
		header: ({ column }) => {
			return renderSnippet(headers.target, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.target, row);
		}
	},
	{
		accessorKey: 'createdBy',
		header: ({ column }) => {
			return renderSnippet(headers.createdBy, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.createdBy, row);
		},
		filterFn: 'arrIncludesSome'
	},
	{
		accessorKey: 'operation',
		header: ({ column }) => {
			return renderSnippet(headers.operation, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.operation, row);
		}
	},
	{
		accessorKey: 'duration',
		header: ({ column }) => {
			return renderSnippet(headers.duration, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.duration, row);
		}
	},
	{
		accessorKey: 'objectSize',
		header: ({ column }) => {
			return renderSnippet(headers.objectSize, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.objectSize, row);
		}
	},
	{
		accessorKey: 'objectCount',
		header: ({ column }) => {
			return renderSnippet(headers.objectCount, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.objectCount, row);
		}
	},
	{
		accessorKey: 'throughputFastest',
		header: ({ column }) => {
			return renderSnippet(headers.throughputFastest, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.throughputFastest, row);
		}
	},
	{
		accessorKey: 'throughputSlowest',
		header: ({ column }) => {
			return renderSnippet(headers.throughputSlowest, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.throughputSlowest, row);
		}
	},
	{
		accessorKey: 'throughputMedian',
		header: ({ column }) => {
			return renderSnippet(headers.throughputMedian, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.throughputMedian, row);
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
	},
	{
		accessorKey: 'actions',
		header: ({ column }) => {
			return renderSnippet(headers.actions, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.actions, row);
		},
		enableHiding: false
	}
];

export { columns, messages };
