<script lang="ts" module>
	import {
		FIO_Input_AccessMode,
		type TestResult,
		TestResult_Status
	} from '$lib/api/bist/v1/bist_pb';
	import { Cells } from '$lib/components/custom/data-table/core';
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
		accessMode,
		jobCount,
		runTime,
		blockSize,
		fileSize,
		ioDepth,
		bandwidth,
		iops,
		latencyMin,
		latencyMax,
		latencyMean,
		createdBy,
		startedAt,
		completedAt,
		actions
	};
</script>

{#snippet row_picker(row: Row<TestResult>)}
	<Cells.RowPicker {row} />
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
	{#if row.original.kind.case === 'fio' && row.original.kind.value?.input}
		{#if row.original.kind.value.target.case === 'cephBlockDevice'}
			<Badge variant="outline">
				{row.original.kind.value.target.value.facilityName}
			</Badge>
		{:else if row.original.kind.value.target.case === 'networkFileSystem'}
			<Badge variant="outline">
				{row.original.kind.value.target.value.endpoint}
			</Badge>
		{/if}
	{/if}
{/snippet}

{#snippet accessMode(row: Row<TestResult>)}
	{#if row.original.kind.case === 'fio' && row.original.kind.value?.input}
		{FIO_Input_AccessMode[row.original.kind.value?.input.accessMode]}
	{/if}
{/snippet}

{#snippet jobCount(row: Row<TestResult>)}
	<div class="flex justify-end">
		{#if row.original.kind.case === 'fio' && row.original.kind.value?.input}
			{row.original.kind.value?.input.jobCount}
		{/if}
	</div>
{/snippet}

{#snippet runTime(row: Row<TestResult>)}
	{#if row.original.kind.case === 'fio' && row.original.kind.value?.input}
		{@const formatted = formatSecond(Number(row.original.kind.value?.input.runTimeSeconds))}
		{formatted.value}
		{formatted.unit}
	{/if}
{/snippet}

{#snippet blockSize(row: Row<TestResult>)}
	<div class="flex justify-end">
		{#if row.original.kind.case === 'fio' && row.original.kind.value?.input}
			{@const formatted = formatCapacity(Number(row.original.kind.value?.input.blockSizeBytes))}
			{formatted.value}
			{formatted.unit}
		{/if}
	</div>
{/snippet}

{#snippet fileSize(row: Row<TestResult>)}
	<div class="flex justify-end">
		{#if row.original.kind.case === 'fio' && row.original.kind.value?.input}
			{@const formatted = formatCapacity(Number(row.original.kind.value?.input.fileSizeBytes))}
			{formatted.value}
			{formatted.unit}
		{/if}
	</div>
{/snippet}

{#snippet ioDepth(row: Row<TestResult>)}
	<div class="flex justify-end">
		{#if row.original.kind.case === 'fio' && row.original.kind.value?.input}
			{row.original.kind.value?.input.ioDepth}
		{/if}
	</div>
{/snippet}

{#snippet createdBy(row: Row<TestResult>)}
	{row.original.createdBy}
{/snippet}

{#snippet bandwidth(row: Row<TestResult>)}
	<div class="flex flex-col items-end gap-1">
		{#if row.original.kind.case === 'fio' && row.original.kind.value?.output?.read}
			<Badge variant="outline">
				Read {(Number(row.original.kind.value.output.read.bandwidthBytes) / 1024 / 1024).toFixed(2)}
				MB/s
			</Badge>
		{/if}
		{#if row.original.kind.case === 'fio' && row.original.kind.value?.output?.write}
			<Badge variant="default">
				Write {(Number(row.original.kind.value.output.write.bandwidthBytes) / 1024 / 1024).toFixed(
					2
				)} MB/s
			</Badge>
		{/if}
		{#if row.original.kind.case === 'fio' && row.original.kind.value?.output?.trim}
			<Badge variant="secondary">
				Trim {(Number(row.original.kind.value.output.trim.bandwidthBytes) / 1024 / 1024).toFixed(2)}
				MB/s
			</Badge>
		{/if}
	</div>
{/snippet}

{#snippet iops(row: Row<TestResult>)}
	<div class="flex flex-col items-end gap-1">
		{#if row.original.kind.case === 'fio' && row.original.kind.value?.output?.read}
			<Badge variant="outline">
				Read: {row.original.kind.value.output.read.ioPerSecond.toFixed(0)}
			</Badge>
		{/if}
		{#if row.original.kind.case === 'fio' && row.original.kind.value?.output?.write}
			<Badge variant="default">
				Write: {row.original.kind.value.output.write.ioPerSecond.toFixed(0)}
			</Badge>
		{/if}
		{#if row.original.kind.case === 'fio' && row.original.kind.value?.output?.trim}
			<Badge variant="secondary">
				Trim: {row.original.kind.value.output.trim.ioPerSecond.toFixed(0)}
			</Badge>
		{/if}
	</div>
{/snippet}

{#snippet latencyMin(row: Row<TestResult>)}
	<div class="flex flex-col items-end gap-1">
		{#if row.original.kind.case === 'fio' && row.original.kind.value?.output?.read?.latency}
			<Badge variant="outline">
				Read: {(
					Number(row.original.kind.value.output.read.latency.minNanoseconds) / 1000000
				).toFixed(3)} ms
			</Badge>
		{/if}
		{#if row.original.kind.case === 'fio' && row.original.kind.value?.output?.write?.latency}
			<Badge variant="default">
				Write: {(
					Number(row.original.kind.value.output.write.latency.minNanoseconds) / 1000000
				).toFixed(3)} ms
			</Badge>
		{/if}
		{#if row.original.kind.case === 'fio' && row.original.kind.value?.output?.trim?.latency}
			<Badge variant="secondary">
				Trim: {(
					Number(row.original.kind.value.output.trim.latency.minNanoseconds) / 1000000
				).toFixed(3)} ms
			</Badge>
		{/if}
	</div>
{/snippet}

{#snippet latencyMax(row: Row<TestResult>)}
	<div class="flex flex-col items-end gap-1">
		{#if row.original.kind.case === 'fio' && row.original.kind.value?.output?.read?.latency}
			<Badge variant="outline">
				Read: {(
					Number(row.original.kind.value.output.read.latency.maxNanoseconds) / 1000000
				).toFixed(3)} ms
			</Badge>
		{/if}
		{#if row.original.kind.case === 'fio' && row.original.kind.value?.output?.write?.latency}
			<Badge variant="default">
				Write: {(
					Number(row.original.kind.value.output.write.latency.maxNanoseconds) / 1000000
				).toFixed(3)} ms
			</Badge>
		{/if}
		{#if row.original.kind.case === 'fio' && row.original.kind.value?.output?.trim?.latency}
			<Badge variant="secondary">
				Trim: {(
					Number(row.original.kind.value.output.trim.latency.maxNanoseconds) / 1000000
				).toFixed(3)} ms
			</Badge>
		{/if}
	</div>
{/snippet}

{#snippet latencyMean(row: Row<TestResult>)}
	<div class="flex flex-col items-end gap-1">
		{#if row.original.kind.case === 'fio' && row.original.kind.value?.output?.read?.latency}
			<Badge variant="outline">
				Read: {(
					Number(row.original.kind.value.output.read.latency.meanNanoseconds) / 1000000
				).toFixed(3)} ms
			</Badge>
		{/if}
		{#if row.original.kind.case === 'fio' && row.original.kind.value?.output?.write?.latency}
			<Badge variant="default">
				Write: {(
					Number(row.original.kind.value.output.write.latency.meanNanoseconds) / 1000000
				).toFixed(3)} ms
			</Badge>
		{/if}
		{#if row.original.kind.case === 'fio' && row.original.kind.value?.output?.trim?.latency}
			<Badge variant="secondary">
				Trim: {(
					Number(row.original.kind.value.output.trim.latency.meanNanoseconds) / 1000000
				).toFixed(3)} ms
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
	{#if row.original.completedAt && Number(timestampDate(row.original.completedAt)) >= 0}
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
{/snippet}

{#snippet actions(row: Row<TestResult>)}
	<Actions testResult={row.original} />
{/snippet}
