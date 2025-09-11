<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';

	import { DataTable } from './snapshot/data-table';

	import type { VirtualMachine, VirtualMachineSnapshot } from '$lib/api/kubevirt/v1/kubevirt_pb';
	import { KubeVirtService } from '$lib/api/kubevirt/v1/kubevirt_pb';
	import * as Loading from '$lib/components/custom/loading';
	import * as Sheet from '$lib/components/ui/sheet';
	import { currentKubernetes } from '$lib/stores';
</script>

<script lang="ts">
	let { virtualMachine }: { virtualMachine: VirtualMachine } = $props();

	const transport: Transport = getContext('transport');
	const KubeVirtClient = createClient(KubeVirtService, transport);
	let snapshots = $state<VirtualMachineSnapshot[]>([]);

	let isSnapshotLoading = $state(true);

	async function fetchSnapshot() {
		try {
			const response = await KubeVirtClient.listVirtualMachineSnapshots({
				scopeUuid: $currentKubernetes?.scopeUuid,
				facilityName: $currentKubernetes?.name,
				vmName: virtualMachine.metadata?.name,
				namespace: virtualMachine.metadata?.namespace,
			});
			snapshots = response.snapshots;
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			isSnapshotLoading = false;
		}
	}

	let isMounted = $state(false);
	onMount(async () => {
		try {
			await fetchSnapshot();
			if (!isSnapshotLoading) {
				isMounted = true;
			}
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
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
				<DataTable virtualMachineSnapshots={snapshots} />
			{/if}
		</Sheet.Content>
	</Sheet.Root>
</div>
