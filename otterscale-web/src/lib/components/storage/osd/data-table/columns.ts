import type { OSD } from '$lib/api/storage/v1/storage_pb';
import { getSortingFunction } from '$lib/components/custom/data-table/core';
import { renderSnippet } from '$lib/components/ui/data-table/index.js';
import type { ColumnDef } from '@tanstack/table-core';
import { cells } from './cells.svelte';
import { headers } from './headers.svelte';
import { m } from '$lib/paraglide/messages';

const messages = {
	name: m.name(),
	state: m.state(),
	in: m.osd_ins(),
	up: m.osd_ups(),
	exists: m.osd_exists(),
	deviceClass: m.device_class(),
	machine: m.machine(),
	placementGroupCount: m.placement_group(),
	usage: m.usage(),
	iops: m.iops()
};

const columns: ColumnDef<OSD>[] = [
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
		accessorKey: 'state',
		header: ({ column }) => {
			return renderSnippet(headers.state, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.state, row);
		}
	},
	{
		id: 'in',
		header: ({ column }) => {
			return renderSnippet(headers.osdIn, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.osdIn, row);
		},
		filterFn: (row, columnId, filterValue: boolean) => {
			if (filterValue === undefined) {
				return true;
			}

			return row.original.in === filterValue;
		},
		enableHiding: false
	},
	{
		id: 'up',
		header: ({ column }) => {
			return renderSnippet(headers.osdUp, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.osdUp, row);
		},
		filterFn: (row, columnId, filterValue: boolean) => {
			if (filterValue === undefined) {
				return true;
			}

			return row.original.up === filterValue;
		},
		enableHiding: false
	},
	{
		accessorKey: 'exists',
		header: ({ column }) => {
			return renderSnippet(headers.exists, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.exists, row);
		},
		filterFn: 'equals'
	},
	{
		accessorKey: 'deviceClass',
		header: ({ column }) => {
			return renderSnippet(headers.deviceClass, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.deviceClass, row);
		},
		filterFn: 'arrIncludesSome'
	},
	{
		accessorKey: 'machine',
		header: ({ column }) => {
			return renderSnippet(headers.machine, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.machine, row);
		},
		filterFn: 'arrIncludesSome'
	},
	{
		accessorKey: 'placementGroupCount',
		header: ({ column }) => {
			return renderSnippet(headers.placementGroupCount, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.placementGroupCount, row);
		}
	},
	{
		accessorKey: 'usage',
		header: ({ column }) => {
			return renderSnippet(headers.usage, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.usage, row);
		},
		sortingFn: (previousRow, nextRow, columnId) =>
			getSortingFunction(
				Number(previousRow.original.usedBytes) / Number(previousRow.original.sizeBytes),
				Number(nextRow.original.usedBytes) / Number(nextRow.original.sizeBytes),
				(p, n) => p < n,
				(p, n) => p === n
			)
	},
	{
		accessorKey: 'iops',
		header: ({ column }) => {
			return renderSnippet(headers.iops, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.iops, row);
		}
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
