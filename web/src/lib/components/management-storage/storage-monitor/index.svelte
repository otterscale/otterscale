<script lang="ts">
	import type { MON } from '$gen/api/storage/v1/storage_pb';
	import { StorageService } from '$gen/api/storage/v1/storage_pb';
	import PageLoading from '$lib/components/otterscale/ui/page-loading.svelte';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import { DataTable } from './data-table';

	const transport: Transport = getContext('transport');
	const storageClient = createClient(StorageService, transport);

	let monitors = $state(writable<MON[]>([]));
	let isMonitorsLoading = $state(true);
	async function fetchMonitors() {
		try {
			const response = await storageClient.listMONs({
				scopeUuid: 'b62d195e-3905-4960-85ee-7673f71eb21e',
				facilityName: 'ceph-mon'
			});
			monitors.set(response.mons);
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			isMonitorsLoading = false;
		}
	}

	let isMounted = $state(false);
	onMount(async () => {
		try {
			await fetchMonitors();
			if (!isMonitorsLoading) {
				isMounted = true;
			}
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

{#if isMounted}
	<DataTable bind:data={monitors} />
{:else}
	<PageLoading />
{/if}
