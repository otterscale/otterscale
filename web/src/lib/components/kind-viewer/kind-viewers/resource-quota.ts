import { type JsonObject, type JsonValue } from '@bufbuild/protobuf';
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
	getQuantityScalar,
	getRatio,
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

function getResourceQuotaDataSchemas(): Record<ResourceQuotaAttribute, DataSchemaType> {
	return {
		Name: 'text',
		Namespace: 'text',
		'CPU Limit': 'number',
		'Memory Limit': 'number',
		'CPU Request': 'number',
		'GPU Request': 'number',
		'Memory Request': 'number',
		raw: 'object'
	};
}

function getResourceQuotaData(
	object: CoreV1ResourceQuota
): Record<ResourceQuotaAttribute, JsonValue> {
	return {
		Name: object?.metadata?.name ?? null,
		Namespace: object?.metadata?.namespace ?? null,
		'CPU Limit': getRatio(
			getQuantityScalar(object?.status?.used?.['limits.cpu'] ?? null),
			getQuantityScalar(object?.status?.hard?.['limits.cpu'] ?? null)
		),
		'Memory Limit': getRatio(
			getQuantityScalar(object?.status?.used?.['limits.memory'] ?? null),
			getQuantityScalar(object?.status?.hard?.['limits.memory'] ?? null)
		),
		'CPU Request': getRatio(
			getQuantityScalar(object?.status?.used?.['requests.cpu'] ?? null),
			getQuantityScalar(object?.status?.hard?.['requests.cpu'] ?? null)
		),
		'GPU Request': getRatio(
			getQuantityScalar(object?.status?.used?.['requests.gpu'] ?? null),
			getQuantityScalar(object?.status?.hard?.['requests.gpu'] ?? null)
		),
		'Memory Request': getRatio(
			getQuantityScalar(object?.status?.used?.['requests.memory'] ?? null),
			getQuantityScalar(object?.status?.hard?.['requests.memory'] ?? null)
		),
		raw: object as JsonObject
	};
}

function getResourceQuotaUISchemas(): Record<ResourceQuotaAttribute, UISchemaType> {
	return {
		Name: 'link',
		Namespace: 'text',
		'CPU Limit': 'ratio',
		'Memory Limit': 'ratio',
		'CPU Request': 'ratio',
		'GPU Request': 'ratio',
		'Memory Request': 'ratio',
		raw: 'object'
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
					} satisfies LinkMetadata
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
						numerator:
							(row.original['raw'] as CoreV1ResourceQuota).status?.used?.['limits.cpu'] ?? null,
						denominator:
							(row.original['raw'] as CoreV1ResourceQuota).status?.hard?.['limits.cpu'] ?? null
					} satisfies RatioMetadata
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
						numerator:
							(row.original['raw'] as CoreV1ResourceQuota).status?.used?.['requests.cpu'] ?? null,
						denominator:
							(row.original['raw'] as CoreV1ResourceQuota).status?.hard?.['requests.cpu'] ?? null
					} satisfies RatioMetadata
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
						numerator:
							(row.original['raw'] as CoreV1ResourceQuota).status?.used?.['limits.memory'] ?? null,
						denominator:
							(row.original['raw'] as CoreV1ResourceQuota).status?.hard?.['limits.memory'] ?? null
					} satisfies RatioMetadata
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
						numerator:
							(row.original['raw'] as CoreV1ResourceQuota).status?.used?.['requests.memory'] ??
							null,
						denominator:
							(row.original['raw'] as CoreV1ResourceQuota).status?.hard?.['requests.memory'] ?? null
					} satisfies RatioMetadata
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
						numerator:
							(row.original['raw'] as CoreV1ResourceQuota).status?.used?.[
								'requests.otterscale.com/vgpu'
							] ?? null,
						denominator:
							(row.original['raw'] as CoreV1ResourceQuota).status?.hard?.[
								'requests.otterscale.com/vgpu'
							] ?? null
					} satisfies RatioMetadata
				}),
			accessorKey: 'GPU Request',
			size: 100
		}
	];
}

export {
	getResourceQuotaColumnDefinitions,
	getResourceQuotaData,
	getResourceQuotaDataSchemas,
	getResourceQuotaUISchemas
};
