function getSortingFunction(
	previous: any,
	next: any,
	isLess: (previous: any, next: any) => boolean,
	isEqual: (previous: any, next: any) => boolean
) {
	if (!(previous || next)) {
		return 0;
	} else if (!previous) {
		return -1;
	} else if (!next) {
		return 1;
	}

	if (isEqual(previous, next)) {
		return 0;
	}

	if (isLess(previous, next)) {
		return -1;
	} else {
		return 1;
	}
}

export { getSortingFunction };
