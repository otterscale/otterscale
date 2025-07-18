<script lang="ts" module>
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import type { Row } from '@tanstack/table-core';
	import { formatTimeAgo } from '$lib/formatter';
	import { type TestResult, TestResult_Status } from '$gen/api/bist/v1/bist_pb'
	import { timestampDate } from '@bufbuild/protobuf/wkt';

	export const cells = {
		_row_picker: _row_picker,
		name: name,
		status: status,
		createdBy: createdBy,
		startedAt: startedAt,
		completedAt: completedAt
	};
</script>

{#snippet _row_picker(row: Row<TestResult>)}
	<Checkbox
		checked={row.getIsSelected()}
		onCheckedChange={(value) => row.toggleSelected(!!value)}
		class="border-secondary-950"
		aria-label="Select row"
	/>
{/snippet}

{#snippet name(row: Row<TestResult>)}
	<p>
		{row.original.name}
	</p>
{/snippet}

{#snippet status(row: Row<TestResult>)}
	<p>
		{TestResult_Status[row.original.status]}
	</p>
{/snippet}

{#snippet createdBy(row: Row<TestResult>)}
	<p>
		{row.original.createdBy}
	</p>
{/snippet}

{#snippet startedAt(row: Row<TestResult>)}
	{#if row.original.startedAt}
		{formatTimeAgo(timestampDate(row.original.startedAt))}
	{/if}
{/snippet}

{#snippet completedAt(row: Row<TestResult>)}
	{#if row.original.completedAt}
		{#if Number(timestampDate(row.original.completedAt)) >= 0}
			{formatTimeAgo(timestampDate(row.original.completedAt))}
		{/if}
	{/if}
{/snippet}


<!-- {#snippet startTime(row: Row<TestResult>)}
	<p>
		{formatTimeAgo(row.original.startTime)}
	</p>
{/snippet} -->
