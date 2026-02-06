import type { Schema, UiSchemaRoot } from '@sjsf/form';

import { deepMerge, deleteByPath, getByPath, setByPath } from './utils';

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

function isMapSchema(schema: K8sOpenAPISchema): boolean {
	return (
		schema.type === 'object' &&
		!!schema.additionalProperties &&
		typeof schema.additionalProperties === 'object'
	);
}

/** Resolves the next key and source schema, handling dotted keys in Maps. */
function resolveNextKey(
	currentSource: K8sOpenAPISchema,
	parts: string[],
	currentIndex: number
): { key: string; nextIndex: number; sourceSchema: K8sOpenAPISchema | undefined } {
	let key = parts[currentIndex];
	let nextIndex = currentIndex;

	if (currentSource.properties?.[key]) {
		return { key, nextIndex, sourceSchema: currentSource.properties[key] };
	}

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

/** Simplifies Quantity schema (oneOf/anyOf: [string, number]) to string type */
function simplifyQuantitySchema(target: Schema, source: K8sOpenAPISchema) {
	const variants = source.oneOf || source.anyOf;
	if (
		Array.isArray(variants) &&
		variants.some((o) => o.type === 'string') &&
		variants.some((o) => o.type === 'number' || o.type === 'integer')
	) {
		target.type = 'string';
		target.title = source.title ?? target.title ?? '';
		delete target.oneOf;
		delete target.anyOf;
	}
}

/** Subsets the full OpenAPI schema to include only the specified paths. */
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

	const pathKeys = Array.isArray(paths) ? paths : Object.keys(paths);
	const pathOptions = Array.isArray(paths) ? {} : paths;

	for (const path of pathKeys) {
		const options: PathOptions = pathOptions[path] || {};
		const customUiOptions = options.uiSchema || {};

		const parts = path.split('.');
		let currentSource: K8sOpenAPISchema = fullSchema;
		let currentTarget: Schema = rootSchema;
		let currentUiTarget: Record<string, unknown> = uiSchema;

		let cumulativePath = '';

		for (let i = 0; i < parts.length; i++) {
			const { key, nextIndex, sourceSchema } = resolveNextKey(currentSource, parts, i);

			i = nextIndex;
			const part = key;
			cumulativePath = cumulativePath ? `${cumulativePath}.${part}` : part;

			if (!sourceSchema) {
				console.warn(`Path segment "${part}" not found in source schema (full path: ${path})`);
				break;
			}

			const sourceProp = sourceSchema;
			const parentSource = currentSource;
			currentSource = sourceProp;

			const isExplicit = pathKeys.includes(cumulativePath);

			if (!currentTarget.properties) {
				currentTarget.properties = {};
			}

			const isLeaf = i === parts.length - 1;

			// Determine if field should be required
			let shouldBeRequired = parentSource.required?.includes(part);
			if (isLeaf && options.required !== undefined) {
				shouldBeRequired = options.required;
			}

			const applyOptions = (target: Schema, src: K8sOpenAPISchema, isLeafNode: boolean) => {
				if (options.title) target.title = options.title;
				else target.title = src.title;

				if (isLeafNode && options.showDescription) target.description = src.description;
				else delete target.description;
			};

			// Hoist Map
			if (isLeaf && isMapSchema(sourceProp)) {
				const formPath = path.replace(/\./g, '_');
				transformationMappings[path] = formPath;

				const valueSchema = (
					typeof sourceProp.additionalProperties === 'object' ? sourceProp.additionalProperties : {}
				) as Schema;

				const newProp: Schema = {
					type: 'array',
					items: {
						type: 'object',
						properties: {
							key: { type: 'string', title: 'Key' },
							value: { ...valueSchema, title: 'Value' }
						}
					}
				};
				applyOptions(newProp, sourceProp, true);

				if (!rootSchema.properties) rootSchema.properties = {};
				rootSchema.properties[formPath] = newProp;

				if (Object.keys(customUiOptions).length > 0) {
					uiSchema[formPath] = deepMerge(uiSchema[formPath] || {}, customUiOptions);
				}

				if (options.required) {
					if (!rootSchema.required) rootSchema.required = [];
					if (!rootSchema.required.includes(formPath)) {
						rootSchema.required.push(formPath);
					}
				}
				break;
			}

			// Apply required for standard fields
			if (shouldBeRequired) {
				if (!currentTarget.required) currentTarget.required = [];
				if (!currentTarget.required.includes(part)) {
					currentTarget.required.push(part);
				}
			}

			if (!currentUiTarget[part]) {
				currentUiTarget[part] = {};
			}
			const partUi = currentUiTarget[part] as Record<string, unknown>;

			if (!isExplicit) {
				partUi['ui:options'] = { label: false };
			}

			if (isLeaf) {
				if (Array.isArray(sourceProp.enum)) {
					partUi['ui:components'] = { stringField: 'enumField' };
				}

				if (Object.keys(customUiOptions).length > 0) {
					currentUiTarget[part] = deepMerge(partUi, customUiOptions);
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
			}

			const isTerminalLeaf =
				isLeaf && !pathKeys.some((k) => k !== cumulativePath && k.startsWith(cumulativePath + '.'));

			if (!currentTarget.properties[part]) {
				if (isTerminalLeaf) {
					const newProp = { ...sourceProp } as Schema;
					if (newProp.type === 'object' && !newProp.properties) {
						newProp.properties = {};
					}
					applyOptions(newProp, sourceProp, true);
					simplifyQuantitySchema(newProp, sourceProp);

					if (options.disabled) {
						newProp.readOnly = true;
					}

					currentTarget.properties[part] = newProp;
				} else {
					if (sourceProp.type === 'array') {
						currentTarget.properties[part] = {
							type: 'array',
							items: { type: 'object', properties: {} }
						};
					} else {
						currentTarget.properties[part] = {
							type: 'object',
							properties: {},
							additionalProperties: false
						};
					}
					applyOptions(currentTarget.properties[part] as Schema, sourceProp, false);
					if (!isExplicit) (currentTarget.properties[part] as Schema).title = '';
				}
			} else if (isLeaf) {
				const target = currentTarget.properties[part] as Schema;
				applyOptions(target, sourceProp, true);
				simplifyQuantitySchema(target, sourceProp);
			}

			// Handle Array Traversal
			if (sourceProp.type === 'array' && sourceProp.items && !isLeaf) {
				currentSource = sourceProp.items;

				const justAdded = currentTarget.properties[part] as Schema;
				if (!justAdded.items || Array.isArray(justAdded.items)) {
					justAdded.items = { type: 'object', properties: {} };
				}
				currentTarget = justAdded.items as Schema;

				if (!partUi.items) {
					partUi.items = {};
				}
				currentUiTarget = partUi.items as Record<string, unknown>;

				continue;
			}

			if (!isLeaf) {
				const nextTarget = currentTarget.properties[part];
				if (typeof nextTarget === 'object' && nextTarget !== null) {
					currentTarget = nextTarget as Schema;
				}
				currentUiTarget = partUi as Record<string, unknown>;
			}
		}
	}

	return { schema: rootSchema, uiSchema, transformationMappings };
}

/** Converts K8s Object to Form data */
export function k8sToFormData(
	data: unknown,
	mappings: Record<string, string>
): Record<string, unknown> {
	if (!data || typeof data !== 'object' || Array.isArray(data)) return {};
	const formData = JSON.parse(JSON.stringify(data));

	for (const [k8sPath, formPath] of Object.entries(mappings)) {
		const originalValue = getByPath(formData, k8sPath);

		if (originalValue && typeof originalValue === 'object' && !Array.isArray(originalValue)) {
			const arrayValue = Object.entries(originalValue).map(([key, value]) => ({
				key,
				value
			}));
			formData[formPath] = arrayValue;
			deleteByPath(formData, k8sPath);
		} else if (originalValue === undefined || originalValue === null) {
			formData[formPath] = [];
			deleteByPath(formData, k8sPath);
		}
	}
	return formData;
}

/** Recursively converts object with numeric keys back to arrays */
export function normalizeArrays(obj: unknown): unknown {
	if (obj === null || obj === undefined) return obj;
	if (Array.isArray(obj)) {
		return obj.map(normalizeArrays);
	}
	if (typeof obj !== 'object') return obj;

	const record = obj as Record<string, unknown>;
	const keys = Object.keys(record);

	// Check if all keys are numeric (indicating it should be an array)
	const allNumeric = keys.length > 0 && keys.every((k) => /^\d+$/.test(k));

	if (allNumeric) {
		const numericKeys = keys.map(Number).sort((a, b) => a - b);

		// Safety check: Prevent sparse arrays with excessively large indices
		if (numericKeys.length > 0 && numericKeys[numericKeys.length - 1] > 10000) {
			const result: Record<string, unknown> = {};
			for (const key of keys) {
				result[key] = normalizeArrays(record[key]);
			}
			return result;
		}

		const arr: unknown[] = [];
		numericKeys.forEach((idx) => {
			arr[idx] = normalizeArrays(record[String(idx)]);
		});
		return arr;
	}

	// Recursively process object properties
	const result: Record<string, unknown> = {};
	for (const key of keys) {
		result[key] = normalizeArrays(record[key]);
	}
	return result;
}

/** Converts Form data back to K8s Object */
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

/** Filters data to only include properties defined in the schema */
export function filterDataBySchema(data: unknown, schema: Schema): Record<string, unknown> {
	if (!data || typeof data !== 'object' || Array.isArray(data)) return {};
	if (!schema || schema.type !== 'object' || !schema.properties) return {};

	const result: Record<string, unknown> = {};
	const record = data as Record<string, unknown>;

	for (const [key, propSchema] of Object.entries(schema.properties)) {
		if (!(key in record)) continue;

		const value = record[key];
		const prop = propSchema as Schema;

		if (
			prop.type === 'object' &&
			prop.properties &&
			value &&
			typeof value === 'object' &&
			!Array.isArray(value)
		) {
			// Recursively filter nested objects
			const filtered = filterDataBySchema(value, prop);
			if (Object.keys(filtered).length > 0) {
				result[key] = filtered;
			}
		} else if (prop.type === 'array' && Array.isArray(value)) {
			// For arrays, filter each item if items schema is an object
			if (prop.items && typeof prop.items === 'object' && !Array.isArray(prop.items)) {
				result[key] = value.map((item) => {
					if (typeof item === 'object' && item !== null) {
						return filterDataBySchema(item, prop.items as Schema);
					}
					return item;
				});
			} else {
				result[key] = value;
			}
		} else {
			// Primitive values or other types
			result[key] = value;
		}
	}

	return result;
}
