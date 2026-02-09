import { getCronJobColumnDefinitions, getCronJobFields, getCronJobValues } from './cronjob.js';
import { getDefaultColumnDefinitions, getDefaultFields, getDefaultValues } from './default.js';
import {
	getResourceQuotaColumnDefinitions, getResourceQuotaFields, getResourceQuotaValues
} from './resource-quota.js';
import type { FieldsType, ValuesType } from '../type.js';
import type { Column, ColumnDef } from '@tanstack/table-core';
import type { APIResource } from '$lib/api/resource/v1/resource_pb.js';

function getFieldsDynamically(kind: string, schema: any): FieldsType {
	switch (kind) {
		case 'CronJob':
			return getCronJobFields(schema);
		case 'ResourceQuota':
			return getResourceQuotaFields(schema);
		default:
			return getDefaultFields(schema);
	}
}

function getValuesDynamically(kind: string, object: any): ValuesType {
	switch (kind) {
		case 'CronJob':
			return getCronJobValues(object);
		case 'ResourceQuota':
			return getResourceQuotaValues(object);
		default:
			return getDefaultValues(object);
	}
}

function getColumnDefinitionsDynamically(kind: string, apiResource: APIResource, fields: FieldsType): ColumnDef<ValuesType>[] {
	switch (kind) {
		case 'CronJob':
			return getCronJobColumnDefinitions(apiResource, fields);
		case 'ResourceQuota':
			return getResourceQuotaColumnDefinitions(apiResource, fields);
		default:
			return getDefaultColumnDefinitions(apiResource, fields);
	}
}

export { getFieldsDynamically, getValuesDynamically, getColumnDefinitionsDynamically };
