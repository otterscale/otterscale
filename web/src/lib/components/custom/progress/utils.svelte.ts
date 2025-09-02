function formatRatio(numerator: number, denominator: number) {
	const value = (numerator * 100) / denominator;
	if (numerator / denominator === 1 || numerator / denominator === 0) return `${value.toFixed(0)}%`;
	else {
		return `${value.toFixed(2)}%`;
	}
}

export { formatRatio };
