import type { SampleValue } from 'prometheus-query';

type Metric = SvelteMap<string, SampleValue[]>;

type Metrics = {
	memory: Metric;
	storage: Metric;
};
