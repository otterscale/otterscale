<script lang="ts">
	import {
		BinaryIcon,
		BookIcon,
		BracesIcon,
		BracketsIcon,
		ChevronDownIcon,
		ChevronFirstIcon,
		ChevronLastIcon,
		ChevronLeftIcon,
		ChevronRightIcon,
		ChevronUpIcon,
		ClockIcon,
		CodeIcon,
		Columns3Icon,
		EraserIcon,
		HashIcon,
		PercentIcon,
		TypeIcon
	} from '@lucide/svelte';
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
	import { compileExpression } from 'filtrex';
	import lodash from 'lodash';
	import { createRawSnippet, type Snippet } from 'svelte';

	import { shortcut } from '$lib/actions/shortcut.svelte';
	import { Button, buttonVariants } from '$lib/components/ui/button/index.js';
	import * as ButtonGroup from '$lib/components/ui/button-group/index.js';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import {
		createSvelteTable,
		FlexRender,
		renderComponent,
		renderSnippet
	} from '$lib/components/ui/data-table';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import * as Empty from '$lib/components/ui/empty/index.js';
	import * as InputGroup from '$lib/components/ui/input-group/index.js';
	import * as Kbd from '$lib/components/ui/kbd/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import * as Pagination from '$lib/components/ui/pagination/index.js';
	import { Select, SelectContent, SelectItem, SelectTrigger } from '$lib/components/ui/select';
	import * as Sheet from '$lib/components/ui/sheet';
	import * as Table from '$lib/components/ui/table';
	import { cn } from '$lib/utils';

	import { getColumnType } from './utils';
	import type { FieldsType, ValuesType } from '../kind-viewer/type';

	let {
		dataset,
		fields,
		columnDefinitions,
		create,
		bulkDelete,
		reload,
		rowActions = createRawSnippet(() => ({
			render: () => ''
		}))
	}: {
		dataset: ValuesType[];
		fields: FieldsType;
		columnDefinitions: ColumnDef<ValuesType>[];
		create?: Snippet;
		bulkDelete?: Snippet<[{ table: TanStackTabke<ValuesType> }]>;
		rowActions?: Snippet<[{ row: Row<ValuesType>; fields: FieldsType; dataset: ValuesType[] }]>;
		reload: Snippet;
	} = $props();

	const columns: ColumnDef<ValuesType>[] = [
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
			cell: ({ row }) =>
				renderSnippet(rowActions, {
					row: row,
					dataset: dataset,
					fields: fields
				}),
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
	let globalFilterInput = $state('');
	let globalFilterError: Error | null = $state(null);

	const extraFunctions = {
		Time: (time: string | number | Date) => new Date(time).getTime(),

		now: () => Date.now(),

		Seconds: (time: number) => time * 1000,
		Minutes: (time: number) => time * 60 * 1000,
		Hours: (time: number) => time * 60 * 60 * 1000,
		Days: (time: number) => time * 24 * 60 * 60 * 1000,
		Years: (time: number) => time * 365 * 24 * 60 * 60 * 1000,

		length: (array: unknown[]) => (array ? array.length : 0)
	};

	const constants = {
		true: true,
		false: false
	};

	let rowSelection = $state<RowSelectionState>({});
	let columnFilters = $state<ColumnFiltersState>([]);
	let columnVisibility = $state<VisibilityState>({});
	let sorting = $state<SortingState>([]);
	let pagination = $state<PaginationState>({ pageIndex: 0, pageSize: 10 });

	let table = createSvelteTable<ValuesType>({
		columns,
		get data() {
			return dataset;
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
		globalFilterFn: (row) => {
			if (!globalFilter) return true;
			try {
				const expression = compileExpression(globalFilter, { extraFunctions, constants });
				const result = expression(row.original);
				return Boolean(result);
			} catch {
				return false;
			}
		}
	});

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

	function handleKeyDown(event: KeyboardEvent) {
		if ((event.ctrlKey || event.metaKey) && event.key === 'Enter') {
			event.preventDefault();
			handleSearch();
		}

		if (event.key === 'Escape') {
			event.preventDefault();
			handleClear();
		}
	}
	function handleSearch() {
		try {
			globalFilterError = null;
			if (globalFilterInput) {
				compileExpression(globalFilterInput, { extraFunctions, constants });
			}
			globalFilter = globalFilterInput;
			table.setGlobalFilter(globalFilterInput);
		} catch (error) {
			globalFilterError = error as Error;
		}
	}
	function handleClear() {
		globalFilterInput = '';
		globalFilter = '';
		globalFilterError = null;
		table.setGlobalFilter('');
	}
</script>

<svelte:window
	use:shortcut={{
		key: '/',
		ctrl: false,
		callback: () => {
			const input = document.getElementById('global_filter');
			if (input) (input as HTMLInputElement).focus();
		}
	}}
/>
<div class="space-y-4">
	<!-- Controllers -->
	<div class="flex w-full items-center gap-2">
		<!-- Filters -->
		<ButtonGroup.Root class="w-full">
			<InputGroup.Root>
				<InputGroup.Addon>
					<CodeIcon size={16} />
				</InputGroup.Addon>
				<InputGroup.Input
					id="global_filter"
					placeholder="Search via Query Language"
					bind:value={globalFilterInput}
					class="peer w-full"
					onkeydown={handleKeyDown}
				/>
				<InputGroup.Addon align="inline-end" class="hidden peer-focus:flex">
					<Kbd.Group>
						<Kbd.Root>ctrl</Kbd.Root>
						<Kbd.Root>⏎</Kbd.Root>
					</Kbd.Group>
				</InputGroup.Addon>
				<InputGroup.Addon align="inline-end" class="peer-focus:hidden">
					<Kbd.Group>
						<Kbd.Root>/</Kbd.Root>
					</Kbd.Group>
				</InputGroup.Addon>
			</InputGroup.Root>
			<Sheet.Root>
				<Sheet.Trigger
					aria-label="Document"
					class={buttonVariants({ variant: 'outline', size: 'icon' })}
				>
					<BookIcon size={16} />
				</Sheet.Trigger>
				<Sheet.Content side="right" class="min-w-[23vw]">
					<Sheet.Header>
						<Sheet.Title>Filtrex Query Language Documentation</Sheet.Title>
						<Sheet.Description>
							Filtrex is a simple, safe, and powerful expression language for filtering and
							searching data. You can use Filtrex queries in the search box to filter table rows
							using custom logic.
						</Sheet.Description>
					</Sheet.Header>
					<div class="overflow-auto p-4 text-sm">
						<h3 class="font-semibold">Basic Syntax</h3>
						<div class="p-4 font-mono">
							<p>Field comparison: age &gt; 18</p>
							<p>String matching: name == "Alice"</p>
							<p>Logical operators: status == "active" and score &gt; 80</p>
							<p>Field names with spaces: 'full name' == "Alice Smith"</p>
						</div>

						<br />

						<h3 class="font-semibold">Operators</h3>
						<div class="p-4 font-mono">
							<p>and, or, not</p>
							<p>==, ~=, &gt;, &gt;=, &lt;, &lt;=</p>
						</div>

						<br />

						<h3 class="font-semibold">Functions</h3>
						<div class="p-4 font-mono">
							<p>length(array) — Get array length</p>
							<p>Time(date) — Convert date to timestamp</p>
							<p>now() — Current timestamp</p>
							<p>Seconds(n), Minutes(n), Hours(n), Days(n), Years(n)</p>
						</div>

						<br />

						<h3 class="font-semibold">Constants</h3>
						<div class="p-4 font-mono">
							<p>true, false</p>
						</div>

						<br />

						<h3 class="font-semibold">Tips</h3>
						<div class="p-4 font-mono">
							<p>Use <kbd>ctrl</kbd> + <kbd>⏎</kbd> to search.</p>
							<p>Press <kbd>/</kbd> to focus the search box.</p>
							<p>Press <kbd>esc</kbd> to clear the filter.</p>
						</div>
					</div>
					<Sheet.Footer>
						<p class="mt-4 text-xs text-muted-foreground">
							For advanced usage, refer to the <a
								href="https://github.com/joewalnes/filtrex"
								target="_blank"
								class="underline">Filtrex documentation</a
							>.
						</p>
					</Sheet.Footer>
				</Sheet.Content>
			</Sheet.Root>
		</ButtonGroup.Root>
		<DropdownMenu.Root>
			<DropdownMenu.Trigger>
				{#snippet child({ props })}
					<Button variant="outline" {...props}>
						<Columns3Icon class="-ms-1 opacity-60" size={16} aria-hidden="true" />
					</Button>
				{/snippet}
			</DropdownMenu.Trigger>
			<DropdownMenu.Content align="end">
				<DropdownMenu.Label>Toggle columns</DropdownMenu.Label>
				{#each table.getAllColumns().filter((column) => column.getCanHide()) as column (column.id)}
					<DropdownMenu.Item
						class={column.getIsVisible()
							? 'text-primary **:text-primary'
							: 'text-muted-foreground/50 **:text-muted-foreground/50'}
						closeOnSelect={false}
						onSelect={() => column.toggleVisibility(!column.getIsVisible())}
					>
						{@const type = lodash.get(fields, `${column.id}.type`)}
						{@const format = lodash.get(fields, `${column.id}.format`)}
						{@const columnType = getColumnType(type, format)}
						{#if columnType === 'boolean'}
							<BinaryIcon />
						{:else if columnType === 'number'}
							<HashIcon />
						{:else if columnType === 'time'}
							<ClockIcon />
						{:else if columnType === 'string'}
							<TypeIcon />
						{:else if columnType === 'array'}
							<BracketsIcon />
						{:else if columnType === 'object'}
							<BracesIcon />
						{:else if columnType === 'ratio'}
							<PercentIcon />
						{/if}
						{column.id}
					</DropdownMenu.Item>
				{/each}
			</DropdownMenu.Content>
		</DropdownMenu.Root>
		<!-- Accessors -->
		<div>
			{@render create?.()}
			{@render bulkDelete?.({ table })}
		</div>
		<div class="ml-auto">
			{@render reload()}
		</div>
	</div>
	{#if globalFilterError}
		<p class="text-xs text-destructive">
			{globalFilterError.message}
		</p>
	{/if}
	<!-- Table -->
	<div class="overflow-hidden rounded-md border bg-background">
		<Table.Root class="table-fixed">
			<Table.Header class="bg-muted">
				{#each table.getHeaderGroups() as headerGroup (headerGroup.id)}
					<Table.Row class="hover:bg-transparent">
						{#each headerGroup.headers as header (header.id)}
							<Table.Head
								style="width: {header.getSize()}px"
								class={cn(lodash.get(header.column.columnDef.meta, 'class'), 'h-11')}
							>
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
											<ChevronUpIcon class="shrink-0 opacity-60" size={16} aria-hidden="true" />
										{:else if header.column.getIsSorted() === 'desc'}
											<ChevronDownIcon class="shrink-0 opacity-60" size={16} aria-hidden="true" />
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
								<Table.Cell
									class={cn(
										getCellAlignment(fields[cell.column.id]),
										lodash.get(cell.column.columnDef.meta, 'class')
									)}
								>
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
										<Columns3Icon size={32} class="opacity-60" aria-hidden="true" />
									</Empty.Media>
									<Empty.Title>No Resources Found</Empty.Title>
									<Empty.Description>
										No resources found. Please adjust your filters or initiate a new resource to
										populate this table.
									</Empty.Description>
								</Empty.Header>
								<Empty.Content>
									<Button onclick={handleClear}>
										<EraserIcon size={16} class="opacity-60" />
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
							<ChevronFirstIcon size={16} aria-hidden="true" />
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
							<ChevronLeftIcon size={16} aria-hidden="true" />
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
							<ChevronRightIcon size={16} aria-hidden="true" />
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
							<ChevronLastIcon size={16} aria-hidden="true" />
						</Button>
					</Pagination.Item>
				</Pagination.Content>
			</Pagination.Root>
		</div>
	</div>
</div>
