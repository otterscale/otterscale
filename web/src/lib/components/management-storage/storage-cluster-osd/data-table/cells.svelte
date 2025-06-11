<script lang="ts" module>
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import * as Progress from '$lib/components/custom/progress';
	import type { Row } from '@tanstack/table-core';
	import type { OSD } from './types';
	import { Badge } from '$lib/components/ui/badge';
	import { formatCapacity } from '$lib/formatter';

	export const cells = {
		_row_picker: _row_picker,
		id: id,
		host: host,
		status: status,
		deviceClass: deviceClass,
		pgs: pgs,
		size: size,
		flags: flags,
		usage: usage
	};
</script>

{#snippet _row_picker(row: Row<OSD>)}
	<Checkbox
		checked={row.getIsSelected()}
		onCheckedChange={(value) => row.toggleSelected(!!value)}
		class="border-secondary-950"
		aria-label="Select row"
	/>
{/snippet}

{#snippet id(row: Row<OSD>)}
	{row.original.id}
{/snippet}

{#snippet host(row: Row<OSD>)}
	{row.original.host}
{/snippet}

{#snippet status(row: Row<OSD>)}
	<span class="flex items-center gap-1">
		{#each row.original.status as status}
			<Badge variant="outline">
				{status}
			</Badge>
		{/each}
	</span>
{/snippet}

{#snippet deviceClass(row: Row<OSD>)}
	<Badge variant="outline">
		{row.original.deviceClass}
	</Badge>
{/snippet}

{#snippet pgs(row: Row<OSD>)}
	<span class="flex justify-end">{row.original.pgs}</span>
{/snippet}

{#snippet size(row: Row<OSD>)}
	{@const size = formatCapacity(row.original.size)}
	<span class="flex items-center justify-end gap-1">
		{size.value}
		{size.unit}
	</span>
{/snippet}

{#snippet flags(row: Row<OSD>)}
	<span class="flex items-center gap-1">
		{#each row.original.flags as flag}
			<Badge variant="outline">
				{flag}
			</Badge>
		{/each}
	</span>
{/snippet}

{#snippet usage(row: Row<OSD>)}
	<Progress.Root numerator={row.original.usage} denominator={100}>
		{#snippet ratio({ numerator, denominator })}
			{(numerator * 100) / denominator}%
		{/snippet}
	</Progress.Root>
{/snippet}
