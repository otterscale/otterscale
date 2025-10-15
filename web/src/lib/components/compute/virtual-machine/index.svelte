<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onDestroy, onMount, setContext } from 'svelte';
	import { writable } from 'svelte/store';

	import { DataTable } from './data-table/index';
	import { Statistics } from './statistics';

	import { InstanceService, type VirtualMachine } from '$lib/api/instance/v1/instance_pb';
	import * as Loading from '$lib/components/custom/loading';
	import { ReloadManager } from '$lib/components/custom/reloader';
</script>

<script lang="ts">
	let { scope, facility, namespace }: { scope: string; facility: string; namespace: string } = $props();

	const transport: Transport = getContext('transport');
	let isMounted = $state(false);

	const virtualMachines = writable<VirtualMachine[]>([]);

	const VirtualMachineClient = createClient(InstanceService, transport);
	const reloadManager = new ReloadManager(() => {
		VirtualMachineClient.listVirtualMachines({
			scope: scope,
			facility: facility,
			namespace: namespace,
		}).then((response) => {
			virtualMachines.set(response.virtualMachines);
		});
	});
	setContext('reloadManager', reloadManager);

	onMount(() => {
		VirtualMachineClient.listVirtualMachines({
			scope: scope,
			facility: facility,
			namespace: namespace,
		})
			.then((response) => {
				virtualMachines.set(response.virtualMachines);
				isMounted = true;
			})
			.catch((error) => {
				console.error('Error during initial data load:', error);
			});

		reloadManager.start();
	});
	onDestroy(() => {
		reloadManager.stop();
	});
</script>

<main class="space-y-4 py-4">
	{#if isMounted}
		<Statistics virtualMachines={$virtualMachines} />
		<DataTable {virtualMachines} {reloadManager} />
	{:else}
		<Loading.DataTable />
	{/if}
</main>
