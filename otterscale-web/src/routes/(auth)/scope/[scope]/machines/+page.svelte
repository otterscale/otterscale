<script lang="ts">
	import Icon from '@iconify/svelte';
	import { page } from '$app/state';
	import * as Alert from '$lib/components/ui/alert';
	import { m } from '$lib/paraglide/messages';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { dynamicPaths } from '$lib/path';
	import { activeScope, breadcrumb } from '$lib/stores';
	import { MachineService, type Machine } from '$lib/api/machine/v1/machine_pb';
	import { DashBoard } from '$lib/components/machines';
	import { PrometheusDriver } from 'prometheus-query';
	import { EnvironmentService } from '$lib/api/environment/v1/environment_pb';

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
				endpoint: response.endpoint,
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

{#if mounted && prometheusDriver}
	<DashBoard client={prometheusDriver} />
{:else}
	loading
{/if}
<!-- 
{#if $activeScope}
	current scope: {$activeScope.uuid}
{/if}

{#if mounted}
	{#each $machinesStore as machine}
		<p>{machine.id}</p>
	{/each}
{/if} -->

<!-- <Alert.Root variant="default">
	<Icon icon="ph:airplane-takeoff" />
	<Alert.Title>{m.migrating()}</Alert.Title>
	<Alert.Description>{m.migrating_description()}</Alert.Description>
</Alert.Root> -->

<!-- <div class="pointer-events-none fixed inset-0 flex flex-col items-center justify-center">
	<Icon icon="ph:barricade" class="text-9xl" />
	{m.current_version({ version: import.meta.env.PACKAGE_VERSION })}
</div> -->
