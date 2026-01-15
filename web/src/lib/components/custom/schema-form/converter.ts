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
	[key: string]: unknown;
}

/** Result of building schema from K8s OpenAPI */
export interface SchemaFormConfig {
	schema: Schema;
	uiSchema: UiSchemaRoot;
	initialValue: Record<string, unknown>;
}

/**
 * Subsets the full OpenAPI schema to include only the specified paths.
 * @param fullSchema The full OpenAPI V3 schema object
 * @param paths Array of dot-notation paths to include (e.g. "metadata.name", "spec.running")
 * @returns An object containing { schema, uiSchema, initialValue }
 */
export function buildSchemaFromK8s(
	fullSchema: K8sOpenAPISchema,
	paths: string[]
): SchemaFormConfig {
	const rootSchema: Schema = {
		type: 'object',
		properties: {},
		required: [],
		title: fullSchema.title ?? '',
		description: fullSchema.description ?? ''
	};

	const uiSchema: UiSchemaRoot = {};

	for (const path of paths) {
		const parts = path.split('.');
		let currentSource: K8sOpenAPISchema = fullSchema;
		let currentTarget: Schema = rootSchema;

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
			}

			// Locate in Source
			if (!currentSource.properties?.[part]) {
				console.warn(`Path segment "${part}" not found in source schema (full path: ${path})`);
				break;
			}
			const sourceProp = currentSource.properties[part];

			// Ensure Target Properties exists
			if (!currentTarget.properties) {
				currentTarget.properties = {};
			}

			const isLeaf = i === parts.length - 1;

			// Prepare Target Property
			if (!currentTarget.properties[part]) {
				if (isLeaf) {
					// Leaf: Full Copy
					currentTarget.properties[part] = { ...sourceProp } as SchemaDefinition;
				} else {
					// Intermediate: Skeleton
					if (sourceProp.type === 'array') {
						currentTarget.properties[part] = {
							type: 'array',
							title: sourceProp.title,
							description: sourceProp.description,
							items: { type: 'object', properties: {} }
						};
					} else {
						currentTarget.properties[part] = {
							type: 'object',
							title: sourceProp.title,
							description: sourceProp.description,
							properties: {}
						};
					}
				}
			} else if (isLeaf) {
				// Exists but we're at leaf - upgrade to full copy
				currentTarget.properties[part] = { ...sourceProp } as SchemaDefinition;
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
		}
	}

	const initialValue = generateInitialValue(rootSchema);

	return { schema: rootSchema, uiSchema, initialValue };
}

function generateInitialValue(schema: Schema): Record<string, unknown> {
	if (schema.type === 'object' && schema.properties) {
		const obj: Record<string, unknown> = {};
		for (const [key, prop] of Object.entries(schema.properties)) {
			if (typeof prop === 'object' && prop !== null) {
				const val = generateInitialValue(prop as Schema);
				if (val !== undefined && Object.keys(val).length > 0) {
					obj[key] = val;
				}
			}
		}
		return obj;
	} else if (schema.type === 'array') {
		// For array, default is empty array
		return {};
	}
	if (schema.default !== undefined) {
		return { value: schema.default };
	}
	return {};
}
