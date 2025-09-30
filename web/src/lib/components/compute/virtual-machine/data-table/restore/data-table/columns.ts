import type { ColumnDef } from '@tanstack/table-core';

import { cells } from './cells.svelte';
import { headers } from './headers.svelte';

import type { VirtualMachine_Restore } from '$lib/api/virtual_machine/v1/virtual_machine_pb';
import { renderSnippet } from '$lib/components/ui/data-table/index.js';
import { m } from '$lib/paraglide/messages';

const messages = {
	name: m.name(),
	namespace: m.namespace(),
	targetName: m.target_name(),
	snapshotName: m.snapshot_name(),
	complete: m.complete(),
	createTime: m.create_time(),
};

const columns: ColumnDef<VirtualMachine_Restore>[] = [
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
		accessorKey: 'namespace',
		header: ({ column }) => {
			return renderSnippet(headers.namespace, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.namespace, row);
		},
	},
	{
		accessorKey: 'targetName',
		header: ({ column }) => {
			return renderSnippet(headers.targetName, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.targetName, row);
		},
	},
	{
		accessorKey: 'snapshotName',
		header: ({ column }) => {
			return renderSnippet(headers.snapshotName, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.snapshotName, row);
		},
	},
	{
		accessorKey: 'complete',
		header: ({ column }) => {
			return renderSnippet(headers.complete, column);
		},
		cell: ({ row }) => {
			return renderSnippet(cells.complete, row);
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
