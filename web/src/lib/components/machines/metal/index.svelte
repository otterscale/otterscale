<script lang="ts" module>
	import { page } from '$app/state';
	import { MachineService, type Machine } from '$lib/api/machine/v1/machine_pb';
	import * as Loading from '$lib/components/custom/loading';
	import { Data } from '$lib/components/machines/metal/data';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';
</script>

<script lang="ts">
	const transport: Transport = getContext('transport');

	const machine = writable<Machine>();
	let isMounted = $state(false);

	const machineClient = createClient(MachineService, transport);

	onMount(async () => {
		try {
			machineClient
				.getMachine({
					id: page.params.id,
				})
				.then((response) => {
					machine.set(response);
					isMounted = true;
				});
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

<main>
	{#if isMounted}
		<Data {machine} />
	{:else}
		<Loading.Data />
	{/if}
</main>
