import type { SampleValue } from 'prometheus-query';

import type { Application } from '$lib/api/application/v1/application_pb';

interface Metrics {
	kv_cache: SampleValue[];
	gpu_cache: SampleValue[];
	requests: SampleValue[];
	time_to_first_token: SampleValue[];
}

interface LargeLanguageModel {
	name: string;
	application: Application;
	metrics: Metrics;
}

export type { LargeLanguageModel };
