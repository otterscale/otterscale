import type { RangeVector } from 'prometheus-query';
import { SvelteMap } from 'svelte/reactivity';

import type { Metric } from './types.d.ts';

function getMapInstanceToMetric(rangeVectors: RangeVector[]): Metric {
	return new SvelteMap(
		rangeVectors.map((cpu) => [(cpu.metric.labels as { pod: string }).pod, cpu.values])
	);
}

export { getMapInstanceToMetric };
