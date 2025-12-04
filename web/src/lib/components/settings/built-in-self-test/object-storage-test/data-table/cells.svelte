<script lang="ts" module>
	import { timestampDate } from '@bufbuild/protobuf/wkt';
	import Icon from '@iconify/svelte';
	import type { Row } from '@tanstack/table-core';

	import {
		InternalObjectService_Type,
		type TestResult,
		TestResult_Status,
		Warp_Input_Operation
	} from '$lib/api/configuration/v1/configuration_pb';
	import { Cells } from '$lib/components/custom/data-table/core';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Table from '$lib/components/custom/table/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { formatCapacity, formatSecond, formatTimeAgo } from '$lib/formatter';

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
		actions
	};
</script>

{#snippet row_picker(row: Row<TestResult>)}
	<Table.Cell alignClass="items-center">
		<Cells.RowPicker {row} />
	</Table.Cell>
{/snippet}

{#snippet name(row: Row<TestResult>)}
	<Table.Cell alignClass="items-start">
		{row.original.name}
	</Table.Cell>
{/snippet}

{#snippet status(row: Row<TestResult>)}
	<Table.Cell alignClass="items-start">
		{#if TestResult_Status[row.original.status] === 'SUCCEEDED'}
			<Icon icon="ph:check" />
		{:else if TestResult_Status[row.original.status] === 'FAILED'}
			<Icon icon="ph:x" />
		{:else}
			<Icon icon="svg-spinners:180-ring-with-bg" />
		{/if}
	</Table.Cell>
{/snippet}

{#snippet target(row: Row<TestResult>)}
	<Table.Cell alignClass="items-start">
		{#if row.original.kind.case === 'warp' && row.original.kind.value?.input}
			{#if row.original.kind.value.target.case === 'internalObjectService'}
				<Badge variant="outline">
					{InternalObjectService_Type[row.original.kind.value.target.value.type]}
				</Badge>
			{:else if row.original.kind.value.target.case === 'externalObjectService'}
				<Badge variant="outline">
					{row.original.kind.value.target.value.host}
				</Badge>
			{/if}
		{/if}
	</Table.Cell>
{/snippet}

{#snippet createdBy(row: Row<TestResult>)}
	<Table.Cell alignClass="items-start">
		{row.original.createdBy}
	</Table.Cell>
{/snippet}

{#snippet operation(row: Row<TestResult>)}
	<Table.Cell alignClass="items-start">
		{#if row.original.kind.case === 'warp' && row.original.kind.value?.input}
			{Warp_Input_Operation[row.original.kind.value.input.operation]}
		{/if}
	</Table.Cell>
{/snippet}
{#snippet duration(row: Row<TestResult>)}
	<Table.Cell alignClass="items-end">
		{#if row.original.kind.case === 'warp' && row.original.kind.value?.input}
			{@const formatted = formatSecond(Number(row.original.kind.value?.input.durationSeconds))}
			{formatted.value}
			{formatted.unit}
		{/if}
	</Table.Cell>
{/snippet}
{#snippet objectSize(row: Row<TestResult>)}
	<Table.Cell alignClass="items-end">
		{#if row.original.kind.case === 'warp' && row.original.kind.value?.input}
			{@const formatted = formatCapacity(Number(row.original.kind.value?.input.objectSizeBytes))}
			{formatted.value}
			{formatted.unit}
		{/if}
	</Table.Cell>
{/snippet}
{#snippet objectCount(row: Row<TestResult>)}
	<Table.Cell alignClass="items-end">
		{#if row.original.kind.case === 'warp' && row.original.kind.value?.input}
			{row.original.kind.value.input.objectCount}
		{/if}
	</Table.Cell>
{/snippet}

{#snippet throughputFastest(row: Row<TestResult>)}
	<Table.Cell alignClass="items-end">
		{#if row.original.kind.case === 'warp' && row.original.kind.value?.output?.get?.bytes}
			<Badge variant="outline">
				Get {(Number(row.original.kind.value.output.get.bytes.fastestPerSecond) / 1000000).toFixed(
					3
				)} MB/s
			</Badge>
		{/if}
		{#if row.original.kind.case === 'warp' && row.original.kind.value?.output?.put?.bytes}
			<Badge variant="default">
				Put {(Number(row.original.kind.value.output.put.bytes.fastestPerSecond) / 1000000).toFixed(
					3
				)} MB/s
			</Badge>
		{/if}
		{#if row.original.kind.case === 'warp' && row.original.kind.value?.output?.delete?.bytes}
			<Badge variant="destructive">
				Delete {(
					Number(row.original.kind.value.output.delete.bytes.fastestPerSecond) / 1000000
				).toFixed(3)} MB/s
			</Badge>
		{/if}
	</Table.Cell>
{/snippet}

{#snippet throughputSlowest(row: Row<TestResult>)}
	<Table.Cell alignClass="items-end">
		{#if row.original.kind.case === 'warp' && row.original.kind.value?.output?.get?.bytes}
			<Badge variant="outline">
				Get {(Number(row.original.kind.value.output.get.bytes.slowestPerSecond) / 1000000).toFixed(
					3
				)} MB/s
			</Badge>
		{/if}
		{#if row.original.kind.case === 'warp' && row.original.kind.value?.output?.put?.bytes}
			<Badge variant="default">
				Put {(Number(row.original.kind.value.output.put.bytes.slowestPerSecond) / 1000000).toFixed(
					3
				)} MB/s
			</Badge>
		{/if}
		{#if row.original.kind.case === 'warp' && row.original.kind.value?.output?.delete?.bytes}
			<Badge variant="destructive">
				Delete {(
					Number(row.original.kind.value.output.delete.bytes.slowestPerSecond) / 1000000
				).toFixed(3)} MB/s
			</Badge>
		{/if}
	</Table.Cell>
{/snippet}

{#snippet throughputMedian(row: Row<TestResult>)}
	<Table.Cell alignClass="items-end">
		{#if row.original.kind.case === 'warp' && row.original.kind.value?.output?.get?.bytes}
			<Badge variant="outline">
				Get {(Number(row.original.kind.value.output.get.bytes.medianPerSecond) / 1000000).toFixed(
					3
				)} MB/s
			</Badge>
		{/if}
		{#if row.original.kind.case === 'warp' && row.original.kind.value?.output?.put?.bytes}
			<Badge variant="default">
				Put {(Number(row.original.kind.value.output.put.bytes.medianPerSecond) / 1000000).toFixed(
					3
				)} MB/s
			</Badge>
		{/if}
		{#if row.original.kind.case === 'warp' && row.original.kind.value?.output?.delete?.bytes}
			<Badge variant="destructive">
				Delete {(
					Number(row.original.kind.value.output.delete.bytes.medianPerSecond) / 1000000
				).toFixed(3)} MB/s
			</Badge>
		{/if}
	</Table.Cell>
{/snippet}

{#snippet startedAt(row: Row<TestResult>)}
	<Table.Cell alignClass="items-start">
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
	</Table.Cell>
{/snippet}

{#snippet completedAt(row: Row<TestResult>)}
	<Table.Cell alignClass="items-start">
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
	</Table.Cell>
{/snippet}

{#snippet actions(data: { row: Row<TestResult>; scope: string; reloadManager: ReloadManager })}
	<Table.Cell alignClass="items-start">
		<Actions testResult={data.row.original} scope={data.scope} reloadManager={data.reloadManager} />
	</Table.Cell>
{/snippet}
