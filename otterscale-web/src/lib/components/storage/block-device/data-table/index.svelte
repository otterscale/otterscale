<script lang="ts" module>
	import type { Image } from '$lib/api/storage/v1/storage_pb';
	import ColumnViewer from '$lib/components/custom/data-table/data-table-column-viewer.svelte';
	import TableEmpty from '$lib/components/custom/data-table/data-table-empty.svelte';
	import * as Filters from '$lib/components/custom/data-table/data-table-filters';
	import TableFooter from '$lib/components/custom/data-table/data-table-footer.svelte';
	import * as Layout from '$lib/components/custom/data-table/data-table-layout';
	import TablePagination from '$lib/components/custom/data-table/data-table-pagination.svelte';
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
	import { type Writable } from 'svelte/store';
	import Actions from './actions.svelte';
	import { columns } from './columns';
	import Create from './create.svelte';
	import { headers } from './headers.svelte';
	import { Snapshot } from './snapshot';
	import Statistics from './statistics.svelte';
</script>

<script lang="ts" generics="TData, TValue">
	let {
		selectedScopeUuid,
		selectedFacility,
		images
	}: { selectedScopeUuid: string; selectedFacility: string; images: Writable<Image[]> } = $props();

	let pagination = $state<PaginationState>({ pageIndex: 0, pageSize: 10 });
	let sorting = $state<SortingState>([]);
	let columnFilters = $state<ColumnFiltersState>([]);
	let columnVisibility = $state<VisibilityState>({});
	let rowSelection = $state<RowSelectionState>({});

	export const table = createSvelteTable({
		get data() {
			return $images;
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
		autoResetAll: false
	});
</script>

<Layout.Root>
	<Layout.Statistics>
		<Statistics {table} />
	</Layout.Statistics>
	<Layout.Controller>
		<Layout.ControllerFilter>
			<Filters.StringFuzzy values={$images.map((row) => row.name)} columnId="name" {table} />
			<Filters.StringMatch
				columnId="poolName"
				alias="Pool Name"
				{table}
				values={$images.map((row) => row.poolName)}
			/>
			<ColumnViewer {table} />
		</Layout.ControllerFilter>
		<Layout.ControllerAction>
			<Create {selectedScopeUuid} {selectedFacility} bind:images />
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
						<Table.Head>
							{@render headers.snapshots()}
						</Table.Head>
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
							<Snapshot {selectedScopeUuid} image={row.original} bind:images />
						</Table.Cell>
						<Table.Cell>
							<Actions {selectedScopeUuid} {selectedFacility} {row} bind:images />
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
