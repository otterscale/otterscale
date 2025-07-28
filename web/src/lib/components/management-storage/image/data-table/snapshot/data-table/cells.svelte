<script lang="ts" module>
	import type { Image_Snapshot } from '$gen/api/storage/v1/storage_pb';
	import TableRowPicker from '$lib/components/custom/data-table/data-table-row-pickers/cell.svelte';
	import * as Progress from '$lib/components/custom/progress/index.js';
	import Icon from '@iconify/svelte';
	import type { Row } from '@tanstack/table-core';

	export const cells = {
		_row_picker,
		name,
		protect,
		usage
	};
</script>

{#snippet _row_picker(row: Row<Image_Snapshot>)}
	<TableRowPicker {row} />
{/snippet}

{#snippet name(row: Row<Image_Snapshot>)}
	{row.original.name}
{/snippet}

{#snippet protect(row: Row<Image_Snapshot>)}
	<div class="flex justify-end">
		{#if row.original.protected}
			<Icon icon="ph:circle" class="text-primary" />
		{:else}
			<Icon icon="ph:x" class="text-destructive" />
		{/if}
	</div>
{/snippet}

{#snippet usage(row: Row<Image_Snapshot>)}
	{@const denominator = Number(row.original.quotaBytes)}
	{@const numerator = Number(row.original.usedBytes)}
	<div class="flex justify-end">
		<Progress.Root {numerator} {denominator}>
			{#snippet ratio({ numerator, denominator })}
				{Math.round((numerator * 100) / denominator)}%
			{/snippet}
		</Progress.Root>
	</div>
{/snippet}
