<script lang="ts" module>
	import {
		getCoreRowModel,
		getFilteredRowModel,
		getPaginationRowModel,
		getSortedRowModel,
		type ColumnFiltersState,
		type PaginationState,
		type RowSelectionState,
		type SortingState,
		type VisibilityState,
	} from '@tanstack/table-core';
	import { type Writable } from 'svelte/store';

	import type { LargeLanguageModel } from '../type';

	import Create from './action-create.svelte';
	import { columns, messages } from './columns';
	import Statistics from './statistics.svelte';

	import { Empty, Filters, Footer, Pagination } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { Reloader, ReloadManager } from '$lib/components/custom/reloader';
	import { createSvelteTable, FlexRender } from '$lib/components/ui/data-table/index.js';
	import * as Table from '$lib/components/ui/table/index.js';
</script>

<script lang="ts">
	let {
		largeLanguageModels,
		reloadManager,
	}: { largeLanguageModels: Writable<LargeLanguageModel[]>; reloadManager: ReloadManager } = $props();

	let pagination = $state<PaginationState>({ pageIndex: 0, pageSize: 9 });
	let sorting = $state<SortingState>([]);
	let columnFilters = $state<ColumnFiltersState>([]);
	let columnVisibility = $state<VisibilityState>({});
	let rowSelection = $state<RowSelectionState>({});

	const table = createSvelteTable({
		get data() {
			return $largeLanguageModels;
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
			},
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
		autoResetPageIndex: false,
	});
</script>

<Layout.Root>
	<Layout.Statistics>
		<Statistics {table} />
	</Layout.Statistics>
	<Layout.Controller>
		<Layout.ControllerFilter>
			<Filters.StringFuzzy
				columnId="name"
				values={$largeLanguageModels.map((row) => row.name)}
				{messages}
				{table}
			/>
			<Filters.Column {table} {messages} />
		</Layout.ControllerFilter>
		<Layout.ControllerAction>
			<Create />
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
