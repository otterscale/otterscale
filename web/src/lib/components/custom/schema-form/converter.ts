import type { Schema, UiSchemaRoot } from '@sjsf/form';
import lodash from 'lodash';
import pluralize from 'pluralize';

import { deepMerge, deleteByPath, getByPath, setByPath } from './utils';

// ── Public Types ───────────────────────────────────────────

/** K8s OpenAPI Schema type */
export interface K8sOpenAPISchema {
	type?: string;
	title?: string;
	description?: string;
	properties?: Record<string, K8sOpenAPISchema>;
	items?: K8sOpenAPISchema;
	required?: string[];
	default?: unknown;
	enum?: string[];
	format?: string;
	nullable?: boolean;
	additionalProperties?: K8sOpenAPISchema | boolean;
	oneOf?: K8sOpenAPISchema[];
	anyOf?: K8sOpenAPISchema[];
	[key: string]: unknown;
}

export interface SchemaFormConfig {
	schema: Schema;
	uiSchema: UiSchemaRoot;
	transformationMappings: Record<string, string>;
}

export interface PathOptions {
	title?: string;
	required?: boolean;
	showDescription?: boolean;
	uiSchema?: Record<string, unknown>;
	disabled?: boolean;
}

// ── Internal Types ─────────────────────────────────────────

interface BuildContext {
	readonly schema: Schema;
	readonly uiSchema: UiSchemaRoot & Record<string, unknown>;
	readonly transformationMappings: Record<string, string>;
	readonly allPaths: string[];
	readonly allOptions: Record<string, PathOptions>;
}

interface TraversalState {
	source: K8sOpenAPISchema;
	schema: Schema;
	uiSchema: Record<string, unknown>;
	cumulativePath: string;
}

// ── Predicates ─────────────────────────────────────────────

const UNSAFE_KEYS = new Set(['__proto__', 'constructor', 'prototype']);

function isMapSchema(schema: K8sOpenAPISchema): boolean {
	return (
		schema.type === 'object' &&
		!!schema.additionalProperties &&
		typeof schema.additionalProperties === 'object'
	);
}

// ── Key Resolution ─────────────────────────────────────────

/**
 * Resolves the next key and its source schema, merging dotted Map keys.
 *
 * @example path: "spec.resourceQuota.hard.requests.cpu"
 * // parts=["spec","resourceQuota","hard","requests","cpu"]
 * //
 * // index=0 → "spec"          (direct property match)
 * // index=1 → "resourceQuota" (direct property match)
 * // index=2 → "hard"          (direct property match)
 * // index=3 → "hard" is a Map (additionalProperties: { anyOf:[integer,string] })
 * //           "requests" not in value schema props, "cpu" also not
 * //           → merged key="requests.cpu", nextIndex=4
 */
function resolveNextKey(
	currentSource: K8sOpenAPISchema,
	parts: string[],
	currentIndex: number
): { key: string; nextIndex: number; sourceSchema: K8sOpenAPISchema | undefined } {
	let key = parts[currentIndex];
	let nextIndex = currentIndex;

	// Case 1: Direct property match
	if (currentSource.properties?.[key]) {
		return { key, nextIndex, sourceSchema: currentSource.properties[key] };
	}

	// Case 2: Map schema — merge dotted segments back into a single key
	if (isMapSchema(currentSource)) {
		const mapValueSchema = currentSource.additionalProperties as K8sOpenAPISchema;

		let peek = currentIndex + 1;
		while (peek < parts.length) {
			const nextPart = parts[peek];
			const isValueMap = isMapSchema(mapValueSchema);
			const nextPartIsProp = !!mapValueSchema.properties?.[nextPart];

			if (!nextPartIsProp && !isValueMap) {
				key += `.${nextPart}`;
				peek++;
			} else {
				break;
			}
		}
		nextIndex = peek - 1;
		return { key, nextIndex, sourceSchema: mapValueSchema };
	}

	return { key, nextIndex, sourceSchema: undefined };
}

// ── Schema Helpers ─────────────────────────────────────────

/**
 * { oneOf: [{type:"string"},{type:"integer"}] }  →  { type: "string" }
 * Lets the form use a single text input for K8s Quantity values ("100m", "2Gi").
 */
function simplifyQuantitySchema(schema: Schema, source: K8sOpenAPISchema): void {
	const variants = source.oneOf || source.anyOf;
	if (
		Array.isArray(variants) &&
		variants.some((v) => v.type === 'string') &&
		variants.some((v) => v.type === 'number' || v.type === 'integer')
	) {
		schema.type = 'string';
		schema.title = source.title ?? schema.title ?? '';
		delete schema.oneOf;
		delete schema.anyOf;
	}
}

/** Applies user-provided title / description overrides onto a schema node. */
function applySchemaOptions(
	schema: Schema,
	source: K8sOpenAPISchema,
	options: PathOptions,
	isLeaf: boolean
): void {
	if (options.title) schema.title = options.title;
	else schema.title = source.title;

	if (isLeaf && options.showDescription) schema.description = source.description;
	else delete schema.description;
}

// ── Map Hoisting ───────────────────────────────────────────

/**
 * Hoists a K8s Map (additionalProperties) to root as an editable array.
 *
 * @example
 * // path: "metadata.labels"  →  formPath: "metadata_labels"
 * // K8s:  { metadata: { labels: { app: "web", env: "prod" } } }
 * // Form: { metadata_labels: [{ key: "app", value: "web" }, { key: "env", value: "prod" }] }
 */
function hoistMapToRoot(
	ctx: BuildContext,
	originalPath: string,
	sourceProp: K8sOpenAPISchema,
	options: PathOptions,
	customUiOptions: Record<string, unknown>
): void {
	const formPath = originalPath.replace(/\./g, '_');
	ctx.transformationMappings[originalPath] = formPath;

	const valueSchema = (
		typeof sourceProp.additionalProperties === 'object' ? sourceProp.additionalProperties : {}
	) as Schema;

	const arrayProp: Schema = {
		type: 'array',
		items: {
			type: 'object',
			properties: {
				key: { type: 'string', title: 'Key' },
				value: { ...valueSchema, title: 'Value' }
			}
		}
	};
	applySchemaOptions(arrayProp, sourceProp, options, true);

	if (!ctx.schema.properties) ctx.schema.properties = {};
	ctx.schema.properties[formPath] = arrayProp;

	if (Object.keys(customUiOptions).length > 0) {
		ctx.uiSchema[formPath] = deepMerge(ctx.uiSchema[formPath] || {}, customUiOptions);
	}

	if (options.required) {
		if (!ctx.schema.required) ctx.schema.required = [];
		if (!ctx.schema.required.includes(formPath)) {
			ctx.schema.required.push(formPath);
		}
	}
}

// ── Required Handling ──────────────────────────────────────

/** Marks field as required (PathOptions overrides K8s source at leaf level). */
function applyRequiredIfNeeded(
	schema: Schema,
	part: string,
	parentSource: K8sOpenAPISchema,
	options: PathOptions,
	isLeaf: boolean
): void {
	let shouldBeRequired = parentSource.required?.includes(part);
	if (isLeaf && options.required !== undefined) {
		shouldBeRequired = options.required;
	}

	if (shouldBeRequired) {
		if (!schema.required) schema.required = [];
		if (!schema.required.includes(part)) {
			schema.required.push(part);
		}
	}
}

// ── UI Schema Building ─────────────────────────────────────

/**
 * Builds UI schema for one path segment. Returns the partUi cursor.
 *
 * @example
 * // Intermediate: { "ui:options": { label: false } }
 * // Leaf with enum: { "ui:components": { stringField: "enumField" } }
 * // Leaf with title: { "ui:title": "Name", "ui:options": { label: true } }
 */
function buildUiSchemaNode(
	uiSchema: Record<string, unknown>,
	part: string,
	sourceProp: K8sOpenAPISchema,
	options: PathOptions,
	customUiOptions: Record<string, unknown>,
	flags: { isLeaf: boolean; isExplicit: boolean }
): Record<string, unknown> {
	if (!uiSchema[part]) {
		uiSchema[part] = {};
	}
	const partUi = uiSchema[part] as Record<string, unknown>;

	if (!flags.isExplicit) {
		partUi['ui:options'] = { label: false };
	}

	if (!flags.isLeaf) return partUi;

	if (Array.isArray(sourceProp.enum)) {
		partUi['ui:components'] = { stringField: 'enumField' };
	}

	if (Object.keys(customUiOptions).length > 0) {
		uiSchema[part] = deepMerge(partUi, customUiOptions);
	}

	if (options.title) {
		partUi['ui:title'] = options.title;
		partUi['ui:options'] = deepMerge((partUi['ui:options'] as Record<string, unknown>) || {}, {
			label: true
		});
	}

	if (options.disabled) {
		partUi['ui:options'] = deepMerge((partUi['ui:options'] as Record<string, unknown>) || {}, {
			shadcn4Text: { disabled: true }
		});
	}

	return partUi;
}

// ── Schema Property Creation ───────────────────────────────

/** Full copy of source schema for a terminal leaf (no deeper paths reference it). */
function createTerminalLeafProperty(sourceProp: K8sOpenAPISchema, options: PathOptions): Schema {
	const prop = { ...sourceProp } as Schema;
	if (prop.type === 'object' && !prop.properties) {
		prop.properties = {};
	}

	applySchemaOptions(prop, sourceProp, options, true);
	simplifyQuantitySchema(prop, sourceProp);

	if (options.disabled) {
		prop.readOnly = true;
	}

	return prop;
}

/** Skeleton schema (type: object/array) for nodes that have deeper children. */
function createIntermediateProperty(
	sourceProp: K8sOpenAPISchema,
	part: string,
	options: PathOptions,
	isExplicit: boolean
): Schema {
	let prop: Schema;

	if (sourceProp.type === 'array') {
		prop = {
			type: 'array',
			items: {
				type: 'object',
				title: formatName(part),
				properties: {}
			}
		};
	} else {
		prop = {
			type: 'object',
			properties: {},
			additionalProperties: false
		};
	}

	applySchemaOptions(prop, sourceProp, options, false);
	if (!isExplicit) prop.title = '';

	return prop;
}

// ── Array Traversal ────────────────────────────────────────

/**
 * Dives all cursors into array `items` (K8s uses "spec.containers.image" not "[0].image").
 */
function advanceThroughArray(
	state: TraversalState,
	part: string,
	partUi: Record<string, unknown>
): void {
	state.source = state.source.items!;

	const arrayProp = state.schema.properties![part] as Schema;
	if (!arrayProp.items || Array.isArray(arrayProp.items)) {
		arrayProp.items = { type: 'object', properties: {} };
	}
	state.schema = arrayProp.items as Schema;

	if (!partUi.items) {
		partUi.items = {};
	}
	state.uiSchema = partUi.items as Record<string, unknown>;
}

// ── Main: buildSchemaFromK8s ───────────────────────────────

/**
 * Extracts a subset of a K8s OpenAPI schema → SJSF-compatible { schema, uiSchema, transformationMappings }.
 *
 * @example
 * const { schema, uiSchema, transformationMappings } = buildSchemaFromK8s(k8sSchema, {
 *   'metadata.name':   { title: 'Name' },
 *   'metadata.labels':  { title: 'Labels' },          // Map → hoisted as array of {key,value}
 *   'spec.replicas':    { title: 'Replicas' },
 *   'spec.containers.image': { title: 'Image' },       // Array traversal (containers is array)
 * });
 * // schema.properties   → { metadata: { properties: { name: ... } }, metadata_labels: [...], spec: ... }
 * // transformationMappings → { "metadata.labels": "metadata_labels" }
 */
export function buildSchemaFromK8s(
	fullSchema: K8sOpenAPISchema,
	paths: string[] | Record<string, PathOptions>
): SchemaFormConfig {
	const schema: Schema = {
		type: 'object',
		properties: {},
		required: [],
		title: fullSchema.title ?? ''
	};

	const uiSchema = {} as UiSchemaRoot & Record<string, unknown>;
	const transformationMappings: Record<string, string> = {};

	const allPaths = Array.isArray(paths) ? paths : Object.keys(paths);
	const allOptions = Array.isArray(paths) ? {} : paths;
	const ctx: BuildContext = { schema, uiSchema, transformationMappings, allPaths, allOptions };

	for (const path of allPaths) {
		processPath(ctx, fullSchema, path);
	}

	return { schema, uiSchema, transformationMappings };
}

/** Walks one dot-delimited path, building schema + uiSchema + mappings. */
function processPath(ctx: BuildContext, fullSchema: K8sOpenAPISchema, path: string): void {
	const options: PathOptions = ctx.allOptions[path] || {};
	const customUiOptions = options.uiSchema || {};
	const parts = path.split('.');

	const state: TraversalState = {
		source: fullSchema,
		schema: ctx.schema,
		uiSchema: ctx.uiSchema as Record<string, unknown>,
		cumulativePath: ''
	};

	for (let i = 0; i < parts.length; i++) {
		const { key: part, nextIndex, sourceSchema } = resolveNextKey(state.source, parts, i);
		i = nextIndex;

		if (UNSAFE_KEYS.has(part)) break;
		state.cumulativePath = state.cumulativePath ? `${state.cumulativePath}.${part}` : part;

		if (!sourceSchema) {
			console.warn(`Path segment "${part}" not found in source schema (full path: ${path})`);
			break;
		}

		const parentSource = state.source;
		state.source = sourceSchema;
		const isLeaf = i === parts.length - 1;
		const isExplicit = ctx.allPaths.includes(state.cumulativePath);

		if (isLeaf && isMapSchema(sourceSchema)) {
			hoistMapToRoot(ctx, path, sourceSchema, options, customUiOptions);
			break;
		}

		if (!state.schema.properties) state.schema.properties = {};
		applyRequiredIfNeeded(state.schema, part, parentSource, options, isLeaf);

		const partUi = buildUiSchemaNode(state.uiSchema, part, sourceSchema, options, customUiOptions, {
			isLeaf,
			isExplicit
		});

		const isTerminalLeaf =
			isLeaf &&
			!ctx.allPaths.some(
				(k) => k !== state.cumulativePath && k.startsWith(state.cumulativePath + '.')
			);

		if (!state.schema.properties[part]) {
			state.schema.properties[part] = isTerminalLeaf
				? createTerminalLeafProperty(sourceSchema, options)
				: createIntermediateProperty(sourceSchema, part, options, isExplicit);
		} else if (isLeaf) {
			const existing = state.schema.properties[part] as Schema;
			applySchemaOptions(existing, sourceSchema, options, true);
			simplifyQuantitySchema(existing, sourceSchema);
		}

		if (sourceSchema.type === 'array' && sourceSchema.items && !isLeaf) {
			advanceThroughArray(state, part, partUi);
			continue;
		}

		if (!isLeaf) {
			const nextTarget = state.schema.properties[part];
			if (typeof nextTarget === 'object' && nextTarget !== null) {
				state.schema = nextTarget as Schema;
			}
			state.uiSchema = partUi;
		}
	}
}

// ── Data Conversion: K8s ↔ Form ───────────────────────────

/**
 * K8s → Form: converts Map objects to key-value arrays based on mappings.
 *
 * @example
 * // mappings: { "metadata.labels": "metadata_labels" }
 * // input:  { metadata: { labels: { app: "web" } } }
 * // output: { metadata: {}, metadata_labels: [{ key: "app", value: "web" }] }
 */
export function k8sToFormData(
	data: unknown,
	mappings: Record<string, string>
): Record<string, unknown> {
	if (!data || typeof data !== 'object' || Array.isArray(data)) return {};
	const formData = JSON.parse(JSON.stringify(data));

	for (const [k8sPath, formPath] of Object.entries(mappings)) {
		const originalValue = getByPath(formData, k8sPath);

		if (originalValue && typeof originalValue === 'object' && !Array.isArray(originalValue)) {
			formData[formPath] = Object.entries(originalValue).map(([key, value]) => ({
				key,
				value
			}));
			deleteByPath(formData, k8sPath);
		} else if (originalValue === undefined || originalValue === null) {
			formData[formPath] = [];
			deleteByPath(formData, k8sPath);
		}
	}
	return formData;
}

/**
 * Form → K8s: converts key-value arrays back to Map objects.
 *
 * @example
 * // mappings: { "metadata.labels": "metadata_labels" }
 * // input:  { metadata_labels: [{ key: "app", value: "web" }] }
 * // output: { metadata: { labels: { app: "web" } } }
 */
export function formDataToK8s(
	formData: unknown,
	mappings: Record<string, string>
): Record<string, unknown> {
	if (!formData || typeof formData !== 'object' || Array.isArray(formData)) return {};
	let k8sData = JSON.parse(JSON.stringify(formData));

	for (const [k8sPath, formPath] of Object.entries(mappings)) {
		const arrayValue = k8sData[formPath];

		if (Array.isArray(arrayValue)) {
			const objectValue = arrayValue.reduce(
				(acc: Record<string, unknown>, item: { key?: string; value?: unknown }) => {
					if (item && item.key) {
						acc[item.key] = item.value;
					}
					return acc;
				},
				{}
			);
			setByPath(k8sData, k8sPath, objectValue);
			delete k8sData[formPath];
		}
	}

	k8sData = normalizeArrays(k8sData) as Record<string, unknown>;

	return k8sData;
}

// ── Array Normalization ────────────────────────────────────

/**
 * { "0": "a", "1": "b" } → ["a", "b"]  (recursively)
 */
export function normalizeArrays(obj: unknown): unknown {
	if (obj === null || obj === undefined) return obj;
	if (Array.isArray(obj)) return obj.map(normalizeArrays);
	if (typeof obj !== 'object') return obj;

	const record = obj as Record<string, unknown>;
	const keys = Object.keys(record);
	const allNumeric = keys.length > 0 && keys.every((k) => /^\d+$/.test(k));

	if (allNumeric) {
		const numericKeys = keys.map(Number).sort((a, b) => a - b);

		// Prevent sparse arrays with excessively large indices
		if (numericKeys.length > 0 && numericKeys[numericKeys.length - 1] > 10000) {
			return recurseObjectValues(record);
		}

		const arr: unknown[] = [];
		numericKeys.forEach((idx) => {
			arr[idx] = normalizeArrays(record[String(idx)]);
		});
		return arr;
	}

	return recurseObjectValues(record);
}

function recurseObjectValues(record: Record<string, unknown>): Record<string, unknown> {
	const result: Record<string, unknown> = {};
	for (const key of Object.keys(record)) {
		result[key] = normalizeArrays(record[key]);
	}
	return result;
}

// ── Data Filtering ─────────────────────────────────────────

/**
 * Keeps only the data fields present in the schema (deep recursive).
 *
 * @example
 * // schema defines: { name: string, image: string }
 * // input:  { name: "web", image: "nginx", extra: true }
 * // output: { name: "web", image: "nginx" }
 */
export function filterDataBySchema(data: unknown, schema: Schema): Record<string, unknown> {
	if (!data || typeof data !== 'object' || Array.isArray(data)) return {};
	if (!schema || schema.type !== 'object' || !schema.properties) return {};

	const result: Record<string, unknown> = {};
	const record = data as Record<string, unknown>;

	for (const [key, propSchema] of Object.entries(schema.properties)) {
		if (!(key in record)) continue;

		const value = record[key];
		const prop = propSchema as Schema;

		if (isFilterableObject(prop, value)) {
			const filtered = filterDataBySchema(value, prop);
			if (Object.keys(filtered).length > 0) {
				result[key] = filtered;
			}
		} else if (prop.type === 'array' && Array.isArray(value)) {
			result[key] = filterArrayItems(value, prop);
		} else {
			result[key] = value;
		}
	}

	return result;
}

function isFilterableObject(prop: Schema, value: unknown): value is Record<string, unknown> {
	return (
		prop.type === 'object' &&
		!!prop.properties &&
		!!value &&
		typeof value === 'object' &&
		!Array.isArray(value)
	);
}

function filterArrayItems(value: unknown[], prop: Schema): unknown[] {
	if (prop.items && typeof prop.items === 'object' && !Array.isArray(prop.items)) {
		return value.map((item) =>
			typeof item === 'object' && item !== null
				? filterDataBySchema(item, prop.items as Schema)
				: item
		);
	}
	return value;
}

// ── Formatting ─────────────────────────────────────────────

function formatName(text: string): string {
	return lodash.startCase(pluralize.singular(text));
}
