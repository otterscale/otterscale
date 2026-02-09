import { type JsonValue } from '@bufbuild/protobuf';
import type { BatchV1CronJob } from '@otterscale/types';
import type { Column, ColumnDef } from '@tanstack/table-core';
import { type Row } from '@tanstack/table-core';
import lodash from 'lodash';

import { resolve } from '$app/paths';
import { page } from '$app/state';
import { DynamicTableCell, DynamicTableHeader } from '$lib/components/dynamic-table';
import ArrayOfObjectCell from '$lib/components/dynamic-table/cells/array-of-object-cell.svelte';
import LinkCell from '$lib/components/dynamic-table/cells/link-cell.svelte';
import { renderComponent } from '$lib/components/ui/data-table';
import type { FieldsType, ValuesType } from '../type';
import type { APIResource } from '$lib/api/resource/v1/resource_pb';

function getCronJobFields(
	schema: any
): FieldsType {
	return {
		Name: lodash.get(schema, 'properties.metadata.properties.name'),
		Namespace: lodash.get(schema, 'properties.metadata.properties.namespace'),
		Schedule: lodash.get(schema, 'properties.spec.properties.schedule'),
		Suspend: lodash.get(schema, 'properties.spec.properties.suspend'),
		Active: lodash.get(schema, 'properties.status.properties.active'),
		'Last Schedule': lodash.get(schema, 'properties.status.properties.lastScheduleTime'),
		'Creation Timestamp': lodash.get(schema, 'properties.metadata.properties.creationTimestamp'),
		Images: {
			description: lodash.get(
				schema,
				'properties.spec.properties.jobTemplate.properties.spec.properties.template.properties.spec.properties.containers.items.properties.image.description'
			),
			type: 'array'
		},
		Configuration: schema
	};
}

function getCronJobValues(object: BatchV1CronJob): ValuesType {
	return {
		Name: object?.metadata?.name as JsonValue,
		Namespace: object?.metadata?.namespace as JsonValue,
		Schedule: object?.spec?.schedule as JsonValue,
		Suspend: object?.spec?.suspend as JsonValue,
		Active: object?.status?.active as JsonValue,
		'Last Schedule': object?.status?.lastScheduleTime as JsonValue,
		'Creation Timestamp': object?.metadata?.creationTimestamp as JsonValue,
		Images: object?.spec?.jobTemplate?.spec?.template?.spec?.containers.map(
			(container) => container.image
		) as JsonValue,
		Configuration: object as JsonValue
	};
}

function getCronJobColumnDefinitions(apiResource: APIResource, fields: FieldsType): ColumnDef<ValuesType>[] {
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
						`/(auth)/${page.params.cluster!}/${apiResource.kind}/${apiResource.resource}?group=${apiResource.group}&version=${apiResource.version}&name=${row.original[column.id]}&namespace=${page.url.searchParams.get('namespace')!}`
					)
				}),
			accessorKey: 'Name'
		},
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
		},
		{
			id: 'Schedule',
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
			accessorKey: 'Schedule'
		},
		{
			id: 'Suspend',
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
			accessorKey: 'Suspend'
		},
		{
			id: 'Active',
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
				renderComponent(ArrayOfObjectCell, {
					keys: {
						title: 'name',
						description: 'uid',
						actions: 'resourceVersion'
					},
					row: row,
					column: column,
					fields: fields
				}),
			accessorKey: 'Active'
		},
		{
			id: 'Last Schedule',
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
			accessorKey: 'Last Schedule'
		},
		{
			id: 'Creation Timestamp',
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
			accessorKey: 'Creation Timestamp'
		},
		{
			id: 'Images',
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
			accessorKey: 'Images',
			meta: {
				class: 'hidden xl:table-cell'
			}
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

export { getCronJobFields, getCronJobValues, getCronJobColumnDefinitions };
