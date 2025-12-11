<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import type { Row } from '@tanstack/table-core';

	import { type Pool, Pool_Application, Pool_Type } from '$lib/api/storage/v1/storage_pb';
	import { Cells } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import * as Progress from '$lib/components/custom/progress';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import { Badge } from '$lib/components/ui/badge';
	import { formatCapacity } from '$lib/formatter';

	import Actions from './cell-actions.svelte';
	import { getPlacementGroupStateClassName, getPlacementGroupStateVariant } from './utils.svelte';

	export const cells = {
		row_picker,
		name,
		type,
		applications,
		placement_group_state,
		usage,
		actions
	};
	export function getPoolApplicationLabel(value: number): string | undefined {
		switch (value) {
			case Pool_Application.BLOCK:
				return 'BLOCK';
			case Pool_Application.FILE:
				return 'FILE';
			case Pool_Application.OBJECT:
				return 'OBJECT';
		}
	}
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
			{#if row.original.type == Pool_Type.ERASURE}
				ERASURE:{row.original.dataChunks}<Icon icon="ph:x" />{row.original.codingChunks}
			{:else if row.original.type == Pool_Type.REPLICATED}
				REPLICATED:<Icon icon="ph:x" />{row.original.replicatedSize}
			{:else}{/if}
		</Badge>
	</Layout.Cell>
{/snippet}

{#snippet applications(row: Row<Pool>)}
	<Layout.Cell class="items-start">
		<span class="flex gap-1">
			{#each row.original.applications as application}
				{@const label = getPoolApplicationLabel(application)}
				{#if label}
					<Badge variant="outline">
						{label}
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
				<Badge
					variant={getPlacementGroupStateVariant(state)}
					class={getPlacementGroupStateClassName(state)}
				>
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
	</Layout.Cell>
{/snippet}

{#snippet actions(data: { row: Row<Pool>; scope: string; reloadManager: ReloadManager })}
	<Layout.Cell class="items-start">
		<Actions pool={data.row.original} scope={data.scope} reloadManager={data.reloadManager} />
	</Layout.Cell>
{/snippet}
