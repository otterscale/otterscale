<script lang="ts" module>
	import type { JsonValue } from '@bufbuild/protobuf';

	import { Progress } from '$lib/components/ui/progress/index.js';

	export type RatioMetadata = {
		numerator: JsonValue;
		denominator: JsonValue;
	};
</script>

<script lang="ts">
	import { type Column, type Row } from '@tanstack/table-core';

	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	let {
		row,
		column,
		metadata
	}: {
		row: Row<Record<string, JsonValue>>;
		column: Column<Record<string, JsonValue>>;
		metadata: RatioMetadata;
	} = $props();

	const data = $derived(row.original[column.id] as number);
</script>

{#if metadata.denominator !== undefined}
	<div class="flex flex-col gap-1">
		<Progress value={data} max={1} class="w-full" />
		<Tooltip.Provider>
			<Tooltip.Root>
				<Tooltip.Trigger class="ml-auto">
					{(Number(data) * 100).toFixed(2)}%
				</Tooltip.Trigger>
				<Tooltip.Content>
					{metadata.numerator}/{metadata.denominator}
				</Tooltip.Content>
			</Tooltip.Root>
		</Tooltip.Provider>
	</div>
{/if}
