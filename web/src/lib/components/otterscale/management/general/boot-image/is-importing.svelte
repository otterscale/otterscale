<script lang="ts">
	import { getContext, onMount } from 'svelte';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { Nexus } from '$gen/api/nexus/v1/nexus_pb';

	let {
		isImportingBootImages = $bindable()
	}: {
		isImportingBootImages: boolean;
	} = $props();

	const transport: Transport = getContext('transport');
	const client = createClient(Nexus, transport);

	async function fetchIsImportingBootImages(isImportingBootImages: boolean) {
		while (true) {
			const response = await client.isImportingBootImages({});
			console.log(`Checking importing status: ${response.importing}`);
			if (!response.importing) {
				isImportingBootImages = false;
				console.log(`Finished`);
				break;
			} else {
				await new Promise((resolve) => setTimeout(resolve, 5000)); // Wait 5 seconds between checks
			}
		}
	}

	onMount(async () => {
		try {
			console.log(`Checking importing status: ${isImportingBootImages}`);
			await new Promise((resolve) => setTimeout(resolve, 3000)); // Wait 3 seconds
			await fetchIsImportingBootImages(isImportingBootImages);
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

{#if isImportingBootImages}
	<span class="flex items-center gap-2 text-sm text-muted-foreground">
		<Icon icon="ph:spinner" class="size-5 animate-spin" />
		Importing...
	</span>
{:else}
	Finished
{/if}
