<script lang="ts" module>
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import type { Row } from '@tanstack/table-core';
	import { formatTimeAgo } from '$lib/formatter';
	import { type TestResult, TestResult_Status, FIO_Input_AccessMode } from '$gen/api/bist/v1/bist_pb'
	import { timestampDate } from '@bufbuild/protobuf/wkt';
	import Icon from '@iconify/svelte';
	import { Badge } from '$lib/components/ui/badge/index.js';

	export const cells = {
		_row_picker: _row_picker,
		name: name,
		status: status,
		target: target,
		accessMode: accessMode,
		jobCount: jobCount,
		runTime: runTime,
		blockSize: blockSize,
		fileSize: fileSize,
		ioDepth: ioDepth,
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
		{#if TestResult_Status[row.original.status] === 'SUCCEEDED'}
			<Icon icon="ph:check" />
		{:else if TestResult_Status[row.original.status] === 'FAILED'}
			<Icon icon="ph:x" />
		{:else}
			<Icon icon="svg-spinners:180-ring-with-bg" />
		{/if}
	</p>
{/snippet}

{#snippet target(row: Row<TestResult>)}
	<p>
		{#if row.original.kind.case === 'fio' &&  row.original.kind.value?.input}
			{#if row.original.kind.value.target.case === 'cephBlockDevice' }
				<Badge variant="outline">
					{row.original.kind.value.target.value.facilityName}
				</Badge>
			{:else if row.original.kind.value.target.case === 'networkFileSystem' }
				<Badge variant="outline">
					{row.original.kind.value.target.value.endpoint}
				</Badge>
			{/if}
        {/if}
	</p>
{/snippet}

{#snippet accessMode(row: Row<TestResult>)}
	<p>
		{#if row.original.kind.case === 'fio' &&  row.original.kind.value?.input}
            {FIO_Input_AccessMode[row.original.kind.value?.input.accessMode]}
        {/if}
	</p>
{/snippet}

{#snippet jobCount(row: Row<TestResult>)}
	<p>
		{#if row.original.kind.case === 'fio' &&  row.original.kind.value?.input}
            {row.original.kind.value?.input.jobCount}
        {/if}
	</p>
{/snippet}

{#snippet runTime(row: Row<TestResult>)}
	<p>
		{#if row.original.kind.case === 'fio' &&  row.original.kind.value?.input}
            {row.original.kind.value?.input.runTime}
        {/if}
	</p>
{/snippet}

{#snippet blockSize(row: Row<TestResult>)}
	<p>
		{#if row.original.kind.case === 'fio' &&  row.original.kind.value?.input}
            {row.original.kind.value?.input.blockSize}
        {/if}
	</p>
{/snippet}

{#snippet fileSize(row: Row<TestResult>)}
	<p>
		{#if row.original.kind.case === 'fio' &&  row.original.kind.value?.input}
            {row.original.kind.value?.input.fileSize}
        {/if}
	</p>
{/snippet}

{#snippet ioDepth(row: Row<TestResult>)}
	<p>
		{#if row.original.kind.case === 'fio' &&  row.original.kind.value?.input}
            {row.original.kind.value?.input.ioDepth}
        {/if}
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


