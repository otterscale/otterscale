<script lang="ts" module>
	import * as Progress from '$lib/components/custom/progress';
	import { Badge } from '$lib/components/ui/badge';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import type { Row } from '@tanstack/table-core';
	import type { Pool } from './types';
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
		dataProtection: dataProtection,
		applications: applications,
		PGStatus: PGStatus,
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

{#snippet dataProtection(row: Row<Pool>)}
	<Badge variant="outline">
		{row.original.dataProtection}
	</Badge>
{/snippet}

{#snippet applications(row: Row<Pool>)}
	<Badge variant="outline">
		{row.original.applications}
	</Badge>
{/snippet}

{#snippet PGStatus(row: Row<Pool>)}
	{row.original.PGStatus}
{/snippet}

{#snippet usage(row: Row<Pool>)}
	<Progress.Root numerator={row.original.usage} denominator={100}>
		{#snippet ratio({ numerator, denominator })}
			{(numerator * 100) / denominator}%
		{/snippet}
	</Progress.Root>
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
