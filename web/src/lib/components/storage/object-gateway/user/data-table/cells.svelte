<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import type { Row } from '@tanstack/table-core';

	import type { User } from '$lib/api/storage/v1/storage_pb';
	import { Cells } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { Key } from '$lib/components/storage/object-gateway/user/key';
	import Badge from '$lib/components/ui/badge/badge.svelte';

	import Actions from './cell-actions.svelte';

	export const cells = {
		row_picker,
		id,
		name,
		suspended,
		keys,
		actions
	};
</script>

{#snippet row_picker(row: Row<User>)}
	<Layout.Cell class="items-center">
		<Cells.RowPicker {row} />
	</Layout.Cell>
{/snippet}

{#snippet id(row: Row<User>)}
	<Layout.Cell class="items-start">
		{row.original.id}
	</Layout.Cell>
{/snippet}

{#snippet name(row: Row<User>)}
	<Layout.Cell class="items-start">
		<Badge variant="outline">{row.original.name}</Badge>
	</Layout.Cell>
{/snippet}

{#snippet suspended(row: Row<User>)}
	<Layout.Cell class="items-end">
		{#if row.original.suspended}
			<Icon icon="ph:circle" class="text-primary" />
		{:else}
			<Icon icon="ph:x" class="text-destructive" />
		{/if}
	</Layout.Cell>
{/snippet}

{#snippet keys(row: Row<User>)}
	<Layout.Cell class="items-end">
		<Key user={row.original} />
	</Layout.Cell>
{/snippet}

{#snippet actions(row: Row<User>)}
	<Layout.Cell class="items-start">
		<Actions user={row.original} />
	</Layout.Cell>
{/snippet}
