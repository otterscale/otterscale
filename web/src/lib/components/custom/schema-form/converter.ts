import { deepMerge, deleteByPath,getByPath, setByPath } from './utils';

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
	transformationMappings: Record<string, string>; // K8s Path -> Form Path
}

export interface PathOptions {
	title?: string;
	required?: boolean;
	showDescription?: boolean;
	uiSchema?: Record<string, unknown>;
}

/**
 * Subsets the full OpenAPI schema to include only the specified paths.
 * @param fullSchema The full OpenAPI V3 schema object
 * @param paths Array of dot-notation paths to include (e.g. "metadata.name", "spec.running")
 * @returns An object containing { schema, uiSchema, transformationMappings }
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

	const pathKeys = Array.isArray(paths) ? paths : Object.keys(paths);
	const pathOptions = Array.isArray(paths) ? {} : paths;

	for (const path of pathKeys) {
		const options: PathOptions = pathOptions[path] || {};
		const customUiOptions = options.uiSchema || {};

		const parts = path.split('.');
		let currentSource: K8sOpenAPISchema = fullSchema;
		let currentTarget: Schema = rootSchema;
		let currentUiTarget: any = uiSchema;

		// Track cumulative path to identify implicit nodes
		let cumulativePath = '';

		for (let i = 0; i < parts.length; i++) {
			const part = parts[i];
			cumulativePath = cumulativePath ? `${cumulativePath}.${part}` : part;

			// Check if this current path node is explicitly requested by the user
			const isExplicit = pathKeys.includes(cumulativePath);

			// Handling implicit array traversal:
			// If currentSource is an Array, we expect to find the property inside 'items'
			if (currentSource.type === 'array' && currentSource.items) {
				currentSource = currentSource.items;

				if (!currentTarget.items) {
					currentTarget.items = { type: 'object', properties: {} };
				}
				// Only proceed if items is a Schema object
				if (typeof currentTarget.items === 'object' && !Array.isArray(currentTarget.items)) {
					currentTarget = currentTarget.items as Schema;
				}

				if (!currentUiTarget.items) {
					currentUiTarget.items = {};
				}
				currentUiTarget = currentUiTarget.items;
			}

			// Locate in Source
			// Capture parent source for checking 'required' later
			const parentSource = currentSource;

			if (!currentSource.properties?.[part]) {
				// Check if current source is a Map (has additionalProperties schema)
				if (
					currentSource.type === 'object' &&
					currentSource.additionalProperties &&
					typeof currentSource.additionalProperties === 'object'
				) {
					// It's a Map, so we are traversing into a key of the map.
					// The schema for this key is the schema defined in additionalProperties.
					currentSource = currentSource.additionalProperties as K8sOpenAPISchema;
					// We successfully "found" the property (conceptually), but it's dynamic.
					// We continue traversal with this new source.
					// Note: We don't have a concrete "sourceProp" object for this dynamic key in the original schema's "properties" map,
					// but currentSource is now pointing to the correct schema for this node.
				} else {
					console.warn(`Path segment "${part}" not found in source schema (full path: ${path})`);
					break;
				}
			} else {
				currentSource = currentSource.properties[part];
			}
			
			const sourceProp = currentSource;
			const isLeaf = i === parts.length - 1;

			// Ensure Target Properties exists
			if (!currentTarget.properties) {
				currentTarget.properties = {};
			}

			// Function to apply common transformations (title, description)
			const applyOptions = (target: Schema, src: K8sOpenAPISchema, isLeafNode: boolean) => {
				// Title: Use options.title if leaf and provided, else source title
				if (options.title) {
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

				// Handle Required
				if (options.required) {
					if (!currentTarget.required) currentTarget.required = [];
					if (!currentTarget.required.includes(part)) {
						currentTarget.required.push(part);
					}
				}
			};

			// Check if this property is a Map (object with additionalProperties)
			const isMap = sourceProp.type === 'object' && !!sourceProp.additionalProperties;

			// HOISTING LOGIC:
			// If it is a Map, we hoist it to avoid nested object update issues in the form library.
			if (isLeaf && isMap) {
				// It's a Map. Hoist it.
				const formPath = path.replace(/\./g, '_'); // e.g. metadata_annotations
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

				// Add to ROOT schema properties
				if (!rootSchema.properties) rootSchema.properties = {};
				rootSchema.properties[formPath] = newProp;

				// Apply custom UI options to Hoisted field (root level in uiSchema)
				if (Object.keys(customUiOptions).length > 0) {
					uiSchema[formPath] = deepMerge(uiSchema[formPath] || {}, customUiOptions);
				}

				// Stop processing this path
				break;
			}

			// Standard Logic (Non-Hoist)
			// Manage UI Schema for non-hoisted
			if (!currentUiTarget[part]) {
				currentUiTarget[part] = {};
			}

			// Logic to hide labels:
			// 1. Implicit intermediate nodes (!isExplicit): Always hide to avoid clutter.
			if (!isExplicit) {
				currentUiTarget[part]['ui:options'] = { label: false };
			}

			if (isLeaf) {
				// SIMPLIFY QUANTITY (oneOf: [string, number]) -> type: string
				// This prevents the UI from rendering a complex selector for CPU/Memory fields.
				if (
					sourceProp.oneOf &&
					sourceProp.oneOf.some((o: any) => o.type === 'string') &&
					sourceProp.oneOf.some((o: any) => o.type === 'number' || o.type === 'integer')
				) {
					currentUiTarget[part] = currentUiTarget[part] || {};
					currentTarget.properties[part] = {
						type: 'string',
						title: sourceProp.title || '',
						description: sourceProp.description || ''
					};
					// Important: Update sourceProp to point to this new string schema so downstream logic works
					// However, we just modified target directly, so we should be good.
				}

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

				// Apply custom UI options to Leaf field
				if (Object.keys(customUiOptions).length > 0) {
					currentUiTarget[part] = deepMerge(currentUiTarget[part] || {}, customUiOptions);
				}

				if (options.title) {
					currentUiTarget[part]['ui:title'] = options.title;

					// Force label to be true if title is provided, to override any implicit hiding
					currentUiTarget[part]['ui:options'] = deepMerge(
						currentUiTarget[part]['ui:options'] || {},
						{ label: true }
					);
				}
			}

			const isTerminalLeaf =
				isLeaf && !pathKeys.some((k) => k !== cumulativePath && k.startsWith(cumulativePath + '.'));

			if (!currentTarget.properties[part]) {
				if (isTerminalLeaf) {
					// Terminal Leaf: Full deep copy because no children are customised below this point
					const newProp = { ...sourceProp } as Schema;
					if (newProp.type === 'object' && !newProp.properties) {
						newProp.properties = {};
					}
					applyOptions(newProp, sourceProp, true);
					currentTarget.properties[part] = newProp;
				} else {
					// Intermediate (Implicit OR Explicit Parent): Skeleton only
					if (sourceProp.type === 'array') {
						currentTarget.properties[part] = {
							type: 'array',
							items: { type: 'object', properties: {} }
						};
					} else {
						currentTarget.properties[part] = {
							type: 'object',
							properties: {},
							additionalProperties: true // Keep this for safety
						};
					}

					// Intermediate nodes: If implicit, remove title to "hide" it visually in some renderers,
					// combined with ui:options: { label: false } above.
					applyOptions(currentTarget.properties[part] as Schema, sourceProp, false);
					if (!isExplicit) {
						(currentTarget.properties[part] as Schema).title = '';
					}
				}
			} else if (isLeaf) {
				const target = currentTarget.properties[part] as Schema;
				applyOptions(target, sourceProp, true);
			}

			// Handle Required
			const isRequired =
				Array.isArray(parentSource.required) && parentSource.required.includes(part);
			if (isRequired) {
				if (path in transformationMappings) {
					const formPath = transformationMappings[path];
					if (!rootSchema.required) rootSchema.required = [];
					if (!rootSchema.required.includes(formPath)) {
						rootSchema.required.push(formPath);
					}
				} else {
					if (!currentTarget.required) currentTarget.required = [];
					if (!currentTarget.required.includes(part)) {
						currentTarget.required.push(part);
					}
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

	return { schema: rootSchema, uiSchema, transformationMappings };
}



/**
 * Converts K8s Object data (Standard) to Form data (Flattened/Mapped)
 */
export function k8sToFormData(
	data: unknown,
	mappings: Record<string, string>
): Record<string, unknown> {
	if (!data || typeof data !== 'object' || Array.isArray(data)) return {};
	// Deep clone simple data
	const formData = JSON.parse(JSON.stringify(data));

	for (const [k8sPath, formPath] of Object.entries(mappings)) {
		const originalValue = getByPath(formData, k8sPath);

		// If the value exists, move it to formPath and transform it
		if (originalValue && typeof originalValue === 'object' && !Array.isArray(originalValue)) {
			// Convert Object {k:v} to Array [{key:k, value:v}]
			const arrayValue = Object.entries(originalValue).map(([key, value]) => ({
				key,
				value
			}));
			// Set at Form Path (Hoisted)
			formData[formPath] = arrayValue;

			// Remove the original nested path to avoid duplication
			deleteByPath(formData, k8sPath);
		} else if (originalValue === undefined || originalValue === null) {
			formData[formPath] = [];
			// Also ensure original path is removed if it existed as null/undefined
			deleteByPath(formData, k8sPath);
		}
	}
	return formData;
}

/**
 * Converts Form data (Flattened/Mapped) back to K8s Object data
 */
export function formDataToK8s(
	formData: unknown,
	mappings: Record<string, string>
): Record<string, unknown> {
	if (!formData || typeof formData !== 'object' || Array.isArray(formData)) return {};
	const k8sData = JSON.parse(JSON.stringify(formData));

	for (const [k8sPath, formPath] of Object.entries(mappings)) {
		const arrayValue = k8sData[formPath]; // Get from Hoisted

		if (Array.isArray(arrayValue)) {
			// Convert Array [{key:k, value:v}] to Object {k:v}
			const objectValue = arrayValue.reduce((acc: Record<string, any>, item: any) => {
				if (item && item.key) {
					acc[item.key] = item.value;
				}
				return acc;
			}, {});
			setByPath(k8sData, k8sPath, objectValue); // Put back to Deep

			// Remove Hoisted key from output
			delete k8sData[formPath];
		}
	}
	return k8sData;
}
