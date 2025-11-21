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

/**
 * Returns a Tailwind CSS background color class based on the given value.
 *
 * @param target - LTB (Lower The Better) or STB (Stronger The Better)
 * @returns A string representing the Tailwind CSS background color class.
 */
const processVariant = {
	error: '*:bg-red-700 dark:*:bg-red-800 bg-red-100 dark:bg-red-900',
	warning: '*:bg-yellow-500 dark:*:bg-yellow-600 bg-yellow-100 dark:bg-yellow-900',
	information: '*:bg-green-700 dark:*:bg-green-800 bg-green-50 dark:bg-green-950'
};

type ProgressTargetType = 'LTB' | 'STB';

function getProgressColor(
	numerator: number,
	denominator: number,
	target: ProgressTargetType
): string {
	const percent = (numerator / denominator) * 100;

	if (target === 'LTB') {
		if (percent > 62) {
			return processVariant.information;
		} else if (percent > 38) {
			return processVariant.warning;
		} else {
			return processVariant.error;
		}
	} else if (target === 'STB') {
		if (percent < 38) {
			return processVariant.information;
		} else if (percent < 62) {
			return processVariant.warning;
		} else {
			return processVariant.error;
		}
	} else {
		return '';
	}
}

export type { ProgressTargetType };
export { formatRatio, getProgressColor };
