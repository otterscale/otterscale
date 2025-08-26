<script lang="ts" module>
	import { type TestResult } from '$lib/api/bist/v1/bist_pb';
	import { Headers, Sorter } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { m } from '$lib/paraglide/messages';
	import type { Column, Table } from '@tanstack/table-core';

	export const messages = {
		name: m.name(),
		status: m.status(),
		target: m.target(),
		accessMode: m.access_mode(),
		jobCount: m.job_count(),
		runTime: m.run_time(),
		blockSize: m.block_size(),
		fileSize: m.file_size(),
		ioDepth: m.io_depth(),
		bandwidth: m.bandwidth(),
		iops: m.iops(),
		latencyMinimum: m.latency_minimum(),
		latencyMaximum: m.latency_maximum(),
		latencyMean: m.latency_mean(),
		createdBy: m.created_by(),
		startedAt: m.started_at(),
		completedAt: m.completed_at()
	};

	export const headers = {
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

{#snippet row_picker(table: Table<TestResult>)}
	<Layout.Header class="justify-center">
		<Layout.HeaderController>
			<Headers.RowPicker {table} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet name(column: Column<TestResult>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.name()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet status(column: Column<TestResult>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.status()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet target(column: Column<TestResult>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.target()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet accessMode(column: Column<TestResult>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.access_mode()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet jobCount(column: Column<TestResult>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.job_count()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet runTime(column: Column<TestResult>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.run_time()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet blockSize(column: Column<TestResult>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.block_size()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet fileSize(column: Column<TestResult>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.file_size()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet ioDepth(column: Column<TestResult>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.io_depth()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet createdBy(column: Column<TestResult>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.created_by()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet bandwidth(column: Column<TestResult>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.bandwidth()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet iops(column: Column<TestResult>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderViewer>{m.iops()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet latencyMin(column: Column<TestResult>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderViewer>
			{m.latency()}
			<p class="text-muted-foreground">{m.minimum()}</p>
		</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet latencyMax(column: Column<TestResult>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderViewer>
			{m.latency()}
			<p class="text-muted-foreground">{m.maximum()}</p>
		</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet latencyMean(column: Column<TestResult>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderViewer>
			{m.latency()}
			<p class="text-muted-foreground">{m.mean()}</p>
		</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet startedAt(column: Column<TestResult>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.started_at()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet completedAt(column: Column<TestResult>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.completed_at()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet actions(column: Column<TestResult>)}{/snippet}
