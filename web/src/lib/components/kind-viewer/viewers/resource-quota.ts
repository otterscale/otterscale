import { type JsonValue } from '@bufbuild/protobuf';
import type { CoreV1ResourceQuota } from '@otterscale/types';
import type { Column } from '@tanstack/table-core';
import { type Row } from '@tanstack/table-core';
import lodash from 'lodash';

import { resolve } from '$app/paths';
import { page } from '$app/state';
import type { APIResource } from '$lib/api/resource/v1/resource_pb';
import { DynamicTableCell, DynamicTableHeader } from '$lib/components/dynamic-table';
import LinkCell from '$lib/components/dynamic-table/cells/link-cell.svelte';
import RatioCell from '$lib/components/dynamic-table/cells/ratio-cell.svelte';
import type { RatioType } from '$lib/components/dynamic-table/cells/type';
import { renderComponent } from '$lib/components/ui/data-table';

function resourceQuotaFieldsMask(
	schema: any
): Record<string, { description: string; type: string; format?: string }> {
	return {
		Name: lodash.get(schema, 'properties.metadata.properties.name'),
		Namespace: lodash.get(schema, 'properties.metadata.properties.namespace'),
		Labels: lodash.get(schema, 'properties.metadata.properties.labels'),
		Annotations: lodash.get(schema, 'properties.metadata.properties.annotations'),
		'CPU Limit': {
			description: lodash.get(
				schema,
				'properties.status.properties.used.properties["limits.cpu"].description'
			),
			type: 'object',
			format: 'ratio'
		},
		'Memory Limit': {
			description: lodash.get(
				schema,
				'properties.status.properties.used.properties["limits.memory"].description'
			),
			type: 'object',
			format: 'ratio'
		},
		'CPU Request': {
			description: lodash.get(
				schema,
				'properties.status.properties.used.properties["requests.cpu"].description'
			),
			type: 'object',
			format: 'ratio'
		},
		'GPU Request': {
			description: lodash.get(
				schema,
				'properties.status.properties.used.properties["requests.otterscale.com/vgpu"].description'
			),
			type: 'object',
			format: 'ratio'
		},
		'Memory Request': {
			description: lodash.get(
				schema,
				'properties.status.properties.used.properties["requests.memory"].description'
			),
			type: 'object',
			format: 'ratio'
		},
		Configuration: schema,
		raw: schema
	};
}

function resourceQuotaObjectMask(
	object: CoreV1ResourceQuota
): Record<string, JsonValue | RatioType> {
	return {
		Name: object?.metadata?.name as JsonValue,
		Namespace: object?.metadata?.namespace as JsonValue,
		Labels: object?.metadata?.labels as JsonValue,
		Annotations: object?.metadata?.annotations as JsonValue,
		'CPU Limit': {
			numerator: (object as any)['status']['used']['limits.cpu'],
			denominator: (object as any)['status']['hard']['limits.cpu']
		} as RatioType,
		'Memory Limit': {
			numerator: (object as any)['status']['used']['limits.memory'],
			denominator: (object as any)['status']['hard']['limits.memory']
		} as RatioType,
		'CPU Request': {
			numerator: (object as any)['status']['used']['requests.cpu'],
			denominator: (object as any)['status']['hard']['requests.cpu']
		} as RatioType,
		'GPU Request': {
			numerator: (object as any)['status']['used']['requests.otterscale.com/vgpu'],
			denominator: (object as any)['status']['hard']['requests.otterscale.com/vgpu']
		} as RatioType,
		'Memory Request': {
			numerator: (object as any)['status']['used']['requests.memory'],
			denominator: (object as any)['status']['hard']['requests.memory']
		} as RatioType,
		Configuration: object as JsonValue,
		raw: object as JsonValue
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
					hyperlink: resolve(
						`/(auth)/${page.params.cluster!}/ResourceQuota/resourcequotas?group=&version=v1&name=${row.original[column.id]}&namespace=${page.url.searchParams.get('namespace')!}`
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
					numerator: (row.original[column.id] as RatioType).numerator,
					denominator: (row.original[column.id] as RatioType).denominator
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
					numerator: (row.original[column.id] as RatioType).numerator,
					denominator: (row.original[column.id] as RatioType).denominator
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
					numerator: (row.original[column.id] as RatioType).numerator,
					denominator: (row.original[column.id] as RatioType).denominator
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
					numerator: (row.original[column.id] as RatioType).numerator,
					denominator: (row.original[column.id] as RatioType).denominator
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
					numerator: (row.original[column.id] as RatioType).numerator,
					denominator: (row.original[column.id] as RatioType).denominator
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
