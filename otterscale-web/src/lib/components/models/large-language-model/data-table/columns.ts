import { renderSnippet } from '$lib/components/ui/data-table/index.js';
import { m } from '$lib/paraglide/messages';
import type { ColumnDef } from '@tanstack/table-core';
import { type LargeLangeageModel } from '../protobuf.svelte';
import { cells } from './cells.svelte';
import { headers } from './headers.svelte';

const messages = {
	name: m.name(),
	version: m.version(),
	parameters: m.parameters(),
	accuracy: m.accuracy(),
	speed: m.speed(),
	architecture: m.architecture(),
	requests: m.requests(),
	uptime: m.uptime()
};

const columns: ColumnDef<LargeLangeageModel>[] = [
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
		accessorKey: 'version',
		header: ({ column }) => {
			return renderSnippet(headers.version, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.version, row);
		}
	},
	{
		accessorKey: 'parameters',
		header: ({ column }) => {
			return renderSnippet(headers.parameters, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.parameters, row);
		}
	},
	{
		accessorKey: 'accuracy',
		header: ({ column }) => {
			return renderSnippet(headers.accuracy, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.accuracy, row);
		}
	},
	{
		accessorKey: 'speed',
		header: ({ column }) => {
			return renderSnippet(headers.speed, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.speed, row);
		}
	},
	{
		accessorKey: 'architecture',
		header: ({ column }) => {
			return renderSnippet(headers.architecture, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.architecture, row);
		}
	},
	{
		accessorKey: 'requests',
		header: ({ column }) => {
			return renderSnippet(headers.requests, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.requests, row);
		}
	},
	{
		accessorKey: 'uptime',
		header: ({ column }) => {
			return renderSnippet(headers.uptime, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.uptime, row);
		}
	}
];

export { columns, messages };
