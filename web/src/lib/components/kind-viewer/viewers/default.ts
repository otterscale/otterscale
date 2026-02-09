import { type JsonValue } from '@bufbuild/protobuf';
import type { Column, ColumnDef } from '@tanstack/table-core';
import { type Row } from '@tanstack/table-core';
import lodash from 'lodash';

import { resolve } from '$app/paths';
import { page } from '$app/state';
import type { APIResource } from '$lib/api/resource/v1/resource_pb';
import { DynamicTableCell, DynamicTableHeader } from '$lib/components/dynamic-table';
import LinkCell from '$lib/components/dynamic-table/cells/link-cell.svelte';
import { renderComponent } from '$lib/components/ui/data-table';
import type { FieldsType, ValuesType } from '../type';

function getDefaultFields(
	schema: any
): FieldsType {
	return {
		Name: lodash.get(schema, 'properties.metadata.properties.name'),
		Namespace: lodash.get(schema, 'properties.metadata.properties.namespace'),
		Labels: lodash.get(schema, 'properties.metadata.properties.labels'),
		Annotations: lodash.get(schema, 'properties.metadata.properties.annotations'),
		Age: lodash.get(schema, 'properties.metadata.properties.creationTimestamp'),
		Configuration: schema
	};
}
function getDefaultValues(object: any): ValuesType {
	return {
		Name: lodash.get(object, 'metadata.name'),
		Namespace: lodash.get(object, 'metadata.namespace'),
		Labels: lodash.get(object, 'metadata.labels'),
		Annotations: lodash.get(object, 'metadata.annotations'),
		Age: lodash.get(object, 'metadata.creationTimestamp'),
		Configuration: object
	};
}
function getDefaultColumnDefinitions(apiResource: APIResource, fields: FieldsType): ColumnDef<ValuesType>[] {
	return [
		{
			id: 'Name',
			header: ({ column }: { column: Column<ValuesType> }) =>
				renderComponent(DynamicTableHeader, {
					column: column,
					fields: fields
				}),
			cell: ({
				column,
				row
			}: {
				column: Column<ValuesType>;
				row: Row<ValuesType>;
			}) =>
				renderComponent(LinkCell, {
					display: String(row.original[column.id]),
					hyperlink: resolve(
						`/(auth)/${page.params.cluster!}/${apiResource.kind}/${apiResource.resource}?group=${apiResource.group}&version=${apiResource.version}&name=${row.original[column.id]}&namespace=${page.url.searchParams.get('namespace') ?? ''}`
					)
				}),
			accessorKey: 'Name'
		},
		...[
			{
				id: 'Namespace',
				header: ({ column }: { column: Column<ValuesType> }) =>
					renderComponent(DynamicTableHeader, {
						column: column,
						fields: fields
					}),
				cell: ({
					column,
					row
				}: {
					column: Column<ValuesType>;
					row: Row<ValuesType>;
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
			header: ({ column }: { column: Column<ValuesType> }) =>
				renderComponent(DynamicTableHeader, {
					column: column,
					fields: fields
				}),
			cell: ({
				column,
				row
			}: {
				column: Column<ValuesType>;
				row: Row<ValuesType>;
			}) =>
				renderComponent(DynamicTableCell, {
					row: row,
					column: column,
					fields: fields
				}),
			accessorFn: (row: ValuesType) =>
				row['Annotations'] ? Object.keys(row['Annotations']).length : null
		},
		{
			id: 'Labels',
			header: ({ column }: { column: Column<ValuesType> }) =>
				renderComponent(DynamicTableHeader, {
					column: column,
					fields: fields
				}),
			cell: ({
				column,
				row
			}: {
				column: Column<ValuesType>;
				row: Row<ValuesType>;
			}) =>
				renderComponent(DynamicTableCell, {
					row: row,
					column: column,
					fields: fields
				}),
			accessorFn: (row: ValuesType) =>
				row['Labels'] ? Object.keys(row['Labels']).length : null
		},
		{
<<<<<<< Updated upstream
			id: 'Age',
			header: ({ column }: { column: Column<Record<string, JsonValue>> }) =>
=======
			id: 'CreateTime',
			header: ({ column }: { column: Column<ValuesType> }) =>
>>>>>>> Stashed changes
				renderComponent(DynamicTableHeader, {
					column: column,
					fields: fields
				}),
			cell: ({
				column,
				row
			}: {
				column: Column<ValuesType>;
				row: Row<ValuesType>;
			}) =>
				renderComponent(DynamicTableCell, {
					row: row,
					column: column,
					fields: fields
				}),
			accessorKey: 'Age'
		},
		{
			id: 'Configuration',
			header: ({ column }: { column: Column<ValuesType> }) =>
				renderComponent(DynamicTableHeader, {
					column: column,
					fields: fields
				}),
			cell: ({
				column,
				row
			}: {
				column: Column<ValuesType>;
				row: Row<ValuesType>;
			}) =>
				renderComponent(DynamicTableCell, {
					row: row,
					column: column,
					fields: fields
				}),
			accessorKey: 'Configuration',
			meta: {
				class: 'hidden xl:table-cell'
			}
		}
	];
}

export { getDefaultColumnDefinitions, getDefaultFields, getDefaultValues };
