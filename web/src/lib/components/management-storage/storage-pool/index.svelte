<script lang="ts">
	import { DataTable } from './data-table';
	import type { Pool } from '$gen/api/storage/v1/storage_pb';
	import { writable, type Writable } from 'svelte/store';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onMount } from 'svelte';
	import { StorageService } from '$gen/api/storage/v1/storage_pb';
	import PageLoading from '$lib/components/otterscale/ui/page-loading.svelte';

	const transport: Transport = getContext('transport');
	const storageClient = createClient(StorageService, transport);

	let pools = $state(writable<Pool[]>([]));
	let isPoolsLoading = $state(true);
	async function fetchPools() {
		try {
			const response = await storageClient.listPools({
				scopeUuid: 'b62d195e-3905-4960-85ee-7673f71eb21e',
				facilityName: 'ceph-mon'
			});
			pools.set(response.pools);
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			isPoolsLoading = false;
		}
	}

	let isMounted = $state(false);
	onMount(async () => {
		try {
			await fetchPools();
			if (!isPoolsLoading) {
				isMounted = true;
			}
		} catch (error) {
			console.error('Error during initial data load:', error);
		}
	});
</script>

{#if isMounted}
	<DataTable bind:data={pools} />
{:else}
	<PageLoading />
{/if}
