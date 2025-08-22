<script lang="ts" module>
	import type { Subvolume } from '$lib/api/storage/v1/storage_pb';
	import { Cells } from '$lib/components/custom/data-table/core';
	import * as Progress from '$lib/components/custom/progress/index.js';
	import { Snapshot } from '$lib/components/storage/file-system/nfs/snapshot';
	import { Badge } from '$lib/components/ui/badge';
	import * as HoverCard from '$lib/components/ui/hover-card';
	import * as Table from '$lib/components/ui/table';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { formatCapacity, formatTimeAgo } from '$lib/formatter';
	import { timestampDate } from '@bufbuild/protobuf/wkt';
	import Icon from '@iconify/svelte';
	import type { Row } from '@tanstack/table-core';
	import { toast } from 'svelte-sonner';
	import Actions from './cell-actions.svelte';

	export const cells = {
		row_picker,
		name,
		poolName,
		usage,
		path,
		mode,
		createTime,
		exportSubvolume,
		snapshots,
		actions
	};
</script>

{#snippet row_picker(row: Row<Subvolume>)}
	<Cells.RowPicker {row} />
{/snippet}

{#snippet name(row: Row<Subvolume>)}
	{row.original.name}
{/snippet}

{#snippet path(row: Row<Subvolume>)}
	<p class="max-w-[200px] overflow-auto text-xs font-light">{row.original.path}</p>
{/snippet}

{#snippet mode(row: Row<Subvolume>)}
	<Badge variant="outline">
		{row.original.mode}
	</Badge>
{/snippet}

{#snippet poolName(row: Row<Subvolume>)}
	<Badge variant="outline">
		{row.original.poolName}
	</Badge>
{/snippet}

{#snippet exportSubvolume(row: Row<Subvolume>)}
	{#if row.original.export}
		<div class="flex items-center gap-1">
			<Badge variant="outline">
				{row.original.export.ip}
			</Badge>
			<HoverCard.Root>
				<HoverCard.Trigger>
					<Icon icon="ph:info" />
				</HoverCard.Trigger>
				<HoverCard.Content class="min-w-[500px]">
					<Table.Root>
						<Table.Body class="text-xs">
							<Table.Row>
								<Table.Head class="text-right">IP</Table.Head>
								<Table.Cell>
									{row.original.export.ip}
								</Table.Cell>
							</Table.Row>
							<Table.Row>
								<Table.Head>
									<div class="flex h-full items-center justify-end gap-1">
										<Icon
											icon="ph:copy"
											class="cursor-pointer"
											onclick={async () => {
												const text = row.original.export?.path ?? '';
												navigator.clipboard.writeText(text).then((response) => {
													toast.success('Path copied to clipboard');
												});
											}}
										/>
										Path
									</div>
								</Table.Head>
								<Table.Cell>
									{row.original.export.path}
								</Table.Cell>
							</Table.Row>
							<Table.Row>
								<Table.Head class="text-right">Clients</Table.Head>
								<Table.Cell>
									{#each row.original.export.clients as client}
										<Badge variant="outline">
											{client}
										</Badge>
									{/each}
								</Table.Cell>
							</Table.Row>
							<Table.Row>
								<Table.Head>
									<div class="flex h-full items-center justify-end gap-1">
										<Icon
											icon="ph:copy"
											class="cursor-pointer"
											onclick={async () => {
												const text = row.original.export?.command ?? '';
												navigator.clipboard.writeText(text).then((response) => {
													toast.success('Command copied to clipboard');
												});
											}}
										/>
										Command
									</div>
								</Table.Head>
								<Table.Cell>
									{row.original.export.command}
								</Table.Cell>
							</Table.Row>
						</Table.Body>
					</Table.Root>
				</HoverCard.Content>
			</HoverCard.Root>
		</div>
	{/if}
{/snippet}

{#snippet usage(row: Row<Subvolume>)}
	<div class="flex justify-end">
		<Progress.Root
			numerator={Number(row.original.usedBytes)}
			denominator={Number(row.original.usedBytes)}
		>
			{#snippet ratio({ numerator, denominator })}
				{Progress.formatRatio(numerator, denominator)}
			{/snippet}
			{#snippet detail({ numerator, denominator })}
				{@const { value: numeratorValue, unit: numeratorUnit } = formatCapacity(numerator)}
				{@const { value: denominatorValue, unit: denominatorUnit } = formatCapacity(denominator)}
				{numeratorValue}
				{numeratorUnit}
				/
				{denominatorValue}
				{denominatorUnit}
			{/snippet}
		</Progress.Root>
	</div>
{/snippet}

{#snippet createTime(row: Row<Subvolume>)}
	{#if row.original.createdAt}
		<Tooltip.Provider>
			<Tooltip.Root>
				<Tooltip.Trigger>
					{formatTimeAgo(timestampDate(row.original.createdAt))}
				</Tooltip.Trigger>
				<Tooltip.Content>
					{timestampDate(row.original.createdAt)}
				</Tooltip.Content>
			</Tooltip.Root>
		</Tooltip.Provider>
	{/if}
{/snippet}

{#snippet snapshots(row: Row<Subvolume>)}
	<Snapshot subvolume={row.original} />
{/snippet}

{#snippet actions(row: Row<Subvolume>)}
	<Actions subvolume={row.original} />
{/snippet}
