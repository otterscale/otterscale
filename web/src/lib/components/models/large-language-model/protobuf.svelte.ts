export interface LargeLangeageModel {
	name: string;
	cache: {
		gpu: number;
		kv: number;
	};
	usageStats: {
		requests: number;
		uptime: number;
	};
}
