import type { RangeVector } from 'prometheus-query';
import { SvelteMap } from 'svelte/reactivity';

import type { Metric } from './types';

function getMapInstanceToMetric(rangeVectors: RangeVector[]): Metric {
	return new SvelteMap(
		rangeVectors.map((cpu) => [(cpu.metric.labels as { instance: string }).instance, cpu.values])
	);
}

export { getMapInstanceToMetric };
