<script lang="ts" module>
	import type { OSD } from '$gen/api/storage/v1/storage_pb';
	import * as Progress from '$lib/components/custom/progress';
	import { Badge } from '$lib/components/ui/badge';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import { formatCapacity } from '$lib/formatter';
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
		placementGroupCount: placementGroupCount,
		usage: usage,
		readBytes: readBytes,
		writeBytes: writeBytes
	};
</script>

{#snippet _row_picker(row: Row<OSD>)}
	<Checkbox
		checked={row.getIsSelected()}
		onCheckedChange={(value) => row.toggleSelected(!!value)}
		class="border-secondary-950"
		aria-label="Select row"
	/>
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
		{#snippet detail({ numerator, denominator })}
			{@const { value: numeratorValue, unit: numeratorUnit } = formatCapacity(
				numerator / (1024 * 1024)
			)}
			{@const { value: denominatorValue, unit: denominatorUnit } = formatCapacity(
				denominator / (1024 * 1024)
			)}
			<span>
				{numeratorValue}
				{numeratorUnit}
			</span>
			<span>/</span>
			<span>
				{denominatorValue}
				{denominatorUnit}
			</span>
		{/snippet}
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
