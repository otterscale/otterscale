<script lang="ts" module>
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import type { Row } from '@tanstack/table-core';
	import { type FlexibleIOTest } from './types';
	import { formatTimeAgo } from '$lib/formatter';

	export const cells = {
		_row_picker: _row_picker,
		name: name,
		rwMode: rwMode,
		fileSize: fileSize,
		numberJobs: numberJobs,
		blockSize: blockSize,
		runtime: runtime,
		createTime: createTime,
		modifyTime: modifyTime,
	};
</script>

{#snippet _row_picker(row: Row<FlexibleIOTest>)}
	<Checkbox
		checked={row.getIsSelected()}
		onCheckedChange={(value) => row.toggleSelected(!!value)}
		class="border-secondary-950"
		aria-label="Select row"
	/>
{/snippet}

{#snippet name(row: Row<FlexibleIOTest>)}
	<p>
		{row.original.name}
	</p>
{/snippet}

{#snippet rwMode(row: Row<FlexibleIOTest>)}
	<p>
		{row.original.rwMode}
	</p>
{/snippet}


{#snippet fileSize(row: Row<FlexibleIOTest>)}
	<p>
		{row.original.fileSize}
	</p>
{/snippet}

{#snippet numberJobs(row: Row<FlexibleIOTest>)}
	<p>
		{row.original.numberJobs}
	</p>
{/snippet}

{#snippet blockSize(row: Row<FlexibleIOTest>)}
	<p>
		{row.original.blockSize}
	</p>
{/snippet}

{#snippet runtime(row: Row<FlexibleIOTest>)}
	<p>
		{row.original.runtime}
	</p>
{/snippet}

{#snippet createTime(row: Row<FlexibleIOTest>)}
	<p>
		{formatTimeAgo(row.original.createTime)}
	</p>
{/snippet}

{#snippet modifyTime(row: Row<FlexibleIOTest>)}
	<p>
		{formatTimeAgo(row.original.modifyTime)}
	</p>
{/snippet}
