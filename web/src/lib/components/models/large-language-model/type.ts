import type { Application } from '$lib/api/application/v1/application_pb';

interface Metrics {
	kv_cache: number;
	gpu_cache: number;
	requests: number;
	time_to_first_token: number;
}

interface LargeLangeageModel {
	name: string;
	application: Application;
	metrics: Metrics;
}

export type { LargeLangeageModel };
