<script lang="ts" module>
	import { InternalObjectService_Type, type TestResult, TestResult_Status, Warp_Input_Operation } from '$lib/api/bist/v1/bist_pb';
	import { Cell as RowPicker } from '$lib/components/custom/data-table/data-table-row-pickers';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { formatCapacity, formatSecond, formatTimeAgo } from '$lib/formatter';
	import { timestampDate } from '@bufbuild/protobuf/wkt';
	import Icon from '@iconify/svelte';
	import type { Row } from '@tanstack/table-core';
	import Actions from './actions.svelte';

	export const cells = {
		row_picker,
		name,
		status,
		target,
		operation,
		duration,
		objectSize, 
		objectCount,
		throughputFastest,
		throughputSlowest,
		throughputMedian,
		createdBy,
		startedAt,
		completedAt,
		actions,
	};
</script>

{#snippet row_picker(row: Row<TestResult>)}
	<RowPicker {row} />
{/snippet}

{#snippet name(row: Row<TestResult>)}
		{row.original.name}
{/snippet}

{#snippet status(row: Row<TestResult>)}
		{#if TestResult_Status[row.original.status] === 'SUCCEEDED'}
			<Icon icon="ph:check" />
		{:else if TestResult_Status[row.original.status] === 'FAILED'}
			<Icon icon="ph:x" />
		{:else}
			<Icon icon="svg-spinners:180-ring-with-bg" />
		{/if}
{/snippet}

{#snippet target(row: Row<TestResult>)}
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
{/snippet}

{#snippet createdBy(row: Row<TestResult>)}
	{row.original.createdBy}
{/snippet}

{#snippet operation(row: Row<TestResult>)}
		{#if row.original.kind.case === 'warp' &&  row.original.kind.value?.input}
			{Warp_Input_Operation[row.original.kind.value.input.operation]}
        {/if}
{/snippet}
{#snippet duration(row: Row<TestResult>)}
	<span class="flex justify-end">
		{#if row.original.kind.case === 'warp' &&  row.original.kind.value?.input}
			{@const formatted = formatSecond(Number(row.original.kind.value?.input.durationSeconds))}
			{formatted.value} {formatted.unit}
        {/if}
	</span>
{/snippet}
{#snippet objectSize(row: Row<TestResult>)}
	<span class="flex justify-end">
		{#if row.original.kind.case === 'warp' &&  row.original.kind.value?.input}
            {@const formatted = formatCapacity(Number(row.original.kind.value?.input.objectSizeBytes))}
            {formatted.value} {formatted.unit}
        {/if}
	</span>
{/snippet}
{#snippet objectCount(row: Row<TestResult>)}
	<span class="flex justify-end">
		{#if row.original.kind.case === 'warp' &&  row.original.kind.value?.input}
			{row.original.kind.value.input.objectCount}
        {/if}
	</span>
{/snippet}

{#snippet throughputFastest(row: Row<TestResult>)}
	<div class="flex flex-col gap-1 items-end">
		{#if row.original.kind.case === 'warp' && row.original.kind.value?.output?.get?.bytes}
			<Badge variant="outline">
				GET {(Number(row.original.kind.value.output.get.bytes.fastestPerSecond) / 1000000).toFixed(3)} MB/s
			</Badge>
		{/if}
		{#if row.original.kind.case === 'warp' && row.original.kind.value?.output?.put?.bytes}
			<Badge variant="default">
				PUT {(Number(row.original.kind.value.output.put.bytes.fastestPerSecond) / 1000000).toFixed(3)} MB/s
			</Badge>
		{/if}
		{#if row.original.kind.case === 'warp' && row.original.kind.value?.output?.delete?.bytes}
			<Badge variant="destructive">
				DELETE {(Number(row.original.kind.value.output.delete.bytes.fastestPerSecond) / 1000000).toFixed(3)} MB/s
			</Badge>
		{/if}
	</div>
{/snippet}

{#snippet throughputSlowest(row: Row<TestResult>)}
	<div class="flex flex-col gap-1 items-end">
		{#if row.original.kind.case === 'warp' && row.original.kind.value?.output?.get?.bytes}
			<Badge variant="outline">
				GET {(Number(row.original.kind.value.output.get.bytes.slowestPerSecond) / 1000000).toFixed(3)} MB/s
			</Badge>
		{/if}
		{#if row.original.kind.case === 'warp' && row.original.kind.value?.output?.put?.bytes}
			<Badge variant="default">
				PUT {(Number(row.original.kind.value.output.put.bytes.slowestPerSecond) / 1000000).toFixed(3)} MB/s
			</Badge>
		{/if}
		{#if row.original.kind.case === 'warp' && row.original.kind.value?.output?.delete?.bytes}
			<Badge variant="destructive">
				DELETE {(Number(row.original.kind.value.output.delete.bytes.slowestPerSecond) / 1000000).toFixed(3)} MB/s
			</Badge>
		{/if}
	</div>
{/snippet}

{#snippet throughputMedian(row: Row<TestResult>)}
	<div class="flex flex-col gap-1 items-end">
		{#if row.original.kind.case === 'warp' && row.original.kind.value?.output?.get?.bytes}
			<Badge variant="outline">
				GET {(Number(row.original.kind.value.output.get.bytes.medianPerSecond) / 1000000).toFixed(3)} MB/s
			</Badge>
		{/if}
		{#if row.original.kind.case === 'warp' && row.original.kind.value?.output?.put?.bytes}
			<Badge variant="default">
				PUT {(Number(row.original.kind.value.output.put.bytes.medianPerSecond) / 1000000).toFixed(3)} MB/s
			</Badge>
		{/if}
		{#if row.original.kind.case === 'warp' && row.original.kind.value?.output?.delete?.bytes}
			<Badge variant="destructive">
				DELETE {(Number(row.original.kind.value.output.delete.bytes.medianPerSecond) / 1000000).toFixed(3)} MB/s
			</Badge>
		{/if}
	</div>
{/snippet}

{#snippet startedAt(row: Row<TestResult>)}
	{#if row.original.startedAt}
		<Tooltip.Provider>
			<Tooltip.Root>
				<Tooltip.Trigger>
					{formatTimeAgo(timestampDate(row.original.startedAt))}
				</Tooltip.Trigger>
				<Tooltip.Content>
					{timestampDate(row.original.startedAt)}
				</Tooltip.Content>
			</Tooltip.Root>
		</Tooltip.Provider>
	{/if}
{/snippet}

{#snippet completedAt(row: Row<TestResult>)}
	{#if row.original.completedAt}
		{#if Number(timestampDate(row.original.completedAt)) >= 0}
			<Tooltip.Provider>
				<Tooltip.Root>
					<Tooltip.Trigger>
						{formatTimeAgo(timestampDate(row.original.completedAt))}
					</Tooltip.Trigger>
					<Tooltip.Content>
						{timestampDate(row.original.completedAt)}
					</Tooltip.Content>
				</Tooltip.Root>
			</Tooltip.Provider>
		{/if}
	{/if}
{/snippet}

{#snippet actions(row: Row<TestResult>)}
	<Actions testResult={row.original}/>
{/snippet}
