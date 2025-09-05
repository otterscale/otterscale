import IPv4CIDRInput from './ipv4.svelte';
import { isNumber } from './is-number';

/** Attempts to parse the provided address into a valid IP. Returns undefined for
 * undefined returns a valid IP in array form for a valid IP and returns a 0 filled array for a incomplete IP.
 *
 * **This is used only for parsing the placeholder**
 *
 * @param IPv4CIDR IP Address string to be parsed can be `0.0.0.0/16` or `0 0 0 0/16` or `0_0_0_0/16` or `0 0 0` (partially complete)
 * @returns
 */
export const safeParseIPv4CIDR = (
	IPv4CIDR: string | undefined,
): [string | null, string | null, string | null, string | null, string | null] | undefined => {
	if (IPv4CIDR === undefined) return undefined;
	let ip = IPv4CIDR.trim();

	ip = ip.replaceAll('_', '.');
	ip = ip.replaceAll(' ', '.');
	ip = ip.replaceAll('/', '.');

	const segments: (string | null)[] = ip.split('.');

	while (segments.length < 5) {
		segments.push(null);
	}

	for (let i = 0; i < segments.length; i++) {
		if (!isNumber(segments[i]) || segments[i] === null) {
			segments[i] = null;
			continue;
		}

		const num = Number.parseInt(segments[i]!);

		if (num < 0 || num > 255) {
			segments[i] = null;
			continue;
		}

		segments[i] = num.toString();
	}

	// @ts-expect-error We know this is 4 we just made sure
	return segments;
};

export { IPv4CIDRInput };

export type * from './types';
