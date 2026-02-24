import { type JsonObject, type JsonValue } from '@bufbuild/protobuf';
import type { TenantOtterscaleIoV1Alpha1Workspace } from '@otterscale/types';
import type { Column, ColumnDef } from '@tanstack/table-core';
import { type Row } from '@tanstack/table-core';

import { resolve } from '$app/paths';
import { page } from '$app/state';
import type { APIResource } from '$lib/api/resource/v1/resource_pb';
import { DynamicTableHeader } from '$lib/components/dynamic-table';
import DynamicTableCell from '$lib/components/dynamic-table/dynamic-table-cell.svelte';
import type {
	ArrayOfObjectItemType,
	ArrayOfObjectMetadata
} from '$lib/components/dynamic-table/dynamic-table-cells/array-of-object-cell.svelte';
import type { LinkMetadata } from '$lib/components/dynamic-table/dynamic-table-cells/link-cell.svelte';
import type { NumberWithPrefixMetadata } from '$lib/components/dynamic-table/dynamic-table-cells/number-with-prefix-cell.svelte';
import {
	type DataSchemaType,
	getQuantityScalar,
	type UISchemaType
} from '$lib/components/dynamic-table/utils';
import { renderComponent } from '$lib/components/ui/data-table';

type WorkspaceAttribute =
	| 'Name'
	| 'Namespace'
	| 'Users'
	| 'CPU Limit'
	| 'CPU Requests'
	| 'Memory Limit'
	| 'Memory Requests'
	| 'GPU Requests'
	| 'Creation Timestamp'
	| 'raw';

function getWorkspaceDataSchemas(): Record<WorkspaceAttribute, DataSchemaType> {
	return {
		Name: 'text',
		Namespace: 'text',
		Users: 'number',
		'CPU Limit': 'number',
		'CPU Requests': 'number',
		'Memory Limit': 'number',
		'Memory Requests': 'number',
		'GPU Requests': 'number',
		'Creation Timestamp': 'time',
		raw: 'object'
	};
}

function getWorkspaceData(
	object: TenantOtterscaleIoV1Alpha1Workspace
): Record<WorkspaceAttribute, JsonValue | bigint> {
	return {
		Name: object?.metadata?.name ?? null,
		Namespace: object?.spec?.namespace ?? null,
		Users: (object?.spec?.users ?? []).length,
		'CPU Limit': getQuantityScalar(object?.spec?.resourceQuota?.hard?.['limits.cpu'] ?? null),
		'CPU Requests': getQuantityScalar(object?.spec?.resourceQuota?.hard?.['requests.cpu'] ?? null),
		'Memory Limit': getQuantityScalar(object?.spec?.resourceQuota?.hard?.['limits.memory'] ?? null),
		'Memory Requests': getQuantityScalar(
			object?.spec?.resourceQuota?.hard?.['requests.memory'] ?? null
		),
		'GPU Requests': getQuantityScalar(
			object?.spec?.resourceQuota?.hard?.['requests.otterscale.com/vgpu'] ?? null
		),
		'Creation Timestamp': object?.metadata?.creationTimestamp as JsonValue,
		raw: object as JsonObject
	};
}

function getWorkspaceUISchemas(): Record<WorkspaceAttribute, UISchemaType> {
	return {
		Name: 'link',
		Namespace: 'link',
		Users: 'array-of-object',
		'CPU Limit': 'number-with-prefix',
		'CPU Requests': 'number-with-prefix',
		'Memory Limit': 'number-with-prefix',
		'Memory Requests': 'number-with-prefix',
		'GPU Requests': 'number-with-prefix',
		'Creation Timestamp': 'time',
		raw: 'object'
	};
}

function getWorkspaceColumnDefinitions(
	apiResource: APIResource,
	uiSchemas: Record<WorkspaceAttribute, UISchemaType>,
	dataSchemas: Record<string, DataSchemaType>
): ColumnDef<Record<WorkspaceAttribute, JsonValue>>[] {
	return [
		{
			id: 'Name',
			header: ({ column }: { column: Column<Record<WorkspaceAttribute, JsonValue>> }) =>
				renderComponent(DynamicTableHeader, {
					column: column,
					dataSchemas: dataSchemas
				}),
			cell: ({
				column,
				row
			}: {
				column: Column<Record<WorkspaceAttribute, JsonValue>>;
				row: Row<Record<WorkspaceAttribute, JsonValue>>;
			}) =>
				renderComponent(DynamicTableCell, {
					row: row,
					column: column,
					uiSchemas: uiSchemas,
					metadata: {
						hyperlink: resolve(
							`/(auth)/${page.params.cluster!}/${apiResource.kind}/${apiResource.resource}?group=${apiResource.group}&version=${apiResource.version}&name=${row.original[column.id as WorkspaceAttribute]}`
						)
					} as LinkMetadata
				}),
			accessorKey: 'Name'
		},
		{
			id: 'Namespace',
			header: ({ column }: { column: Column<Record<WorkspaceAttribute, JsonValue>> }) =>
				renderComponent(DynamicTableHeader, {
					column: column,
					dataSchemas: dataSchemas
				}),
			cell: ({
				column,
				row
			}: {
				column: Column<Record<WorkspaceAttribute, JsonValue>>;
				row: Row<Record<WorkspaceAttribute, JsonValue>>;
			}) =>
				renderComponent(DynamicTableCell, {
					row: row,
					column: column,
					uiSchemas: uiSchemas,
					metadata: {
						hyperlink: resolve(
							`/(auth)/${page.params.cluster!}/Namespace/namespaces?group=&version=v1&name=${row.original['Namespace']}`
						)
					} satisfies LinkMetadata
				}),
			accessorKey: 'Namespace'
		},
		{
			id: 'Users',
			header: ({ column }: { column: Column<Record<WorkspaceAttribute, JsonValue>> }) =>
				renderComponent(DynamicTableHeader, {
					column: column,
					dataSchemas: dataSchemas
				}),
			cell: ({
				column,
				row
			}: {
				column: Column<Record<WorkspaceAttribute, JsonValue>>;
				row: Row<Record<WorkspaceAttribute, JsonValue>>;
			}) =>
				renderComponent(DynamicTableCell, {
					row: row,
					column: column,
					uiSchemas: uiSchemas,
					metadata: {
						items: (row.original.raw as TenantOtterscaleIoV1Alpha1Workspace).spec?.users?.map(
							(user) =>
								({
									title: user?.name,
									description: user?.subject,
									actions: user?.role
								}) as ArrayOfObjectItemType
						)
					} satisfies ArrayOfObjectMetadata
				}),
			accessorKey: 'Users'
		},
		{
			id: 'CPU Limit',
			header: ({ column }: { column: Column<Record<WorkspaceAttribute, JsonValue>> }) =>
				renderComponent(DynamicTableHeader, {
					column: column,
					dataSchemas: dataSchemas
				}),
			cell: ({
				column,
				row
			}: {
				column: Column<Record<WorkspaceAttribute, JsonValue>>;
				row: Row<Record<WorkspaceAttribute, JsonValue>>;
			}) =>
				renderComponent(DynamicTableCell, {
					row: row,
					column: column,
					uiSchemas: uiSchemas,
					metadata: {
						prefix: 'binary'
					} satisfies NumberWithPrefixMetadata
				}),
			accessorKey: 'CPU Limit'
		},
		{
			id: 'CPU Requests',
			header: ({ column }: { column: Column<Record<WorkspaceAttribute, JsonValue>> }) =>
				renderComponent(DynamicTableHeader, {
					column: column,
					dataSchemas: dataSchemas
				}),
			cell: ({
				column,
				row
			}: {
				column: Column<Record<WorkspaceAttribute, JsonValue>>;
				row: Row<Record<WorkspaceAttribute, JsonValue>>;
			}) =>
				renderComponent(DynamicTableCell, {
					row: row,
					column: column,
					uiSchemas: uiSchemas,
					metadata: {
						prefix: 'binary'
					} satisfies NumberWithPrefixMetadata
				}),
			accessorKey: 'CPU Requests'
		},
		{
			id: 'Memory Limit',
			header: ({ column }: { column: Column<Record<WorkspaceAttribute, JsonValue>> }) =>
				renderComponent(DynamicTableHeader, {
					column: column,
					dataSchemas: dataSchemas
				}),
			cell: ({
				column,
				row
			}: {
				column: Column<Record<WorkspaceAttribute, JsonValue>>;
				row: Row<Record<WorkspaceAttribute, JsonValue>>;
			}) =>
				renderComponent(DynamicTableCell, {
					row: row,
					column: column,
					uiSchemas: uiSchemas,
					metadata: {
						prefix: 'binary'
					} satisfies NumberWithPrefixMetadata
				}),
			accessorKey: 'Memory Limit'
		},
		{
			id: 'Memory Requests',
			header: ({ column }: { column: Column<Record<WorkspaceAttribute, JsonValue>> }) =>
				renderComponent(DynamicTableHeader, {
					column: column,
					dataSchemas: dataSchemas
				}),
			cell: ({
				column,
				row
			}: {
				column: Column<Record<WorkspaceAttribute, JsonValue>>;
				row: Row<Record<WorkspaceAttribute, JsonValue>>;
			}) =>
				renderComponent(DynamicTableCell, {
					row: row,
					column: column,
					uiSchemas: uiSchemas,
					metadata: {
						prefix: 'binary'
					} satisfies NumberWithPrefixMetadata
				}),
			accessorKey: 'Memory Requests'
		},
		{
			id: 'GPU Requests',
			header: ({ column }: { column: Column<Record<WorkspaceAttribute, JsonValue>> }) =>
				renderComponent(DynamicTableHeader, {
					column: column,
					dataSchemas: dataSchemas
				}),
			cell: ({
				column,
				row
			}: {
				column: Column<Record<WorkspaceAttribute, JsonValue>>;
				row: Row<Record<WorkspaceAttribute, JsonValue>>;
			}) =>
				renderComponent(DynamicTableCell, {
					row: row,
					column: column,
					uiSchemas: uiSchemas,
					metadata: {
						prefix: 'binary'
					} satisfies NumberWithPrefixMetadata
				}),
			accessorKey: 'GPU Requests'
		},
		{
			id: 'Creation Timestamp',
			header: ({ column }: { column: Column<Record<WorkspaceAttribute, JsonValue>> }) =>
				renderComponent(DynamicTableHeader, {
					column: column,
					dataSchemas: dataSchemas
				}),
			cell: ({
				column,
				row
			}: {
				column: Column<Record<WorkspaceAttribute, JsonValue>>;
				row: Row<Record<WorkspaceAttribute, JsonValue>>;
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

export {
	getWorkspaceColumnDefinitions,
	getWorkspaceData,
	getWorkspaceDataSchemas,
	getWorkspaceUISchemas
};
