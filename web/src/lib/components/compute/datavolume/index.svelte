<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	import { Data } from './data';

	// import { ApplicationService, type Application } from '$lib/api/application/v1/application_pb';
	import { KubeVirtService, type VirtualMachine } from '$lib/api/kubevirt/v1/kubevirt_pb';
	import * as Loading from '$lib/components/custom/loading';
</script>

<script lang="ts">
	let {
		scopeUuid,
		facilityName,
		namespace,
		virtualMachineName,
	}: { scopeUuid: string; facilityName: string; namespace: string; virtualMachineName: string } = $props();

	const transport: Transport = getContext('transport');
	const client = createClient(KubeVirtService, transport);

	const virtualMachine = writable<VirtualMachine>();

	let isMounted = $state(false);
	onMount(async () => {
		try {
			// getVirtualMachine, getDataVolume
			client
				.getVirtualMachine({
					scopeUuid: scopeUuid,
					facilityName: facilityName,
					namespace: namespace,
					name: virtualMachineName,
				})
				.then((response) => {
					virtualMachine.set(response);
					isMounted = true;
				});
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

<main>
	{#if isMounted}
		<Data {virtualMachine} />
	{:else}
		<Loading.Data />
	{/if}
</main>
