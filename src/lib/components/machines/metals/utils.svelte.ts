import type { RangeVector } from 'prometheus-query';
import { SvelteMap } from 'svelte/reactivity';

import type { Metric } from './types';

function getMapInstanceToMetric(rangeVectors: RangeVector[]): Metric {
	return new SvelteMap(
		rangeVectors.map((rangeVector) => [
			(rangeVector.metric.labels as { instance: string }).instance,
			rangeVector.values
		])
	);
}

export { getMapInstanceToMetric };
