import type { ColumnDef } from '@tanstack/table-core';

import { cells } from './cells.svelte';
import { headers } from './headers.svelte';

import type { Application_Service } from '$lib/api/application/v1/application_pb';
import { renderSnippet } from '$lib/components/ui/data-table/index.js';
import { m } from '$lib/paraglide/messages';

const messages = {
	name: m.name(),
	type: m.type(),
	clusterIp: m.cluster_ip(),
	port: m.ports(),
	createTime: m.create_time(),
};

const columns: ColumnDef<Application_Service>[] = [
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
		accessorKey: 'type',
		header: ({ column }) => {
			return renderSnippet(headers.type, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.type, row);
		},
	},
	{
		accessorKey: 'clusterIp',
		header: ({ column }) => {
			return renderSnippet(headers.clusterIp, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.clusterIp, row);
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
	{
		accessorKey: 'createTime',
		header: ({ column }) => {
			return renderSnippet(headers.createTime, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.createTime, row);
		},
	},
];

export { columns, messages };
