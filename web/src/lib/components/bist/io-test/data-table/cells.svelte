<script lang="ts" module>
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import type { Row } from '@tanstack/table-core';
	import { formatTimeAgo } from '$lib/formatter';
	import { type TestResult, TestResult_Status, FIO_Input_AccessMode } from '$gen/api/bist/v1/bist_pb'
	import { timestampDate } from '@bufbuild/protobuf/wkt';
	import Icon from '@iconify/svelte';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { formatByte, formatSecond } from '$lib/formatter';

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
		bandwidth: bandwidth,
		iops: iops,
		latencyMin: latencyMin,
		latencyMax: latencyMax,
		latencyMean: latencyMean,
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
			{@const formatted = formatSecond(Number(row.original.kind.value?.input.runTime))}
            {formatted.value} {formatted.unit}
        {/if}
	</p>
{/snippet}

{#snippet blockSize(row: Row<TestResult>)}
	<p>
		{#if row.original.kind.case === 'fio' &&  row.original.kind.value?.input}
            {@const formatted = formatByte(Number(row.original.kind.value?.input.blockSize))}
            {formatted.value} {formatted.unit}
        {/if}
	</p>
{/snippet}

{#snippet fileSize(row: Row<TestResult>)}
	<p>
		{#if row.original.kind.case === 'fio' &&  row.original.kind.value?.input}
            {@const formatted = formatByte(Number(row.original.kind.value?.input.fileSize))}
            {formatted.value} {formatted.unit}
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

{#snippet bandwidth(row: Row<TestResult>)}
	{#if row.original.kind.case === 'fio' && row.original.kind.value?.output?.read}
		<p>
			<Badge variant="outline">
				R: {(Number(row.original.kind.value.output.read.bandwidthBytes) / 1024 / 1024).toFixed(2)} MB/s
			</Badge>
		</p>
	{/if}
	{#if row.original.kind.case === 'fio' && row.original.kind.value?.output?.write}
		<p class={row.original.kind.value?.output?.read ? 'mt-1' : ''}>
			<Badge variant="outline">
				W: {(Number(row.original.kind.value.output.write.bandwidthBytes) / 1024 / 1024).toFixed(2)} MB/s
			</Badge>
		</p>
	{/if}
	{#if row.original.kind.case === 'fio' && row.original.kind.value?.output?.trim}
		<p class={(row.original.kind.value?.output?.read || row.original.kind.value?.output?.write) ? 'mt-1' : ''}>
			<Badge variant="outline">
				T: {(Number(row.original.kind.value.output.trim.bandwidthBytes) / 1024 / 1024).toFixed(2)} MB/s
			</Badge>
		</p>
	{/if}
{/snippet}

{#snippet iops(row: Row<TestResult>)}
	{#if row.original.kind.case === 'fio' && row.original.kind.value?.output?.read}
		<p>
			<Badge variant="outline">
				R: {row.original.kind.value.output.read.ioPerSecond.toFixed(0)}
			</Badge>	
		</p>
	{/if}
	{#if row.original.kind.case === 'fio' && row.original.kind.value?.output?.write}
		<p class={row.original.kind.value?.output?.read ? 'mt-1' : ''}>
			<Badge variant="outline">
				W: {row.original.kind.value.output.write.ioPerSecond.toFixed(0)}
			</Badge>
		</p>
	{/if}
	{#if row.original.kind.case === 'fio' && row.original.kind.value?.output?.trim}
		<p class={(row.original.kind.value?.output?.read || row.original.kind.value?.output?.write) ? 'mt-1' : ''}>
			<Badge variant="outline">
				T: {row.original.kind.value.output.trim.ioPerSecond.toFixed(0)}
			</Badge>
		</p>
	{/if}
{/snippet}

{#snippet latencyMin(row: Row<TestResult>)}
	{#if row.original.kind.case === 'fio' && row.original.kind.value?.output?.read?.latency}
		<p>
			<Badge variant="outline">
				R: 
				{(Number(row.original.kind.value.output.read.latency.minNanoseconds) / 1000000).toFixed(3)} ms
			</Badge>
		</p>
	{/if}
	{#if row.original.kind.case === 'fio' && row.original.kind.value?.output?.write?.latency}
		<p class={row.original.kind.value?.output?.read ? 'mt-1' : ''}>
			<Badge variant="outline">
				W: 
				{(Number(row.original.kind.value.output.write.latency.minNanoseconds) / 1000000).toFixed(3)} ms
			</Badge>
		</p>
	{/if}
	{#if row.original.kind.case === 'fio' && row.original.kind.value?.output?.trim?.latency}
		<p class={(row.original.kind.value?.output?.read || row.original.kind.value?.output?.write) ? 'mt-1' : ''}>
			<Badge variant="outline">
				T: 
				{(Number(row.original.kind.value.output.trim.latency.minNanoseconds) / 1000000).toFixed(3)} ms
			</Badge>
		</p>
	{/if}
{/snippet}

{#snippet latencyMax(row: Row<TestResult>)}
	{#if row.original.kind.case === 'fio' && row.original.kind.value?.output?.read?.latency}
		<p>
			<Badge variant="outline">
				R: 
				{(Number(row.original.kind.value.output.read.latency.maxNanoseconds) / 1000000).toFixed(3)} ms
			</Badge>
		</p>
	{/if}
	{#if row.original.kind.case === 'fio' && row.original.kind.value?.output?.write?.latency}
		<p class={row.original.kind.value?.output?.read ? 'mt-1' : ''}>
			<Badge variant="outline">
				W: 
				{(Number(row.original.kind.value.output.write.latency.maxNanoseconds) / 1000000).toFixed(3)} ms
			</Badge>
		</p>
	{/if}
	{#if row.original.kind.case === 'fio' && row.original.kind.value?.output?.trim?.latency}
		<p class={(row.original.kind.value?.output?.read || row.original.kind.value?.output?.write) ? 'mt-1' : ''}>
			<Badge variant="outline">
				T: 
				{(Number(row.original.kind.value.output.trim.latency.maxNanoseconds) / 1000000).toFixed(3)} ms
			</Badge>
		</p>
	{/if}
{/snippet}


{#snippet latencyMean(row: Row<TestResult>)}
	{#if row.original.kind.case === 'fio' && row.original.kind.value?.output?.read?.latency}
		<p>
			<Badge variant="outline">
				R: 
				{(Number(row.original.kind.value.output.read.latency.meanNanoseconds) / 1000000).toFixed(3)} ms
			</Badge>
		</p>
	{/if}
	{#if row.original.kind.case === 'fio' && row.original.kind.value?.output?.write?.latency}
		<p class={row.original.kind.value?.output?.read ? 'mt-1' : ''}>
			<Badge variant="outline">
				W: 
				{(Number(row.original.kind.value.output.write.latency.meanNanoseconds) / 1000000).toFixed(3)} ms
			</Badge>
		</p>
	{/if}
	{#if row.original.kind.case === 'fio' && row.original.kind.value?.output?.trim?.latency}
		<p class={(row.original.kind.value?.output?.read || row.original.kind.value?.output?.write) ? 'mt-1' : ''}>
			<Badge variant="outline">
				T: 
				{(Number(row.original.kind.value.output.trim.latency.meanNanoseconds) / 1000000).toFixed(3)} ms
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
