<script lang="ts" module>
	import { timestampDate } from '@bufbuild/protobuf/wkt';
	import Icon from '@iconify/svelte';
	import type { Row } from '@tanstack/table-core';

	import type { Secret } from '$lib/api/application/v1/application_pb';
	import CopyButton from '$lib/components/custom/copy-button/copy-button.svelte';
	import { Cells } from '$lib/components/custom/data-table/core';
	import * as Table from '$lib/components/custom/table/index.js';
	import { TagGroup } from '$lib/components/tag-group';
	import { Badge } from '$lib/components/ui/badge';
	import * as Dialog from '$lib/components/ui/dialog';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { formatTimeAgo } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';

	export const cells = {
		row_picker,
		name,
		namespace,
		type,
		immutable,
		labels,
		created_at
	};
</script>

{#snippet row_picker(row: Row<Secret>)}
	<Table.Cell alignClass="items-center">
		<Cells.RowPicker {row} />
	</Table.Cell>
{/snippet}

{#snippet name(row: Row<Secret>)}
	<Table.Cell alignClass="items-start">
		<Dialog.Root>
			<Dialog.Trigger>
				<p class="cursor-default underline hover:no-underline">
					{row.original.name}
				</p>
			</Dialog.Trigger>
			<Dialog.Content class="overflow-y-auto">
				<Dialog.Header>
					<Dialog.Title>
						<h1 class="text-center text-2xl">{m.secrets()}</h1>
					</Dialog.Title>
				</Dialog.Header>
				<div class="grid gap-4">
					{#each Object.entries(row.original.data) as [key, value] (key)}
						{@const secret = btoa(String.fromCharCode(...value))}
						<div class="group flex w-full items-center gap-2">
							<span
								class="group-hover:bg-text-card rounded-full bg-muted p-2 transition-colors duration-200 group-hover:bg-muted-foreground"
							>
								<Icon icon="ph:key" class="size-5" />
							</span>
							<div class="w-full space-y-1">
								<h3 class="text-xs font-medium text-muted-foreground">{key}</h3>
								<span class="flex items-center gap-1">
									<p class="max-w-sm truncate text-sm">{secret}</p>
									<CopyButton
										class="ml-auto size-4 opacity-0 transition-opacity duration-200 group-hover:opacity-100"
										text={secret}
									/>
								</span>
							</div>
						</div>
					{/each}
				</div>
			</Dialog.Content>
		</Dialog.Root>
	</Table.Cell>
{/snippet}

{#snippet namespace(row: Row<Secret>)}
	<Table.Cell alignClass="items-start">
		{row.original.namespace}
	</Table.Cell>
{/snippet}

{#snippet type(row: Row<Secret>)}
	<Table.Cell alignClass="items-start">
		<Badge variant="outline">
			{row.original.type}
		</Badge>
	</Table.Cell>
{/snippet}

{#snippet immutable(row: Row<Secret>)}
	{@const value = row.original.immutable}
	<Table.Cell alignClass="items-end">
		<Icon icon={value ? 'ph:check' : 'ph:x'} class={value ? 'text-green-500' : 'text-red-500'} />
	</Table.Cell>
{/snippet}

{#snippet labels(row: Row<Secret>)}
	<Table.Cell alignClass="items-end">
		<TagGroup
			items={Object.entries(row.original.labels).map(([key, value]) => ({
				title: `${key}: ${value}`,
				icon: 'ph:tag'
			}))}
		/>
	</Table.Cell>
{/snippet}

{#snippet created_at(row: Row<Secret>)}
	{#if row.original.createdAt}
		<Table.Cell alignClass="items-end">
			<Tooltip.Provider>
				<Tooltip.Root>
					<Tooltip.Trigger>
						{formatTimeAgo(timestampDate(row.original.createdAt))}
					</Tooltip.Trigger>
					<Tooltip.Content>
						{timestampDate(row.original.createdAt)}
					</Tooltip.Content>
				</Tooltip.Root>
			</Tooltip.Provider>
		</Table.Cell>
	{/if}
{/snippet}
