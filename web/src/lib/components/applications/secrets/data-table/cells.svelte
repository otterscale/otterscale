<script lang="ts" module>
	import { timestampDate } from '@bufbuild/protobuf/wkt';
	import Icon from '@iconify/svelte';
	import type { Row } from '@tanstack/table-core';

	import type { Secret } from '$lib/api/application/v1/application_pb';
	import CopyButton from '$lib/components/custom/copy-button/copy-button.svelte';
	import { Cells } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
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
	<Layout.Cell class="items-center">
		<Cells.RowPicker {row} />
	</Layout.Cell>
{/snippet}

{#snippet name(row: Row<Secret>)}
	<Layout.Cell class="items-start">
		<Dialog.Root>
			<Dialog.Header>
				<Dialog.Trigger>
					<p class="cursor-default underline hover:no-underline">
						{row.original.name}
					</p>
				</Dialog.Trigger>
			</Dialog.Header>
			<Dialog.Content class="overflow-y-auto">
				<Dialog.Title>
					<h1 class="text-center text-2xl">{m.secrets()}</h1>
				</Dialog.Title>
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
	</Layout.Cell>
{/snippet}

{#snippet namespace(row: Row<Secret>)}
	<Layout.Cell class="items-start">
		{row.original.namespace}
	</Layout.Cell>
{/snippet}

{#snippet type(row: Row<Secret>)}
	<Layout.Cell class="items-start">
		<Badge variant="outline">
			{row.original.type}
		</Badge>
	</Layout.Cell>
{/snippet}

{#snippet immutable(row: Row<Secret>)}
	<Layout.Cell class="items-end">
		{@const value = row.original.immutable}
		<Icon icon={value ? 'ph:check' : 'ph:x'} class={value ? 'text-green-500' : 'text-red-500'} />
	</Layout.Cell>
{/snippet}

{#snippet labels(row: Row<Secret>)}
	<Layout.Cell class="items-end">
		<TagGroup
			items={Object.entries(row.original.labels).map(([key, value]) => ({
				title: `${key}: ${value}`,
				icon: 'ph:tag'
			}))}
		/>
	</Layout.Cell>
{/snippet}

{#snippet created_at(row: Row<Secret>)}
	{#if row.original.createdAt}
		<Layout.Cell class="items-end">
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
		</Layout.Cell>
	{/if}
{/snippet}
