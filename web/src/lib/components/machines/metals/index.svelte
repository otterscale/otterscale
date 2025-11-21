<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import { PrometheusDriver } from 'prometheus-query';
	import { getContext, onDestroy, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import { env } from '$env/dynamic/public';
	import { EnvironmentService } from '$lib/api/environment/v1/environment_pb';
	import { type Machine, MachineService } from '$lib/api/machine/v1/machine_pb';
	import * as Loading from '$lib/components/custom/loading';
	import { ReloadManager } from '$lib/components/custom/reloader';

	import { DataTable } from './data-table/index';
	import { Statistics } from './statistics';
	import type { Metrics } from './types';
	import { getMapInstanceToMetric } from './utils.svelte';
</script>

<script lang="ts">
	const transport: Transport = getContext('transport');
	const machineClient = createClient(MachineService, transport);
	const environmentService = createClient(EnvironmentService, transport);

	const machines = writable<Machine[]>([]);
	async function fetchMachines() {
		machineClient
			.listMachines({})
			.then((response) => {
				machines.set(response.machines);
				isMounted = true;
			})
			.catch((error) => {
				console.error('Error during initial data load:', error);
			});
		return $machines;
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

		const memoryResponse = await prometheusDriver.rangeQuery(
			`
				(
					node_memory_MemTotal_bytes - node_memory_MemFree_bytes
					-
					(node_memory_Cached_bytes + node_memory_Buffers_bytes + node_memory_SReclaimable_bytes)
				)
				/
				node_memory_MemTotal_bytes
			`,
			Date.now() - 10 * 60 * 1000,
			Date.now(),
			2 * 60
		);
		const memorySampleVectors = getMapInstanceToMetric(memoryResponse.result);

		const storageResponse = await prometheusDriver.rangeQuery(
			`
				node_load1
			`,
			Date.now() - 10 * 60 * 1000,
			Date.now(),
			2 * 60
		);
		const storageSampleVectors = getMapInstanceToMetric(storageResponse.result);

		metrics = { memory: memorySampleVectors, storage: storageSampleVectors };
	}

	async function fetch() {
		try {
			await Promise.all([fetchMachines(), fetchMetrics()]);
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
	{#if isMounted}
		<Statistics machines={$machines} />
		<DataTable {machines} {metrics} {reloadManager} />
	{:else}
		<Loading.DataTable />
	{/if}
</main>
