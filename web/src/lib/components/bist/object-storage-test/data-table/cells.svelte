<script lang="ts" module>
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import type { Row } from '@tanstack/table-core';
	import { formatTimeAgo } from '$lib/formatter';
	import { type TestResult, TestResult_Status, InternalObjectService_Type, Warp_Input_Operation } from '$gen/api/bist/v1/bist_pb'
	import { timestampDate } from '@bufbuild/protobuf/wkt';
	import Icon from '@iconify/svelte';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { formatByte, formatSecond } from '$lib/formatter';

	export const cells = {
		_row_picker: _row_picker,
		name: name,
		status: status,
		target: target,
		operation: operation,
		duration: duration,
		objectSize: objectSize, 
		objectCount: objectCount,
		throughputFastest: throughputFastest,
		throughputSlowest: throughputSlowest,
		throughputMedian: throughputMedian,
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
		{#if row.original.kind.case === 'warp' &&  row.original.kind.value?.input}
			{#if row.original.kind.value.target.case === 'internalObjectService' }
				<Badge variant="outline">
					{InternalObjectService_Type[row.original.kind.value.target.value.type]}-{row.original.kind.value.target.value.facilityName}
				</Badge>
			{:else if row.original.kind.value.target.case === 'externalObjectService' }
				<Badge variant="outline">
					{row.original.kind.value.target.value.endpoint}
				</Badge>
			{/if}
        {/if}
	</p>
{/snippet}

{#snippet createdBy(row: Row<TestResult>)}
	<p>
		{row.original.createdBy}
	</p>
{/snippet}

{#snippet operation(row: Row<TestResult>)}
	<p>
		{#if row.original.kind.case === 'warp' &&  row.original.kind.value?.input}
			{Warp_Input_Operation[row.original.kind.value.input.operation]}
        {/if}
	</p>
{/snippet}
{#snippet duration(row: Row<TestResult>)}
	<p>
		{#if row.original.kind.case === 'warp' &&  row.original.kind.value?.input}
			{@const formatted = formatSecond(Number(row.original.kind.value?.input.durationSeconds))}
			{formatted.value} {formatted.unit}
        {/if}
	</p>
{/snippet}
{#snippet objectSize(row: Row<TestResult>)}
	<p>
		{#if row.original.kind.case === 'warp' &&  row.original.kind.value?.input}
            {@const formatted = formatByte(Number(row.original.kind.value?.input.objectSizeBytes))}
            {formatted.value} {formatted.unit}
        {/if}
	</p>
{/snippet}
{#snippet objectCount(row: Row<TestResult>)}
	<p>
		{#if row.original.kind.case === 'warp' &&  row.original.kind.value?.input}
			{row.original.kind.value.input.objectCount}
        {/if}
	</p>
{/snippet}

{#snippet throughputFastest(row: Row<TestResult>)}
	{#if row.original.kind.case === 'warp' && row.original.kind.value?.output?.get?.bytes}
		<p>
			<Badge variant="outline">
				G: {(Number(row.original.kind.value.output.get.bytes.fastestPerSecond) / 1000000).toFixed(3)} MB/s
			</Badge>
		</p>
	{/if}
	{#if row.original.kind.case === 'warp' && row.original.kind.value?.output?.put?.bytes}
		<p class={row.original.kind.value?.output?.put ? 'mt-1' : ''}>
			<Badge variant="outline">
				P: {(Number(row.original.kind.value.output.put.bytes.fastestPerSecond) / 1000000).toFixed(3)} MB/s
			</Badge>
		</p>
	{/if}
	{#if row.original.kind.case === 'warp' && row.original.kind.value?.output?.delete?.bytes}
		<p class={(row.original.kind.value?.output?.get || row.original.kind.value?.output?.put) ? 'mt-1' : ''}>
			<Badge variant="outline">
				D: {(Number(row.original.kind.value.output.delete.bytes.fastestPerSecond) / 1000000).toFixed(3)} MB/s
			</Badge>
		</p>
	{/if}
{/snippet}

{#snippet throughputSlowest(row: Row<TestResult>)}
	{#if row.original.kind.case === 'warp' && row.original.kind.value?.output?.get?.bytes}
		<p>
			<Badge variant="outline">
				G: {(Number(row.original.kind.value.output.get.bytes.slowestPerSecond) / 1000000).toFixed(3)} MB/s
			</Badge>
		</p>
	{/if}
	{#if row.original.kind.case === 'warp' && row.original.kind.value?.output?.put?.bytes}
		<p class={row.original.kind.value?.output?.put ? 'mt-1' : ''}>
			<Badge variant="outline">
				P: {(Number(row.original.kind.value.output.put.bytes.slowestPerSecond) / 1000000).toFixed(3)} MB/s
			</Badge>
		</p>
	{/if}
	{#if row.original.kind.case === 'warp' && row.original.kind.value?.output?.delete?.bytes}
		<p class={(row.original.kind.value?.output?.get || row.original.kind.value?.output?.put) ? 'mt-1' : ''}>
			<Badge variant="outline">
				D: {(Number(row.original.kind.value.output.delete.bytes.slowestPerSecond) / 1000000).toFixed(3)} MB/s
			</Badge>
		</p>
	{/if}
{/snippet}

{#snippet throughputMedian(row: Row<TestResult>)}
	{#if row.original.kind.case === 'warp' && row.original.kind.value?.output?.get?.bytes}
		<p>
			<Badge variant="outline">
				G: {(Number(row.original.kind.value.output.get.bytes.medianPerSecond) / 1000000).toFixed(3)} MB/s
			</Badge>
		</p>
	{/if}
	{#if row.original.kind.case === 'warp' && row.original.kind.value?.output?.put?.bytes}
		<p class={row.original.kind.value?.output?.put ? 'mt-1' : ''}>
			<Badge variant="outline">
				P: {(Number(row.original.kind.value.output.put.bytes.medianPerSecond) / 1000000).toFixed(3)} MB/s
			</Badge>
		</p>
	{/if}
	{#if row.original.kind.case === 'warp' && row.original.kind.value?.output?.delete?.bytes}
		<p class={(row.original.kind.value?.output?.get || row.original.kind.value?.output?.put) ? 'mt-1' : ''}>
			<Badge variant="outline">
				D: {(Number(row.original.kind.value.output.delete.bytes.medianPerSecond) / 1000000).toFixed(3)} MB/s
			</Badge>
		</p>
	{/if}
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
