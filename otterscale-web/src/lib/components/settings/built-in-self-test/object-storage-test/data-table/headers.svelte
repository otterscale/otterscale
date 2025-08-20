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

{#snippet createdBy(column: Column<TestResult>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>CREATOR</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet operation(column: Column<TestResult>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>OPERATION</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet duration(column: Column<TestResult>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>DURATION</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet objectSize(column: Column<TestResult>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>OBJECT SIZE</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet objectCount(column: Column<TestResult>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>OBJECT COUNT</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet throughputFastest(column: Column<TestResult>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderViewer>
			THROUGHPUT
			<p class="text-muted-foreground">FASTEST</p>
		</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet throughputSlowest(column: Column<TestResult>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderViewer>
			THROUGHPUT
			<p class="text-muted-foreground">SLOWEST</p>
		</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet throughputMedian(column: Column<TestResult>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderViewer>
			THROUGHPUT
			<p class="text-muted-foreground">MEDIAN</p>
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

{#snippet actions(column: Column<TestResult>)}{/snippet}
