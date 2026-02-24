<script lang="ts">
	import type { JsonValue } from '@bufbuild/protobuf';
	import type { Column, Row } from '@tanstack/table-core';
	import { type WithElementRef } from 'bits-ui';
	import type { HTMLAttributes } from 'svelte/elements';

	import ArrayCell from './dynamic-table-cells/array-cell.svelte';
	import ArrayOfEnumerationCell from './dynamic-table-cells/array-of-enumeration-cell.svelte';
	import ArrayOfObjectCell, {
		type ArrayOfObjectMetadata
	} from './dynamic-table-cells/array-of-object-cell.svelte';
	import BooleanCell from './dynamic-table-cells/boolean-cell.svelte';
	import EmptyCell from './dynamic-table-cells/empty-cell.svelte';
	import LinkCell, { type LinkMetadata } from './dynamic-table-cells/link-cell.svelte';
	import NumberCell from './dynamic-table-cells/number-cell.svelte';
	import NumberWithPrefixCell, {
		type NumberWithPrefixMetadata
	} from './dynamic-table-cells/number-with-prefix-cell.svelte';
	import ObjectCell from './dynamic-table-cells/object-cell.svelte';
	import ObjectOfKeyValueCell, {
		type ObjectOfKeyValueMetadata
	} from './dynamic-table-cells/object-of-key-value-cell.svelte';
	import RatioCell, { type RatioMetadata } from './dynamic-table-cells/ratio-cell.svelte';
	import TextCell from './dynamic-table-cells/text-cell.svelte';
	import TimeCell from './dynamic-table-cells/time-cell.svelte';
	import type { UISchemaType } from './utils';

	let {
		ref = $bindable(null),
		uiSchemas,
		row,
		column,
		metadata
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> & {
		uiSchemas: Record<string, UISchemaType>;
		row: Row<Record<string, JsonValue>>;
		column: Column<Record<string, JsonValue>>;
		metadata?:
			| ArrayOfObjectMetadata
			| LinkMetadata
			| RatioMetadata
			| ObjectOfKeyValueMetadata
			| NumberWithPrefixMetadata;
	} = $props();

	const uiSchema = $derived(uiSchemas[column.id]);
</script>

{#if uiSchema === 'array'}
	<ArrayCell {row} {column} />
{:else if uiSchema === 'array-of-enumeration'}
	<ArrayOfEnumerationCell {row} {column} />
{:else if uiSchema === 'array-of-object'}
	<ArrayOfObjectCell {row} {column} metadata={metadata as ArrayOfObjectMetadata} />
{:else if uiSchema === 'boolean'}
	<BooleanCell {row} {column} />
{:else if uiSchema === 'link'}
	<LinkCell {row} {column} metadata={metadata as LinkMetadata} />
{:else if uiSchema === 'number'}
	<NumberCell {row} {column} />
{:else if uiSchema === 'number-with-prefix'}
	<NumberWithPrefixCell {row} {column} metadata={metadata as NumberWithPrefixMetadata} />
{:else if uiSchema === 'object'}
	<ObjectCell {row} {column} />
{:else if uiSchema === 'object-of-key-value'}
	<ObjectOfKeyValueCell {row} {column} metadata={metadata as ObjectOfKeyValueMetadata} />
{:else if uiSchema === 'ratio'}
	<RatioCell {row} {column} metadata={metadata as RatioMetadata} />
{:else if uiSchema === 'text'}
	<TextCell {row} {column} />
{:else if uiSchema === 'time'}
	<TimeCell {row} {column} />
{:else}
	<EmptyCell />
{/if}
