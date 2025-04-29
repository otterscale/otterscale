<script lang="ts">
	import { getContext, onMount } from 'svelte';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { writable } from 'svelte/store';
	import { Nexus, type Machine } from '$gen/api/nexus/v1/nexus_pb';

	let {
		id
	}: {
		id: string;
	} = $props();

	let status = $state('');
	let statusMessgae = $state('');

	const transport: Transport = getContext('transportNEW');
	const client = createClient(Nexus, transport);

	const machineStore = writable<Machine>();
	async function fetchMachineStatus() {
		while (true) {
			console.log(`Checking machine status...`);

			try {
				const response = await client.getMachine({ id: id });
				machineStore.set(response);

				status = $machineStore.status;
				statusMessgae = $machineStore.statusMessage;
			} catch (error) {
				console.error('Error fetching machine:', error);
			}

			if ($machineStore.status.toLowerCase() == 'deployed') {
				console.log(`Finished`);
				break;
			} else {
				await new Promise((resolve) => setTimeout(resolve, 30000)); // Wait 5 seconds between checks
			}
		}
	}

	onMount(async () => {
		try {
			await fetchMachineStatus();
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

<main class="text-xs font-extralight">
	{#if status.toLowerCase() == 'deployed'}
		{statusMessgae}
	{:else}
		<span class="flex items-center gap-1">
			<Icon icon="ph:spinner" class="animate-spin" />
			{statusMessgae}
		</span>
	{/if}
</main>
