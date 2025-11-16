<script lang="ts" module>
	import {
		type ColumnFiltersState,
		getCoreRowModel,
		getFilteredRowModel,
		getPaginationRowModel,
		getSortedRowModel,
		type PaginationState,
		type RowSelectionState,
		type SortingState,
		type VisibilityState
	} from '@tanstack/table-core';
	import { type Writable } from 'svelte/store';

	import { page } from '$app/state';
	import type { Machine } from '$lib/api/machine/v1/machine_pb';
	import { Empty, Filters, Footer, Pagination } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { Reloader, ReloadManager } from '$lib/components/custom/reloader';
	import { createSvelteTable, FlexRender } from '$lib/components/ui/data-table/index.js';
	import * as Table from '$lib/components/ui/table/index.js';

	import { getColumns, messages } from './columns';
</script>

<script lang="ts">
	let { machines, reloadManager }: { machines: Writable<Machine[]>; reloadManager: ReloadManager } =
		$props();

	let pagination = $state<PaginationState>({ pageIndex: 0, pageSize: 8 });
	let sorting = $state<SortingState>([]);
	let columnFilters = $state<ColumnFiltersState>([]);
	let columnVisibility = $state<VisibilityState>({});
	let rowSelection = $state<RowSelectionState>({});

	const table = createSvelteTable({
		get data() {
			return $machines;
		},
		get columns() {
			return getColumns(reloadManager);
		},

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
	<Layout.Controller>
		<Layout.ControllerFilter>
			<Filters.StringFuzzy
				columnId="fqdn_ip"
				values={$machines.map((row) => row.fqdn)}
				{messages}
				{table}
			/>
			<Filters.StringMatch
				columnId="powerState"
				values={$machines.flatMap((row) => row.powerState)}
				{messages}
				{table}
			/>
			<Filters.StringMatch
				columnId="status"
				values={$machines.flatMap((row) => row.status)}
				{messages}
				{table}
			/>
			<Filters.Column {messages} {table} />
		</Layout.ControllerFilter>
		<Layout.ControllerAction>
			<Reloader
				bind:checked={reloadManager.state}
				onCheckedChange={() => {
					if (reloadManager.state) {
						reloadManager.restart();
					} else {
						reloadManager.stop();
					}
				}}
			/>
		</Layout.ControllerAction>
	</Layout.Controller>
	<Layout.Viewer>
		<Table.Root>
			<Table.Header>
				{#each table.getHeaderGroups() as headerGroup (headerGroup.id)}
					<Table.Row>
						{#each headerGroup.headers as header (header.id)}
							{#if header.column.id !== 'gpu' || page.data['feature-states.mdl-general']}
								<Table.Head>
									{#if !header.isPlaceholder}
										<FlexRender
											content={header.column.columnDef.header}
											context={header.getContext()}
										/>
									{/if}
								</Table.Head>
							{/if}
						{/each}
					</Table.Row>
				{/each}
			</Table.Header>
			<Table.Body>
				{#each table.getRowModel().rows as row (row.id)}
					<Table.Row data-state={row.getIsSelected() && 'selected'}>
						{#each row.getVisibleCells() as cell (cell.id)}
							{#if cell.column.id !== 'gpu' || page.data['feature-states.mdl-general']}
								<Table.Cell>
									<FlexRender content={cell.column.columnDef.cell} context={cell.getContext()} />
								</Table.Cell>
							{/if}
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
