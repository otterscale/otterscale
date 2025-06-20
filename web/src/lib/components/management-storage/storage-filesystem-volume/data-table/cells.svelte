<script lang="ts" module>
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';

	import { Badge } from '$lib/components/ui/badge';
	import { formatTimeAgo } from '$lib/formatter';
	import Icon from '@iconify/svelte';
	import type { Row } from '@tanstack/table-core';
	import { type Volume } from './types';

	export const cells = {
		_row_picker: _row_picker,
		name: name,
		enabled: enabled,
		permission: permission,
		createTime: createTime
	};
</script>

{#snippet _row_picker(row: Row<Volume>)}
	<Checkbox
		checked={row.getIsSelected()}
		onCheckedChange={(value) => row.toggleSelected(!!value)}
		class="border-secondary-950"
		aria-label="Select row"
	/>
{/snippet}

{#snippet name(row: Row<Volume>)}
	{row.original.name}
{/snippet}

{#snippet enabled(row: Row<Volume>)}
	{#if row.original.enabled}
		<Icon icon="ph:check" />
	{:else}
		<Icon icon="ph:x" />
	{/if}
{/snippet}

{#snippet permission(row: Row<Volume>)}
	<Badge variant="outline">
		{row.original.permission}
	</Badge>
{/snippet}

{#snippet createTime(row: Row<Volume>)}
	{formatTimeAgo(row.original.createTime)}
{/snippet}
