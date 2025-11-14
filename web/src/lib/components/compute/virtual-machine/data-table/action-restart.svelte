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
	const isRestarting = $derived(virtualMachine.status === 'Restarting');

	$effect(() => {
		// If we are in a loading state and the status has changed from what it was when we clicked,
		// we can stop showing the loading indicator.
		if (loading && virtualMachine.status !== statusAtClick) {
			loading = false;
		}
	});

	async function restartVM() {
		const request = {
			scope: scope,
			,
			name: virtualMachine.name,
			namespace: virtualMachine.namespace
		};

		toast.promise(() => virtualMachineClient.restartVirtualMachine(request), {
			loading: `Restarting virtual machine ${request.name}...`,
			success: () => `Successfully restarted virtual machine ${request.name}.`,
			error: (e) => {
				const msg = `Failed to restart virtual machine ${request.name}.`;
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

		await restartVM();
	}
</script>

<div class="flex items-center justify-end gap-1">
	<button
		onclick={handleClick}
		disabled={loading || !isRunning}
		class="flex items-center gap-1 disabled:pointer-events-none disabled:cursor-not-allowed disabled:opacity-50"
	>
		{#if loading || isRestarting}
			<Icon icon="ph:spinner-gap" class="animate-spin" />
			{m.please_wait()}
		{:else}
			<Icon icon="ph:arrow-clockwise" class="text-accent-foreground" />
			{m.restarts()}
		{/if}
	</button>
</div>
