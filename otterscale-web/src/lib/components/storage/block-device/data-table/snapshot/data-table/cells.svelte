<script lang="ts" module>
	import type { Image_Snapshot } from '$lib/api/storage/v1/storage_pb';
	import TableRowPicker from '$lib/components/custom/data-table/data-table-row-pickers/cell.svelte';
	import * as Progress from '$lib/components/custom/progress/index.js';
	import { formatCapacityV2 as formatCapacity } from '$lib/formatter';
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
	<Progress.Root {numerator} {denominator}>
		{#snippet ratio({ numerator, denominator })}
			{Progress.formatRatio(numerator, denominator)}
		{/snippet}
		{#snippet detail({ numerator, denominator })}
			{@const { value: numeratorValue, unit: numeratorUnit } = formatCapacity(numerator)}
			{@const { value: denominatorValue, unit: denominatorUnit } = formatCapacity(denominator)}
			{numeratorValue}
			{numeratorUnit}/{denominatorValue}
			{denominatorUnit}
		{/snippet}
	</Progress.Root>
{/snippet}
