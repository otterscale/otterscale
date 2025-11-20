export const isFlexibleBooleanTrue = (envVar: string | undefined): boolean => {
	return ['true', '1', 'yes', 'on'].includes((envVar || '').toLowerCase());
};
