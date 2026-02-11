import { type JsonValue } from '@bufbuild/protobuf';
import type { Column, ColumnDef } from '@tanstack/table-core';
import { type Row } from '@tanstack/table-core';

import { resolve } from '$app/paths';
import { page } from '$app/state';
import type { APIResource } from '$lib/api/resource/v1/resource_pb';
import { DynamicTableHeader } from '$lib/components/dynamic-table';
import DynamicTableCell from '$lib/components/dynamic-table/dynamic-table-cell.svelte';
import type { LinkMetadata } from '$lib/components/dynamic-table/dynamic-table-cells/link-cell.svelte';
import type { ObjectOfKeyValueMetadata } from '$lib/components/dynamic-table/dynamic-table-cells/object-of-key-value.svelte';
import { type DataSchemaType, type UISchemaType } from '$lib/components/dynamic-table/utils';
import { renderComponent } from '$lib/components/ui/data-table';

type DefaultAttribute =
	| 'Name'
	| 'Namespace'
	| 'Labels'
	| 'Annotations'
	| 'Creation Timestamp'
	| 'raw';

// Determine metadata type
function getDefaultUISchemas(): Record<DefaultAttribute, UISchemaType> {
	return {
		Name: 'link' as UISchemaType,
		Namespace: 'text' as UISchemaType,
		Labels: 'object-of-key-value' as UISchemaType,
		Annotations: 'object-of-key-value' as UISchemaType,
		'Creation Timestamp': 'time' as UISchemaType,
		raw: 'object' as UISchemaType
	};
}

// Determine data type
function getDefaultDataSchemas(): Record<DefaultAttribute, DataSchemaType> {
	return {
		Name: 'text' as DataSchemaType,
		Namespace: 'text' as DataSchemaType,
		Labels: 'number' as DataSchemaType,
		Annotations: 'number' as DataSchemaType,
		'Creation Timestamp': 'time' as DataSchemaType,
		raw: 'object' as DataSchemaType
	};
}

function getDefaultData(object: any): Record<DefaultAttribute, JsonValue> {
	return {
		Name: object?.metadata?.name as JsonValue,
		Namespace: object?.metadata?.namespace as JsonValue,
		Labels: Object.keys(object?.metadata?.labels).length as JsonValue,
		Annotations: Object.keys(object?.metadata?.annotations).length as JsonValue,
		'Creation Timestamp': object?.metadata?.creationTimestamp as JsonValue,
		raw: object as JsonValue
	};
}

function getDefaultColumnDefinitions(
	apiResource: APIResource,
	uiSchemas: Record<DefaultAttribute, UISchemaType>,
	dataSchemas: Record<DefaultAttribute, DataSchemaType>
): ColumnDef<Record<DefaultAttribute, JsonValue>>[] {
	return [
		{
			id: 'Name',
			header: ({ column }: { column: Column<Record<DefaultAttribute, JsonValue>> }) =>
				renderComponent(DynamicTableHeader, {
					column: column,
					dataSchemas: dataSchemas
				}),
			cell: ({
				column,
				row
			}: {
				column: Column<Record<DefaultAttribute, JsonValue>>;
				row: Row<Record<DefaultAttribute, JsonValue>>;
			}) =>
				renderComponent(DynamicTableCell, {
					row: row,
					column: column,
					uiSchemas: uiSchemas,
					metadata: {
						hyperlink: resolve(
							`/(auth)/${page.params.cluster!}/${apiResource.kind}/${apiResource.resource}?group=${apiResource.group}&version=${apiResource.version}&name=${row.original[column.id as DefaultAttribute]}&namespace=${page.url.searchParams.get('namespace') ?? ''}`
						)
					} as LinkMetadata
				}),
			accessorKey: 'Name'
		},
		...[
			{
				id: 'Namespace',
				header: ({ column }: { column: Column<Record<DefaultAttribute, JsonValue>> }) =>
					renderComponent(DynamicTableHeader, {
						column: column,
						dataSchemas: dataSchemas
					}),
				cell: ({
					column,
					row
				}: {
					column: Column<Record<DefaultAttribute, JsonValue>>;
					row: Row<Record<DefaultAttribute, JsonValue>>;
				}) =>
					renderComponent(DynamicTableCell, {
						row: row,
						column: column,
						uiSchemas: uiSchemas
					}),
				accessorKey: 'Namespace'
			}
		].filter(() => apiResource!.namespaced),
		{
			id: 'Annotations',
			header: ({ column }: { column: Column<Record<DefaultAttribute, JsonValue>> }) =>
				renderComponent(DynamicTableHeader, {
					column: column,
					dataSchemas: dataSchemas
				}),
			cell: ({
				column,
				row
			}: {
				column: Column<Record<DefaultAttribute, JsonValue>>;
				row: Row<Record<DefaultAttribute, JsonValue>>;
			}) =>
				renderComponent(DynamicTableCell, {
					row: row,
					column: column,
					uiSchemas: uiSchemas,
					metadata: (row.original.raw as any).metadata.annotations as ObjectOfKeyValueMetadata
				}),
			accessorFn: (row: Record<DefaultAttribute, JsonValue>) =>
				row['Annotations'] ? Object.keys(row['Annotations'] as object).length : null
		},
		{
			id: 'Labels',
			header: ({ column }: { column: Column<Record<DefaultAttribute, JsonValue>> }) =>
				renderComponent(DynamicTableHeader, {
					column: column,
					dataSchemas: dataSchemas
				}),
			cell: ({
				column,
				row
			}: {
				column: Column<Record<DefaultAttribute, JsonValue>>;
				row: Row<Record<DefaultAttribute, JsonValue>>;
			}) =>
				renderComponent(DynamicTableCell, {
					row: row,
					column: column,
					uiSchemas: uiSchemas,
					metadata: (row.original.raw as any).metadata.labels as ObjectOfKeyValueMetadata
				}),
			accessorFn: (row: Record<DefaultAttribute, JsonValue>) =>
				row['Labels'] ? Object.keys(row['Labels'] as object).length : null
		},
		{
			id: 'Creation Timestamp',
			header: ({ column }: { column: Column<Record<DefaultAttribute, JsonValue>> }) =>
				renderComponent(DynamicTableHeader, {
					column: column,
					dataSchemas: dataSchemas
				}),
			cell: ({
				column,
				row
			}: {
				column: Column<Record<DefaultAttribute, JsonValue>>;
				row: Row<Record<DefaultAttribute, JsonValue>>;
			}) =>
				renderComponent(DynamicTableCell, {
					row: row,
					column: column,
					uiSchemas: uiSchemas
				}),
			accessorKey: 'Creation Timestamp'
		}
	];
}

export { getDefaultColumnDefinitions, getDefaultData, getDefaultDataSchemas, getDefaultUISchemas };
