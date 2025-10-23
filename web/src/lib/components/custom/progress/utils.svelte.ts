function formatRatio(numerator: number, denominator: number) {
	const value = (numerator * 100) / denominator;
	if (numerator / denominator === 1 || numerator / denominator === 0) return `${value.toFixed(0)}%`;
	else {
		// When value is less than 0.01 and greater than 0, round up to 0.01
		if (value > 0 && value < 0.01) {
			return '0.01%';
		}
		// Otherwise, round to two decimal places
		return `${value.toFixed(2)}%`;
	}
}

export { formatRatio };
