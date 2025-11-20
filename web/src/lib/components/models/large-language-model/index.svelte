<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import { PrometheusDriver } from 'prometheus-query';
	import { getContext, onDestroy, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import { env } from '$env/dynamic/public';
	import { type Application, ApplicationService } from '$lib/api/application/v1/application_pb';
	import { EnvironmentService } from '$lib/api/environment/v1/environment_pb';
	import * as Loading from '$lib/components/custom/loading';
	import { ReloadManager } from '$lib/components/custom/reloader';

	import { DataTable } from './data-table/index';
	import ExtensionsAlert from './extensions-alert.svelte';
	import { type LargeLanguageModel } from './type';
</script>

<script lang="ts">
	let { scope }: { scope: string } = $props();

	const largeLanguageModels = writable<LargeLanguageModel[]>([]);

	const transport: Transport = getContext('transport');
	const applicationClient = createClient(ApplicationService, transport);
	const environmentService = createClient(EnvironmentService, transport);

	let prometheusDriver = $state<PrometheusDriver | null>(null);

	const applications = writable<Application[]>([]);
	const models = $derived($applications.filter((application) => application.labels['model-name']));

	async function fetchApplications() {
		const response = await applicationClient.listApplications({
			scope: scope
		});
		applications.set(response.applications);
	}

	async function fetchPrometheus() {
		try {
			const response = await environmentService.getPrometheus({});
			prometheusDriver = new PrometheusDriver({
				endpoint: `${env.PUBLIC_API_URL}/prometheus`,
				baseURL: response.baseUrl
			});
		} catch (error) {
			console.error('Failed to initialize Prometheus driver:', error);
		}
	}

	async function fetchMetrics() {
		if (!prometheusDriver) return {};

		const [gpuCacheResponse, kvCacheResponse, timeToFirstTokenResponse, requestLatencyResponse] =
			await Promise.all([
				prometheusDriver.instantQuery(`vllm:gpu_cache_usage_perc{scope_uuid="${scope}"}`),
				prometheusDriver.instantQuery(`vllm:kv_cache_usage_perc{scope_uuid="${scope}"}`),
				prometheusDriver.instantQuery(
					`vllm:time_to_first_token_seconds_sum{scope_uuid="${scope}"}`
				),
				prometheusDriver.instantQuery(`vllm:e2e_request_latency_seconds_sum{scope_uuid="${scope}"}`)
			]);

		const gpuCacheUsageByPod = new Map(
			gpuCacheResponse.result.map((instantVector) => [
				(instantVector.metric.labels as { pod?: string }).pod ?? '',
				instantVector.value.value
			])
		);

		const kvCacheUsageByPod = new Map(
			kvCacheResponse.result.map((instantVector) => [
				(instantVector.metric.labels as { pod?: string }).pod ?? '',
				instantVector.value.value
			])
		);

		const timeToFirstTokenByPod = new Map(
			timeToFirstTokenResponse.result.map((instantVector) => [
				(instantVector.metric.labels as { pod?: string }).pod ?? '',
				instantVector.value.value
			])
		);

		const requestLatencyByPod = new Map(
			requestLatencyResponse.result.map((instantVector) => [
				(instantVector.metric.labels as { pod?: string }).pod ?? '',
				instantVector.value.value
			])
		);

		return {
			gpuCacheUsageByPod,
			kvCacheUsageByPod,
			timeToFirstTokenByPod,
			requestLatencyByPod
		};
	}

	async function fetch() {
		try {
			await fetchApplications();
			await fetchPrometheus();

			const {
				gpuCacheUsageByPod = new Map(),
				kvCacheUsageByPod = new Map(),
				timeToFirstTokenByPod = new Map(),
				requestLatencyByPod = new Map()
			} = await fetchMetrics();

			largeLanguageModels.set(
				models.map(
					(model) =>
						({
							name: model.labels['model-name'],
							application: model,
							metrics: {
								gpu_cache: gpuCacheUsageByPod.get(model.labels['model-name']) ?? 0,
								kv_cache: kvCacheUsageByPod.get(model.labels['model-name']) ?? 0,
								requests: requestLatencyByPod.get(model.labels['model-name']) ?? 0,
								time_to_first_token: timeToFirstTokenByPod.get(model.labels['model-name']) ?? 0
							}
						}) as LargeLanguageModel
				)
			);
		} catch (error) {
			console.error('Failed to fetch large language model data:', error);
		}
	}

	const reloadManager = new ReloadManager(fetch);

	let isLoaded = $state(false);
	onMount(async () => {
		await fetch();
		isLoaded = true;
	});
	onDestroy(() => {
		reloadManager.stop();
	});
</script>

<main class="space-y-4 py-4">
	<ExtensionsAlert {scope} />
	{#if isLoaded}
		<DataTable {largeLanguageModels} {scope} {reloadManager} />
	{:else}
		<Loading.DataTable />
	{/if}
</main>
