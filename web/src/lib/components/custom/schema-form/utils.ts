/**
 * Check if value is empty (null, undefined, empty object, or array of empty objects)
 */
function isEmptyValue(value: any): boolean {
	if (value === null || value === undefined) return true;
	if (Array.isArray(value)) {
		return value.length === 0 || value.every((item) => isEmptyValue(item));
	}
	if (typeof value === 'object') {
		return Object.keys(value).length === 0;
	}
	return false;
}

/**
 * Deep merge two objects.
 * Skips merging if source value is empty to preserve existing data.
 */
export function deepMerge(target: any, source: any): any {
	if (typeof target !== 'object' || target === null) return source;
	if (typeof source !== 'object' || source === null) return target;

	// Arrays: merge element-by-element to preserve fields not in source
	if (Array.isArray(target) || Array.isArray(source)) {
		// If source is empty or contains only empty objects, keep target
		if (isEmptyValue(source)) {
			return target;
		}
		// If both are arrays, deep merge each element by index
		if (Array.isArray(target) && Array.isArray(source)) {
			const result = [...target];
			for (let i = 0; i < source.length; i++) {
				if (i < result.length) {
					if (isEmptyValue(source[i])) {
						// Keep original element if source element is empty
						continue;
					}
					result[i] = deepMerge(result[i], source[i]);
				} else {
					result[i] = source[i];
				}
			}
			return result;
		}
		return source;
	}

	const output = { ...target };
	Object.keys(source).forEach((key) => {
		// Guard against prototype pollution
		if (key === '__proto__' || key === 'constructor' || key === 'prototype') return;

		const sourceValue = source[key];
		const targetValue = target[key];

		// Skip if source value is empty and target has data
		if (isEmptyValue(sourceValue) && targetValue !== undefined && !isEmptyValue(targetValue)) {
			return;
		}

		if (typeof sourceValue === 'object' && sourceValue !== null && key in target) {
			output[key] = deepMerge(targetValue, sourceValue);
		} else {
			output[key] = sourceValue;
		}
	});
	return output;
}

/**
 * Access nested property by dot path
 */
export function getByPath(obj: any, path: string): any {
	return path.split('.').reduce((acc, part) => acc && acc[part], obj);
}

/**
 * Set nested property by dot path, creating objects as needed
 */
export function setByPath(obj: any, path: string, value: any): void {
	const parts = path.split('.');
	// Guard against prototype pollution
	if (
		parts.some((part) => part === '__proto__' || part === 'constructor' || part === 'prototype')
	) {
		return;
	}
	let current = obj;
	for (let i = 0; i < parts.length - 1; i++) {
		const part = parts[i];
		if (!current[part]) current[part] = {};
		current = current[part];
	}
	current[parts[parts.length - 1]] = value;
}

/**
 * Deletes nested property by dot path
 */
export function deleteByPath(obj: any, path: string): void {
	const parts = path.split('.');
	// Guard against prototype pollution
	if (
		parts.some((part) => part === '__proto__' || part === 'constructor' || part === 'prototype')
	) {
		return;
	}
	let current = obj;
	for (let i = 0; i < parts.length - 1; i++) {
		const part = parts[i];
		if (!current[part]) return;
		current = current[part];
	}
	delete current[parts[parts.length - 1]];
}
