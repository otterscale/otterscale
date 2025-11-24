import type { RangeVector } from 'prometheus-query';
import { SvelteMap } from 'svelte/reactivity';

import type { Metric } from './types';

function getMapNameToMetric(rangeVectors: RangeVector[]): Metric {
	return new SvelteMap(
		rangeVectors.map((rangeVector) => [
			(rangeVector.metric.labels as { name: string }).name,
			rangeVector.values
		])
	);
}

export { getMapNameToMetric };
