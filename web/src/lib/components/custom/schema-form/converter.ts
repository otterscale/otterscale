import type { Schema, UiSchemaRoot } from '@sjsf/form';
import type { SchemaDefinition } from '@sjsf/form/core';

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
	[key: string]: unknown;
}

/** Result of building schema from K8s OpenAPI */
export interface SchemaFormConfig {
	schema: Schema;
	uiSchema: UiSchemaRoot;
	initialValue: Record<string, unknown>;
	mapPaths: string[];
}

export interface PathOptions {
	title?: string;
	showDescription?: boolean;
}

/**
 * Subsets the full OpenAPI schema to include only the specified paths.
 * @param fullSchema The full OpenAPI V3 schema object
 * @param paths Array of dot-notation paths to include (e.g. "metadata.name", "spec.running")
 * @returns An object containing { schema, uiSchema, initialValue, mapPaths }
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
	const mapPaths: string[] = [];

	const pathKeys = Array.isArray(paths) ? paths : Object.keys(paths);
	const pathOptions = Array.isArray(paths) ? {} : paths;

	for (const path of pathKeys) {
		const options = pathOptions[path] || {};
		const parts = path.split('.');
		let currentSource: K8sOpenAPISchema = fullSchema;
		let currentTarget: Schema = rootSchema;
		let currentUiTarget: any = uiSchema;

		for (let i = 0; i < parts.length; i++) {
			const part = parts[i];

			// Handling implicit array traversal:
			// If currentSource is an Array, we expect to find the property inside 'items'
			if (currentSource.type === 'array' && currentSource.items) {
				currentSource = currentSource.items;

				if (!currentTarget.items) {
					currentTarget.items = { type: 'object', properties: {} };
				}
				// Only proceed if items is a Schema object (not boolean or array)
				if (typeof currentTarget.items === 'object' && !Array.isArray(currentTarget.items)) {
					currentTarget = currentTarget.items as Schema;
				}

				if (!currentUiTarget.items) {
					currentUiTarget.items = {};
				}
				currentUiTarget = currentUiTarget.items;
			}

			// Locate in Source
			if (!currentSource.properties?.[part]) {
				console.warn(`Path segment "${part}" not found in source schema (full path: ${path})`);
				break;
			}
			const sourceProp = currentSource.properties[part];
			const isLeaf = i === parts.length - 1;

			// Update UI Schema
			if (!currentUiTarget[part]) {
				currentUiTarget[part] = {};
			}

			if (isLeaf) {
				if (Array.isArray(sourceProp.enum)) {
					currentUiTarget[part]['ui:components'] = { stringField: 'enumField' };
				} else if (
					sourceProp.type === 'array' &&
					sourceProp.items &&
					!Array.isArray(sourceProp.items) &&
					Array.isArray(sourceProp.items.enum)
				) {
					currentUiTarget[part]['ui:components'] = { arrayField: 'multiEnumField' };
				}
			}

			// Ensure Target Properties exists
			if (!currentTarget.properties) {
				currentTarget.properties = {};
			}

			// Function to apply common transformations (title, description)
			const applyOptions = (target: Schema, src: K8sOpenAPISchema, isLeafNode: boolean) => {
				// Title: Use options.title if leaf and provided, else source title
				if (isLeafNode && options.title) {
					target.title = options.title;
				} else {
					target.title = src.title;
				}

				// Description: Default to NONE. Only show if leaf and showDescription is true.
				if (isLeafNode && options.showDescription) {
					target.description = src.description;
				} else {
					delete target.description;
				}
			};

			// Check if this property is a Map (object with additionalProperties)
			const isMap = sourceProp.type === 'object' && !!sourceProp.additionalProperties;

			// Prepare Target Property
			if (!currentTarget.properties[part]) {
				if (isLeaf) {
					let newProp: Schema;

					if (isMap) {
						// Transform Map to Array of Key-Value
						mapPaths.push(path);
						const valueSchema = (typeof sourceProp.additionalProperties === 'object'
							? sourceProp.additionalProperties
							: {}) as Schema;

						newProp = {
							type: 'array',
							items: {
								type: 'object',
								properties: {
									key: { type: 'string', title: 'Key' },
									value: { ...valueSchema, title: 'Value' }
								}
							}
						};
					} else {
						// Leaf: Full Copy
						newProp = { ...sourceProp } as Schema;
						// Ensure properties is initialized for objects, otherwise some form generators won't render the map
						if (newProp.type === 'object' && !newProp.properties) {
							newProp.properties = {};
						}
					}

					applyOptions(newProp, sourceProp, true);
					currentTarget.properties[part] = newProp;
				} else {
					// Intermediate: Skeleton
					if (sourceProp.type === 'array') {
						currentTarget.properties[part] = {
							type: 'array',
							items: { type: 'object', properties: {} }
						};
					} else {
						currentTarget.properties[part] = {
							type: 'object',
							properties: {}
						};
					}
					// Apply intermediate options (mostly just ensuring no description)
					applyOptions(currentTarget.properties[part] as Schema, sourceProp, false);
				}
			} else if (isLeaf) {
				// Exists but we're at leaf - upgrade to full copy (or transformed map)
				if (isMap) {
					mapPaths.push(path);
					const valueSchema = (typeof sourceProp.additionalProperties === 'object'
						? sourceProp.additionalProperties
						: {}) as Schema;
					
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
					currentTarget.properties[part] = newProp;
				} else {
					const newProp = { ...sourceProp } as Schema;
					if (newProp.type === 'object' && !newProp.properties) {
						newProp.properties = {};
					}
					applyOptions(newProp, sourceProp, true);
					currentTarget.properties[part] = newProp;
				}
			}

			// Handle Required (for both leaf and intermediate)
			const isRequired =
				Array.isArray(currentSource.required) && currentSource.required.includes(part);
			if (isRequired) {
				if (!currentTarget.required) currentTarget.required = [];
				if (!currentTarget.required.includes(part)) {
					currentTarget.required.push(part);
				}
			}

			// Advance
			currentSource = sourceProp;
			const targetProp = currentTarget.properties[part];
			if (typeof targetProp === 'object' && targetProp !== null) {
				currentTarget = targetProp as Schema;
			}
			currentUiTarget = currentUiTarget[part];
		}
	}

	return { schema: rootSchema, uiSchema, initialValue: {}, mapPaths };
}

/**
 * Access nested property by dot path
 */
function getByPath(obj: any, path: string): any {
	return path.split('.').reduce((acc, part) => acc && acc[part], obj);
}

/**
 * Set nested property by dot path, creating objects as needed
 */
function setByPath(obj: any, path: string, value: any): void {
	const parts = path.split('.');
	let current = obj;
	for (let i = 0; i < parts.length - 1; i++) {
		const part = parts[i];
		if (!current[part]) current[part] = {};
		current = current[part];
	}
	current[parts[parts.length - 1]] = value;
}

/**
 * Converts K8s Object data (Standard) to Form data (Array-based Maps)
 */

/**
 * Converts K8s Object data (Standard) to Form data (Array-based Maps)
 */
export function k8sToFormData(data: unknown, mapPaths: string[]): Record<string, unknown> {
	if (!data || typeof data !== 'object' || Array.isArray(data)) return {};
	// Deep clone simple data
	const formData = JSON.parse(JSON.stringify(data));

	for (const path of mapPaths) {
		const originalValue = getByPath(formData, path);
		if (originalValue && typeof originalValue === 'object' && !Array.isArray(originalValue)) {
			// Convert Object {k:v} to Array [{key:k, value:v}]
			const arrayValue = Object.entries(originalValue).map(([key, value]) => ({
				key,
				value
			}));
			setByPath(formData, path, arrayValue);
		} else if (originalValue === undefined || originalValue === null) {
			// Ensure it is initialized as array if missing (optional, but good for UI)
			setByPath(formData, path, []);
		}
	}
	return formData;
}

/**
 * Converts Form data (Array-based Maps) back to K8s Object data
 */
export function formDataToK8s(formData: unknown, mapPaths: string[]): Record<string, unknown> {
	if (!formData || typeof formData !== 'object' || Array.isArray(formData)) return {};
	const k8sData = JSON.parse(JSON.stringify(formData));

	for (const path of mapPaths) {
		const arrayValue = getByPath(k8sData, path);
		if (Array.isArray(arrayValue)) {
			// Convert Array [{key:k, value:v}] to Object {k:v}
			const objectValue = arrayValue.reduce((acc: Record<string, any>, item: any) => {
				if (item && item.key) {
					acc[item.key] = item.value;
				}
				return acc;
			}, {});
			setByPath(k8sData, path, objectValue);
		}
	}
	return k8sData;
}

