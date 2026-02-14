// ── Predicates ─────────────────────────────────────────────

const UNSAFE_KEYS = new Set(['__proto__', 'constructor', 'prototype']);

/** null, undefined, {}, [], or array of only empty values. */
function isEmptyValue(value: any): boolean {
	if (value === null || value === undefined) return true;
	if (Array.isArray(value)) return value.length === 0 || value.every(isEmptyValue);
	if (typeof value === 'object') return Object.keys(value).length === 0;
	return false;
}

// ── Deep Merge ─────────────────────────────────────────────

/**
 * Recursively merges `source` into `target`. Skips empty source values to preserve existing data.
 *
 * @example
 * deepMerge({ a: 1, b: { x: 1 } }, { b: { y: 2 }, c: 3 })
 * // → { a: 1, b: { x: 1, y: 2 }, c: 3 }
 *
 * deepMerge([{ name: "a" }], [{}])
 * // → [{ name: "a" }]  (empty source element is skipped)
 */
export function deepMerge(target: any, source: any): any {
	if (typeof target !== 'object' || target === null) return source;
	if (typeof source !== 'object' || source === null) return target;

	if (Array.isArray(target) || Array.isArray(source)) {
		return mergeArrays(target, source);
	}

	return mergeObjects(target, source);
}

function mergeArrays(target: any, source: any): any {
	if (isEmptyValue(source)) return target;

	if (Array.isArray(target) && Array.isArray(source)) {
		const result = [...target];
		for (let i = 0; i < source.length; i++) {
			if (i < result.length) {
				if (isEmptyValue(source[i])) continue;
				result[i] = deepMerge(result[i], source[i]);
			} else {
				result[i] = source[i];
			}
		}
		return result;
	}

	return source;
}

function mergeObjects(target: any, source: any): any {
	const output = { ...target };

	for (const key of Object.keys(source)) {
		if (UNSAFE_KEYS.has(key)) continue;

		const srcVal = source[key];
		const tgtVal = target[key];

		if (isEmptyValue(srcVal) && tgtVal !== undefined && !isEmptyValue(tgtVal)) continue;

		output[key] =
			typeof srcVal === 'object' && srcVal !== null && key in target
				? deepMerge(tgtVal, srcVal)
				: srcVal;
	}

	return output;
}

// ── Dot-Path Accessors ─────────────────────────────────────

/**
 * @example getByPath({ a: { b: 1 } }, "a.b") // → 1
 */
export function getByPath(obj: any, path: string): any {
	return path.split('.').reduce((acc, part) => acc && acc[part], obj);
}

/**
 * @example
 * const o = {};
 * setByPath(o, "a.b.c", 42) // → o = { a: { b: { c: 42 } } }
 */
export function setByPath(obj: any, path: string, value: any): void {
	const parts = path.split('.');
	if (parts.some((p) => UNSAFE_KEYS.has(p))) return;

	let current = obj;
	for (let i = 0; i < parts.length - 1; i++) {
		if (!current[parts[i]]) current[parts[i]] = {};
		current = current[parts[i]];
	}
	current[parts[parts.length - 1]] = value;
}

/**
 * @example
 * const o = { a: { b: 1, c: 2 } };
 * deleteByPath(o, "a.b") // → o = { a: { c: 2 } }
 */
export function deleteByPath(obj: any, path: string): void {
	const parts = path.split('.');
	if (parts.some((p) => UNSAFE_KEYS.has(p))) return;

	let current = obj;
	for (let i = 0; i < parts.length - 1; i++) {
		if (!current[parts[i]]) return;
		current = current[parts[i]];
	}
	delete current[parts[parts.length - 1]];
}
