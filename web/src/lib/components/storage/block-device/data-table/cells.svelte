<script lang="ts" module>
	import type { Row } from '@tanstack/table-core';

	import Actions from './cell-actions.svelte';

	import type { Image } from '$lib/api/storage/v1/storage_pb';
	import { Cells } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import * as Progress from '$lib/components/custom/progress/index.js';
	import { Snapshot } from '$lib/components/storage/block-device/snapshot';
	import { Badge } from '$lib/components/ui/badge';
	import { formatCapacity } from '$lib/formatter';

	export const cells = {
		row_picker,
		name,
		poolName,
		usage,
		snapshots,
		actions
	};
</script>

{#snippet row_picker(row: Row<Image>)}
	<Layout.Cell class="items-center">
		<Cells.RowPicker {row} />
	</Layout.Cell>
{/snippet}

{#snippet name(row: Row<Image>)}
	<Layout.Cell class="items-start">
		{row.original.name}
	</Layout.Cell>
{/snippet}

{#snippet poolName(row: Row<Image>)}
	<Layout.Cell class="items-start">
		<Badge variant="outline">{row.original.poolName}</Badge>
	</Layout.Cell>
{/snippet}

{#snippet usage(row: Row<Image>)}
	{@const denominator = Number(row.original.quotaBytes)}
	{@const numerator = Number(row.original.usedBytes)}
	<Layout.Cell class="items-end">
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
	</Layout.Cell>
{/snippet}

{#snippet snapshots(row: Row<Image>)}
	<Layout.Cell class="items-end">
		<Snapshot image={row.original} />
	</Layout.Cell>
{/snippet}

{#snippet actions(row: Row<Image>)}
	<Layout.Cell class="items-start">
		<Actions image={row.original} />
	</Layout.Cell>
{/snippet}
