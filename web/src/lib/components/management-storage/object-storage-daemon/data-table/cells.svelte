<script lang="ts" module>
	import { goto } from '$app/navigation';
	import type { OSD } from '$gen/api/storage/v1/storage_pb';
	import TableRowPicker from '$lib/components/custom/data-table/data-table-row-pickers/cell.svelte';
	import * as Progress from '$lib/components/custom/progress';
	import { Badge } from '$lib/components/ui/badge';
	import { formatCapacity } from '$lib/formatter';
	import Icon from '@iconify/svelte';
	import type { Row } from '@tanstack/table-core';
	import { LineChart } from 'layerchart';

	const renderContext: 'svg' | 'canvas' = 'canvas';
	const debug = false;
	function generateRandomTimeSeriesData(days = 7, minValue = 50, maxValue = 100) {
		const data = [];
		const now = new Date();

		for (let i = 0; i < days; i++) {
			data.push({
				date: new Date(now.getTime() - (days - 1 - i) * 24 * 60 * 60 * 1000),
				value: Math.floor(Math.random() * (maxValue - minValue) + minValue)
			});
		}

		return data;
	}

	export const cells = {
		_row_picker: _row_picker,
		id: id,
		name: name,
		stateUp: stateUp,
		stateIn: stateIn,
		exists: exists,
		deviceClass: deviceClass,
		machine: machine,
		placementGroupCount: placementGroupCount,
		usage: usage,
		readBytes: readBytes,
		writeBytes: writeBytes
	};
</script>

{#snippet _row_picker(row: Row<OSD>)}
	<TableRowPicker {row} />
{/snippet}

{#snippet id(row: Row<OSD>)}
	{row.original.id}
{/snippet}

{#snippet name(row: Row<OSD>)}
	{row.original.name}
{/snippet}

{#snippet stateUp(row: Row<OSD>)}
	<Badge variant="outline">
		{row.original.up}
	</Badge>
{/snippet}

{#snippet stateIn(row: Row<OSD>)}
	<Badge variant="outline">
		{row.original.in}
	</Badge>
{/snippet}

{#snippet exists(row: Row<OSD>)}
	<Badge variant="outline">
		{row.original.exists}
	</Badge>
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
	<Progress.Root
		numerator={Number(row.original.usedBytes)}
		denominator={Number(row.original.sizeBytes)}
	>
		{#snippet ratio({ numerator, denominator })}
			{((numerator * 100) / denominator).toFixed(2)}%
		{/snippet}
	</Progress.Root>
{/snippet}

{#snippet readBytes(row: Row<OSD>)}
	<span>
		<div class="h-[50px] w-[100px]">
			<LineChart data={generateRandomTimeSeriesData()} x="date" y="value" {renderContext} {debug} />
		</div>
	</span>
{/snippet}

{#snippet writeBytes(row: Row<OSD>)}
	<span>
		<div class="h-[50px] w-[100px]">
			<LineChart data={generateRandomTimeSeriesData()} x="date" y="value" {renderContext} {debug} />
		</div>
	</span>
{/snippet}
