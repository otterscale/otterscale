<script lang="ts" module>
	import type { Row } from '@tanstack/table-core';

	import type { Image } from '$lib/api/storage/v1/storage_pb';
	import { Cells } from '$lib/components/custom/data-table/core';
	import * as Progress from '$lib/components/custom/progress/index.js';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Table from '$lib/components/custom/table/index.js';
	import { Snapshot } from '$lib/components/storage/block-device/snapshot';
	import { Badge } from '$lib/components/ui/badge';
	import { formatCapacity } from '$lib/formatter';

	import Actions from './cell-actions.svelte';

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
	<Table.Cell alignClass="items-center">
		<Cells.RowPicker {row} />
	</Table.Cell>
{/snippet}

{#snippet name(row: Row<Image>)}
	<Table.Cell alignClass="items-start">
		{row.original.name}
	</Table.Cell>
{/snippet}

{#snippet poolName(row: Row<Image>)}
	<Table.Cell alignClass="items-start">
		<Badge variant="outline">{row.original.poolName}</Badge>
	</Table.Cell>
{/snippet}

{#snippet usage(row: Row<Image>)}
	{@const denominator = Number(row.original.quotaBytes)}
	{@const numerator = Number(row.original.usedBytes)}
	<Table.Cell alignClass="items-end">
		<Progress.Root {numerator} {denominator} target="STB">
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
	</Table.Cell>
{/snippet}

{#snippet snapshots(data: { row: Row<Image>; scope: string; reloadManager: ReloadManager })}
	<Table.Cell alignClass="items-end">
		<Snapshot image={data.row.original} scope={data.scope} reloadManager={data.reloadManager} />
	</Table.Cell>
{/snippet}

{#snippet actions(data: { row: Row<Image>; scope: string; reloadManager: ReloadManager })}
	<Table.Cell alignClass="items-start">
		<Actions image={data.row.original} scope={data.scope} reloadManager={data.reloadManager} />
	</Table.Cell>
{/snippet}
