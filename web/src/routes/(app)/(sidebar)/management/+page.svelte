<script lang="ts">
	import { createClient, type Transport } from '@connectrpc/connect';
	import { Nexus, type Error } from '$gen/api/nexus/v1/nexus_pb';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import { Monitor } from '$lib/components/otterscale/index';

	const transport: Transport = getContext('transportNEW');
	const client = createClient(Nexus, transport);

	const errorsStore = writable<Error[]>([]);
	const errorsLoading = writable(true);
	async function fetchScopes() {
		try {
			const response = await client.verifyEnvironment({});
			errorsStore.set(response.errors);
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			errorsLoading.set(false);
		}
	}
	// const machinesStore = writable<Machine[]>([]);
	// const machinesLoading = writable(true);
	// async function fetchMachines() {
	// 	try {
	// 		const response = await client.listMachines({});
	// 		machinesStore.set(response.machines);
	// 	} catch (error) {
	// 		console.error('Error fetching:', error);
	// 	} finally {
	// 		machinesLoading.set(false);
	// 	}
	// }

	let mounted = false;
	onMount(async () => {
		try {
			await fetchScopes();
			// await fetchMachines();
		} catch (error) {
			console.error('Error during initial data load:', error);
		}

		mounted = true;
	});
</script>

<Monitor errors={$errorsStore} />
