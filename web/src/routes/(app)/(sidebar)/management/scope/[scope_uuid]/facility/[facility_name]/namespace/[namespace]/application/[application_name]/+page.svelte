<script lang="ts">
	import { page } from '$app/state';
	import { ManagementApplication } from '$lib/components/otterscale/index';
	import Icon from '@iconify/svelte';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import { Nexus, type Application } from '$gen/api/nexus/v1/nexus_pb';

	const transport: Transport = getContext('transportNEW');
	const client = createClient(Nexus, transport);

	const applicationStore = writable<Application>();
	const applicationIsLoading = writable(true);
	async function fetchApplications() {
		try {
			const response = await client.getApplication({
				scopeUuid: page.params.scope_uuid,
				facilityName: page.params.facility_name,
				namespace: page.params.namespace,
				name: page.params.application_name
			});
			applicationStore.set(response);
		} catch (error) {
			console.error('Error fetching machine:', error);
		} finally {
			applicationIsLoading.set(false);
		}
	}

	let mounted = false;
	onMount(async () => {
		try {
			await fetchApplications();
		} catch (error) {
			console.error('Error during initial data load:', error);
		}

		mounted = true;
	});
</script>

<main class="h-[calc(100vh_-_theme(spacing.16))]">
	{#if mounted}
		<ManagementApplication application={$applicationStore} />
	{:else}
		<div class="flex h-full w-full items-center justify-center gap-2 text-sm text-muted-foreground">
			<Icon icon="ph:spinner" class="size-8 animate-spin" />
			Loading...
		</div>
	{/if}
</main>
