import { type JsonValue } from '@bufbuild/protobuf';
import type { CoreV1ResourceQuota } from '@otterscale/types';
import type { Column, ColumnDef } from '@tanstack/table-core';
import { type Row } from '@tanstack/table-core';

import { resolve } from '$app/paths';
import { page } from '$app/state';
import type { APIResource } from '$lib/api/resource/v1/resource_pb';
import { DynamicTableHeader } from '$lib/components/dynamic-table';
import DynamicTableCell from '$lib/components/dynamic-table/dynamic-table-cell.svelte';
import type { LinkMetadata } from '$lib/components/dynamic-table/dynamic-table-cells/link-cell.svelte';
import { type RatioMetadata } from '$lib/components/dynamic-table/dynamic-table-cells/ratio-cell.svelte';
import {
	type DataSchemaType,
	quantityToScalar,
	type UISchemaType
} from '$lib/components/dynamic-table/utils';
import { renderComponent } from '$lib/components/ui/data-table';

type ResourceQuotaAttribute =
	| 'Name'
	| 'Namespace'
	| 'CPU Limit'
	| 'Memory Limit'
	| 'CPU Request'
	| 'GPU Request'
	| 'Memory Request'
	| 'raw';

function getResourceQuotaUISchemas(): Record<ResourceQuotaAttribute, UISchemaType> {
	return {
		Name: 'link' as UISchemaType,
		Namespace: 'text' as UISchemaType,
		'CPU Limit': 'ratio' as UISchemaType,
		'Memory Limit': 'ratio' as UISchemaType,
		'CPU Request': 'ratio' as UISchemaType,
		'GPU Request': 'ratio' as UISchemaType,
		'Memory Request': 'ratio' as UISchemaType,
		raw: 'object' as UISchemaType
	};
}

function getResourceQuotaDataSchemas(): Record<ResourceQuotaAttribute, DataSchemaType> {
	return {
		Name: 'text' as DataSchemaType,
		Namespace: 'text' as DataSchemaType,
		'CPU Limit': 'number' as DataSchemaType,
		'Memory Limit': 'number' as DataSchemaType,
		'CPU Request': 'number' as DataSchemaType,
		'GPU Request': 'number' as DataSchemaType,
		'Memory Request': 'number' as DataSchemaType,
		raw: 'object' as DataSchemaType
	};
}

function getRatio(
	numerator: string | number | undefined,
	denominator: string | number | undefined
): JsonValue {
	if (numerator === undefined || denominator === undefined) return null;
	return (
		Number(quantityToScalar(String(numerator))) / Number(quantityToScalar(String(denominator)))
	);
}

// Ratio might lose some accuracy when resource is over 10Pi
function getResourceQuotaData(
	object: CoreV1ResourceQuota
): Record<ResourceQuotaAttribute, JsonValue> {
	return {
		Name: object?.metadata?.name as JsonValue,
		Namespace: object?.metadata?.namespace as JsonValue,
		'CPU Limit': getRatio(
			object?.status?.used?.['limits.cpu'],
			object?.status?.hard?.['limits.cpu']
		) as JsonValue,
		'Memory Limit': getRatio(
			object?.status?.used?.['limits.memory'],
			object?.status?.hard?.['limits.memory']
		) as JsonValue,
		'CPU Request': getRatio(
			object?.status?.used?.['requests.cpu'],
			object?.status?.hard?.['requests.cpu']
		) as JsonValue,
		'GPU Request': getRatio(
			object?.status?.used?.['requests.gpu'],
			object?.status?.hard?.['requests.gpu']
		) as JsonValue,
		'Memory Request': getRatio(
			object?.status?.used?.['requests.memory'],
			object?.status?.hard?.['requests.memory']
		) as JsonValue,
		raw: object as JsonValue
	};
}

function getResourceQuotaColumnDefinitions(
	apiResource: APIResource,
	uiSchemas: Record<ResourceQuotaAttribute, UISchemaType>,
	dataSchemas: Record<ResourceQuotaAttribute, DataSchemaType>
): ColumnDef<Record<ResourceQuotaAttribute, JsonValue>>[] {
	return [
		{
			id: 'Name',
			header: ({ column }: { column: Column<Record<ResourceQuotaAttribute, JsonValue>> }) =>
				renderComponent(DynamicTableHeader, {
					column: column,
					dataSchemas: dataSchemas
				}),
			cell: ({
				column,
				row
			}: {
				column: Column<Record<ResourceQuotaAttribute, JsonValue>>;
				row: Row<Record<ResourceQuotaAttribute, JsonValue>>;
			}) =>
				renderComponent(DynamicTableCell, {
					row: row,
					column: column,
					uiSchemas: uiSchemas,
					metadata: {
						hyperlink: resolve(
							`/(auth)/${page.params.cluster!}/${apiResource.kind}/${apiResource.resource}?group=${apiResource.group}&version=${apiResource.version}&name=${row.original[column.id as ResourceQuotaAttribute]}&namespace=${page.url.searchParams.get('namespace') ?? ''}`
						)
					} as LinkMetadata
				}),
			accessorKey: 'Name'
		},
		{
			id: 'Namespace',
			header: ({ column }: { column: Column<Record<ResourceQuotaAttribute, JsonValue>> }) =>
				renderComponent(DynamicTableHeader, {
					column: column,
					dataSchemas: dataSchemas
				}),
			cell: ({
				column,
				row
			}: {
				column: Column<Record<ResourceQuotaAttribute, JsonValue>>;
				row: Row<Record<ResourceQuotaAttribute, JsonValue>>;
			}) =>
				renderComponent(DynamicTableCell, {
					row: row,
					column: column,
					uiSchemas: uiSchemas
				}),
			accessorKey: 'Namespace'
		},
		{
			id: 'CPU Limit',
			header: ({ column }: { column: Column<Record<ResourceQuotaAttribute, JsonValue>> }) =>
				renderComponent(DynamicTableHeader, {
					column: column,
					dataSchemas: dataSchemas
				}),
			cell: ({
				column,
				row
			}: {
				column: Column<Record<ResourceQuotaAttribute, JsonValue>>;
				row: Row<Record<ResourceQuotaAttribute, JsonValue>>;
			}) =>
				renderComponent(DynamicTableCell, {
					row: row,
					column: column,
					uiSchemas: uiSchemas,
					metadata: {
						numerator: (row.original['raw'] as CoreV1ResourceQuota).status?.used?.['limits.cpu'],
						denominator: (row.original['raw'] as CoreV1ResourceQuota).status?.hard?.['limits.cpu']
					} as RatioMetadata
				}),
			accessorKey: 'CPU Limit',
			size: 100
		},
		{
			id: 'CPU Request',
			header: ({ column }: { column: Column<Record<ResourceQuotaAttribute, JsonValue>> }) =>
				renderComponent(DynamicTableHeader, {
					column: column,
					dataSchemas: dataSchemas
				}),
			cell: ({
				column,
				row
			}: {
				column: Column<Record<ResourceQuotaAttribute, JsonValue>>;
				row: Row<Record<ResourceQuotaAttribute, JsonValue>>;
			}) =>
				renderComponent(DynamicTableCell, {
					row: row,
					column: column,
					uiSchemas: uiSchemas,
					metadata: {
						numerator: (row.original['raw'] as CoreV1ResourceQuota).status?.used?.['requests.cpu'],
						denominator: (row.original['raw'] as CoreV1ResourceQuota).status?.hard?.['requests.cpu']
					} as RatioMetadata
				}),
			accessorKey: 'CPU Request',
			size: 100
		},
		{
			id: 'Memory Limit',
			header: ({ column }: { column: Column<Record<ResourceQuotaAttribute, JsonValue>> }) =>
				renderComponent(DynamicTableHeader, {
					column: column,
					dataSchemas: dataSchemas
				}),
			cell: ({
				column,
				row
			}: {
				column: Column<Record<ResourceQuotaAttribute, JsonValue>>;
				row: Row<Record<ResourceQuotaAttribute, JsonValue>>;
			}) =>
				renderComponent(DynamicTableCell, {
					row: row,
					column: column,
					uiSchemas: uiSchemas,
					metadata: {
						numerator: (row.original['raw'] as CoreV1ResourceQuota).status?.used?.['limits.memory'],
						denominator: (row.original['raw'] as CoreV1ResourceQuota).status?.hard?.[
							'limits.memory'
						]
					} as RatioMetadata
				}),
			accessorKey: 'Memory Limit',
			size: 100
		},
		{
			id: 'Memory Request',
			header: ({ column }: { column: Column<Record<ResourceQuotaAttribute, JsonValue>> }) =>
				renderComponent(DynamicTableHeader, {
					column: column,
					dataSchemas: dataSchemas
				}),
			cell: ({
				column,
				row
			}: {
				column: Column<Record<ResourceQuotaAttribute, JsonValue>>;
				row: Row<Record<ResourceQuotaAttribute, JsonValue>>;
			}) =>
				renderComponent(DynamicTableCell, {
					row: row,
					column: column,
					uiSchemas: uiSchemas,
					metadata: {
						numerator: (row.original['raw'] as CoreV1ResourceQuota).status?.used?.[
							'requests.memory'
						],
						denominator: (row.original['raw'] as CoreV1ResourceQuota).status?.hard?.[
							'requests.memory'
						]
					} as RatioMetadata
				}),
			accessorKey: 'Memory Request',
			size: 100
		},
		{
			id: 'GPU Request',
			header: ({ column }: { column: Column<Record<ResourceQuotaAttribute, JsonValue>> }) =>
				renderComponent(DynamicTableHeader, {
					column: column,
					dataSchemas: dataSchemas
				}),
			cell: ({
				column,
				row
			}: {
				column: Column<Record<ResourceQuotaAttribute, JsonValue>>;
				row: Row<Record<ResourceQuotaAttribute, JsonValue>>;
			}) =>
				renderComponent(DynamicTableCell, {
					row: row,
					column: column,
					uiSchemas: uiSchemas,
					metadata: {
						numerator: (row.original['raw'] as CoreV1ResourceQuota).status?.used?.[
							'requests.otterscale.com/vgpu'
						],
						denominator: (row.original['raw'] as CoreV1ResourceQuota).status?.hard?.[
							'requests.otterscale.com/vgpu'
						]
					} as RatioMetadata
				}),
			accessorKey: 'GPU Request',
			size: 100
		},
		{
			id: 'raw',
			header: ({ column }: { column: Column<Record<ResourceQuotaAttribute, JsonValue>> }) =>
				renderComponent(DynamicTableHeader, {
					column: column,
					dataSchemas: dataSchemas
				}),
			cell: ({
				column,
				row
			}: {
				column: Column<Record<ResourceQuotaAttribute, JsonValue>>;
				row: Row<Record<ResourceQuotaAttribute, JsonValue>>;
			}) =>
				renderComponent(DynamicTableCell, {
					row: row,
					column: column,
					uiSchemas: uiSchemas
				}),
			accessorKey: 'raw',
			meta: {
				class: 'hidden xl:table-cell'
			}
		}
	];
}

export {
	getResourceQuotaColumnDefinitions,
	getResourceQuotaData,
	getResourceQuotaDataSchemas,
	getResourceQuotaUISchemas
};
