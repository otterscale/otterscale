import type { RangeVector } from 'prometheus-query';
import { SvelteMap } from 'svelte/reactivity';

import type { Metric } from './types';

function getMapCephDaemonToMetric(rangeVectors: RangeVector[]): Metric {
	return new SvelteMap(
		rangeVectors.map((rangeVector) => [
			(rangeVector.metric.labels as { ceph_daemon: string }).ceph_daemon,
			rangeVector.values
		])
	);
}

export { getMapCephDaemonToMetric };
