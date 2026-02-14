import type { ColumnDef } from '@tanstack/table-core';

import type { Machine } from '$lib/api/machine/v1/machine_pb';
import { getSortingFunction } from '$lib/components/custom/data-table/core';
import type { ReloadManager } from '$lib/components/custom/reloader';
import { renderSnippet } from '$lib/components/ui/data-table/index.js';
import { m } from '$lib/paraglide/messages';

import type { Metrics } from '../types';
import { cells } from './cells.svelte';
import { headers } from './headers.svelte';

const messages = {
	fqdn_ip: m.fqdn(),
	power_state: m.power(),
	status: m.status(),
	cores_arch: m.core(),
	ram: m.ram(),
	disk: m.disk(),
	storage: m.storage(),
	gpu: m.gpu(),
	scope: m.scope(),
	memory_metric: m.memory(),
	cpu_metric: m.cpu(),
	storage_metric: m.storage(),
	tags: m.tags()
};

function getColumns(metrics: Metrics, reloadManager: ReloadManager): ColumnDef<Machine>[] {
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
			accessorKey: 'fqdn_ip',
			header: ({ column }) => {
				return renderSnippet(headers.fqdn_ip, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.fqdn_ip, row);
			},
			sortingFn: (previousRow, nextRow) =>
				getSortingFunction(
					previousRow.original.fqdn,
					nextRow.original.fqdn,
					(p: string, n: string) => p.localeCompare(n) < 0,
					(p, n) => p === n
				),
			filterFn: (row, filterValue: string | undefined) => {
				if (filterValue === undefined) {
					return true;
				}

				return row.original.fqdn.includes(filterValue);
			}
		},
		{
			accessorKey: 'scope',
			header: ({ column }) => {
				return renderSnippet(headers.scope, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.scope, row);
			}
		},
		{
			accessorKey: 'power_state',
			header: ({ column }) => {
				return renderSnippet(headers.power_state, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.power_state, row);
			},
			filterFn: (row, columnId, filterValue) => {
				if (!filterValue) {
					return true;
				}

				return row.original.powerState ? filterValue.includes(row.original.powerState) : false;
			}
		},
		{
			accessorKey: 'status',
			header: ({ column }) => {
				return renderSnippet(headers.status, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.status, row);
			},
			filterFn: 'arrIncludesSome'
		},
		{
			accessorKey: 'cores_arch',
			header: ({ column }) => {
				return renderSnippet(headers.cores_arch, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.cores_arch, row);
			},
			sortingFn: (previousRow, nextRow) =>
				getSortingFunction(
					previousRow.original.cpuCount,
					nextRow.original.cpuCount,
					(p, n) => p < n,
					(p, n) => p === n
				)
		},
		{
			accessorKey: 'ram',
			header: ({ column }) => {
				return renderSnippet(headers.ram, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.ram, row);
			},
			sortingFn: (previousRow, nextRow) =>
				getSortingFunction(
					previousRow.original.memoryMb,
					nextRow.original.memoryMb,
					(p, n) => p < n,
					(p, n) => p === n
				)
		},
		{
			accessorKey: 'disk',
			header: ({ column }) => {
				return renderSnippet(headers.disk, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.disk, row);
			},
			sortingFn: (previousRow, nextRow) =>
				getSortingFunction(
					previousRow.original.blockDevices.length,
					nextRow.original.blockDevices.length,
					(p, n) => p < n,
					(p, n) => p === n
				)
		},
		{
			accessorKey: 'storage',
			header: ({ column }) => {
				return renderSnippet(headers.storage, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.storage, row);
			},
			sortingFn: (previousRow, nextRow) =>
				getSortingFunction(
					previousRow.original.storageMb,
					nextRow.original.storageMb,
					(p, n) => p < n,
					(p, n) => p === n
				)
		},
		{
			accessorKey: 'gpu',
			header: ({ column }) => {
				return renderSnippet(headers.gpu, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.gpu, row);
			}
		},

		{
			accessorKey: 'tags',
			header: ({ column }) => {
				return renderSnippet(headers.tags, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.tags, { row, reloadManager });
			},
			sortingFn: (previousRow, nextRow) =>
				getSortingFunction(
					previousRow.original.tags.length,
					nextRow.original.tags.length,
					(p, n) => p < n,
					(p, n) => p === n
				)
		},
		{
			accessorKey: 'cpu_metric',
			header: ({ column }) => {
				return renderSnippet(headers.cpu_metric, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.cpu_metric, { row, metrics });
			}
		},
		{
			accessorKey: 'memory_metric',
			header: ({ column }) => {
				return renderSnippet(headers.memory_metric, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.memory_metric, { row, metrics });
			}
		},
		{
			accessorKey: 'storage_metric',
			header: ({ column }) => {
				return renderSnippet(headers.storage_metric, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.storage_metric, { row, metrics });
			}
		},
		{
			accessorKey: 'actions',
			header: ({ column }) => {
				return renderSnippet(headers.actions, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.actions, { row, reloadManager });
			},
			enableHiding: false
		}
	];
}

export { getColumns, messages };
