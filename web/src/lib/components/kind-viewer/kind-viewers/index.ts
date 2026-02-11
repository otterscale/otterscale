import type { JsonValue } from '@bufbuild/protobuf';
import type { ColumnDef } from '@tanstack/table-core';

import type { APIResource } from '$lib/api/resource/v1/resource_pb.js';
import type { DataSchemaType, UISchemaType } from '$lib/components/dynamic-table/utils.js';

import {
	getCronJobColumnDefinitions,
	getCronJobData,
	getCronJobDataSchemas,
	getCronJobUISchemas
} from './cronjob.js';
import {
	getDefaultColumnDefinitions,
	getDefaultData,
	getDefaultDataSchemas,
	getDefaultUISchemas
} from './default.js';
import {
	getResourceQuotaColumnDefinitions,
	getResourceQuotaData,
	getResourceQuotaDataSchemas,
	getResourceQuotaUISchemas
} from './resource-quota.js';
import {
	getWorkspaceColumnDefinitions,
	getWorkspaceData,
	getWorkspaceDataSchemas,
	getWorkspaceUISchemas
} from './workspaces.js';

function getDataSchemas(kind: string): Record<string, DataSchemaType> {
	switch (kind) {
		case 'CronJob':
			return getCronJobDataSchemas();
		case 'ResourceQuota':
			return getResourceQuotaDataSchemas();
		case 'Workspace':
			return getWorkspaceDataSchemas();
		default:
			return getDefaultDataSchemas();
	}
}

function getUISchemas(kind: string): Record<string, UISchemaType> {
	switch (kind) {
		case 'CronJob':
			return getCronJobUISchemas();
		case 'ResourceQuota':
			return getResourceQuotaUISchemas();
		case 'Workspace':
			return getWorkspaceUISchemas();
		default:
			return getDefaultUISchemas();
	}
}

function getDataSet(kind: string, object: any): Record<string, JsonValue> {
	switch (kind) {
		case 'CronJob':
			return getCronJobData(object);
		case 'ResourceQuota':
			return getResourceQuotaData(object);
		case 'Workspace':
			return getWorkspaceData(object);
		default:
			return getDefaultData(object);
	}
}

function getColumnDefinitions(
	kind: string,
	apiResource: APIResource,
	uiSchemas: Record<string, UISchemaType>,
	dataSchemas: Record<string, DataSchemaType>
): ColumnDef<Record<string, JsonValue>>[] {
	switch (kind) {
		case 'CronJob':
			return getCronJobColumnDefinitions(apiResource, uiSchemas, dataSchemas);
		case 'ResourceQuota':
			return getResourceQuotaColumnDefinitions(apiResource, uiSchemas, dataSchemas);
		case 'Workspace':
			return getWorkspaceColumnDefinitions(apiResource, uiSchemas, dataSchemas);
		default:
			return getDefaultColumnDefinitions(apiResource, uiSchemas, dataSchemas);
	}
}

export { getColumnDefinitions, getDataSchemas, getDataSet, getUISchemas };
