<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import { PrometheusDriver } from 'prometheus-query';
	import { getContext, onDestroy, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import { EnvironmentService } from '$lib/api/environment/v1/environment_pb';
	import { type Machine, MachineService } from '$lib/api/machine/v1/machine_pb';
	import * as Loading from '$lib/components/custom/loading';
	import { ReloadManager } from '$lib/components/custom/reloader';

	import { DataTable } from './data-table/index';
	import type { Metrics } from './types';
	import { getMapInstanceToMetric } from './utils.svelte';
</script>

<script lang="ts">
	const transport: Transport = getContext('transport');
	const machineClient = createClient(MachineService, transport);
	const environmentService = createClient(EnvironmentService, transport);

	const machines = writable<Machine[]>([]);
	async function fetchMachines() {
		const response = await machineClient.listMachines({});
		machines.set(response.machines);
	}

	let prometheusDriver = $state<PrometheusDriver | null>(null);
	let metrics = $state({} as Metrics);
	async function fetchMetrics() {
		if (!prometheusDriver) {
			const response = await environmentService.getPrometheus({});
			prometheusDriver = new PrometheusDriver({
				endpoint: '/prometheus',
				baseURL: response.baseUrl,
				headers: {
					'x-proxy-target': 'api'
				}
			});
		}

		const cpuResponse = await prometheusDriver.rangeQuery(
			`
			avg by (instance) (
				1 - rate(node_cpu_seconds_total{juju_model=~".*",mode="idle"}[2m])
			)
			`,
			Date.now() - 10 * 60 * 1000,
			Date.now(),
			2 * 60
		);
		const cpuSampleVectors = getMapInstanceToMetric(cpuResponse.result);

		const memoryResponse = await prometheusDriver.rangeQuery(
			`
				node_memory_MemTotal_bytes
				- node_memory_MemFree_bytes
				- node_memory_Cached_bytes
				- node_memory_Buffers_bytes
				- node_memory_SReclaimable_bytes
			`,
			Date.now() - 10 * 60 * 1000,
			Date.now(),
			2 * 60
		);
		const memorySampleVectors = getMapInstanceToMetric(memoryResponse.result);

		const storageResponse = await prometheusDriver.rangeQuery(
			`
				1 - sum by (instance) (node_filesystem_avail_bytes) / sum by (instance) (node_filesystem_size_bytes)
			`,
			Date.now() - 10 * 60 * 1000,
			Date.now(),
			2 * 60
		);
		const storageSampleVectors = getMapInstanceToMetric(storageResponse.result);

		metrics = { cpu: cpuSampleVectors, memory: memorySampleVectors, storage: storageSampleVectors };
	}

	async function fetch() {
		try {
			await Promise.all([fetchMachines(), fetchMetrics()]);
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
		<DataTable {machines} {metrics} {reloadManager} />
	{:else}
		<Loading.DataTable />
	{/if}
</main>
