import { type JsonObject, type JsonValue } from '@bufbuild/protobuf';
import type { BatchV1CronJob } from '@otterscale/types';
import type { Column, ColumnDef } from '@tanstack/table-core';
import { type Row } from '@tanstack/table-core';

import { resolve } from '$app/paths';
import { page } from '$app/state';
import type { APIResource } from '$lib/api/resource/v1/resource_pb';
import { DynamicTableHeader } from '$lib/components/dynamic-table';
import DynamicTableCell from '$lib/components/dynamic-table/dynamic-table-cell.svelte';
import {
	type ArrayOfObjectItemsType,
	type ArrayOfObjectItemType,
	type ArrayOfObjectMetadata
} from '$lib/components/dynamic-table/dynamic-table-cells/array-of-object-cell.svelte';
import type { LinkMetadata } from '$lib/components/dynamic-table/dynamic-table-cells/link-cell.svelte';
import { type DataSchemaType, type UISchemaType } from '$lib/components/dynamic-table/utils';
import { renderComponent } from '$lib/components/ui/data-table';

type CronJobAttribute =
	| 'Name'
	| 'Namespace'
	| 'Schedule'
	| 'Time Zone'
	| 'Suspend'
	| 'Active'
	| 'Last Schedule'
	| 'Creation Timestamp'
	| 'Containers'
	| 'Images'
	| 'raw';

function getCronJobDataSchemas(): Record<CronJobAttribute, DataSchemaType> {
	return {
		Name: 'text',
		Namespace: 'text',
		Schedule: 'text',
		'Time Zone': 'text',
		Suspend: 'boolean',
		Active: 'number',
		'Last Schedule': 'time',
		'Creation Timestamp': 'time',
		Containers: 'number',
		Images: 'number',
		raw: 'object'
	};
}

function getCronJobData(object: BatchV1CronJob): Record<CronJobAttribute, JsonValue> {
	return {
		Name: object?.metadata?.name ?? null,
		Namespace: object?.metadata?.namespace ?? null,
		Schedule: object?.spec?.schedule ?? null,
		'Time Zone': object?.spec?.timeZone ?? null,
		Suspend: object?.spec?.suspend ?? null,
		Active: object?.status?.active?.length ?? null,
		'Last Schedule': object?.status?.lastScheduleTime ?? null,
		'Creation Timestamp': object?.metadata?.creationTimestamp ?? null,
		Containers: object?.spec?.jobTemplate?.spec?.template?.spec?.containers?.length ?? null,
		Images:
			new Set(
				object?.spec?.jobTemplate?.spec?.template?.spec?.containers?.map(
					(container) => container.image
				)
			).size ?? null,
		raw: (object as JsonObject) ?? null
	};
}

function getCronJobUISchemas(): Record<CronJobAttribute, UISchemaType> {
	return {
		Name: 'link',
		Namespace: 'text',
		Schedule: 'text',
		'Time Zone': 'text',
		Suspend: 'boolean',
		Active: 'array-of-object',
		'Last Schedule': 'time',
		'Creation Timestamp': 'time',
		Containers: 'array-of-object',
		Images: 'array-of-object',
		raw: 'object'
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
					} satisfies LinkMetadata
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
			id: 'Time Zone',
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
			accessorKey: 'Time Zone'
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
					uiSchemas: uiSchemas,
					metadata: {
						items:
							(row.original.raw as BatchV1CronJob).status?.active?.map(
								(job) =>
									({
										title: job?.name,
										description: job?.uid,
										actions: job?.resourceVersion
									}) as ArrayOfObjectItemType
							) ?? []
					} satisfies ArrayOfObjectMetadata
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
		},
		{
			id: 'Containers',
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
						items: (
							row.original.raw as BatchV1CronJob
						).spec?.jobTemplate.spec?.template.spec?.containers.map((container) => ({
							title: container.name,
							description: container.command?.join(' '),
							actions: container.image,
							raw: container
						})) as ArrayOfObjectItemsType
					} satisfies ArrayOfObjectMetadata
				}),
			accessorKey: 'Containers',
			meta: {
				class: 'hidden xl:table-cell'
			}
		},
		{
			id: 'Images',
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
						items: (
							row.original.raw as BatchV1CronJob
						).spec?.jobTemplate.spec?.template.spec?.containers.map((container) => ({
							title: container.image,
							description: container.name,
							raw: {
								image: container.image,
								imagePullPolicy: container.imagePullPolicy,
								container: container.name
							}
						})) as ArrayOfObjectItemsType
					} satisfies ArrayOfObjectMetadata
				}),
			accessorKey: 'Images',
			meta: {
				class: 'hidden xl:table-cell'
			}
		}
	];
}

export { getCronJobColumnDefinitions, getCronJobData, getCronJobDataSchemas, getCronJobUISchemas };
