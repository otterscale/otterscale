export interface LargeLangeageModel {
	name: string;
	gpu_cache: number;
	kv_cache: number;
	requests: number;
	time_to_first_token: number;
}
