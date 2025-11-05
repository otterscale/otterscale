import type { SampleValue } from 'prometheus-query';

import type { Application_Pod } from '$lib/api/application/v1/application_pb';
import type { Model } from '$lib/api/model/v1/model_pb';

type Meta = {
	isRowAction?: boolean;
};

interface Metrics {
	kv_cache: SampleValue[];
	gpu_cache: SampleValue[];
	requests: SampleValue[];
	time_to_first_token: SampleValue[];
}

interface Pod extends Application_Pod {
	metrics: Metrics;
}

interface LargeLanguageModel extends Model {
	pods: Pod[];
}

export type { LargeLanguageModel, Meta };
