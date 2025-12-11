<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import { PrometheusDriver } from 'prometheus-query';
	import { getContext, onDestroy, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import { env } from '$env/dynamic/public';
	import { EnvironmentService } from '$lib/api/environment/v1/environment_pb';
	import { type ObjectStorageDaemon, StorageService } from '$lib/api/storage/v1/storage_pb';
	import * as Loading from '$lib/components/custom/loading';
	import { ReloadManager } from '$lib/components/custom/reloader';

	import { DataTable } from './data-table';
	import type { Metrics } from './types';
	import { getMapCephDaemonToMetric } from './utils.svelte';
</script>

<script lang="ts">
	let {
		scope
	}: {
		scope: string;
	} = $props();

	const transport: Transport = getContext('transport');
	const storageClient = createClient(StorageService, transport);
	const environmentService = createClient(EnvironmentService, transport);

	const objectStorageDaemons = writable([] as ObjectStorageDaemon[]);
	async function fetchObjectStorageDaemons() {
		const response = await storageClient.listObjectStorageDaemons({ scope: scope });
		objectStorageDaemons.set(response.objectStorageDaemons);
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

		const readResponse = await prometheusDriver.rangeQuery(
			`irate(ceph_osd_op_r{juju_model="${scope}"}[1h])`,
			Date.now() - 10 * 60 * 1000,
			Date.now(),
			2 * 60
		);
		const readSampleVectors = getMapCephDaemonToMetric(readResponse.result);

		const writeResponse = await prometheusDriver.rangeQuery(
			`irate(ceph_osd_op_w{juju_model="${scope}"}[1h])`,
			Date.now() - 10 * 60 * 1000,
			Date.now(),
			2 * 60
		);
		const writeSampleVectors = getMapCephDaemonToMetric(writeResponse.result);

		const inputResponse = await prometheusDriver.rangeQuery(
			`irate(ceph_osd_op_w_in_bytes{juju_model="${scope}"}[1h])`,
			Date.now() - 10 * 60 * 1000,
			Date.now(),
			2 * 60
		);
		const inputSampleVectors = getMapCephDaemonToMetric(inputResponse.result);

		const outputResponse = await prometheusDriver.rangeQuery(
			`irate(ceph_osd_op_r_out_bytes{juju_model="${scope}"}[1h])`,
			Date.now() - 10 * 60 * 1000,
			Date.now(),
			2 * 60
		);
		const outputSampleVectors = getMapCephDaemonToMetric(outputResponse.result);

		metrics = {
			read: readSampleVectors,
			write: writeSampleVectors,
			input: inputSampleVectors,
			output: outputSampleVectors
		};
	}

	async function fetch() {
		try {
			await Promise.all([fetchObjectStorageDaemons(), fetchMetrics()]);
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
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
		<DataTable {objectStorageDaemons} {metrics} {scope} {reloadManager} />
	{:else}
		<Loading.DataTable />
	{/if}
</main>
