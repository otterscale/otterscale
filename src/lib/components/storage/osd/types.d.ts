import type { SampleValue } from 'prometheus-query';

type Metric = SvelteMap<string, SampleValue[]>;

type Metrics = {
	input: Metric;
	output: Metric;
	read: Metric;
	write: Metric;
};
