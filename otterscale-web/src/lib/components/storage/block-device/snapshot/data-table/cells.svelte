<script lang="ts" module>
	import type { Image_Snapshot } from '$lib/api/storage/v1/storage_pb';
	import { Cells } from '$lib/components/custom/data-table/core';
	import * as Progress from '$lib/components/custom/progress/index.js';
	import { formatCapacity } from '$lib/formatter';
	import Icon from '@iconify/svelte';
	import type { Row } from '@tanstack/table-core';
	import Actions from './cell-actions.svelte';

	export const cells = {
		row_picker,
		name,
		protect,
		usage,
		actions
	};
</script>

{#snippet row_picker(row: Row<Image_Snapshot>)}
	<Cells.RowPicker {row} />
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
	</div>
{/snippet}

{#snippet actions(row: Row<Image_Snapshot>)}
	<Actions snapshot={row.original} />
{/snippet}
