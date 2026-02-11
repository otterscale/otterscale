import { type JsonValue } from '@bufbuild/protobuf';
import type { BatchV1CronJob } from '@otterscale/types';
import type { Column, ColumnDef } from '@tanstack/table-core';
import { type Row } from '@tanstack/table-core';

import { resolve } from '$app/paths';
import { page } from '$app/state';
import type { APIResource } from '$lib/api/resource/v1/resource_pb';
import { DynamicTableHeader } from '$lib/components/dynamic-table';
import DynamicTableCell from '$lib/components/dynamic-table/dynamic-table-cell.svelte';
import { type ArrayOfObjectMetadata } from '$lib/components/dynamic-table/dynamic-table-cells/array-of-object-cell.svelte';
import type { LinkMetadata } from '$lib/components/dynamic-table/dynamic-table-cells/link-cell.svelte';
import { type DataSchemaType, type UISchemaType } from '$lib/components/dynamic-table/utils';
import { renderComponent } from '$lib/components/ui/data-table';

type CronJobAttribute =
	| 'Name'
	| 'Namespace'
	| 'Schedule'
	| 'Suspend'
	| 'Active'
	| 'Last Schedule'
	| 'Creation Timestamp'
	| 'raw';

function getCronJobUISchemas(): Record<CronJobAttribute, UISchemaType> {
	return {
		Name: 'link' as UISchemaType,
		Namespace: 'text' as UISchemaType,
		Schedule: 'text' as UISchemaType,
		Suspend: 'boolean' as UISchemaType,
		Active: 'array-of-object' as UISchemaType,
		'Last Schedule': 'time' as UISchemaType,
		'Creation Timestamp': 'time' as UISchemaType,
		raw: 'object' as UISchemaType
	};
}

function getCronJobDataSchemas(): Record<CronJobAttribute, DataSchemaType> {
	return {
		Name: 'text' as DataSchemaType,
		Namespace: 'text' as DataSchemaType,
		Schedule: 'text' as DataSchemaType,
		Suspend: 'boolean' as DataSchemaType,
		Active: 'number' as DataSchemaType,
		'Last Schedule': 'time' as DataSchemaType,
		'Creation Timestamp': 'time' as DataSchemaType,
		raw: 'object' as DataSchemaType
	};
}

function getCronJobData(object: BatchV1CronJob): Record<CronJobAttribute, JsonValue> {
	return {
		Name: object?.metadata?.name as JsonValue,
		Namespace: object?.metadata?.namespace as JsonValue,
		Schedule: object?.spec?.schedule as JsonValue,
		Suspend: object?.spec?.suspend as JsonValue,
		Active: object?.status?.active?.length as JsonValue,
		'Last Schedule': object?.status?.lastScheduleTime as JsonValue,
		'Creation Timestamp': object?.metadata?.creationTimestamp as JsonValue,
		raw: object as JsonValue
	};
}

function getCronJobColumnDefinitions(
	apiResource: APIResource,
	uiSchemas: Record<CronJobAttribute, UISchemaType>,
	dataSchemas: Record<CronJobAttribute, DataSchemaType>
): ColumnDef<Record<CronJobAttribute, JsonValue>>[] {
	return [
		{
			id: 'Name',
			header: ({ column }: { column: Column<Record<CronJobAttribute, JsonValue>> }) =>
				renderComponent(DynamicTableHeader, {
					column: column,
					dataSchemas: dataSchemas
				}),
			cell: ({
				column,
				row
			}: {
				column: Column<Record<CronJobAttribute, JsonValue>>;
				row: Row<Record<CronJobAttribute, JsonValue>>;
			}) =>
				renderComponent(DynamicTableCell, {
					row: row,
					column: column,
					uiSchemas: uiSchemas,
					metadata: {
						hyperlink: resolve(
							`/(auth)/${page.params.cluster!}/${apiResource.kind}/${apiResource.resource}?group=${apiResource.group}&version=${apiResource.version}&name=${row.original[column.id as CronJobAttribute]}&namespace=${page.url.searchParams.get('namespace')!}`
						)
					} as LinkMetadata
				}),
			accessorKey: 'Name'
		},
		{
			id: 'Namespace',
			header: ({ column }: { column: Column<Record<CronJobAttribute, JsonValue>> }) =>
				renderComponent(DynamicTableHeader, {
					column: column,
					dataSchemas: dataSchemas
				}),
			cell: ({
				column,
				row
			}: {
				column: Column<Record<CronJobAttribute, JsonValue>>;
				row: Row<Record<CronJobAttribute, JsonValue>>;
			}) =>
				renderComponent(DynamicTableCell, {
					row: row,
					column: column,
					uiSchemas: uiSchemas
				}),
			accessorKey: 'Namespace'
		},
		{
			id: 'Schedule',
			header: ({ column }: { column: Column<Record<CronJobAttribute, JsonValue>> }) =>
				renderComponent(DynamicTableHeader, {
					column: column,
					dataSchemas: dataSchemas
				}),
			cell: ({
				column,
				row
			}: {
				column: Column<Record<CronJobAttribute, JsonValue>>;
				row: Row<Record<CronJobAttribute, JsonValue>>;
			}) =>
				renderComponent(DynamicTableCell, {
					row: row,
					column: column,
					uiSchemas: uiSchemas
				}),
			accessorKey: 'Schedule'
		},
		{
			id: 'Suspend',
			header: ({ column }: { column: Column<Record<CronJobAttribute, JsonValue>> }) =>
				renderComponent(DynamicTableHeader, {
					column: column,
					dataSchemas: dataSchemas
				}),
			cell: ({
				column,
				row
			}: {
				column: Column<Record<CronJobAttribute, JsonValue>>;
				row: Row<Record<CronJobAttribute, JsonValue>>;
			}) =>
				renderComponent(DynamicTableCell, {
					row: row,
					column: column,
					uiSchemas: uiSchemas
				}),
			accessorKey: 'Suspend'
		},
		{
			id: 'Active',
			header: ({ column }: { column: Column<Record<CronJobAttribute, JsonValue>> }) =>
				renderComponent(DynamicTableHeader, {
					column: column,
					dataSchemas: dataSchemas
				}),
			cell: ({
				column,
				row
			}: {
				column: Column<Record<CronJobAttribute, JsonValue>>;
				row: Row<Record<CronJobAttribute, JsonValue>>;
			}) =>
				renderComponent(DynamicTableCell, {
					row: row,
					column: column,
					metadata: {
						items: (row.original.raw as BatchV1CronJob).status?.active?.map((job) => ({
							title: job?.name,
							description: job?.uid,
							actions: job?.resourceVersion
						}))
					} as ArrayOfObjectMetadata,
					uiSchemas: uiSchemas
				}),
			accessorKey: 'Active'
		},
		{
			id: 'Last Schedule',
			header: ({ column }: { column: Column<Record<CronJobAttribute, JsonValue>> }) =>
				renderComponent(DynamicTableHeader, {
					column: column,
					dataSchemas: dataSchemas
				}),
			cell: ({
				column,
				row
			}: {
				column: Column<Record<CronJobAttribute, JsonValue>>;
				row: Row<Record<CronJobAttribute, JsonValue>>;
			}) =>
				renderComponent(DynamicTableCell, {
					row: row,
					column: column,
					uiSchemas: uiSchemas
				}),
			accessorKey: 'Last Schedule'
		},
		{
			id: 'Creation Timestamp',
			header: ({ column }: { column: Column<Record<CronJobAttribute, JsonValue>> }) =>
				renderComponent(DynamicTableHeader, {
					column: column,
					dataSchemas: dataSchemas
				}),
			cell: ({
				column,
				row
			}: {
				column: Column<Record<CronJobAttribute, JsonValue>>;
				row: Row<Record<CronJobAttribute, JsonValue>>;
			}) =>
				renderComponent(DynamicTableCell, {
					row: row,
					column: column,
					uiSchemas: uiSchemas
				}),
			accessorKey: 'Creation Timestamp'
		}
		// {
		// 	id: 'Images',
		// 	header: ({ column }: { column: Column<Record<CronJobAttribute, JsonValue>> }) =>
		// 		renderComponent(DynamicTableHeader, {
		// 			column: column,
		// 			uiSchemas: uiSchemas
		// 		}),
		// 	cell: ({
		// 		column,
		// 		row
		// 	}: {
		// 		column: Column<Record<CronJobAttribute, JsonValue>>;
		// 		row: Row<Record<CronJobAttribute, JsonValue>>;
		// 	}) =>
		// 		renderComponent(DefaultObjectCell, {
		// 			row: row,
		// 			column: column
		// 		}),
		// 	accessorKey: 'Images',
		// 	meta: {
		// 		class: 'hidden xl:table-cell'
		// 	}
		// },
	];
}

export { getCronJobColumnDefinitions, getCronJobData, getCronJobDataSchemas, getCronJobUISchemas };
