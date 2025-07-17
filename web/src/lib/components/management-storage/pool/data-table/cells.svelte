<script lang="ts" module>
	import type { Pool } from '$gen/api/storage/v1/storage_pb';
	import TableRowPicker from '$lib/components/custom/data-table/data-table-row-pickers/cell.svelte';
	import * as Progress from '$lib/components/custom/progress';
	import { Badge } from '$lib/components/ui/badge';
	import type { Row } from '@tanstack/table-core';

	export const cells = {
		_row_picker: _row_picker,
		name: name,
		applications: applications,
		placement_group_state: placement_group_state,
		usage: usage
	};
</script>

{#snippet _row_picker(row: Row<Pool>)}
	<TableRowPicker {row} />
{/snippet}

{#snippet name(row: Row<Pool>)}
	{row.original.name}
{/snippet}

{#snippet applications(row: Row<Pool>)}
	{#each row.original.applications as application}
		{#if application}
			<Badge variant="outline">
				{application}
			</Badge>
		{/if}
	{/each}
{/snippet}

{#snippet placement_group_state(row: Row<Pool>)}
	{#each Object.entries(row.original.placementGroupState) as [state, number]}
		<Badge variant="outline">
			{state}:{number}
		</Badge>
	{/each}
{/snippet}

{#snippet usage(row: Row<Pool>)}
	<div class="flex justify-end">
		<Progress.Root
			numerator={Number(row.original.usedBytes)}
			denominator={Number(row.original.quotaBytes)}
		>
			{#snippet ratio({ numerator, denominator })}
				{(numerator * 100) / denominator}%
			{/snippet}
		</Progress.Root>
	</div>
{/snippet}
