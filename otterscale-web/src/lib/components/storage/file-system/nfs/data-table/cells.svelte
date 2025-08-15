<script lang="ts" module>
	import type { Subvolume } from '$lib/api/storage/v1/storage_pb';
	import TableRowPicker from '$lib/components/custom/data-table/data-table-row-pickers/cell.svelte';
	import * as Progress from '$lib/components/custom/progress/index.js';
	import { Badge } from '$lib/components/ui/badge';
	import * as HoverCard from '$lib/components/ui/hover-card';
	import * as Table from '$lib/components/ui/table';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { formatCapacityV2 as formatCapacity, formatTimeAgo } from '$lib/formatter';
	import { timestampDate } from '@bufbuild/protobuf/wkt';
	import Icon from '@iconify/svelte';
	import type { Row } from '@tanstack/table-core';
	import { toast } from 'svelte-sonner';

	export const cells = {
		_row_picker,
		name,
		poolName,
		usage,
		path,
		mode,
		createTime,
		Export
	};
</script>

{#snippet _row_picker(row: Row<Subvolume>)}
	<TableRowPicker {row} />
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

{#snippet Export(row: Row<Subvolume>)}
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
