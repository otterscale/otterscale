<script lang="ts">
	import Icon from '@iconify/svelte';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { ManagementGeneral } from '$lib/components/otterscale';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import {
		Nexus,
		type Configuration,
		type Configuration_BootImageSelection,
		type Tag
	} from '$gen/api/nexus/v1/nexus_pb';

	import { PageLoading } from '$lib/components/otterscale/ui/index';

	const transport: Transport = getContext('transport');
	const client = createClient(Nexus, transport);

	const configurationStore = writable<Configuration>();
	const configurationLoading = writable(true);
	async function fetchConfiguration() {
		try {
			const response = await client.getConfiguration({});
			configurationStore.set(response);
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			configurationLoading.set(false);
		}
	}

	const tagsStore = writable<Tag[]>();
	const tagsLoading = writable(true);
	async function fetchTags() {
		try {
			const response = await client.listTags({});
			tagsStore.set(response.tags);
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			tagsLoading.set(false);
		}
	}

	// const bootImageSelectionsStore = writable<Configuration_BootImageSelection[]>([]);
	// const bootImageSelectionsLoading = writable(true);
	// async function fetchBootImageSelections() {
	// 	try {
	// 		const response = await client.listBootImageSelections({});
	// 		bootImageSelectionsStore.set(response.bootImageSelections);
	// 	} catch (error) {
	// 		console.error('Error fetching:', error);
	// 	} finally {
	// 		bootImageSelectionsLoading.set(false);
	// 	}
	// }

	let mounted = false;
	onMount(async () => {
		try {
			await fetchConfiguration();
			await fetchTags();
		} catch (error) {
			console.error('Error during initial data load:', error);
		}

		mounted = true;
	});
</script>

{#if mounted}
	<ManagementGeneral configuration={$configurationStore} tags={$tagsStore} />
{:else}
	<PageLoading />
{/if}
