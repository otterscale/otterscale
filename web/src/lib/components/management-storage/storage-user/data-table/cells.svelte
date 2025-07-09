<script lang="ts" module>
	import type { User } from '$gen/api/storage/v1/storage_pb';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import Icon from '@iconify/svelte';
	import type { Row } from '@tanstack/table-core';

	export const cells = {
		_row_picker: _row_picker,
		id: id,
		name: name,
		suspended: suspended
	};
</script>

{#snippet _row_picker(row: Row<User>)}
	<Checkbox
		checked={row.getIsSelected()}
		onCheckedChange={(value) => row.toggleSelected(!!value)}
		class="border-secondary-950"
		aria-label="Select row"
	/>
{/snippet}

{#snippet id(row: Row<User>)}
	{row.original.id}
{/snippet}

{#snippet name(row: Row<User>)}
	{row.original.name}
{/snippet}

{#snippet suspended(row: Row<User>)}
	{#if row.original.suspended}
		<Icon icon="ph:x" class="text-destructive" />
	{:else}
		<Icon icon="ph:check" class="text-success" />
	{/if}
{/snippet}
