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
		'Memory Limit': {
			type: 'object',
			format: 'ratio'
		},
		'CPU Request': {
			type: 'object',
			format: 'ratio'
		},
		'GPU Request': {
			type: 'object',
			format: 'ratio'
		},
		'Memory Request': {
			type: 'object',
			format: 'ratio'
		},
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
			numerator: object['status']['used']['limits.cpu'],
			denominator: object['status']['hard']['limits.cpu']
		},
		'Memory Limit': {
			numerator: object['status']['used']['limits.memory'],
			denominator: object['status']['hard']['limits.memory']
		},
		'CPU Request': {
			numerator: object['status']['used']['requests.cpu'],
			denominator: object['status']['hard']['requests.cpu']
		},
		'GPU Request': {
			numerator: object['status']['used']['requests.otterscale.com/vgpu'],
			denominator: object['status']['hard']['requests.otterscale.com/vgpu']
		},
		'Memory Request': {
			numerator: object['status']['used']['requests.memory'],
			denominator: object['status']['hard']['requests.memory']
		},
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
			accessorKey: 'CPU Limit'
		},
		{
			id: 'Memory Limit',
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
			accessorKey: 'Memory Limit'
		},
		{
			id: 'CPU Request',
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
			accessorKey: 'CPU Request'
		},
		{
			id: 'GPU Request',
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
			accessorKey: 'GPU Request'
		},
		{
			id: 'Memory Request',
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
			accessorKey: 'Memory Request'
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
