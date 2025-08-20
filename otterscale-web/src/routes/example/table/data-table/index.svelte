<script lang="ts">
	import ColumnViewer from '$lib/components/custom/data-table/data-table-filters/column.svelte';
	import TableEmpty from '$lib/components/custom/data-table/data-table-empty.svelte';
	import * as Filter from '$lib/components/custom/data-table/data-table-filters';
	import TableFooter from '$lib/components/custom/data-table/data-table-footer.svelte';
	import * as Layout from '$lib/components/custom/data-table/data-table-layout';
	import TablePagination from '$lib/components/custom/data-table/data-table-pagination.svelte';
	import Button from '$lib/components/ui/button/button.svelte';
	import { createSvelteTable, FlexRender } from '$lib/components/ui/data-table/index.js';
	import * as Table from '$lib/components/ui/table/index.js';
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
	import { writable } from 'svelte/store';
	import Actions from './actions.svelte';
	import { columns } from './columns';
	import type { TableRow } from './type';
</script>

<script lang="ts" generics="TData, TValue">
	let {
		selectedScope,
		selectedFacility,
		dataset
	}: { selectedScope: string; selectedFacility: string; dataset: TableRow[] } = $props();

	let data = $state(writable(dataset));

	let pagination = $state<PaginationState>({ pageIndex: 0, pageSize: 10 });
	let sorting = $state<SortingState>([]);
	let columnFilters = $state<ColumnFiltersState>([]);
	let columnVisibility = $state<VisibilityState>({});
	let rowSelection = $state<RowSelectionState>({});

	const table = createSvelteTable({
		get data() {
			return $data;
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
		}
	});
</script>

<Layout.Root>
	<Layout.Statistics></Layout.Statistics>
	<Layout.Controller>
		<Layout.ControllerFilter>
			<Filter.StringFuzzy columnId="name" values={$data.map((row) => row.name)} {table} />
			<Filter.StringMatch columnId="name" values={$data.map((row) => row.name)} {table} />
			<Filter.BooleanMatch
				columnId="isVerified"
				values={$data.map((row) => row.isVerified)}
				{table}
			/>
			<Filter.NumberRange columnId="id" values={$data.map((row) => row.id)} {table} />
			<ColumnViewer {table} />
		</Layout.ControllerFilter>
		<Layout.ControllerAction>
			<Button class="h-10">Create</Button>
		</Layout.ControllerAction>
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
						<Table.Head></Table.Head>
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

						<Table.Cell>
							<Actions {selectedScope} {selectedFacility} row={row.original} bind:data />
						</Table.Cell>
					</Table.Row>
				{:else}
					<Table.Row>
						<Table.Cell colspan={columns.length}>
							<TableEmpty />
						</Table.Cell>
					</Table.Row>
				{/each}
			</Table.Body>
		</Table.Root>
	</Layout.Viewer>
	<Layout.Footer>
		<TableFooter {table} />
		<TablePagination {table} />
	</Layout.Footer>
</Layout.Root>
