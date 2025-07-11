<script lang="ts" module>
	import { StorageService } from '$gen/api/storage/v1/storage_pb';
	import { DataTable as DataTableLoading } from '$lib/components/custom/loading';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext } from 'svelte';
	import { DataTable } from './data-table';
	import Pickers from './pickers.svelte';
</script>

<script lang="ts">
	const transport: Transport = getContext('transport');
	const storageClient = createClient(StorageService, transport);

	let selectedScope = $state('b62d195e-3905-4960-85ee-7673f71eb21e');
	let selectedFacility = $state('ceph-mon');
</script>

<main class="space-y-4">
	<Pickers bind:selectedScope bind:selectedFacility />

	{#await storageClient.listBuckets({ scopeUuid: selectedScope, facilityName: selectedFacility })}
		<DataTableLoading />
	{:then response}
		{@const buckets = response.buckets}
		<DataTable {selectedScope} {selectedFacility} {buckets} />
	{:catch}
		<DataTableLoading />
	{/await}
</main>
