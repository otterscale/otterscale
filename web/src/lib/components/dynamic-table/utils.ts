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

function getColumnType(type: JsonValue, format: JsonValue) {
	if (type === 'boolean') {
		return 'boolean';
	} else if (type === 'number' || type === 'integer') {
		return 'number';
	} else if (type === 'string' && (format === 'date' || format === 'date-time')) {
		return 'time';
	} else if (type === 'string') {
		return 'string';
	} else if (type === 'array') {
		return 'array';
	} else if (type === 'object' && format === 'ratio') {
		return 'ratio';
	} else if (type === 'object') {
		return 'object';
	} else {
		return undefined;
	}
}

export { format, getColumnType, getRelativeTime };
