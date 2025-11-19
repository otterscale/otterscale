<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import type { Row } from '@tanstack/table-core';

	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import type { ObjectStorageDaemon } from '$lib/api/storage/v1/storage_pb';
	import { Cells } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import * as Progress from '$lib/components/custom/progress';
	import { Badge } from '$lib/components/ui/badge';
	import { formatCapacity } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';

	import Actions from './cell-actions.svelte';

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

{#snippet row_picker(row: Row<ObjectStorageDaemon>)}
	<Layout.Cell class="items-center">
		<Cells.RowPicker {row} />
	</Layout.Cell>
{/snippet}

{#snippet name(row: Row<ObjectStorageDaemon>)}
	<Layout.Cell class="items-start">
		{row.original.name}
	</Layout.Cell>
{/snippet}

{#snippet state(row: Row<ObjectStorageDaemon>)}
	<Layout.Cell class="flex-row items-center">
		{#if row.original.in}
			<Badge variant="outline">{m.osd_in()}</Badge>
		{/if}
		{#if row.original.up}
			<Badge variant="outline">{m.osd_up()}</Badge>
		{/if}
	</Layout.Cell>
{/snippet}

{#snippet osdUp()}{/snippet}

{#snippet osdIn()}{/snippet}

{#snippet exists(row: Row<ObjectStorageDaemon>)}
	<Layout.Cell class="items-start">
		{#if !row.original.exists}
			<Icon icon="ph:x" class="text-destructive" />
		{:else}
			<Icon icon="ph:circle" class="text-primary" />
		{/if}
	</Layout.Cell>
{/snippet}

{#snippet machine(row: Row<ObjectStorageDaemon>)}
	<Layout.Cell class="items-start">
		<div class="flex items-center gap-1">
			<Badge variant="outline">
				{row.original.machine?.hostname}
			</Badge>
			<Icon
				icon="ph:arrow-square-out"
				class="hover:cursor-pointer"
				onclick={() => {
					goto(
						resolve('/(auth)/machines/metal/[id]', {
							id: row.original.machine?.id ?? ''
						})
					);
				}}
			/>
		</div>
	</Layout.Cell>
{/snippet}

{#snippet deviceClass(row: Row<ObjectStorageDaemon>)}
	<Layout.Cell class="items-start">
		<Badge variant="outline">
			{row.original.deviceClass}
		</Badge>
	</Layout.Cell>
{/snippet}

{#snippet placementGroupCount(row: Row<ObjectStorageDaemon>)}
	<Layout.Cell class="items-end">
		{row.original.placementGroupCount}
	</Layout.Cell>
{/snippet}

{#snippet usage(row: Row<ObjectStorageDaemon>)}
	<Layout.Cell class="items-end">
		<Progress.Root
			numerator={Number(row.original.usedBytes)}
			denominator={Number(row.original.sizeBytes)}
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

{#snippet actions(data: { row: Row<ObjectStorageDaemon>; scope: string })}
	<Layout.Cell class="items-start">
		<Actions osd={data.row.original} scope={data.scope} />
	</Layout.Cell>
{/snippet}
