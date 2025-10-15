<script lang="ts" module>
	import type { Column, Table } from '@tanstack/table-core';

	import { type TestResult } from '$lib/api/configuration/v1/configuration_pb';
	import { Headers, Sorter } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { m } from '$lib/paraglide/messages';

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
		latencyMinimum,
		latencyMaximum,
		latencyMean,
		createdBy,
		startedAt,
		completedAt,
		actions,
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

{#snippet target()}
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

{#snippet iops()}
	<Layout.Header class="justify-end">
		<Layout.HeaderViewer>{m.iops()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet latencyMinimum()}
	<Layout.Header class="justify-end">
		<Layout.HeaderViewer>
			{m.latency()}
			<p class="text-muted-foreground">{m.minimum()}</p>
		</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet latencyMaximum()}
	<Layout.Header class="justify-end">
		<Layout.HeaderViewer>
			{m.latency()}
			<p class="text-muted-foreground">{m.maximum()}</p>
		</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet latencyMean()}
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

{#snippet actions()}{/snippet}
