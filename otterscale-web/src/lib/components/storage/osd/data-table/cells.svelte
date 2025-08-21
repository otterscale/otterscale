<script lang="ts" module>
	import { goto } from '$app/navigation';
	import type { OSD } from '$lib/api/storage/v1/storage_pb';
	import { Cell as RowPicker } from '$lib/components/custom/data-table/data-table-row-pickers';
	import * as Progress from '$lib/components/custom/progress';
	import { Badge } from '$lib/components/ui/badge';
	import { formatCapacity } from '$lib/formatter';
	import Icon from '@iconify/svelte';
	import type { Row } from '@tanstack/table-core';
	import Actions from './cell-actions.svelte';
	import { m } from '$lib/paraglide/messages';

	export const cells = {
		row_picker,
		name,
		state,
		osdUp,
		osdIn,
		exists,
		deviceClass,
		machine,
		placementGroupCount,
		usage,
		iops,
		actions
	};
</script>

{#snippet row_picker(row: Row<OSD>)}
	<RowPicker {row} />
{/snippet}

{#snippet name(row: Row<OSD>)}
	{row.original.name}
{/snippet}

{#snippet state(row: Row<OSD>)}
	<div class="flex items-center gap-1">
		{#if row.original.in}
			<Badge variant="outline">{m.storage_osds_in()}</Badge>
		{/if}
		{#if row.original.up}
			<Badge variant="outline">{m.storage_osds_up()}</Badge>
		{/if}
	</div>
{/snippet}

{#snippet osdUp(row: Row<OSD>)}{/snippet}

{#snippet osdIn(row: Row<OSD>)}{/snippet}

{#snippet exists(row: Row<OSD>)}
	{#if !row.original.exists}
		<Icon icon="ph:x" class="text-destructive" />
	{:else}
		<Icon icon="ph:circle" class="text-primary" />
	{/if}
{/snippet}

{#snippet machine(row: Row<OSD>)}
	<div class="flex items-center gap-1">
		<Badge variant="outline">
			{row.original.machine?.hostname}
		</Badge>
		<Icon
			icon="ph:arrow-square-out"
			class="hover:cursor-pointer"
			onclick={() => {
				goto(`/management/machine/${row.original.machine?.id}`);
			}}
		/>
	</div>
{/snippet}

{#snippet deviceClass(row: Row<OSD>)}
	<Badge variant="outline">
		{row.original.deviceClass}
	</Badge>
{/snippet}

{#snippet placementGroupCount(row: Row<OSD>)}
	<span class="flex justify-end">{row.original.placementGroupCount}</span>
{/snippet}

{#snippet usage(row: Row<OSD>)}
	<div class="flex justify-end">
		<Progress.Root
			numerator={Number(row.original.usedBytes)}
			denominator={Number(row.original.sizeBytes)}
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

{#snippet iops(row: Row<OSD>)}{/snippet}

{#snippet actions(row: Row<OSD>)}
	<Actions osd={row.original} />
{/snippet}
