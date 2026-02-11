<script lang="ts" module>
	import type { JsonValue } from '@bufbuild/protobuf';

	export type ObjectOfKeyValueMetadata = { [key: string]: JsonValue };
</script>

<script lang="ts">
	import {
		FileCheckIcon,
		FileClockIcon,
		FileCodeIcon,
		FileDigit,
		FileIcon,
		FileTextIcon,
		GridIcon,
		ListIcon
	} from '@lucide/svelte/icons';
	import { type Column, type Row } from '@tanstack/table-core';

	import * as Code from '$lib/components/custom/code/index.js';
	import { Button } from '$lib/components/ui/button';
	import * as Collapsible from '$lib/components/ui/collapsible';
	import * as Item from '$lib/components/ui/item';
	import * as Sheet from '$lib/components/ui/sheet/index.js';
	import * as Table from '$lib/components/ui/table/index.js';
	import * as Tabs from '$lib/components/ui/tabs/index.js';

	import { format } from '../utils';

	let {
		row,
		column,
		metadata
	}: {
		row: Row<Record<string, JsonValue>>;
		column: Column<Record<string, JsonValue>>;
		metadata: ObjectOfKeyValueMetadata;
	} = $props();

	const data = $derived(row.original[column.id] as number);
</script>

{#if !data}
	<Sheet.Root>
		<Sheet.Trigger>
			<Button variant="ghost" disabled>
				<FileCodeIcon />
			</Button>
		</Sheet.Trigger>
	</Sheet.Root>
{:else}
	<Sheet.Root>
		<Sheet.Trigger>
			<Button variant="ghost" class="hover:underline">
				{data}
			</Button>
		</Sheet.Trigger>
		<Sheet.Content
			side="right"
			class="flex h-full max-w-[62vw] min-w-[50vw] flex-col gap-0 overflow-y-auto p-4"
		>
			<Sheet.Header class="shrink-0 space-y-4">
				<Sheet.Title>{column.id}</Sheet.Title>
			</Sheet.Header>
			<Tabs.Root value="grid" class="p-4">
				<Tabs.List class="ml-auto">
					<Tabs.Trigger value="grid">
						<GridIcon />
					</Tabs.Trigger>
					<Tabs.Trigger value="table">
						<ListIcon />
					</Tabs.Trigger>
				</Tabs.List>
				<Tabs.Content value="grid">
					<div class="space-y-0">
						{#each Object.entries(metadata) as [key, value], index (index)}
							{#if typeof value === 'string'}
								{@const formatted = format(value)}
								{@const isExpandable = formatted.split('\n').length > 2}
								<Collapsible.Root
									class="rounded-lg transition-colors duration-200 hover:bg-muted/50"
								>
									<Collapsible.Trigger
										disabled={!isExpandable}
										class="w-full transition-colors duration-200 hover:underline"
									>
										<Item.Root size="sm">
											<Item.Media variant="icon">
												{#if !value}
													<FileIcon />
												{:else if ['true', 'false'].includes(value)}
													<FileCheckIcon />
												{:else if !isNaN(Number(value))}
													<FileDigit />
												{:else if !isNaN(Date.parse(value))}
													<FileClockIcon />
												{:else}
													<FileTextIcon />
												{/if}
											</Item.Media>
											<Item.Content class="min-w-0 flex-1 text-left">
												<Item.Title class="w-full">
													{key}
												</Item.Title>
												<Item.Description class="wrap-break-words breaks-all">
													{value}
												</Item.Description>
											</Item.Content>
										</Item.Root>
									</Collapsible.Trigger>
									<Collapsible.Content
										class="overflow-hidden transition-all duration-300 ease-in-out"
									>
										<Code.Root
											lang="json"
											hideLines
											code={formatted}
											class="border-none bg-transparent px-8"
										/>
									</Collapsible.Content>
								</Collapsible.Root>
							{/if}
						{/each}
					</div>
				</Tabs.Content>
				<Tabs.Content value="table">
					<div>
						<Table.Root class="[&_td:first-child]:rounded-l-lg [&_td:last-child]:rounded-r-lg">
							<Table.Header>
								<Table.Row class="hover:[&>th,td]:bg-transparent!">
									<Table.Head>Key</Table.Head>
									<Table.Head>Value</Table.Head>
								</Table.Row>
							</Table.Header>
							<Table.Body>
								{#each Object.entries(metadata) as [key, value], index (index)}
									<Table.Row class="border-none align-top">
										<Table.Cell class="align-top">{key}</Table.Cell>
										<Table.Cell class="align-top">
											<p
												class="wrap-break-words max-w-3xl text-sm leading-normal font-normal text-balance break-all text-muted-foreground"
											>
												{value}
											</p>
										</Table.Cell>
									</Table.Row>
								{/each}
							</Table.Body>
						</Table.Root>
					</div>
				</Tabs.Content>
			</Tabs.Root>
		</Sheet.Content>
	</Sheet.Root>
{/if}
