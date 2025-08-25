<script lang="ts" module>
	import { PoolType, type Pool } from '$lib/api/storage/v1/storage_pb';
	import { Cell as RowPicker } from '$lib/components/custom/data-table/data-table-row-pickers';
	import * as Progress from '$lib/components/custom/progress';
	import { Badge } from '$lib/components/ui/badge';
	import { formatCapacity } from '$lib/formatter';
	import Icon from '@iconify/svelte';
	import type { Row } from '@tanstack/table-core';
	import Actions from './cell-actions.svelte';
	import { getPlacementGroupStateVariant } from './utils.svelte';

	export const cells = {
		row_picker,
		name,
		type,
		applications,
		placement_group_state,
		usage,
		iops,
		actions
	};
</script>

{#snippet row_picker(row: Row<Pool>)}
	<RowPicker {row} />
{/snippet}

{#snippet name(row: Row<Pool>)}
	<div class="flex items-center gap-1">
		{row.original.name}
		{#if row.original.updating}
			<Icon icon="ph:spinner-gap" class="size-5 animate-spin" />
		{/if}
	</div>
{/snippet}

{#snippet type(row: Row<Pool>)}
	<Badge variant="outline">
		{#if row.original.poolType == PoolType.ERASURE}
			ERASURE:{row.original.dataChunks}<Icon icon="ph:x" />{row.original.codingChunks}
		{:else if row.original.poolType == PoolType.REPLICATED}
			REPLICATED:<Icon icon="ph:x" />{row.original.replicatedSize}
		{:else}{/if}
	</Badge>
{/snippet}

{#snippet applications(row: Row<Pool>)}
	<span class="flex gap-1">
		{#each row.original.applications as application}
			{#if application}
				<Badge variant="outline">
					{application}
				</Badge>
			{/if}
		{/each}
	</span>
{/snippet}

{#snippet placement_group_state(row: Row<Pool>)}
	<span class="flex flex-col gap-1">
		{#each Object.entries(row.original.placementGroupState) as [state, number]}
			<Badge variant={getPlacementGroupStateVariant(state)}>
				{state}:{number}
			</Badge>
		{/each}
	</span>
{/snippet}

{#snippet usage(row: Row<Pool>)}
	<div class="flex justify-end">
		<Progress.Root
			numerator={Number(row.original.usedBytes)}
			denominator={Number(row.original.quotaBytes)}
			highIsGood={false}
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
	</div>
{/snippet}

{#snippet iops(row: Row<Pool>)}{/snippet}

{#snippet actions(row: Row<Pool>)}
	<Actions pool={row.original} />
{/snippet}
