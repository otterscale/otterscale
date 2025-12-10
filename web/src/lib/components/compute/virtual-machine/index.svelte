<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import { PrometheusDriver } from 'prometheus-query';
	import { getContext, onDestroy, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import { env } from '$env/dynamic/public';
	import { EnvironmentService } from '$lib/api/environment/v1/environment_pb';
	import { InstanceService, type VirtualMachine } from '$lib/api/instance/v1/instance_pb';
	import * as Loading from '$lib/components/custom/loading';
	import { ReloadManager } from '$lib/components/custom/reloader';

	import { DataTable } from './data-table/index';
	import ExtensionsAlert from './extensions-alert.svelte';
	import type { Metrics } from './types';
	import { getMapNameToMetric } from './utils.svelte';
</script>

<script lang="ts">
	let { scope }: { scope: string } = $props();

	const transport: Transport = getContext('transport');
	const VirtualMachineClient = createClient(InstanceService, transport);
	const environmentService = createClient(EnvironmentService, transport);

	const virtualMachines = writable<VirtualMachine[]>([]);
	async function fetchVirtualMachines() {
		const response = await VirtualMachineClient.listVirtualMachines({
			scope: scope
		});
		virtualMachines.set(response.virtualMachines);
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

		const cpuResponse = await prometheusDriver.rangeQuery(
			`rate(kubevirt_vmi_cpu_usage_seconds_total{juju_model="${scope}"}[5m])`,
			Date.now() - 10 * 60 * 1000,
			Date.now(),
			2 * 60
		);
		const cpuSampleVectors = getMapNameToMetric(cpuResponse.result);
		const memoryResponse = await prometheusDriver.rangeQuery(
			`kubevirt_vmi_memory_resident_bytes{juju_model="${scope}"}`,
			Date.now() - 10 * 60 * 1000,
			Date.now(),
			2 * 60
		);
		const memorySampleVectors = getMapNameToMetric(memoryResponse.result);

		const storageReadResponse = await prometheusDriver.rangeQuery(
			`rate(kubevirt_vmi_storage_read_traffic_bytes_total{juju_model="${scope}"}[5m])`,
			Date.now() - 10 * 60 * 1000,
			Date.now(),
			2 * 60
		);
		const storageReadSampleVectors = getMapNameToMetric(storageReadResponse.result);

		const storageWriteResponse = await prometheusDriver.rangeQuery(
			`rate(kubevirt_vmi_storage_write_traffic_bytes_total{juju_model="${scope}"}[5m])`,
			Date.now() - 10 * 60 * 1000,
			Date.now(),
			2 * 60
		);
		const storageWriteSampleVectors = getMapNameToMetric(storageWriteResponse.result);

		metrics = {
			cpu: cpuSampleVectors,
			memory: memorySampleVectors,
			storageRead: storageReadSampleVectors,
			storageWrite: storageWriteSampleVectors
		};
	}

	async function fetch() {
		try {
			await Promise.all([fetchVirtualMachines(), fetchMetrics()]);
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
	<ExtensionsAlert {scope} />
	{#if isMounted}
		<DataTable {virtualMachines} {metrics} {scope} {reloadManager} />
	{:else}
		<Loading.DataTable />
	{/if}
</main>
