<script lang="ts" module>
	import { timestampDate } from '@bufbuild/protobuf/wkt';
	import type { Row } from '@tanstack/table-core';

	import type { Job } from '$lib/api/application/v1/application_pb';
	import { Cells } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import { Spinner } from '$lib/components/ui/spinner';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { formatTimeAgo } from '$lib/formatter';

	import { getJobStatus } from '../utils';

	export const cells = {
		row_picker,
		name,
		namespace,
		active,
		ready,
		succeeded,
		failed,
		terminating,
		status,
		conditions,
		startedAt,
		completedAt
	};
</script>

{#snippet row_picker(row: Row<Job>)}
	<Layout.Cell class="items-center">
		<Cells.RowPicker {row} />
	</Layout.Cell>
{/snippet}

{#snippet name(row: Row<Job>)}
	<Layout.Cell class="items-start">
		{row.original.name}
	</Layout.Cell>
{/snippet}

{#snippet namespace(row: Row<Job>)}
	<Layout.Cell class="items-start">
		{row.original.namespace}
	</Layout.Cell>
{/snippet}

{#snippet active(row: Row<Job>)}
	<Layout.Cell class="items-end">
		{row.original.active}
	</Layout.Cell>
{/snippet}

{#snippet ready(row: Row<Job>)}
	<Layout.Cell class="items-end">
		{row.original.ready}
	</Layout.Cell>
{/snippet}

{#snippet succeeded(row: Row<Job>)}
	<Layout.Cell class="items-end">
		{row.original.succeeded}
	</Layout.Cell>
{/snippet}

{#snippet failed(row: Row<Job>)}
	<Layout.Cell class="items-end">
		{row.original.failed}
	</Layout.Cell>
{/snippet}

{#snippet terminating(row: Row<Job>)}
	<Layout.Cell class="items-end">
		{row.original.terminating}
	</Layout.Cell>
{/snippet}

{#snippet status(row: Row<Job>)}
	<Layout.Cell class="items-start">
		{@const status = getJobStatus(row.original)}
		{#if status === 'Running'}
			<span class="flex items-center gap-1 text-muted-foreground">
				<Spinner />
				{status}
			</span>
		{:else if ['Failed', 'FailureTarget'].includes(status)}
			<p class="text-destructive">
				{status}
			</p>
		{:else}
			{status}
		{/if}
	</Layout.Cell>
{/snippet}

{#snippet conditions(row: Row<Job>)}
	{#if row.original.conditions}
		{@const trueConditions = row.original.conditions.filter(
			(condition) => condition.status === 'True'
		)}
		{#if trueConditions.length > 0}
			<Layout.Cell class="items-start">
				<div class="flex flex-wrap gap-1">
					{#each trueConditions as trueCondition, index (index)}
						<Tooltip.Provider>
							<Tooltip.Root>
								<Tooltip.Trigger>
									<Badge
										variant="outline"
										class={['Failed', 'FailureTarget'].includes(trueCondition.type)
											? 'border-destructive/50 text-destructive'
											: ''}
									>
										{trueCondition.type}
									</Badge>
								</Tooltip.Trigger>
								<Tooltip.Content>
									{#if trueCondition.message}
										{trueCondition.message}
									{:else}
										{trueCondition.type}
									{/if}
								</Tooltip.Content>
							</Tooltip.Root>
						</Tooltip.Provider>
					{/each}
				</div>
			</Layout.Cell>
		{/if}
	{/if}
{/snippet}

{#snippet startedAt(row: Row<Job>)}
	{#if row.original.startedAt}
		<Layout.Cell class="items-end">
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
		</Layout.Cell>
	{/if}
{/snippet}

{#snippet completedAt(row: Row<Job>)}
	{#if row.original.completedAt}
		<Layout.Cell class="items-end">
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
		</Layout.Cell>
	{/if}
{/snippet}
