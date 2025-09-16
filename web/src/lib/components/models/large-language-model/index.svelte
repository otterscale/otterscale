<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import { InstantVector, PrometheusDriver } from 'prometheus-query';
	import { getContext, onDestroy, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import { DataTable } from './data-table/index';
	import { type LargeLangeageModel } from './protobuf.svelte';

	import { ApplicationService } from '$lib/api/application/v1/application_pb';
	import * as Loading from '$lib/components/custom/loading';
	import { ReloadManager } from '$lib/components/custom/reloader';
</script>

<script lang="ts">
	let { scopeUuid, facilityName }: { scopeUuid: string; facilityName: string } = $props();

	const largeLanguageModels = writable<LargeLangeageModel[]>([]);

	const transport: Transport = getContext('transport');
	const applicationClient = createClient(ApplicationService, transport);

	const prometheusDriver = new PrometheusDriver({
		endpoint: 'http://192.168.41.100:30091',
		baseURL: '/api/v1',
	});

	let gpuCacheUsage = $state([] as InstantVector[]);
	let kvCacheUsage = $state([] as InstantVector[]);
	let timeToFirstTokenSeconds = $state([] as InstantVector[]);
	let requestLatencySeconds = $state([] as InstantVector[]);

	let modelNames = [] as string[];

	async function fetch() {
		let gpuCacheUsageByPod = $state({} as Map<string, number>);
		let kvCacheUsageByPod = $state({} as Map<string, number>);
		let timeToFirstTokenByPod = $state({} as Map<string, number>);
		let requestLatencyByPod = $state({} as Map<string, number>);

		await applicationClient
			.listApplications({
				scopeUuid: scopeUuid,
				facilityName: facilityName,
			})
			.then((response) => {
				modelNames = [
					...response.applications
						.filter((application) => application.labels['model-name'])
						.map((application) => application.labels['model-name']),
				];
			})
			.catch((error) => {
				console.error('Error during initial data load:', error);
			});
		await prometheusDriver.instantQuery(`vllm:gpu_cache_usage_perc{scope_uuid="${scopeUuid}"}`).then((response) => {
			gpuCacheUsage = response.result;
			gpuCacheUsageByPod = new Map(
				gpuCacheUsage.map((instantVector) => [
					(instantVector.metric.labels as { pod?: string }).pod ?? '',
					instantVector.value.value,
				]),
			);
			modelNames = [...modelNames, ...Array.from(gpuCacheUsageByPod.keys())];
		});
		await prometheusDriver.instantQuery(`vllm:kv_cache_usage_perc{scope_uuid="${scopeUuid}"}`).then((response) => {
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

		await largeLanguageModels.set(
			modelNames.map(
				(modelName) =>
					({
						name: modelName,
						cache: {
							gpu: gpuCacheUsageByPod.get(modelName) ?? 0,
							kv: kvCacheUsageByPod.get(modelName) ?? 0,
						},
						usageStats: {
							requests: requestLatencyByPod.get(modelName) ?? 0,
							uptime: timeToFirstTokenByPod.get(modelName) ?? 0,
						},
					}) as LargeLangeageModel,
			),
		);
	}

	const reloadManager = new ReloadManager(fetch);
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
