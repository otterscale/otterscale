<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import { PrometheusDriver, SampleValue } from 'prometheus-query';
	import { getContext, onDestroy, onMount, setContext } from 'svelte';
	import type { SvelteMap } from 'svelte/reactivity';
	import { writable } from 'svelte/store';

	import { DataTable } from './data-table/index';
	import { type LargeLanguageModel } from './type';
	import { getGatewayURL, getMetricsMap } from './utils.svelte';

	import { env } from '$env/dynamic/public';
	import { ApplicationService, type Application } from '$lib/api/application/v1/application_pb';
	import { EnvironmentService } from '$lib/api/environment/v1/environment_pb';
	import * as Loading from '$lib/components/custom/loading';
	import { ReloadManager } from '$lib/components/custom/reloader';
</script>

<script lang="ts">
	let { scopeUuid, facilityName }: { scopeUuid: string; facilityName: string } = $props();

	const transport: Transport = getContext('transport');
	const applicationClient = createClient(ApplicationService, transport);
	const environmentService = createClient(EnvironmentService, transport);
	let prometheusDriver = $state<PrometheusDriver | null>(null);

	const largeLanguageModels = writable<LargeLanguageModel[]>([]);

	const applications = writable<Application[]>([]);
	const modelNames = writable<string[]>([]);
	const endToEndRequestLatencyMap = writable({} as SvelteMap<string | undefined, SampleValue[]>);
	const gpuCacheMap = writable({} as SvelteMap<string | undefined, SampleValue[]>);
	const kvCacheMap = writable({} as SvelteMap<string | undefined, SampleValue[]>);
	const timeToFirstTokenMap = writable({} as SvelteMap<string | undefined, SampleValue[]>);

	const modelServices = $derived(
		$applications.filter((application) =>
			application.labels['helm.sh/chart']
				? application.labels['helm.sh/chart'].startsWith('llm-d-modelservice')
				: false,
		),
	);

	async function fetchModels() {
		await applicationClient
			.listApplications({
				scope: scopeUuid,
				facility: facilityName,
			})
			.then((response) => {
				applications.set(response.applications);
			});

		await fetch(getGatewayURL($applications), {
			method: 'GET',
			headers: { 'Content-Type': 'application/json' },
		}).then((response) => {
			if (!response.ok) {
				throw new Error('Network response was not ok');
			}

			response
				.json()
				.then((r) => {
					modelNames.set((r.data as { id: string }[]).map((model) => model['id']));
				})
				.catch((error) => {
					console.error(`Failed to fetch models:`, error);
				});
		});

		await environmentService
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

		largeLanguageModels.set(
			modelServices.map(
				(model) =>
					({
						name: model.labels['model-name'],
						application: model,
						metrics: {
							gpu_cache: $gpuCacheMap.get('inference-dev'),
							kv_cache: $kvCacheMap.get('inference-dev'),
							requests: $endToEndRequestLatencyMap.get('inference-dev'),
							time_to_first_token: $timeToFirstTokenMap.get('inference-dev'),
						},
					}) as LargeLanguageModel,
			),
		);
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
	{#if isMounted}
		<DataTable {largeLanguageModels} {reloadManager} />
	{:else}
		<Loading.DataTable />
	{/if}
</main>
