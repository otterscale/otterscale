<script lang="ts" module>
	import type { Pool } from '$gen/api/storage/v1/storage_pb';
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
		name: name,
		applications: applications,
		placement_group_state: placement_group_state,
		usage: usage,
		readBytes: readBytes,
		writeBytes: writeBytes
	};
</script>

{#snippet _row_picker(row: Row<Pool>)}
	<Checkbox
		checked={row.getIsSelected()}
		onCheckedChange={(value) => row.toggleSelected(!!value)}
		class="border-secondary-950"
		aria-label="Select row"
	/>
{/snippet}

{#snippet name(row: Row<Pool>)}
	{row.original.name}
{/snippet}

{#snippet applications(row: Row<Pool>)}
	{#each row.original.applications as application}
		{#if application}
			<Badge variant="outline">
				{application}
			</Badge>
		{/if}
	{/each}
{/snippet}

{#snippet placement_group_state(row: Row<Pool>)}
	{#each Object.entries(row.original.placementGroupState) as [state, number]}
		<Badge variant="outline">
			{state}:{number}
		</Badge>
	{/each}
{/snippet}

{#snippet usage(row: Row<Pool>)}
	<div class="flex justify-end">
		<Progress.Root
			numerator={Number(row.original.usedBytes)}
			denominator={Number(row.original.quotaBytes)}
		>
			{#snippet ratio({ numerator, denominator })}
				{(numerator * 100) / denominator}%
			{/snippet}
		</Progress.Root>
	</div>
{/snippet}

{#snippet readBytes(row: Row<Pool>)}
	<div class="h-[50px] w-[100px]">
		<LineChart data={generateRandomTimeSeriesData()} x="date" y="value" {renderContext} {debug} />
	</div>
{/snippet}

{#snippet writeBytes(row: Row<Pool>)}
	<div class="h-[50px] w-[100px]">
		<LineChart data={generateRandomTimeSeriesData()} x="date" y="value" {renderContext} {debug} />
	</div>
{/snippet}
