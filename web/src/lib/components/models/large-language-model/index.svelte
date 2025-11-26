<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import { PrometheusDriver } from 'prometheus-query';
	import { getContext, onDestroy, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import { env } from '$env/dynamic/public';
	import { EnvironmentService } from '$lib/api/environment/v1/environment_pb';
	import { type Model, ModelService } from '$lib/api/model/v1/model_pb.ts';
	import * as Loading from '$lib/components/custom/loading';
	import { ReloadManager } from '$lib/components/custom/reloader';

	import { DataTable } from './data-table/index';
	import ExtensionsAlert from './extensions-alert.svelte';
	import type { Metrics } from './types.d.ts';
	import { getMapInstanceToMetric } from './utils.svelte.ts';
</script>

<script lang="ts">
	let { scope }: { scope: string } = $props();

	const transport: Transport = getContext('transport');
	const modelClient = createClient(ModelService, transport);
	const environmentService = createClient(EnvironmentService, transport);

	const models = writable<Model[]>([]);
	async function fetchModels() {
		const response = await modelClient.listModels({
			scope: scope
		});
		models.set(response.models);
	}

	let prometheusDriver = $state<PrometheusDriver | null>(null);
	let metrics = $state({} as Metrics);
	async function fetchMetrics() {
		if (!prometheusDriver) {
			const response = await environmentService.getPrometheus({});
			prometheusDriver = new PrometheusDriver({
				endpoint: `${env.PUBLIC_API_URL}/prometheus`,
				baseURL: response.baseUrl
			});
		}

		const kvCacheResponse = await prometheusDriver.rangeQuery(
			`vllm:kv_cache_usage_perc{juju_model="${scope}"}`,
			Date.now() - 10 * 60 * 1000,
			Date.now(),
			2 * 60
		);
		const kvCacheSampleVectors = getMapInstanceToMetric(kvCacheResponse.result);

		const timeToFirstTokenResponse = await prometheusDriver.rangeQuery(
			`vllm:time_to_first_token_seconds_sum{juju_model="${scope}"}`,
			Date.now() - 10 * 60 * 1000,
			Date.now(),
			2 * 60
		);
		const timeToFirstTokenSampleVectors = getMapInstanceToMetric(timeToFirstTokenResponse.result);

		const requestLatencyResponse = await prometheusDriver.rangeQuery(
			`vllm:e2e_request_latency_seconds_sum{juju_model="${scope}"}`,
			Date.now() - 10 * 60 * 1000,
			Date.now(),
			2 * 60
		);
		const requestLatencySampleVectors = getMapInstanceToMetric(requestLatencyResponse.result);

		metrics = {
			kvCache: kvCacheSampleVectors,
			requestLatency: requestLatencySampleVectors,
			timeToFirstToken: timeToFirstTokenSampleVectors
		};
	}

	async function fetch() {
		try {
			await Promise.all([fetchModels(), fetchMetrics()]);
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	}

	const reloadManager = new ReloadManager(fetch, false);

	let isMounted = $state(false);
	onMount(async () => {
		await fetch();
		isMounted = true;
	});
	onDestroy(() => {
		reloadManager.stop();
	});
</script>

<main class="space-y-4 py-4">
	<ExtensionsAlert {scope} />
	{#if isMounted}
		<DataTable {models} namespace="default" {metrics} {scope} {reloadManager} />
	{:else}
		<Loading.DataTable />
	{/if}
</main>
