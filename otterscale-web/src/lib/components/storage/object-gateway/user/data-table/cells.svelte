<script lang="ts">
	import type { User } from '$lib/api/storage/v1/storage_pb';
	import { Cell as RowPicker } from '$lib/components/custom/data-table/data-table-row-pickers';
	import { Key } from '$lib/components/storage/object-gateway/user/key';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import Icon from '@iconify/svelte';
	import type { Row } from '@tanstack/table-core';
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
	<RowPicker {row} />
{/snippet}

{#snippet id(row: Row<User>)}
	{row.original.id}
{/snippet}

{#snippet name(row: Row<User>)}
	<Badge variant="outline">{row.original.name}</Badge>
{/snippet}

{#snippet suspended(row: Row<User>)}
	<div class="flex justify-end">
		{#if row.original.suspended}
			<Icon icon="ph:circle" class="text-primary" />
		{:else}
			<Icon icon="ph:x" class="text-destructive" />
		{/if}
	</div>
{/snippet}

{#snippet keys(row: Row<User>)}
	<Key user={row.original} />
{/snippet}

{#snippet actions(row: Row<User>)}
	<Actions user={row.original} />
{/snippet}
