<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onDestroy, onMount, setContext } from 'svelte';
	import { writable } from 'svelte/store';

	import { DataTable } from './data-table/index';
	import { Statistics } from './statistics';

	import { KubeVirtService, type VirtualMachine } from '$lib/api/kubevirt/v1/kubevirt_pb';
	import * as Loading from '$lib/components/custom/loading';
	import { ReloadManager } from '$lib/components/custom/reloader';
</script>

<script lang="ts">
	let { scopeUuid, facilityName, namespace }: { scopeUuid: string; facilityName: string; namespace: string } =
		$props();

	const transport: Transport = getContext('transport');
	let isMounted = $state(false);

	const virtualMachines = writable<VirtualMachine[]>([]);

	const KubeVirtClient = createClient(KubeVirtService, transport);
	const reloadManager = new ReloadManager(() => {
		KubeVirtClient.listVirtualMachines({
			scopeUuid: scopeUuid,
			facilityName: facilityName,
			namespace: namespace,
		}).then((response) => {
			virtualMachines.set(response.virtualMachines);
		});
	});
	setContext('reloadManager', reloadManager);

	onMount(() => {
		KubeVirtClient.listVirtualMachines({
			scopeUuid: scopeUuid,
			facilityName: facilityName,
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
