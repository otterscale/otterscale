<script lang="ts" module>
	import { StorageService } from '$gen/api/storage/v1/storage_pb';
	import { DataTable as DataTableLoading } from '$lib/components/custom/loading';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext } from 'svelte';
	import { DataTable } from './data-table';
	import Picker from './pickers/index.svelte';
</script>

<script lang="ts">
	const transport: Transport = getContext('transport');
	const storageClient = createClient(StorageService, transport);

	let selectedScope = $state('b62d195e-3905-4960-85ee-7673f71eb21e');
	let selectedFacility = $state('ceph-mon');
	let selectedVolume = $state('');
</script>

<main class="space-y-4">
	<Picker bind:selectedScope bind:selectedFacility bind:selectedVolume />

	{#await storageClient.listSubvolumeGroups( { scopeUuid: selectedScope, facilityName: selectedFacility, volumeName: selectedVolume } )}
		<DataTableLoading />
	{:then response}
		{@const subvolumeGroups = response.subvolumeGroups}
		<DataTable {selectedScope} {selectedFacility} {selectedVolume} {subvolumeGroups} />
	{:catch}
		<DataTableLoading />
	{/await}
</main>
