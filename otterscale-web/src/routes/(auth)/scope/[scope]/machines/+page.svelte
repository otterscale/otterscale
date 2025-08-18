<script lang="ts">
	import { page } from '$app/state';
	import { env } from '$env/dynamic/public';
	import { EnvironmentService } from '$lib/api/environment/v1/environment_pb';
	import { MachineService, type Machine } from '$lib/api/machine/v1/machine_pb';
	import { Dashboard } from '$lib/components/machines/dashboard';
	import { dynamicPaths } from '$lib/path';
	import { breadcrumb } from '$lib/stores';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { PrometheusDriver } from 'prometheus-query';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	// Set breadcrumb navigation
	breadcrumb.set({ parents: [], current: dynamicPaths.machines(page.params.scope) });

	const transport: Transport = getContext('transport');
	const machineClient = createClient(MachineService, transport);
	const environmentService = createClient(EnvironmentService, transport);
	const machinesStore = writable<Machine[]>([]);
	let prometheusDriver: PrometheusDriver | null = null;

	async function initializePrometheusDriver(): Promise<PrometheusDriver> {
		try {
			const response = await environmentService.getPrometheus({});
			return new PrometheusDriver({
				endpoint: `${env.PUBLIC_API_URL}/prometheus`,
				baseURL: response.baseUrl
			});
		} catch (error) {
			console.error('Error initializing Prometheus driver:', error);
			throw error;
		}
	}

	async function fetchMachines(): Promise<void> {
		try {
			const response = await machineClient.listMachines({});
			machinesStore.set(response.machines);
		} catch (error) {
			console.error('Error fetching machines:', error);
		}
	}

	let mounted = false;

	onMount(async () => {
		try {
			prometheusDriver = await initializePrometheusDriver();
			await fetchMachines();
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
		mounted = true;
	});
</script>

{#if mounted && prometheusDriver && $machinesStore.length > 0}
	{@const filteredMachines = $machinesStore.filter((machine) =>
		machine.workloadAnnotations?.['juju-machine-id']?.includes('-machine-')
	)}
	{@const allMachine = {
		fqdn: filteredMachines.map((machine) => machine.fqdn).join('|'),
		id: 'All Machine'
	} as Machine}
	{@const machines = [allMachine, ...filteredMachines]}
	<Dashboard client={prometheusDriver} {machines} />
{:else}
	<div class="flex items-center justify-center p-8">
		<Icon icon="mdi:loading" class="animate-spin text-2xl" />
		<span class="ml-2">Loading machines...</span>
	</div>
{/if}
