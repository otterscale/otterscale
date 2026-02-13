import { type JsonValue } from '@bufbuild/protobuf';
import type { Column } from '@tanstack/table-core';
import { type Row } from '@tanstack/table-core';
import lodash from 'lodash';

import { resolve } from '$app/paths';
import { page } from '$app/state';
import type { APIResource } from '$lib/api/resource/v1/resource_pb';
import { DynamicTableCell, DynamicTableHeader } from '$lib/components/dynamic-table';
import LinkCell from '$lib/components/dynamic-table/cells/link-cell.svelte';
import { renderComponent } from '$lib/components/ui/data-table';

function simpleappFieldsMask(
	schema: any
): Record<string, { description: string; type: string; format?: string }> {
	return {
		Name: lodash.get(schema, 'properties.metadata.properties.name'),
		Namespace: lodash.get(schema, 'properties.metadata.properties.namespace'),
		Ready: lodash.get(schema, 'properties.status.properties.ready') || { 
			description: 'Ready status', 
			type: 'boolean' 
		},
		Replicas: lodash.get(schema, 'properties.spec.properties.deploymentSpec.properties.replicas') || { 
			description: 'Number of replicas', 
			type: 'integer' 
		},
		'Service Type': lodash.get(schema, 'properties.spec.properties.serviceSpec.properties.type') || { 
			description: 'Service type', 
			type: 'string' 
		},
		Age: lodash.get(schema, 'properties.metadata.properties.creationTimestamp'),
		Configuration: schema
	};
}

function simpleappObjectMask(object: any): Record<string, JsonValue> {
	// Parse Ready status from conditions array
	// JSONPath: .status.conditions[?(@.type=="Ready")].status
	let ready: JsonValue = null;
	if (object?.status?.conditions && Array.isArray(object.status.conditions)) {
		const readyCondition = object.status.conditions.find(
			(condition: any) => condition?.type === 'Ready'
		);
		ready = readyCondition?.status === 'True' ? true : (readyCondition?.status === 'False' ? false : null);
	}

	return {
		Name: object?.metadata?.name as JsonValue,
		Namespace: object?.metadata?.namespace as JsonValue,
		Ready: ready,
		Replicas: object?.spec?.deploymentSpec?.replicas as JsonValue,
		'Service Type': object?.spec?.serviceSpec?.type as JsonValue,
		Age: object?.metadata?.creationTimestamp as JsonValue,
		Configuration: object as JsonValue
	};
}

function simpleappColumnDefinitions(apiResource: APIResource, fields: any) {
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
						`/(auth)/${page.params.cluster!}/SimpleApp/simpleapps?group=apps.otterscale.io&version=v1alpha1&name=${row.original[column.id]}&namespace=${page.url.searchParams.get('namespace')!}`
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
			id: 'Ready',
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
			accessorKey: 'Ready'
		},
		{
			id: 'Replicas',
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
			accessorKey: 'Replicas'
		},
		{
			id: 'Service Type',
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
			accessorKey: 'Service Type'
		},
		{
			id: 'Age',
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
			accessorKey: 'Age'
		}
	];
}

export { simpleappColumnDefinitions, simpleappFieldsMask, simpleappObjectMask };
