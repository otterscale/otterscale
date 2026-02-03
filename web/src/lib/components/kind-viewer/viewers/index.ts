// import type { Component } from 'svelte';

import type { JsonValue } from '@bufbuild/protobuf';
import type { ColumnDef } from '@tanstack/table-core';

import type { APIResource } from '$lib/api/resource/v1/resource_pb.js';

type ColumnDefinitionsType = (
	apiResource: APIResource,
	fields: any
) => ColumnDef<Record<string, JsonValue>>[];
type FieldsMaskType = (
	schema: any
) => Record<string, { description: string; type: string; format: string }>;
type ObjectMaskType = (object: any) => Record<string, JsonValue>;

import { cronjobColumnDefinitions, cronjobFieldsMask, cronjobObjectMask } from './cronjob.js';
import { defaultColumnDefinitions, defaultFieldsMask, defaultObjectMask } from './default.js';

function getFieldsGetter(kind: string): FieldsMaskType {
	if (kind === 'CronJob') {
		return cronjobFieldsMask as FieldsMaskType;
	}
	return defaultFieldsMask as FieldsMaskType;
}

function getObjectGetter(kind: string): ObjectMaskType {
	if (kind === 'CronJob') {
		return cronjobObjectMask as ObjectMaskType;
	}
	return defaultObjectMask as ObjectMaskType;
}

function getColumnDefinitionsGetter(kind: string): ColumnDefinitionsType {
	if (kind === 'CronJob') {
		return cronjobColumnDefinitions as unknown as ColumnDefinitionsType;
	}
	return defaultColumnDefinitions as unknown as ColumnDefinitionsType;
}

export { getColumnDefinitionsGetter, getFieldsGetter, getObjectGetter };
