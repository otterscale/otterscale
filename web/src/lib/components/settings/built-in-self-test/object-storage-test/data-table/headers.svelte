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

{#snippet createdBy(column: Column<TestResult>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.created_by()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet operation(column: Column<TestResult>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.operation()}</Layout.HeaderViewer>
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
		<Layout.HeaderViewer>{m.duration()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet objectSize(column: Column<TestResult>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.object_size()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet objectCount(column: Column<TestResult>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.object_count()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet throughputFastest()}
	<Layout.Header class="justify-end">
		<Layout.HeaderViewer>
			{m.throughput()}
			<p class="text-muted-foreground">{m.fastest()}</p>
		</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet throughputSlowest()}
	<Layout.Header class="justify-end">
		<Layout.HeaderViewer>
			{m.throughput()}
			<p class="text-muted-foreground">{m.slowest()}</p>
		</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet throughputMedian()}
	<Layout.Header class="justify-end">
		<Layout.HeaderViewer>
			{m.throughput()}
			<p class="text-muted-foreground">{m.median()}</p>
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
