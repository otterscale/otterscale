<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import { InstantVector, PrometheusDriver } from 'prometheus-query';
	import { getContext, onDestroy, onMount, setContext } from 'svelte';
	import { writable } from 'svelte/store';

	import { DataTable } from './data-table/index';
	import { type LargeLanguageModel } from './type';

	import { env } from '$env/dynamic/public';
	import { ApplicationService, type Application } from '$lib/api/application/v1/application_pb';
	import { EnvironmentService } from '$lib/api/environment/v1/environment_pb';
	import * as Loading from '$lib/components/custom/loading';
	import { ReloadManager } from '$lib/components/custom/reloader';
</script>

<script lang="ts">
	let { scopeUuid, facilityName }: { scopeUuid: string; facilityName: string } = $props();

	const largeLanguageModels = writable<LargeLanguageModel[]>([]);

	const transport: Transport = getContext('transport');
	const applicationClient = createClient(ApplicationService, transport);
	const environmentService = createClient(EnvironmentService, transport);

	let prometheusDriver = $state<PrometheusDriver | null>(null);

	const applications = writable<Application[]>([]);
	const models = $derived($applications.filter((application) => application.labels['model-name']));

	async function fetch() {
		let gpuCacheUsage = $state([] as InstantVector[]);
		let gpuCacheUsageByPod = $state({} as Map<string, number>);

		let kvCacheUsage = $state([] as InstantVector[]);
		let kvCacheUsageByPod = $state({} as Map<string, number>);

		let timeToFirstTokenSeconds = $state([] as InstantVector[]);
		let timeToFirstTokenByPod = $state({} as Map<string, number>);

		let requestLatencySeconds = $state([] as InstantVector[]);
		let requestLatencyByPod = $state({} as Map<string, number>);

		await applicationClient
			.listApplications({
				scope: scopeUuid,
				facility: facilityName,
			})
			.then((response) => {
				applications.set(response.applications);
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
				.instantQuery(`vllm:gpu_cache_usage_perc{scope_uuid="${scopeUuid}"}`)
				.then((response) => {
					gpuCacheUsage = response.result;
					gpuCacheUsageByPod = new Map(
						gpuCacheUsage.map((instantVector) => [
							(instantVector.metric.labels as { pod?: string }).pod ?? '',
							instantVector.value.value,
						]),
					);
				});
			await prometheusDriver
				.instantQuery(`vllm:kv_cache_usage_perc{scope_uuid="${scopeUuid}"}`)
				.then((response) => {
					kvCacheUsage = response.result;
					kvCacheUsageByPod = new Map(
						kvCacheUsage.map((instantVector) => [
							(instantVector.metric.labels as { pod?: string }).pod ?? '',
							instantVector.value.value,
						]),
					);
				});
			await prometheusDriver
				.instantQuery(`vllm:time_to_first_token_seconds_sum{scope_uuid="${scopeUuid}"}`)
				.then((response) => {
					timeToFirstTokenSeconds = response.result;
					timeToFirstTokenByPod = new Map(
						timeToFirstTokenSeconds.map((instantVector) => [
							(instantVector.metric.labels as { pod?: string }).pod ?? '',
							instantVector.value.value,
						]),
					);
				});
			await prometheusDriver
				.instantQuery(`vllm:e2e_request_latency_seconds_sum{scope_uuid="${scopeUuid}"}`)
				.then((response) => {
					requestLatencySeconds = response.result;
					requestLatencyByPod = new Map(
						requestLatencySeconds.map((instantVector) => [
							(instantVector.metric.labels as { pod?: string }).pod ?? '',
							instantVector.value.value,
						]),
					);
				});
		}

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
							time_to_first_token: timeToFirstTokenByPod.get(model.labels['model-name']) ?? 0,
						},
					}) as LargeLanguageModel,
			),
		);
	}

	const reloadManager = new ReloadManager(() => {
		fetch();
	});
	setContext('reloadManager', reloadManager);
	let isMounted = $state(false);

	onMount(async () => {
		await fetch();
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
