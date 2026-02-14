import { type JsonValue } from '@bufbuild/protobuf';
import type { BatchV1CronJob } from '@otterscale/types';
import type { Column } from '@tanstack/table-core';
import { type Row } from '@tanstack/table-core';
import lodash from 'lodash';

import { resolve } from '$app/paths';
import { page } from '$app/state';
import type { APIResource } from '$lib/api/resource/v1/resource_pb';
import { DynamicTableCell, DynamicTableHeader } from '$lib/components/dynamic-table';
import ArrayOfObjectCell from '$lib/components/dynamic-table/cells/array-of-object-cell.svelte';
import LinkCell from '$lib/components/dynamic-table/cells/link-cell.svelte';
import { renderComponent } from '$lib/components/ui/data-table';

function cronjobFieldsMask(
	schema: any
): Record<string, { description: string; type: string; format?: string }> {
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

function cronjobObjectMask(object: BatchV1CronJob): Record<string, JsonValue> {
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

function cronjobColumnDefinitions(apiResource: APIResource, fields: any) {
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
				renderComponent(LinkCell, {
					display: String(row.original[column.id]),
					hyperlink: resolve(
						`/(auth)/${page.params.cluster!}/CronJob/cronjobs?group=batch&version=v1&name=${row.original[column.id]}&namespace=${page.url.searchParams.get('namespace')!}`
					)
				}),
			accessorKey: 'Name'
		},
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
		},
		{
			id: 'Schedule',
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
			accessorKey: 'Schedule'
		},
		{
			id: 'Suspend',
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
			accessorKey: 'Suspend'
		},
		{
			id: 'Active',
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
			accessorKey: 'Last Schedule'
		},
		{
			id: 'Creation Timestamp',
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
			accessorKey: 'Creation Timestamp'
		},
		{
			id: 'Images',
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
			accessorKey: 'Images',
			meta: {
				class: 'hidden xl:table-cell'
			}
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
			accessorKey: 'Configuration',
			meta: {
				class: 'hidden xl:table-cell'
			}
		}
	];
}

export { cronjobColumnDefinitions, cronjobFieldsMask, cronjobObjectMask };
