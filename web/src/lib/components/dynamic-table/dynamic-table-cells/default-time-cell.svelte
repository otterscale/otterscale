<script lang="ts">
	import type { JsonValue } from '@bufbuild/protobuf';
	import { type Column, type Row } from '@tanstack/table-core';

	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { now } from '$lib/stores/now';

	import { getRelativeTime } from '../utils';

	let {
		row,
		column
	}: {
		row: Row<Record<string, JsonValue>>;
		column: Column<Record<string, JsonValue>>;
	} = $props();

	const data = $derived(new Date(row.original[column.id] as string));
</script>

{#if data && !isNaN(data.getTime())}
	{@const { value, unit } = getRelativeTime($now, data.getTime())}
	<Tooltip.Provider>
		<Tooltip.Root>
			<Tooltip.Trigger>
				{value}
				{unit}
			</Tooltip.Trigger>
			<Tooltip.Content>
				{new Intl.DateTimeFormat('en-CA', {
					year: 'numeric',
					month: '2-digit',
					day: '2-digit',
					hour: '2-digit',
					minute: '2-digit',
					second: '2-digit',
					hour12: false,
					timeZoneName: 'longOffset'
				}).format(data)}
			</Tooltip.Content>
		</Tooltip.Root>
	</Tooltip.Provider>
{/if}
