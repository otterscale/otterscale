<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import type { VirtualMachine } from '$lib/api/kubevirt/v1/kubevirt_pb';
	import { KubeVirtService } from '$lib/api/kubevirt/v1/kubevirt_pb';
	import { m } from '$lib/paraglide/messages';
	import { currentKubernetes } from '$lib/stores';
</script>

<script lang="ts">
	let { virtualMachine }: { virtualMachine: VirtualMachine } = $props();

	const transport: Transport = getContext('transport');
	const KubeVirtClient = createClient(KubeVirtService, transport);
	let loading = $state(false);
	let statusAtClick = $state('');
	const isRunning = $derived(virtualMachine.statusPhase === 'Running');

	$effect(() => {
		// If we are in a loading state and the status has changed from what it was when we clicked,
		// we can stop showing the loading indicator.
		if (loading && virtualMachine.statusPhase !== statusAtClick) {
			loading = false;
		}
	});

	async function startVM() {
		const request = {
			scopeUuid: $currentKubernetes?.scopeUuid,
			facilityName: $currentKubernetes?.name,
			name: virtualMachine.metadata?.name,
			namespace: virtualMachine.metadata?.namespace,
		};

		await toast.promise(() => KubeVirtClient.startVirtualMachine(request), {
			loading: `Starting virtual machine ${request.name}...`,
			success: () => `Successfully started virtual machine ${request.name}.`,
			error: (e) => {
				const msg = `Failed to start virtual machine ${request.name}.`;
				toast.error(msg, {
					description: (e as ConnectError).message.toString(),
					duration: Number.POSITIVE_INFINITY,
				});
				return msg;
			},
		});
	}
	async function stopVM() {
		const request = {
			scopeUuid: $currentKubernetes?.scopeUuid,
			facilityName: $currentKubernetes?.name,
			name: virtualMachine.metadata?.name,
			namespace: virtualMachine.metadata?.namespace,
		};

		await toast.promise(() => KubeVirtClient.stopVirtualMachine(request), {
			loading: `Stopping virtual machine ${request.name}...`,
			success: () => `Successfully stopped virtual machine ${request.name}.`,
			error: (e) => {
				const msg = `Failed to stop virtual machine ${request.name}.`;
				toast.error(msg, {
					description: (e as ConnectError).message.toString(),
					duration: Number.POSITIVE_INFINITY,
				});
				return msg;
			},
		});
	}

	async function handleClick() {
		// Store the status at the moment of the click
		statusAtClick = virtualMachine.statusPhase ?? '';
		loading = true;

		try {
			if (isRunning) {
				await stopVM();
			} else {
				await startVM();
			}
		} finally {
			loading = false;
		}
	}
</script>

<div class="flex items-center justify-end gap-1">
	<button onclick={handleClick} disabled={loading} class="flex items-center gap-1">
		{#if loading}
			<Icon icon="ph:spinner-gap" class="animate-spin" />
			{m.please_wait()}
		{:else if isRunning}
			<!-- <Icon icon="ph:stop" /> {m.vm_stop()} -->
			<Icon icon="ph:power" class="text-destructive" />
			{m.vm_stop()}
		{:else}
			<!-- <Icon icon="ph:play" /> {m.vm_start()} -->
			<Icon icon="ph:power" class="text-accent-foreground" />
			{m.vm_start()}
		{/if}
	</button>
</div>
