import type { SampleValue } from 'prometheus-query';

type Metric = SvelteMap<string, SampleValue[]>;

type Metrics = {
	cpu: Metric;
	memory: Metric;
	storage: Metric;
};
