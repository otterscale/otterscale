<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import type { Row } from '@tanstack/table-core';

	import { type Pool, PoolType } from '$lib/api/storage/v1/storage_pb';
	import { Cells } from '$lib/components/custom/data-table/core';
	import * as Progress from '$lib/components/custom/progress';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Table from '$lib/components/custom/table/index.js';
	import { Badge } from '$lib/components/ui/badge';
	import { formatCapacity } from '$lib/formatter';

	import Actions from './cell-actions.svelte';
	import { getPlacementGroupStateVariant } from './utils.svelte';

	export const cells = {
		row_picker,
		name,
		type,
		applications,
		placement_group_state,
		usage,
		actions
	};
</script>

{#snippet row_picker(row: Row<Pool>)}
	<Table.Cell alignClass="items-center">
		<Cells.RowPicker {row} />
	</Table.Cell>
{/snippet}

{#snippet name(row: Row<Pool>)}
	<Table.Cell alignClass="items-start">
		<div class="flex items-center gap-1">
			{row.original.name}
			{#if row.original.updating}
				<Icon icon="ph:spinner-gap" class="size-5 animate-spin" />
			{/if}
		</div>
	</Table.Cell>
{/snippet}

{#snippet type(row: Row<Pool>)}
	<Table.Cell alignClass="items-start">
		<Badge variant="outline">
			{#if row.original.poolType == PoolType.ERASURE}
				ERASURE:{row.original.dataChunks}<Icon icon="ph:x" />{row.original.codingChunks}
			{:else if row.original.poolType == PoolType.REPLICATED}
				REPLICATED:<Icon icon="ph:x" />{row.original.replicatedSize}
			{:else}{/if}
		</Badge>
	</Table.Cell>
{/snippet}

{#snippet applications(row: Row<Pool>)}
	<Table.Cell alignClass="items-start">
		<span class="flex gap-1">
			{#each row.original.applications as application}
				{#if application}
					<Badge variant="outline">
						{application}
					</Badge>
				{/if}
			{/each}
		</span>
	</Table.Cell>
{/snippet}

{#snippet placement_group_state(row: Row<Pool>)}
	<Table.Cell alignClass="items-start">
		<span class="flex flex-col gap-1">
			{#each Object.entries(row.original.placementGroupState) as [state, number]}
				<Badge variant={getPlacementGroupStateVariant(state)}>
					{state}:{number}
				</Badge>
			{/each}
		</span>
	</Table.Cell>
{/snippet}

{#snippet usage(row: Row<Pool>)}
	<Table.Cell alignClass="items-end">
		<Progress.Root
			numerator={Number(row.original.usedBytes)}
			denominator={Number(row.original.usedBytes + row.original.maxBytes)}
			target="STB"
		>
			{#snippet ratio({ numerator, denominator })}
				{Progress.formatRatio(numerator, denominator)}
			{/snippet}
			{#snippet detail({ numerator, denominator })}
				{@const { value: numeratorValue, unit: numeratorUnit } = formatCapacity(numerator)}
				{@const { value: denominatorValue, unit: denominatorUnit } = formatCapacity(denominator)}
				{numeratorValue}
				{numeratorUnit}/{denominatorValue}
				{denominatorUnit}
			{/snippet}
		</Progress.Root>
	</Table.Cell>
{/snippet}

{#snippet actions(data: { row: Row<Pool>; scope: string; reloadManager: ReloadManager })}
	<Table.Cell alignClass="items-start">
		<Actions pool={data.row.original} scope={data.scope} reloadManager={data.reloadManager} />
	</Table.Cell>
{/snippet}
