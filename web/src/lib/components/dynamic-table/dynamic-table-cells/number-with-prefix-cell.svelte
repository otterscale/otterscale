<script lang="ts" module>
	import type { JsonValue } from '@bufbuild/protobuf';

	export type PrefixType = 'decimal' | 'binary';

	export type NumberWithPrefixMetadata = {
		prefix: PrefixType;
	};
</script>

<script lang="ts">
	import { type Column, type Row } from '@tanstack/table-core';

	import { formatWithBinaryPrefix, formatWithDecimalPrefix } from '../utils';

	let {
		row,
		column,
		metadata
	}: {
		row: Row<Record<string, JsonValue>>;
		column: Column<Record<string, JsonValue>>;
		metadata: NumberWithPrefixMetadata;
	} = $props();

	if (!metadata) {
		throw Error(
			`Expected metadata of ${column.id} for NumberWithPrefixCell, but got metadata:`,
			metadata
		);
	}

	const data = $derived(
		row.original[column.id] ? BigInt(row.original[column.id] as number) : undefined
	);
</script>

{#if data}
	{#if metadata.prefix === 'decimal'}
		{@const { value, unit } = formatWithDecimalPrefix(data)}
		{`${value.toFixed(0)} ${unit}`}
	{:else if metadata.prefix === 'binary'}
		{@const { value, unit } = formatWithBinaryPrefix(data)}
		{`${value.toFixed(0)} ${unit}`}
	{/if}
{/if}
