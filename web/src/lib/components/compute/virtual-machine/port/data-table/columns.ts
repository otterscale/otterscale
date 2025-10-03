import type { ColumnDef } from '@tanstack/table-core';

import { cells } from './cells.svelte';
import { headers } from './headers.svelte';

import type { Application_Service_Port } from '$lib/api/application/v1/application_pb';
import { renderSnippet } from '$lib/components/ui/data-table/index.js';
import { m } from '$lib/paraglide/messages';

const messages = {
	name: m.name(),
	protocol: m.protocol(),
	port: m.ports(),
	nodePort: m.node_port(),
};

const columns: ColumnDef<Application_Service_Port>[] = [
	{
		id: 'select',
		header: ({ table }) => {
			return renderSnippet(headers.row_picker, table);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.row_picker, row);
		},
		enableSorting: false,
		enableHiding: false,
	},
	{
		accessorKey: 'name',
		header: ({ column }) => {
			return renderSnippet(headers.name, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.name, row);
		},
	},
	{
		accessorKey: 'protocol',
		header: ({ column }) => {
			return renderSnippet(headers.protocol, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.protocol, row);
		},
	},
	{
		accessorKey: 'port',
		header: ({ column }) => {
			return renderSnippet(headers.port, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.port, row);
		},
	},
	// {
	// 	accessorKey: 'nodePort',
	// 	header: ({ column }) => {
	// 		return renderSnippet(headers.nodePort, column);
	// 	},
	// 	cell: ({ row }) => {
	// 		return renderSnippet(cells.nodePort, row);
	// 	},
	// },
];

export { columns, messages };
