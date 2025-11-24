import type { SampleValue } from 'prometheus-query';
import type { SvelteMap } from 'svelte/reactivity';

type Meta = {
	isRowAction?: boolean;
};

type Metric = SvelteMap<string, SampleValue[]>;

type Metrics = {
	gpuCache: Metric;
	kvCache: Metric;
	requestLatency: Metric;
	timeToFirstToken: Metric;
};

export type { Meta, Metric, Metrics };
