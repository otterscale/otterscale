import type { SampleValue } from 'prometheus-query';
import type { SvelteMap } from 'svelte/reactivity';

type ModeSource = 'local' | 'cloud';

type Meta = {
	isRowAction?: boolean;
};

type Metric = SvelteMap<string, SampleValue[]>;

type Metrics = {
	requestLatency: Metric;
	timeToFirstToken: Metric;
};

export type { Meta, Metric, Metrics, ModeSource };
