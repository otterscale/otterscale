import type { ColumnDef } from '@tanstack/table-core';

import { cells } from './cells.svelte';
import { headers } from './headers.svelte';

import type { VirtualMachine } from '$lib/api/virtual_machine/v1/virtual_machine_pb';
import { renderSnippet } from '$lib/components/ui/data-table/index.js';
import { m } from '$lib/paraglide/messages';

const messages = {
	name: m.name(),
	status: m.status(),
	namespace: m.namespace(),
	machineId: m.machine(),
	instanceTypeName: m.instance_type(),
	instancePhase: m.instance_phase(),
	disk: m.disk(),
	port: m.ports(),
	createTime: m.create_time(),
};

const columns: ColumnDef<VirtualMachine>[] = [
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
		accessorKey: 'status',
		header: ({ column }) => {
			return renderSnippet(headers.status, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.status, row);
		},
		filterFn: 'arrIncludesSome',
	},
	{
		accessorKey: 'namespace',
		header: ({ column }) => {
			return renderSnippet(headers.namespace, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.namespace, row);
		},
		filterFn: 'arrIncludesSome',
	},
	{
		accessorKey: 'machineId',
		header: ({ column }) => {
			return renderSnippet(headers.machineId, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.machineId, row);
		},
		filterFn: 'arrIncludesSome',
	},

	{
		accessorKey: 'instanceTypeName',
		header: ({ column }) => {
			return renderSnippet(headers.instanceTypeName, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.instanceTypeName, row);
		},
		filterFn: 'arrIncludesSome',
	},
	{
		accessorKey: 'instancePhase',
		header: ({ column }) => {
			return renderSnippet(headers.instancePhase, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.instancePhase, row);
		},
		filterFn: 'arrIncludesSome',
	},
	{
		accessorKey: 'disk',
		header: ({ column }) => {
			return renderSnippet(headers.disk, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.disk, row);
		},
		filterFn: 'arrIncludesSome',
	},
	{
		accessorKey: 'port',
		header: ({ column }) => {
			return renderSnippet(headers.port, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.port, row);
		},
		filterFn: 'arrIncludesSome',
	},
	{
		accessorKey: 'createTime',
		header: ({ column }) => {
			return renderSnippet(headers.createTime, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.createTime, row);
		},
		filterFn: 'arrIncludesSome',
	},
	{
		accessorKey: 'actions',
		header: ({ column }) => {
			return renderSnippet(headers.actions, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.actions, row);
		},
		enableHiding: false,
	},
];

export { columns, messages };
