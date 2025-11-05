import type { RangeVector, SampleValue } from 'prometheus-query';
import { SvelteMap } from 'svelte/reactivity';

function getMetricsMap(vectors: RangeVector[]) {
	return new SvelteMap(
		vectors.map((vector) => [(vector.metric.labels as { pod?: string }).pod, vector.values as SampleValue[]]),
	);
}

export { getMetricsMap };
