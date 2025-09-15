<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount, onDestroy, setContext } from 'svelte';
	import { writable } from 'svelte/store';

	import { DataTable } from './snapshot/data-table';

	import type { VirtualMachine, VirtualMachineSnapshot } from '$lib/api/kubevirt/v1/kubevirt_pb';
	import { KubeVirtService } from '$lib/api/kubevirt/v1/kubevirt_pb';
	import * as Loading from '$lib/components/custom/loading';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Sheet from '$lib/components/ui/sheet';
	import { currentKubernetes } from '$lib/stores';
</script>

<script lang="ts">
	let { virtualMachine }: { virtualMachine: VirtualMachine } = $props();

	let snapshots = writable<VirtualMachineSnapshot[]>([]);

	const transport: Transport = getContext('transport');
	const KubeVirtClient = createClient(KubeVirtService, transport);
	const reloadManager = new ReloadManager(() => {
		KubeVirtClient.listVirtualMachineSnapshots({
			scopeUuid: $currentKubernetes?.scopeUuid,
			facilityName: $currentKubernetes?.name,
			vmName: virtualMachine.metadata?.name,
			namespace: virtualMachine.metadata?.namespace,
		}).then((response) => {
			snapshots.set(response.snapshots);
		});
	});
	setContext('reloadManager', reloadManager);

	let isMounted = $state(false);
	onMount(() => {
		KubeVirtClient.listVirtualMachineSnapshots({
			scopeUuid: $currentKubernetes?.scopeUuid,
			facilityName: $currentKubernetes?.name,
			vmName: virtualMachine.metadata?.name,
			namespace: virtualMachine.metadata?.namespace,
		})
			.then((response) => {
				snapshots.set(response.snapshots);
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

<div class="flex items-center justify-end gap-1">
	<Sheet.Root>
		<Sheet.Trigger class="flex items-center gap-1">
			<Icon icon="mdi:backup-restore" /> Snapshot
		</Sheet.Trigger>

		<Sheet.Content class="min-w-[70vw] p-4">
			{#if !isMounted}
				<Loading.Report />
			{:else}
				<DataTable {virtualMachine} virtualMachineSnapshots={$snapshots} />
			{/if}
		</Sheet.Content>
	</Sheet.Root>
</div>
