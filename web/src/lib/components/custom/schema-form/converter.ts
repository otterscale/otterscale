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

/** Shared build context threaded through all helper functions. */
interface BuildContext {
	readonly rootSchema: Schema;
	readonly uiSchema: UiSchemaRoot;
	readonly transformationMappings: Record<string, string>;
	readonly allPaths: string[];
	readonly allOptions: Record<string, PathOptions>;
}

/** Mutable cursors for traversing source & target schemas in parallel. */
interface TraversalState {
	/** Cursor into the K8s source schema tree */
	source: K8sOpenAPISchema;
	/** Cursor into the JSON Schema being built */
	target: Schema;
	/** Cursor into the UI Schema being built */
	uiTarget: Record<string, unknown>;
	/** Dot-delimited path accumulated so far (e.g. "spec.containers") */
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
 * Resolves the next key and source schema, handling dotted keys in Maps.
 *
 * K8s Map keys (e.g. annotation keys like "app.kubernetes.io/name") may
 * contain dots. When the current source is a Map schema, this function
 * greedily merges consecutive path parts into a single key until it hits
 * a recognized property or another Map boundary.
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

	// Case 3: Not found in source schema
	return { key, nextIndex, sourceSchema: undefined };
}

// ── Schema Helpers ─────────────────────────────────────────

/**
 * Simplifies K8s "Quantity" schemas (oneOf/anyOf: [string, number/integer])
 * to a single string type for form rendering (handles values like "100m", "2Gi").
 */
function simplifyQuantitySchema(target: Schema, source: K8sOpenAPISchema): void {
	const variants = source.oneOf || source.anyOf;
	if (
		Array.isArray(variants) &&
		variants.some((v) => v.type === 'string') &&
		variants.some((v) => v.type === 'number' || v.type === 'integer')
	) {
		target.type = 'string';
		target.title = source.title ?? target.title ?? '';
		delete target.oneOf;
		delete target.anyOf;
	}
}

/** Applies PathOptions (title, description) onto a target schema node. */
function applySchemaOptions(
	target: Schema,
	source: K8sOpenAPISchema,
	options: PathOptions,
	isLeaf: boolean
): void {
	if (options.title) target.title = options.title;
	else target.title = source.title;

	if (isLeaf && options.showDescription) target.description = source.description;
	else delete target.description;
}

// ── Map Hoisting ───────────────────────────────────────────

/**
 * Hoists a K8s Map field (additionalProperties) to the root-level schema
 * as an array of `{key, value}` pairs, making it editable via form controls.
 *
 * Records the path mapping in `transformationMappings` so that
 * `k8sToFormData` and `formDataToK8s` can convert between representations.
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

	if (!ctx.rootSchema.properties) ctx.rootSchema.properties = {};
	ctx.rootSchema.properties[formPath] = arrayProp;

	if (Object.keys(customUiOptions).length > 0) {
		ctx.uiSchema[formPath] = deepMerge(ctx.uiSchema[formPath] || {}, customUiOptions);
	}

	if (options.required) {
		if (!ctx.rootSchema.required) ctx.rootSchema.required = [];
		if (!ctx.rootSchema.required.includes(formPath)) {
			ctx.rootSchema.required.push(formPath);
		}
	}
}

// ── Required Handling ──────────────────────────────────────

/** Adds `part` to the target's required array when appropriate. */
function applyRequiredIfNeeded(
	target: Schema,
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
		if (!target.required) target.required = [];
		if (!target.required.includes(part)) {
			target.required.push(part);
		}
	}
}

// ── UI Schema Building ─────────────────────────────────────

/**
 * Builds the UI schema node for a single path segment.
 * Returns the `partUi` reference used for deeper traversal.
 */
function buildUiSchemaNode(
	uiTarget: Record<string, unknown>,
	part: string,
	sourceProp: K8sOpenAPISchema,
	options: PathOptions,
	customUiOptions: Record<string, unknown>,
	flags: { isLeaf: boolean; isExplicit: boolean }
): Record<string, unknown> {
	if (!uiTarget[part]) {
		uiTarget[part] = {};
	}
	const partUi = uiTarget[part] as Record<string, unknown>;

	// Hide labels for intermediate (non-user-specified) path segments
	if (!flags.isExplicit) {
		partUi['ui:options'] = { label: false };
	}

	if (!flags.isLeaf) return partUi;

	// ── Leaf-only UI configuration ──

	if (Array.isArray(sourceProp.enum)) {
		partUi['ui:components'] = { stringField: 'enumField' };
	}

	if (Object.keys(customUiOptions).length > 0) {
		uiTarget[part] = deepMerge(partUi, customUiOptions);
	}

	if (options.title) {
		partUi['ui:title'] = options.title;
		partUi['ui:options'] = deepMerge(
			(partUi['ui:options'] as Record<string, unknown>) || {},
			{ label: true }
		);
	}

	if (options.disabled) {
		partUi['ui:options'] = deepMerge(
			(partUi['ui:options'] as Record<string, unknown>) || {},
			{ shadcn4Text: { disabled: true } }
		);
	}

	return partUi;
}

// ── Schema Property Creation ───────────────────────────────

/** Creates a full schema property for a terminal leaf (complete copy from source). */
function createTerminalLeafProperty(
	sourceProp: K8sOpenAPISchema,
	options: PathOptions
): Schema {
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

/** Creates a skeleton schema property for an intermediate (non-terminal) node. */
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
 * Advances all three traversal cursors through an array's `items` level.
 *
 * K8s paths treat arrays transparently (`spec.containers.image`
 * rather than `spec.containers[0].image`), so we "dive into" the
 * array's items schema for the next iteration.
 */
function advanceThroughArray(
	state: TraversalState,
	part: string,
	partUi: Record<string, unknown>
): void {
	state.source = state.source.items!;

	const arrayProp = state.target.properties![part] as Schema;
	if (!arrayProp.items || Array.isArray(arrayProp.items)) {
		arrayProp.items = { type: 'object', properties: {} };
	}
	state.target = arrayProp.items as Schema;

	if (!partUi.items) {
		partUi.items = {};
	}
	state.uiTarget = partUi.items as Record<string, unknown>;
}

// ── Main: buildSchemaFromK8s ───────────────────────────────

/**
 * Subsets the full K8s OpenAPI schema to include only the specified paths,
 * producing:
 * - `schema`:  JSON Schema Draft-07 compatible schema for form rendering
 * - `uiSchema`: UI configuration for SJSF form controls
 * - `transformationMappings`: Map ↔ Array conversion lookup for data transforms
 */
export function buildSchemaFromK8s(
	fullSchema: K8sOpenAPISchema,
	paths: string[] | Record<string, PathOptions>
): SchemaFormConfig {
	const rootSchema: Schema = {
		type: 'object',
		properties: {},
		required: [],
		title: fullSchema.title ?? ''
	};

	const uiSchema: UiSchemaRoot = {};
	const transformationMappings: Record<string, string> = {};

	const allPaths = Array.isArray(paths) ? paths : Object.keys(paths);
	const allOptions = Array.isArray(paths) ? {} : paths;
	const ctx: BuildContext = { rootSchema, uiSchema, transformationMappings, allPaths, allOptions };

	for (const path of allPaths) {
		processPath(ctx, fullSchema, path);
	}

	return { schema: rootSchema, uiSchema, transformationMappings };
}

/**
 * Processes a single dot-delimited path, walking the source K8s schema
 * and building up the target JSON Schema, UI Schema, and transformation
 * mappings for that path.
 */
function processPath(
	ctx: BuildContext,
	fullSchema: K8sOpenAPISchema,
	path: string
): void {
	const options: PathOptions = ctx.allOptions[path] || {};
	const customUiOptions = options.uiSchema || {};
	const parts = path.split('.');

	const state: TraversalState = {
		source: fullSchema,
		target: ctx.rootSchema,
		uiTarget: ctx.uiSchema,
		cumulativePath: ''
	};

	for (let i = 0; i < parts.length; i++) {
		// 1. Resolve key (handles dotted Map keys like annotation keys)
		const { key: part, nextIndex, sourceSchema } = resolveNextKey(state.source, parts, i);
		i = nextIndex;

		// 2. Reject prototype-pollution keys
		if (UNSAFE_KEYS.has(part)) break;

		// 3. Track cumulative path
		state.cumulativePath = state.cumulativePath ? `${state.cumulativePath}.${part}` : part;

		// 4. Bail if path segment not found in source
		if (!sourceSchema) {
			console.warn(`Path segment "${part}" not found in source schema (full path: ${path})`);
			break;
		}

		const parentSource = state.source;
		state.source = sourceSchema;
		const isLeaf = i === parts.length - 1;
		const isExplicit = ctx.allPaths.includes(state.cumulativePath);

		// 5. Hoist Map fields → root-level array of {key, value}
		if (isLeaf && isMapSchema(sourceSchema)) {
			hoistMapToRoot(ctx, path, sourceSchema, options, customUiOptions);
			break;
		}

		// 6. Ensure properties object exists & handle required
		if (!state.target.properties) state.target.properties = {};
		applyRequiredIfNeeded(state.target, part, parentSource, options, isLeaf);

		// 7. Build UI schema node
		const partUi = buildUiSchemaNode(
			state.uiTarget, part, sourceSchema,
			options, customUiOptions,
			{ isLeaf, isExplicit }
		);

		// 8. Create or update the schema property
		const isTerminalLeaf =
			isLeaf &&
			!ctx.allPaths.some(
				(k) => k !== state.cumulativePath && k.startsWith(state.cumulativePath + '.')
			);

		if (!state.target.properties[part]) {
			state.target.properties[part] = isTerminalLeaf
				? createTerminalLeafProperty(sourceSchema, options)
				: createIntermediateProperty(sourceSchema, part, options, isExplicit);
		} else if (isLeaf) {
			const existing = state.target.properties[part] as Schema;
			applySchemaOptions(existing, sourceSchema, options, true);
			simplifyQuantitySchema(existing, sourceSchema);
		}

		// 9. Advance cursors: dive into array items or descend normally
		if (sourceSchema.type === 'array' && sourceSchema.items && !isLeaf) {
			advanceThroughArray(state, part, partUi);
			continue;
		}

		if (!isLeaf) {
			const nextTarget = state.target.properties[part];
			if (typeof nextTarget === 'object' && nextTarget !== null) {
				state.target = nextTarget as Schema;
			}
			state.uiTarget = partUi;
		}
	}
}

// ── Data Conversion: K8s ↔ Form ───────────────────────────

/** Converts K8s Object → Form data (Maps become key-value arrays). */
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

/** Converts Form data → K8s Object (key-value arrays become Maps). */
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

	// Normalize numeric-keyed objects back to arrays
	k8sData = normalizeArrays(k8sData) as Record<string, unknown>;

	return k8sData;
}

// ── Array Normalization ────────────────────────────────────

/** Recursively converts objects with all-numeric keys back to arrays. */
export function normalizeArrays(obj: unknown): unknown {
	if (obj === null || obj === undefined) return obj;
	if (Array.isArray(obj)) return obj.map(normalizeArrays);
	if (typeof obj !== 'object') return obj;

	const record = obj as Record<string, unknown>;
	const keys = Object.keys(record);
	const allNumeric = keys.length > 0 && keys.every((k) => /^\d+$/.test(k));

	if (allNumeric) {
		const numericKeys = keys.map(Number).sort((a, b) => a - b);

		// Safety check: Prevent sparse arrays with excessively large indices
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

/** Recursively normalizes all values in a plain object. */
function recurseObjectValues(record: Record<string, unknown>): Record<string, unknown> {
	const result: Record<string, unknown> = {};
	for (const key of Object.keys(record)) {
		result[key] = normalizeArrays(record[key]);
	}
	return result;
}

// ── Data Filtering ─────────────────────────────────────────

/** Filters data to only include properties defined in the schema. */
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

/** Checks if a value is a non-array object whose schema has sub-properties to recurse into. */
function isFilterableObject(prop: Schema, value: unknown): value is Record<string, unknown> {
	return (
		prop.type === 'object' &&
		!!prop.properties &&
		!!value &&
		typeof value === 'object' &&
		!Array.isArray(value)
	);
}

/** Filters each array element against the items schema when applicable. */
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
