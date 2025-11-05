<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import { PrometheusDriver, SampleValue } from 'prometheus-query';
	import { getContext, onDestroy, onMount, setContext } from 'svelte';
	import type { SvelteMap } from 'svelte/reactivity';
	import { writable } from 'svelte/store';

	import { DataTable } from './data-table/index';
	import ExtensionsAlert from './extensions-alert.svelte';
	import { type LargeLanguageModel } from './type';
	import { getMetricsMap } from './utils.svelte';

	import { env } from '$env/dynamic/public';
	import { EnvironmentService } from '$lib/api/environment/v1/environment_pb';
	import { ModelService, type Model } from '$lib/api/model/v1/model_pb';
	import * as Loading from '$lib/components/custom/loading';
	import { ReloadManager } from '$lib/components/custom/reloader';
</script>

<script lang="ts">
	let { scope, facility }: { scope: string; facility: string } = $props();

	const transport: Transport = getContext('transport');
	const modelClient = createClient(ModelService, transport);
	const environmentClient = createClient(EnvironmentService, transport);
	let prometheusDriver = $state<PrometheusDriver | null>(null);

	const largeLanguageModels = writable<LargeLanguageModel[]>([]);

	const models = writable<Model[]>([]);

	const endToEndRequestLatencyMap = writable({} as SvelteMap<string | undefined, SampleValue[]>);
	const gpuCacheMap = writable({} as SvelteMap<string | undefined, SampleValue[]>);
	const kvCacheMap = writable({} as SvelteMap<string | undefined, SampleValue[]>);
	const timeToFirstTokenMap = writable({} as SvelteMap<string | undefined, SampleValue[]>);

	async function fetchModels() {
		modelClient
			.listModels({ scope, facility, namespace: 'default' })
			.then((response) => {
				models.set(response.models);
			})
			.catch((error) => (console.error(`Failed to fetch models for namespace default:`, error), []));

		await environmentClient
			.getPrometheus({})
			.then((response) => {
				prometheusDriver = new PrometheusDriver({
					endpoint: `${env.PUBLIC_API_URL}/prometheus`,
					baseURL: response.baseUrl,
				});
			})
			.catch((error) => {
				console.error('Failed to initialize Prometheus driver:', error);
			});

		if (prometheusDriver) {
			await prometheusDriver
				.rangeQuery(`vllm:e2e_request_latency_seconds_sum`, Date.now() - 10 * 60 * 1000, Date.now(), 2 * 60)
				.then((response) => {
					endToEndRequestLatencyMap.set(getMetricsMap(response.result));
				});
			await prometheusDriver
				.rangeQuery(`vllm:gpu_cache_usage_perc`, Date.now() - 10 * 60 * 1000, Date.now(), 2 * 60)
				.then((response) => {
					gpuCacheMap.set(getMetricsMap(response.result));
				});
			await prometheusDriver
				.rangeQuery(`vllm:kv_cache_usage_perc`, Date.now() - 10 * 60 * 1000, Date.now(), 2 * 60)
				.then((response) => {
					kvCacheMap.set(getMetricsMap(response.result));
				});
			await prometheusDriver
				.rangeQuery(`vllm:time_to_first_token_seconds_sum`, Date.now() - 10 * 60 * 1000, Date.now(), 2 * 60)
				.then((response) => {
					timeToFirstTokenMap.set(getMetricsMap(response.result));
				});
		}

		largeLanguageModels.set([
			...($models.map(
				(model) =>
					({
						...model,
						pods: model.pods.map((pod) => ({
							...pod,
							metrics: {
								gpu_cache: $gpuCacheMap.get('vllm-llama-3-2-1b-instruct-deployment-85c77654c7-st9qk'),
								kv_cache: $kvCacheMap.get('vllm-llama-3-2-1b-instruct-deployment-85c77654c7-st9qk'),
								requests: $endToEndRequestLatencyMap.get(
									'vllm-llama-3-2-1b-instruct-deployment-85c77654c7-st9qk',
								),
								time_to_first_token: $timeToFirstTokenMap.get(
									'vllm-llama-3-2-1b-instruct-deployment-85c77654c7-st9qk',
								),
							},
						})),
					}) as LargeLanguageModel,
			) as LargeLanguageModel[]),
			...($models.map(
				(model) =>
					({
						...model,
						pods: [
							...model.pods.map((pod) => ({
								...pod,
								metrics: {
									gpu_cache: $gpuCacheMap.get(
										'vllm-llama-3-2-1b-instruct-deployment-85c77654c7-st9qk',
									),
									kv_cache: $kvCacheMap.get('vllm-llama-3-2-1b-instruct-deployment-85c77654c7-st9qk'),
									requests: $endToEndRequestLatencyMap.get(
										'vllm-llama-3-2-1b-instruct-deployment-85c77654c7-st9qk',
									),
									time_to_first_token: $timeToFirstTokenMap.get(
										'vllm-llama-3-2-1b-instruct-deployment-85c77654c7-st9qk',
									),
								},
							})),
							...model.pods.map((pod) => ({
								...pod,
								metrics: {
									gpu_cache: $gpuCacheMap.get(
										'vllm-llama-3-2-1b-instruct-deployment-85c77654c7-st9qk',
									),
									kv_cache: $kvCacheMap.get('vllm-llama-3-2-1b-instruct-deployment-85c77654c7-st9qk'),
									requests: $endToEndRequestLatencyMap.get(
										'vllm-llama-3-2-1b-instruct-deployment-85c77654c7-st9qk',
									),
									time_to_first_token: $timeToFirstTokenMap.get(
										'vllm-llama-3-2-1b-instruct-deployment-85c77654c7-st9qk',
									),
								},
							})),
						],
					}) as LargeLanguageModel,
			) as LargeLanguageModel[]),
		]);
	}

	const reloadManager = new ReloadManager(() => {
		fetchModels();
	});
	setContext('reloadManager', reloadManager);

	let isMounted = $state(false);
	onMount(async () => {
		await fetchModels();
		isMounted = true;
		reloadManager.start();
	});
	onDestroy(() => {
		reloadManager.stop();
	});
</script>

<main class="space-y-4 py-4">
	<ExtensionsAlert {scope} {facility} />
	{#if isMounted}
		<DataTable {largeLanguageModels} {reloadManager} />
	{:else}
		<Loading.DataTable />
	{/if}
</main>
