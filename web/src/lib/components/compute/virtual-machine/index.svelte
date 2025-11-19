<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onDestroy, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import { InstanceService, type VirtualMachine } from '$lib/api/instance/v1/instance_pb';
	import * as Loading from '$lib/components/custom/loading';
	import { ReloadManager } from '$lib/components/custom/reloader';

	import { DataTable } from './data-table/index';
	import ExtensionsAlert from './extensions-alert.svelte';
	import { Statistics } from './statistics';
</script>

<script lang="ts">
	let { scope }: { scope: string } = $props();

	const transport: Transport = getContext('transport');
	const VirtualMachineClient = createClient(InstanceService, transport);

	const virtualMachines = writable<VirtualMachine[]>([]);
	async function fetch() {
		VirtualMachineClient.listVirtualMachines({
			scope: scope
		})
			.then((response) => {
				virtualMachines.set(response.virtualMachines);
				isMounted = true;
			})
			.catch((error) => {
				console.error('Error during initial data load:', error);
			});
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
	<ExtensionsAlert {scope} />
	{#if isMounted}
		<Statistics virtualMachines={$virtualMachines} />
		<DataTable {virtualMachines} {scope} {reloadManager} />
	{:else}
		<Loading.DataTable />
	{/if}
</main>
