import { SvelteSet } from 'svelte/reactivity';

const CATEGORY_CLEAN = new Set(['active', 'clean']);
const CATEGORY_WORKING = new Set([
	'activating',
	'backfill_wait',
	'backfilling',
	'creating',
	'deep',
	'degraded',
	'forced_backfill',
	'forced_recovery',
	'peering',
	'peered',
	'recovering',
	'recovery_wait',
	'repair',
	'scrubbing',
	'snaptrim',
	'snaptrim_wait'
]);
const CATEGORY_WARNING = new Set([
	'backfill_toofull',
	'backfill_unfound',
	'down',
	'incomplete',
	'inconsistent',
	'recovery_toofull',
	'recovery_unfound',
	'remapped',
	'snaptrim_error',
	'stale',
	'undersized'
]);
const VALID = CATEGORY_CLEAN.union(CATEGORY_WORKING).union(CATEGORY_WARNING);

function parsePlacementGroupStates(placementGroupState: string): SvelteSet<string> {
	return new SvelteSet(
		placementGroupState
			.replace(/[^a-z_]+/g, ' ')
			.trim()
			.split(' ')
	);
}
function getPlacementGroupStateVariant(placementGroupState: string) {
	const states = parsePlacementGroupStates(placementGroupState);
	if (!states.isSubsetOf(VALID)) {
		return 'secondary';
	} else if (states.intersection(CATEGORY_WARNING).size > 0) {
		return 'outline';
	} else if (states.intersection(CATEGORY_WORKING).size > 0) {
		return 'default';
	} else if (states.intersection(CATEGORY_CLEAN).size > 0) {
		return 'outline';
	}
	return 'secondary';
}

function getPlacementGroupStateClassName(placementGroupState: string) {
	const states = parsePlacementGroupStates(placementGroupState);
	if (states.intersection(CATEGORY_WARNING).size > 0) {
		return 'text-destructive border-destructive/30';
	}
	return '';
}

export { getPlacementGroupStateClassName, getPlacementGroupStateVariant };
