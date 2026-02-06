import { type JsonValue } from '@bufbuild/protobuf';
import type { Column } from '@tanstack/table-core';
import { type Row } from '@tanstack/table-core';
import lodash from 'lodash';

import type { APIResource } from '$lib/api/resource/v1/resource_pb';
import { DynamicTableCell, DynamicTableHeader } from '$lib/components/dynamic-table';
import { renderComponent } from '$lib/components/ui/data-table';

function defaultFieldsMask(
	schema: any
): Record<string, { description: string; type: string; format?: string }> {
	return {
		Name: lodash.get(schema, 'properties.metadata.properties.name'),
		Namespace: lodash.get(schema, 'properties.metadata.properties.namespace'),
		Labels: lodash.get(schema, 'properties.metadata.properties.labels'),
		Annotations: lodash.get(schema, 'properties.metadata.properties.annotations'),
		CreateTime: lodash.get(schema, 'properties.metadata.properties.creationTimestamp'),
		Configuration: schema
	};
}
function defaultObjectMask(object: any): Record<string, JsonValue> {
	return {
		Name: lodash.get(object, 'metadata.name'),
		Namespace: lodash.get(object, 'metadata.namespace'),
		Labels: lodash.get(object, 'metadata.labels'),
		Annotations: lodash.get(object, 'metadata.annotations'),
		CreateTime: lodash.get(object, 'metadata.creationTimestamp'),
		Configuration: object
	};
}
function defaultColumnDefinitions(apiResource: APIResource, fields: any) {
	return [
		{
			id: 'Name',
			header: ({ column }: { column: Column<Record<string, JsonValue>> }) =>
				renderComponent(DynamicTableHeader, {
					column: column,
					fields: fields
				}),
			cell: ({
				column,
				row
			}: {
				column: Column<Record<string, JsonValue>>;
				row: Row<Record<string, JsonValue>>;
			}) =>
				renderComponent(DynamicTableCell, {
					row: row,
					column: column,
					fields: fields
				}),
			accessorKey: 'Name'
		},
		...[
			{
				id: 'Namespace',
				header: ({ column }: { column: Column<Record<string, JsonValue>> }) =>
					renderComponent(DynamicTableHeader, {
						column: column,
						fields: fields
					}),
				cell: ({
					column,
					row
				}: {
					column: Column<Record<string, JsonValue>>;
					row: Row<Record<string, JsonValue>>;
				}) =>
					renderComponent(DynamicTableCell, {
						row: row,
						column: column,
						fields: fields
					}),
				accessorKey: 'Namespace'
			}
		].filter(() => apiResource.namespaced),
		{
			id: 'Annotations',
			header: ({ column }: { column: Column<Record<string, JsonValue>> }) =>
				renderComponent(DynamicTableHeader, {
					column: column,
					fields: fields
				}),
			cell: ({
				column,
				row
			}: {
				column: Column<Record<string, JsonValue>>;
				row: Row<Record<string, JsonValue>>;
			}) =>
				renderComponent(DynamicTableCell, {
					row: row,
					column: column,
					fields: fields
				}),
			accessorFn: (row: Record<string, JsonValue>) =>
				row['Annotations'] ? Object.keys(row['Annotations']).length : null
		},
		{
			id: 'Labels',
			header: ({ column }: { column: Column<Record<string, JsonValue>> }) =>
				renderComponent(DynamicTableHeader, {
					column: column,
					fields: fields
				}),
			cell: ({
				column,
				row
			}: {
				column: Column<Record<string, JsonValue>>;
				row: Row<Record<string, JsonValue>>;
			}) =>
				renderComponent(DynamicTableCell, {
					row: row,
					column: column,
					fields: fields
				}),
			accessorFn: (row: Record<string, JsonValue>) =>
				row['Labels'] ? Object.keys(row['Labels']).length : null
		},
		{
			id: 'CreateTime',
			header: ({ column }: { column: Column<Record<string, JsonValue>> }) =>
				renderComponent(DynamicTableHeader, {
					column: column,
					fields: fields
				}),
			cell: ({
				column,
				row
			}: {
				column: Column<Record<string, JsonValue>>;
				row: Row<Record<string, JsonValue>>;
			}) =>
				renderComponent(DynamicTableCell, {
					row: row,
					column: column,
					fields: fields
				}),
			accessorKey: 'CreateTime'
		},
		{
			id: 'Configuration',
			header: ({ column }: { column: Column<Record<string, JsonValue>> }) =>
				renderComponent(DynamicTableHeader, {
					column: column,
					fields: fields
				}),
			cell: ({
				column,
				row
			}: {
				column: Column<Record<string, JsonValue>>;
				row: Row<Record<string, JsonValue>>;
			}) =>
				renderComponent(DynamicTableCell, {
					row: row,
					column: column,
					fields: fields
				}),
			accessorKey: 'Configuration'
		}
	];
}

export { defaultColumnDefinitions, defaultFieldsMask, defaultObjectMask };
