import { type JsonValue } from '@bufbuild/protobuf';
import type { BatchV1CronJob } from '@otterscale/types';
import type { Column } from '@tanstack/table-core';
import { type Row } from '@tanstack/table-core';
import lodash from 'lodash';

import { page } from '$app/state';
import type { APIResource } from '$lib/api/resource/v1/resource_pb';
import { DynamicTableCell, DynamicTableHeader } from '$lib/components/dynamic-table';
import LinkCell from '$lib/components/dynamic-table/cells/link-cell.svelte';
import RatioCell from '$lib/components/dynamic-table/cells/ratio-cell.svelte';
import { renderComponent } from '$lib/components/ui/data-table';

function resourceQuotaFieldsMask(schema: any): Record<string, { type: string; format: string }> {
	return {
		Name: lodash.get(schema, 'properties.metadata.properties.name'),
		Namespace: lodash.get(schema, 'properties.metadata.properties.namespace'),
		Labels: lodash.get(schema, 'properties.metadata.properties.labels'),
		Annotations: lodash.get(schema, 'properties.metadata.properties.annotations'),
		'CPU Limit': {
			type: 'object',
			format: 'ratio'
		},
		'Creation Timestamp': lodash.get(schema, 'properties.metadata.properties.creationTimestamp'),
		Configuration: schema
	};
}

function resourceQuotaObjectMask(object: BatchV1CronJob): Record<string, JsonValue | undefined> {
	return {
		Name: object?.metadata?.name,
		Namespace: object?.metadata?.namespace,
		Labels: lodash.get(object, 'metadata.labels'),
		Annotations: lodash.get(object, 'metadata.annotations'),
		'CPU Limit': {
			numerator: Number(object['status']['used']['limits.cpu']),
			denominator: Number(object['status']['hard']['limits.cpu'])
		},
		'Creation Timestamp': object?.metadata?.creationTimestamp,
		Configuration: object as JsonValue
	};
}

function resourceQuotaColumnDefinitions(apiResource: APIResource, fields: any) {
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
					hyperlink: `/(auth)/${page.params.cluster!}/ResourceQuota/resourcequotas?group=&version=v1&name=${row.original[column.id]}&namespace=${page.url.searchParams.get('namespace')!}`
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
			id: 'CPU Limit',
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
				renderComponent(RatioCell, {
					numerator: row.original[column.id].numerator,
					denominator: row.original[column.id].denominator
				}),
			accessorKey: 'Creation Timestamp'
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

export { resourceQuotaColumnDefinitions, resourceQuotaFieldsMask, resourceQuotaObjectMask };
