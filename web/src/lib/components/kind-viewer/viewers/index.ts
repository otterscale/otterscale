import type { JsonValue } from '@bufbuild/protobuf';
import type { ColumnDef } from '@tanstack/table-core';

import type { APIResource } from '$lib/api/resource/v1/resource_pb.js';

type ColumnDefinitionsType = (
	apiResource: APIResource,
	fields: any
) => ColumnDef<Record<string, JsonValue>>[];
type FieldsMaskType = (
	schema: any
) => Record<string, { description: string; type: string; format?: string }>;
type ObjectMaskType = (object: any) => Record<string, JsonValue>;

import { cronjobColumnDefinitions, cronjobFieldsMask, cronjobObjectMask } from './cronjob.js';
import { defaultColumnDefinitions, defaultFieldsMask, defaultObjectMask } from './default.js';
import { podColumnDefinitions, podFieldsMask, podObjectMask } from './pod.js';
import {
	resourceQuotaColumnDefinitions,
	resourceQuotaFieldsMask,
	resourceQuotaObjectMask
} from './resource-quota.js';

function getFieldsGetter(kind: string): FieldsMaskType {
	switch (kind) {
		case 'CronJob':
			return cronjobFieldsMask;
		case 'Pod':
			return podFieldsMask;
		case 'ResourceQuota':
			return resourceQuotaFieldsMask;
		default:
			return defaultFieldsMask;
	}
}

function getObjectGetter(kind: string): ObjectMaskType {
	switch (kind) {
		case 'CronJob':
			return cronjobObjectMask;
		case 'Pod':
			return podObjectMask;
		case 'ResourceQuota':
			return resourceQuotaObjectMask;
		default:
			return defaultObjectMask;
	}
}

function getColumnDefinitionsGetter(kind: string): ColumnDefinitionsType {
	switch (kind) {
		case 'CronJob':
			return cronjobColumnDefinitions;
		case 'Pod':
			return podColumnDefinitions;
		case 'ResourceQuota':
			return resourceQuotaColumnDefinitions;
		default:
			return defaultColumnDefinitions;
	}
}

export { getColumnDefinitionsGetter, getFieldsGetter, getObjectGetter };
