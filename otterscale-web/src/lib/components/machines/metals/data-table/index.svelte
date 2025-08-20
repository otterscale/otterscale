<script lang="ts" generics="TData, TValue">
	import type { Machine } from '$lib/api/machine/v1/machine_pb';
	import {
		Empty,
		Filters,
		Footer,
		getSortingFunction,
		Layout,
		Pagination
	} from '$lib/components/custom/data-table';
	import {
		createSvelteTable,
		FlexRender,
		renderSnippet
	} from '$lib/components/ui/data-table/index.js';
	import * as Table from '$lib/components/ui/table/index.js';
	import type { ColumnDef } from '@tanstack/table-core';
	import {
		getCoreRowModel,
		getFilteredRowModel,
		getPaginationRowModel,
		getSortedRowModel,
		type ColumnFiltersState,
		type PaginationState,
		type RowSelectionState,
		type SortingState,
		type VisibilityState
	} from '@tanstack/table-core';
	import { type Writable } from 'svelte/store';
	import { cells } from './cells.svelte';
	import { headers } from './headers.svelte';

	let { machines }: { machines: Writable<Machine[]> } = $props();

	const columns: ColumnDef<Machine>[] = [
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
			sortingFn: (previousRow, nextRow, columnId) =>
				getSortingFunction(
					previousRow.original.fqdn,
					nextRow.original.fqdn,
					(p: string, n: string) => p.localeCompare(n) < 0,
					(p, n) => p === n
				),
			filterFn: (row, columnId, filterValue: string | undefined) => {
				if (filterValue === undefined) {
					return true;
				}

				return row.original.fqdn.includes(filterValue);
			}
		},
		{
			accessorKey: 'powerState',
			header: ({ column }) => {
				return renderSnippet(headers.powerState, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.powerState, row);
			},
			filterFn: 'arrIncludesSome'
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
			sortingFn: (previousRow, nextRow, columnId) =>
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
			sortingFn: (previousRow, nextRow, columnId) =>
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
			sortingFn: (previousRow, nextRow, columnId) =>
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
			sortingFn: (previousRow, nextRow, columnId) =>
				getSortingFunction(
					previousRow.original.storageMb,
					nextRow.original.storageMb,
					(p, n) => p < n,
					(p, n) => p === n
				)
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
			accessorKey: 'tags',
			header: ({ column }) => {
				return renderSnippet(headers.tags, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.tags, row);
			},
			sortingFn: (previousRow, nextRow, columnId) =>
				getSortingFunction(
					previousRow.original.tags.length,
					nextRow.original.tags.length,
					(p, n) => p < n,
					(p, n) => p === n
				)
		},
		{
			accessorKey: 'actions',
			header: ({ column }) => {
				return renderSnippet(headers.actions, column);
			},
			cell: ({ row }) => {
				return renderSnippet(cells.actions, row);
			}
		}
	];

	let pagination = $state<PaginationState>({ pageIndex: 0, pageSize: 10 });
	let sorting = $state<SortingState>([]);
	let columnFilters = $state<ColumnFiltersState>([]);
	let columnVisibility = $state<VisibilityState>({});
	let rowSelection = $state<RowSelectionState>({});
	const table = createSvelteTable({
		get data() {
			return $machines;
		},
		columns,
		getCoreRowModel: getCoreRowModel(),
		getPaginationRowModel: getPaginationRowModel(),
		getSortedRowModel: getSortedRowModel(),
		getFilteredRowModel: getFilteredRowModel(),
		state: {
			get pagination() {
				return pagination;
			},
			get sorting() {
				return sorting;
			},
			get columnFilters() {
				return columnFilters;
			},
			get columnVisibility() {
				return columnVisibility;
			},
			get rowSelection() {
				return rowSelection;
			}
		},
		onPaginationChange: (updater) => {
			if (typeof updater === 'function') {
				pagination = updater(pagination);
			} else {
				pagination = updater;
			}
		},
		onSortingChange: (updater) => {
			if (typeof updater === 'function') {
				sorting = updater(sorting);
			} else {
				sorting = updater;
			}
		},
		onColumnFiltersChange: (updater) => {
			if (typeof updater === 'function') {
				columnFilters = updater(columnFilters);
			} else {
				columnFilters = updater;
			}
		},
		onColumnVisibilityChange: (updater) => {
			if (typeof updater === 'function') {
				columnVisibility = updater(columnVisibility);
			} else {
				columnVisibility = updater;
			}
		},
		onRowSelectionChange: (updater) => {
			if (typeof updater === 'function') {
				rowSelection = updater(rowSelection);
			} else {
				rowSelection = updater;
			}
		},
		autoResetPageIndex: false
	});
</script>

<Layout.Root>
	<Layout.Statistics></Layout.Statistics>
	<Layout.Controller>
		<Layout.ControllerFilter>
			<Filters.StringFuzzy
				values={$machines.map((row) => row.fqdn)}
				columnId="fqdn_ip"
				alias="FQDN"
				{table}
			/>
			<Filters.StringMatch
				values={$machines.flatMap((row) => row.powerState)}
				columnId="powerState"
				alias="Power"
				{table}
			/>
			<Filters.StringMatch
				values={$machines.flatMap((row) => row.status)}
				columnId="status"
				{table}
			/>
			<Filters.Column {table} />
		</Layout.ControllerFilter>
	</Layout.Controller>
	<Layout.Viewer>
		<Table.Root>
			<Table.Header>
				{#each table.getHeaderGroups() as headerGroup (headerGroup.id)}
					<Table.Row>
						{#each headerGroup.headers as header (header.id)}
							<Table.Head>
								{#if !header.isPlaceholder}
									<FlexRender
										content={header.column.columnDef.header}
										context={header.getContext()}
									/>
								{/if}
							</Table.Head>
						{/each}
					</Table.Row>
				{/each}
			</Table.Header>
			<Table.Body>
				{#each table.getRowModel().rows as row (row.id)}
					<Table.Row data-state={row.getIsSelected() && 'selected'}>
						{#each row.getVisibleCells() as cell (cell.id)}
							<Table.Cell>
								<FlexRender content={cell.column.columnDef.cell} context={cell.getContext()} />
							</Table.Cell>
						{/each}
					</Table.Row>
				{:else}
					<Empty {table} />
				{/each}
			</Table.Body>
		</Table.Root>
	</Layout.Viewer>
	<Layout.Footer>
		<Footer {table} />
		<Pagination {table} />
	</Layout.Footer>
</Layout.Root>
