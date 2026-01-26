/**
 * Deep merge two objects.
 */
export function deepMerge(target: any, source: any): any {
	if (typeof target !== 'object' || target === null) return source;
	if (typeof source !== 'object' || source === null) return target;

	const output = { ...target };
	Object.keys(source).forEach((key) => {
		// Guard against prototype pollution
		if (key === '__proto__' || key === 'constructor' || key === 'prototype') return;
		if (typeof source[key] === 'object' && source[key] !== null && key in target) {
			output[key] = deepMerge(target[key], source[key]);
		} else {
			output[key] = source[key];
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
