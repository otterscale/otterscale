/*
    Installed from @ieedan/std
*/

import { isNumber } from './is-number';
import { Err, Ok, type Result } from './result';

export type Octets = [number, number, number, number, number];

export type IPv4CIDR =
    | Octets
    | `${number}.${number}.${number}.${number}/${number}`
    | `${number} ${number} ${number} ${number}/${number}`
    | `${number}_${number}_${number}_${number}/${number}`;

export type ParseError = {
    octet?: number;
    message: string;
};

/** Parses the ip address from a string in the format of `0.0.0.0/16` or `0 0 0 0/16` or `0_0_0_0/16` into an array of octets
 *
 * @param address
 * @returns
 *
 * ## Usage
 *
 * ```ts
 * parse("192.168.100.10").unwrap(); // [192, 168, 100, 10]
 * ```
 */
export function parse(address: string): Result<Octets, ParseError> {
    let newAddress = address.trim();

    newAddress = newAddress.replaceAll(' ', '.');
    newAddress = newAddress.replaceAll('_', '.');
    newAddress = newAddress.replaceAll('/', '.');

    const octets = newAddress.split('.');

    if (octets.length !== 5)
        return Err({ message: `'${address}' is invalid as it doesn't contain 4 octets.` });

    const final: Octets = [0, 0, 0, 0, 0];

    for (let i = 0; i < octets.length; i++) {
        const octet = octets[i];

        if (!isNumber(octet)) return Err({ octet: i + 1, message: `'${octet}' is not a number.` });

        const num = Number.parseInt(octet);

        if (num < 0 || num > 255) return Err({ octet: i + 1, message: `'${octet}' is out of range.` });

        final[i] = num;
    }

    return Ok(final);
}

/** Validates the provided address
 *
 * @param address
 * @returns
 *
 * ## Usage
 *
 * ```ts
 * validate("192.168.100.10"); // true
 * validate([192, 168, 100, 10]); // true
 *
 * validate("192.168.100.256"); // false
 * validate([192, 168, 100, 256]); // false
 * ```
 */
export function validate(address: IPv4CIDR): boolean {
    if (typeof address === 'string') return parse(address).isOk();

    for (let i = 0; i < address.length; i++) {
        const octet = address[i];

        if (octet < 0 || octet > 255) return false;
    }

    return true;
}

/** Formats the provided address to a string with the provided separator
 *
 * @param address
 * @param separator @default "."
 * @returns
 *
 * ## Usage
 *
 * ```ts
 * formatToString([192, 168, 100, 10]); // "192.168.100.10"
 * ```
 */
export function formatToString(
    address: IPv4CIDR,
    separator: '.' | '_' | ' ' = '.'
): Result<string, string> {
    if (Array.isArray(address)) return Ok(address.slice(0, 4).join(separator) + '/' + address[4]);

    const parsed = parse(address);

    if (parsed.isErr()) return Err(parsed.unwrapErr().message);

    return formatToString(parsed.unwrap(), separator);
}