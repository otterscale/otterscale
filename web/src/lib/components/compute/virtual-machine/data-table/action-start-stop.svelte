<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import { InstanceService, type VirtualMachine } from '$lib/api/instance/v1/instance_pb';
	import { m } from '$lib/paraglide/messages';
	import { currentKubernetes } from '$lib/stores';
</script>

<script lang="ts">
	let { virtualMachine }: { virtualMachine: VirtualMachine } = $props();

	const transport: Transport = getContext('transport');
	const virtualMachineClient = createClient(InstanceService, transport);
	let loading = $state(false);
	let statusAtClick = $state('Unknown');
	const isRunning = $derived(virtualMachine.status === 'Running');
	const isStarting = $derived(virtualMachine.status === 'Starting');

	$effect(() => {
		// If we are in a loading state and the status has changed from what it was when we clicked,
		// we can stop showing the loading indicator.
		if (loading && virtualMachine.status !== statusAtClick) {
			loading = false;
		}
	});

	async function startVM() {
		const request = {
			scope: $currentKubernetes?.scope,
			facility: $currentKubernetes?.name,
			name: virtualMachine.name,
			namespace: virtualMachine.namespace
		};

		toast.promise(() => virtualMachineClient.startVirtualMachine(request), {
			loading: `Starting virtual machine ${request.name}...`,
			success: () => `Successfully started virtual machine ${request.name}.`,
			error: (e) => {
				const msg = `Failed to start virtual machine ${request.name}.`;
				toast.error(msg, {
					description: (e as ConnectError).message.toString(),
					duration: Number.POSITIVE_INFINITY
				});
				return msg;
			}
		});
	}
	async function stopVM() {
		const request = {
			scope: $currentKubernetes?.scope,
			facility: $currentKubernetes?.name,
			name: virtualMachine.name,
			namespace: virtualMachine.namespace
		};

		toast.promise(() => virtualMachineClient.stopVirtualMachine(request), {
			loading: `Stopping virtual machine ${request.name}...`,
			success: () => `Successfully stopped virtual machine ${request.name}.`,
			error: (e) => {
				const msg = `Failed to stop virtual machine ${request.name}.`;
				toast.error(msg, {
					description: (e as ConnectError).message.toString(),
					duration: Number.POSITIVE_INFINITY
				});
				return msg;
			}
		});
	}

	async function handleClick() {
		// Store the status at the moment of the click
		statusAtClick = virtualMachine.status;
		loading = true;

		if (isRunning) {
			await stopVM();
		} else {
			await startVM();
		}
	}
</script>

<div class="flex items-center justify-end gap-1">
	<button
		onclick={handleClick}
		disabled={loading || isStarting}
		class="flex items-center gap-1 disabled:pointer-events-none disabled:cursor-not-allowed disabled:opacity-50"
	>
		{#if loading || isStarting}
			<Icon icon="ph:spinner-gap" class="animate-spin" />
			{m.please_wait()}
		{:else if isRunning}
			<Icon icon="ph:power" class="text-destructive" />
			{m.vm_stop()}
		{:else}
			<Icon icon="ph:power" class="text-accent-foreground" />
			{m.vm_start()}
		{/if}
	</button>
</div>
