<script lang="ts">
	import type { JsonValue } from '@bufbuild/protobuf';
	import type { Column, Row } from '@tanstack/table-core';
	import { type WithElementRef } from 'bits-ui';
	import type { HTMLAttributes } from 'svelte/elements';

	import ArrayOfObjectCell, {
		type ArrayOfObjectMetadata
	} from './dynamic-table-cells/array-of-object-cell.svelte';
	import DefaultArrayCell from './dynamic-table-cells/default-array-cell.svelte';
	import DefaultBooleanCell from './dynamic-table-cells/default-boolean-cell.svelte';
	import DefaultEmptyCell from './dynamic-table-cells/default-empty-cell.svelte';
	import DefaultNumberCell from './dynamic-table-cells/default-number-cell.svelte';
	import DefaultObjectCell from './dynamic-table-cells/default-object-cell.svelte';
	import DefaultTextCell from './dynamic-table-cells/default-text-cell.svelte';
	import DefaultTimeCell from './dynamic-table-cells/default-time-cell.svelte';
	import LinkCell, { type LinkMetadata } from './dynamic-table-cells/link-cell.svelte';
	import ObjectOfKeyValue, {
		type ObjectOfKeyValueMetadata
	} from './dynamic-table-cells/object-of-key-value.svelte';
	import RatioCell, { type RatioMetadata } from './dynamic-table-cells/ratio-cell.svelte';
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
		metadata?: ArrayOfObjectMetadata | LinkMetadata | RatioMetadata | ObjectOfKeyValueMetadata;
	} = $props();

	const uiSchema = $derived(uiSchemas[column.id]);
</script>

{#if uiSchema === 'array-of-object'}
	<ArrayOfObjectCell {row} {column} metadata={metadata as ArrayOfObjectMetadata} />
{:else if uiSchema === 'array'}
	<DefaultArrayCell {row} {column} />
{:else if uiSchema === 'boolean'}
	<DefaultBooleanCell {row} {column} />
{:else if uiSchema === 'number'}
	<DefaultNumberCell {row} {column} />
{:else if uiSchema === 'object'}
	<DefaultObjectCell {row} {column} />
{:else if uiSchema === 'text'}
	<DefaultTextCell {row} {column} />
{:else if uiSchema === 'time'}
	<DefaultTimeCell {row} {column} />
{:else if uiSchema === 'link'}
	<LinkCell {row} {column} metadata={metadata as LinkMetadata} />
{:else if uiSchema === 'object-of-key-value'}
	{console.log(metadata)}
	<ObjectOfKeyValue {row} {column} metadata={metadata as ObjectOfKeyValueMetadata} />
{:else if uiSchema === 'ratio'}
	<RatioCell {row} {column} metadata={metadata as RatioMetadata} />
{:else}
	<DefaultEmptyCell />
{/if}
