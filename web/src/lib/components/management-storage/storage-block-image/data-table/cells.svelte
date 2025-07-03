<script lang="ts" module>
	import * as Progress from '$lib/components/custom/progress/index.js';
	import { Badge } from '$lib/components/ui/badge';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import { formatCapacity } from '$lib/formatter';
	import { Snapshot } from './snapshot';
	import type { Row } from '@tanstack/table-core';
	import type { Image } from '$gen/api/storage/v1/storage_pb';

	export const cells = {
		_row_picker: _row_picker,
		name: name,
		poolName: poolName,
		usage: usage,
		snapshots: snapshots
	};
</script>

{#snippet _row_picker(row: Row<Image>)}
	<Checkbox
		checked={row.getIsSelected()}
		onCheckedChange={(value) => row.toggleSelected(!!value)}
		class="border-secondary-950"
		aria-label="Select row"
	/>
{/snippet}

{#snippet name(row: Row<Image>)}
	{row.original.name}
{/snippet}

{#snippet poolName(row: Row<Image>)}
	<Badge variant="outline">{row.original.poolName}</Badge>
{/snippet}

{#snippet usage(row: Row<Image>)}
	{@const denominator = Number(row.original.quotaBytes)}
	{@const numerator = Number(row.original.usedBytes)}
	<div class="flex justify-end">
		<Progress.Root {numerator} {denominator}>
			{#snippet ratio({ numerator, denominator })}
				{Math.round((numerator * 100) / denominator)}%
			{/snippet}
			{#snippet detail({ numerator, denominator })}
				{numerator}/{denominator}
			{/snippet}
		</Progress.Root>
	</div>
{/snippet}

{#snippet snapshots(selectedScope, selectedFacility, row: Row<Image>)}
	<div class="flex justify-end">
		<span class="flex items-center gap-1">
			{row.original.snapshots.length}
			<Snapshot {selectedScope} {selectedFacility} data={row.original} />
		</span>
	</div>
{/snippet}
