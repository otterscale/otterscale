<script lang="ts">
	import { page } from '$app/state';
	import { MachineService, type Machine } from '$lib/api/machine/v1/machine_pb';
	import { Data } from '$lib/components/machines/metal/data';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	const transport: Transport = getContext('transport');
	const client = createClient(MachineService, transport);

	const machine = writable<Machine>();

	let isMounted = $state(false);
	onMount(async () => {
		try {
			client
				.getMachine({
					id: page.params.id
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
		Loading
	{/if}
</main>
