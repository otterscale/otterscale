import { getLocale } from './paraglide/runtime';
import { m } from '$lib/paraglide/messages';

const TIME_DIVISIONS = [
	{ amount: 60, name: 'seconds' },
	{ amount: 60, name: 'minutes' },
	{ amount: 24, name: 'hours' },
	{ amount: 7, name: 'days' },
	{ amount: 4.34524, name: 'weeks' },
	{ amount: 12, name: 'months' },
	{ amount: Number.POSITIVE_INFINITY, name: 'years' },
] as const;

export function formatTimeAgo(date: Date): string {
	const formatter = new Intl.RelativeTimeFormat(getLocale(), {
		numeric: 'auto',
	});

	let duration = (date.getTime() - Date.now()) / 1000;

	for (const division of TIME_DIVISIONS) {
		if (Math.abs(duration) < division.amount) {
			return formatter.format(Math.round(duration), division.name);
		}
		duration /= division.amount;
	}

	return formatter.format(Math.round(duration), 'years');
}

export function formatDuration(duration: number): { value: number; unit: string } {
	if (duration === 0) return { value: 0, unit: m.second() };

	const years = duration / (365 * 24 * 3600);
	if (years >= 1) return { value: years, unit: m.year() };

	const weeks = (duration % (365 * 24 * 3600)) / (7 * 24 * 3600);
	if (weeks >= 1) return { value: weeks, unit: m.week() };

	const days = (duration % (7 * 24 * 3600)) / (24 * 3600);
	if (days >= 1) return { value: days, unit: m.day() };

	const hours = (duration % (24 * 3600)) / 3600;
	if (hours >= 1) return { value: hours, unit: m.hour() };

	const minutes = (duration % 3600) / 60;
	if (minutes >= 1) return { value: minutes, unit: m.minute() };

	const seconds = duration % 60;
	return { value: seconds, unit: m.second() };
}

export function formatCapacity(capacity: number | bigint): { value: number; unit: string } {
	const B = Number(capacity);
	const KB = B / 1024;
	const MB = KB / 1024;
	const GB = MB / 1024;
	const TB = GB / 1024;

	if (TB >= 1) {
		return { value: Math.round(TB * 100) / 100, unit: 'TB' };
	} else if (GB >= 1) {
		return { value: Math.round(GB * 100) / 100, unit: 'GB' };
	} else if (MB >= 1) {
		return { value: Math.round(MB * 100) / 100, unit: 'MB' };
	} else {
		return { value: Math.round(KB * 100) / 100, unit: 'KB' };
	}
}

export function formatIO(bytes: number | bigint): { value: number; unit: string } {
	const B = Number(bytes);
	const KB = B / 1024;
	const MB = KB / 1024;
	const GB = MB / 1024;
	const TB = GB / 1024;

	if (TB >= 1) {
		return { value: Math.round(TB * 100) / 100, unit: 'TB/s' };
	} else if (GB >= 1) {
		return { value: Math.round(GB * 100) / 100, unit: 'GB/s' };
	} else if (MB >= 1) {
		return { value: Math.round(MB * 100) / 100, unit: 'MB/s' };
	} else if (KB >= 1) {
		return { value: Math.round(KB * 100) / 100, unit: 'KB/s' };
	} else {
		return { value: Math.round(B * 100) / 100, unit: 'B/s' };
	}
}

export function formatBigNumber(number: number | bigint) {
	return number.toLocaleString('en-US');
}

// export const formatTime = (v: Date | number): string => {
//     return dayjs(v).format('HH:mm');
// };

export function formatLatency(second: number): { value: number; unit: string } {
	const millisecond = second * 1000;
	const microsecond = millisecond * 1000;
	const nanosecond = microsecond * 1000;

	if (second >= 1) {
		return { value: Math.round(second * 100) / 100, unit: 's' };
	} else if (millisecond >= 1) {
		return { value: Math.round(millisecond * 100) / 100, unit: 'ms' };
	} else if (microsecond >= 1) {
		return { value: Math.round(microsecond * 100) / 100, unit: 'us' };
	} else {
		return { value: Math.round(nanosecond * 100) / 100, unit: 'ns' };
	}
}

export function formatLatencyNano(nanosecond: number): { value: number; unit: string } {
	const second = nanosecond / 1000000000;
	const millisecond = nanosecond / 1000000;
	const microsecond = nanosecond / 1000;

	if (second >= 1) {
		return { value: Math.round(second * 100) / 100, unit: 's' };
	} else if (millisecond >= 1) {
		return { value: Math.round(millisecond * 100) / 100, unit: 'ms' };
	} else if (microsecond >= 1) {
		return { value: Math.round(microsecond * 100) / 100, unit: 'us' };
	} else {
		return { value: Math.round(nanosecond * 100) / 100, unit: 'ns' };
	}
}

export function formatSecond(second: number): { value: string; unit: string } {
	const minute = second / 60;
	const hour = minute / 60;
	const day = hour / 24;

	if (day >= 1) {
		return { value: `${Math.round(day * 100) / 100}`, unit: 'd' };
	} else if (hour >= 1) {
		return { value: `${Math.round(hour * 100) / 100}`, unit: 'h' };
	} else if (minute >= 1) {
		return { value: `${Math.round(minute * 100) / 100}`, unit: 'm' };
	} else {
		return { value: `${Math.round(second * 100) / 100}`, unit: 's' };
	}
}

/**
 * Returns a Tailwind CSS background color class based on the given value.
 *
 * @param value - The input value, either as a percentage (0~100) or a decimal (0~1).
 * @param isPercent - Whether the value is in percentage format (default: true).
 * @param highIsGood - Whether higher values should be green (true) or red (false) (default: true).
 * @returns A string representing the Tailwind CSS background color class:
 *   When highIsGood is true:
 *   - '*:bg-green-700' for value > 62%
 *   - '*:bg-yellow-500' for value > 38% and <= 62%
 *   - '*:bg-red-700' for value <= 38%
 *   When highIsGood is false:
 *   - '*:bg-red-700' for value > 62%
 *   - '*:bg-yellow-500' for value > 38% and <= 62%
 *   - '*:bg-green-700' for value <= 38%
 */
export function formatProgressColor(value: number, isPercent: boolean = true, highIsGood: boolean = true): string {
	const percent = isPercent ? value : value * 100;

	if (highIsGood) {
		if (percent > 62) {
			return '*:bg-green-700';
		} else if (percent > 38) {
			return '*:bg-yellow-500';
		} else {
			return '*:bg-red-700';
		}
	} else {
		if (percent > 62) {
			return '*:bg-red-700';
		} else if (percent > 38) {
			return '*:bg-yellow-500';
		} else {
			return '*:bg-green-700';
		}
	}
}
