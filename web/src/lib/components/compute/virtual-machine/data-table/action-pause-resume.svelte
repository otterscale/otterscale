<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import type { VirtualMachine } from '$lib/api/kubevirt/v1/kubevirt_pb';
	import { KubeVirtService, VirtualMachine_status } from '$lib/api/kubevirt/v1/kubevirt_pb';
	import { m } from '$lib/paraglide/messages';
	import { currentKubernetes } from '$lib/stores';
</script>

<script lang="ts">
	let { virtualMachine }: { virtualMachine: VirtualMachine } = $props();

	const transport: Transport = getContext('transport');
	const KubeVirtClient = createClient(KubeVirtService, transport);
	let loading = $state(false);
	let statusAtClick = $state<VirtualMachine_status>(VirtualMachine_status.UNKNOWN);
	const isRunning = $derived(virtualMachine.statusPhase === VirtualMachine_status.RUNNING);
	const isPaused = $derived(virtualMachine.statusPhase === VirtualMachine_status.PAUSED);
	const isShutdown = $derived(virtualMachine.statusPhase === VirtualMachine_status.STOPPED);

	$effect(() => {
		// If we are in a loading state and the status has changed from what it was when we clicked,
		// we can stop showing the loading indicator.
		if (loading && virtualMachine.statusPhase !== statusAtClick) {
			loading = false;
		}
	});

	async function resumeVM() {
		const request = {
			scopeUuid: $currentKubernetes?.scopeUuid,
			facilityName: $currentKubernetes?.name,
			name: virtualMachine.metadata?.name,
			namespace: virtualMachine.metadata?.namespace,
		};

		await toast.promise(() => KubeVirtClient.resumeVirtualMachine(request), {
			loading: `Resuming virtual machine ${request.name}...`,
			success: () => `Successfully resumed virtual machine ${request.name}.`,
			error: (e) => {
				const msg = `Failed to resume virtual machine ${request.name}.`;
				toast.error(msg, {
					description: (e as ConnectError).message.toString(),
					duration: Number.POSITIVE_INFINITY,
				});
				return msg;
			},
		});
	}
	async function pauseVM() {
		const request = {
			scopeUuid: $currentKubernetes?.scopeUuid,
			facilityName: $currentKubernetes?.name,
			name: virtualMachine.metadata?.name,
			namespace: virtualMachine.metadata?.namespace,
		};

		await toast.promise(() => KubeVirtClient.pauseVirtualMachine(request), {
			loading: `Pausing virtual machine ${request.name}...`,
			success: () => `Successfully paused virtual machine ${request.name}.`,
			error: (e) => {
				const msg = `Failed to pause virtual machine ${request.name}.`;
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
		statusAtClick = virtualMachine.statusPhase;
		loading = true;

		if (isRunning) {
			await pauseVM();
		} else if (isPaused) {
			await resumeVM();
		}
	}
</script>

<div class="flex items-center justify-end gap-1">
	<button
		onclick={handleClick}
		disabled={loading || isShutdown || (!isRunning && !isPaused)}
		class="flex items-center gap-1 disabled:pointer-events-none disabled:cursor-not-allowed disabled:opacity-50"
	>
		{#if loading}
			<Icon icon="ph:spinner-gap" class="animate-spin" />
			{m.please_wait()}
		{:else if isRunning}
			<Icon icon="ph:pause" class="text-orange-400" /> {m.vm_pause()}
		{:else if isPaused}
			<Icon icon="ph:play" class="text-accent-foreground" /> {m.vm_resume()}
		{:else if isShutdown}
			<Icon icon="ph:play" class="text-accent-foreground" /> {m.vm_resume()}
		{:else}
			<Icon icon="ph:pause" class="text-orange-400" /> {m.vm_pause()}
		{/if}
	</button>
</div>
