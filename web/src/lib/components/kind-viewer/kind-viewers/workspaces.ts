import { type JsonValue } from '@bufbuild/protobuf';
import type { TenantOtterscaleIoV1Alpha1Workspace } from '@otterscale/types';
import type { Column, ColumnDef } from '@tanstack/table-core';
import { type Row } from '@tanstack/table-core';

import { resolve } from '$app/paths';
import { page } from '$app/state';
import type { APIResource } from '$lib/api/resource/v1/resource_pb';
import { DynamicTableHeader } from '$lib/components/dynamic-table';
import DynamicTableCell from '$lib/components/dynamic-table/dynamic-table-cell.svelte';
import type { ArrayOfObjectMetadata } from '$lib/components/dynamic-table/dynamic-table-cells/array-of-object-cell.svelte';
import type { LinkMetadata } from '$lib/components/dynamic-table/dynamic-table-cells/link-cell.svelte';
import { type DataSchemaType, type UISchemaType } from '$lib/components/dynamic-table/utils';
import { renderComponent } from '$lib/components/ui/data-table';

type WorkspaceAttribute = 'Name' | 'Namespace' | 'Users' | 'Creation Timestamp' | 'raw';

function getWorkspaceUISchemas(): Record<WorkspaceAttribute, UISchemaType> {
	return {
		Name: 'link' as UISchemaType,
		Namespace: 'text' as UISchemaType,
		Users: 'array-of-object' as UISchemaType,
		'Creation Timestamp': 'time' as UISchemaType,
		raw: 'object' as UISchemaType
	};
}

function getWorkspaceDataSchemas(): Record<WorkspaceAttribute, DataSchemaType> {
	return {
		Name: 'text' as DataSchemaType,
		Namespace: 'text' as DataSchemaType,
		Users: 'number' as DataSchemaType,
		'Creation Timestamp': 'time' as DataSchemaType,
		raw: 'object' as DataSchemaType
	};
}

function getWorkspaceData(
	object: TenantOtterscaleIoV1Alpha1Workspace
): Record<WorkspaceAttribute, JsonValue> {
	return {
		Name: object?.metadata?.name as JsonValue,
		Namespace: object?.spec?.namespace as JsonValue,
		Users: object?.spec?.users?.length as JsonValue,
		'Creation Timestamp': object?.metadata?.creationTimestamp as JsonValue,
		raw: object as JsonValue
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
					uiSchemas: uiSchemas
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
							(user) => ({
								title: user?.name,
								description: user?.subject,
								actions: user?.role
							})
						)
					} as ArrayOfObjectMetadata
				}),
			accessorKey: 'Users'
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
