<script lang="ts" module>
	import { timestampDate } from '@bufbuild/protobuf/wkt';
	import type { Row } from '@tanstack/table-core';

	import type { Job } from '$lib/api/application/v1/application_pb';
	import { Cells } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { formatTimeAgo } from '$lib/formatter';

	export const cells = {
		row_picker,
		name,
		namespace,
		active,
		ready,
		succeeded_or_failed,
		terminating,
		lastCondition,
		startedAt,
		completedAt
	};

	export function getJobStatus(job: Job) {
		if (job.terminating) {
			return 'Terminating';
		}

		if (job.active > 0) {
			return 'Active';
		}

		if (job.succeeded > 0) {
			return 'Succeeded';
		}

		if (job.failed > 0) {
			return 'Failed';
		}

		return 'Pending';
	}
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

{#snippet succeeded_or_failed(row: Row<Job>)}
	<Layout.Cell class="items-end">
		{row.original.succeeded}/{row.original.failed}
	</Layout.Cell>
{/snippet}

{#snippet terminating(row: Row<Job>)}
	<Layout.Cell class="items-end">
		{row.original.terminating}
	</Layout.Cell>
{/snippet}

{#snippet lastCondition(row: Row<Job>)}
	<Layout.Cell class="items-start">
		{#if row.original.lastCondition && row.original.lastCondition.type === 'Failed'}
			<div class="space-y-1">
				<h4 class="text-destructive">
					{getJobStatus(row.original)}:
					{row.original.lastCondition.reason}
				</h4>
				<div class="flex gap-1">
					<Tooltip.Provider>
						<Tooltip.Root>
							<Tooltip.Trigger>
								<p class="max-w-50 truncate text-muted-foreground">
									{row.original.lastCondition.message}
								</p>
							</Tooltip.Trigger>
							<Tooltip.Content>
								{row.original.lastCondition.message}
							</Tooltip.Content>
						</Tooltip.Root>
					</Tooltip.Provider>
				</div>
			</div>
		{:else}
			{getJobStatus(row.original)}
		{/if}
	</Layout.Cell>
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
