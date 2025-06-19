<script lang="ts" module>
	import * as Progress from '$lib/components/custom/progress';
	import { Badge } from '$lib/components/ui/badge';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import { formatCapacity } from '$lib/formatter';
	import type { Row } from '@tanstack/table-core';
	import { ActionViewDevice } from './action-view-device';
	import type { ObjectStorageDaemon } from './types';

	export const cells = {
		_row_picker: _row_picker,
		id: id,
		host: host,
		devices: devices,
		status: status,
		deviceClass: deviceClass,
		pgs: pgs,
		size: size,
		flags: flags,
		usage: usage
	};
</script>

{#snippet _row_picker(row: Row<ObjectStorageDaemon>)}
	<Checkbox
		checked={row.getIsSelected()}
		onCheckedChange={(value) => row.toggleSelected(!!value)}
		class="border-secondary-950"
		aria-label="Select row"
	/>
{/snippet}

{#snippet id(row: Row<ObjectStorageDaemon>)}
	{row.original.id}
{/snippet}

{#snippet host(row: Row<ObjectStorageDaemon>)}
	{row.original.host}
{/snippet}

{#snippet devices(row: Row<ObjectStorageDaemon>)}
	<span class="flex items-center justify-end gap-1">
		{row.original.devices}
		<ActionViewDevice row={row.original} />
	</span>
{/snippet}

{#snippet status(row: Row<ObjectStorageDaemon>)}
	<span class="flex items-center gap-1">
		{#each row.original.status as status}
			<Badge variant="outline">
				{status}
			</Badge>
		{/each}
	</span>
{/snippet}

{#snippet deviceClass(row: Row<ObjectStorageDaemon>)}
	<Badge variant="outline">
		{row.original.deviceClass}
	</Badge>
{/snippet}

{#snippet pgs(row: Row<ObjectStorageDaemon>)}
	<span class="flex justify-end">{row.original.pgs}</span>
{/snippet}

{#snippet size(row: Row<ObjectStorageDaemon>)}
	{@const size = formatCapacity(row.original.size)}
	<span class="flex items-center justify-end gap-1">
		{size.value}
		{size.unit}
	</span>
{/snippet}

{#snippet flags(row: Row<ObjectStorageDaemon>)}
	<span class="flex items-center gap-1">
		{#each row.original.flags as flag}
			<Badge variant="outline">
				{flag}
			</Badge>
		{/each}
	</span>
{/snippet}

{#snippet usage(row: Row<ObjectStorageDaemon>)}
	<Progress.Root numerator={row.original.usage} denominator={100}>
		{#snippet ratio({ numerator, denominator })}
			{(numerator * 100) / denominator}%
		{/snippet}
	</Progress.Root>
{/snippet}
