import { renderSnippet } from '$lib/components/ui/data-table/index.js';
import { m } from '$lib/paraglide/messages';
import type { ColumnDef } from '@tanstack/table-core';
import { type LargeLangeageModel } from '../protobuf.svelte';
import { cells } from './cells.svelte';
import { headers } from './headers.svelte';

const messages = {
	name: m.models_large_language_model_headers_name(),
	version: m.models_large_language_model_headers_version(),
	parameters: m.models_large_language_model_headers_parameters(),
	accuracy: m.models_large_language_model_headers_accuracy(),
	speed: m.models_large_language_model_headers_speed(),
	architecture: m.models_large_language_model_headers_architecture(),
	requests: m.models_large_language_model_headers_requests(),
	uptime: m.models_large_language_model_headers_uptime()
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
