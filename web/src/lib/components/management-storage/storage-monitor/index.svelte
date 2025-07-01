<script lang="ts">
	import type { OSD } from '$gen/api/storage/v1/storage_pb';
	import { StorageService } from '$gen/api/storage/v1/storage_pb';
	import PageLoading from '$lib/components/otterscale/ui/page-loading.svelte';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import { DataTable } from './data-table';

	const transport: Transport = getContext('transport');
	const storageClient = createClient(StorageService, transport);

	let osds = $state(writable<OSD[]>([]));
	let isOSDsLoading = $state(true);
	async function fetchOSDs() {
		try {
			const response = await storageClient.listOSDs({
				scopeUuid: 'b62d195e-3905-4960-85ee-7673f71eb21e',
				facilityName: 'ceph-mon'
			});
			osds.set(response.osds);
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			isOSDsLoading = false;
		}
	}

	let isMounted = $state(false);
	onMount(async () => {
		try {
			await fetchOSDs();
			if (!isOSDsLoading) {
				isMounted = true;
			}
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

{#if isMounted}
	<DataTable bind:data={osds} />
{:else}
	<PageLoading />
{/if}
