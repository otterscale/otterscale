<script lang="ts" module>
	import type { JsonValue } from '@bufbuild/protobuf';

	export type ArrayOfObjectItemType = {
		title: JsonValue;
		description: JsonValue;
		actions: JsonValue;
	};

	export type ArrayOfObjectItemsType = ArrayOfObjectItemType[];

	export type ArrayOfObjectMetadata = { items: ArrayOfObjectItemsType };
</script>

<script lang="ts">
	import Braces from '@lucide/svelte/icons/braces';
	import ChevronFirst from '@lucide/svelte/icons/chevron-first';
	import ChevronLast from '@lucide/svelte/icons/chevron-last';
	import ChevronLeft from '@lucide/svelte/icons/chevron-left';
	import ChevronRight from '@lucide/svelte/icons/chevron-right';
	import File from '@lucide/svelte/icons/file';
	import {
		type Column,
		type ColumnDef,
		getCoreRowModel,
		getFilteredRowModel,
		getPaginationRowModel,
		type PaginationState,
		type Row
	} from '@tanstack/table-core';

	import * as CodeBlock from '$lib/components/custom/code/index.js';
	import { Button } from '$lib/components/ui/button';
	import * as Collapsible from '$lib/components/ui/collapsible';
	import { createSvelteTable } from '$lib/components/ui/data-table';
	import * as Empty from '$lib/components/ui/empty/index.js';
	import * as Item from '$lib/components/ui/item';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Select, SelectContent, SelectItem, SelectTrigger } from '$lib/components/ui/select';
	import * as Sheet from '$lib/components/ui/sheet/index.js';
	import { onMount } from 'svelte';

	let {
		row,
		column,
		metadata
	}: {
		row: Row<Record<string, JsonValue>>;
		column: Column<Record<string, JsonValue>>;
		metadata: ArrayOfObjectMetadata;
	} = $props();

	const data = $derived(row.original[column.id] as number);

	let pagination = $state<PaginationState>({ pageIndex: 0, pageSize: 10 });

	const columns: ColumnDef<ArrayOfObjectItemType>[] = [
		{
			id: 'item',
			accessorFn: (row) => row,
			cell: ({ row }) => row.original
		}
	];

	let table = createSvelteTable<ArrayOfObjectItemType>({
		columns,
		get data() {
			return metadata.items;
		},
		getCoreRowModel: getCoreRowModel(),
		getFilteredRowModel: getFilteredRowModel(),
		getPaginationRowModel: getPaginationRowModel(),
		onPaginationChange: (updater) => {
			if (typeof updater === 'function') {
				pagination = updater(pagination);
			} else {
				pagination = updater;
			}
		},
		state: {
			get pagination() {
				return pagination;
			}
		}
	});

	onMount(() => {
		if (metadata === undefined) {
			console.warn(
				`Expected metadata of ${column.id} for ArrayOfObjectCell, but got metadata:`,
				metadata
			);
		}
	});
</script>

<Sheet.Root>
	<Sheet.Trigger>
		<Button variant="ghost" class="hover:underline">
			{data}
		</Button>
	</Sheet.Trigger>
	{#if metadata.items}
		<Sheet.Content
			side="right"
			class="flex h-full max-w-[50vw] min-w-[38vw] flex-col gap-0 overflow-y-auto p-4"
		>
			<Sheet.Header class="shrink-0 space-y-4">
				<Sheet.Title>{column.id}</Sheet.Title>
			</Sheet.Header>

			{#if table.getRowModel().rows.length > 0}
				<div class="flex-1 space-y-0 overflow-y-auto">
					{#each table.getRowModel().rows as row (row.id)}
						{@const item = row.original}
						<Collapsible.Root class="rounded-lg transition-colors duration-200 hover:bg-muted/50">
							<Collapsible.Trigger class="w-full transition-colors duration-200 hover:underline">
								<Item.Root size="sm">
									<Item.Media variant="icon">
										<File />
									</Item.Media>
									<Item.Content class="min-w-0 flex-1 text-left">
										<Item.Title class="w-full">
											{item.title}
										</Item.Title>
										<Item.Description class="wrap-break-words breaks-all">
											{item.description}
										</Item.Description>
									</Item.Content>
									<Item.Actions>
										{item.actions}
									</Item.Actions>
								</Item.Root>
							</Collapsible.Trigger>
							<Collapsible.Content class="overflow-hidden transition-all duration-300 ease-in-out">
								<CodeBlock.Root
									lang="json"
									hideLines
									code={JSON.stringify(item, null, 4)}
									class="border-none bg-transparent px-8"
								/>
							</Collapsible.Content>
						</Collapsible.Root>
					{/each}
				</div>

				<!-- Pagination -->
				<div class="mt-4 flex shrink-0 items-center justify-between gap-4 border-t pt-4">
					<div class="flex items-center gap-2">
						<Label class="text-sm">Rows per page</Label>
						<Select
							type="single"
							value={table.getState().pagination.pageSize.toString()}
							onValueChange={(value) => table.setPageSize(Number(value))}
						>
							<SelectTrigger class="w-fit">
								{table.getState().pagination.pageSize}
							</SelectTrigger>
							<SelectContent>
								{#each [5, 10, 25, 50] as size (size)}
									<SelectItem value={size.toString()}>{size}</SelectItem>
								{/each}
							</SelectContent>
						</Select>
					</div>

					<p class="text-sm text-muted-foreground">
						<span class="text-foreground">
							{table.getState().pagination.pageIndex * table.getState().pagination.pageSize +
								1}-{Math.min(
								(table.getState().pagination.pageIndex + 1) * table.getState().pagination.pageSize,
								table.getFilteredRowModel().rows.length
							)}
						</span>
						of
						<span class="text-foreground">{table.getFilteredRowModel().rows.length}</span>
					</p>

					<div class="flex items-center gap-1">
						<Button
							size="icon"
							variant="outline"
							onclick={() => table.firstPage()}
							disabled={!table.getCanPreviousPage()}
							aria-label="Go to first page"
						>
							<ChevronFirst size={16} />
						</Button>
						<Button
							size="icon"
							variant="outline"
							onclick={() => table.previousPage()}
							disabled={!table.getCanPreviousPage()}
							aria-label="Go to previous page"
						>
							<ChevronLeft size={16} />
						</Button>
						<Button
							size="icon"
							variant="outline"
							onclick={() => table.nextPage()}
							disabled={!table.getCanNextPage()}
							aria-label="Go to next page"
						>
							<ChevronRight size={16} />
						</Button>
						<Button
							size="icon"
							variant="outline"
							onclick={() => table.lastPage()}
							disabled={!table.getCanNextPage()}
							aria-label="Go to last page"
						>
							<ChevronLast size={16} />
						</Button>
					</div>
				</div>
			{:else}
				<Empty.Root class="m-4 bg-muted/50">
					<Empty.Header>
						<Empty.Media variant="icon">
							<Braces size={36} />
						</Empty.Media>
						<Empty.Title>No Data</Empty.Title>
						<Empty.Description>
							No data is currently available for this resource.
							<br />
							To populate this resource, please add properties or values through the resource editor.
						</Empty.Description>
					</Empty.Header>
				</Empty.Root>
			{/if}
		</Sheet.Content>
	{/if}
</Sheet.Root>
