<script lang="ts">
	import { type JsonObject, type JsonValue } from '@bufbuild/protobuf';
	import Binary from '@lucide/svelte/icons/binary';
	import Braces from '@lucide/svelte/icons/braces';
	import ChevronDown from '@lucide/svelte/icons/chevron-down';
	import ChevronFirst from '@lucide/svelte/icons/chevron-first';
	import ChevronLast from '@lucide/svelte/icons/chevron-last';
	import ChevronLeft from '@lucide/svelte/icons/chevron-left';
	import ChevronRight from '@lucide/svelte/icons/chevron-right';
	import ChevronUp from '@lucide/svelte/icons/chevron-up';
	import Clock from '@lucide/svelte/icons/clock';
	import Columns3 from '@lucide/svelte/icons/columns-3';
	import Eraser from '@lucide/svelte/icons/eraser';
	import Hash from '@lucide/svelte/icons/hash';
	import Type from '@lucide/svelte/icons/type';
	import {
		type ColumnDef,
		type ColumnFiltersState,
		getCoreRowModel,
		getFacetedUniqueValues,
		getFilteredRowModel,
		getPaginationRowModel,
		getSortedRowModel,
		type PaginationState,
		type Row,
		type RowSelectionState,
		type SortingState,
		type Table as TanStackTabke,
		type VisibilityState
	} from '@tanstack/table-core';
	import jsep from 'jsep';
	import lodash from 'lodash';
	import { createRawSnippet, type Snippet } from 'svelte';

	import { Button } from '$lib/components/ui/button/index.js';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import {
		createSvelteTable,
		FlexRender,
		renderComponent,
		renderSnippet
	} from '$lib/components/ui/data-table';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import * as Empty from '$lib/components/ui/empty/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import * as Pagination from '$lib/components/ui/pagination/index.js';
	import { Select, SelectContent, SelectItem, SelectTrigger } from '$lib/components/ui/select';
	import * as Table from '$lib/components/ui/table';
	import { cn } from '$lib/utils';

	import DynamicalTableQuery, { evaluate } from './dynamical-table-query.svelte';

	let {
		objects,
		fields,
		columnDefinitions,
		create,
		bulkDelete,
		reload,
		rowActions = createRawSnippet(() => ({
			render: () => ''
		}))
	}: {
		objects: Record<string, JsonValue>[];
		fields: Record<string, JsonValue>;
		columnDefinitions: ColumnDef<Record<string, JsonValue>>[];
		create?: Snippet;
		bulkDelete?: Snippet<[{ table: TanStackTabke<Record<string, JsonValue>> }]>;
		rowActions?: Snippet<[{ row: Row<Record<string, JsonValue>> }]>;
		reload: Snippet;
	} = $props();

	const columns: ColumnDef<Record<string, JsonValue>>[] = [
		{
			id: 'select',
			header: ({ table }) =>
				renderComponent(Checkbox, {
					class: 'm-1',
					'aria-label': 'Select all',
					checked: table.getIsAllPageRowsSelected(),
					indeterminate: table.getIsSomePageRowsSelected() && !table.getIsAllPageRowsSelected(),
					onCheckedChange: (value) => table.toggleAllPageRowsSelected(!!value)
				}),
			cell: ({ row }) =>
				renderComponent(Checkbox, {
					class: 'm-1',
					'aria-label': 'Select row',
					checked: row.getIsSelected(),
					onCheckedChange: (value) => row.toggleSelected(!!value)
				}),
			enableHiding: false,
			enableSorting: false,
			size: 30
		},
		...columnDefinitions,
		{
			id: 'actions',
			cell: ({ row }) => renderSnippet(rowActions, { row: row }),
			header: () =>
				renderSnippet(
					createRawSnippet(() => {
						return {
							render: () => `<span class="sr-only">Actions</span>`
						};
					}),
					{}
				),
			enableHiding: false,
			enableSorting: false,
			size: 40
		}
	];

	let globalFilter = $state('');

	let rowSelection = $state<RowSelectionState>({});
	let columnFilters = $state<ColumnFiltersState>([]);
	let columnVisibility = $state<VisibilityState>({});
	let sorting = $state<SortingState>([]);
	let pagination = $state<PaginationState>({ pageIndex: 0, pageSize: 10 });

	let table = createSvelteTable<JsonObject>({
		columns,
		get data() {
			return objects;
		},
		getCoreRowModel: getCoreRowModel(),
		getFacetedUniqueValues: getFacetedUniqueValues(),
		getFilteredRowModel: getFilteredRowModel(),
		getPaginationRowModel: getPaginationRowModel(),
		getSortedRowModel: getSortedRowModel(),
		onColumnFiltersChange: (updater) => {
			if (typeof updater === 'function') {
				columnFilters = updater(columnFilters);
			} else {
				columnFilters = updater;
			}
		},
		onGlobalFilterChange: (updater) => {
			if (typeof updater === 'function') {
				globalFilter = updater(globalFilter);
			} else {
				globalFilter = updater;
			}
		},
		onColumnVisibilityChange: (updater) => {
			if (typeof updater === 'function') {
				columnVisibility = updater(columnVisibility);
			} else {
				columnVisibility = updater;
			}
		},
		onPaginationChange: (updater) => {
			if (typeof updater === 'function') {
				pagination = updater(pagination);
			} else {
				pagination = updater;
			}
		},
		onRowSelectionChange: (updater) => {
			if (typeof updater === 'function') {
				rowSelection = updater(rowSelection);
			} else {
				rowSelection = updater;
			}
		},
		onSortingChange: (updater) => {
			if (typeof updater === 'function') {
				sorting = updater(sorting);
			} else {
				sorting = updater;
			}
		},
		state: {
			get globalFilter() {
				return globalFilter;
			},
			get columnFilters() {
				return columnFilters;
			},
			get columnVisibility() {
				return columnVisibility;
			},
			get pagination() {
				return pagination;
			},
			get rowSelection() {
				return rowSelection;
			},
			get sorting() {
				return sorting;
			}
		},
		globalFilterFn: (row, _, filterValue: string) => {
			if (!filterValue) return true;
			try {
				const ast = jsep(filterValue);
				return evaluate(ast, row.original);
			} catch (error) {
				console.error('Parse error:', error);
				return true;
			}
		}
	});

	function handleResetFilter() {
		expression = '';
	}

	// eslint-disable-next-line
	function getAlignment(field: any): 'start' | 'center' | 'end' {
		if (
			field?.type === 'integer' ||
			field?.type === 'number' ||
			(field?.type === 'string' && field?.format === 'date-time')
		) {
			return 'end';
		} else if (field?.type === 'boolean' || field?.type === 'object') {
			return 'center';
		} else {
			return 'start';
		}
	}
	// eslint-disable-next-line
	function getHeaderAlignment(field: any): string {
		const alignment = getAlignment(field);
		switch (alignment) {
			case 'start':
				return 'justify-start';
			case 'center':
				return 'justify-center';
			case 'end':
			default:
				return 'justify-end';
		}
	}
	// eslint-disable-next-line
	function getCellAlignment(field: any): string {
		const alignment = getAlignment(field);
		switch (alignment) {
			case 'start':
				return 'text-start';
			case 'center':
				return 'text-center';
			case 'end':
			default:
				return 'text-end';
		}
	}

	let expression = $state('');
</script>

<div class="space-y-4">
	<!-- Controllers -->
	<!-- Accessors -->
	<div class="flex w-full items-center gap-4">
		<div>
			{@render create?.()}
			{@render bulkDelete?.({ table })}
		</div>
		<div class="ml-auto">
			{@render reload()}
		</div>
	</div>
	<!-- Filters -->
	<div class="flex w-full items-center gap-4">
		<DynamicalTableQuery bind:expression {table} />
		<DropdownMenu.Root>
			<DropdownMenu.Trigger>
				{#snippet child({ props })}
					<Button variant="outline" {...props}>
						<Columns3 class="-ms-1 opacity-60" size={16} aria-hidden="true" />
						View
					</Button>
				{/snippet}
			</DropdownMenu.Trigger>
			<DropdownMenu.Content align="end">
				<DropdownMenu.Label>Toggle columns</DropdownMenu.Label>
				{#each table.getAllColumns().filter((column) => column.getCanHide()) as column (column.id)}
					<DropdownMenu.CheckboxItem
						checked={column.getIsVisible()}
						closeOnSelect={false}
						onCheckedChange={(value) => column.toggleVisibility(!!value)}
					>
						{@const type = lodash.get(fields, `${column.id}.type`)}
						{@const format = lodash.get(fields, `${column.id}.format`)}
						{#if type === 'boolean'}
							<Binary />
						{:else if type === 'number' || type === 'integer'}
							<Hash />
						{:else if type === 'string' && (format === 'date' || format === 'date-time')}
							<Clock />
						{:else if type === 'string'}
							<Type />
						{:else if type === 'array'}
							<Braces />
						{:else if type === 'object'}
							<Braces />
						{/if}
						{column.id}
					</DropdownMenu.CheckboxItem>
				{/each}
			</DropdownMenu.Content>
		</DropdownMenu.Root>
	</div>
	<!-- Table -->
	<div class="overflow-hidden rounded-md border bg-background">
		<Table.Root class="table-fixed">
			<Table.Header>
				{#each table.getHeaderGroups() as headerGroup (headerGroup.id)}
					<Table.Row class="hover:bg-transparent">
						{#each headerGroup.headers as header (header.id)}
							<Table.Head style="width: {header.getSize()}px" class="h-11">
								{#if !header.isPlaceholder && header.column.getCanSort()}
									<div
										class={cn(
											header.column.getCanSort() &&
												'flex h-full cursor-pointer items-center justify-between gap-2 select-none',
											getHeaderAlignment(fields[header.column.id])
										)}
										onclick={header.column.getToggleSortingHandler()}
										onkeydown={(e) => {
											if (header.column.getCanSort() && (e.key === 'Enter' || e.key === ' ')) {
												e.preventDefault();
												header.column.getToggleSortingHandler()?.(e);
											}
										}}
										{...header.column.getCanSort()
											? {
													tabindex: 0,
													role: 'button',
													'aria-pressed': header.column.getIsSorted() ? 'true' : 'false'
												}
											: {}}
									>
										<FlexRender
											content={header.column.columnDef.header}
											context={header.getContext()}
										/>
										{#if header.column.getIsSorted() === 'asc'}
											<ChevronUp class="shrink-0 opacity-60" size={16} aria-hidden="true" />
										{:else if header.column.getIsSorted() === 'desc'}
											<ChevronDown class="shrink-0 opacity-60" size={16} aria-hidden="true" />
										{/if}
									</div>
								{:else if !header.isPlaceholder && !header.column.getCanSort()}
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
				{#if table.getRowModel().rows?.length}
					{#each table.getRowModel().rows as row (row.id)}
						<Table.Row data-state={row.getIsSelected() && 'selected'}>
							{#each row.getVisibleCells() as cell (cell.id)}
								<Table.Cell class={getCellAlignment(fields[cell.column.id])}>
									<FlexRender content={cell.column.columnDef.cell} context={cell.getContext()} />
								</Table.Cell>
							{/each}
						</Table.Row>
					{/each}
				{:else}
					<Table.Row>
						<Table.Cell colspan={columns.length} class="h-full text-center">
							<Empty.Root>
								<Empty.Header>
									<Empty.Media variant="icon">
										<Columns3 size={32} class="opacity-60" aria-hidden="true" />
									</Empty.Media>
									<Empty.Title>No Resources Found</Empty.Title>
									<Empty.Description>
										No resources found. Please adjust your filters or initiate a new resource to
										populate this table.
									</Empty.Description>
								</Empty.Header>
								<Empty.Content>
									<Button onclick={handleResetFilter}>
										<Eraser size={16} class="opacity-60" />
										Reset
									</Button>
								</Empty.Content>
							</Empty.Root>
						</Table.Cell>
					</Table.Row>
				{/if}
			</Table.Body>
		</Table.Root>
	</div>

	<!-- Pagination -->
	<div class="flex items-center justify-between gap-8">
		<!-- Results -->
		<div class="flex items-center gap-3">
			<Label class="max-sm:sr-only">Rows per page</Label>
			<Select
				type="single"
				value={table.getState().pagination.pageSize.toString()}
				onValueChange={(value) => {
					table.setPageSize(Number(value));
				}}
			>
				<SelectTrigger class="w-fit whitespace-nowrap">
					{table.getState().pagination.pageSize.toString() ?? 'Select number of results'}
				</SelectTrigger>
				<SelectContent
					class="[&_*[role=option]]:ps-2 [&_*[role=option]]:pe-8 [&_*[role=option]>span]:start-auto [&_*[role=option]>span]:end-2"
				>
					{#each [5, 10, 25, 50] as pageSize (pageSize)}
						<SelectItem value={pageSize.toString()}>
							{pageSize}
						</SelectItem>
					{/each}
				</SelectContent>
			</Select>
		</div>

		<!-- Page -->
		<div class="flex grow justify-end text-sm whitespace-nowrap text-muted-foreground">
			<p class="text-sm whitespace-nowrap text-muted-foreground" aria-live="polite">
				<span class="text-foreground">
					{table.getState().pagination.pageIndex * table.getState().pagination.pageSize +
						1}-{Math.min(
						Math.max(
							table.getState().pagination.pageIndex * table.getState().pagination.pageSize +
								table.getState().pagination.pageSize,
							0
						),
						table.getRowCount()
					)}
				</span>
				of
				<span class="text-foreground">
					{table.getRowCount().toString()}
				</span>
			</p>
		</div>

		<!-- Controller -->
		<div>
			<Pagination.Root count={table.getRowCount()}>
				<Pagination.Content>
					<!-- First page button -->
					<Pagination.Item>
						<Button
							size="icon"
							variant="outline"
							class="disabled:pointer-events-none disabled:opacity-50"
							onclick={() => table.firstPage()}
							disabled={!table.getCanPreviousPage()}
							aria-label="Go to first page"
						>
							<ChevronFirst size={16} aria-hidden="true" />
						</Button>
					</Pagination.Item>
					<!-- Previous page button -->
					<Pagination.Item>
						<Button
							size="icon"
							variant="outline"
							class="disabled:pointer-events-none disabled:opacity-50"
							onclick={() => table.previousPage()}
							disabled={!table.getCanPreviousPage()}
							aria-label="Go to previous page"
						>
							<ChevronLeft size={16} aria-hidden="true" />
						</Button>
					</Pagination.Item>
					<!-- Next page button -->
					<Pagination.Item>
						<Button
							size="icon"
							variant="outline"
							class="disabled:pointer-events-none disabled:opacity-50"
							onclick={() => table.nextPage()}
							disabled={!table.getCanNextPage()}
							aria-label="Go to next page"
						>
							<ChevronRight size={16} aria-hidden="true" />
						</Button>
					</Pagination.Item>
					<!-- Last page button -->
					<Pagination.Item>
						<Button
							size="icon"
							variant="outline"
							class="disabled:pointer-events-none disabled:opacity-50"
							onclick={() => table.lastPage()}
							disabled={!table.getCanNextPage()}
							aria-label="Go to last page"
						>
							<ChevronLast size={16} aria-hidden="true" />
						</Button>
					</Pagination.Item>
				</Pagination.Content>
			</Pagination.Root>
		</div>
	</div>
</div>
