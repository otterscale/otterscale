import type { JsonValue } from '@bufbuild/protobuf';

function findSuffix(quantity: string): string {
	let ix = quantity.length - 1;
	while (ix >= 0 && !/[.0-9]/.test(quantity.charAt(ix))) {
		ix--;
	}
	return ix === -1 ? '' : quantity.substring(ix + 1);
}

function quantityToScalar(quantity: string): number | bigint {
	if (!quantity) {
		return 0;
	}
	const suffix = findSuffix(quantity);
	if (suffix === '') {
		const num = Number(quantity).valueOf();
		if (isNaN(num)) {
			throw new Error('Unknown quantity ' + quantity);
		}
		return num;
	}
	switch (suffix) {
		case 'n':
			return Number(quantity.substr(0, quantity.length - 1)).valueOf() / 1_000_000_000.0;
		case 'u':
			return Number(quantity.substr(0, quantity.length - 1)).valueOf() / 1_000_000.0;
		case 'm':
			return Number(quantity.substr(0, quantity.length - 1)).valueOf() / 1000.0;
		case 'k':
			return BigInt(quantity.substr(0, quantity.length - 1)) * BigInt(1000);
		case 'M':
			return BigInt(quantity.substr(0, quantity.length - 1)) * BigInt(1000 * 1000);
		case 'G':
			return BigInt(quantity.substr(0, quantity.length - 1)) * BigInt(1000 * 1000 * 1000);
		case 'T':
			return (
				BigInt(quantity.substr(0, quantity.length - 1)) * BigInt(1000 * 1000 * 1000) * BigInt(1000)
			);
		case 'P':
			return (
				BigInt(quantity.substr(0, quantity.length - 1)) *
				BigInt(1000 * 1000 * 1000) *
				BigInt(1000 * 1000)
			);
		case 'E':
			return (
				BigInt(quantity.substr(0, quantity.length - 1)) *
				BigInt(1000 * 1000 * 1000) *
				BigInt(1000 * 1000 * 1000)
			);
		case 'Ki':
			return BigInt(quantity.substr(0, quantity.length - 2)) * BigInt(1024);
		case 'Mi':
			return BigInt(quantity.substr(0, quantity.length - 2)) * BigInt(1024 * 1024);
		case 'Gi':
			return BigInt(quantity.substr(0, quantity.length - 2)) * BigInt(1024 * 1024 * 1024);
		case 'Ti':
			return (
				BigInt(quantity.substr(0, quantity.length - 2)) * BigInt(1024 * 1024 * 1024) * BigInt(1024)
			);
		case 'Pi':
			return (
				BigInt(quantity.substr(0, quantity.length - 2)) *
				BigInt(1024 * 1024 * 1024) *
				BigInt(1024 * 1024)
			);
		case 'Ei':
			return (
				BigInt(quantity.substr(0, quantity.length - 2)) *
				BigInt(1024 * 1024 * 1024) *
				BigInt(1024 * 1024 * 1024)
			);
		default:
			throw new Error(`Unknown suffix: ${suffix}`);
	}
}

function formatWithDecimalPrefix(value: bigint): { value: number; unit: string } {
	const units = [
		{ value: BigInt(1e18), symbol: 'E' },
		{ value: BigInt(1e15), symbol: 'P' },
		{ value: BigInt(1e12), symbol: 'T' },
		{ value: BigInt(1e9), symbol: 'G' },
		{ value: BigInt(1e6), symbol: 'M' },
		{ value: BigInt(1e3), symbol: 'k' },
		{ value: BigInt(1), symbol: '' }
	];
	for (const unit of units) {
		if (value >= unit.value) {
			const number = Number(value) / Number(unit.value);
			return {
				value: number,
				unit: unit.symbol
			};
		}
	}
	return {
		value: Number(value),
		unit: ''
	};
}

function formatWithBinaryPrefix(value: bigint): { value: number; unit: string } {
	const units = [
		{ value: BigInt(2) ** BigInt(60), symbol: 'Ei' },
		{ value: BigInt(2) ** BigInt(50), symbol: 'Pi' },
		{ value: BigInt(2) ** BigInt(40), symbol: 'Ti' },
		{ value: BigInt(2) ** BigInt(30), symbol: 'Gi' },
		{ value: BigInt(2) ** BigInt(20), symbol: 'Mi' },
		{ value: BigInt(2) ** BigInt(10), symbol: 'Ki' },
		{ value: BigInt(1), symbol: '' }
	];
	for (const unit of units) {
		if (value >= unit.value) {
			const number = Number(value) / Number(unit.value);
			return {
				value: number,
				unit: unit.symbol
			};
		}
	}
	return {
		value: Number(value),
		unit: ''
	};
}

function getQuantityScalar(quantity: string | number | null): number | bigint | null {
	if (quantity === null) return null;

	return quantityToScalar(String(quantity));
}

function getRatio(
	numerator: number | bigint | null,
	denominator: number | bigint | null
): JsonValue {
	if (numerator === null || denominator === null) return null;
	return Number(numerator) / Number(denominator);
}

function format(value: string) {
	try {
		return JSON.stringify(JSON.parse(value), null, 4);
	} catch {
		return value;
	}
}

function getRelativeTime(now: number, timestamp: number) {
	const milliseconds = Math.max(timestamp, 0);

	const seconds = Math.floor((now - milliseconds) / 1000);
	if (seconds < 60) return { value: seconds, unit: 'second' };

	const minutes = Math.floor(seconds / 60);
	if (minutes < 60) return { value: minutes, unit: 'minute' };

	const hours = Math.floor(minutes / 60);
	if (hours < 24) return { value: hours, unit: 'hour' };

	const days = Math.floor(hours / 24);
	if (days < 7) return { value: days, unit: 'day' };

	const weeks = Math.floor(days / 7);
	if (weeks < 5) return { value: weeks, unit: 'week' };

	const months = Math.floor(days / 30);
	if (months < 12) return { value: months, unit: 'month' };

	const years = Math.floor(days / 365);
	return { value: years, unit: 'year' };
}

type UISchemaType =
	| 'boolean'
	| 'number'
	| 'number-with-prefix'
	| 'time'
	| 'text'
	| 'array'
	| 'array-of-enumeration'
	| 'array-of-object'
	| 'object'
	| 'object-of-key-value'
	| 'link'
	| 'ratio'
	| undefined;
function getDefaultUISchema(type: JsonValue, format?: JsonValue): UISchemaType {
	if (type === 'boolean') {
		return 'boolean';
	}
	if (type === 'number' || type === 'integer') {
		return 'number';
	}
	if (type === 'string' && (format === 'date' || format === 'date-time')) {
		return 'time';
	}
	if (type === 'string') {
		return 'text';
	}
	if (type === 'array') {
		return 'array';
	}
	if (type === 'object') {
		return 'object';
	}
	return undefined;
}

type DataSchemaType = 'boolean' | 'number' | 'time' | 'text' | 'array' | 'object' | undefined;
function getDefaultDataSchema(type: JsonValue | undefined, format?: JsonValue): DataSchemaType {
	if (type === 'boolean') {
		return 'boolean';
	}
	if (type === 'number' || type === 'integer') {
		return 'number';
	}
	if (type === 'string' && (format === 'date' || format === 'date-time')) {
		return 'time';
	}
	if (type === 'string') {
		return 'text';
	}
	if (type === 'array') {
		return 'array';
	}
	if (type === 'object') {
		return 'object';
	}
	return undefined;
}

export {
	format,
	formatWithBinaryPrefix,
	formatWithDecimalPrefix,
	getDefaultDataSchema,
	getDefaultUISchema,
	getQuantityScalar,
	getRatio,
	getRelativeTime,
	quantityToScalar
};
export type { DataSchemaType, UISchemaType };
