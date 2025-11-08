<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import type { Row } from '@tanstack/table-core';

	import Actions from './cell-actions.svelte';
	import { getPlacementGroupStateVariant } from './utils.svelte';

	import { PoolType, type Pool } from '$lib/api/storage/v1/storage_pb';
	import { Cells } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import * as Progress from '$lib/components/custom/progress';
	import { Badge } from '$lib/components/ui/badge';
	import { formatCapacity } from '$lib/formatter';

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
	<Layout.Cell class="items-center">
		<Cells.RowPicker {row} />
	</Layout.Cell>
{/snippet}

{#snippet name(row: Row<Pool>)}
	<Layout.Cell class="items-start">
		<div class="flex items-center gap-1">
			{row.original.name}
			{#if row.original.updating}
				<Icon icon="ph:spinner-gap" class="size-5 animate-spin" />
			{/if}
		</div>
	</Layout.Cell>
{/snippet}

{#snippet type(row: Row<Pool>)}
	<Layout.Cell class="items-start">
		<Badge variant="outline">
			{#if row.original.poolType == PoolType.ERASURE}
				ERASURE:{row.original.dataChunks}<Icon icon="ph:x" />{row.original.codingChunks}
			{:else if row.original.poolType == PoolType.REPLICATED}
				REPLICATED:<Icon icon="ph:x" />{row.original.replicatedSize}
			{:else}{/if}
		</Badge>
	</Layout.Cell>
{/snippet}

{#snippet applications(row: Row<Pool>)}
	<Layout.Cell class="items-start">
		<span class="flex gap-1">
			{#each row.original.applications as application}
				{#if application}
					<Badge variant="outline">
						{application}
					</Badge>
				{/if}
			{/each}
		</span>
	</Layout.Cell>
{/snippet}

{#snippet placement_group_state(row: Row<Pool>)}
	<Layout.Cell class="items-start">
		<span class="flex flex-col gap-1">
			{#each Object.entries(row.original.placementGroupState) as [state, number]}
				<Badge variant={getPlacementGroupStateVariant(state)}>
					{state}:{number}
				</Badge>
			{/each}
		</span>
	</Layout.Cell>
{/snippet}

{#snippet usage(row: Row<Pool>)}
	<Layout.Cell class="items-end">
		<Progress.Root
			numerator={Number(row.original.usedBytes)}
			denominator={Number(row.original.usedBytes + row.original.maxBytes)}
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
	</Layout.Cell>
{/snippet}

{#snippet iops()}{/snippet}

{#snippet actions(row: Row<Pool>)}
	<Layout.Cell class="items-start">
		<Actions pool={row.original} />
	</Layout.Cell>
{/snippet}
