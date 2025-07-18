<script lang="ts" module>
	import type { MON } from '$gen/api/storage/v1/storage_pb';
	import ColumnViewer from '$lib/components/custom/data-table/data-table-column-viewer.svelte';
	import TableEmpty from '$lib/components/custom/data-table/data-table-empty.svelte';
	import FuzzyFilter from '$lib/components/custom/data-table/data-table-filters/fuzzy-filter.svelte';
	import BooleanyFilter from '$lib/components/custom/data-table/data-table-filters/boolean-filter.svelte';
	import TableFooter from '$lib/components/custom/data-table/data-table-footer.svelte';
	import TablePagination from '$lib/components/custom/data-table/data-table-pagination.svelte';
	import * as Layout from '$lib/components/custom/data-table/layout';
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
	import { columns } from './columns';
	import Statistics from './statistics.svelte';
</script>

<script lang="ts" generics="TData, TValue">
	let { monitors }: { monitors: MON[] } = $props();

	let data = $state(writable(monitors));

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
	<Layout.Statistics>
		<Statistics {table} />
	</Layout.Statistics>
	<Layout.Controller>
		<Layout.ControllerFilter>
			<FuzzyFilter columnId="name" {table} />
			<FuzzyFilter columnId="publicAddress" alias="Address" {table} />
			<BooleanyFilter
				{table}
				values={$data.map((row) => row.leader)}
				columnId="leader"
				descriptor={(value) => (value ? 'Leading' : 'Not Leading')}
			/>
			<ColumnViewer {table} />
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
