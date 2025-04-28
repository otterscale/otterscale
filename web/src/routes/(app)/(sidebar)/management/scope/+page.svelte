<script lang="ts">
	import Icon from '@iconify/svelte';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { ManagementFacilityController } from '$lib/components/otterscale';
	import { Nexus, type Scope } from '$gen/api/nexus/v1/nexus_pb';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';

	const transport: Transport = getContext('transportNEW');
	const client = createClient(Nexus, transport);

	const scopesStore = writable<Scope[]>([]);
	const scopesLoading = writable(true);
	async function fetchScopes() {
		try {
			const response = await client.listScopes({});
			scopesStore.set(response.scopes);
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			scopesLoading.set(false);
		}
	}

	let mounted = false;
	onMount(async () => {
		try {
			await fetchScopes();
		} catch (error) {
			console.error('Error during initial data load:', error);
		}

		mounted = true;
	});
</script>

<main class="h-[calc(100vh_-_theme(spacing.16))]">
	{#if mounted}
		<ManagementFacilityController scopes={$scopesStore} />
	{:else}
		<div class="flex h-full w-full items-center justify-center gap-2 text-sm text-muted-foreground">
			<Icon icon="ph:spinner" class="size-8 animate-spin" />
			Loading...
		</div>
	{/if}
</main>
