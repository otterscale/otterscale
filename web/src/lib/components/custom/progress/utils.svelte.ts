function formatRatio(numerator: number, denominator: number) {
	const value = (numerator * 100) / denominator;
	if (numerator / denominator === 1 || numerator / denominator === 0) return `${value.toFixed(0)}%`;
	else {
		// Round up (ceiling) to two decimal places
		const roundedValue = Math.ceil(value * 100) / 100;
		return `${roundedValue.toFixed(2)}%`;
	}
}

export { formatRatio };
