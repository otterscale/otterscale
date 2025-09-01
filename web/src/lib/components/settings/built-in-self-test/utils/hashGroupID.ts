/**
 * Generates a deterministic hash from a string input.
 * @param str Input string
 * @returns Hash number
 */
export function hashCode(str: string): number {
	let hash = 0;
	for (let i = 0; i < str.length; i++) {
		hash = (hash << 5) - hash + str.charCodeAt(i);
		hash |= 0; // Convert to 32bit integer
	}
	return hash;
}

/**
 * Encodes an object into a hash string (uid-<hash>).
 * @param obj Object to encode
 * @returns Encoded string
 */
export function encode(obj: Record<string, any>): string {
	const str = Object.values(obj).join('-');
	const hash = hashCode(str);
	return `${Math.abs(hash)}`;
}
