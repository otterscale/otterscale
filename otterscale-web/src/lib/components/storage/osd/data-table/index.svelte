<script lang="ts" module>
	import type { OSD } from '$lib/api/storage/v1/storage_pb';
	import { Empty, Filters, Footer, Layout, Pagination } from '$lib/components/custom/data-table';
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
	import type { PrometheusDriver } from 'prometheus-query';
	import { getContext } from 'svelte';
	import { type Writable } from 'svelte/store';
	import { columns } from './columns';
</script>

<script lang="ts" generics="TData, TValue">
	const prometheusDriver: Writable<PrometheusDriver> = getContext('prometheusDriver');

	let {
		objectStorageDaemons
	}: {
		objectStorageDaemons: Writable<OSD[]>;
	} = $props();

	let pagination = $state<PaginationState>({ pageIndex: 0, pageSize: 10 });
	let sorting = $state<SortingState>([]);
	let columnFilters = $state<ColumnFiltersState>([]);
	let columnVisibility = $state<VisibilityState>({});
	let rowSelection = $state<RowSelectionState>({});

	const table = createSvelteTable({
		get data() {
			return $objectStorageDaemons;
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
	<Layout.Statistics>
		<!-- <Statistics {table} /> -->
	</Layout.Statistics>
	<Layout.Controller>
		<Layout.ControllerFilter>
			<Filters.StringFuzzy
				values={$objectStorageDaemons.map((row) => row.name)}
				columnId="name"
				{table}
			/>
			<Filters.BooleanMatch
				columnId="_in"
				alias="In"
				{table}
				values={$objectStorageDaemons.map((row) => row.in)}
				descriptor={(value) => (value ? 'In' : 'Out')}
			/>
			<Filters.BooleanMatch
				columnId="_up"
				alias="Up"
				{table}
				values={$objectStorageDaemons.map((row) => row.up)}
				descriptor={(value) => (value ? 'Up' : 'Down')}
			/>
			<Filters.BooleanMatch
				columnId="exists"
				{table}
				values={$objectStorageDaemons.map((row) => row.exists)}
				descriptor={(value) => (value ? 'Exists' : 'Not Exists')}
			/>
			<Filters.StringMatch
				columnId="deviceClass"
				alias="Device Class"
				{table}
				values={$objectStorageDaemons.map((row) => row.deviceClass)}
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
							{#if !header.column.columnDef.id?.startsWith('_')}
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
							{#if !cell.column.columnDef.id?.startsWith('_')}
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
