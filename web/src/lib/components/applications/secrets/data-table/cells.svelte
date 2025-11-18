<script lang="ts" module>
	import { timestampDate } from '@bufbuild/protobuf/wkt';
	import Icon from '@iconify/svelte';
	import type { Row } from '@tanstack/table-core';

	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import type { Secret } from '$lib/api/application/v1/application_pb';
	import { Cells } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { TagGroup } from '$lib/components/tag-group';
	import { Badge } from '$lib/components/ui/badge';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { formatTimeAgo } from '$lib/formatter';

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
		<a
			class="underline hover:no-underline"
			href={resolve('/(auth)/scope/[scope]/applications/workloads/[namespace]/[application_name]', {
				scope: page.params.scope!,
				namespace: row.original.namespace!,
				application_name: row.original.name!
			})}
		>
			{row.original.name}
		</a>
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
