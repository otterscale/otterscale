<script lang="ts" module>
	import { type TestResult } from '$lib/api/bist/v1/bist_pb';
	import * as Layout from '$lib/components/custom/data-table/data-table-layout';
	import { Header as RowPicker } from '$lib/components/custom/data-table/data-table-row-pickers';
	import Sorter from '$lib/components/custom/data-table/data-table-sorter.svelte';
	import type { Column, Table } from '@tanstack/table-core';

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
			<RowPicker {table} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet name(column: Column<TestResult>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>NAME</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet status(column: Column<TestResult>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>STATUS</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet target(column: Column<TestResult>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>TARGET</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet accessMode(column: Column<TestResult>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>ACCESS MODE</Layout.HeaderViewer>
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
		<Layout.HeaderViewer>JOB COUNT</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet runTime(column: Column<TestResult>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>RUN TIME</Layout.HeaderViewer>
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
		<Layout.HeaderViewer>BLOCK SIZE</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet fileSize(column: Column<TestResult>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>FILE SIZE</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet ioDepth(column: Column<TestResult>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>IO DEPTH</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet createdBy(column: Column<TestResult>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>CREATOR</Layout.HeaderViewer>
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
		<Layout.HeaderViewer>BANDWIDTH</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet iops(column: Column<TestResult>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderViewer>IOPS</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet latencyMin(column: Column<TestResult>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderViewer>
			LATENCY
			<p class="text-muted-foreground">
				MINIMUM
			</p>
		</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet latencyMax(column: Column<TestResult>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderViewer>
			LATENCY
			<p class="text-muted-foreground">
				MAXIMUM
			</p>
		</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet latencyMean(column: Column<TestResult>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderViewer>
			LATENCY
			<p class="text-muted-foreground">
				MEAN
			</p>
		</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet startedAt(column: Column<TestResult>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>STARTED</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet completedAt(column: Column<TestResult>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>COMPLETED</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet actions(column: Column<TestResult>)}
{/snippet}